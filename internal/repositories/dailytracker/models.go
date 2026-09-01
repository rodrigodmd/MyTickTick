package dailytracker

import (
	"time"

	"myticktick/internal/core/domain"
)

// DailyTrackerDB representa el modelo de tracker diario para GORM
type DailyTrackerDB struct {
	ID         uint    `gorm:"primaryKey"`
	UserID     uint    `gorm:"index"`
	Name       string  `gorm:"not null"`
	LimitValue float64 `gorm:"not null"`
	// LimitType: "max" (cumple si valor <= LimitValue) o "min" (cumple si valor >= LimitValue)
	LimitType string `gorm:"not null;default:max"`
	Unit      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToDomain convierte DailyTrackerDB a DailyTracker del dominio
func (d *DailyTrackerDB) ToDomain() *domain.DailyTracker {
	return &domain.DailyTracker{
		ID:         d.ID,
		UserID:     d.UserID,
		Name:       d.Name,
		LimitValue: d.LimitValue,
		LimitType:  d.LimitType,
		Unit:       d.Unit,
		IsActive:   d.IsActive,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}
