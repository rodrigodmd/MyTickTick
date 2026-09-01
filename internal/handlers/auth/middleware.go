package auth

import (
	"context"
	"net/http"

	"myticktick/internal/core/token"
)

// ContextKey es la clave para guardar el userID en el context de la request.
type contextKey string

const ContextUserID contextKey = "mtt_user_id"

// RequestUserID devuelve el userID de la request (inyectado por RequireAuth).
func RequestUserID(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(ContextUserID).(uint)
	return id, ok
}

// RequireAuth es el middleware de protección de rutas.
// Lee la cookie de sesión, valida el token JWT y, si es válido, inyecta
// el userID en el context de la request. Si no hay cookie o el token es
// inválido/expirado, devuelve 401 sin procesar el handler.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionCookie)
		if err != nil || cookie.Value == "" {
			writeUnauthorized(w, "missing session")
			return
		}

		claims, err := token.Validate(cookie.Value)
		if err != nil {
			writeUnauthorized(w, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeUnauthorized(w http.ResponseWriter, reason string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"unauthorized","reason":"` + reason + `"}`))
}
