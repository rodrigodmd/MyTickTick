package monthlytask

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"myticktick/internal/core/domain"
)

// Repository implementa el repositorio para tareas mensuales
type Repository struct {
	db *gorm.DB
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create crea una nueva tarea mensual y persiste en la base de datos
func (r *Repository) Create(ctx context.Context, task *domain.MonthlyTask) error {
	taskDB := &MonthlyTaskDB{
		UserID:      task.UserID,
		Name:        task.Name,
		Description: task.Description,
		IsActive:    task.IsActive,
	}

	if err := r.db.WithContext(ctx).Create(taskDB).Error; err != nil {
		return fmt.Errorf("failed to create monthly task: %w", err)
	}

	task.ID = taskDB.ID
	task.CreatedAt = taskDB.CreatedAt
	task.UpdatedAt = taskDB.UpdatedAt

	return nil
}

// FindByID busca una tarea mensual por ID
func (r *Repository) FindByID(ctx context.Context, id uint) (*domain.MonthlyTask, error) {
	var taskDB MonthlyTaskDB
	if err := r.db.WithContext(ctx).First(&taskDB, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find monthly task by id: %w", err)
	}

	return taskDB.ToDomain(), nil
}

// FindByUserID busca todas las tareas mensuales de un usuario.
// Si userID es 0, devuelve todas las tareas (app single-user por ahora).
func (r *Repository) FindByUserID(ctx context.Context, userID uint) ([]*domain.MonthlyTask, error) {
	var tasksDB []MonthlyTaskDB

	query := r.db.WithContext(ctx).Order("created_at ASC")
	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&tasksDB).Error; err != nil {
		return nil, fmt.Errorf("failed to find monthly tasks: %w", err)
	}

	tasks := make([]*domain.MonthlyTask, 0, len(tasksDB))
	for i := range tasksDB {
		tasks = append(tasks, tasksDB[i].ToDomain())
	}

	return tasks, nil
}

// Update actualiza una tarea mensual
func (r *Repository) Update(ctx context.Context, task *domain.MonthlyTask) error {
	// Cargar la fila existente para preservar created_at
	var existing MonthlyTaskDB
	if err := r.db.WithContext(ctx).First(&existing, task.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return fmt.Errorf("failed to find monthly task for update: %w", err)
	}

	existing.UserID = task.UserID
	existing.Name = task.Name
	existing.Description = task.Description
	existing.IsActive = task.IsActive

	if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return fmt.Errorf("failed to update monthly task: %w", err)
	}

	task.CreatedAt = existing.CreatedAt
	task.UpdatedAt = existing.UpdatedAt

	return nil
}

// Delete elimina una tarea mensual
func (r *Repository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&MonthlyTaskDB{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete monthly task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
