//POSTGRES `bball` QUERIES: MIGRATED TO POSTGRES FROM MARIADB 08/06/2025

package pgdb

var QPlayerStore = `
select 
	a.player_id,
	a.player as player,
	lower(a.player) as plr,
	case 
		when a.lg_id = 0 then 'nba'
		when a.lg_id = 1 then 'wnba'
	end as lg,
	c.rs_max, 
	c.rs_min,
	coalesce(d.po_max, 0) as po_max,
	coalesce(d.po_min, 0) as po_min,
	b.teams
from lg.plr a
inner join (
	select 
		player_id, 
		string_agg(distinct team_id::text, ',') as teams
	from stats.pbox
	group by player_id
) b on b.player_id = a.player_id
inner join (
	select player_id, min(season_id) as rs_min, max(season_id) as rs_max
	from api.plr_agg
	where left(cast(season_id as varchar(5)), 1) = '2'
	and right(cast(season_id as varchar(5)), 4) != '9999'
	group by player_id
) c on c.player_id = a.player_id
left join (
	select player_id, min(season_id) as po_min, max(season_id) as po_max
	from api.plr_agg
	where left(cast(season_id as varchar(5)), 1) = '4'
	and right(cast(season_id as varchar(5)), 4) != '9999'
	group by player_id
) d on d.player_id = a.player_id
`

var VerifyTeamSzn = `
select 1
from api.plr_agg
where season_id = $1
and team_id = $2
and player_id = $3
and stat_type = 'tot'
`
var PlayerTeamBool = `
select 1
from api.plr_agg
where player_id = $1
and team_id = $2
`

// query all seasons in database, populates global seasons struct
var AllSeasons = `
select szn_id, szn_desc, wszn_desc
from lg.szn	
where left(cast(szn_id as varchar(5)), 1) in ('2', '4', '9')
order by right(cast(szn_id as varchar(5)), 4) desc, 
left(cast(szn_id as varchar(5)), 1)	
`

/*
select each player and the min max reg/post season stats. used to populate global
player store
9/23/2025: updated playoff min/max coalesce calls to replace null with 0 rather
than reg season min/max. this fixes issue where pressing random player with
playoff season checkbox checked would get players without a playoff game
*/
var PlayersSeason = `
select 
	a.player_id,
	lower(a.player) as plr,
	case 
		when a.lg_id = 0 then 'nba'
		when a.lg_id = 1 then 'wnba'
	end,
	b.rs_max, 
	b.rs_min,
	coalesce(c.po_max, 0),
	coalesce(c.po_min, 0)
from lg.plr a
inner join (
	select player_id, min(season_id) as rs_min, max(season_id) as rs_max
	from api.plr_agg
	where left(cast(season_id as varchar(5)), 1) = '2'
	and right(cast(season_id as varchar(5)), 4) != '9999'
	group by player_id
) b on b.player_id = a.player_id
left join (
	select player_id, min(season_id) as po_min, max(season_id) as po_max
	from api.plr_agg
	where left(cast(season_id as varchar(5)), 1) = '4'
	and right(cast(season_id as varchar(5)), 4) != '9999'
	group by player_id
) c on c.player_id = a.player_id
`

/*
query database for all teams, used to populate global teams store
*/
var Teams = `
select
	case
		when lg_id = 0 then 'NBA'
		when lg_id = 1 then 'WNBA'
	end,
	team_id, team, team_long
from lg.team
where team_id > 0
`

// player dash from api table from passed player and season
var PlayerDash = `select * from api.plr_agg where player_id = $1 and season_id = $2`

/*
get the player dash for most recent night's games top scorer
*/
var TeamTopScorerDash = `
with tstot as ( 
select * 
from api.plr_agg 
where team_id = $1 and season_id = $2
order by points desc
limit 1)
select * from tstot 
union
select a.* 
from api.plr_agg a
inner join tstot b 
	on a.team_id = b.team_id 
	and a.season_id = b.season_id
	and a.player_id = b.player_id
where a.stat_type = 'avg'
`

/*
team and top player stats from most recent night's games
*/
var RecGameTopScorers = `
select * from (
	select distinct on (a.game_id, a.team_id)
	a.game_id,  
	a.team_id, 
	d.player_id, 
	e.player, 
	case 
		when c.lg_id = 0 then 'NBA'
		when c.lg_id = 1 then 'WNBA'
	end as lg, 
	c.team,
	c.team_long,
	to_char(a.gdate, 'MM/DD/YYYY'),
	a.matchup, 
	a.wl, 
	a.pts as tm_pts, 
	f.pts as opp_pts,
	d.pts as plr_pts,
	d.ast as plr_ast,
	d.reb as plr_reb
	from stats.tbox a
	inner join (
		select max(gdate) as md
		from stats.tbox
	) b on a.gdate = b.md
	inner join lg.team c on c.team_id = a.team_id
	inner join stats.pbox d on d.game_id = a.game_id and d.team_id = a.team_id
	inner join lg.plr e on e.player_id = d.player_id
	inner join (
		select game_id, team_id, pts
		from stats.tbox
	) f on f.game_id = a.game_id and f.team_id <> a.team_id
	order by a.game_id, a.team_id, d.pts desc, (d.ast + d.reb + d.stl + d.blk) desc)
order by plr_pts desc
`

