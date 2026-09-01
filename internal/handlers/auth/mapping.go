package auth

import "myticktick/internal/core/domain"

// ToResponseUser convierte un usuario del dominio a ResponseUser
func ToResponseUser(user *domain.User) ResponseUser {
	return ResponseUser{
		ID:       user.ID,
		Username: user.Username,
	}
}
