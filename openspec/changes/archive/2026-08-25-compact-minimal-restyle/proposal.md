# Proposal: Compact Minimal Restyle

## Why

La app vista desde el móvil se percibe "inflada": tipografía grande, paddings y gaps holgados, y un tema "Vibrant Gradient" (gradientes + glows) que resta densidad. El objetivo es que la app se sienta compacta, con letras más chicas y un look minimal, al estilo de las apps de referencia (flat, un solo acento, lista densa y legible). El foco es legibilidad y densidad en pantallas chicas, sin tocar la funcionalidad.

## What Changes

- **Densidad**: reducir el tamaño base de tipografía y el espaciado vertical (gaps, paddings) para que entren más items por pantalla en mobile.
- **Tipografía más chica**: bajar escalas (títulos, etiquetas, meta-texto) y apretar `line-height` sin perder legibilidad.
- **Tema flat/minimal**: reemplazar los gradientes y glows por superficies planas (elevation por tono, no por sombra) y radios/paddings más chicos.
- **Un solo acento con propósito**: concentrar el color en un único acento (naranja) para estados activos y prioridad, en lugar de morado + rosa + glows en varias superficies.
- **Controles compactos**: cards, botones, badges/pills, checkboxes y formularios más chicos y apretados; primary actions con mejor jerarquía en móvil.
- **Sin cambios de comportamiento**: el CRUD de tareas, trackers, dashboard, historial y autenticación quedan exactamente igual; solo cambia la presentación.

## Capabilities

### New Capabilities

- (none)

### Modified Capabilities

- (none)

> **Nota**: este cambio es **solo de presentación (CSS/estilo)**. No modifica ningún requirement de comportamiento, por lo que no requiere delta de specs. El change declara `skip_specs: true` en `.openspec.yaml` (cambio de UI, no de comportamiento), conforme a la guía de OpenSpec para evitar inventar requirements que no existen.

## Impact

- **Frontend (lo único afectado)**: `web/css/style.css` (tema, variables, densidad), y de forma puntual `web/html/*.html` / `web/js/*` **solo si** se ajustan clases o la estructura mínima de un control (p. ej. orden de botones en un formulario). Se conserva el set de clases y el DOM funcional.
- **Backend / API**: sin cambios. Mismas rutas, mismos flujos, mismos datos.
- **Comportamiento observable**: la app sigue haciendo lo mismo; cambia el aspecto (tamaños, colores, espaciado) en mobile y desktop.
