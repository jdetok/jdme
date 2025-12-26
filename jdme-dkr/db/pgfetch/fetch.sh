#!/bin/bash

LOGD=/app/log
LOGF=$LOGD/fetch_$(date +'%m%d%y_%H%M%S').log
EXEC=bin/fetch
PROC=scripts/dly.sql
BACKUP_FILE="/dump/bball_$(date +%m%d%Y).sql"

export PATH=$PATH:/usr/local/go/bin
export PGPASSWORD=$PG_PASS

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++ DAILY BBALL ETL STARTED
++ $(date)
++ LOGFILE: $LOGF
" >> $LOGF

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++ RUNNING GO ETL CLI APPLICATION
" >> $LOGF

# run go app
# ./$EXEC -envf skip -mode daily || exit 1
# ./"$EXEC" -envf skip -mode daily 2>&1 | tee -a "$LOGF"
./"$EXEC" -envf skip -mode custom -szn 2025 2>&1 | tee -a "$LOGF"

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++ GO ETL CLI APPLICATION RAN SUCCESSFULLY
++ $(date)
++ ATTEMPTING NIGHTLY POSTGRES PROCEDURES TO UPDATE API TABLES
" >> $LOGF

# update sql views and tables with new stats data
# psql -h postgres -U postgres -d $PG_DB < ./$PROC 2>&1 || exit 1
psql -h $PG_HOST -U $PG_USER -d $PG_DB -v ON_ERROR_STOP=1 -f $PROC 2>&1 | tee -a "$LOGF"

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++ RAN PSQL PROCEDURES, CREATING DB BACKUP
++ $(date)" >> $LOGF

# run pg_dump and compress
pg_dump -h $PG_HOST -U $PG_USER $PG_DB | gzip > "$BACKUP_FILE.gz"

# email log
./$EXEC -mode email -logf $LOGF || exit 1

echo "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
++ EMAIL SENT - SCRIPT COMPLETE
++ $(date)" >> $LOGF