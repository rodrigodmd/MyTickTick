# Tasks: Compact Minimal Restyle

## 1. Baseline y tokens

- [x] 1.1 Tomar capturas "antes" (mobile ~390px) de home, immediate-tasks, trackers, dashboard y auth para comparar después.
- [x] 1.2 En `web/css/style.css`, definir el nuevo acento naranja como token `--accent` (reutilizando el valor `#f97316` de referencia) y los tokens de densidad (`--gap`, `--radius`) base en `:root`.
- [x] 1.3 Reducir el tamaño base de tipografía (`body`) a ~14px y bajar `line-height` a ~1.45, manteniendo legibilidad.

## 2. Tipografía y espaciado compactos

- [x] 2.1 Re-escalar hacia abajo: `.page-title` (~1.2rem), etiquetas/`.form-label` (~0.82rem), inputs/selects (~0.88rem), meta-texto (~0.78rem).
- [x] 2.2 Reducir paddings de `.page`, `.card`/`.home-card`, `.form-input`, botones `.btn`, y el gap entre filas/items de lista.
- [x] 2.3 Bajar radios (`--radius` a 10px; botones/inputs a 6px) para un look más tight.
- [x] 2.4 Ajustar los paddings de header/nav (`site-header`, `site-nav`, `nav-toggle`) y del footer para que la app ocupe menos altura en mobile.

## 3. Flat: quitar gradientes y glows

- [x] 3.1 Reemplazar usos de `--gradient-primary`/`--gradient-secondary` (brand, títulos, botones primarios, nav activa) por **color plano** del acento o de superficie.
- [x] 3.2 Quitar "glows": `--shadow-hover` (borde luminoso) y `box-shadow` de halo en hover/focus; cambiar el focus a un anillo plano y discreto.
- [x] 3.3 Cambiar la elevación de cards a "un tono arriba" + borde sutil (sin glow), y el fondo body a color plano (quitar los `radial-gradient` de fondo).
- [x] 3.4 Buscar y eliminar referencias rezadas: `background-image: var(--gradient-*)`, `background-clip: text`, `box-shadow` de glow, en todo `web/css/style.css`.

## 4. Acento con propósito y controles

- [x] 4.1 Usar el acento naranja solo para: nav activa, primary action, prioridad alta, foco de inputs y estados activos.
- [x] 4.2 Atenuar las demás prioridades (media/baja) a tonos neutros/grises; ajustar badges/pills para que sean chicos (padding chico, texto ~0.72rem).
- [x] 4.3 Compactar checkboxes/toggles, chips de icono y filas de lista para mayor densidad (mantener tap target ≥ ~36px).
- [x] 4.4 En formularios (immediate-tasks, trackers, auth): hacer la primary action prominente (ancho completo en auth) y separar la acción destructiva (Eliminar) para evitar el wrap/orfandad.

## 5. Verificación visual y QA

- [x] 5.1 Recorrer las 5 vistas (home, immediate-tasks, monthly-tasks, trackers, dashboard) + auth en mobile ~390px y capturar "después".
- [x] 5.2 Comparar "antes/después" contra las apps de referencia: densidad, tipografía, acento único, flat. Ajustar tokens si algo queda demasiado chico o desalineado.
- [x] 5.3 Verificar tablet/desktop (media queries) para que el compactado no rompa layouts anchos.
- [x] 5.4 Smoke test de comportamiento (no debe cambiar): crear/completar/eliminar tarea, registrar tracker, ver dashboard y login siguen funcionando. (Solo CSS — no se modificó HTML/JS de comportamiento.)
