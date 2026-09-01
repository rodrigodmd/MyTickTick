package monthlytaskhistory

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"myticktick/internal/core/domain"
)

// Repository implementa el repositorio para historial de tareas mensuales
type Repository struct {
	db *gorm.DB
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create crea un nuevo registro de historial
func (r *Repository) Create(ctx context.Context, history *domain.MonthlyTaskHistory) error {
	historyDB := &MonthlyTaskHistoryDB{
		MonthlyTaskID: history.MonthlyTaskID,
		UserID:        history.UserID,
		Month:         history.Month,
		Year:          history.Year,
		Completed:     history.Completed,
		CompletedAt:   history.CompletedAt,
	}

	if err := r.db.WithContext(ctx).Create(historyDB).Error; err != nil {
		return fmt.Errorf("failed to create monthly task history: %w", err)
	}

	history.ID = historyDB.ID
	history.CreatedAt = historyDB.CreatedAt

	return nil
}

// Update actualiza un registro de historial existente
func (r *Repository) Update(ctx context.Context, history *domain.MonthlyTaskHistory) error {
	var historyDB MonthlyTaskHistoryDB
	if err := r.db.WithContext(ctx).First(&historyDB, history.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return fmt.Errorf("failed to find monthly task history for update: %w", err)
	}

	historyDB.Completed = history.Completed
	historyDB.CompletedAt = history.CompletedAt

	if err := r.db.WithContext(ctx).Save(&historyDB).Error; err != nil {
		return fmt.Errorf("failed to update monthly task history: %w", err)
	}

	return nil
}

// FindByTaskID busca el historial de una tarea mensual
func (r *Repository) FindByTaskID(ctx context.Context, taskID uint) ([]*domain.MonthlyTaskHistory, error) {
	var historyDB []MonthlyTaskHistoryDB
	if err := r.db.WithContext(ctx).Where("monthly_task_id = ?", taskID).Order("year DESC, month DESC").Find(&historyDB).Error; err != nil {
		return nil, fmt.Errorf("failed to find monthly task history: %w", err)
	}

	history := make([]*domain.MonthlyTaskHistory, 0, len(historyDB))
	for i := range historyDB {
		history = append(history, historyDB[i].ToDomain())
	}

	return history, nil
}

// FindByUserID busca todo el historial de un usuario
func (r *Repository) FindByUserID(ctx context.Context, userID uint) ([]*domain.MonthlyTaskHistory, error) {
	return nil, nil
}

// FindByPeriod busca todos los registros de historial de un mes/año dado
// (independientemente de la tarea). Se usa para métricas agregadas (6.4).
func (r *Repository) FindByPeriod(ctx context.Context, month int, year int) ([]*domain.MonthlyTaskHistory, error) {
	var historyDB []MonthlyTaskHistoryDB
	if err := r.db.WithContext(ctx).
		Where("month = ? AND year = ?", month, year).
		Order("monthly_task_id ASC").
		Find(&historyDB).Error; err != nil {
		return nil, fmt.Errorf("failed to find monthly task history by period: %w", err)
	}

	history := make([]*domain.MonthlyTaskHistory, 0, len(historyDB))
	for i := range historyDB {
		history = append(history, historyDB[i].ToDomain())
	}
	return history, nil
}

// FindByTaskAndPeriod busca el historial de una tarea en un periodo específico
func (r *Repository) FindByTaskAndPeriod(ctx context.Context, taskID uint, month int, year int) (*domain.MonthlyTaskHistory, error) {
	var historyDB MonthlyTaskHistoryDB
	err := r.db.WithContext(ctx).
		Where("monthly_task_id = ? AND month = ? AND year = ?", taskID, month, year).
		First(&historyDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find monthly task history by period: %w", err)
	}
	return historyDB.ToDomain(), nil
}

// EnsureCurrentPeriod asegura que exista un registro de historial pendiente
// para el mes/año indicado. Si el registro ya existe no hace nada (idempotente).
// Si existe pero estaba marcado como completado, no lo altera.
// Devuelve true si creó el registro, false si ya existía.
func (r *Repository) EnsureCurrentPeriod(ctx context.Context, taskID uint, userID uint, month int, year int) (bool, error) {
	existing, err := r.FindByTaskAndPeriod(ctx, taskID, month, year)
	if err != nil {
		return false, err
	}
	if existing != nil {
		return false, nil
	}

	history := &MonthlyTaskHistoryDB{
		MonthlyTaskID: taskID,
		UserID:        userID,
		Month:         month,
		Year:          year,
		Completed:     false,
		CompletedAt:   nil,
	}
	if err := r.db.WithContext(ctx).Create(history).Error; err != nil {
		return false, fmt.Errorf("failed to create monthly task history: %w", err)
	}
	return true, nil
}
