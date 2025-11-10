package api

import (
	"fmt"
	"net/http"

	"github.com/jdetok/go-api-jdeko.me/pkg/clnd"
	"github.com/jdetok/go-api-jdeko.me/pkg/resp"
)

// /player handler
/*
- create PlayerQuery pq struct to hold query string params as strings
- create PQueryIds iq to hold them as uint64s
- get all as strings, remove accents from player names
- call ValidatePlayerSzn
	- pass pq strings, get iq uints back
	- attempts to convert teamId and seasonId to ints
	- calls RandomPlayerId if player name is "random"
	- if passed as a player id it assigns that to iq.PId and moves on
	- otherwise searches passed player name again the in-mem player slice
	- finally calls HandleSeasonId to validate the seasonId
- call GetPlayerDash
	- pass iq (result of ValidatePlayerSzn)
*/
func (app *App) HndlPlayer(w http.ResponseWriter, r *http.Request) {
	app.Lg.LogHTTP(r)

	// pq holds all query parameters as strings
	var pq resp.PlayerQuery
	var iq resp.PQueryIds
	// var iq PQueryIds
	var rp resp.RespPlayerDash
	// var tId uint64

	// get and convert team from query string
	pq.Team = r.URL.Query().Get("team")

	// get season from query string
	pq.Season = r.URL.Query().Get("season")

	//get league from query string
	pq.League = r.URL.Query().Get("league")

	// get player from query string
	pq.Player = clnd.ConvToASCII(r.URL.Query().Get("player"))

	// validate player & get playerid/season id
	iq, err := resp.ValidatePlayerSzn(app.Store.Players, &app.Store.CurrentSzns, &pq, &rp.ErrorMsg)
	if err != nil {
		msg := fmt.Sprintf("validate player %s", pq.Player)
		app.Lg.HTTPErr(w, err, http.StatusUnprocessableEntity, msg)
	}

	// query the player & build JSON response, returned as []byte to write
	js, err := rp.GetPlayerDash(app.DB, &iq)
	if err != nil {
		msg := fmt.Sprintf("server failed to return player dash for %s", pq.Player)
		app.Lg.HTTPErr(w, err, http.StatusInternalServerError, msg)
	}
	// write JSON response
	app.JSONWriter(w, js)
}
