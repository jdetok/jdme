-- docker exec -i mdb mariadb nba < /mnt/web/jdeko.me/sql-final/procedures/run_nightly_procs.sql

CALL sp_insert_career_totals();
CALL sp_insert_career_avgs();
CALL sp_insert_top_scorers();
CALL sp_insert_season_totals();
CALL sp_insert_season_avgs();
