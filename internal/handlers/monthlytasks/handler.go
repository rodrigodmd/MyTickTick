package monthlytasks

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"myticktick/internal/core/services/monthlytasks"
)

// MonthlyTaskHandler maneja las solicitudes de tareas mensuales
type MonthlyTaskHandler struct {
	service *monthlytasks.Service
}

// NewMonthlyTaskHandler crea una nueva instancia del handler
func NewMonthlyTaskHandler(service *monthlytasks.Service) *MonthlyTaskHandler {
	return &MonthlyTaskHandler{service: service}
}

// Create maneja POST /api/monthly-tasks
func (h *MonthlyTaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
		return
	}

	// Crear la tarea desde el request
	task, err := h.service.Create(r.Context(), req.UserID, req.Name, req.Description)
	if err != nil {
		http.Error(w, `{"error":"failed to create monthly task"}`, http.StatusInternalServerError)
		return
	}

	// Devolver la tarea creada
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToResponse(task))
}

// List maneja GET /api/monthly-tasks
func (h *MonthlyTaskHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to list monthly tasks"}`, http.StatusInternalServerError)
		return
	}

	// Convertir a formato de respuesta consistente con Create
	responses := make([]*MonthlyTaskResponse, 0, len(tasks))
	for _, task := range tasks {
		responses = append(responses, ToResponse(task))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// Get maneja GET /api/monthly-tasks/{id}
func (h *MonthlyTaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraer ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, `{"error":"invalid task id"}`, http.StatusBadRequest)
		return
	}

	task, err := h.service.Get(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to get monthly task"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToResponse(task))
}

// Update maneja PUT /api/monthly-tasks/{id}
func (h *MonthlyTaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraer ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, `{"error":"invalid task id"}`, http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
		return
	}

	task, err := h.service.Update(r.Context(), uint(id), req.Name, req.Description)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to update monthly task"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToResponse(task))
}

// Delete maneja DELETE /api/monthly-tasks/{id}
func (h *MonthlyTaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraer ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, `{"error":"invalid task id"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to delete monthly task"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetHistory maneja GET /api/monthly-tasks/{id}/history
func (h *MonthlyTaskHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraer ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, `{"error":"invalid task id"}`, http.StatusBadRequest)
		return
	}

	// Verificar que la tarea exista para dar 404 consistente con Get
	if _, err := h.service.Get(r.Context(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to get monthly task"}`, http.StatusInternalServerError)
		return
	}

	history, err := h.service.GetHistory(r.Context(), uint(id))
	if err != nil {
		http.Error(w, `{"error":"failed to get history"}`, http.StatusInternalServerError)
		return
	}

	responses := make([]*HistoryResponse, 0, len(history))
	for _, item := range history {
		responses = append(responses, ToHistoryResponse(item))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// Activate maneja POST /api/monthly-tasks/activate
// Dispara la activación automática del mes en curso (tarea 6.1).
// También se ejecuta sola al arrancar el servidor (main.go).
func (h *MonthlyTaskHandler) Activate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	created, err := h.service.ActivateMonthlyTasks(r.Context(), time.Now().UTC())
	if err != nil {
		http.Error(w, `{"error":"failed to activate monthly tasks"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"activated":      true,
		"recordsCreated": created,
	})
}

// RecordCompletion maneja PUT /api/monthly-tasks/{id}/completion
// Marca la tarea como cumplida para el mes/año indicado (default: mes en curso).
func (h *MonthlyTaskHandler) RecordCompletion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, `{"error":"invalid task id"}`, http.StatusBadRequest)
		return
	}

	var req CompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Body vacío: se usa el mes en curso.
		req = CompletionRequest{}
	}

	now := time.Now().UTC()
	month, year := int(now.Month()), now.Year()
	if req.Month != nil {
		month = *req.Month
	}
	if req.Year != nil {
		year = *req.Year
	}
	if month < 1 || month > 12 || year < 1 {
		http.Error(w, `{"error":"invalid month/year"}`, http.StatusBadRequest)
		return
	}

	history, err := h.service.RecordCompletion(r.Context(), uint(id), month, year, &now)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to record completion"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToHistoryResponse(history))
}

// CompletionRequest representa la solicitud para marcar una tarea como cumplida
type CompletionRequest struct {
	Month *int `json:"month,omitempty"`
	Year  *int `json:"year,omitempty"`
}

// MapURLs configura las rutas para este handler
func (h *MonthlyTaskHandler) MapURLs(mux *http.ServeMux) {
	// 3.1 - Crear tarea mensual
	mux.HandleFunc("POST /api/monthly-tasks", h.Create)
	// 3.2 - Listar tareas mensuales
	mux.HandleFunc("GET /api/monthly-tasks", h.List)
	// 3.3 - Obtener tarea mensual por ID
	mux.HandleFunc("GET /api/monthly-tasks/{id}", h.Get)
	// 3.4 - Actualizar tarea mensual
	mux.HandleFunc("PUT /api/monthly-tasks/{id}", h.Update)
	// 3.5 - Eliminar tarea mensual
	mux.HandleFunc("DELETE /api/monthly-tasks/{id}", h.Delete)
	// 3.6 - Historial de tarea mensual
	mux.HandleFunc("GET /api/monthly-tasks/{id}/history", h.GetHistory)
	// 6.1 - Activación automática del mes en curso (idempotente)
	mux.HandleFunc("POST /api/monthly-tasks/activate", h.Activate)
	// 6.1 - Marcar tarea como cumplida para un mes
	mux.HandleFunc("PUT /api/monthly-tasks/{id}/completion", h.RecordCompletion)
}
