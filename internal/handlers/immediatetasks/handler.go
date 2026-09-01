package immediatetasks

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ImmediateTasksHandler maneja las peticiones de tareas inmediatas
type ImmediateTasksHandler struct {
	// repository interface{} // ImmediateTaskRepository
}

// NewHandler crea una nueva instancia de ImmediateTasksHandler
func NewHandler() *ImmediateTasksHandler {
	return &ImmediateTasksHandler{}
}

// Create maneja la creación de una tarea inmediata
func (h *ImmediateTasksHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Implementation pendiente: crear tarea en repositorio
	task := &ImmediateTaskResponse{
		ID:   1, // Simulado para MVP
		Name: req.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// List maneja la lista de tareas inmediatas
func (h *ImmediateTasksHandler) List(w http.ResponseWriter, r *http.Request) {
	// Implementation pendiente: obtener tareas del repositorio ordenadas por dueDate
	tasks := []*ImmediateTaskResponse{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Get maneja la obtención de una tarea inmediata por ID
func (h *ImmediateTasksHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	// Implementation pendiente: obtener tarea por ID del repositorio
	task := &ImmediateTaskResponse{
		ID:   uint(id),
		Name: "Tarea de ejemplo",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Update maneja la actualización de una tarea inmediata
func (h *ImmediateTasksHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Implementation pendiente: actualizar tarea en repositorio
	task := &ImmediateTaskResponse{
		ID:   uint(id),
		Name: req.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Delete maneja la eliminación de una tarea inmediata
func (h *ImmediateTasksHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	// Implementation pendiente: eliminar tarea del repositorio

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tarea eliminada exitosamente",
		"id":      id,
	})
}
