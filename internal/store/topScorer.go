package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
)

// structs defined in commmon.go
type TopScorePlayer struct {
	MetaG GameMeta `json:"game_meta"`
	Meta PlayerMeta `json:"player_meta"`
	Box BoxStats `json:"box_stats"`
	Shooting ShootingStats `json:"shooting_stats"`
	// Stats Stats `json:"stats"`
}

// wrap top scrorer in struct to name the json object
type TopScorers struct {
	Players []TopScorePlayer `json:"players"`
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
			&tsp.MetaG.SeasonID, &tsp.MetaG.GameId, &tsp.MetaG.GameDate,
			&tsp.Meta.Player, &tsp.Meta.Team, &tsp.Meta.TeamName,
			&tsp.Box.Minutes, &tsp.Box.Points, &tsp.Box.Assists,
			&tsp.Box.Rebounds, &tsp.Box.Steals, &tsp.Box.Blocks,
			&tsp.Shooting.FgMade, &tsp.Shooting.FgAtpt, &tsp.Shooting.FgPct,
			&tsp.Shooting.Fg3Made, &tsp.Shooting.Fg3Atpt, &tsp.Shooting.Fg3Pct,
			&tsp.Shooting.FtMade, &tsp.Shooting.FtAtpt, &tsp.Shooting.FtPct)
		
		tsp.Meta.MakeCaptions()
		tsp.Meta.MakeHeadshotUrl()
		ts.Players = append(ts.Players, tsp)
	}
}