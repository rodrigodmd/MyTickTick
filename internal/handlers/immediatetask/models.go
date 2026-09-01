package immediatetask

import (
	"time"

	"myticktick/internal/core/domain"
)

// CreateRequest representa la solicitud para crear una tarea inmediata
type CreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"dueDate,omitempty"`
	Priority    string `json:"priority,omitempty"`
	UserID      uint   `json:"userId"`
}

// UpdateRequest representa la solicitud para actualizar una tarea inmediata (PUT)
type UpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"dueDate,omitempty"`
	Priority    string `json:"priority,omitempty"`
	IsCompleted bool   `json:"isCompleted"`
}

// ImmediateTaskResponse representa la respuesta de una tarea inmediata
type ImmediateTaskResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"userId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	DueDate     string  `json:"dueDate"`
	IsCompleted bool    `json:"isCompleted"`
	Priority    string  `json:"priority,omitempty"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   *string `json:"updatedAt"`
}

// ToResponse convierte una entidad ImmediateTask a ImmediateTaskResponse.
// Las fechas se formatean en UTC para que el día no se deslice con la zona
// local del proceso (mismo patrón que la fix de TrackerEntry).
func ToResponse(task *domain.ImmediateTask) *ImmediateTaskResponse {
	var updatedAt *string
	if !task.UpdatedAt.IsZero() {
		v := task.UpdatedAt.UTC().Format(time.RFC3339)
		updatedAt = &v
	}
	return &ImmediateTaskResponse{
		ID:          task.ID,
		UserID:      task.UserID,
		Name:        task.Name,
		Description: task.Description,
		DueDate:     task.DueDate.UTC().Format(time.RFC3339),
		IsCompleted: task.IsCompleted,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   updatedAt,
	}
}
