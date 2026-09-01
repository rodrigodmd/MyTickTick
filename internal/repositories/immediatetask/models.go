package immediatetask

import (
	"time"

	"myticktick/internal/core/domain"
)

// ImmediateTaskDB representa el modelo de tarea inmediata para GORM
type ImmediateTaskDB struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"index"`
	Name        string    `gorm:"not null"`
	Description string
	DueDate     time.Time `gorm:"index"`
	IsCompleted bool
	Priority    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ToDomain convierte ImmediateTaskDB a ImmediateTask del dominio
func (i *ImmediateTaskDB) ToDomain() *domain.ImmediateTask {
	return &domain.ImmediateTask{
		ID:          i.ID,
		UserID:      i.UserID,
		Name:        i.Name,
		Description: i.Description,
		DueDate:     i.DueDate,
		IsCompleted: i.IsCompleted,
		Priority:    i.Priority,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}
