package domain

import "time"

// User representa un usuario del sistema
type User struct {
	ID       uint
	Username string
	Password string // hash bcrypt del password (no texto plano)
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ImmediateTask representa una tarea inmediata con fecha límite
type ImmediateTask struct {
	ID          uint
	UserID      uint
	Name        string
	Description string
	DueDate     time.Time
	IsCompleted bool
	Priority    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// DailyTracker representa un tracker de métrica diaria con límite unilateral
type DailyTracker struct {
	ID          uint
	UserID      uint
	Name        string
	LimitValue  float64 // el límite (ej. peso máx 85, sueño mín 6)
	LimitType   string  // "max" = cumple si valor <= LimitValue; "min" = cumple si valor >= LimitValue
	Unit        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TrackerEntry representa una entrada diaria en un tracker
type TrackerEntry struct {
	ID        uint
	TrackerID uint
	UserID    uint
	Value     float64
	EntryDate time.Time
	Notes     string
	IsMet     bool
	Deviation float64
	CreatedAt time.Time
}

// MonthlyTask representa una tarea mensual recurrente
type MonthlyTask struct {
	ID        uint
	UserID    uint
	Name      string
	Description string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MonthlyTaskHistory representa el historial de cumplimiento de una tarea mensual
type MonthlyTaskHistory struct {
	ID            uint
	MonthlyTaskID uint
	UserID        uint
	Month         int
	Year          int
	Completed     bool
	CompletedAt   *time.Time
	CreatedAt     time.Time
}
