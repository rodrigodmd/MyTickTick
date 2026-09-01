package monthlytasks

import (
	"encoding/json"
	"time"

	"myticktick/internal/core/domain"
)

// CreateRequest representa la solicitud para crear una tarea mensual
type CreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	UserID      uint   `json:"userId"`
}

// UpdateRequest representa la solicitud para actualizar una tarea mensual (PUT)
type UpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// HistoryResponse representa un registro del historial de cumplimiento
type HistoryResponse struct {
	ID          uint    `json:"id"`
	MonthlyTaskID uint  `json:"monthlyTaskId"`
	Month       int     `json:"month"`
	Year        int     `json:"year"`
	Completed   bool    `json:"completed"`
	CompletedAt *string `json:"completedAt"`
}

// ToHistoryResponse convierte una entidad MonthlyTaskHistory a HistoryResponse
func ToHistoryResponse(h *domain.MonthlyTaskHistory) *HistoryResponse {
	var completedAt *string
	if h.CompletedAt != nil {
		v := h.CompletedAt.Format(time.RFC3339)
		completedAt = &v
	}
	return &HistoryResponse{
		ID:            h.ID,
		MonthlyTaskID: h.MonthlyTaskID,
		Month:         h.Month,
		Year:          h.Year,
		Completed:     h.Completed,
		CompletedAt:   completedAt,
	}
}

// MonthlyTaskResponse representa la respuesta de una tarea mensual
type MonthlyTaskResponse struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"userId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ToResponse convierte una entidad MonthlyTask a MonthlyTaskResponse
func ToResponse(task *domain.MonthlyTask) *MonthlyTaskResponse {
	return &MonthlyTaskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		Name:        task.Name,
		Description: task.Description,
		IsActive:    task.IsActive,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}
}

// FromRequest crea una entidad MonthlyTask desde CreateRequest
func FromRequest(req CreateRequest) *domain.MonthlyTask {
	now := time.Now()
	return &domain.MonthlyTask{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ParseJSON parsea JSON a CreateRequest
func ParseJSON(data []byte) (CreateRequest, error) {
	var req CreateRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		return CreateRequest{}, err
	}
	return req, nil
}
