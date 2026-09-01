package migrations

import (
	"log"

	"gorm.io/gorm"
	"myticktick/internal/repositories/user"
	"myticktick/internal/repositories/monthlytask"
	"myticktick/internal/repositories/monthlytaskhistory"
	"myticktick/internal/repositories/immediatetask"
	"myticktick/internal/repositories/dailytracker"
	"myticktick/internal/repositories/trackerentry"
)

// RunMigrations ejecuta las migraciones de base de datos
// Crea todas las tablas necesarias si no existen
func RunMigrations(db *gorm.DB) {
	log.Println("Running database migrations...")

	// Migración manual idempotente: rename del modelo de thresholds
	// (target_value + threshold) a límite unilateral (limit_value + limit_type).
	// AutoMigrate no renombra columnas ni agrega NOT NULL con datos existentes,
	// así que hacemos esta migración explícita antes del AutoMigrate.
	migrateTrackersToLimitModel(db)

	// AutoMigrate crea las tablas si no existen
	// y agrega columnas nuevas si se agregan a los modelos
	db.AutoMigrate(
		&user.UserDB{},
		&monthlytask.MonthlyTaskDB{},
		&monthlytaskhistory.MonthlyTaskHistoryDB{},
		&immediatetask.ImmediateTaskDB{},
		&dailytracker.DailyTrackerDB{},
		&trackerentry.TrackerEntryDB{},
	)

	log.Println("Database migrations completed successfully")
}

// hasColumn indica si la tabla tiene una columna dada.
func hasColumn(db *gorm.DB, table, column string) bool {
	var count int
	db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_name = ? AND column_name = ?", table, column).Scan(&count)
	return count > 0
}

// migrateTrackersToLimitModel migra daily_tracker_dbs del modelo de thresholds
// simétricos (target_value, threshold) al modelo de límite unilateral
// (limit_value, limit_type). Idempotente: si ya está migrada, no hace nada.
func migrateTrackersToLimitModel(db *gorm.DB) {
	const table = "daily_tracker_dbs"

	// Si ya existe limit_value, la migración ya corrió; solo limpiar columnas viejas.
	if hasColumn(db, table, "limit_value") {
		// Asegurar default en limit_type por si el AutoMigrate anterior no lo puso
		db.Exec(`ALTER TABLE "` + table + `" ALTER COLUMN "limit_type" SET DEFAULT 'max'`)
		dropIfExists(db, table, "target_value")
		dropIfExists(db, table, "threshold")
		return
	}

	if !hasColumn(db, table, "target_value") {
		// Tabla recién creada (o ya sin columnas viejas): no hay nada que migrar.
		return
	}

	log.Println("Migrando daily_tracker_dbs a modelo de límite unilateral...")

	// 1) Si un intento de migración anterior dejó limit_value/limit_type a medias,
	//    limpiarlos para poder renombrar sin conflictos.
	dropIfExists(db, table, "limit_value")
	dropIfExists(db, table, "limit_type")

	// 2) Renombrar target_value -> limit_value y threshold -> limit_type.
	db.Exec(`ALTER TABLE "` + table + `" RENAME COLUMN "target_value" TO "limit_value"`)
	db.Exec(`ALTER TABLE "` + table + `" RENAME COLUMN "threshold" TO "limit_type"`)

	// 3) limit_type ahora contiene el viejo threshold (numérico) -> asignar 'max'
	//    a todos los trackers existentes (conservación: antes se penalizaba en
	//    ambas direcciones; ahora asumimos "máximo" para peso y el usuario puede
	//    cambiarlo a "mínimo" para sueño).
	db.Exec(`ALTER TABLE "` + table + `" ALTER COLUMN "limit_type" TYPE text USING 'max'`)
	db.Exec(`ALTER TABLE "` + table + `" ALTER COLUMN "limit_type" SET DEFAULT 'max'`)
	db.Exec(`UPDATE "` + table + `" SET "limit_type" = 'max' WHERE "limit_type" IS NULL OR "limit_type" = ''`)

	// 4) Asegurar NOT NULL con valores válidos.
	db.Exec(`UPDATE "` + table + `" SET "limit_value" = 0 WHERE "limit_value" IS NULL`)
	db.Exec(`ALTER TABLE "` + table + `" ALTER COLUMN "limit_value" SET NOT NULL`)
	db.Exec(`ALTER TABLE "` + table + `" ALTER COLUMN "limit_type" SET NOT NULL`)

	log.Println("Migración de trackers a límite unilateral completada")
}

// dropIfExists elimina una columna si existe.
func dropIfExists(db *gorm.DB, table, column string) {
	if hasColumn(db, table, column) {
		db.Exec(`ALTER TABLE "` + table + `" DROP COLUMN "` + column + `"`)
	}
}
