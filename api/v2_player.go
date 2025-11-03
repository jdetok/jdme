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
		fmt.Printf("integer player id requested: %d\n", plrIdInt)
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

// new endpoint for use with new player store data structure
func (app *App) HndlPlayerV2(w http.ResponseWriter, r *http.Request) {
	app.Lg.LogHTTP(r)
	var exists bool
	seasonQ, sznErr := app.SeasonFromQ(r)
	if sznErr != nil {
		http.Error(w, sznErr.Error(), http.StatusUnprocessableEntity)
	}

	// 1610612744
	teamQ := r.URL.Query().Get("team")
	fmt.Println(teamQ)
	playerQ := app.PlayerFromQ(r)
	if plrId, ok := playerQ.(uint64); ok {
		if teamQ != "" {
			tmId := app.MStore.Maps.GetTeamIDUintCC(teamQ)

			exists = app.MStore.Maps.PlrSznTmExists(plrId, tmId, seasonQ)
		} else {
			exists = app.MStore.Maps.PlrIdSznExists(plrId, seasonQ)
		}

	} else {
		exists = app.MStore.Maps.PlrSznExists(playerQ.(string), seasonQ)
	}

	var wErr error

	if exists {
		_, wErr = fmt.Fprintf(w, "player %v team %v exists in season %d\n", playerQ, teamQ, seasonQ)
	} else {
		_, wErr = fmt.Fprintf(w, "player %v does not exist in season %d\n", playerQ, seasonQ)
	}

	if wErr != nil {
		http.Error(w,
			fmt.Sprintf("failed to write HTTP response\n**%s", wErr),
			http.StatusInternalServerError)
	}
	app.Lg.Infof("served /v2/player request")
}
