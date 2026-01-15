#!/bin/bash
set -e

BKPDIR="/var/lib/postgresql/dump"
SQLDIR="/var/lib/postgresql/sql"

# Defaults
FORCE_DUMP=1

usage() {
  echo "Usage: $0 [--from-dump]"
  echo "  --from-dump  Restore from latest dump regardless of PGDATA"
  exit 1
}

# Parse options
while [[ $# -gt 0 ]]; do
    case "$1" in
        --from-dump)
            FORCE_DUMP=1
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            # leave remaining args for postgres entrypoint
            break
            ;;
    esac
done

restore_from_dump() {
    echo "Restoring Postgres from latest backup..."
    LATEST_BKP=$(ls -t "$BKPDIR"/*.sql* 2>/dev/null | head -1)

    if [ -z "$LATEST_BKP" ]; then
        echo "No backup found in $BKPDIR"
        return 1
    fi

    echo "Using backup: $LATEST_BKP"
    export PGPASSWORD="$POSTGRES_PASSWORD"

    if [[ "$LATEST_BKP" == *.gz ]]; then
        gunzip -c "$LATEST_BKP" | psql -U "$POSTGRES_USER" "$POSTGRES_DB"
    else
        psql -U "$POSTGRES_USER" "$POSTGRES_DB" -f "$LATEST_BKP"
    fi
}

build_from_sql() {
    echo "Building Postgres database from $SQLDIR..."
    export PGPASSWORD="$POSTGRES_PASSWORD"

    for f in "$SQLDIR"/*.sql; do
        [ -f "$f" ] || continue
        echo "Running $f..."
        psql -U "$POSTGRES_USER" "$POSTGRES_DB" -f "$f"
    done
}

# Main logic
if [ "$FORCE_DUMP" -eq 1 ]; then
    echo "--from-dump flag detected"
    restore_from_dump
else
    if [ -z "$(ls -A "$PGDATA")" ]; then
        echo "PGDATA is empty, first initialization"
        if ! restore_from_dump; then
            build_from_sql
        fi
    else
        echo "PGDATA exists, starting Postgres normally"
    fi
fi

# Finally, exec the official Postgres entrypoint with remaining args
exec docker-entrypoint.sh "$@"
