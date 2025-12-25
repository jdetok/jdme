/*
to be called nightly after go etl process
*/
create or replace procedure sp_nightly_call()
language plpgsql
as $$
begin
	-- load seasons
	raise notice e'inserting seasons...\n';
	call lg.sp_szn_load();
	raise notice e'seasons insert complete: %\n', fn_cntstr('lg.szn');

	-- load all teams
	raise notice e'inserting all nba/wnba teams...\n';
	call lg.sp_team_all_load();
	raise notice e'team insert complete: %\n', fn_cntstr('lg.team');

	-- load tbox table with team box stats
	raise notice e'inserting team box stats into stats.tbox...\n';
	call stats.sp_tbox();
	raise notice e'tbox insert complete: %\n', fn_cntstr('stats.tbox');

	-- load all players
	raise notice e'inserting all nba/wnba players...\n';
	call lg.sp_plr_all_load();
	raise notice e'player insert complete: %\n', fn_cntstr('lg.plr');

	-- load pbox table with player box scores after inserting player causing issue
	raise notice e'inserting player box stats into stats.pbox...\n';
	call stats.sp_pbox();
	raise notice e'pbox insert complete: %\n', fn_cntstr('stats.pbox');

	-- load api.plr_agg table with pbox stats 
	raise notice e'inserting season/career stat aggregations into api.plr_agg...\n';
	call api.sp_plr_agg();
	raise notice e'player agg insert complete: %\n', fn_cntstr('api.plr_agg');
end; $$;
-- call sp_nightly_call();