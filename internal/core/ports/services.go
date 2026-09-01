package ports

import (
	"context"
	"time"

	"myticktick/internal/core/domain"
)

// AuthServicePort define el contrato para el servicio de autenticación
type AuthServicePort interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, error)
}

// TaskServicePort define el contrato para el servicio de tareas
type TaskServicePort interface {
	CreateTask(ctx context.Context, task *domain.ImmediateTask) error
	GetTask(ctx context.Context, id uint) (*domain.ImmediateTask, error)
	GetTasksByUser(ctx context.Context, userID uint) ([]*domain.ImmediateTask, error)
	UpdateTask(ctx context.Context, task *domain.ImmediateTask) error
	DeleteTask(ctx context.Context, id uint) error
}

// TrackerServicePort define el contrato para el servicio de trackers
type TrackerServicePort interface {
	CreateTracker(ctx context.Context, tracker *domain.DailyTracker) error
	GetTracker(ctx context.Context, id uint) (*domain.DailyTracker, error)
	GetTrackersByUser(ctx context.Context, userID uint) ([]*domain.DailyTracker, error)
	UpdateTracker(ctx context.Context, tracker *domain.DailyTracker) error
	DeleteTracker(ctx context.Context, id uint) error
	RecordEntry(ctx context.Context, entry *domain.TrackerEntry) error
}

// MonthlyTaskServicePort define el contrato para el servicio de tareas mensuales
type MonthlyTaskServicePort interface {
	CreateTask(ctx context.Context, task *domain.MonthlyTask) error
	GetTask(ctx context.Context, id uint) (*domain.MonthlyTask, error)
	GetTasksByUser(ctx context.Context, userID uint) ([]*domain.MonthlyTask, error)
	UpdateTask(ctx context.Context, task *domain.MonthlyTask) error
	DeleteTask(ctx context.Context, id uint) error
	GetHistory(ctx context.Context, taskID uint) ([]*domain.MonthlyTaskHistory, error)
	RecordCompletion(ctx context.Context, taskID uint, month int, year int, completedAt *time.Time) (*domain.MonthlyTaskHistory, error)
	ActivateMonthlyTasks(ctx context.Context, now time.Time) error
}
