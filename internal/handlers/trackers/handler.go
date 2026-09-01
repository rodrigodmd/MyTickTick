package trackers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"myticktick/internal/core/domain"
	trackerService "myticktick/internal/core/services/trackers"
)

// TrackersHandler maneja las solicitudes de trackers
type TrackersHandler struct {
	service *trackerService.Service
}

// NewHandler crea una nueva instancia del handler
func NewHandler(service *trackerService.Service) *TrackersHandler {
	return &TrackersHandler{service: service}
}

// Create maneja POST /api/trackers
func (h *TrackersHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	tracker, err := h.service.Create(r.Context(), FromCreateRequest(req))
	if err != nil {
		http.Error(w, `{"error":"failed to create tracker"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToResponse(tracker))
}

// List maneja GET /api/trackers
func (h *TrackersHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// userID 0: app single-user, devuelve todos
	trackers, err := h.service.List(r.Context(), 0)
	if err != nil {
		http.Error(w, `{"error":"failed to list trackers"}`, http.StatusInternalServerError)
		return
	}

	responses := make([]*TrackerResponse, 0, len(trackers))
	for _, tracker := range trackers {
		responses = append(responses, ToResponse(tracker))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// Get maneja GET /api/trackers/{id}
func (h *TrackersHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraer ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, `{"error":"invalid tracker id"}`, http.StatusBadRequest)
		return
	}

	tracker, err := h.service.Get(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error":"tracker not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to get tracker"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToResponse(tracker))
}

// Update maneja PUT /api/trackers/{id}
func (h *TrackersHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraer ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, `{"error":"invalid tracker id"}`, http.StatusBadRequest)
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

	tracker, err := h.service.Update(r.Context(), FromUpdateRequest(uint(id), req))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error":"tracker not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to update tracker"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToResponse(tracker))
}

// Delete maneja DELETE /api/trackers/{id}
func (h *TrackersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraer ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, `{"error":"invalid tracker id"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error":"tracker not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to delete tracker"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateRecord maneja POST /api/trackers/{id}/records (crear registro)
func (h *TrackersHandler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	entry, status, ok := h.processRecord(w, r, false, http.MethodPost)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ToEntryResponse(entry))
}

// UpsertRecord maneja PUT /api/trackers/{id}/records (crear o re-setear el
// valor del día). Si ya hay un registro para esa fecha lo actualiza en vez de
// duplicarlo, permitiendo corregir un valor marcado por error.
func (h *TrackersHandler) UpsertRecord(w http.ResponseWriter, r *http.Request) {
	entry, status, ok := h.processRecord(w, r, true, http.MethodPut)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ToEntryResponse(entry))
}

// processRecord valida método, ID y body, y delega a CreateEntry o UpsertEntry.
// Devuelve el entry resultante, el status HTTP y ok=false si ya escribió un
// error HTTP en w.
func (h *TrackersHandler) processRecord(w http.ResponseWriter, r *http.Request, upsert bool, method string) (*domain.TrackerEntry, int, bool) {
	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil, 0, false
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, `{"error":"invalid tracker id"}`, http.StatusBadRequest)
		return nil, 0, false
	}

	var req RecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return nil, 0, false
	}

	entry := FromRecordRequest(uint(id), req)

	var result *domain.TrackerEntry
	if upsert {
		result, err = h.service.UpsertEntry(r.Context(), entry)
	} else {
		result, err = h.service.CreateEntry(r.Context(), entry)
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error":"tracker not found"}`, http.StatusNotFound)
			return nil, 0, false
		}
		http.Error(w, `{"error":"failed to save record"}`, http.StatusInternalServerError)
		return nil, 0, false
	}

	status := http.StatusCreated
	if upsert && result.ID == 0 {
		// Upsert que terminó en "ya existía" (ID preservado por el repo) -> 200.
		status = http.StatusOK
	}
	return result, status, true
}

// GetHistory maneja GET /api/trackers/{id}/history
func (h *TrackersHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraer ID de la URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, `{"error":"invalid tracker id"}`, http.StatusBadRequest)
		return
	}

	// Verificar que el tracker exista para dar 404 consistente con Get
	if _, err := h.service.Get(r.Context(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error":"tracker not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to get tracker"}`, http.StatusInternalServerError)
		return
	}

	history, err := h.service.GetHistory(r.Context(), uint(id))
	if err != nil {
		http.Error(w, `{"error":"failed to get history"}`, http.StatusInternalServerError)
		return
	}

	responses := make([]*TrackerEntryResponse, 0, len(history))
	for _, entry := range history {
		responses = append(responses, ToEntryResponse(entry))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}
