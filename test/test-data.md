# Test Data for MyTickTick (esquema actual)

Usuario de test: `testapi_<timestamp>` (lo crea el script; contraseña `Test123!`)

## Monthly Tasks
- Name: "Pagar cuota"
  - Description: "Cuota mensual de gimnasio"
  - isActive: true
- Name: "Lectura mensual"
  - Description: "Leer un libro este mes"
  - isActive: true

## Immediate Tasks
- Name: "Pagar facturas"
  - Description: "Pagar facturas de servicios"
  - dueDate: +3 días (RFC3339)
  - priority: high
  - isCompleted: false
- Name: "Comprar comida"
  - Description: "Comprar comida para la semana"
  - dueDate: +1 día
  - priority: medium
  - isCompleted: true (se marca con PUT isCompleted)

## Trackers
- Name: "Peso" (límite máximo)
  - limitValue: 85
  - limitType: "max"
  - unit: "kg"
- Name: "Sueño" (límite mínimo)
  - limitValue: 6
  - limitType: "min"
  - unit: "h"

## Tracker Entries (y desviación esperada)
| Tracker | value | fecha  | esperado isMet | esperado deviation |
|---------|-------|--------|----------------|--------------------|
| Peso    | 84    | hoy-1  | true           | 0                  |
| Peso    | 87    | hoy    | false          | 2                  |
| Sueño   | 7     | hoy-1  | true           | 0                  |
| Sueño   | 5     | hoy    | false          | 1                  |

## Monthly completions / history
- "Pagar cuota": completada este mes (PUT /completion) → aparece en GET /history con completed=true
- POST /api/monthly-tasks/activate → idempotente (2da llamada crea 0 registros)

## Auth
- POST /api/register → 201 + cookie mtt_session httpOnly
- POST /api/register (duplicado) → 409
- POST /api/login credenciales malas → 401
- POST /api/login recordarme=false → Max-Age 604800 (7 días)
- POST /api/login recordarme=true  → Max-Age 315360000 (10 años)
- Rutas protegidas sin cookie → 401
- POST /api/logout → cookie borrada (Max-Age 0)
