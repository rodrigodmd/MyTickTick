# MyTickTick

Tracker personal de tareas y métricas diarias. Backend Go + GORM + PostgreSQL, frontend Vanilla JS servido por el mismo binario.

## Estructura del proyecto

```
MyTickTick/
├── cmd/
│   └── api/
│       └── main.go                  # Punto de entrada del servidor HTTP
├── internal/
│   ├── core/                        # Núcleo (Domain + Ports + Services)
│   │   ├── domain/
│   │   │   └── entities.go          # Entidades de negocio
│   │   ├── ports/
│   │   │   ├── repositories.go      # Interfaces de repositorios
│   │   │   └── services.go          # Interfaces de servicios
│   │   ├── services/
│   │   │   ├── auth/                # Servicio de autenticación
│   │   │   ├── compliance/          # Cálculo de cumplimiento
│   │   │   ├── immediatetask/       # Tareas inmediatas
│   │   │   ├── monthlytasks/        # Tareas mensuales
│   │   │   └── trackers/            # Trackers diarios
│   │   └── token/                   # JWT (HMAC-SHA256 stdlib)
│   ├── handlers/                    # Adaptadores de entrada (HTTP)
│   │   ├── auth/                    # Auth + middleware + cookies
│   │   ├── compliance/              # Métricas de cumplimiento
│   │   ├── immediatetask/           # Tareas inmediatas
│   │   ├── immediatetasks/          # (legacy)
│   │   ├── monthlytasks/            # Tareas mensuales
│   │   ├── static/                  # Servidor de web app
│   │   └── trackers/                # Trackers diarios
│   └── repositories/                # Adaptadores de salida (GORM)
│       ├── database/                # Conexión + migraciones
│       ├── user/
│       ├── immediatetask/
│       ├── monthlytask/
│       ├── monthlytaskhistory/
│       ├── dailytracker/
│       └── trackerentry/
├── web/                             # Frontend servido (HTML/CSS/JS)
│   ├── html/                        # Páginas: index, dashboard, tasks, trackers, login
│   ├── css/
│   └── js/
├── docs/                            # Decisiones de diseño
├── test/                            # Scripts de prueba API
├── cmd/api/main.go                  # Binario
├── go.mod
├── Dockerfile                       # Multi-stage build (Go → Alpine)
├── docker-compose.yml               # Stack de desarrollo
├── docker-compose.prod.yml          # Stack de producción
├── .env.example                     # Variables de entorno (producción)
├── Makefile                         # Comandos de despliegue
└── BACKUP.md                        # Proceso de backup
```

## Cómo levantarlo

### Desarrollo

Requiere: Docker, Docker Compose.

```bash
# Levanta PostgreSQL + API en modo dev
make dev

# Acceder
open http://localhost:8080

# Detener
make dev-down
```

Variables de entorno para dev (defaults en `database.go`):
- `DB_HOST` (default: localhost)
- `DB_PORT` (default: 5432)
- `DB_USER` (default: myticktick)
- `DB_PASSWORD` (default: myticktick)
- `DB_NAME` (default: myticktick)
- `MYTICKTICK_TOKEN_SECRET` (default: valor de dev — no usar en prod)
- `PORT` (default: 8080)

### Producción

Requiere: Docker, Docker Compose.

```bash
# 1. (Opcional) Generar el secreto JWT en ~/.myticktick/token-secret
#    `make up` lo genera solo la primera vez si no existe, así que puedes saltar este paso.
make gen-secret

# 2. (Opcional) Crear .env solo si querés overrides de puertos/DB
#    cp .env.example .env
#    Por defecto no hace falta: los defaults del compose ya sirven.

# 3. Levantar (build + up)
make up

# Verificar
curl http://localhost:8090/api/health
# → {"status":"ok"}

# Detener (datos persisten en volumen)
make down

# Logs
make logs
```

**Secreto JWT** (`MYTICKTICK_TOKEN_SECRET`): es la clave HMAC con la que se firman tus cookies de sesión. Es la única variable sin default. Por defecto lo guarda `make gen-secret` en `~/.myticktick/token-secret` (fuera del repo, chmod 600) y `make up` lo inyecta solo. Debe ser estable entre reinicios: si lo regenerás, todas las sesiones expiran.

**Variables de entorno** (opcionales, en `.env` si las querés cambiar; el `:-` en `docker-compose.prod.yml` da los defaults):

- `MYTICKTICK_TOKEN_SECRET` — ver arriba. Si la ponés en `.env` se usa en vez del archivo de home.
- `API_PORT` — Puerto expuesto en el host (default: 8090)
- `POSTGRES_USER` — Usuario de la DB (default: myticktick)
- `POSTGRES_PASSWORD` — Contraseña de la DB (default: myticktick)
- `POSTGRES_DB` — Nombre de la BD (default: myticktick)
- `POSTGRES_PORT` — Puerto expuesto de PostgreSQL (default: 5432)
- `BACKUP_DIR` — Directorio para dumps SQL (opcional, ver BACKUP.md)

### Servicio de boot (systemd)

Para que MyTickTick se levante solo al arrancar la máquina, instalá una unidad systemd:

1. Crear el archivo `/etc/systemd/system/myticktick.service`:

```ini
[Unit]
Description=MyTickTick (API + PostgreSQL)
After=network-online.target docker.service
Wants=network-online.target
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
User=rodri
WorkingDirectory=/home/rodri/Projects/MyTickTick
ExecStart=/usr/bin/make up
ExecStop=/usr/bin/make down
Environment=HOME=/home/rodri
Environment=PATH=/usr/local/bin:/usr/bin:/bin

[Install]
WantedBy=multi-user.target
```

2. Activar:
```bash
sudo systemctl daemon-reload
sudo systemctl enable myticktick
sudo systemctl start myticktick
```

3. Verificar:
```bash
systemctl status myticktick
curl http://localhost:8090/api/health
# → {"status":"ok"}
```

4. Desactivar (si querés quitarlo del boot):
```bash
sudo systemctl disable myticktick
sudo systemctl stop myticktick
sudo rm /etc/systemd/system/myticktick.service
sudo systemctl daemon-reload
```

**Notas:**
- `Type=oneshot` + `RemainAfterExit=yes`: `make up` arranca los contenedores (que corren por su cuenta) y systemd marca el servicio como "active".
- Los contenedores tienen `restart: unless-stopped` en el compose, así que si uno muere Docker lo reinicia.
- `make down` (ExecStop) detiene los contenedores sin borrar los datos del volumen.

### Backups

Ver [BACKUP.md](BACKUP.md).
