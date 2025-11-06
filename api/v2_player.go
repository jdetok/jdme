package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jdetok/go-api-jdeko.me/pkg/clnd"
	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
)

// read player from query string and clean the value (remove accents, lowercase)
func (app *App) PlayerFromQ(r *http.Request) any {
	pStr := r.URL.Query().Get("player")
	// check if integer
	plrIdInt, err := strconv.ParseUint(pStr, 10, 64)
	if err == nil {
		return plrIdInt
	}

	p_lwr := strings.ToLower(pStr)
	p_cln := clnd.RemoveDiacritics(p_lwr)

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

func (app *App) TeamFromQ(r *http.Request) (uint64, error) {
	var teamId uint64 = 0
	var err error
	t := r.URL.Query().Get("team")
	if t != "" {
		teamId, err = app.MStore.Maps.GetTeamIDUintCC(t)
		if err != nil {
			return teamId, err
		}
	}
	return teamId, nil
}

// new endpoint for use with new player store data structure
func (app *App) HndlPlayerV2(w http.ResponseWriter, r *http.Request) {
	app.Lg.LogHTTP(r)
	// var exists bool

	// get team from query string
	teamQ, err := app.TeamFromQ(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	// get player from query string
	playerQ := app.PlayerFromQ(r)

	var plrId uint64
	if _, ok := playerQ.(uint64); ok {
		plrId = playerQ.(uint64)
	} else {
		plrId, err = app.MStore.Maps.GetPlrIdFromName(playerQ.(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
	}

	// get season from query string
	seasonQ, err := app.SeasonFromQ(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	} // get most recent season for 0 or 88888

	fmt.Println(seasonQ)
	// use TmSznPlr
	var rp Resp
	iq := PQueryIds{PId: plrId, TId: teamQ, SId: seasonQ}

	fmt.Printf("%d | %d | %d\n", iq.PId, iq.SId, iq.TId)

	js, err := rp.GetPlayerDashV2(app.DB, app.MStore.Maps, &iq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	app.JSONWriter(w, js)
	app.Lg.Infof("served /v2/player request")
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

	if iq.SId == 0 || iq.SId == 88888 {
		maxSzn, err := sm.GetSznFromPlrId(iq.PId)
		if err != nil {
			return err
		}
		iq.SId = maxSzn
	}
	args := []any{iq.PId, iq.TId, iq.SId}

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
