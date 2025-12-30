#!/bin/bash

LOGD=/app/log
LOGF=$LOGD/fetch_$(date +'%m%d%y_%H%M%S').log
EXEC=bin/fetch
PROC=scripts/dly.sql
BACKUP_FILE="/dump/bball_$(date +%m%d%Y_%H%M%S).sql"

export PATH=$PATH:/usr/local/go/bin
export PGPASSWORD=$PG_PASS

set -euo pipefail

# defaults
SEASON=""
MODE="daily"
DB_DUMP=1
SEND_EMAIL=1

usage() {
  echo "CLI options: $0"
  echo "| -s <season>    | fetch stats for a specific season (e.g. 2025)"
  echo "| -nd|--no-dump  | Skip pg_dump backup"
  echo "| -ne|--no-email | Don't send confirmation email"
  exit 1
}

while [[ $# -gt 0 ]]; do
   case "$1" in
    -s|-szn)
      SEASON="$2"
      MODE="custom"
      shift 2
      ;;
    --no-dump|-nd)
      DB_DUMP=0
      shift
      ;;
    --no-email|-ne)
      SEND_EMAIL=0
      shift
      ;;  
    -h|--help)
      usage
      ;;
    *)
      echo "Unknown option: $1"
      usage
      ;;
  esac
done

echo -e "++ $(date) | DAILY BBALL ETL STARTED\n++++ LOGGING TO: $LOGF\n" | tee -a $LOGF

echo "++ $(date) | RUNNING GO ETL CLI APPLICATION" | tee -a $LOGF
if [[ -n "$SEASON" ]]; then
    ./"$EXEC" -envf skip -logf cli -mode $MODE -szn $SEASON 2>&1 | tee -a $LOGF
else
    ./"$EXEC" -envf skip -logf cli -mode $MODE 2>&1 | tee -a $LOGF
fi
echo -e "++ $(date) | GO ETL CLI APPLICATION RAN SUCCESSFULLY\n" | tee -a $LOGF

echo "++ $(date) | RUNNING PSQL PROCEDURES IN $PROC" | tee -a $LOGF
psql -h $PG_HOST -U $PG_USER -d $PG_DB -v ON_ERROR_STOP=1 -f $PROC 2>&1 | tee -a $LOGF
echo -e "++ $(date) | RAN PSQL PROCEDURES\n" | tee -a $LOGF

# run pg_dump and compress
if [[ $DB_DUMP == 1 ]]; then
  echo "++ $(date) | CREATING DB BACKUP" | tee -a $LOGF
  pg_dump -h $PG_HOST -U $PG_USER $PG_DB | gzip > "$BACKUP_FILE.gz"
  echo -e "++ $(date) | DB BACKUP CREATED AT $BACKUP_FILE\n" | tee -a $LOGF
else
  echo -e "++ $(date) | SKIPPING DB BACKUP (--no-dump)\n" | tee -a $LOGF
fi

# email log
if [[ $SEND_EMAIL == 1 ]]; then
  echo "++ $(date) | SENDING EMAIL"
  ./$EXEC -mode email -attach $LOGF -logf cli | tee -a $LOGF
  echo -e "++ $(date) | EMAIL SENT\n"
else
  echo -e "++ $(date) | SKIPPING EMAIL (--no-email)\n"
fi

echo "++ $(date) | SCRIPT COMPLETE" | tee -a $LOGF