# MyTickTick Makefile
# Reemplaza start.sh. Uso: make <target>

COMPOSE_PROD := docker compose -f docker-compose.prod.yml
COMPOSE_DEV  := docker compose -f docker-compose.yml

# El secreto JWT vive FUERA del repo (en tu home).
# SECR es solo una ruta (no se lee a parse-time). Cada target de prod la
# inyecta al shell SOLO si el archivo existe (ver la receta _prod-env).
SECR := $(HOME)/.myticktick/token-secret
# Inyecta MYTICKTICK_TOKEN_SECRET al entorno del recipe si el archivo existe.
# Si el usuario ya la definió en el entorno, se respeta (no se sobreescribe).
PROD_ENV := if [ -z "$$MYTICKTICK_TOKEN_SECRET" ] && [ -f "$(SECR)" ]; then export MYTICKTICK_TOKEN_SECRET=$$(cat "$(SECR)"); fi

.DEFAULT_GOAL := help

## --- Producción (docker-compose.prod.yml) ---

## Levantar servicios de producción (auto-genera el secreto si no existe)
up:
	@SECR="$(SECR)"; \
	if [ -z "$$MYTICKTICK_TOKEN_SECRET" ] && [ ! -f "$$SECR" ]; then \
		echo "Generando MYTICKTICK_TOKEN_SECRET en $$SECR ..."; \
		mkdir -p "$$(dirname "$$SECR")" && openssl rand -hex 32 > "$$SECR" && chmod 600 "$$SECR"; \
	fi; \
	if [ -f "$$SECR" ]; then export MYTICKTICK_TOKEN_SECRET=$$(cat "$$SECR"); fi; \
	if [ -z "$$MYTICKTICK_TOKEN_SECRET" ]; then echo "ERROR: no se pudo resolver el secreto JWT (make gen-secret)"; exit 1; fi; \
	echo "Iniciando MyTickTick (producción)..."; \
	$(COMPOSE_PROD) up -d --build; \
	echo; \
	$(COMPOSE_PROD) ps; \
	echo "API + Frontend: http://localhost:$${API_PORT:-8090}"; \
	echo "Logs: make logs"

## Detener servicios de producción (datos NO se borran)
down:
	@echo "Deteniendo servicios de producción..."
	@$(PROD_ENV); $(COMPOSE_PROD) down
	@echo "Para borrar la BD: docker volume rm myticktick_postgres_data"

## Ver logs en tiempo real
logs:
	@$(PROD_ENV); $(COMPOSE_PROD) logs -f --tail=100

## Reiniciar contenedores
restart:
	@$(PROD_ENV); $(COMPOSE_PROD) restart

## Ver estado
status:
	@$(PROD_ENV); $(COMPOSE_PROD) ps

## Generar (o regenerar) el secreto JWT en ~/.myticktick/token-secret
## Regenerar cambia el valor → todas las sesiones expiran.
gen-secret:
	@mkdir -p "$(dir $(SECR))"
	@openssl rand -hex 32 > "$(SECR)"
	@chmod 600 "$(SECR)"
	@echo "Secreto guardado en $(SECR)"
	@echo "make up lo inyecta solo; no hace falta ponerlo en .env."

## --- Desarrollo (docker-compose.yml) ---

## Levantar stack de desarrollo (Postgres + API, usa el secreto dev por defecto)
dev:
	$(COMPOSE_DEV) up -d --build
	@echo "Dev stack listo: http://localhost:8080"

## Detener stack de desarrollo
dev-down:
	$(COMPOSE_DEV) down

## --- Utilidades ---

## Build sin correr
build:
	@$(PROD_ENV); $(COMPOSE_PROD) build

## Compilar el cliente TUI sin correrlo
tui-build:
	go build -o bin/myticktick-tui ./cmd/tui

## Compilar y correr el cliente TUI contra el server (MTT_URL, default http://localhost:8080)
tui:
	go build -o bin/myticktick-tui ./cmd/tui
	./bin/myticktick-tui

## Ver help
help:
	@echo "MyTickTick — targets disponibles:"
	@echo
	@echo "  Producción:"
	@echo "    make up        Levantar (auto-genera el secreto si no existe)"
	@echo "    make down      Detener"
	@echo "    make logs      Logs en tiempo real"
	@echo "    make restart   Reiniciar"
	@echo "    make status    Ver estado"
	@echo "    make gen-secret  Generar/regenerar el secreto JWT (~/.myticktick/token-secret)"
	@echo
	@echo "  Desarrollo:"
	@echo "    make dev       Levantar stack dev"
	@echo "    make dev-down  Detener stack dev"
	@echo
	@echo "  Utilidades:"
	@echo "    make build     Build sin correr"
	@echo "    make tui       Compilar y correr el cliente TUI"
	@echo "    make tui-build Solo compilar el cliente TUI"
	@echo "    make help      Esta ayuda"

.PHONY: up down logs restart status gen-secret dev dev-down build tui tui-build help
