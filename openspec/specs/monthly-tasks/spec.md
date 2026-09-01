# monthly-tasks Specification

## Purpose
Gestionar tareas mensuales recurrentes que se activan automáticamente el día 1 de cada mes, con historial de cumplimiento.

## Requirements

### Requirement: Tareas mensuales recurrentes
El sistema SHALL permitir crear tareas que se repiten mensualmente y se activan automáticamente el día 1 de cada mes.

#### Scenario: Crear tarea mensual
- **WHEN** el usuario crea una tarea mensual como "Pagar cuota de hijo"
- **THEN** el sistema la activa automáticamente el día 1 de cada mes

#### Scenario: Ver lista de tareas mensuales
- **WHEN** el usuario abre la sección de tareas mensuales
- **THEN** el sistema muestra todas las tareas mensuales configuradas con su estado actual

### Requirement: Historial de cumplimiento mensual
El sistema SHALL mantener un historial de si cada tarea mensual fue cumplida en meses anteriores.

#### Scenario: Ver historial de tarea
- **WHEN** el usuario consulta una tarea mensual
- **THEN** el sistema muestra el historial de cumplimiento de los últimos meses

#### Scenario: Marcar tarea como cumplida
- **WHEN** el usuario marca una tarea mensual como cumplida (toggle o swipe, sin confirmación)
- **THEN** el sistema registra la fecha y el estado de cumplimiento, con opción de "Deshacer"

### Requirement: Estado de tareas mensuales
El sistema SHALL rastrear el estado de cada tarea mensual (pendiente, completada, vencida).

#### Scenario: Tarea pendiente
- **WHEN** llega el día 1 del mes y la tarea no está completada
- **THEN** el sistema la marca como pendiente

#### Scenario: Tarea completada
- **WHEN** el usuario completa la tarea antes del fin del mes
- **THEN** el sistema la marca como completada y la registra en el historial
