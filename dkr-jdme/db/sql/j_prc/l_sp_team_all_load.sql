create or replace procedure lg.sp_team_all_load()
language plpgsql
as $$
begin
	-- drop current teams except for 0 placeholder
	-- delete from lg.team where team_id > 0;

	-- insert all teams from intake.gm_team & intake.[w]player
	insert into lg.team
		select
		    0,
		    a.team_id,
		    a.team_abbreviation,
		    lower(a.team_abbreviation) || '_' || c.team_code,
		    a.team_name,
		    c.team_city,
		    c.team_name
		from intake.gm_team a
		inner join (
		    select 
		        team_id as t_id,
		        max(season_id) as s_id
		    from intake.gm_team
		    group by team_id
		) b on b.t_id = a.team_id and b.s_id = a.season_id
		inner join intake.player c
		    on c.team_id = a.team_id
		-- player from year greater than current season year
		where c.from_year <= right(cast(a.season_id as varchar(5)), 4)
		and c.to_year >= right(cast(a.season_id as varchar(5)), 4)
		and c.team_id > 0 -- no team_id = 0
		group by a.season_id, a.team_id, a.team_abbreviation, 
		    c.team_code, a.team_name, c.team_city, c.team_name 
		-- ============================================================
		union -- ++++++++++++++++++++++++++++++++++++++++++++++++++++++
		-- WNBA QUERY
		select
		    1,
		    a.team_id,
		    a.team_abbreviation,
		    lower(a.team_abbreviation) || '_' || c.team_code,
		    a.team_name,
		    c.team_city,
		    c.team_name
		from intake.gm_team a
		inner join (
		    select 
		        team_id as t_id,
		        max(season_id) as s_id
		    from intake.gm_team
		    group by team_id
		) b on b.t_id = a.team_id and b.s_id = a.season_id
		inner join intake.wplayer c
		    on c.team_id = a.team_id
		-- player from year greater than current season year
		where c.from_year <= right(cast(a.season_id as varchar(5)), 4)
		and c.to_year >= right(cast(a.season_id as varchar(5)), 4)
		and c.team_id > 0 -- no team_id = 0
		group by a.season_id, a.team_id, a.team_abbreviation, 
		    c.team_code, a.team_name, c.team_city, c.team_name
		order by team_id
	on conflict(team_id) do nothing;
end; $$;