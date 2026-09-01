package immediatetasks

// CreateRequest representa la solicitud para crear una tarea inmediata
type CreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"dueDate"`
	Priority    string `json:"priority,omitempty"`
}

// UpdateRequest representa la solicitud para actualizar una tarea inmediata
type UpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"dueDate"`
	Priority    string `json:"priority,omitempty"`
	IsCompleted bool   `json:"isCompleted"`
}

// ImmediateTaskResponse representa la respuesta de una tarea inmediata
type ImmediateTaskResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
	Priority    string `json:"priority"`
	IsCompleted bool   `json:"isCompleted"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
