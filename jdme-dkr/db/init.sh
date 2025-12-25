#!/bin/bash
set -e

BKPDIR="/var/lib/postgresql/dump"

if [ -z $(ls -A $PGDATA) ]; then
    echo "data directory empty, restoring from backup..."

    LATEST_BKP=$(ls -t $BKPDIR/*.sql 2>/dev/null | head -1)

    if [ -z $LATEST_BKP ]; then
        echo "no backups found in $BKPDIR, building from entrypoint"
    else
        echo "restoring backup $LATEST_BKP"
        if [[ $LATEST_BKP == *.gz ]]; then
            gunzip -c $LATEST_BKP | psql -U postgres
        else
            psql -U postgres -f $LATEST_BKP
        fi
    fi
fi

exec docker-entrypoint.sh "$@"