// top $3 players by points, season id, lg_cde (nba or wnba), limit # as arguments
var LeagueTopScorers = `
select 
	a.player_id,
	b.player,
	e.szn,
	max(c.team) as team,
	sum(a.pts) as points
from stats.pbox a
inner join lg.plr b on b.player_id = a.player_id
inner join lg.team c on c.team_id = a.team_id
inner join lg.league d on d.lg_id = b.lg_id
inner join lg.szn e on e.szn_id = a.szn_id
where a.szn_id = $1
and d.lg_cde = $2
group by a.player_id, b.player, b.lg_id, e.szn
order by points desc
limit $3
`

var TstTeamPlayer = `
select 
	player_id, team_id, lg, 
	max(szn_id) as szn_id, max(szn_desc) as szn_desc, max(wszn_desc) as wszn_desc,
	max(stype) as stype, max(player) as player, max(team) as team, max(team_long) as team_long, 
	sum(gp) as gp, sum(minutes) as minutes, sum(points) as points, sum(assists) as assists, 
	sum(rebounds) as rebounds, sum(steals) as steals, sum(blocks) as blocks, 
	sum(fgm) as fgm, sum(fga) as fga, 
	round(avg(to_number(substring(fgp, 0, 5), '999.99')), 2) || '%' as fgp, 
	sum(f3m) as f3m, sum(f3a) as f3a, 
	round(avg(to_number(substring(f3p, 0, 5), '999.99')), 2) || '%' as f3p,
	sum(ftm) as ftm, sum(fta) as fta, 
	round(avg(to_number(substring(ftp, 0, 5), '999.99')), 2) || '%' as ftp
from api.v_plr_szn_tot
where player_id = $1 
and team_id = $2
group by player_id, team_id, lg
union
select 
	player_id, team_id, lg, 
	max(szn_id) as szn_id, max(szn_desc) as szn_desc, max(wszn_desc) as wszn_desc,
	max(stype) as stype, max(player) as player, max(team) as team, max(team_long) as team_long, 
	sum(gp) as gp, round(avg(minutes), 2) as minutes, round(avg(points), 2) as points, round(avg(assists), 2) as assists, 
	round(avg(rebounds), 2) as rebounds, round(avg(steals), 2) as steals, round(avg(blocks), 2) as blocks, 
	round(avg(fgm), 2) as fgm, round(avg(fga), 2) as fga, 
	round(avg(to_number(substring(fgp, 0, 5), '999.99')), 2) || '%' as fgp, 
	round(avg(f3m), 2) as f3m, round(avg(f3a), 2) as f3a, 
	round(avg(to_number(substring(f3p, 0, 5), '999.99')), 2) || '%' as f3p,
	round(avg(ftm), 2) as ftm, round(avg(fta), 2) as fta, 
	round(avg(to_number(substring(ftp, 0, 5), '999.99')), 2) || '%' as ftp
from api.v_plr_szn_avg
where player_id = $1
and team_id = $2
group by player_id, team_id, lg
`

var PlTmSzn = `
select 
    (select player from lg.plr where player_id = $1) as player,
    (select team from lg.team where team_id = $2) as team,
    (select szn from lg.szn where szn_id = $3) as season
`
var TeamSznRecords = `
with team_results as (
    select
        d.lg, 
        t.szn_id,
        case
            when d.lg = 'NBA' then szn
            when d.lg = 'WNBA' then wszn
        end as season,
        case
            when d.lg = 'NBA' then szn_desc
            when d.lg = 'WNBA' then wszn_desc
        end as season_desc,
        t.team_id,
        b.team,
        b.team_long,
        count(distinct case when t.wl = 'W' then t.game_id end) as wins,
        count(distinct case when t.wl = 'L' then t.game_id end) as losses
    from stats.tbox t
    inner join lg.team b on b.team_id = t.team_id
    inner join lg.szn c on c.szn_id = t.szn_id
    inner join lg.league d on d.lg_id = b.lg_id
    where t.szn_id in ($1, $2)
    group by d.lg, t.szn_id, 
             case when d.lg = 'NBA' then szn when d.lg = 'WNBA' then wszn end,
             case when d.lg = 'NBA' then szn_desc when d.lg = 'WNBA' then wszn_desc end,
             t.team_id, b.team, b.team_long
)
select lg, szn_id, season, season_desc, team_id, team, team_long, wins, losses
from (
    select *,
        row_number() over (
            partition by lg, szn_id
            order by wins desc, losses asc, team_id
        ) as rank
    from team_results
) ranked
order by lg, rank
`
