package trackers

import (
	"time"

	"myticktick/internal/core/domain"
)

// CreateRequest representa la solicitud para crear un tracker
type CreateRequest struct {
	Name       string  `json:"name"`
	LimitValue float64 `json:"limitValue"`
	LimitType  string  `json:"limitType"` // "min" o "max"
	Unit       string  `json:"unit,omitempty"`
}

// UpdateRequest representa la solicitud para actualizar un tracker
type UpdateRequest struct {
	Name       string  `json:"name"`
	LimitValue float64 `json:"limitValue"`
	LimitType  string  `json:"limitType"` // "min" o "max"
	Unit       string  `json:"unit,omitempty"`
	IsActive   bool    `json:"isActive"`
}

// RecordRequest representa la solicitud para crear un registro de tracker
type RecordRequest struct {
	Value float64 `json:"value"`
	Notes string  `json:"notes,omitempty"`
	Date  string  `json:"date"` // YYYY-MM-DD
}

// TrackerResponse representa la respuesta de un tracker
type TrackerResponse struct {
	ID         uint    `json:"id"`
	UserID     uint    `json:"userId"`
	Name       string  `json:"name"`
	LimitValue float64 `json:"limitValue"`
	LimitType  string  `json:"limitType"`
	Unit       string  `json:"unit"`
	IsActive   bool    `json:"isActive"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

// TrackerEntryResponse representa la respuesta de una entrada de tracker
type TrackerEntryResponse struct {
	ID        uint    `json:"id"`
	TrackerID uint    `json:"trackerId"`
	Value     float64 `json:"value"`
	EntryDate string  `json:"entryDate"` // YYYY-MM-DD
	Notes     string  `json:"notes"`
	IsMet     bool    `json:"isMet"`
	Deviation float64 `json:"deviation"`
	CreatedAt string  `json:"createdAt"`
}

// ToResponse convierte una entidad DailyTracker a TrackerResponse
func ToResponse(tracker *domain.DailyTracker) *TrackerResponse {
	return &TrackerResponse{
		ID:         tracker.ID,
		UserID:     tracker.UserID,
		Name:       tracker.Name,
		LimitValue: tracker.LimitValue,
		LimitType:  tracker.LimitType,
		Unit:       tracker.Unit,
		IsActive:   tracker.IsActive,
		CreatedAt:  tracker.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  tracker.UpdatedAt.Format(time.RFC3339),
	}
}

// ToEntryResponse convierte una entidad TrackerEntry a TrackerEntryResponse
func ToEntryResponse(entry *domain.TrackerEntry) *TrackerEntryResponse {
	return &TrackerEntryResponse{
		ID:        entry.ID,
		TrackerID: entry.TrackerID,
		Value:     entry.Value,
		EntryDate: entry.EntryDate.UTC().Format("2006-01-02"),
		Notes:     entry.Notes,
		IsMet:     entry.IsMet,
		Deviation: entry.Deviation,
		CreatedAt: entry.CreatedAt.Format(time.RFC3339),
	}
}

// FromCreateRequest crea una entidad DailyTracker desde CreateRequest
func FromCreateRequest(req CreateRequest) *domain.DailyTracker {
	return &domain.DailyTracker{
		UserID:     0, // app single-user por ahora
		Name:       req.Name,
		LimitValue: req.LimitValue,
		LimitType:  normalizeLimitType(req.LimitType),
		Unit:       req.Unit,
		IsActive:   true,
	}
}

// FromUpdateRequest crea una entidad DailyTracker desde UpdateRequest
func FromUpdateRequest(id uint, req UpdateRequest) *domain.DailyTracker {
	return &domain.DailyTracker{
		ID:         id,
		UserID:     0, // app single-user por ahora
		Name:       req.Name,
		LimitValue: req.LimitValue,
		LimitType:  normalizeLimitType(req.LimitType),
		Unit:       req.Unit,
		IsActive:   req.IsActive,
	}
}

// FromRecordRequest crea una entidad TrackerEntry desde RecordRequest
func FromRecordRequest(trackerID uint, req RecordRequest) *domain.TrackerEntry {
	entryDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		// Fallback a hoy (midnight UTC para evitar el drift de zona horaria)
		now := time.Now()
		entryDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	}

	return &domain.TrackerEntry{
		TrackerID: trackerID,
		UserID:    0, // app single-user por ahora
		Value:     req.Value,
		EntryDate: entryDate,
		Notes:     req.Notes,
	}
}

// normalizeLimitType normaliza el tipo de límite a "min" o "max" (default: max).
func normalizeLimitType(s string) string {
	switch s {
	case "min":
		return "min"
	case "max":
		return "max"
	default:
		return "max"
	}
}
