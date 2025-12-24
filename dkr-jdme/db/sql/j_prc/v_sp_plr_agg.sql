/*
THE VIEWS & PROCEDURES DEFINED IN THIS SCRIPT ARE RESPONSIBLE FOR INSERTING
AGGREGATE PLAYER STATS INTO TABLE api.plr_agg
FOR EACH PLAYER, STATS ARE AGGREGATED BY INDIVIDUAL REGULAR/POST SEASON, 
CAREER REGULAR/POST SEASON, AND CAREER COMBINED REG/POST SEASON. 
EACH PLAYER WILL HAVE TWO ROWS PER SEASON/AGG TYPE: ONE WITH AVERAGE (PER GAME)
STATS AND ONE WITH SEASON TOTALS. S
*/
/*
drop view if exists api.v_plr_szn_tot;
drop view if exists api.v_plr_szn_avg;
drop view if exists api.v_plr_rp_tot;
drop view if exists api.v_plr_rp_avg;
drop view if exists api.v_plr_cc_tot;
drop view if exists api.v_plr_cc_avg;
*/

-- ============================================================================
/* 
STORED PROCEDURE TO INSERT THE RESULTS OF THE VIEWS ABOVE INTO API TABLE
HAD TO ADD ON CONFLICT DO NOTHING TO EACH INSERT
EDITED AND MOVED EACH VIEW TO ITS OWN FILE. FOR CAREER VIEWS, JOINED TO TEAM ON
TEAM WITH MOST GAMES/MINUTES PLAYED
*/ 
create or replace procedure api.sp_plr_agg()
language plpgsql
as $$
begin
	raise notice e'deleting existing values in api.plr_agg\n';
    truncate api.plr_agg;

	-- season aggs
	raise notice 'inserting season totals';
    insert into api.plr_agg 
		select * from api.v_plr_szn_tot
	on conflict(player_id, season_id, stat_type) do nothing;
	raise notice e'season totals complete: %\n', public.fn_cntstr('api.plr_agg');

	raise notice 'inserting season avgs';
	insert into api.plr_agg 
		select * from api.v_plr_szn_avg
	on conflict(player_id, season_id, stat_type) do nothing;
	raise notice e'season avgs complete: %\n', public.fn_cntstr('api.plr_agg');

	-- reg season/playoff aggs
	raise notice 'inserting rs/playoff totals';
    insert into api.plr_agg 
		select * from api.v_plr_rp_tot
	on conflict(player_id, season_id, stat_type) do nothing;
	raise notice e'regszn/playoff totals complete: %\n', public.fn_cntstr('api.plr_agg');

	raise notice 'inserting rs/playoff avgs';
	insert into api.plr_agg 
		select * from api.v_plr_rp_avg
	on conflict(player_id, season_id, stat_type) do nothing;
	raise notice e'regszn/playoff avgs complete: %\n', public.fn_cntstr('api.plr_agg');

	-- combined reg season/playoff aggs
	raise notice 'inserting combined rs/playoff totals';
	insert into api.plr_agg
		select * from api.v_plr_cc_tot
	on conflict(player_id, season_id, stat_type) do nothing;
	raise notice e'combined totals complete: %\n', public.fn_cntstr('api.plr_agg');

	raise notice 'inserting combined rs/playoff avgs';
	insert into api.plr_agg 
		select * from api.v_plr_cc_avg
	on conflict(player_id, season_id, stat_type) do nothing;
	raise notice e'combined avgs complete: %\n', public.fn_cntstr('api.plr_agg');
end; $$;