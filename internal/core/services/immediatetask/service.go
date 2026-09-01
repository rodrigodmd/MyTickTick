package immediatetask

import (
	"context"

	"myticktick/internal/core/domain"
	"myticktick/internal/core/ports"
)

// Service maneja la lógica de negocio de tareas inmediatas
type Service struct {
	repository ports.TaskRepository
}

// NewService crea una nueva instancia de Service
func NewService(repo ports.TaskRepository) *Service {
	return &Service{repository: repo}
}

// Create crea una nueva tarea inmediata
func (s *Service) Create(ctx context.Context, task *domain.ImmediateTask) (*domain.ImmediateTask, error) {
	err := s.repository.Create(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// List obtiene todas las tareas inmediatas de un usuario
func (s *Service) List(ctx context.Context, userID uint) ([]*domain.ImmediateTask, error) {
	return s.repository.FindByUserID(ctx, userID)
}

// Get obtiene una tarea inmediata por ID
func (s *Service) Get(ctx context.Context, id uint) (*domain.ImmediateTask, error) {
	return s.repository.FindByID(ctx, id)
}

// Update actualiza una tarea inmediata
func (s *Service) Update(ctx context.Context, task *domain.ImmediateTask) (*domain.ImmediateTask, error) {
	err := s.repository.Update(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Delete elimina una tarea inmediata
func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repository.Delete(ctx, id)
}
