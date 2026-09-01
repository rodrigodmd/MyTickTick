# Design: MyTickTick App

## Context

Proyecto nuevo sin infraestructura previa. El usuario quiere una app de productividad simple pero con buen historial.

## Goals / Non-Goals

**Goals:**
- Backend en Go con GORM para ORM
- PostgreSQL como base de datos
- Docker para levantar la infraestructura
- API REST para comunicación con frontend
- Web app responsive para móvil y desktop
- Dashboard con gráficos para visualización en desktop

**Non-Goals:**
- Multi-usuario (por ahora solo para un usuario)
- Sincronización en tiempo real
- Notificaciones push
- Google OAuth (pendiente para futuro)

## Decisions

### Backend: Go + GORM

**Decisión:** Usar Go como lenguaje de backend con GORM como ORM.

**Rationale:**
- Go es simple y rápido para APIs REST
- GORM ya fue usado y gustó al usuario
- Fácil de mantener y extender

**Alternativas consideradas:**
- Sin ORM (más control pero más complejo)
- sqlc (type-safe pero menos conocido)

### Base de Datos: PostgreSQL en Docker

**Decisión:** PostgreSQL como base de datos, levantado con Docker.

**Rationale:**
- Robusto y confiable
- Buen soporte en Go
- Fácil de levantar con Docker
- Permite migrar fácilmente si se necesita más potencia

**Alternativas consideradas:**
- SQLite (más simple pero menos robusto)
- MySQL (similar a PostgreSQL pero menos features)

### Arquitectura: API REST

**Decisión:** API REST para comunicación entre frontend y backend.

**Rationale:**
- Simple y estándar
- Fácil de documentar y consumir
- Bueno para apps CRUD como esta

**Alternativas consideradas:**
- GraphQL (más flexible pero más complejo)
- gRPC (más performante pero menos accesible)

### Frontend: Web App Responsive

**Decisión:** Web app con diseño responsive para móvil y desktop.

**Rationale:**
- Una sola base de código
- Se adapta a ambos dispositivos
- No requiere app nativa

**Alternativas consideradas:**
- App nativa móvil (más compleja de mantener)
- PWA ( Progressive Web App)

### Autenticación: username + password, JWT en cookie, "recordarme"

**Decisión:** Login simple con **username** (sin email) y password. El backend almacena el password con hash (bcrypt). En el login se emite un **JWT** que se entrega como **cookie httpOnly**:
- Con **"Recordarme"**: expiración muy larga (p. ej. 10 años) → sesión persistente.
- Sin "Recordarme": expiración corta (p. ej. 7 días).

Un **middleware de autenticación** protege todas las rutas `/api/*` de tareas, trackers y métricas; quedan públicas solo `/api/health`, `/api/register`, `/api/login` y `/api/logout`.

**Rationale:**
- El usuario pidió login simple con recordatorio persistente
- JWT en cookie httpOnly: sin exponer el token al JS, el navegador lo reenvía solo, y el "recordarme" es trivial (solo cambia la expiración)
- Hash bcrypt en el servidor (hoy el password se guarda en texto plano, lo que el plan corrige)

**Alternativas consideradas:**
- Sesión server-side (tabla de sesiones): más estado en el servidor
- Token expuesto al frontend (localStorage): menos seguro (XSS)
- Google OAuth: fuera de alcance por ahora

### Interacción de listados: toggle + swipe, sin confirmaciones, con "undo"

**Decisión:** Completar y eliminar tareas se hace por **toggle** (clic en el check) o **swipe** (deslizar la fila), **sin diálogos de confirmación**. Tras completar o eliminar aparece un **toast con "Deshacer"** durante unos segundos para revertir la acción.

**Rationale:**
- El usuario no quiere botones + `confirm()` molesto ("que te pregunte al pedo")
- Toggle y swipe son el estándar en apps de tareas y funcionan bien en móvil
- El "undo" (toast) protege contra errores sin interrumpir el flujo

**Alternativas consideradas:**
- `confirm()` nativo: interrumpe el flujo (descartado por el usuario)
- Modales de confirmación: igual de intrusivos
- Sin revertir (peligroso para borrar historial)

### Tema visual: Vibrant Gradient (dark + moderno)

**Decisión:** El frontend servido usa el tema **Vibrant Gradient**: fondo oscuro, acentos con gradientes y estética moderna.

**Paleta de referencia** (`frontend/css/variables.css`):
- **Fondo**: `#0a0a0a` (primario), `#111111` (tarjetas), `#1a1a1a` (hover)
- **Gradientes**: primary `#667eea → #764ba2`, secondary `#f093fb → #f5576c`
- **Acentos**: púrpura `#8b5cf6`, rosa `#ec4899`, azul `#3b82f6`, verde `#10b981`
- **Texto**: `#ffffff` (primario), `#a0a0a0` (secundario)
- **Bordes**: `rgba(139, 92, 246, 0.2)` con glow `0 0 20px rgba(139, 92, 246, 0.3)`

**Rationale:**
- Elegido explícitamente por el usuario durante el análisis (opción "Vibrant Gradient / Tema 3")
- Estética moderna y consistente entre vistas de auth y dashboard

**Alternativas consideradas:**
- Tema claro (mockup original: fondo `#f5f5f5`, tarjetas `#ffffff`) — descartado

**Nota de implementación:** la web servida (`web/`) debe usar esta paleta. El mockup original con tema claro queda obsoleto.

## Risks / Trade-offs

**Riesgo:** PostgreSQL en Docker consume más recursos que SQLite
→ **Mitigación:** Monitorizar uso de recursos y ajustar si es necesario

**Riesgo:** Web app puede ser menos rápida que app nativa en móvil
→ **Mitigación:** Optimizar carga y usar caching

**Riesgo:** GORM puede ser menos performante que queries raw
→ **Mitigación:** Usar GORM para desarrollo, evaluar performance después

## Migration Plan

1. **Setup inicial:**
   - Crear Dockerfile para backend
   - Configurar docker-compose para PostgreSQL
   - Implementar modelos de datos con GORM (User, ImmediateTask, DailyTracker, TrackerEntry, MonthlyTask, MonthlyTaskHistory)

2. **Desarrollo de API:**
   - Implementar endpoints de tareas mensuales (CRUD + historial)
   - Implementar endpoints de tareas inmediatas (CRUD)
   - Implementar endpoints de trackers (CRUD + records)
   - Implementar endpoints de historial y métricas

3. **Servicios de negocio:**
   - Implementar MonthlyTaskService con activación automática el día 1
   - Implementar TrackerService con cálculo de límite unilateral (mín/máx) y desviación
   - Implementar TaskService para tareas inmediatas
   - Implementar AuthService (username + password hash, JWT, middleware de protección)

4. **Desarrollo de frontend:**
   - Crear estructura de web app
   - Implementar vistas para cada tipo de tarea
   - Implementar dashboard con gráficos
   - Implementar toggle + swipe + toast "Deshacer" en las listas (reemplaza botones + `confirm()`)
   - Implementar vistas de login/registro con "Recordarme"

5. **Testing:**
   - Testear endpoints de API
   - Testear frontend en móvil y desktop
   - Testear dashboard de gráficos
   - Testear lógica de negocio (límite unilateral mín/máx, desviación, activación mensual)
   - Testear autenticación (registro, login, recordarme, protección de rutas)

## Open Questions

- (Resuelta) Librería de gráficos: **Chart.js**, ya integrada en el dashboard.
- (Resuelta) Autenticación: login simple username + password con JWT en cookie y "recordarme"; Google OAuth queda pendiente para el futuro.
