create or replace view v_szn_avgs as
select 
	a.player, 
	b.team,
    a.lg,
    a.active,
    e.season_id,
    case 
        when a.lg = "WNBA"
        then e.wseason_desc
        else e.season_desc
    end as season,
	avg(c.pts) as pts, 
	avg(c.ast) as ast,
	avg(c.reb) as reb,
	avg(d.fgm) as fgm,
	avg(d.fg3m) as fg3m,
	avg(d.ftm) as ftm
	from player a
	inner join team b on b.team_id = a.team_id
	inner join p_box c on c.player_id = a.player_id
	inner join p_shtg d 
		on d.player_id = a.player_id and d.game_id = c.game_id
	inner join season e on e.season_id = c.season_id
 	where a.lg <> "GNBA"
	and LEFT(e.season_id, 1) = 2
	group by a.player, b.team, a.lg, e.season_id, season
	order by pts desc;