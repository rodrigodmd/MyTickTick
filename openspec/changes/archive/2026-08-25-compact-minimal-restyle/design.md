# Design: Compact Minimal Restyle

## Context

El frontend que el backend sirve en producción es la carpeta `web/` (`internal/handlers/static/handler.go` → `NewStaticHandler("web")` en `cmd/api/main.go`). El CSS relevante es `web/css/style.css` (1166 líneas, **mobile-first**: base para pantallas chicas + media queries hacia arriba). Usa el tema "Vibrant Gradient": variables en `:root` (colores, `--radius: 16px`, `--gap: 16px`, gradientes `--gradient-primary/secondary`, glows `--shadow-hover`) aplicados a header, cards, botones, formularios, listas, modales, dashboard y auth.

Las apps de referencia (que el usuario quiere imitar) comparten un estilo **flat/minimal**: fondo casi negro, una sola superficie "un tono arriba", **un único acento naranja** reservado para lo activo/importante, tipografía chica y consistente, spacing apretado pero rítmico, radios chicos, y sin gradientes ni sombras de "brillo".

Esto es un cambio **solo de presentación**: el DOM, el JS y la API no cambian (ver proposal.md).

## Goals / Non-Goals

**Goals:**
- Mayor **densidad en mobile**: más items visibles por pantalla, sin scroll excesivo.
- **Tipografía más chica** y jerarquía clara (título > etiqueta > meta).
- **Look flat/minimal** con **un solo acento (naranja)** usado con propósito.
- Cambios **centralizados** (tokens en `:root`) para que sean fáciles de afinar y revertir.
- Mantener **mobile-first** y no romper layouts en tablet/desktop.

**Non-Goals:**
- No cambiar comportamiento: CRUD, trackers, dashboard, historial y auth quedan idénticos.
- No reescribir HTML/JS: el DOM y la lógica se conservan; solo se tocan clases/estructura mínima si un control lo exige (p. ej. jerarquía de botones en un formulario).
- No introducir un framework de UI ni un build step (se sigue con CSS plano servido estático).
- No tocar el backend, la API ni la base de datos.

## Decisions

**D1 — Frontend objetivo: `web/` (no `frontend/`).**
El backend sirve `web/`. `frontend/` es una copia anterior/legada que no se sirve en producción. El restyle se aplica sobre `web/`. *Alternativa descartada*: restylear `frontend/` (no tendría efecto en la app que el usuario usa).

**D2 — Densidad vía tokens centralizados, no por regla suelta.**
Se redefine la base tipográfica y el espaciado en `:root` y en las reglas base (`body`, `.page`, `.card`, `.form-*`, listas, botones), de modo que la mayor parte del compactado fluya por variables. *Razón*: el CSS actual repite paddings/gaps en muchos selectores; centralizar evita cambios fragmentados y facilita afinar y revertir. *Alternativa descartada*: editar cada selector individual (frágil y difícil de revertir).

**D3 — Escala de tipografía compacta (propuesta base, a validar visualmente).**
Bajar el tamaño base de ~16px a ~14px (móvil) y re-escalar títulos/etiquetas/meta hacia abajo, con `line-height` algo más apretado (p. ej. 1.4–1.5). Propuesta indicativa:
- Título de página: ~1.5rem → ~1.125rem
- Etiquetas/labels: ~0.9–0.95rem → ~0.8rem
- Meta-texto (fechas, badges): ~0.8rem (se mantiene pequeño)
- Texto de inputs/selects: ~1rem → ~0.9rem
*Alternativa descartada*: solo reducir `line-height` (mejora poco la densidad real).

**D4 — Spacing y radios compactos (propuesta base).**
Reducir `--gap` (~16px → ~10px), paddings de cards/inputs/botones y radios (`--radius: 16px` → ~10–12px; botones/inputs a ~8–10px). Menos espacio entre filas de una card y entre cards. *Alternativa descartada*: mantener radios grandes (contradice el look minimal de referencia).

**D5 — Flat: quitar gradientes y glows; elevation por tono.**
Reemplazar `--gradient-primary/secondary` y los "glows" (`--shadow-hover`, `box-shadow` de borde) por **color plano** y elevación por un tono de superficie ligeramente más claro + borde sutil. *Razón*: los gradientes/glows son la principal fuente de "ruido" vs. las referencias flat. *Alternativa descartada*: mantener un glow sutil (sigue alejándose del look minimal pedido).

**D6 — Un solo acento (naranja) con propósito.**
Concentrar el color en un único acento naranja (el proyecto ya define `--accent-orange: #f59e0b` en `variables.css` como referencia) para: nav activa, primary action, prioridad alta, foco de inputs y estados "activos". Los demás colores de prioridad (media/baja) usan tonos neutros/atenidos. *Razón*: las referencias usan un solo acento para "lo importante" y el resto es gris. *Alternativa considerada*: mantener acento morado plano (el usuario no respondió el clarify; se asume naranja por alineación con las referencias pasadas — queda registrado como supuesto y es fácil de revertir por ser un token).

**D7 — Jerarquía de controles en móvil.**
En formularios, hacer la primary action más prominente (p. ej. ancho completo en móvil) y separar la acción destructiva (Eliminar) para que no quede "huérfana" apretada bajo las otras. *Alternativa descartada*: dejar todo inline (produce el wrap/orfandad visible en las capturas).

## Risks / Trade-offs

- **[Algo se ve demasiado chico en mobile]** → La escala (D3) es una base *indicativa*; se valida contra capturas reales en el paso de verificación y se ajusta por token. Se conserva `min-height` de tap targets razonable (~40px) para no sacrificar usabilidad por densidad.
- **[Restylear solo `web/` deja `frontend/` desactualizado]** → Se documenta que `frontend/` no se sirve; se deja intacto (o se alinea si el usuario lo pide). Riesgo bajo: no afecta la app en producción.
- **[Quitar gradientes/glows cambia el "carácter" del tema]** → Es el objetivo explícito (look minimal). Se mantiene coherencia usando un mismo acento y superficies planas en todas las secciones.
- **[El acento naranja es una suposición (clarify sin respuesta)]** → Queda registrado como supuesto en proposal.md y es trivial de revertir (es un token). Si el usuario prefiere morado plano, solo cambia el valor del acento.
- **[Algunas clases están atadas a gradientes en varios selectores]** → Se hace el cambio de forma sistemática por token y se revisa que no queden `background-image: var(--gradient-*)` ni `box-shadow` de glow rezando.

## Open Questions

- ¿Acento naranja (asumido) o morado plano? — **Asumido naranja** por alineación con las referencias; revertible por token. No bloquea el plan ni las tasks.
- ¿También alinear el frontend legado `frontend/`? — **No** (no se sirve); se decide solo si el usuario lo pide.
