package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"myticktick/internal/core/token"
)

// TestRequireAuth_RejectsSinCookie verifica que sin cookie la ruta devuelve 401.
func TestRequireAuth_RejectsSinCookie(t *testing.T) {
	handler := RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("sin cookie debe devolver 401, got %d", rec.Code)
	}
}

// TestRequireAuth_AcceptaCookieValida verifica que con cookie válida la ruta
// procese y el userID quede disponible en el context.
func TestRequireAuth_AcceptaCookieValida(t *testing.T) {
	var gotUserID uint
	sawUserID := false

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id, ok := RequestUserID(r.Context()); ok {
			gotUserID = id
			sawUserID = true
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := RequireAuth(inner)

	exp := time.Now().Add(time.Hour)
	tok, _ := token.Generate(7, "rodri", exp)

	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookie, Value: tok})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("con cookie válida debe devolver 200, got %d", rec.Code)
	}
	if !sawUserID || gotUserID != 7 {
		t.Errorf("el userID del context debe ser 7, got %d (ok=%v)", gotUserID, sawUserID)
	}
}

// TestRequireAuth_RejectsTokenInvalido verifica que un token basura dé 401.
func TestRequireAuth_RejectsTokenInvalido(t *testing.T) {
	handler := RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookie, Value: "no-es-un-jwt"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("token inválido debe devolver 401, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid") {
		t.Errorf("el body debe indicar invalid, got %q", rec.Body.String())
	}
}
