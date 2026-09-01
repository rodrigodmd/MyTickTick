package trackerentry

import (
	"time"

	"myticktick/internal/core/domain"
)

// TrackerEntryDB representa el modelo de entrada de tracker para GORM
type TrackerEntryDB struct {
	ID        uint      `gorm:"primaryKey"`
	TrackerID uint      `gorm:"index"`
	UserID    uint      `gorm:"index"`
	Value     float64   `gorm:"not null"`
	EntryDate time.Time `gorm:"index"`
	Notes     string
	IsMet     bool
	Deviation float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToDomain convierte TrackerEntryDB a TrackerEntry del dominio
func (t *TrackerEntryDB) ToDomain() *domain.TrackerEntry {
	return &domain.TrackerEntry{
		ID:        t.ID,
		TrackerID: t.TrackerID,
		UserID:    t.UserID,
		Value:     t.Value,
		EntryDate: t.EntryDate,
		Notes:     t.Notes,
		IsMet:     t.IsMet,
		Deviation: t.Deviation,
		CreatedAt: t.CreatedAt,
	}
}
