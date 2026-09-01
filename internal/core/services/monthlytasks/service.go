package monthlytasks

import (
	"context"
	"time"

	"myticktick/internal/core/domain"
	"myticktick/internal/core/ports"
)

// Service maneja la lógica de negocio de tareas mensuales
type Service struct {
	taskRepo    ports.MonthlyTaskRepository
	historyRepo ports.MonthlyTaskHistoryRepository
}

// NewService crea una nueva instancia de Service
func NewService(taskRepo ports.MonthlyTaskRepository, historyRepo ports.MonthlyTaskHistoryRepository) *Service {
	return &Service{taskRepo: taskRepo, historyRepo: historyRepo}
}

// Create crea una nueva tarea mensual
func (s *Service) Create(ctx context.Context, userID uint, name, description string) (*domain.MonthlyTask, error) {
	task := &domain.MonthlyTask{
		UserID:      userID,
		Name:        name,
		Description: description,
		IsActive:    true,
	}
	err := s.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// List obtiene todas las tareas mensuales
func (s *Service) List(ctx context.Context) ([]*domain.MonthlyTask, error) {
	return s.taskRepo.FindByUserID(ctx, 0)
}

// Get obtiene una tarea mensual por ID
func (s *Service) Get(ctx context.Context, id uint) (*domain.MonthlyTask, error) {
	return s.taskRepo.FindByID(ctx, id)
}

// Update actualiza una tarea mensual
func (s *Service) Update(ctx context.Context, id uint, name, description string) (*domain.MonthlyTask, error) {
	task, err := s.taskRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	task.Name = name
	task.Description = description
	err = s.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Delete elimina una tarea mensual
func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.taskRepo.Delete(ctx, id)
}

// GetHistory obtiene el historial de una tarea mensual
func (s *Service) GetHistory(ctx context.Context, id uint) ([]*domain.MonthlyTaskHistory, error) {
	return s.historyRepo.FindByTaskID(ctx, id)
}

// ActivateMonthlyTasks activa automáticamente las tareas mensuales al arrancar el mes.
//
// Comportamiento (spec monthly-tasks / "Tareas mensuales recurrentes"):
//  1. Asegura que todas las tareas mensuales activas tengan un registro
//     de historial para el mes/año actual en estado pendiente.
//  2. Re-activa cualquier tarea que quede inactiva (por error o uso manual).
//
// Es idempotente: invocarla varias veces no duplica registros.
// Devuelve la cantidad de registros pendientes creados en esta ejecución.
func (s *Service) ActivateMonthlyTasks(ctx context.Context, now time.Time) (int, error) {
	tasks, err := s.taskRepo.FindByUserID(ctx, 0)
	if err != nil {
		return 0, err
	}

	month := int(now.Month())
	year := now.Year()

	created := 0
	for _, task := range tasks {
		// 1) Re-activar si quedó desactivada.
		if !task.IsActive {
			task.IsActive = true
			if err := s.taskRepo.Update(ctx, task); err != nil {
				return 0, err
			}
		}

		// 2) Asegurar registro pendiente del mes en curso.
		didCreate, err := s.historyRepo.EnsureCurrentPeriod(ctx, task.ID, task.UserID, month, year)
		if err != nil {
			return 0, err
		}
		if didCreate {
			created++
		}
	}
	return created, nil
}

// RecordCompletion marca la tarea mensual como completada para el mes/año indicado
// (spec monthly-tasks / "Marcar tarea como cumplida"). Crea el registro si no existe.
func (s *Service) RecordCompletion(ctx context.Context, taskID uint, month int, year int, completedAt *time.Time) (*domain.MonthlyTaskHistory, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	existing, err := s.historyRepo.FindByTaskAndPeriod(ctx, taskID, month, year)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		h := &domain.MonthlyTaskHistory{
			MonthlyTaskID: taskID,
			UserID:        task.UserID,
			Month:         month,
			Year:          year,
			Completed:     true,
			CompletedAt:   completedAt,
		}
		if err := s.historyRepo.Create(ctx, h); err != nil {
			return nil, err
		}
		return h, nil
	}

	existing.Completed = true
	existing.CompletedAt = completedAt
	if err := s.historyRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}
