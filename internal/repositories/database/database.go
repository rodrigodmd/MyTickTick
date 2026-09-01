package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewConnection crea una nueva conexión a la base de datos PostgreSQL usando GORM
// Puede recibir parámetros directos o leer de variables de entorno
func NewConnection(host string, port int, user string, password string, dbname string) (*gorm.DB, error) {
	// Si se pasan parámetros vacíos, leer de variables de entorno
	if host == "" {
		host = getEnv("DB_HOST", "localhost")
	}
	if port == 0 {
		port = getEnvAsInt("DB_PORT", 5432)
	}
	if user == "" {
		user = getEnv("DB_USER", "myticktick")
	}
	if password == "" {
		password = getEnv("DB_PASSWORD", "myticktick")
	}
	if dbname == "" {
		dbname = getEnv("DB_NAME", "myticktick")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// getEnv devuelve el valor de una variable de entorno o un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt devuelve el valor de una variable de entorno como int o un valor por defecto
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		fmt.Sscanf(value, "%d", &defaultValue)
	}
	return defaultValue
}
