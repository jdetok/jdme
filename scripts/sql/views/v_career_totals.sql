create or replace view v_career_totals as
select 
	a.player, 
	b.team,
    a.lg,
	sum(c.pts) as pts, 
	sum(c.ast) as ast,
	sum(c.reb) as reb,
	sum(d.fgm) as fgm,
	sum(d.fg3m) as fg3m,
	sum(d.ftm) as ftm
	from player a
	inner join team b on b.team_id = a.team_id
	inner join p_box c on c.player_id = a.player_id
	inner join p_shtg d 
		on d.player_id = a.player_id and d.game_id = c.game_id
	inner join season e on e.season_id = c.season_id
 	where a.lg <> "GNBA"
	and LEFT(e.season_id, 1) = 2
	group by a.player, b.team, a.lg
	order by pts desc;