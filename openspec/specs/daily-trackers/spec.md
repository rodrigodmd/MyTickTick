# daily-trackers Specification

## Purpose
Permitir tracking diario de métricas personales con un límite unilateral (mínimo o máximo) y registro de desviación.

## Requirements

### Requirement: Crear trackers diarios
El sistema SHALL permitir crear trackers para métricas diarias como peso, sueño, agua, etc., cada uno con un límite unilateral.

#### Scenario: Crear tracker de peso (límite máximo)
- **WHEN** el usuario crea un tracker de peso con límite máximo 85kg
- **THEN** el sistema lo guarda como tracker diario con limitType=max y limitValue=85

#### Scenario: Crear tracker de sueño (límite mínimo)
- **WHEN** el usuario crea un tracker de sueño con límite mínimo 6h
- **THEN** el sistema lo guarda como tracker diario con limitType=min y limitValue=6

### Requirement: Límite unilateral
El sistema SHALL evaluar el cumplimiento contra un solo límite (mínimo o máximo), no contra un rango simétrico.

#### Scenario: Límite máximo de peso
- **WHEN** el tracker de peso tiene límite máximo 85kg
- **THEN** los valores <= 85kg se marcan como cumplidos y los > 85kg como no cumplidos
- **THEN** estar por debajo del límite (ej. 80kg) siempre se considera cumplido

#### Scenario: Límite mínimo de sueño
- **WHEN** el tracker de sueño tiene límite mínimo 6h
- **THEN** los valores >= 6h se marcan como cumplidos y los < 6h como no cumplidos
- **THEN** dormir más del límite (ej. 9h) siempre se considera cumplido

### Requirement: Registro de desviación
El sistema SHALL registrar cuánto el valor real superó el límite cuando no se cumple. La desviación es siempre >= 0 (0 si se cumplió).

#### Scenario: Peso por encima del límite máximo
- **WHEN** el usuario registra 87kg (límite máximo 85kg)
- **THEN** el sistema marca como no cumplido y registra desviación de 2kg

#### Scenario: Sueño por debajo del límite mínimo
- **WHEN** el usuario registra 5h de sueño (límite mínimo 6h)
- **THEN** el sistema marca como no cumplido y registra desviación de 1h

#### Scenario: Valor dentro del cumplimiento
- **WHEN** el usuario registra 84kg (límite máximo 85kg) o 7h (límite mínimo 6h)
- **THEN** el sistema marca como cumplido y registra desviación de 0

### Requirement: Historial de trackers
El sistema SHALL mantener historial de todos los registros diarios de cada tracker.

#### Scenario: Ver historial de tracker
- **WHEN** el usuario consulta el historial de un tracker
- **THEN** el sistema muestra todos los registros anteriores con fechas
