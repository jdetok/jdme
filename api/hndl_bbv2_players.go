package api

import (
	"net/http"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/resp"
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
	// app.Lg.HTTPf(r)
	app.Lg.HTTPf(r)
	var err error

	rp := *resp.NewRespPlayerDash(r)

	var lgQ int
	lgQ, err = resp.LgFromQ(r)
	if err != nil {
		app.ErrHTTP(w, err, &rp.Meta, http.StatusUnprocessableEntity)
		return
	}

	// get season from query string
	sl, err := app.Store.CurrentSzns.LgSznsByMonth(time.Now())
	if err != nil {
		app.ErrHTTP(w, err, &rp.Meta, http.StatusUnprocessableEntity)
	}
	seasonQ, err := resp.SeasonFromQ(r, sl.WSznId, sl.WPSznId)
	if err != nil {
		app.ErrHTTP(w, err, &rp.Meta, http.StatusUnprocessableEntity)
		return
	} // get most recent season for 0 or 88888

	// get team from query string
	teamQ, err := resp.TeamFromQ(r, app.MStore.Maps)
	if err != nil {
		app.ErrHTTP(w, err, &rp.Meta, http.StatusUnprocessableEntity)
		return
	}

	var tmId uint64
	if tmIdCheck, ok := teamQ.(uint64); ok {
		tmId = tmIdCheck
	} else { // teamQ is a string, get teamId from the passed abbr
		tmIdFromAbbr, err := app.MStore.Maps.GetLgTmIdFromAbbr(teamQ.(string), lgQ)
		if err != nil {
			app.ErrHTTP(w, err, &rp.Meta, http.StatusUnprocessableEntity)
			return
		}
		tmId = tmIdFromAbbr
	}

	// get player from query string
	var plrId uint64
	playerQ, err := resp.PlayerFromQ(r, app.MStore.Maps)
	if err != nil {
		app.ErrHTTP(w, err, &rp.Meta, http.StatusUnprocessableEntity)
		return
	}
	plrId = playerQ

	if seasonQ == 0 || seasonQ == 88888 {
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

	// ensure requested args are valid
	stp, err := app.MStore.Maps.ValiSznTmPlr(plrId, tmId, seasonQ)
	if err != nil {
		app.ErrHTTP(w, err, &rp.Meta, http.StatusUnprocessableEntity)
		return
	}

	iq := resp.PQueryIds{PId: stp.PId, TId: stp.TId, SId: stp.SId}

	if err := rp.BuildPlayerRespV2(app.DB, app.MStore.Maps, &iq); err != nil {
		app.ErrHTTP(w, err, &rp.Meta, http.StatusUnprocessableEntity)
		return
	}
	app.WriteJson(w, &rp)
}
