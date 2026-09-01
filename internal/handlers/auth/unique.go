package auth

import "strings"

// IsUniqueViolation indica si el error es una violación de unique index.
// GORM/Postgres devuelve "duplicate key value violates unique constraint"
// en el error de driver; lo detectamos por mensaje para no acoplarnos al driver.
func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint")
}
