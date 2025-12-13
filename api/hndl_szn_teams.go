package api

import (
	"encoding/json"
	"net/http"

	"github.com/jdetok/go-api-jdeko.me/pkg/resp"
)

// /seasons handler
// FOR SEASONS SELECTOR - CALLED ON PAGE LOAD
func (app *App) HndlSeasons(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	season := r.URL.Query().Get("szn")
	w.Header().Set("Content-Type", "application/json")
	if season == "" { // send all szns when szn is not in q str, used most often
		json.NewEncoder(w).Encode(app.Store.Seasons)
	} else {
		for _, szn := range app.Store.Seasons {
			if season == szn.SeasonId { // validate szn from q string
				json.NewEncoder(w).Encode(map[string]string{
					"szn": season,
				})
			}
		}
	}
}

// /teams handler
// FOR TEAMS SELECTOR - CALLED ON PAGE LOAD
func (app *App) HndlTeams(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	team := r.URL.Query().Get("team")
	w.Header().Set("Content-Type", "application/json")
	if team == "" { // send all teams when team is not in q str, used most often
		json.NewEncoder(w).Encode(app.Store.Teams)
	} else { // read & valid team from q string, not yet used 8/6
		for _, tm := range app.Store.Teams {
			if team == tm.TeamAbbr {

				tm.LogoUrl = resp.MakeTeamLogoUrl(tm.League, tm.TeamId)
				json.NewEncoder(w).Encode(tm)
			}
		}
	}
}
