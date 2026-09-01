# Arquitectura - MyTickTick

## Diagrama de Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                         CLIENTE                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                 Web App (Responsive)                │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │   │
│  │  │  Móvil   │  │ Desktop  │  │   Dashboard      │ │   │
│  │  │  Vista   │  │  Vista   │  │   (Gráficos)     │ │   │
│  │  └──────────┘  └──────────┘  └──────────────────┘ │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │ HTTPS
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                         BACKEND                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   Go Server                         │   │
│  │  ┌──────────────────────────────────────────────┐  │   │
│  │  │  HTTP Handlers (net/http)                    │  │   │
│  │  │  /api/health                                  │  │   │
│  │  │  /api/monthly-tasks                          │  │   │
│  │  │  /api/immediate-tasks                        │  │   │
│  │  │  /api/trackers                               │  │   │
│  │  └──────────────────────────────────────────────┘  │   │
│  │                            │                        │   │
│  │  ┌──────────────────────────────────────────────┐  │   │
│  │  │  GORM ORM                                    │  │   │
│  │  │  - MonthlyTask                               │  │   │
│  │  │  - MonthlyTaskHistory                        │  │   │
│  │  │  - ImmediateTask                             │  │   │
│  │  │  - DailyTracker                              │  │   │
│  │  │  - TrackerRecord                             │  │   │
│  │  └──────────────────────────────────────────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │ PostgreSQL Protocol
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      BASE DE DATOS                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              PostgreSQL 16 (Docker)                 │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐│   │
│  │  │  monthly_   │  │ immediate_  │  │  daily_     ││   │
│  │  │  tasks      │  │ tasks       │  │ trackers    ││   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘│   │
│  │  ┌─────────────┐  ┌─────────────┐                 ││   │
│  │  │ monthly_    │  │ tracker_    │                 ││   │
│  │  │ task_history│  │ records     │                 ││   │
│  │  └─────────────┘  └─────────────┘                 ││   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Componentes

### Frontend
- **Web App Responsive**: Una sola base de código para móvil y desktop
- **Vistas**:
  - Tareas mensuales
  - Tareas inmediatas
  - Trackers diarios
  - Dashboard (solo desktop)

### Backend
- **Lenguaje**: Go
- **ORM**: GORM
- **API**: REST
- **Puerto**: 8080 (interno del contenedor; expuesto en el host en 8090)

### Base de Datos
- **Sistema**: PostgreSQL 16
- **Infraestructura**: Docker container
- **Persistencia**: Volume de Docker
- **Puerto**: 5432

## Flujo de Datos

1. **Usuario** interactúa con la web app
2. **Web App** hace peticiones HTTP al backend Go
3. **Backend Go** procesa la petición y usa GORM para acceder a la BD
4. **PostgreSQL** almacena/recupera los datos
5. **Respuesta** viaja de vuelta al usuario

## Tecnologías

| Capa | Tecnología |
|------|-----------|
| Frontend | HTML/CSS/JS (Vanilla o framework) |
| Backend | Go 1.21+ |
| ORM | GORM |
| Base de Datos | PostgreSQL 16 |
| Infraestructura | Docker, docker-compose |
