package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jdetok/go-api-jdeko.me/pkg/clnd"
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

	// get season from query string
	seasonQ, sznErr := app.SeasonFromQ(r)
	if sznErr != nil {
		http.Error(w, sznErr.Error(), http.StatusUnprocessableEntity)
	}

	// get team from query string
	teamQ, err := app.TeamFromQ(r)
	if err != nil {
		http.Error(w, sznErr.Error(), http.StatusUnprocessableEntity)
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

	var rp Resp
	iq := PQueryIds{PId: plrId, TId: teamQ, SId: seasonQ}

	fmt.Printf("%d | %d | %d\n", iq.PId, iq.SId, iq.TId)

	js, err := rp.GetPlayerDash(app.DB, &iq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	app.JSONWriter(w, js)
	app.Lg.Infof("served /v2/player request")
}
