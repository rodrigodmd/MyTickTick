package trackers

import "net/http"

// MapURLs configura las rutas para este handler
func (h *TrackersHandler) MapURLs(mux *http.ServeMux) {
	// 5.1 - Crear tracker
	mux.HandleFunc("POST /api/trackers", h.Create)
	// 5.2 - Listar trackers
	mux.HandleFunc("GET /api/trackers", h.List)
	// 5.3 - Obtener tracker por ID
	mux.HandleFunc("GET /api/trackers/{id}", h.Get)
	// 5.4 - Actualizar tracker
	mux.HandleFunc("PUT /api/trackers/{id}", h.Update)
	// 5.5 - Eliminar tracker
	mux.HandleFunc("DELETE /api/trackers/{id}", h.Delete)
	// 5.6 - Registrar valor diario
	mux.HandleFunc("POST /api/trackers/{id}/records", h.CreateRecord)
	// 5.6b - Crear o re-setear el valor del día (corregir un registro)
	mux.HandleFunc("PUT /api/trackers/{id}/records", h.UpsertRecord)
	// 5.7 - Historial del tracker
	mux.HandleFunc("GET /api/trackers/{id}/history", h.GetHistory)
}
