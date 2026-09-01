package auth

import (
	"net/http"
	"time"
)

// SessionCookie es el nombre de la cookie que guarda el token JWT.
const SessionCookie = "mtt_session"

// SetSessionCookie establece la cookie de sesión con el token JWT.
// remember=true -> Max-Age de 10 años; false -> 7 días.
// Siempre httpOnly (no accesible desde JS) y SameSite=Lax (protección CSRF básica).
func SetSessionCookie(w http.ResponseWriter, token string, remember bool, secure bool) {
	var maxAge int64
	if remember {
		maxAge = int64((10 * 365 * 24 * time.Hour) / time.Second)
	} else {
		maxAge = int64((7 * 24 * time.Hour) / time.Second)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   int(maxAge),
		HttpOnly: true,
		Secure:   secure, // true solo si la petición vino por HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie borra la cookie de sesión (logout).
func ClearSessionCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}
