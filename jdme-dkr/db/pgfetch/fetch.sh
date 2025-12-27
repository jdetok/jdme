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

usage() {
  echo "Usage: $0 [-s szn]"
  echo "  -s <season>   Run fetch for a specific season (e.g. 2025)"
  exit 1
}

while [[ $# > 0 ]]; do
  case "$1" in
    -s|-szn)
      SEASON="$2"
      MODE="custom"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++ DAILY BBALL ETL STARTED
++ $(date)
++ LOGFILE: $LOGF
" | tee -a "$LOGF"

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++ RUNNING GO ETL CLI APPLICATION
" | tee -a "$LOGF"

# run go app
if [[ -n "$SEASON" ]]; then
    ./"$EXEC" -envf skip -mode $MODE -szn $SEASON 2>&1 | tee -a "$LOGF"
else
    ./"$EXEC" -envf skip -mode $MODE 2>&1 | tee -a "$LOGF"
fi

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++ GO ETL CLI APPLICATION RAN SUCCESSFULLY
++ $(date)
++ ATTEMPTING NIGHTLY POSTGRES PROCEDURES TO UPDATE API TABLES
" | tee -a "$LOGF"

# update sql views and tables with new stats data
# psql -h postgres -U postgres -d $PG_DB < ./$PROC 2>&1 || exit 1
psql -h $PG_HOST -U $PG_USER -d $PG_DB -v ON_ERROR_STOP=1 -f $PROC 2>&1 | tee -a "$LOGF"

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++ RAN PSQL PROCEDURES, CREATING DB BACKUP
++ $(date)" | tee -a "$LOGF"

# run pg_dump and compress
pg_dump -h $PG_HOST -U $PG_USER $PG_DB | gzip > "$BACKUP_FILE.gz"

# email log
./$EXEC -mode email -logf $LOGF || exit 1

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++ EMAIL SENT - SCRIPT COMPLETE
++ $(date)" | tee -a "$LOGF"