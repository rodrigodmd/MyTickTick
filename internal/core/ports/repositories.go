package ports

import (
	"context"
	"errors"
	"time"

	"myticktick/internal/core/domain"
)

// Errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UserRepository define el contrato para operaciones con usuarios
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
}

// TaskRepository define el contrato para operaciones con tareas
type TaskRepository interface {
	Create(ctx context.Context, task *domain.ImmediateTask) error
	FindByID(ctx context.Context, id uint) (*domain.ImmediateTask, error)
	FindByUserID(ctx context.Context, userID uint) ([]*domain.ImmediateTask, error)
	Update(ctx context.Context, task *domain.ImmediateTask) error
	Delete(ctx context.Context, id uint) error
}

// TrackerRepository define el contrato para operaciones con trackers
type TrackerRepository interface {
	Create(ctx context.Context, tracker *domain.DailyTracker) error
	FindByID(ctx context.Context, id uint) (*domain.DailyTracker, error)
	FindByUserID(ctx context.Context, userID uint) ([]*domain.DailyTracker, error)
	Update(ctx context.Context, tracker *domain.DailyTracker) error
	Delete(ctx context.Context, id uint) error
}

// TrackerEntryRepository define el contrato para entradas de trackers
type TrackerEntryRepository interface {
	Create(ctx context.Context, entry *domain.TrackerEntry) error
	Update(ctx context.Context, entry *domain.TrackerEntry) error
	FindByTrackerID(ctx context.Context, trackerID uint) ([]*domain.TrackerEntry, error)
	FindByTrackerAndDate(ctx context.Context, trackerID uint, entryDate time.Time) (*domain.TrackerEntry, error)
	FindByTrackerAndRange(ctx context.Context, trackerID uint, from, to time.Time) ([]*domain.TrackerEntry, error)
	FindByUserIDAndDate(ctx context.Context, userID uint, entryDate time.Time) ([]*domain.TrackerEntry, error)
}

// MonthlyTaskRepository define el contrato para operaciones con tareas mensuales
type MonthlyTaskRepository interface {
	Create(ctx context.Context, task *domain.MonthlyTask) error
	FindByID(ctx context.Context, id uint) (*domain.MonthlyTask, error)
	FindByUserID(ctx context.Context, userID uint) ([]*domain.MonthlyTask, error)
	Update(ctx context.Context, task *domain.MonthlyTask) error
	Delete(ctx context.Context, id uint) error
}

// MonthlyTaskHistoryRepository define el contrato para historial de tareas mensuales
type MonthlyTaskHistoryRepository interface {
	Create(ctx context.Context, history *domain.MonthlyTaskHistory) error
	Update(ctx context.Context, history *domain.MonthlyTaskHistory) error
	FindByTaskID(ctx context.Context, taskID uint) ([]*domain.MonthlyTaskHistory, error)
	FindByUserID(ctx context.Context, userID uint) ([]*domain.MonthlyTaskHistory, error)
	FindByTaskAndPeriod(ctx context.Context, taskID uint, month int, year int) (*domain.MonthlyTaskHistory, error)
	FindByPeriod(ctx context.Context, month int, year int) ([]*domain.MonthlyTaskHistory, error)
	EnsureCurrentPeriod(ctx context.Context, taskID uint, userID uint, month int, year int) (created bool, err error)
}
