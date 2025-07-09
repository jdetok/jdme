package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
)

type PlayerMeta struct {
	PlayerId string `json:"playerId"`
	TeamId string `json:"teamId"`
	League string `json:"league"`
	SeasonID string `json:"seasonId"`
	GameId string `json:"gameId"`
	GameDate string `json:"gameDate"`
	Player string  `json:"players"`
	Team string  `json:"team"`
	TeamName string  `json:"teamName"`
	Caption string  `json:"caption"`
	CaptionShort string  `json:"captionShort"`
	HeadshotUrl string `json:"headshotUrl"`
}

// idea: break out box and shooting
type Stats struct {
	Minutes string `json:"minutes"`
	Points string `json:"points"`
	Assists string `json:"assists"`
	Rebounds string `json:"rebounds"`
	Steals string `json:"steals"`
	Blocks string `json:"blocks"`
	FgMade string `json:"fgMade"`
	FgAtpt string `json:"fgAtpt"`
	FgPct string `json:"fgPct"`
	Fg3Made string `json:"fg3Made"`
	Fg3Atpt string `json:"fg3Atpt"`
	Fg3Pct string `json:"fg3Pct"`
	FtMade string `json:"ftMade"`
	FtAtpt string `json:"ftAtpt"`
	FtPct string `json:"ftPct"`
}

type TopScorePlayer struct {
	Meta PlayerMeta `json:"meta"`
	Stats Stats `json:"stats"`
}

// wrap top scrorer in struct to name the json object
type TopScorers struct {
	Players []TopScorePlayer `json:"players"`
	// Players map[string]TopScorePlayer `json:"players"`
}

func (ts *TopScorers) GetTopScorers(db *sql.DB) ([]byte, error) {
	rows, err := db.Query(mariadb.TopScorer.Q)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ts.MakeTopScorers(rows)
	js, err := json.Marshal(ts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return js, nil
}

// scans sql rows to appropriate struct field, runs meta funcs
func (ts *TopScorers) MakeTopScorers(rows *sql.Rows) {//(TopScorers, error) 
	for rows.Next() {
		var tsp TopScorePlayer

		rows.Scan(&tsp.Meta.PlayerId, &tsp.Meta.TeamId, &tsp.Meta.League,
			&tsp.Meta.SeasonID, &tsp.Meta.GameId, &tsp.Meta.GameDate,
			&tsp.Meta.Player, &tsp.Meta.Team, &tsp.Meta.TeamName,
			&tsp.Stats.Minutes, &tsp.Stats.Points, &tsp.Stats.Assists,
			&tsp.Stats.Rebounds, &tsp.Stats.Steals, &tsp.Stats.Blocks,
			&tsp.Stats.FgMade, &tsp.Stats.FgAtpt, &tsp.Stats.FgPct,
			&tsp.Stats.Fg3Made, &tsp.Stats.Fg3Atpt, &tsp.Stats.Fg3Pct,
			&tsp.Stats.FtMade, &tsp.Stats.FtAtpt, &tsp.Stats.FtPct)
		
		tsp.Meta.MakeCaptions()
		tsp.Meta.MakeHeadshotUrl()
		ts.Players = append(ts.Players, tsp)
	}
}

func (pm *PlayerMeta) MakeCaptions() {
	pm.Caption = fmt.Sprintf("%s - %s", pm.Player, pm.TeamName)
	pm.CaptionShort = fmt.Sprintf("%s - %s", pm.Player, pm.Team)
}

func (pm *PlayerMeta) MakeHeadshotUrl() {
	lg := strings.ToLower(pm.League)
	pm.HeadshotUrl = fmt.Sprintf(
		`https://cdn.%s.com/headshots/%s/latest/1040x760/%s.png`, 
		lg, lg, pm.PlayerId)
}


/*
game_id int,
	game_date date,
	season_id int,
	team_id int,
	team varchar(3),
	team_name varchar(255),
	player_id int,	
	player varchar(255),
	minutes int,
	points int,
	assists int,
	rebounds int,
	steals int,
	blocks int,
	fgm int,
	fga int,
	fgp varchar(10),
	fg3m int,
	fg3a int,
	fg3p varchar(10),
	ftm int,
	fta int,
	ftp varchar(10),
*/
/*
type TopScorer struct {
	GameId string `json:"gameId"`
	GameDate string `json:"gameDate"`
	SeasonID string `json:"seasonId"`
	League string `json:"leageue"`
	TeamId string `json:"teamId"`
	Team string `json:"team"`
	TeamName string `json:"teamName"`
	PlayerId string `json:"playerId"`
	Player string `json:"player"`
	Minutes string `json:"minutes"`
	Points string `json:"points"`
	Assists string `json:"assists"`
	Rebounds string `json:"rebounds"`
	Steals string `json:"steals"`
	Blocks string `json:"blocks"`
	FgMade string `json:"fgMade"`
	FgAtpt string `json:"fgAtpt"`
	FgPct string `json:"fgPct"`
	Fg3Made string `json:"fg3Made"`
	Fg3Atpt string `json:"fg3Atpt"`
	Fg3Pct string `json:"fg3Pct"`
	FtMade string `json:"ftMade"`
	FtAtpt string `json:"ftAtpt"`
	FtPct string `json:"ftPct"`
}

func (tsp *TopScorePlayer) MakeTopScorer(rows *sql.Rows) (TopScorers, error) {
	for rows.Next() {
		rows.Scan(&tsp.GameId, &tsp.GameDate, &tsp.SeasonID, &tsp.League, &tsp.TeamId, &tsp.Team,
			&tsp.TeamName, &tsp.PlayerId, &tsp.Player, &tsp.Minutes, &tsp.Pointsp,
			&tsp.Assistsp, &tsp.Rebounds, &tsp.Steals, &tsp.Blocks, &tsp.FgMade,
			&tsp.FgAtpt, &tsp.FgPct, &tsp.Fg3Made, &tsp.Fg3Atpt, &tsp.Fg3Pct,
			&tsp.FtMade, &tsp.FtAtpt, &tsp.FtPct)
	}
}

func (ts *TopScorer) GetTopScorer(db *sql.DB) ([]byte, error){
	rows, err := db.Query(mariadb.TopScorer.Q)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	ts.MakeTopScorer(rows)
	ts.MakeHeadshotUrl()
	js, err := json.Marshal(ts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return js, nil
}


*/