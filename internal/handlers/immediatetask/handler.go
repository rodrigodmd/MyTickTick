package immediatetask

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"myticktick/internal/core/domain"
	"myticktick/internal/core/services/immediatetask"
)

// ImmediateTaskHandler maneja las solicitudes de tareas inmediatas
type ImmediateTaskHandler struct {
	service *immediatetask.Service
}

// NewHandler crea una nueva instancia del handler
func NewHandler(service *immediatetask.Service) *ImmediateTaskHandler {
	return &ImmediateTaskHandler{service: service}
}

// Create maneja POST /api/immediate-tasks
func (h *ImmediateTaskHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var dueDate time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			http.Error(w, `{"error":"invalid dueDate, expected RFC3339"}`, http.StatusBadRequest)
			return
		}
		dueDate = parsed
	} else {
		dueDate = time.Now()
	}

	task := &domain.ImmediateTask{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
		DueDate:     dueDate,
		Priority:    req.Priority,
		IsCompleted: false,
	}

	created, err := h.service.Create(r.Context(), task)
	if err != nil {
		http.Error(w, `{"error":"failed to create immediate task"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToResponse(created))
}

// List maneja GET /api/immediate-tasks
func (h *ImmediateTaskHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// userID 0: app single-user, devuelve todas
	tasks, err := h.service.List(r.Context(), 0)
	if err != nil {
		http.Error(w, `{"error":"failed to list immediate tasks"}`, http.StatusInternalServerError)
		return
	}

	responses := make([]*ImmediateTaskResponse, 0, len(tasks))
	for _, task := range tasks {
		responses = append(responses, ToResponse(task))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// Get maneja GET /api/immediate-tasks/{id}
func (h *ImmediateTaskHandler) Get(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, `{"error":"failed to get immediate task"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToResponse(task))
}

// Update maneja PUT /api/immediate-tasks/{id}
func (h *ImmediateTaskHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var dueDate time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			http.Error(w, `{"error":"invalid dueDate, expected RFC3339"}`, http.StatusBadRequest)
			return
		}
		dueDate = parsed
	} else {
		dueDate = time.Now()
	}

	task := &domain.ImmediateTask{
		ID:          uint(id),
		UserID:      1, // app single-user
		Name:        req.Name,
		Description: req.Description,
		DueDate:     dueDate,
		Priority:    req.Priority,
		IsCompleted: req.IsCompleted,
	}

	updated, err := h.service.Update(r.Context(), task)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to update immediate task"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToResponse(updated))
}

// Delete maneja DELETE /api/immediate-tasks/{id}
func (h *ImmediateTaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, `{"error":"failed to delete immediate task"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MapURLs configura las rutas para este handler
func (h *ImmediateTaskHandler) MapURLs(mux *http.ServeMux) {
	// 4.1 - Crear tarea inmediata
	mux.HandleFunc("POST /api/immediate-tasks", h.Create)
	// 4.2 - 4.5 pendientes (implementación en sus respectivas tareas)
	mux.HandleFunc("GET /api/immediate-tasks", h.List)
	mux.HandleFunc("GET /api/immediate-tasks/{id}", h.Get)
	mux.HandleFunc("PUT /api/immediate-tasks/{id}", h.Update)
	mux.HandleFunc("DELETE /api/immediate-tasks/{id}", h.Delete)
}
