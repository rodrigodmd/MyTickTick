#!/usr/bin/env bash
# 12.1 - Test de endpoints de API con datos de prueba
# Uso: BASE=http://localhost:8080 bash test/test-api.sh
# Requiere: curl, jq. El backend debe estar corriendo.
set -u

BASE="${BASE:-http://localhost:8080}"
JAR="$(mktemp)"
PASS=0; FAIL=0
USERNAME="testapi_$(date +%s)"
PASSWORD='Test123!'

check() { # <nombre> <real> <esperado>
  local name="$1" actual="$2" expected="$3"
  if [ "$actual" = "$expected" ]; then
    PASS=$((PASS+1)); echo "  PASS: $name"
  else
    FAIL=$((FAIL+1)); echo "  FAIL: $name (esperado=$expected, real=$actual)"
  fi
}

json() { # <expresión jq...> (admite --arg etc. antes del filtro)
  local out
  if out="$(jq -c "$@" 2>/dev/null)"; then echo "$out"; else echo "NO-JSON"; fi
}
status_of() { # <archivo de respuesta con headers> -> status code
  head -1 "$1" | awk '{print $2}'
}
setcookie_maxage() { # <archivo headers> -> Max-Age de mtt_session
  grep -i '^set-cookie:.*mtt_session' "$1" | grep -o 'Max-Age=[0-9]*' | head -1 | cut -d= -f2
}

req() { # <method> <path> [json body] -> escribe headers+body en $OUT; $RC=status
  local method="$1" path="$2" body="${3:-}"
  OUT="$(mktemp)"
  if [ -n "$body" ]; then
    curl -s -X "$method" "$BASE$path" -H 'Content-Type: application/json' \
      -H 'Cookie: mtt_session='"$TOKEN" \
      -d "$body" -D - -o "${OUT}.body" > "$OUT" 2>/dev/null
  else
    curl -s -X "$method" "$BASE$path" \
      -H 'Cookie: mtt_session='"$TOKEN" \
      -D - -o "${OUT}.body" > "$OUT" 2>/dev/null
  fi
  RC="$(status_of "$OUT")"
}

echo "=== MyTickTick API test (base=$BASE) ==="

# ---------- 1. Health (público) ----------
echo "[1] GET /api/health"
OUT="$(mktemp)"; curl -s "$BASE/api/health" -D - -o "${OUT}.body" > "$OUT"
check "status 200" "$(status_of "$OUT")" "200"
check "body ok" "$(json '.status' < "${OUT}.body")" '"ok"'

# ---------- 2. Auth ----------
echo "[2] POST /api/register"
OUT="$(mktemp)"
curl -s -X POST "$BASE/api/register" -H 'Content-Type: application/json' \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" -D - -o "${OUT}.body" > "$OUT"
check "status 201" "$(status_of "$OUT")" "201"
check "emite cookie mtt_session" "$(grep -ci 'set-cookie:.*mtt_session' "$OUT")" "1"
check "cookie httpOnly" "$(grep -io 'httponly' "$OUT" | head -1 | tr 'A-Z' 'a-z')" "httponly"

echo "[3] POST /api/register (duplicado)"
OUT="$(mktemp)"
curl -s -X POST "$BASE/api/register" -H 'Content-Type: application/json' \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" -D - -o "${OUT}.body" > "$OUT"
check "status 409" "$(status_of "$OUT")" "409"

echo "[4] POST /api/login (credenciales malas)"
OUT="$(mktemp)"
curl -s -X POST "$BASE/api/login" -H 'Content-Type: application/json' \
  -d "{\"username\":\"$USERNAME\",\"password\":\"clave-mala\"}" -D - -o "${OUT}.body" > "$OUT"
check "status 401" "$(status_of "$OUT")" "401"

echo "[5] POST /api/login (sin recordarme)"
OUT="$(mktemp)"
curl -s -X POST "$BASE/api/login" -H 'Content-Type: application/json' \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\",\"remember\":false}" -D - -o "${OUT}.body" > "$OUT"
check "status 200" "$(status_of "$OUT")" "200"
check "Max-Age 7 dias (604800)" "$(setcookie_maxage "$OUT")" "604800"

echo "[6] POST /api/login (con recordarme)"
OUT="$(mktemp)"
curl -s -X POST "$BASE/api/login" -H 'Content-Type: application/json' \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\",\"remember\":true}" -D - -o "${OUT}.body" > "$OUT"
check "status 200" "$(status_of "$OUT")" "200"
check "Max-Age 10 años (315360000)" "$(setcookie_maxage "$OUT")" "315360000"
TOKEN="$(grep -io 'mtt_session=[^;]*' "$OUT" | head -1 | cut -d= -f2-)"
check "cookie tiene token" "$([ -n "$TOKEN" ] && echo yes)" "yes"

echo "[7] Rutas protegidas sin sesión"
for path in /api/monthly-tasks /api/immediate-tasks /api/trackers /api/metrics; do
  OUT="$(mktemp)"
  curl -s "$BASE$path" -D - -o "${OUT}.body" > "$OUT"
  check "$path sin cookie -> 401" "$(status_of "$OUT")" "401"
