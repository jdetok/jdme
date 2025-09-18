//POSTGRES `bball` QUERIES: MIGRATED TO POSTGRES FROM MARIADB 08/06/2025

package pgdb

type Query struct {
	Args []string // arguments to accept
	Q    string   // query
}

// query all seasons in database, populates global seasons struct
var AllSeasons = Query{
	Args: []string{},
	Q: `
	select szn_id, szn_desc, wszn_desc
	from lg.szn	
	where left(cast(szn_id as varchar(5)), 1) in ('2', '4', '9')
	order by right(cast(szn_id as varchar(5)), 4) desc, 
	left(cast(szn_id as varchar(5)), 1)	
	`,
}

// player dash from api table from passed player and season
var PlayerDash = Query{
	Q: `select * from api.plr_agg where player_id = $1 and season_id = $2`,
}

/*
get the player dash for most recent night's games top scorer
*/
var TeamTopScorerDash = Query{
	Q: `
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
	`,
}

/*
team and top player stats from most recent night's games
*/
var RecGameTopScorers = Query{
	Q: `
	select * from (
		select distinct on (a.game_id)
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
		d.pts as plr_pts
		from stats.tbox a
		inner join (
			select max(gdate) as md
			from stats.tbox
		) b on a.gdate = b.md
		inner join lg.team c on c.team_id = a.team_id
		inner join stats.pbox d on d.game_id = a.game_id and d.team_id = a.team_id
		inner join lg.plr e on e.player_id = d.player_id
		order by a.game_id, d.pts desc, (d.ast + d.reb + d.stl + d.blk) desc)
	order by plr_pts desc
	`,
}

/*
select each player and the min max reg/post season stats. used to populate global
player store
*/
var PlayersSeason = Query{
	Q: `
	select 
		a.player_id,
		lower(a.player) as plr,
		case 
			when a.lg_id = 0 then 'nba'
			when a.lg_id = 1 then 'wnba'
		end,
		b.rs_max, 
		b.rs_min,
		coalesce(c.po_max, b.rs_max),
		coalesce(c.po_min, b.rs_min)
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
	`,
}

/*
query database for all teams, used to populate global teams store
*/
var Teams = Query{
	Args: []string{},
	Q: `
	select
		case
			when lg_id = 0 then 'NBA'
			when lg_id = 1 then 'WNBA'
		end,
		team_id, team, team_long
	from lg.team
	where team_id > 0
	`,
}

// top five players by points, season id and lg_cde (nba or wnba) as arguments
var LgTop5 = Query{
	Args: []string{},
	Q: `
	select 
		a.player_id,
		b.player,
		max(c.team) as team,
		sum(a.pts) as points
	from stats.pbox a
	inner join lg.plr b on b.player_id = a.player_id
	inner join lg.team c on c.team_id = a.team_id
	inner join lg.league d on d.lg_id = b.lg_id
	where a.szn_id = $1
	and d.lg_cde = $2
	group by a.player_id, b.player, b.lg_id
	order by points desc
	limit 5
	`,
}
