/*
to be called at the end of db init (in docker entrypoint dir)
sp_rebuild() also exists, which is identical but truncates tables first
	- used for testing build sequence on an existing database
*/
create or replace procedure sp_build()
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

	/* INSERT A ROW INTO lg.plr FOR WNBA PLAYER ANGEL ROBINSON WITH PLAYER ID 
	202270 this player had the ID 202270 in 2014 and 202657 in all years after
	this was causing an error with loading stats.pbox table
	create a new record with identical data except player id
	MUST BE RUN AFTER THE lg.sp_player_all_load*/
	raise notice e'inserting 202270 copy of 202657, won''t work without...\n';
	insert into lg.plr 
		(lg_id, player_id, plr_cde, player, last_first, from_year, to_year)
	select
		1, 202270, 
		playercode, display_first_last, display_last_comma_first,
		from_year, to_year
	from intake.wplayer
	where person_id = 202657;
	raise notice e'player id 202270 insert complete: %\n', fn_cntstr('lg.plr');

	-- load pbox table with player box scores after inserting player causing issue
	raise notice e'inserting player box stats into stats.pbox...\n';
	call stats.sp_pbox();
	raise notice e'pbox insert complete: %\n', fn_cntstr('stats.pbox');

	-- load api.plr_agg table with pbox stats 
	raise notice e'inserting season/career stat aggregations into api.plr_agg...\n';
	call api.sp_plr_agg();
	raise notice e'player agg insert complete: %\n', fn_cntstr('api.plr_agg');
end; $$;
-- call sp_build();