package auth

import "net/http"

// MapURLs configura las rutas para este handler
func (h *AuthHandler) MapURLs(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/register", h.Register)
	mux.HandleFunc("POST /api/login", h.Login)
	mux.HandleFunc("POST /api/logout", h.Logout)
}
