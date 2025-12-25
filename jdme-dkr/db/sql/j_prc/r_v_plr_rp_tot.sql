-- career reg/pl totals
create or replace view api.v_plr_rp_tot as
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
    'tot' as "stype", 
    b.player, 
    c.team as "team", 
    c.team_long as "team_long",
    count(distinct a.game_id) as "gp", 
    sum(a.mins) as "minutes",
    sum(a.pts) as "points", 
    sum(a.ast) as "assists", 
	sum(a.reb) as "rebounds", 
	sum(a.stl) as "steals", 
	sum(a.blk) as "blocks", 
	sum(a.fgm) as "fgm", 
	sum(a.fga) as "fga",
	coalesce(
		cast(round(avg(a.fgp) * 100, 2) as varchar(10)) || '%', '0%')
	as "fgp",
	sum(a.f3m) as "f3m", 
	sum(a.f3a) as "f3a",
	coalesce(
		cast(round(avg(a.f3p) * 100, 2) as varchar(10)) || '%', '0%')
	as "f3p",
	sum(a.ftm) as "ftm", 
	sum(a.fta) as "fta",
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