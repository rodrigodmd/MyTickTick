package immediatetasks

import "net/http"

// MapURLs configura las rutas para este handler
func (h *ImmediateTasksHandler) MapURLs(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/immediate-tasks", h.Create)
	mux.HandleFunc("GET /api/immediate-tasks", h.List)
	mux.HandleFunc("GET /api/immediate-tasks/{id}", h.Get)
	mux.HandleFunc("PUT /api/immediate-tasks/{id}", h.Update)
	mux.HandleFunc("DELETE /api/immediate-tasks/{id}", h.Delete)
}
