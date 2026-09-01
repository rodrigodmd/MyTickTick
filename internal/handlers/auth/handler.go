package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"myticktick/internal/core/services/auth"
)

// AuthHandler maneja las peticiones de autenticación.
type AuthHandler struct {
	authService *auth.AuthService
}

// NewHandler crea una nueva instancia de AuthHandler.
func NewHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// secureDetects si la petición vino por HTTPS (para marcar la cookie Secure).
func secureDetects(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	// Detrás de proxy: X-Forwarded-Proto https
	return r.Header.Get("X-Forwarded-Proto") == "https"
}

// Register maneja POST /api/register. Crea el usuario y, si las credenciales
// son válidas, inicia la sesión emitiendo la cookie (UX: registro → adentro).
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RequestRegister

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, `{"error":"username and password are required"}`, http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	// Iniciar sesión de inmediato: emitir la cookie de sesión.
	session, err := h.authService.Login(r.Context(), req.Username, req.Password, false)
	if err != nil {
		// El usuario quedó creado pero no se pudo firmar la sesión.
		// Devolvemos el usuario para que el cliente pueda hacer login manual.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ResponseMessage{
			Message: "Usuario creado exitosamente",
		})
		return
	}

	secure := secureDetects(r)
	SetSessionCookie(w, session.Token, false, secure)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Usuario creado exitosamente",
		"user":    ToResponseUser(user),
	})
}

// Login maneja POST /api/login. Emite la cookie de sesión si las credenciales son válidas.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req RequestLogin

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	session, err := h.authService.Login(r.Context(), req.Username, req.Password, req.Remember)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) || errors.Is(err, auth.ErrUserNotFound) {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"error":"login failed"}`, http.StatusInternalServerError)
		return
	}

	secure := secureDetects(r)
	SetSessionCookie(w, session.Token, req.Remember, secure)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Login exitoso",
		"user":    ToResponseUser(session.User),
	})
}

// Logout maneja POST /api/logout. Borra la cookie de sesión.
// (El token JWT es stateless: borrar la cookie termina la sesión del cliente.
// Si se necesita revocación server-side, agregar un blacklist o session store.)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ClearSessionCookie(w, secureDetects(r))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseMessage{
		Message: "Sesión cerrada",
	})
}

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case IsUniqueViolation(err):
		http.Error(w, `{"error":"username already exists"}`, http.StatusConflict)
	default:
		http.Error(w, `{"error":"failed to register user"}`, http.StatusInternalServerError)
	}
}
