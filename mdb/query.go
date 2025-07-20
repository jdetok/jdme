package mdb

type Query struct {
	Args []string // arguments to accept
	Q    string   // query
}

// GetPlayerDash
var Player = Query{
	Q: `
		select a.*, b.season_desc, b.wseason_desc
		from api_player_stats a
		join season b on b.season_id = a.season_id
		where a.player_id = ? and a.season_id = ?
	`,
}

// GetPlayerDash
var TeamSeasonTopP = Query{
	Q: `
		select a.*, b.season_desc, b.wseason_desc
		from api_player_stats a
		join season b on b.season_id = a.season_id
		where team_id = ? and a.season_id = ?
		order by a.points desc
		limit 2;
	`,
}

var RecentGamePlayers = Query{
	Q: `
	select a.game_id, a.team_id, e.player_id, f.player, b.lg, c.team, c.team_name,
	b.game_date, b.matchup, b.final, b.ot, a.pts, e.pts
	from t_box a
	inner join game b on b.game_id = a.game_id
	inner join team c on c.team_id = a.team_id
	inner join p_box d on d.game_id = a.game_id
		and d.team_id = a.team_id
	inner join (
		select game_id, player_id, team_id, pts
		from p_box
		group by game_id, team_id, player_id
		order by pts desc
		limit 1
	) e on e.game_id = a.game_id and e.team_id = a.team_id
	inner join player f on f.player_id = e.player_id and f.team_id = a.team_id
	where b.game_date = (
		select max(game_date) from game 
		where left(season_id, 1) in ('2', '4')
		and lg in ('NBA', 'WNBA')
	)
	and left(a.season_id, 1) in ('2', '4')
	and b.lg in ('NBA', 'WNBA')
	group by a.game_id, a.team_id
	order by e.pts desc
	`,
}

var PlayersSeason = Query{
	Q: `
	select a.player_id, lower(a.player), lower(a.league), 
		max(a.season_id) as rs_max, 
		min(a.season_id) as rs_min, 
		ifnull(b.po_max, 40001) as po_max, ifnull(b.po_min, 40001) as po_min
	from api_player_stats a
	left join (
		select player_id, player, league, max(season_id) as po_max, min(season_id) as po_min
		from api_player_stats
		where left(season_id, 1) = 4
		group by player_id, player, league, left(season_id, 1)
	) b on b.player_id = a.player_id
	where left(a.season_id, 1) = 2
	group by a.player_id, a.player, a.league, left(a.season_id, 1);
	`,
}

// GetSeasons
var RSeasons = Query{
	Args: []string{},
	Q: `
	select season_id, season_desc, wseason_desc
	from season
	where (
		left(season_id, 1) in ('2', '4')
		and right(season_id, 4) >= 2000
	) or season_id > 99990 -- agg seasons 
	order by right(season_id, 4) desc, left(season_id, 1)
	`,
}

var Teams = Query{
	Args: []string{},
	Q: `
	select a.lg, a.team_id, a.team, a.team_name
	from team a
	inner join ( 
		select season_id, team_id
		from t_box
		where left(season_id, 1) = '2'
		and right(season_id, 4) >= '2000'
		group by season_id, team_id
		) b on b.team_id = a.team_id
	where a.lg in ('NBA', 'WNBA')
	group by a.lg, a.team_id, a.team, a.team_name
	`,
}
