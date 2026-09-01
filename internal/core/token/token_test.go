package token

import (
	"testing"
	"time"
)

// TestGenerateValidate round-trip: generar un token y validar sus claims.
func TestGenerateValidate(t *testing.T) {
	exp := time.Now().Add(1 * time.Hour)
	tok, err := Generate(42, "rodri", exp)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	claims, err := Validate(tok)
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}
	if claims.Username != "rodri" {
		t.Errorf("Username = %q, want %q", claims.Username, "rodri")
	}
}

// TestValidate_RejectsTampered verifica que un token manipulado sea rechazado.
func TestValidate_RejectsTampered(t *testing.T) {
	tok, _ := Generate(1, "user", time.Now().Add(time.Hour))

	// Corromper el payload
	parts := []byte(tok)
	parts[len(parts)-5] ^= 0xFF
	if _, err := Validate(string(parts)); err == nil {
		t.Error("un token corrompido debe ser rechazado")
	}
}

// TestValidate_RejectsExpired verifica que un token expirado sea rechazado.
func TestValidate_RejectsExpired(t *testing.T) {
	tok, _ := Generate(1, "user", time.Now().Add(-1*time.Hour)) // ya expiró
	if _, err := Validate(tok); err != ErrTokenExpired {
		t.Errorf("un token expirado debe devolver ErrTokenExpired, got %v", err)
	}
}

// TestValidate_RejectsForeignSecret verifica que un token firmado con otra
// clave no pase validación (protege contra falsificación con secreto distinto).
func TestValidate_RejectsForeignSecret(t *testing.T) {
	origSecret := Secret
	defer func() { Secret = origSecret }()

	// Firmar con un secreto...
	Secret = []byte("secreto-A")
	tokA, _ := Generate(1, "user", time.Now().Add(time.Hour))

	// ...y validar con otro secreto diferente.
	Secret = []byte("secreto-B")
	if _, err := Validate(tokA); err == nil {
		t.Error("un token firmado con otro secreto debe ser rechazado")
	}
}

// TestExpirationFor verifica las duraciones de "Recordarme" vs sesión normal.
func TestExpirationFor(t *testing.T) {
	remember := ExpirationFor(true)
	session := ExpirationFor(false)

	// "Recordarme" debe durar ~10 años
	if d := remember.Sub(time.Now()); d < 9*365*24*time.Hour || d > 11*365*24*time.Hour {
		t.Errorf("Recordarme debe durar ~10 años, got %v", d)
	}
	// Sesión normal debe durar ~7 días
	if d := session.Sub(time.Now()); d < 6*24*time.Hour || d > 8*24*time.Hour {
		t.Errorf("sesión normal debe durar ~7 días, got %v", d)
	}
	// Recordarme debe durar más que sesión normal
	if !remember.After(session) {
		t.Error("Recordarme debe expirar más tarde que la sesión normal")
	}
}
