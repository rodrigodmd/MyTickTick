package monthlytask

import (
	"time"

	"myticktick/internal/core/domain"
)

// MonthlyTaskDB representa el modelo de tarea mensual para GORM
type MonthlyTaskDB struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"index"`
	Name        string    `gorm:"not null"`
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ToDomain convierte MonthlyTaskDB a MonthlyTask del dominio
func (m *MonthlyTaskDB) ToDomain() *domain.MonthlyTask {
	return &domain.MonthlyTask{
		ID:          m.ID,
		UserID:      m.UserID,
		Name:        m.Name,
		Description: m.Description,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
