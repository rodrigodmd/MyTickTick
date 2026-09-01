# Proposal: MyTickTick App

## Why

TickTick no satisface completamente mis necesidades de productividad. Necesito una app más simple pero con mejor historial y capacidades de tracking específicas que TickTick no muestra bien.

## What Changes

- **Nueva aplicación de gestión de tareas** desde cero
- **Tres tipos de tareas**:
  - Tareas mensuales recurrentes con historial de cumplimiento
  - Tareas inmediatas para próximos días
  - Trackers diarios con límite unilateral (mínimo o máximo)
- **Dashboard de gráficos** para visualizar cumplimiento histórico
- **Login simple** (username + password) con JWT y "recordarme" (sesión persistente), rutas protegidas
- **Listas interactivas**: completar/eliminar con toggle o swipe, sin confirmaciones, con "undo"
- **Backend en Go** con GORM y PostgreSQL en Docker
- **Web app responsive** para usar desde móvil y desktop
- **Tema visual dark + moderno** (Vibrant Gradient) con acentos de gradiente

## Capabilities

### New Capabilities

- `monthly-tasks`: Gestión de tareas mensuales recurrentes con historial
- `immediate-tasks`: Gestión de tareas para próximos días
- `daily-trackers`: Trackers diarios con límite unilateral (mín o máx) y registro de desviación
- `compliance-history`: Historial y métricas de cumplimiento
- `dashboard`: Visualización gráfica de métricas (desktop)
- `api-rest`: API REST para backend (incluye autenticación y protección de rutas)
- `data-model`: Modelos de datos para usuario, tareas, historial y trackers

### Modified Capabilities

- (none)

## Impact

- **Nuevo backend**: Go + GORM + PostgreSQL
- **Autenticación**: username + password (hash), JWT en cookie httpOnly, "recordarme" para sesión persistente; todas las rutas protegidas
- **Nuevo frontend**: Web app responsive
- **Interacción de listados**: toggle y swipe sin confirmaciones, con "undo" (reemplaza los botones + `confirm()`)
- **Infraestructura**: Docker para levantar PostgreSQL
- **Persistencia**: Base de datos relacional con historial completo
