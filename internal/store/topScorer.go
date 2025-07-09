package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jdetok/go-api-jdeko.me/internal/mariadb"
)

// structs defined in commmon.go
type TopScorePlayer struct {
	MetaG GameMeta `json:"gameMeta"`
	Meta PlayerMeta `json:"playerMeta"`
	Stats Stats `json:"stats"`
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