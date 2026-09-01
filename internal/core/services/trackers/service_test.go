package trackers

import (
	"testing"
)

// TestEvaluateLimit cubre el cálculo de límite unilateral (min y max)
// con desviación siempre >= 0.
func TestEvaluateLimit(t *testing.T) {
	tests := []struct {
		name        string
		limitType   string
		limitValue  float64
		value       float64
		wantIsMet   bool
		wantDev     float64
	}{
		// Límite máximo: peso, no superar el máximo
		{"max: bajo límite -> cumple, dev 0", "max", 85, 80, true, 0},
		{"max: en límite -> cumple, dev 0", "max", 85, 85, true, 0},
		{"max: sobre límite -> no cumple, dev = excedente", "max", 85, 90, false, 5},
		{"max: muy sobre límite -> dev proporcional", "max", 85, 100, false, 15},

		// Límite mínimo: sueño, no bajar del mínimo
		{"min: sobre mínimo -> cumple, dev 0", "min", 6, 8, true, 0},
		{"min: en mínimo -> cumple, dev 0", "min", 6, 6, true, 0},
		{"min: bajo mínimo -> no cumple, dev = déficit", "min", 6, 4, false, 2},
		{"min: muy bajo mínimo -> dev proporcional", "min", 6, 1, false, 5},

		// Tipo vacío o inválido: se asume max
		{"vacío: sobre límite -> no cumple (asume max)", "", 85, 90, false, 5},
		{"vacío: bajo límite -> cumple (asume max)", "", 85, 80, true, 0},
		{"inválido: se trata como max", "weird", 85, 90, false, 5},

		// Decimales
		{"max decimal: 85.5 vs 86", "max", 85.5, 86, false, 0.5},
		{"min decimal: 6.5 vs 6", "min", 6.5, 6, false, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isMet, dev := EvaluateLimit(tt.limitType, tt.limitValue, tt.value)
			if isMet != tt.wantIsMet {
				t.Errorf("IsMet = %v, want %v", isMet, tt.wantIsMet)
			}
			if dev != tt.wantDev {
				t.Errorf("Deviation = %v, want %v", dev, tt.wantDev)
			}
			if dev < 0 {
				t.Errorf("Deviation nunca debe ser negativa, got %v", dev)
			}
		})
	}
}

// TestEvaluateLimit_SymmetryVerif verifica que el comportamiento anterior
// (target ± threshold simétrico) ya NO se aplica: un valor por debajo del
// límite máximo también "cumple" (solo se penaliza por exceder el límite).
func TestEvaluateLimit_UnilateralBehavior(t *testing.T) {
	// Con el modelo anterior: target=85, threshold=5 -> 80 NO cumplía.
	// Con el modelo unilateral (max 85): 80 SÍ cumple.
	isMet, dev := EvaluateLimit("max", 85, 80)
	if !isMet {
		t.Error("con límite unilateral max, un valor bajo el límite debe cumplir")
	}
	if dev != 0 {
		t.Errorf("deviación debe ser 0 cuando se cumple, got %v", dev)
	}

	// Con el modelo anterior: target=6, threshold=1 -> 4 cumplía (dentro del rango).
	// Con el modelo unilateral (min 6): 4 NO cumple (bajo el mínimo).
	isMet, dev = EvaluateLimit("min", 6, 4)
	if isMet {
		t.Error("con límite unilateral min, un valor bajo el mínimo no debe cumplir")
	}
	if dev != 2 {
		t.Errorf("desviación debe ser 2 (6-4), got %v", dev)
	}
}
