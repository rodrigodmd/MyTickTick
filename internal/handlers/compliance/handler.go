package compliance

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"myticktick/internal/core/services/compliance"
)

// Handler expone las métricas de cumplimiento (6.4).
type Handler struct {
	service *compliance.Service
}

// NewHandler crea una instancia del handler.
func NewHandler(service *compliance.Service) *Handler {
	return &Handler{service: service}
}

// Metrics maneja GET /api/metrics
//
// Query params opcionales:
//
//	month=1..12, year=YYYY  → fija el mes/año para métricas mensuales y
//	                           el rango de los trackers (default: mes en curso)
//	from=YYYY-MM-DD, to=YYYY-MM-DD → rango explícito para trackers
func (h *Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now().UTC()
	month, year := int(now.Month()), now.Year()

	if v := r.URL.Query().Get("month"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 12 {
			http.Error(w, `{"error":"invalid month (1-12)"}`, http.StatusBadRequest)
			return
		}
		month = n
	}
	if v := r.URL.Query().Get("year"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 9999 {
			http.Error(w, `{"error":"invalid year"}`, http.StatusBadRequest)
			return
		}
		year = n
	}

	// Rango para trackers: por defecto el mes/año indicado.
	from, to := monthBounds(month, year)
	if v := r.URL.Query().Get("from"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			http.Error(w, `{"error":"invalid from (YYYY-MM-DD)"}`, http.StatusBadRequest)
			return
		}
		from = t
	}
	if v := r.URL.Query().Get("to"); v != "" {
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			http.Error(w, `{"error":"invalid to (YYYY-MM-DD)"}`, http.StatusBadRequest)
			return
		}
		to = t
	}

	ctx := r.Context()
	monthly, err := h.service.MonthlyForPeriod(ctx, month, year)
	if err != nil {
		http.Error(w, `{"error":"failed to compute monthly metrics"}`, http.StatusInternalServerError)
		return
	}
	series, err := h.service.MonthlySeries(ctx, year)
	if err != nil {
		http.Error(w, `{"error":"failed to compute monthly series"}`, http.StatusInternalServerError)
		return
	}
	trackers, err := h.service.TrackersInRange(ctx, from, to)
	if err != nil {
		http.Error(w, `{"error":"failed to compute tracker metrics"}`, http.StatusInternalServerError)
		return
	}

	resp := compliance.Overview{
		Monthly:       *monthly,
		MonthlySeries: series,
		Trackers:      trackers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// monthBounds devuelve inicio (00:00) y fin (23:59:59) del mes/año en UTC.
func monthBounds(month, year int) (time.Time, time.Time) {
	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	to := time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 23, 59, 59, 0, time.UTC)
	return from, to
}