done

# ---------- 3. Monthly tasks ----------
echo "[8] Monthly tasks CRUD"
req POST /api/monthly-tasks '{"name":"Pagar cuota","description":"Cuota mensual de gimnasio"}'
check "POST create -> 201" "$RC" "201"
MT_ID="$(json '.id' < "$OUT.body")"
check "devuelve id" "$([ -n "$MT_ID" ] && echo yes)" "yes"

req GET /api/monthly-tasks
check "GET list -> 200" "$RC" "200"
check "lista contiene la tarea" "$(json --arg id "$MT_ID" '[.[] | select(.id == ($id|tonumber))] | length' < "$OUT.body")" "1"

req GET "/api/monthly-tasks/$MT_ID"
check "GET por id -> 200" "$RC" "200"
check "nombre correcto" "$(json '.name' < "$OUT.body")" '"Pagar cuota"'

req PUT "/api/monthly-tasks/$MT_ID" '{"name":"Pagar cuota 2","description":"actualizada"}'
check "PUT update -> 200" "$RC" "200"
check "nombre actualizado" "$(json '.name' < "$OUT.body")" '"Pagar cuota 2"'

req GET "/api/monthly-tasks/999999"
check "GET inexistente -> 404" "$RC" "404"

# ---------- 4. Monthly completion + history + activate ----------
echo "[9] Completion / historial / activación"
req PUT "/api/monthly-tasks/$MT_ID/completion" '{"month":8,"year":2026}'
check "PUT completion -> 200" "$RC" "200"
check "completed=true" "$(json '.completed' < "$OUT.body")" "true"

req GET "/api/monthly-tasks/$MT_ID/history"
check "GET history -> 200" "$RC" "200"
check "history contiene mes 8/2026 cumplido" \
  "$(json '[.[] | select(.month==8 and .year==2026 and .completed==true)] | length' < "$OUT.body")" "1"

req POST /api/monthly-tasks/activate
check "POST activate -> 200" "$RC" "200"
FIRST="$(json '.recordsCreated' < "$OUT.body")"
req POST /api/monthly-tasks/activate
check "activate idempotente (2da vez = 0)" "$(json '.recordsCreated' < "$OUT.body")" "0"

# ---------- 5. Immediate tasks ----------
echo "[10] Immediate tasks CRUD"
DUE3="$(date -u -d '+3 days' '+%Y-%m-%dT23:59:59Z')"
DUE1="$(date -u -d '+1 day' '+%Y-%m-%dT23:59:59Z')"
req POST /api/immediate-tasks "{\"name\":\"Pagar facturas\",\"description\":\"servicios\",\"dueDate\":\"$DUE3\",\"priority\":\"high\"}"
check "POST create -> 201" "$RC" "201"
IT_ID="$(json '.id' < "$OUT.body")"
check "priority high" "$(json '.priority' < "$OUT.body")" '"high"'

req POST /api/immediate-tasks "{\"name\":\"Comprar comida\",\"dueDate\":\"$DUE1\",\"priority\":\"medium\"}"
IT2_ID="$(json '.id' < "$OUT.body")"

req GET /api/immediate-tasks
check "GET list -> 200" "$RC" "200"
check "lista contiene ambas" "$(json --arg a "$IT_ID" --arg b "$IT2_ID" '[.[] | select(.id==($a|tonumber) or .id==($b|tonumber))] | length' < "$OUT.body")" "2"

req GET /api/immediate-tasks
check "orden por dueDate" "$(json 'map(.dueDate) | . == (sort)' < "$OUT.body")" "true"

req PUT "/api/immediate-tasks/$IT2_ID" '{"name":"Comprar comida","description":"","dueDate":"","isCompleted":true}'
check "PUT isCompleted -> 200" "$RC" "200"
check "isCompleted=true" "$(json '.isCompleted' < "$OUT.body")" "true"

req DELETE "/api/immediate-tasks/$IT2_ID"
check "DELETE -> 204" "$RC" "204"
req GET "/api/immediate-tasks/$IT2_ID"
check "GET tras DELETE -> 404" "$RC" "404"

# ---------- 6. Trackers + límites unilaterales + desviación ----------
echo "[11] Trackers (límite unilateral y desviación)"
TODAY="$(date -u +%Y-%m-%d)"
YESTERDAY="$(date -u -d 'yesterday' +%Y-%m-%d)"

req POST /api/trackers '{"name":"Peso","limitValue":85,"limitType":"max","unit":"kg"}'
check "POST tracker max -> 201" "$RC" "201"
PESO="$(json '.id' < "$OUT.body")"
check "limitType=max" "$(json '.limitType' < "$OUT.body")" '"max"'

req POST /api/trackers '{"name":"Sueño","limitValue":6,"limitType":"min","unit":"h"}'
check "POST tracker min -> 201" "$RC" "201"
SUENO="$(json '.id' < "$OUT.body")"
check "limitType=min" "$(json '.limitType' < "$OUT.body")" '"min"'

