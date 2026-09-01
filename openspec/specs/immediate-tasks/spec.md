# immediate-tasks Specification

## Purpose
Gestionar tareas inmediatas para próximos días con fechas límite y prioridades.

## Requirements

### Requirement: Crear tareas inmediatas
El sistema SHALL permitir crear tareas para cumplir en los próximos días.

#### Scenario: Crear tarea con fecha límite
- **WHEN** el usuario crea una tarea inmediata con fecha límite
- **THEN** el sistema la guarda con la fecha especificada

#### Scenario: Ver tareas próximas
- **WHEN** el usuario abre la sección de tareas inmediatas
- **THEN** el sistema muestra las tareas ordenadas por fecha límite

### Requirement: Priorización de tareas
El sistema SHALL permitir asignar prioridad a las tareas inmediatas.

#### Scenario: Asignar prioridad alta
- **WHEN** el usuario marca una tarea como alta prioridad
- **THEN** el sistema la muestra primero en la lista

#### Scenario: Ver tareas por prioridad
- **WHEN** el usuario filtra por prioridad
- **THEN** el sistema muestra solo las tareas con esa prioridad

### Requirement: Recordatorios de tareas próximas
El sistema SHALL alertar al usuario sobre tareas con fecha límite próxima.

#### Scenario: Recordatorio de tarea próxima
- **WHEN** una tarea tiene fecha límite en los próximos 24 horas
- **THEN** el sistema muestra un recordatorio visual

### Requirement: Interacción de lista
El sistema SHALL permitir completar y eliminar tareas de forma fluida, sin diálogos de confirmación, y con posibilidad de revertir.

#### Scenario: Completar con toggle
- **WHEN** el usuario toca el toggle de una tarea
- **THEN** el sistema la marca como completada de inmediato, sin pedir confirmación

#### Scenario: Completar o eliminar con swipe
- **WHEN** el usuario desliza (swipe) una tarea
- **THEN** el sistema la completa (swipe hacia una dirección) o la elimina (swipe hacia la otra), sin pedir confirmación

#### Scenario: Deshacer (undo)
- **WHEN** el usuario completa o elimina una tarea
- **THEN** el sistema muestra una notificación con acción "Deshacer" durante unos segundos que revierte la operación
