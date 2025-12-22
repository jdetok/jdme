package memd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
)

type TeamRecords struct {
	NBARecords  []TeamRecord `json:"nba_team_records"`
	WNBARecords []TeamRecord `json:"wnba_team_records"`
}

type TeamRecord struct {
	League     string `json:"league"`
	SeasonId   int    `json:"season_id"`
	Season     string `json:"season"`
	SeasonDesc string `json:"season_desc"`
	TeamId     uint64 `json:"team_id"`
	Team       string `json:"team"`
	TeamLong   string `json:"team_long"`
	Wins       uint16 `json:"wins"`
	Losses     uint16 `json:"losses"`
}

// get from memory
func TeamRecordsJSON(team_recs *TeamRecords) ([]byte, error) {
	js, err := json.Marshal(team_recs)
	if err != nil {
		msg := "error marshaling team records"
		return nil, fmt.Errorf("%s\n%v", msg, err)
	}
	return js, nil
}

type TopPlayers struct {
	PlayerId uint64 `json:"player_id"`
	Player   string `json:"player"`
	Season   string `json:"season"`
	Team     string `json:"team"`
	Points   uint32 `json:"points"`
}

type LgTopPlayers struct {
	NBATop  []TopPlayers `json:"nba"`
	WNBATop []TopPlayers `json:"wnba"`
}

type RecentGames struct {
	TopScorers []PlayerBasic `json:"top_scorers"`
	Games      []RecentGame  `json:"recent_games"`
}

type PlayerBasic struct {
	PlayerId uint64 `json:"player_id"`
	TeamId   uint64 `json:"team_id"`
	Player   string `json:"player"`
	League   string `json:"league"`
	Points   uint16 `json:"points"`
	Assists  uint16 `json:"assists"`
	Rebounds uint16 `json:"rebounds"`
}

type RecentGame struct {
	GameId    uint64 `json:"game_id"`
	TeamId    uint64 `json:"team_id"`
	PlayerId  uint64 `json:"player_id"`
	Player    string `json:"player"`
	League    string `json:"league"`
	Team      string `json:"team"`
	TeamName  string `json:"team_name"`
	GameDate  string `json:"game_date"`
	Matchup   string `json:"matchup"`
	WinLoss   string `json:"wl"`
	Points    uint16 `json:"points"`
	OppPoints uint16 `json:"opp_points"`
}

// use the LgTop5 query to get top 5 players per league
// set a slice of strings with both leagues to loop through. NBA is first in the slice -
// this must be maintained for the logic to work. at the end of the for loop the sId
// variable is set to the WNBA season - it's declared as the NBA season before the loop begins
func QueryTopLgPlayers(db pgdb.DB, cs *CurrentSeasons, numPl string) (LgTopPlayers, error) {

	var lt LgTopPlayers

	// current seasons by league
	sl, err := cs.LgSznsByMonth(time.Now())
	if err != nil {
		return lt, err
	}

	// query appropriate season for each league
	var lgs = [2]string{"nba", "wnba"}
	var sId string = strconv.Itoa(sl.SznId)
	for _, lg := range lgs {
		// query database
		r, err := db.Query(pgdb.LeagueTopScorers, sId, lg, numPl)
		if err != nil {
			msg := fmt.Sprintf(
				"failed to query database for top 5 lg players: sznId: %s | lg: %s\n",
				sId, lg)
			return lt, fmt.Errorf("%s\n%v", msg, err)
		}

		// create a Top5 struct for each row, append to appropriate NBA/WNBA member
		for r.Next() {
			var t TopPlayers
			r.Scan(&t.PlayerId, &t.Player, &t.Season, &t.Team, &t.Points)
			switch lg {
			case "nba":
				lt.NBATop = append(lt.NBATop, t)
			case "wnba":
				lt.WNBATop = append(lt.WNBATop, t)
			}
		}

		// after first run set wnba season
		sId = strconv.Itoa(sl.WSznId)
	}
	return lt, nil
}

// marshal LgTop5 struct into JSON []byte
func MarshalTopPlayers(lt *LgTopPlayers) ([]byte, error) {
	js, err := json.Marshal(lt)
	if err != nil {
		return nil, err
	}
	return js, nil
}

// returns json of the top scorer (regardless of team) stats from each of most
// recent night's games. used on page load and to populate recent top scorers table
func (rgs *RecentGames) GetRecentGames(db pgdb.DB) ([]byte, error) {
	rows, err := db.Query(pgdb.RecGameTopScorers)
	if err != nil {
		msg := "recent games query failed"
		return nil, fmt.Errorf("%s\n%v", msg, err)
	}

	rgs.ScanRecentGamesRows(rows)
	js, err := json.Marshal(rgs)
	if err != nil {
		msg := "json marshal failed"
		return nil, fmt.Errorf("%s\n%v", msg, err)
	}
	return js, nil
}

// accepts a sql.Rows pointer and scans it to a RecentGames struct
func (rgs *RecentGames) ScanRecentGamesRows(rows *sql.Rows) {
	for rows.Next() {
		var rg RecentGame
		var ps PlayerBasic
		rows.Scan(&rg.GameId, &rg.TeamId, &rg.PlayerId,
			&rg.Player, &rg.League, &rg.Team,
			&rg.TeamName, &rg.GameDate, &rg.Matchup,
			&rg.WinLoss, &rg.Points, &rg.OppPoints,
			&ps.Points, &ps.Assists, &ps.Rebounds)

		ps.PlayerId = rg.PlayerId
		ps.TeamId = rg.TeamId
		ps.Player = rg.Player
		ps.League = rg.League
		rgs.TopScorers = append(rgs.TopScorers, ps)
		rgs.Games = append(rgs.Games, rg)
	}
}