req PUT "/api/trackers/$PESO" '{"name":"Peso","limitValue":90,"limitType":"max","unit":"kg","isActive":true}'
check "PUT tracker -> 200" "$RC" "200"
check "limitValue actualizado" "$(json '.limitValue' < "$OUT.body")" "90"
req PUT "/api/trackers/$PESO" '{"name":"Peso","limitValue":85,"limitType":"max","unit":"kg","isActive":true}'

# max: <= límite cumple, > límite desvía
req POST "/api/trackers/$PESO/records" "{\"value\":84,\"date\":\"$YESTERDAY\"}"
check "max 84/85 -> isMet" "$(json '.isMet' < "$OUT.body")" "true"
check "max 84/85 -> desviación 0" "$(json '.deviation' < "$OUT.body")" "0"

req POST "/api/trackers/$PESO/records" "{\"value\":87,\"date\":\"$TODAY\"}"
check "max 87/85 -> no cumple" "$(json '.isMet' < "$OUT.body")" "false"
check "max 87/85 -> desviación 2" "$(json '.deviation' < "$OUT.body")" "2"

# min: >= límite cumple, < límite desvía
req POST "/api/trackers/$SUENO/records" "{\"value\":7,\"date\":\"$YESTERDAY\"}"
check "min 7/6 -> isMet" "$(json '.isMet' < "$OUT.body")" "true"
check "min 7/6 -> desviación 0" "$(json '.deviation' < "$OUT.body")" "0"

req POST "/api/trackers/$SUENO/records" "{\"value\":5,\"date\":\"$TODAY\"}"
check "min 5/6 -> no cumple" "$(json '.isMet' < "$OUT.body")" "false"
check "min 5/6 -> desviación 1" "$(json '.deviation' < "$OUT.body")" "1"

# upsert: corregir el valor del día sin duplicar
req PUT "/api/trackers/$SUENO/records" "{\"value\":6,\"date\":\"$TODAY\"}"
check "PUT upsert -> 2xx" "$([[ "$RC" == 20* ]] && echo 2xx)" "2xx"
check "upsert 6/6 -> isMet" "$(json '.isMet' < "$OUT.body")" "true"

req GET "/api/trackers/$SUENO/history"
check "GET history -> 200" "$RC" "200"
check "history 2 registros" "$(json 'length' < "$OUT.body")" "2"
check "registro del día corregido a 6" "$(json '[.[] | select(.entryDate=="'"$TODAY"'") | .value] | length == 1 and .[0]==6' < "$OUT.body")" "true"

req GET /api/trackers
check "GET list trackers -> 200" "$RC" "200"
check "lista contiene ambos" "$(json --arg a "$PESO" --arg b "$SUENO" '[.[] | select(.id==($a|tonumber) or .id==($b|tonumber))] | length' < "$OUT.body")" "2"

req DELETE "/api/trackers/$PESO"
check "DELETE tracker -> 204" "$RC" "204"
req GET "/api/trackers/$PESO"
check "GET tras DELETE -> 404" "$RC" "404"

# ---------- 7. Métricas ----------
echo "[12] GET /api/metrics"
req GET /api/metrics
check "status 200" "$RC" "200"
check "tiene monthly" "$(json 'has("monthly")' < "$OUT.body")" "true"
check "tiene monthlySeries" "$(json 'has("monthlySeries")' < "$OUT.body")" "true"
check "tiene trackers" "$(json 'has("trackers")' < "$OUT.body")" "true"
req GET '/api/metrics?month=13&year=2026'
check "mes inválido -> 400" "$RC" "400"

# ---------- 8. Logout ----------
echo "[13] POST /api/logout"
OUT="$(mktemp)"
curl -s -X POST "$BASE/api/logout" -H 'Cookie: mtt_session='"$TOKEN" -D - -o "${OUT}.body" > "$OUT"
check "status 200" "$(status_of "$OUT")" "200"
check "cookie borrada (Max-Age=0)" "$(grep -i 'set-cookie:.*mtt_session' "$OUT" | grep -o 'Max-Age=[0-9-]*' | head -1 | cut -d= -f2)" "0"

# ---------- Limpieza ----------
echo "[*] Limpieza de datos de prueba"
req DELETE "/api/immediate-tasks/$IT_ID"
req DELETE "/api/trackers/$SUENO"
req DELETE "/api/monthly-tasks/$MT_ID"
rm -f "$JAR"

# Si hay Docker, también limpiar historial órfano y el usuario de test (best effort)
if command -v docker >/dev/null 2>&1 && docker ps --format '{{.Names}}' | grep -q '^myticktick_postgres$'; then
  docker exec myticktick_postgres psql -U myticktick -d myticktick -qc \
    "delete from monthly_task_history_dbs where monthly_task_id in ($MT_ID);
     delete from tracker_entry_dbs where tracker_id in ($PESO, $SUENO);
     delete from user_dbs where username='$USERNAME';" >/dev/null 2>&1 \
    && echo "  DB: historial/entradas/usuario de test eliminados" || true
fi

echo
echo "=== RESULTADO: $PASS passed, $FAIL failed ==="
[ "$FAIL" -eq 0 ]
