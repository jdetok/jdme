package mariadb

// this file should just be string literals of queries to pass to the Select function

type Query struct {
	Args []string // arguments to accept
	Q    string   // query
}

type Queries struct {
	DbQueries []Query
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
	where b.game_date = (select max(game_date) from game)
	group by a.game_id, a.team_id
	order by e.pts desc
	`,
}

var PlayersSeason = Query{
	Q: `
	select a.player_id, a.player, a.league, 
		max(a.season_id) as rs_max, 
		min(a.season_id) as rs_min, 
		b.po_max, b.po_min
	from api_player_stats a
	inner join (
		select player_id, player, league, max(season_id) as po_max, min(season_id) as po_min
		from api_player_stats
		where left(season_id, 1) = 4
		group by player_id, player, league, left(season_id, 1)
	) b on b.player_id = a.player_id
	where left(a.season_id, 1) = 2
	group by a.player_id, a.player, a.league, left(a.season_id, 1)
	`,
}
var PlayersSeasonOld = Query{
	Q: `
	select player_id, player, league, max(season_id), min(season_id)
	from api_player_stats
	where left(season_id, 1) = 2
	group by player_id, player, league, left(season_id, 1);
	`,
}

var Szn = Query{
	Q: `select season_desc, wseason_desc from season where season_id = ?`,
}

var Player = Query{
	Q: `
		select a.*, b.season_desc, b.wseason_desc
		from api_player_stats a
		join season b on b.season_id = a.season_id
		where a.player_id = ? and a.season_id = ?
	`,
}

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

var PlayerRec0712 = Query{
	Q: `
		select * from api_player_stats 
		where player_id = ? and season_id = ?
	`,
}

var RecentGames = Query{
	Q: `
	select a.game_id, a.game_date, a.final, a.ot 
	from game a
	where a.game_date = 
	(select max(z.game_date)
	from game z)
	`,
}

var TopScorer = Query{
	Q: `
	select * 
	from top_scorers 
	order by points desc, (assists + rebounds + steals + blocks) desc 
	limit 1;
	`,
}

var Test = Query{
	Q: `
select 
	a.player, 
	b.team,
	b.team_name,
    a.lg,
    a.active,
    e.season_id,
    case 
        when a.lg = "WNBA"
        then e.wseason_desc
        else e.season_desc
    end as season,
	sum(c.pts) as pts, 
	sum(c.ast) as ast,
	sum(c.reb) as reb,
	sum(d.fgm) as fgm,
	sum(d.fg3m) as fg3m,
	sum(d.ftm) as ftm
	from player a
	inner join team b on b.team_id = a.team_id
	inner join p_box c on c.player_id = a.player_id
	inner join p_shtg d 
		on d.player_id = a.player_id and d.game_id = c.game_id
	inner join season e on e.season_id = c.season_id
 	where a.lg = ?
	and e.season_id = ?
	and b.team = ?
	group by a.player, b.team, a.lg, e.season_id, season
	order by pts desc
`,
}

// var Avgs25 = Query{
// 	"select * from v_nba_rs25_avgs"
// }

var Players = Query{
	Args: []string{},
	Q: `
	select player_id, player, lg 
	from player 
	where lg in ("NBA", "WNBA") 
	group by player_id, player, lg
	`,
}

var PlayersOld = Query{
	Args: []string{},
	Q: `
	select a.player_id
	from player a
	where a.player = ?
	limit 1
	`,
}

var Seasons = Query{
	Args: []string{},
	Q: `
	select season_id, season_desc, wseason_desc
	from season
	where left(season_id, 1) in ('2', '4')
	and right(season_id, 4) >= 2000
	-- and right(season_id, 4) >= year(sysdate()) - 15
	order by right(season_id, 4) desc, left(season_id, 1)
	`,
}
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
	SELECT a.lg, a.team_id, a.team, a.team_name
	FROM team a
	INNER JOIN ( 
		SELECT season_id, team_id
		FROM t_box
		WHERE LEFT(season_id, 1) = '2'
		AND RIGHT(season_id, 4) >= '2000'
		GROUP BY season_id, team_id
		) b ON b.team_id = a.team_id
	WHERE a.lg in ('NBA', 'WNBA')
	GROUP BY a.lg, a.team_id, a.team, a.team_name
	`,
}

// -- and a.lg = ?
var LgPlayerStat = Query{
	Args: []string{"lg", "player"},
	Q: `
	select a.player, b.team, 
		sum(c.pts) as pts, 
		sum(c.ast) as ast,
		sum(c.reb) as reb,
		sum(d.fgm) as fgm,
		sum(d.fg3m) as fg3m,
		sum(d.ftm) as ftm-- ,
		-- avg(d.fg_pct) as fg_pct,
		-- avg(d.fg3_pct) as fg3_pct,
		-- avg(d.ft_pct) as ft_pct
		
	from player a
	inner join team b on b.team_id = a.team_id
	inner join p_box c on c.player_id = a.player_id
	inner join p_shtg d 
		on d.player_id = a.player_id and d.game_id = c.game_id
	inner join season e on e.season_id = c.season_id
	where a.active = 1
	and a.lg = ?
	and a.player_id = ?
	and e.season like "%RS"
	group by a.player, b.team	
	order by pts desc
	`,
}

var LgPlayerAvg = Query{
	Args: []string{"lg", "player"},
	Q: `
	select 
	a.player,
	b.team, 
	round(avg(c.pts), 2) as pts,
	round(avg(c.ast), 2) as ast,
	round(avg(c.reb), 2) as reb,
	round(avg(d.fgm), 2) as fgm,
	round(avg(d.fg3m), 2) as fg3m,
	round(avg(d.ftm), 2) as ftm,
	round(avg(d.fg_pct), 2) as fg_pct,
	round(avg(d.fg3_pct), 2) as fg3_pct,
	round(avg(d.ft_pct), 2) as ft_pct
	
	from player a
	inner join team b on b.team_id = a.team_id
	inner join p_box c on c.player_id = a.player_id
	inner join p_shtg d on d.player_id = a.player_id and d.game_id = c.game_id
	inner join season e on e.season_id = c.season_id
	
	where a.active = 1
	and e.season like "%RS"
	and a.lg = ?
	and a.player_id = ?
	
	group by a.player, b.team	
	order by pts desc;
	`,
}
