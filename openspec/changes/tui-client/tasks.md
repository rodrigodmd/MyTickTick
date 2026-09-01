# Tasks: TUI Client

## 1. Setup

- [ ] 1.1 Agregar dependencias charm (bubbletea, lipgloss, bubbles) al `go.mod` del repo
- [ ] 1.2 Crear estructura `cmd/tui/` con `main.go` mínimo que arranca un `tea.Program` de prueba
- [ ] 1.3 Agregar target `make tui` al Makefile (compilar y correr el cliente)

## 2. Cliente HTTP y sesión

- [ ] 2.1 Implementar `session.go`: load/save/discard de la cookie en `~/.myticktick/session` (permisos 0600)
- [ ] 2.2 Implementar `client.go`: base URL desde `MTT_URL` (default `http://localhost:8080`), cookie de sesión, timeout, detección de 401
- [ ] 2.3 Implementar pantalla de `login.go` (usuario + contraseña), con manejo de credenciales inválidas y reintentos
- [ ] 2.4 Wire: al login exitoso guardar sesión; en 401 descartar sesión y volver a login

## 3. Estructura de app

- [ ] 3.1 Implementar `app.go`: estado global (pantalla activa, data cacheada), navegación entre pantallas por teclas
- [ ] 3.2 Implementar pantalla de error de conexión (server no disponible) con opción de reintentar

## 4. Pantallas

- [ ] 4.1 Implementar `tasks.go`: listar, crear, editar, toggle completada, eliminar tareas inmediatas
- [ ] 4.2 Implementar `monthly.go`: listar, toggle cumplimiento del mes, historial, activar
- [ ] 4.3 Implementar `trackers.go`: listar, cargar/corregir valor del día (upsert), historial con desviación
- [ ] 4.4 Implementar `metrics.go`: métricas por mes/año, serie mensual, trackers en rango

## 5. Verificación

- [ ] 5.1 Correr el TUI contra el server real (`make up` + `make tui`) y verificar login y las 4 pantallas
- [ ] 5.2 Verificar sesión persistente: cerrar y reabrir el TUI sin relóguear
- [ ] 5.3 Verificar server caído: error claro de conexión en vez de spin infinito
- [ ] 5.4 Verificar 401: sesión inválida dispara re-login
- [ ] 5.5 `make build` sigue compilando el server sin las dependencias del TUI
