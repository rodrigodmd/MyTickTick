package immediatetask

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"myticktick/internal/core/domain"
)

// Repository implementa el repositorio para tareas inmediatas
type Repository struct {
	db *gorm.DB
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create crea una nueva tarea inmediata
func (r *Repository) Create(ctx context.Context, task *domain.ImmediateTask) error {
	taskDB := &ImmediateTaskDB{
		UserID:      task.UserID,
		Name:        task.Name,
		Description: task.Description,
		DueDate:     task.DueDate,
		IsCompleted: task.IsCompleted,
		Priority:    task.Priority,
	}

	if err := r.db.WithContext(ctx).Create(taskDB).Error; err != nil {
		return fmt.Errorf("failed to create immediate task: %w", err)
	}

	task.ID = taskDB.ID
	task.CreatedAt = taskDB.CreatedAt
	task.UpdatedAt = taskDB.UpdatedAt

	return nil
}

// FindByID busca una tarea inmediata por ID
func (r *Repository) FindByID(ctx context.Context, id uint) (*domain.ImmediateTask, error) {
	var taskDB ImmediateTaskDB
	if err := r.db.WithContext(ctx).First(&taskDB, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to find immediate task: %w", err)
	}

	return taskDB.ToDomain(), nil
}

// FindByUserID busca todas las tareas inmediatas de un usuario ordenadas por fecha límite
// Si userID es 0, devuelve todas las tareas (app single-user)
func (r *Repository) FindByUserID(ctx context.Context, userID uint) ([]*domain.ImmediateTask, error) {
	var tasksDB []ImmediateTaskDB

	query := r.db.WithContext(ctx)
	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Order("due_date ASC").Find(&tasksDB).Error; err != nil {
		return nil, fmt.Errorf("failed to find immediate tasks: %w", err)
	}

	tasks := make([]*domain.ImmediateTask, 0, len(tasksDB))
	for i := range tasksDB {
		tasks = append(tasks, tasksDB[i].ToDomain())
	}

	return tasks, nil
}

// Update actualiza una tarea inmediata
func (r *Repository) Update(ctx context.Context, task *domain.ImmediateTask) error {
	// Cargar la fila existente para preservar created_at
	var existing ImmediateTaskDB
	if err := r.db.WithContext(ctx).First(&existing, task.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return fmt.Errorf("failed to find immediate task for update: %w", err)
	}

	existing.UserID = task.UserID
	existing.Name = task.Name
	existing.Description = task.Description
	existing.DueDate = task.DueDate
	existing.IsCompleted = task.IsCompleted
	existing.Priority = task.Priority

	if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return fmt.Errorf("failed to update immediate task: %w", err)
	}

	task.CreatedAt = existing.CreatedAt
	task.UpdatedAt = existing.UpdatedAt

	return nil
}

// Delete elimina una tarea inmediata
func (r *Repository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&ImmediateTaskDB{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete immediate task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
