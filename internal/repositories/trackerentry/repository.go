package trackerentry

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"myticktick/internal/core/domain"
)

// Repository implementa el repositorio para entradas de trackers
type Repository struct {
	db *gorm.DB
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create crea una nueva entrada de tracker
func (r *Repository) Create(ctx context.Context, entry *domain.TrackerEntry) error {
	entryDB := &TrackerEntryDB{
		TrackerID: entry.TrackerID,
		UserID:    entry.UserID,
		Value:     entry.Value,
		EntryDate: entry.EntryDate,
		Notes:     entry.Notes,
		IsMet:     entry.IsMet,
		Deviation: entry.Deviation,
	}

	if err := r.db.WithContext(ctx).Create(entryDB).Error; err != nil {
		return fmt.Errorf("failed to create tracker entry: %w", err)
	}

	entry.ID = entryDB.ID
	entry.CreatedAt = entryDB.CreatedAt

	return nil
}

// Update modifica una entrada existente (ID obligatorio).
func (r *Repository) Update(ctx context.Context, entry *domain.TrackerEntry) error {
	if entry.ID == 0 {
		return fmt.Errorf("tracker entry ID is required for update")
	}

	updates := map[string]any{
		"value":      entry.Value,
		"entry_date": entry.EntryDate,
		"notes":      entry.Notes,
		"is_met":     entry.IsMet,
		"deviation":  entry.Deviation,
	}
	if err := r.db.WithContext(ctx).
		Model(&TrackerEntryDB{}).
		Where("id = ?", entry.ID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update tracker entry: %w", err)
	}
	return nil
}

// FindByTrackerID busca todas las entradas de un tracker ordenadas por fecha (más recientes primero)
func (r *Repository) FindByTrackerID(ctx context.Context, trackerID uint) ([]*domain.TrackerEntry, error) {
	var entriesDB []TrackerEntryDB
	if err := r.db.WithContext(ctx).
		Where("tracker_id = ?", trackerID).
		Order("entry_date DESC, id DESC").
		Find(&entriesDB).Error; err != nil {
		return nil, fmt.Errorf("failed to find tracker entries: %w", err)
	}

	entries := make([]*domain.TrackerEntry, 0, len(entriesDB))
	for i := range entriesDB {
		entries = append(entries, entriesDB[i].ToDomain())
	}
	return entries, nil
}

// FindByTrackerAndDate busca la entrada única de un tracker para una fecha
// exacta. Devuelve gorm.ErrRecordNotFound si no existe.
func (r *Repository) FindByTrackerAndDate(ctx context.Context, trackerID uint, entryDate time.Time) (*domain.TrackerEntry, error) {
	var entryDB TrackerEntryDB
	if err := r.db.WithContext(ctx).
		Where("tracker_id = ? AND entry_date = ?", trackerID, entryDate).
		First(&entryDB).Error; err != nil {
		return nil, err // gorm.ErrRecordNotFound o error de DB
	}
	return entryDB.ToDomain(), nil
}

// FindByTrackerAndRange busca entradas de un tracker dentro de un rango de
// fechas (inclusive), ordenadas de la más antigua a la más reciente.
func (r *Repository) FindByTrackerAndRange(ctx context.Context, trackerID uint, from, to time.Time) ([]*domain.TrackerEntry, error) {
	var entriesDB []TrackerEntryDB
	if err := r.db.WithContext(ctx).
		Where("tracker_id = ? AND entry_date >= ? AND entry_date <= ?", trackerID, from, to).
		Order("entry_date ASC, id ASC").
		Find(&entriesDB).Error; err != nil {
		return nil, fmt.Errorf("failed to find tracker entries in range: %w", err)
	}

	entries := make([]*domain.TrackerEntry, 0, len(entriesDB))
	for i := range entriesDB {
		entries = append(entries, entriesDB[i].ToDomain())
	}
	return entries, nil
}

// FindByUserIDAndDate busca las entradas de un usuario en una fecha
func (r *Repository) FindByUserIDAndDate(ctx context.Context, userID uint, entryDate time.Time) ([]*domain.TrackerEntry, error) {
	var entriesDB []TrackerEntryDB
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND entry_date = ?", userID, entryDate).
		Find(&entriesDB).Error; err != nil {
		return nil, fmt.Errorf("failed to find tracker entries by user and date: %w", err)
	}

	entries := make([]*domain.TrackerEntry, 0, len(entriesDB))
	for i := range entriesDB {
		entries = append(entries, entriesDB[i].ToDomain())
	}
	return entries, nil
}
