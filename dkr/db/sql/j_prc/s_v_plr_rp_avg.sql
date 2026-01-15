-- career reg/pl avgs

create or replace view api.v_plr_rp_avg as
with plr_tm_gp as ( -- games played player for each team
	select player_id, team_id, count(distinct game_id) as gp, sum(mins) as mins
	from stats.pbox
	group by player_id, team_id
),
ranked as ( -- rank gp to get team player played longest for 
	select *, 
	rank() over (partition by player_id order by gp desc, mins desc) as rnk
	from plr_tm_gp
)
select 
	a.player_id as "player_id", 
    c.team_id as "team_id", -- joined to the team with most games played
    d.lg, 
    e.szn_id, 
    e.szn_desc, 
    e.wszn_desc, 
    'avg' as "stype", 
    b.player, 
    c.team as "team", 
    c.team_long as "team_long",
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
inner join ranked aa on aa.player_id = a.player_id and aa.rnk = 1
inner join lg.plr b on b.player_id = a.player_id
inner join lg.team c on c.team_id = aa.team_id
inner join lg.league d on d.lg_id = b.lg_id
inner join lg.szn e -- reg. season and playoff aggregates
	on e.szn_id = cast(left(cast(a.szn_id as varchar(5)), 1) || '9999' as int)
where b.lg_id < 2 
group by a.player_id, c.team_id, d.lg, e.szn_id, b.player, e.szn_desc, 
	e.wszn_desc, c.team, c.team_long;