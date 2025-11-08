package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/clnd"
	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
)

// /v2/players handler
// REQUIRED ARGS: season, player
// season must be a five digit integer e.g. 22025
// 88888 will retrieve player's most recent season
// player and team can be passed as either strings or integer
// player as string should be player name e.g. "lebron james"
// player="random" can be passed to get a random player filtered by other params
// player as int should be player id e.g. 2544
// team as string should be team abbreviation e.g. "lal"
// team as int should be team id i.e. 1610617247
// OPTIONAL ARGS: league, team
func (app *App) HndlPlayerV2(w http.ResponseWriter, r *http.Request) {
	app.Lg.LogHTTP(r)
	var err error
	var rp Resp

	var lgQ int
	lgQ, err = app.LgFromQ(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	// get season from query string
	seasonQ, err := app.SeasonFromQ(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	} // get most recent season for 0 or 88888

	// get team from query string
	teamQ, err := app.TeamFromQ(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	var tmId uint64
	if tmIdCheck, ok := teamQ.(uint64); ok {
		tmId = tmIdCheck
	} else { // teamQ is a string, get teamId from the passed abbr
		tmIdFromAbbr, err := app.MStore.Maps.GetLgTmIdFromAbbr(teamQ.(string), lgQ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		tmId = tmIdFromAbbr
	}

	// get player from query string
	var plrId uint64
	playerQ := app.PlayerFromQ(r)
	if plrIdCheck, ok := playerQ.(uint64); ok {
		plrId = plrIdCheck
	} else { // playerQ is a string - search for a corresponding playerId
		plrIdUint, err := app.MStore.Maps.GetPlrIdFromName(playerQ.(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		plrId = plrIdUint
	}

	if seasonQ == 0 || seasonQ == 88888 {
		sl := app.Store.CurrentSzns.LgSznsByMonth(time.Now())
		switch lgQ {
		case 0:
			seasonQ = sl.SznId
		case 1:
			seasonQ = sl.WSznId
		default:
			seasonQ = sl.SznId
		}
	}

	// handle random player by league
	if plrId == 77777 {
		rPlrId := app.MStore.Maps.RandomPlrIdV2(tmId, seasonQ, lgQ)
		plrId = rPlrId
	}
	fmt.Printf("validating %d | %d | %d\n", plrId, tmId, seasonQ)
	stp, err := app.MStore.Maps.ValiSznTmPlr(plrId, tmId, seasonQ)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
	fmt.Printf("validated %d | %d | %d\n", stp.PId, stp.TId, stp.SId)
	iq := PQueryIds{PId: stp.PId, TId: stp.TId, SId: stp.SId}

	fmt.Printf("%d | %d | %d\n", iq.PId, iq.SId, iq.TId)

	js, err := rp.GetPlayerDashV2(app.DB, app.MStore.Maps, &iq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	app.JSONWriter(w, js)
	app.Lg.Infof("served /v2/player request")
}

// read player from query string and clean the value (remove accents, lowercase)
func (app *App) PlayerFromQ(r *http.Request) any {
	pStr := r.URL.Query().Get("player")

	// check if integer
	plrIdInt, err := strconv.ParseUint(pStr, 10, 64)
	if err == nil {
		return plrIdInt
	}

	// clean string & remove accents on letters (all standard ascii)
	p_lwr := strings.ToLower(pStr)
	p_cln := clnd.ConvToASCII(p_lwr)

	fmt.Printf("player request (raw): %s | cleaned: %s\n", pStr, p_cln)
	return p_cln
}

// accept http request, get the "season" passed in the query string, return as int
func (app *App) SeasonFromQ(r *http.Request) (int, error) {
	s := r.URL.Query().Get("season")
	s_int, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("INVALID SEASON: could not convert %s to an int\n%s",
			s, err.Error())
	}

	return s_int, nil
}

// returns team arg from query string
// as uint64 if passed a teamId or as string if passed abbr
// returns a 0 if error occurs or team is not included in query string
// if the returned value is a string, the caller must use the league,
// either from the league argument or from the player's league, to get the
// team id as a uint64
func (app *App) TeamFromQ(r *http.Request) (any, error) {
	t := r.URL.Query().Get("team")
	if t != "" {
		// handle request for team abbr
		if _, err := strconv.Atoi(t); err != nil {
			if tmId, ok := app.MStore.Maps.TmAbbrId[t]; !ok {
				return t, fmt.Errorf("couldn't process request for team %v", t)
			} else {
				return tmId, nil // return string team abbr
			}
		}
		// handle request for team id
		teamId, err := app.MStore.Maps.GetTeamIDUintCC(t)
		if err != nil {
			return uint64(0), err
		}
		return teamId, nil
	}
	return uint64(0), nil
}

func (app *App) LgFromQ(r *http.Request) (int, error) {
	lg := r.URL.Query().Get("league")
	lgId, err := strconv.Atoi(lg)
	if err != nil {
		switch lg {
		case "all", "":
			return 10, nil
		case "nba":
			return 0, nil
		case "wnba":
			return 1, nil
		default:
			return 99999, err
		}
	}
	return lgId, nil
}

// after verifying player exists, query db for their stats
func (r *Resp) GetPlayerDashV2(db *sql.DB, sm *memd.StMaps, iq *PQueryIds) ([]byte, error) {

	// query player, scan to structs, call struct functions
	// appends RespObj to r.Results
	if err := r.BuildPlayerRespV2(db, sm, iq); err != nil {
		msg := fmt.Sprintf("failed to query playerId %d seasonId %d", iq.PId, iq.SId)
		return nil, fmt.Errorf("%s\n%v", msg, err)
	}

	// marshall Resp struct to JSON, return as []byte
	js, err := json.Marshal(r)
	if err != nil {
		msg := "failed to marshal structs to json"
		return nil, fmt.Errorf("%s\n%v", msg, err)
	}
	return js, nil
}

func (r *Resp) BuildPlayerRespV2(db *sql.DB, sm *memd.StMaps, iq *PQueryIds) error {
	pOrT := "plr"
	q := pgdb.TmSznPlr
	args := []any{iq.PId, iq.TId, iq.SId}

	if iq.SId == 0 || iq.SId == 88888 {
		maxSzn, err := sm.GetSznFromPlrId(iq.PId)
		if err != nil {
			return err
		}
		iq.SId = maxSzn
	}

	if iq.TId == 0 {
		q = pgdb.PlayerDash
		args = []any{iq.PId, iq.SId}
		fmt.Println("team 0 args:", args)
	}

	if err := r.ProcessRows(db, pOrT, q, args...); err != nil {
		return err
	}
	return nil
}
