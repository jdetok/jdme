-- season totals
create or replace view api.v_plr_szn_tot as
select 
    a.player_id, max(a.team_id) as "team_id", d.lg, 
    a.szn_id, e.szn_desc, e.wszn_desc, 'tot' as "stype", 
    b.player, max(c.team) as "team", max(c.team_long) as "team_long",
    count(distinct a.game_id) as "gp", sum(a.mins) as "minutes",
	sum(a.pts) as "points", sum(a.ast) as "assists", 
	sum(a.reb) as "rebounds", sum(a.stl) as "steals", sum(a.blk) as "blocks", 
	sum(a.fgm) as "fgm", sum(a.fga) as "fga",
	coalesce(
		cast(round(avg(a.fgp) * 100, 2) as varchar(10)) || '%', '0%')
	as "fgp",
	sum(a.f3m) as "f3m", sum(a.f3a) as "f3a",
	coalesce(
		cast(round(avg(a.f3p) * 100, 2) as varchar(10)) || '%', '0%')
	as "f3p",
	sum(a.ftm) as "ftm", sum(a.fta) as "fta",
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