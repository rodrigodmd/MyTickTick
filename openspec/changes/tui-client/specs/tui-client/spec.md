# tui-client Specification

## Purpose

Proveer un cliente de terminal interactivo para operar MyTickTick (tareas, mensuales, trackers, métricas) contra la API REST existente, sin abrir browser y con sesión persistente local.

## ADDED Requirements

### Requirement: Conexión a la API y configuración
El cliente TUI SHALL comunicarse exclusivamente con la API REST de MyTickTick sobre HTTP, utilizando una base URL configurable. No SHALL acceder directamente a la base de datos ni reimplementar lógica de negocio del servidor.

La base URL SHALL ser configurable mediante la variable de entorno `MTT_URL`, con valor por defecto `http://localhost:8080`.

#### Scenario: Conectar con el valor por defecto
- **WHEN** el usuario corre el TUI sin definir `MTT_URL`
- **THEN** el cliente usa `http://localhost:8080` como base URL

#### Scenario: Conectar con base URL personalizada
- **WHEN** el usuario define `MTT_URL=http://10.0.0.5:8080` antes de correr el TUI
- **THEN** el cliente usa esa base URL para todas sus peticiones

#### Scenario: Servidor no disponible
- **WHEN** el cliente no puede conectar al servidor indicado
- **THEN** muestra un error claro indicando que no hay un MyTickTick server accesible en la URL configurada, sin entrar en reintento infinito silencioso

### Requirement: Autenticación y sesión persistente
El cliente TUI SHALL autenticarse al servidor con usuario y contraseña, y SHALL persistir la sesión localmente para no exigir el login en cada arranque. Cuando la sesión ya no sea válida, el cliente SHALL solicitar el login nuevamente.

#### Scenario: Login inicial
- **WHEN** el usuario arranca el TUI sin sesión guardada
- **THEN** el cliente solicita usuario y contraseña, envía el login al servidor y, si es correcto, guarda la sesión localmente para usos futuros

#### Scenario: Login con credenciales inválidas
- **WHEN** el servidor rechaza las credenciales
- **THEN** el cliente muestra un error de credenciales inválidas y permite reintentar sin cerrar el TUI

#### Scenario: Sesión reutilizada en el próximo arranque
- **WHEN** el usuario cierra el TUI con una sesión válida guardada y vuelve a abrirlo
- **THEN** el cliente entra directamente a la primera pantalla sin pedir credenciales de nuevo

#### Scenario: Sesión expirada o inválida
- **WHEN** el servidor rechaza una petición por sesión inválida o expirada
- **THEN** el cliente descarta la sesión guardada y solicita un login nuevo

### Requirement: Gestión de tareas inmediatas
El cliente TUI SHALL permitir listar, crear, editar y eliminar tareas inmediatas, y togglear su estado de completada, operando sobre la API existente.

#### Scenario: Listar tareas inmediatas
- **WHEN** el usuario abre la pantalla de tareas
- **THEN** el cliente muestra todas las tareas inmediatas con su estado de completada, prioridad y fecha de vencimiento

#### Scenario: Crear una tarea inmediata
- **WHEN** el usuario crea una tarea indicando su nombre
- **THEN** el cliente la crea en el servidor y la muestra en la lista actualizada

#### Scenario: Completar o reabrir una tarea
- **WHEN** el usuario togglea el estado de una tarea
- **THEN** el cliente actualiza su estado de completada en el servidor y lo refleja en la lista

#### Scenario: Eliminar una tarea
- **WHEN** el usuario elimina una tarea
- **THEN** el cliente la elimina en el servidor y desaparece de la lista

### Requirement: Gestión de tareas mensuales
El cliente TUI SHALL permitir listar tareas mensuales, togglear su cumplimiento del período actual, ver su historial y activarlas, operando sobre la API existente.

#### Scenario: Listar tareas mensuales
- **WHEN** el usuario abre la pantalla de mensuales
- **THEN** el cliente muestra cada tarea mensual con su estado de cumplimiento del mes en curso

#### Scenario: Togglear cumplimiento de una mensual
- **WHEN** el usuario marca una mensual como cumplida o la revierte
- **THEN** el cliente registra el cambio en el servidor y actualiza el estado mostrado

#### Scenario: Ver historial de una mensual
- **WHEN** el usuario consulta el historial de una tarea mensual
- **THEN** el cliente muestra su registro de cumplimiento mes a mes

### Requirement: Gestión de trackers
El cliente TUI SHALL permitir listar trackers, cargar o corregir el valor del día, y ver el historial de cada tracker, operando sobre la API existente.

#### Scenario: Listar trackers
- **WHEN** el usuario abre la pantalla de trackers
- **THEN** el cliente muestra cada tracker con su límite, tipo (mín/máx), unidad y el valor cargado para el día actual

#### Scenario: Cargar el valor del día
- **WHEN** el usuario ingresa un valor para un tracker
- **THEN** el cliente lo registra en el servidor para la fecha de hoy y muestra el resultado con su cumplimiento respecto al límite

#### Scenario: Corregir el valor del día
- **WHEN** el usuario ingresa de nuevo un valor para un tracker que ya tiene registro hoy
- **THEN** el cliente reemplaza el valor del día en el servidor en lugar de duplicarlo

#### Scenario: Ver historial de un tracker
- **WHEN** el usuario consulta el historial de un tracker
- **THEN** el cliente muestra los valores registrados con sus fechas y su desviación respecto al límite

### Requirement: Visualización de métricas
El cliente TUI SHALL mostrar las métricas de cumplimiento del sistema: cumplimiento de tareas mensuales por mes y año, la serie mensual de cumplimiento, y las métricas de trackers en un rango.

#### Scenario: Ver métricas del mes en curso
- **WHEN** el usuario abre la pantalla de métricas sin filtro
- **THEN** el cliente muestra el cumplimiento de mensuales del mes en curso, la serie mensual y las métricas de trackers del período

#### Scenario: Filtrar métricas por mes y año
- **WHEN** el usuario selecciona un mes y año específicos
- **THEN** el cliente consulta las métricas para ese período y las muestra

### Requirement: Manejo de errores de operación
El cliente TUI SHALL reportar de forma legible los errores devueltos por el servidor y los fallos de red, permitiendo al usuario continuar operando sin perder su estado actual.

#### Scenario: Error de validación del servidor
- **WHEN** el servidor rechaza una operación por datos inválidos
- **THEN** el cliente muestra el motivo del rechazo y mantiene la vista actual

#### Scenario: Fallo de red durante una operación
- **WHEN** una petición de escritura falla por red
- **THEN** el cliente muestra el error de conexión y permite reintentar la operación
