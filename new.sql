select 
    a.player_id, a.team_id, d.lg, 
    29999 as "szn_id_agg", ('Career with ' || c.team_long) as "szn_long", ('Career with ' || c.team_long) as "wszn_long", 'tot' as "stype", 
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
where b.lg_id < 2
and a.player_id = $1
and a.team_id = $2
and a.szn_id::text like $3 || '%'
group by a.player_id, a.team_id,  d.lg, szn_id_agg, b.player, szn_long, wszn_long;
union
select 
    a.player_id, a.team_id, d.lg, a.szn_id, 
     e.szn_desc, e.wszn_desc, 'avg' as "stype", 
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
and a.player_id = $1
and a.team_id = $2
and left(cast(a.szn_id as varchar(5)), 1) = $3
group by a.player_id, a.team_id, d.lg, a.szn_id, b.player, e.szn_desc, e.wszn_desc;


with szntype as (
	select 
		case 
			when $3 = '2' then 'Regular Season'
			when $3 = '4' then 'Post Seasons'
			else ''
		end
)
select 
    a.player_id,
    a.team_id,
    d.lg,
    ($3 || '9999')::int as szn_id_agg,
    ('Career ' || (select * from szntype) || ' with ' || max(c.team)) as szn_long,
    ('Career ' || (select * from szntype) || ' with ' || max(c.team)) as wszn_long,
    'tot' as stype,
    b.player,
    max(c.team) as team,
    max(c.team_long) as team_long,
    count(distinct a.game_id) as gp,
    sum(a.mins) as minutes,
    sum(a.pts) as points,
    sum(a.ast) as assists,
    sum(a.reb) as rebounds,
    sum(a.stl) as steals,
    sum(a.blk) as blocks,
    sum(a.fgm) as fgm,
    sum(a.fga) as fga,
    coalesce(cast(round(avg(a.fgp)*100,2) as varchar)||'%', '0%') as fgp,
    sum(a.f3m) as f3m,
    sum(a.f3a) as f3a,
    coalesce(cast(round(avg(a.f3p)*100,2) as varchar)||'%', '0%') as f3p,
    sum(a.ftm) as ftm,
    sum(a.fta) as fta,
    coalesce(cast(round(avg(a.ftp)*100,2) as varchar)||'%', '0%') as ftp
from stats.pbox a
inner join lg.plr b on b.player_id = a.player_id
inner join lg.team c on c.team_id = a.team_id
inner join lg.league d on d.lg_id = b.lg_id
where b.lg_id < 2
  and a.player_id = $1
  and a.team_id = $2
  and a.szn_id::text like $3 || '%'
group by a.player_id, a.team_id, d.lg, szn_id_agg, b.player
union
select 
    a.player_id,
    a.team_id,
    d.lg,
    ($3 || '9999')::int as szn_id_agg,
    ('Career ' || (select * from szntype) || ' with ' || max(c.team)) as szn_long,
    ('Career ' || (select * from szntype) || ' with ' || max(c.team)) as wszn_long,
    'avg' as stype,
    b.player,
    max(c.team) as team,
    max(c.team_long) as team_long,
    count(distinct a.game_id) as gp,
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
where b.lg_id < 2
  and a.player_id = $1
  and a.team_id = $2
  and a.szn_id::text like $3 || '%'
group by a.player_id, a.team_id, d.lg, szn_id_agg, b.player
;
