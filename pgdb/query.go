package pgdb

type Query struct {
	Args []string // arguments to accept
	Q    string   // query
}

// GetSeasons
var AllSeasons = Query{
	Args: []string{},
	Q: `
	select szn_id, szn_desc, wszn_desc
	from lg.szn	
	where left(cast(szn_id as varchar(5)), 1) in ('2', '4')
	and right(cast(szn_id as varchar(5)), 4) != '9999'
	order by right(cast(szn_id as varchar(5)), 4) desc, 
	left(cast(szn_id as varchar(5)), 1)
	`,
}

// GetPlayerDash
var PlayerDash = Query{
	Q: `select * from api.plr_agg where player_id = $1 and season_id = $2`,
}

// GetPlayerDash
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
		a.gdate,
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
