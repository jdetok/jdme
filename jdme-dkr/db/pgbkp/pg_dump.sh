#!/bin/bash
# dump Postgres database to mounted volume
export PGPASSWORD=$POSTGRES_PASSWORD

# timestamped backup file
BACKUP_FILE="/dump/bball_$(date +%m%d%Y).sql"

# run pg_dump and compress
# pg_dump -h postgres -U $POSTGRES_USER $POSTGRES_DB | gzip > "$BACKUP_FILE"
pg_dump -h postgres -U $POSTGRES_USER $POSTGRES_DB > $BACKUP_FILE

# optional: remove backups older than 14 days
find /dump -type f -mtime +14 -delete
