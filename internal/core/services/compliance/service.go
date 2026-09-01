package compliance

import (
	"context"
	"time"

	"myticktick/internal/core/ports"
)

// MonthlyMetrics resume la tasa de cumplimiento de tareas mensuales
// para un mes/año dado (spec compliance-history / "tasa de cumplimiento
// de tareas mensuales: porcentaje de tareas completadas vs totales").
type MonthlyMetrics struct {
	Month     int     `json:"month"`
	Year      int     `json:"year"`
	Total     int     `json:"total"`
	Completed int     `json:"completed"`
	Rate      float64 `json:"rate"` // 0.0 - 1.0
}

// TrackerMetrics resume el cumplimiento de un tracker dentro de un rango
// de fechas (spec compliance-history / "tasa de cumplimiento de un
// tracker: porcentaje de días cumplidos vs totales").
type TrackerMetrics struct {
	TrackerID       uint    `json:"trackerId"`
	Name            string  `json:"name"`
	Total           int     `json:"total"`
	Met             int     `json:"met"`
	Rate            float64 `json:"rate"`            // 0.0 - 1.0
	AvgDeviation    float64 `json:"avgDeviation"`    // desviación media (con signo)
	AvgAbsDeviation float64 `json:"avgAbsDeviation"` // desviación media (sin signo)
}

// Overview es la respuesta agregada de GET /api/metrics.
type Overview struct {
	Monthly       MonthlyMetrics   `json:"monthly"`
	MonthlySeries []MonthlyMetrics `json:"monthlySeries"`
	Trackers      []TrackerMetrics `json:"trackers"`
}

// Service agrega métricas de cumplimiento sobre tareas mensuales y trackers.
type Service struct {
	taskRepo    ports.MonthlyTaskRepository
	historyRepo ports.MonthlyTaskHistoryRepository
	trackerRepo ports.TrackerRepository
	entryRepo   ports.TrackerEntryRepository
}

// NewService crea una instancia de Service de métricas.
func NewService(
	taskRepo ports.MonthlyTaskRepository,
	historyRepo ports.MonthlyTaskHistoryRepository,
	trackerRepo ports.TrackerRepository,
	entryRepo ports.TrackerEntryRepository,
) *Service {
	return &Service{
		taskRepo:    taskRepo,
		historyRepo: historyRepo,
		trackerRepo: trackerRepo,
		entryRepo:   entryRepo,
	}
}

// MonthlyForPeriod calcula la tasa de cumplimiento de un mes/año:
// registros de historial de todas las tareas vs los completados.
func (s *Service) MonthlyForPeriod(ctx context.Context, month int, year int) (*MonthlyMetrics, error) {
	history, err := s.historyRepo.FindByPeriod(ctx, month, year)
	if err != nil {
		return nil, err
	}

	total := len(history)
	completed := 0
	for _, h := range history {
		if h.Completed {
			completed++
		}
	}

	return &MonthlyMetrics{
		Month:     month,
		Year:      year,
		Total:     total,
		Completed: completed,
		Rate:      ratio(completed, total),
	}, nil
}

// MonthlySeries calcula la tasa de cumplimiento mes a mes para un año
// (12 meses). Se usa en el gráfico de barras del dashboard (11.3).
func (s *Service) MonthlySeries(ctx context.Context, year int) ([]MonthlyMetrics, error) {
	series := make([]MonthlyMetrics, 0, 12)
	for month := 1; month <= 12; month++ {
		m, err := s.MonthlyForPeriod(ctx, month, year)
		if err != nil {
			return nil, err
		}
		series = append(series, *m)
	}
	return series, nil
}

// TrackersInRange calcula el cumplimiento de cada tracker dentro de un
// rango de fechas (inclusive).
func (s *Service) TrackersInRange(ctx context.Context, from, to time.Time) ([]TrackerMetrics, error) {
	trackers, err := s.trackerRepo.FindByUserID(ctx, 0)
	if err != nil {
		return nil, err
	}

	metrics := make([]TrackerMetrics, 0, len(trackers))
	for _, tracker := range trackers {
		entries, err := s.entryRepo.FindByTrackerAndRange(ctx, tracker.ID, from, to)
		if err != nil {
			return nil, err
		}

		total := len(entries)
		met := 0
		sumDev, sumAbs := 0.0, 0.0
		for _, e := range entries {
			if e.IsMet {
				met++
			}
			sumDev += e.Deviation
			sumAbs += absFloat(e.Deviation)
		}

		metrics = append(metrics, TrackerMetrics{
			TrackerID:       tracker.ID,
			Name:            tracker.Name,
			Total:           total,
			Met:             met,
			Rate:            ratio(met, total),
			AvgDeviation:    safeAvg(sumDev, total),
			AvgAbsDeviation: safeAvg(sumAbs, total),
		})
	}
	return metrics, nil
}

// Overview arma la respuesta agregada para el mes/año y año indicados.
func (s *Service) Overview(ctx context.Context, now time.Time) (*Overview, error) {
	month := int(now.Month())
	year := now.Year()

	monthly, err := s.MonthlyForPeriod(ctx, month, year)
	if err != nil {
		return nil, err
	}
	series, err := s.MonthlySeries(ctx, year)
	if err != nil {
		return nil, err
	}
	from, to := monthBounds(now)
	trackers, err := s.TrackersInRange(ctx, from, to)
	if err != nil {
		return nil, err
	}

	return &Overview{
		Monthly:       *monthly,
		MonthlySeries: series,
		Trackers:      trackers,
	}, nil
}

// ratio devuelve a/b, o 0 si b es 0.
func ratio(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

// safeAvg devuelve la media de un acumulado sobre n elementos, o 0 si n es 0.
func safeAvg(sum float64, n int) float64 {
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}

// absFloat devuelve el valor absoluto.
func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// monthBounds devuelve el inicio y el fin (23:59:59) del mes de t en UTC.
func monthBounds(t time.Time) (time.Time, time.Time) {
	y, m, d := t.Year(), t.Month(), t.Day()
	from := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
	return from, to
}
