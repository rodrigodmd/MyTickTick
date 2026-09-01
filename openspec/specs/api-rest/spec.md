# api-rest Specification

## Purpose
Proveer API REST para comunicación entre frontend y backend.

## Requirements

### Requirement: Endpoints de tareas mensuales
El sistema SHALL proveer endpoints REST para CRUD de tareas mensuales.

#### Scenario: Crear tarea mensual
- **WHEN** se hace POST /api/monthly-tasks con datos de tarea
- **THEN** el sistema crea la tarea y devuelve el recurso creado

#### Scenario: Listar tareas mensuales
- **WHEN** se hace GET /api/monthly-tasks
- **THEN** el sistema devuelve todas las tareas mensuales

#### Scenario: Actualizar tarea mensual
- **WHEN** se hace PUT /api/monthly-tasks/{id} con nuevos datos
- **THEN** el sistema actualiza la tarea y devuelve el recurso actualizado

#### Scenario: Eliminar tarea mensual
- **WHEN** se hace DELETE /api/monthly-tasks/{id}
- **THEN** el sistema elimina la tarea

### Requirement: Endpoints de tareas inmediatas
El sistema SHALL proveer endpoints REST para CRUD de tareas inmediatas.

#### Scenario: Crear tarea inmediata
- **WHEN** se hace POST /api/immediate-tasks con datos de tarea
- **THEN** el sistema crea la tarea y devuelve el recurso creado

#### Scenario: Listar tareas inmediatas
- **WHEN** se hace GET /api/immediate-tasks
- **THEN** el sistema devuelve todas las tareas inmediatas ordenadas por fecha límite

### Requirement: Endpoints de trackers
El sistema SHALL proveer endpoints REST para CRUD de trackers diarios.

#### Scenario: Crear tracker
- **WHEN** se hace POST /api/trackers con datos de tracker
- **THEN** el sistema crea el tracker y devuelve el recurso creado

#### Scenario: Registrar valor diario
- **WHEN** se hace POST /api/trackers/{id}/records con valor y fecha
- **THEN** el sistema guarda el registro y calcula si se cumplió el límite unilateral (mínimo o máximo) del tracker

### Requirement: Endpoints de autenticación
El sistema SHALL proveer endpoints de registro, login y logout basados en username y password.

#### Scenario: Registrar usuario
- **WHEN** se hace POST /api/register con username único y password
- **THEN** el sistema almacena el usuario con password hasheado y devuelve confirmación

#### Scenario: Iniciar sesión
- **WHEN** se hace POST /api/login con username y password correctos
- **THEN** el sistema emite un JWT y lo entrega como cookie httpOnly (con expiración larga si se pidió "recordarme", corta en caso contrario)

#### Scenario: Iniciar sesión con credenciales inválidas
- **WHEN** se hace POST /api/login con username o password incorrectos
- **THEN** el sistema devuelve 401 y no emite cookie

#### Scenario: Cerrar sesión
- **WHEN** se hace POST /api/logout
- **THEN** el sistema invalida la sesión y borra la cookie de token

### Requirement: Protección de rutas
El sistema SHALL requerir una sesión válida (JWT) para acceder a todos los endpoints de tareas, trackers y métricas.

#### Scenario: Petición sin sesión
- **WHEN** se hace una petición a /api/monthly-tasks, /api/immediate-tasks, /api/trackers o /api/metrics sin cookie válida
- **THEN** el sistema devuelve 401 y no expone datos

#### Scenario: Petición con sesión válida
- **WHEN** se hace una petición a esos endpoints con un JWT válido en la cookie
- **THEN** el sistema procesa la petición normalmente

#### Scenario: Endpoints públicos
- **WHEN** se hace una petición a /api/health, /api/register o /api/login
- **THEN** el sistema la procesa sin requerir sesión

### Requirement: Endpoints de historial
El sistema SHALL proveer endpoints REST para consultar historial.

#### Scenario: Obtener historial de tarea mensual
- **WHEN** se hace GET /api/monthly-tasks/{id}/history
- **THEN** el sistema devuelve el historial de cumplimiento de esa tarea

#### Scenario: Obtener historial de tracker
- **WHEN** se hace GET /api/trackers/{id}/history
- **THEN** el sistema devuelve todos los registros del tracker
