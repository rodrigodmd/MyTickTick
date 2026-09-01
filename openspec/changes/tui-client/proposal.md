# Proposal: TUI Client

## Why

MyTickTick hoy solo se usa desde el navegador web, que exige el server arriba y un browser abierto. Un cliente TUI local (Go, bubbletea) permite usar toda la app desde la terminal: ver el estado de hoy, cargar valores de trackers, togglear tareas, revisar métricas — sin salir del flujo de trabajo en terminal y sin abrir browser.

## What Changes

- **Nuevo cliente TUI** en Go (`cmd/tui`) dentro del mismo repo, que consume la API REST existente sobre HTTP (no toca la base de datos ni `internal/core`).
- **Pantallas v1**:
  - Tareas inmediatas: listar, toggle completada, crear, editar, borrar
  - Tareas mensuales: listar, toggle de cumplimiento del mes, historial, activar
  - Trackers: listar, cargar/corregir valor del día, historial
  - Métricas: cumplimiento por mes/año, serie mensual, trackers en rango
- **Autenticación en el TUI**: login con usuario/contraseña (POST /api/login), cookie de sesión persistida localmente para no reloguear; re-login cuando la sesión expira.
- **Config mínima**: base URL por env `MTT_URL` (default `http://localhost:8080`); sesión y estado en `~/.myticktick/`.
- **Makefile**: target `make tui` para correr el cliente.
- **Nuevas dependencias Go**: bubbletea, lipgloss, bubbles (solo en el binario TUI; el server no las usa).
- **Sin cambios en la API existente**: el TUI es un consumidor de las rutas actuales; no se agregan, modifican ni quitan endpoints.

## Capabilities

### New Capabilities

- `tui-client`: Cliente terminal interactivo para operar MyTickTick (tareas, mensuales, trackers, métricas) contra la API REST, con sesión persistente local.

### Modified Capabilities

- (none — la TUI consume la API tal cual; `api-rest` y demás specs no cambian de comportamiento)

## Impact

- **Código nuevo**: `cmd/tui/` (main, app bubbletea, cliente HTTP con cookie, pantallas por dominio).
- **Dependencias**: `github.com/charmbracelet/bubbletea`, `lipgloss`, `bubbles` agregadas a `go.mod` (solo usadas por `cmd/tui`).
- **Build**: `make build` sigue compilando el server; `make tui` compila/corre el cliente.
- **Runtime**: requiere el server API arriba (`make up`); sin server, el TUI muestra error claro de conexión.
- **Sin impacto en datos**: el TUI no abre conexiones a Postgres ni importa `internal/core`; toda la lógica de negocio sigue viviendo en el server.
