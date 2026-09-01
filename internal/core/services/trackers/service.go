package trackers

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"myticktick/internal/core/domain"
	"myticktick/internal/core/ports"
)

// Service maneja la lógica de negocio de trackers
type Service struct {
	trackerRepo ports.TrackerRepository
	entryRepo   ports.TrackerEntryRepository
}

// NewService crea una nueva instancia de Service
func NewService(trackerRepo ports.TrackerRepository, entryRepo ports.TrackerEntryRepository) *Service {
	return &Service{trackerRepo: trackerRepo, entryRepo: entryRepo}
}

// Create crea un nuevo tracker
func (s *Service) Create(ctx context.Context, tracker *domain.DailyTracker) (*domain.DailyTracker, error) {
	err := s.trackerRepo.Create(ctx, tracker)
	if err != nil {
		return nil, err
	}
	return tracker, nil
}

// List obtiene todos los trackers de un usuario
func (s *Service) List(ctx context.Context, userID uint) ([]*domain.DailyTracker, error) {
	return s.trackerRepo.FindByUserID(ctx, userID)
}

// Get obtiene un tracker por ID
func (s *Service) Get(ctx context.Context, id uint) (*domain.DailyTracker, error) {
	return s.trackerRepo.FindByID(ctx, id)
}

// Update actualiza un tracker
func (s *Service) Update(ctx context.Context, tracker *domain.DailyTracker) (*domain.DailyTracker, error) {
	err := s.trackerRepo.Update(ctx, tracker)
	if err != nil {
		return nil, err
	}
	return tracker, nil
}

// Delete elimina un tracker
func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.trackerRepo.Delete(ctx, id)
}

// CreateEntry crea una nueva entrada en el tracker y calcula el cumplimiento
// contra el límite unilateral (mínimo o máximo).
func (s *Service) CreateEntry(ctx context.Context, entry *domain.TrackerEntry) (*domain.TrackerEntry, error) {
	// Obtener el tracker para verificar el límite
	tracker, err := s.trackerRepo.FindByID(ctx, entry.TrackerID)
	if err != nil {
		return nil, err
	}

	// Calcular si se cumplió el límite unilateral:
	// - max: cumple si value <= LimitValue (ej. peso, no superar el máximo)
	// - min: cumple si value >= LimitValue (ej. sueño, no bajar del mínimo)
	entry.IsMet, entry.Deviation = EvaluateLimit(tracker.LimitType, tracker.LimitValue, entry.Value)

	err = s.entryRepo.Create(ctx, entry)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// UpsertEntry crea o actualiza la entrada de un tracker para una fecha.
// Permite "re-setear el valor del día": si ya existe un registro para esa
// fecha se actualiza (value, notes, is_met, deviation) y no se duplica.
func (s *Service) UpsertEntry(ctx context.Context, entry *domain.TrackerEntry) (*domain.TrackerEntry, error) {
	// Obtener el tracker para verificar el límite
	tracker, err := s.trackerRepo.FindByID(ctx, entry.TrackerID)
	if err != nil {
		return nil, err
	}

	entry.IsMet, entry.Deviation = EvaluateLimit(tracker.LimitType, tracker.LimitValue, entry.Value)

	// ¿Ya hay un registro para esa fecha?
	existing, err := s.entryRepo.FindByTrackerAndDate(ctx, entry.TrackerID, entry.EntryDate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe: crearlo.
			return s.CreateEntry(ctx, entry)
		}
		return nil, err
	}

	// Existe: actualizarlo (mismo ID, mismo tracker, nueva fecha opcional).
	existing.Value = entry.Value
	existing.EntryDate = entry.EntryDate
	existing.Notes = entry.Notes
	existing.IsMet, existing.Deviation = entry.IsMet, entry.Deviation

	if err := s.entryRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// EvaluateLimit evalúa un valor contra el límite unilateral de un tracker.
// Devuelve (isMet, deviation) donde deviation es siempre >= 0
// (0 si se cumplió, o cuánto se superó el límite si no).
// Si limitType es inválido o vacío, se asume "max".
func EvaluateLimit(limitType string, limitValue, value float64) (isMet bool, deviation float64) {
	switch limitType {
	case "min":
		if value >= limitValue {
			return true, 0
		}
		return false, limitValue - value
	case "max", "":
		if value <= limitValue {
			return true, 0
		}
		return false, value - limitValue
	default:
		// Tipo desconocido: tratar como max para no perder datos
		if value <= limitValue {
			return true, 0
		}
		return false, value - limitValue
	}
}

// GetHistory obtiene el historial de un tracker
func (s *Service) GetHistory(ctx context.Context, trackerID uint) ([]*domain.TrackerEntry, error) {
	return s.entryRepo.FindByTrackerID(ctx, trackerID)
}

// CalculateComplianceMetrics calcula métricas agregadas de cumplimiento
func (s *Service) CalculateComplianceMetrics(ctx context.Context, trackerID uint) (float64, error) {
	entries, err := s.entryRepo.FindByTrackerID(ctx, trackerID)
	if err != nil {
		return 0, err
	}

	if len(entries) == 0 {
		return 0, nil
	}

	completed := 0
	for _, entry := range entries {
		if entry.IsMet {
			completed++
		}
	}

	return float64(completed) / float64(len(entries)), nil
}
