-- season avgs
create or replace view api.v_plr_szn_avg as
select 
    a.player_id, max(a.team_id) as "team_id", d.lg, 
    a.szn_id, e.szn_desc, e.wszn_desc, 'avg' as "stype", 
    b.player, max(c.team) as "team", max(c.team_long) as "team_long",
    count(distinct a.game_id) as "gp", 
    round(avg(a.mins), 2) as "minutes",
    round(avg(a.pts), 2) as "points", round(avg(a.ast), 2) as "assists", 
	round(avg(a.reb), 2) as "rebounds", round(avg(a.stl), 2) as "steals", 
    round(avg(a.blk), 2) as "blocks", 
	round(avg(a.fgm), 2) as "fgm", round(avg(a.fga), 2) as "fga",
	coalesce(
		cast(round(avg(a.fgp) * 100, 2) as varchar(10)) || '%', '0%')
	as "fgp",
	round(avg(a.f3m), 2) as "f3m", round(avg(a.f3a), 2) as "f3a",
	coalesce(
		cast(round(avg(a.f3p) * 100, 2) as varchar(10)) || '%', '0%')
	as "f3p",
	round(avg(a.ftm), 2) as "ftm", round(avg(a.fta), 2) as "fta",
	coalesce(
		cast(round(avg(a.ftp) * 100, 2) as varchar(10)) || '%', '0%')
	as "ftp"
from stats.pbox a
inner join lg.plr b on b.player_id = a.player_id
inner join lg.team c on c.team_id = a.team_id
inner join lg.league d on d.lg_id = b.lg_id
inner join lg.szn e on e.szn_id = a.szn_id
where b.lg_id < 2 
group by a.player_id, d.lg, a.szn_id, b.player, e.szn_desc, e.wszn_desc
order by a.szn_id desc;