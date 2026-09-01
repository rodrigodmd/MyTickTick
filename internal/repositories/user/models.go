package user

import (
	"time"
	
	"myticktick/internal/core/domain"
)

// UserDB representa el modelo de usuario para GORM
type UserDB struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToDomain convierte UserDB a User del dominio
func (u *UserDB) ToDomain() *domain.User {
	return &domain.User{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
