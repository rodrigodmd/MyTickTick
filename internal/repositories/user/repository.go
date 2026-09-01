package user

import (
	"context"
	"fmt"
	
	"gorm.io/gorm"
	"myticktick/internal/core/domain"
	"myticktick/internal/core/ports"
)

// UserRepository implementa ports.UserRepository para PostgreSQL con GORM
type UserRepository struct {
	db *gorm.DB
}

// NewRepository crea una nueva instancia del repositorio de usuarios
func NewRepository(db *gorm.DB) ports.UserRepository {
	return &UserRepository{db: db}
}

// Create crea un nuevo usuario
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	userDB := &UserDB{
		Username: user.Username,
		Password: user.Password,
	}
	
	if err := r.db.WithContext(ctx).Create(userDB).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	user.ID = userDB.ID
	user.CreatedAt = userDB.CreatedAt
	user.UpdatedAt = userDB.UpdatedAt
	
	return nil
}

// FindByID busca un usuario por ID
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var userDB UserDB
	
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	
	return userDB.ToDomain(), nil
}

// FindByUsername busca un usuario por username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var userDB UserDB
	
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&userDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ports.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	
	return userDB.ToDomain(), nil
}
