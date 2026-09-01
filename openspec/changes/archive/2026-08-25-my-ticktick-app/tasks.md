## 0. Fase de Exploración (MVP)

### 0.1 Setup Inicial

- [x] 0.1.1 Setup rápido: inicializar proyecto Go con estructura básica
- [x] 0.1.2 Setup rápido: levantar PostgreSQL en Docker (docker-compose mínimo)
- [x] 0.1.3 Definir arquitectura: crear diagrama de arquitectura del proyecto

### 0.2 Explorar Frontend

- [x] 0.2.1 Investigar: comparar Vanilla JS vs React vs Vue para tu caso
- [x] 0.2.2 Prototipar: crear página simple con Vanilla JS + CSS
- [x] 0.2.3 Prototipar: crear misma página con React (para comparar)
- [x] 0.2.4 Decidir: elegir framework basado en los prototipos
- [x] 0.2.5 Definir UI: crear mockup básico de las 3 vistas principales

### 0.3 Explorar Login Simple

- [x] 0.3.1 Investigar: entender autenticación básica en Go (usuario/contraseña)
- [x] 0.3.2 Definir modelo: crear modelo User con username único y password hash (bcrypt)
- [x] 0.3.3 Implementar: formulario de registro (sign up) con username + password
- [x] 0.3.4 Implementar: formulario de login con "Recordarme"
- [x] 0.3.5 Implementar: endpoint /register que guarda el usuario con password hasheado
- [x] 0.3.6 Implementar: endpoint /login que emite JWT en cookie httpOnly (expiración larga con "Recordarme", corta sin)
- [x] 0.3.7 Implementar: endpoint /logout que borra la cookie de sesión
- [x] 0.3.8 Documentar: anotar que Google OAuth es pendiente para futuro
- [x] 0.3.9 Implementar: middleware de autenticación que protege todas las rutas /api/* (salvo /api/health, /api/register, /api/login, /api/logout) y devuelve 401 sin sesión

### 0.4 MVP Básico

- [x] 0.4.1 Definir modelos: crear modelos Go básicos para MonthlyTask, ImmediateTask, DailyTracker
- [x] 0.4.2 MVP backend: implementar endpoint GET /api/health
- [x] 0.4.3 MVP frontend: crear página de prueba que conecte con el backend
- [x] 0.4.4 Integrar: conectar login simple (usuario/contraseña) al MVP y proteger rutas

### 0.5 Iterar y Documentar

- [x] 0.5.1 Iterar: revisar arquitectura y ajustar si es necesario
- [x] 0.5.2 Iterar: revisar estética UI y ajustar si es necesario
- [x] 0.5.3 Documentar: anotar decisiones de arquitectura y estilo para la fase de implementación

## 1. Setup de Infraestructura

- [x] 1.1 Crear Dockerfile para el backend en Go
- [x] 1.2 Crear docker-compose.yml para PostgreSQL
- [x] 1.3 Configurar conexión a base de datos en Go
- [x] 1.4 Crear scripts de migración de base de datos

## 2. Modelos de Datos

- [x] 2.1 Implementar modelo MonthlyTask con GORM
- [x] 2.2 Implementar modelo MonthlyTaskHistory con GORM
- [x] 2.3 Implementar modelo ImmediateTask con GORM
- [x] 2.4 Implementar modelo DailyTracker con GORM (limitValue + limitType min/max)
- [x] 2.5 Implementar modelo TrackerEntry con GORM
- [x] 2.6 Crear migrations para todos los modelos

## 3. API REST - Tareas Mensuales

- [x] 3.1 Implementar endpoint POST /api/monthly-tasks
- [x] 3.2 Implementar endpoint GET /api/monthly-tasks
- [x] 3.3 Implementar endpoint GET /api/monthly-tasks/{id}
- [x] 3.4 Implementar endpoint PUT /api/monthly-tasks/{id}
- [x] 3.5 Implementar endpoint DELETE /api/monthly-tasks/{id}
- [x] 3.6 Implementar endpoint GET /api/monthly-tasks/{id}/history

## 4. API REST - Tareas Inmediatas

- [x] 4.1 Implementar endpoint POST /api/immediate-tasks
- [x] 4.2 Implementar endpoint GET /api/immediate-tasks
- [x] 4.3 Implementar endpoint GET /api/immediate-tasks/{id}
- [x] 4.4 Implementar endpoint PUT /api/immediate-tasks/{id}
- [x] 4.5 Implementar endpoint DELETE /api/immediate-tasks/{id}

## 5. API REST - Trackers Diarios

- [x] 5.1 Implementar endpoint POST /api/trackers
- [x] 5.2 Implementar endpoint GET /api/trackers
- [x] 5.3 Implementar endpoint GET /api/trackers/{id}
- [x] 5.4 Implementar endpoint PUT /api/trackers/{id}
- [x] 5.5 Implementar endpoint DELETE /api/trackers/{id}
- [x] 5.6 Implementar endpoint POST /api/trackers/{id}/records
- [x] 5.7 Implementar endpoint GET /api/trackers/{id}/history

## 6. API REST - Lógica de Negocio

- [x] 6.1 Implementar activación automática de tareas mensuales el día 1
- [x] 6.2 Implementar cálculo de cumplimiento contra límite unilateral (mín o máx) en trackers
- [x] 6.3 Implementar cálculo de desviación en tracker records
- [x] 6.4 Implementar agregación de métricas de cumplimiento

## 7. Frontend - Estructura

- [x] 7.1 Crear estructura básica de web app
- [x] 7.2 Configurar routing para diferentes secciones
- [x] 7.3 Implementar layout responsive base
- [x] 7.4 Configurar conexión con API REST
- [x] 7.5 Implementar vistas de login y registro con "Recordarme" y logout

## 8. Frontend - Tareas Mensuales

- [x] 8.1 Implementar vista de lista de tareas mensuales
- [x] 8.2 Implementar formulario de creación de tarea mensual
- [x] 8.3 Implementar formulario de edición de tarea mensual
- [x] 8.4 Implementar vista de historial de tarea mensual
- [x] 8.5 Implementar marcar como completada con toggle o swipe (sin confirmación) + "Deshacer"

## 9. Frontend - Tareas Inmediatas

- [x] 9.1 Implementar vista de lista de tareas inmediatas
- [x] 9.2 Implementar formulario de creación de tarea inmediata
- [x] 9.3 Implementar selector de fecha límite
- [x] 9.4 Implementar selector de prioridad
- [x] 9.5 Implementar ordenamiento por fecha límite
- [x] 9.6 Implementar completar/eliminar por toggle o swipe (sin confirmación) + toast "Deshacer"

## 10. Frontend - Trackers Diarios

- [x] 10.1 Implementar vista de lista de trackers
- [x] 10.2 Implementar formulario de creación de tracker
- [x] 10.3 Implementar configuración de límite (mínimo o máximo) y su valor en el tracker
- [x] 10.4 Implementar registro de valor diario
- [x] 10.5 Implementar visualización de cumplimiento/no cumplimiento
- [x] 10.6 Implementar visualización de desviación (cuánto superó el límite)
- [x] 10.7 Implementar toggle/swipe + "Deshacer" para registrar y revertir valores diarios

## 11. Frontend - Dashboard

- [x] 11.1 Implementar vista de dashboard en desktop
- [x] 11.2 Integrar librería de gráficos (Chart.js o similar)
- [x] 11.3 Implementar gráfico de cumplimiento mensual de tareas
- [x] 11.4 Implementar gráfico de evolución de trackers
- [x] 11.5 Implementar filtros de fecha para gráficos
- [x] 11.6 Implementar vista de métricas de cumplimiento

## 12. Testing

- [x] 12.1 Testear endpoints de API con datos de prueba
- [x] 12.2 Testear frontend en desktop
- [x] 12.3 Testear frontend en móvil
- [x] 12.4 Testear dashboard de gráficos
- [x] 12.5 Testear activación automática de tareas mensuales
- [x] 12.6 Testear límite unilateral (mín/máx) y desviación
- [x] 12.7 Testear autenticación: registro, login, recordarme, logout y 401 en rutas protegidas

## 13. Despliegue

- [x] 13.1 Configurar docker-compose para producción (docker-compose.prod.yml)
- [x] 13.2 Documentar variables de entorno (README + .env.example)
- [x] 13.3 Crear Makefile de despliegue (up, down, logs, restart, status, gen-secret) reemplazando start.sh
- [x] 13.4 Documentar proceso de backup de base de datos (BACKUP.md)
- [x] 13.5 Documentar servicio systemd para autarranque (README)
