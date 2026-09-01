// Package token emite y valida tokens JWT (HMAC-SHA256) con la stdlib.
// Se usa para la sesión del login: el token viaja como cookie httpOnly.
package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

var (
	// ErrInvalidToken indica que el token no es válido (firma, expiración o formato).
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired indica que el token expiró.
	ErrTokenExpired = errors.New("token expired")
)

const (
	RememberTTL = 10 * 365 * 24 * time.Hour // "Recordarme": 10 años
	SessionTTL  = 7 * 24 * time.Hour        // sin "Recordarme": 7 días
)

// Secret es la clave HMAC para firmar/validar tokens.
// En producción se debe configurar el env MYTICKTICK_TOKEN_SECRET; el valor
// por defecto es solo para desarrollo.
var Secret = secretFromEnv()

// secretFromEnv lee MYTICKTICK_TOKEN_SECRET del entorno y usa el valor de
// desarrollo si no está definido.
func secretFromEnv() []byte {
	if v := os.Getenv("MYTICKTICK_TOKEN_SECRET"); v != "" {
		return []byte(v)
	}
	return []byte("myticktick-dev-secret-change-in-prod")
}

// Claims son los claims del token.
type Claims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"sub"`
	IssuedAt int64  `json:"iat"`
	Expiry   int64  `json:"exp"`
}

// B64 encodea/decodifica base64url sin padding.
func b64encode(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }
func b64decode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// ExpirationFor devuelve la expiración según el flag "Recordarme".
func ExpirationFor(remember bool) time.Time {
	now := time.Now()
	if remember {
		return now.Add(RememberTTL)
	}
	return now.Add(SessionTTL)
}

// Generate emite un token JWT firmado para el usuario.
func Generate(userID uint, username string, expiry time.Time) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		IssuedAt: time.Now().Unix(),
		Expiry:   expiry.Unix(),
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	signingInput := b64encode(headerJSON) + "." + b64encode(claimsJSON)
	return signingInput + "." + sign(signingInput), nil
}

// sign calcula HMAC-SHA256 sobre la entrada y lo encodea base64url.
func sign(input string) string {
	mac := hmac.New(sha256.New, Secret)
	mac.Write([]byte(input))
	return b64encode(mac.Sum(nil))
}

// Validate verifica firma y expiración, y devuelve los claims.
func Validate(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	expected := sign(signingInput)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, ErrInvalidToken
	}

	claimsJSON, err := b64decode(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if claims.Expiry < time.Now().Unix() {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}
