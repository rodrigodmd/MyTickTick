package monthlytaskhistory

import (
	"time"

	"myticktick/internal/core/domain"
)

// MonthlyTaskHistoryDB representa el modelo de historial de tarea mensual para GORM
type MonthlyTaskHistoryDB struct {
	ID            uint      `gorm:"primaryKey"`
	MonthlyTaskID uint      `gorm:"index"`
	UserID        uint      `gorm:"index"`
	Month         int
	Year          int
	Completed     bool
	CompletedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ToDomain convierte MonthlyTaskHistoryDB a MonthlyTaskHistory del dominio
func (m *MonthlyTaskHistoryDB) ToDomain() *domain.MonthlyTaskHistory {
	return &domain.MonthlyTaskHistory{
		ID:            m.ID,
		MonthlyTaskID: m.MonthlyTaskID,
		UserID:        m.UserID,
		Month:         m.Month,
		Year:          m.Year,
		Completed:     m.Completed,
		CompletedAt:   m.CompletedAt,
		CreatedAt:     m.CreatedAt,
	}
}
