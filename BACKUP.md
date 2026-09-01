# Backup de Base de Datos

## Script de Backup

```bash
# Backup manual
docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U postgres myticktick > backup_$(date +%Y%m%d).sql

# Backup automático (cron)
0 2 * * * docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U postgres myticktick > /backups/backup_$(date +%Y%m%d).sql
```

## Restaurar Backup

```bash
docker-compose -f docker-compose.prod.yml exec -T postgres psql -U postgres -d myticktick < backup_20260820.sql
```

## Notas
- Los backups se guardan en el volumen `postgres_data`
- Para backups persistentes, montar un volumen adicional en `/backups`
- Retener backups de los últimos 30 días como mínimo
