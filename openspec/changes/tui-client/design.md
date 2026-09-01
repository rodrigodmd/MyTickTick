# Design: TUI Client

## Context

MyTickTick es una app single-user (user_id=0) con API REST en Go (net/http + GORM + Postgres) y frontend web. El TUI es un **consumidor nuevo** de esa API; no se modifica el server. La API expone auth por cookie (`mtt_session`, JWT httpOnly) y rutas CRUD para `immediate-tasks`, `monthly-tasks`, `trackers` (con `records` upsert) y `metrics`.

Restricciones que guían el diseño:
- El repo ya es Go 1.26.5; el TUI va en el mismo módulo (`cmd/tui`).
- El server puede no estar arriba; el TUI debe degradar con un error claro.
- `internal/core` (domain, repos, services) vive del lado server; el TUI no lo importa.

## Goals / Non-Goals

**Goals:**
- Cliente TUI completo de la API en Go (bubbletea) dentro del mismo repo.
- Sesión persistente local (no relóguear cada arranque).
- Cubrir: tareas inmediatas, mensuales, trackers, métricas. (El dashboard queda fuera de esta iteración.)
- Manejo legible de errores de red y de servidor.

**Non-Goals:**
- No modificar la API ni el server (cero cambios en endpoints).
- No acceder a Postgres desde el TUI.
- No reemplazar al frontend web (coexisten).
- No soportar múltiples usuarios ni perfiles (single-user, user_id=0).
- No emular la UI de gráficos del web dashboard (las métricas se ven como tablas/series de texto).

## Decisions

### D1. El TUI habla HTTP con el server, no toca la DB
- **Decisión:** cliente HTTP puro contra la API REST existente.
- **Razón:** la lógica de negocio (activación idempotente de mensuales, cálculo `isMet`/`deviation`, métricas de compliance) ya vive en el server. Reimplementarla en el TUI crearía divergencia. HTTP local es gratis y la auth ya está resuelta.
- **Alternativa descartada:** binario TUI que abre Postgres directamente y linka `internal/core/services`. Rápido pero duplica caminos de negocio y rompe la fuente de verdad única.

### D2. UI en bubbletea + lipgloss (charm)
- **Decisión:** `github.com/charmbracelet/bubbletea` (model-update-view), `lipgloss` (estilos), `bubbles` (textinput, list, viewport).
- **Razón:** mismo lenguaje del repo, un solo binario, ecosistema maduro para TUIs.
- **Alternativa descartada:** Textual (Python) — prototipado rápido pero agrega un stack fuera del repo; ratatui (Rust) — overkill y lenguaje nuevo en el stack.

### D3. Sesión persistida en `~/.myticktick/session`
- **Decisión:** al login exitoso, el TUI guarda el valor de la cookie `mtt_session` en `~/.myticktick/session` (permisos 0600). En arranque, si existe, la reutiliza; si el server responde 401, la descarta y pide login.
- **Razón:** evita reloguear cada vez; respeta el modelo de sesión existente (cookie JWT stateless).
- **Alternativa descartada:** guardar usuario/contraseña en disco (más inseguro y el server no lo pide).

### D4. Base URL por env `MTT_URL` (default `http://localhost:8080`)
- **Decisión:** leer `MTT_URL`; si no está, usar el default. No se inventa config.toml para v1.
- **Razón:** mínimo y consistente con cómo el server lee `PORT`.
- **Alternativa descartada:** config.toml — más fricción que la que resuelve en v1.

### D5. Estructura de código (módulo `cmd/tui`)
```
cmd/tui/
  main.go          # flag/env, arranque tea.Program
  client.go        # cliente HTTP: base URL, cookie de sesión, peticiones por dominio
  session.go       # load/save/discard de ~/.myticktick/session
  app.go           # estado global de la app (pantalla activa, data cacheada)
  tasks.go         # tareas inmediatas (list, create, toggle, edit, delete)
  monthly.go       # mensuales (list, toggle, history, activate)
  trackers.go      # trackers (list, upsert del día, history)
  metrics.go       # métricas (mes/año, serie, trackers en rango)
  login.go         # pantalla de login
```
- **Razón:** una pantalla por archivo, cliente HTTP aislado, fácil de extender.

### D6. Modelo de datos del TUI
- **Decisión:** el TUI define sus propios structs Go con tags JSON que reflejan las respuestas de la API (`immediate-tasks`, `monthly-tasks`, `trackers`, `metrics`). No importa `internal/handlers`.
- **Razón:** mantiene la separación server/cliente; los shapes son públicos (JSON) y estables.

## Risks / Trade-offs

- [El server no está arriba] → el TUI muestra "no hay server en `MTT_URL`" y permite reintentar; no entra en spin infinito.
- [Sesión expirada a mitad de uso] → ante 401, el TUI descarta la sesión guardada y vuelve a la pantalla de login sin perder el contexto.
- [Divergencia de shapes JSON entre server y TUI] → el TUI no importa los types del server; si el server cambia una respuesta, el TUI puede romper. Mitigación: los shapes son estables y el TUI se prueba contra el server real; se puede agregar un test de contrato en el futuro.
- [bubbletea añade dependencias al go.mod] → solo las usa `cmd/tui`; el binario del server no las enlaza (separación por package).

## Migration Plan

No aplica al server (cero cambios). Para el usuario:
1. `make up` (server arriba, ya documentado).
2. `make tui` → arranca el cliente; primer uso pide login.
3. Rollback: no hay migración de datos; si el TUI falla, se sigue usando la web. El TUI es aditivo.

## Open Questions

- ¿`make tui` debe `go run` o `go build` + ejecutar el binario? (Detail de Makefile; se decide en tasks.)
- ¿El dashboard TUI se agrega en una iteración futura? (Fuera de scope de este change; se captura como nuevo change cuando se quiera.)
