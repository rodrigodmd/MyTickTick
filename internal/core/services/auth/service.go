package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"myticktick/internal/core/domain"
	"myticktick/internal/core/ports"
	"myticktick/internal/core/token"
)

// AuthService implementa el servicio de autenticación:
// registro con password hasheado (bcrypt) y login con emisión de token.
type AuthService struct {
	userRepo ports.UserRepository
}

// NewService crea una nueva instancia de AuthService.
func NewService(userRepo ports.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Session es el resultado de un login exitoso.
type Session struct {
	User       *domain.User
	Token      string
	Expiration time.Time
}

// Register registra un nuevo usuario con username único y password hasheado.
func (s *AuthService) Register(ctx context.Context, username, password string) (*domain.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Username: username,
		Password: string(hash), // campo "Password" = hash (bcrypt), ver domain.User
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login autentica un usuario con username y password y emite un token.
// remember=true -> expiración larga ("Recordarme"); false -> corta.
func (s *AuthService) Login(ctx context.Context, username, password string, remember bool) (*Session, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return nil, ports.ErrInvalidCredentials
		}
		return nil, err
	}

	if !CheckPassword(user.Password, password) {
		return nil, ports.ErrInvalidCredentials
	}

	exp := token.ExpirationFor(remember)
	tok, err := token.Generate(user.ID, username, exp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &Session{User: user, Token: tok, Expiration: exp}, nil
}

// CheckPassword compara un bcrypt hash con el password en texto plano.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Exponer errores del paquete ports para que los handlers los usen
var (
	ErrUserNotFound       = ports.ErrUserNotFound
	ErrInvalidCredentials = ports.ErrInvalidCredentials
)
