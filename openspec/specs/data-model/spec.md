# data-model Specification

## Purpose
Definir los modelos de datos para tareas mensuales, tareas inmediatas, trackers y historial.

## Requirements

### Requirement: Modelo de usuario
El sistema SHALL tener un modelo User con username único y password hasheado.

#### Scenario: Crear modelo User
- **WHEN** se define el modelo User
- **THEN** incluye campos: id, username (único), passwordHash, createdAt, updatedAt
- **THEN** no incluye email: la identificación es solo por username

### Requirement: Modelo de tarea mensual
El sistema SHALL tener un modelo MonthlyTask con campos para nombre, descripción, y configuración de repetición mensual.

#### Scenario: Crear modelo MonthlyTask
- **WHEN** se define el modelo MonthlyTask
- **THEN** incluye campos: id, userId, name, description, isActive, createdAt, updatedAt

### Requirement: Modelo de historial de tarea mensual
El sistema SHALL tener un modelo MonthlyTaskHistory para registrar el cumplimiento de cada tarea mensual por mes.

#### Scenario: Crear modelo MonthlyTaskHistory
- **WHEN** se define el modelo MonthlyTaskHistory
- **THEN** incluye campos: id, monthlyTaskId, userId, month, year, completed, completedAt, createdAt

### Requirement: Modelo de tarea inmediata
El sistema SHALL tener un modelo ImmediateTask con campos para nombre, descripción, fecha límite y prioridad.

#### Scenario: Crear modelo ImmediateTask
- **WHEN** se define el modelo ImmediateTask
- **THEN** incluye campos: id, name, description, dueDate, priority, isCompleted, createdAt, updatedAt

### Requirement: Modelo de tracker diario
El sistema SHALL tener un modelo DailyTracker con un límite unilateral (mínimo o máximo) y unidad.

#### Scenario: Crear modelo DailyTracker
- **WHEN** se define el modelo DailyTracker
- **THEN** incluye campos: id, name, limitValue, limitType (min | max), unit, isActive, createdAt, updatedAt
- **THEN** limitType min significa "cumple si el valor es >= limitValue" (ej. sueño); max significa "cumple si el valor es <= limitValue" (ej. peso)

### Requirement: Modelo de registro de tracker
El sistema SHALL tener un modelo TrackerEntry para registrar valores diarios de cada tracker.

#### Scenario: Crear modelo TrackerEntry
- **WHEN** se define el modelo TrackerEntry
- **THEN** incluye campos: id, trackerId, userId, value, entryDate, notes, isMet, deviation, createdAt
- **THEN** deviation es siempre >= 0: cuánto el valor superó el límite (0 si se cumplió)
