package dailytracker

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"myticktick/internal/core/domain"
)

// Repository implementa el repositorio para trackers diarios
type Repository struct {
	db *gorm.DB
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create crea un nuevo tracker diario
func (r *Repository) Create(ctx context.Context, tracker *domain.DailyTracker) error {
	trackerDB := &DailyTrackerDB{
		UserID:     tracker.UserID,
		Name:       tracker.Name,
		LimitValue: tracker.LimitValue,
		LimitType:  tracker.LimitType,
		Unit:       tracker.Unit,
		IsActive:   tracker.IsActive,
	}

	if err := r.db.WithContext(ctx).Create(trackerDB).Error; err != nil {
		return fmt.Errorf("failed to create daily tracker: %w", err)
	}

	tracker.ID = trackerDB.ID
	tracker.CreatedAt = trackerDB.CreatedAt
	tracker.UpdatedAt = trackerDB.UpdatedAt

	return nil
}

// FindByID busca un tracker por ID
func (r *Repository) FindByID(ctx context.Context, id uint) (*domain.DailyTracker, error) {
	var trackerDB DailyTrackerDB
	if err := r.db.WithContext(ctx).First(&trackerDB, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find daily tracker by id: %w", err)
	}

	return trackerDB.ToDomain(), nil
}

// FindByUserID busca todos los trackers de un usuario.
// Si userID es 0, devuelve todos los trackers (app single-user por ahora).
func (r *Repository) FindByUserID(ctx context.Context, userID uint) ([]*domain.DailyTracker, error) {
	var trackersDB []DailyTrackerDB

	query := r.db.WithContext(ctx).Order("created_at ASC")
	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&trackersDB).Error; err != nil {
		return nil, fmt.Errorf("failed to find daily trackers: %w", err)
	}

	trackers := make([]*domain.DailyTracker, 0, len(trackersDB))
	for i := range trackersDB {
		trackers = append(trackers, trackersDB[i].ToDomain())
	}

	return trackers, nil
}

// Update actualiza un tracker
func (r *Repository) Update(ctx context.Context, tracker *domain.DailyTracker) error {
	// Cargar la fila existente para preservar created_at
	var existing DailyTrackerDB
	if err := r.db.WithContext(ctx).First(&existing, tracker.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return fmt.Errorf("failed to find daily tracker for update: %w", err)
	}

	existing.UserID = tracker.UserID
	existing.Name = tracker.Name
	existing.LimitValue = tracker.LimitValue
	existing.LimitType = tracker.LimitType
	existing.Unit = tracker.Unit
	existing.IsActive = tracker.IsActive

	if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return fmt.Errorf("failed to update daily tracker: %w", err)
	}

	tracker.CreatedAt = existing.CreatedAt
	tracker.UpdatedAt = existing.UpdatedAt

	return nil
}

// Delete elimina un tracker
func (r *Repository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&DailyTrackerDB{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete daily tracker: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
