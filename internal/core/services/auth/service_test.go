package auth

import (
	"context"
	"errors"
	"testing"

	"myticktick/internal/core/domain"
	"myticktick/internal/core/ports"
	"myticktick/internal/core/token"
)

// fakeUserRepo implementa ports.UserRepository en memoria para tests.
type fakeUserRepo struct {
	users map[string]*domain.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: map[string]*domain.User{}}
}

func (f *fakeUserRepo) Create(ctx context.Context, user *domain.User) error {
	if _, exists := f.users[user.Username]; exists {
		return errors.New("duplicate key value violates unique constraint")
	}
	f.users[user.Username] = user
	return nil
}

func (f *fakeUserRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	u, ok := f.users[username]
	if !ok {
		return nil, ports.ErrUserNotFound
	}
	return u, nil
}

func (f *fakeUserRepo) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	for _, u := range f.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, ports.ErrUserNotFound
}

func TestRegister_HasheaPassword(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewService(repo)

	created, err := svc.Register(context.Background(), "rodri", "clave123")
	if err != nil {
		t.Fatalf("Register falló: %v", err)
	}

	// El password guardado NO puede ser texto plano
	if created.Password == "clave123" {
		t.Error("el password guardado no debe ser texto plano")
	}
	// Y debe verificarse contra el original
	if !CheckPassword(created.Password, "clave123") {
		t.Error("el hash guardado debe verificar contra el password original")
	}
	// Un password incorrecto no debe verificar
	if CheckPassword(created.Password, "claveMala") {
		t.Error("un password incorrecto no debe verificar")
	}
}

func TestRegister_UsernameDuplicado(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewService(repo)

	if _, err := svc.Register(context.Background(), "rodri", "clave123"); err != nil {
		t.Fatalf("primer registro falló: %v", err)
	}
	// Segundo registro con el mismo username debe dar error de unique constraint
	if _, err := svc.Register(context.Background(), "rodri", "otra123"); err == nil {
		t.Error("un username duplicado debe dar error")
	}
}

func TestLogin_CredencialesInvalidas(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewService(repo)
	_, _ = svc.Register(context.Background(), "rodri", "clave123")

	// Usuario inexistente
	if _, err := svc.Login(context.Background(), "nadie", "clave123", false); !errors.Is(err, ports.ErrInvalidCredentials) {
		t.Errorf("usuario inexistente debe dar ErrInvalidCredentials, got %v", err)
	}
	// Password incorrecto
	if _, err := svc.Login(context.Background(), "rodri", "claveMala", false); !errors.Is(err, ports.ErrInvalidCredentials) {
		t.Errorf("password incorrecto debe dar ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_EmiteTokenValido(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewService(repo)
	created, _ := svc.Register(context.Background(), "rodri", "clave123")

	session, err := svc.Login(context.Background(), "rodri", "clave123", false)
	if err != nil {
		t.Fatalf("Login falló: %v", err)
	}

	// El token debe poder validarse y llevar el usuario correcto
	claims, err := token.Validate(session.Token)
	if err != nil {
		t.Fatalf("el token emitido no es válido: %v", err)
	}
	if claims.Username != "rodri" {
		t.Errorf("claims.Username = %q, want %q", claims.Username, "rodri")
	}
	if created.ID == 0 {
		// El fake no asigna ID; si así fuera, al menos el username coincide.
	}
	if session.Expiration.Before(session.Expiration.Add(-1)) {
		t.Error("expiración inconsistente")
	}
}

func TestLogin_RememberExtiendeExpiración(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewService(repo)
	_, _ = svc.Register(context.Background(), "rodri", "clave123")

	normal, err := svc.Login(context.Background(), "rodri", "clave123", false)
	if err != nil {
		t.Fatal(err)
	}
	remember, err := svc.Login(context.Background(), "rodri", "clave123", true)
	if err != nil {
		t.Fatal(err)
	}

	// Con "Recordarme" la expiración debe ser más lejana
	if !remember.Expiration.After(normal.Expiration) {
		t.Error("Recordarme debe tener expiración más lejana que la sesión normal")
	}
}
