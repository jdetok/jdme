package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// read player from query string and clean the value (remove accents, lowercase)
func (app *App) PlayerFromQ(r *http.Request) string {
	p := r.URL.Query().Get("player")
	p_lwr := strings.ToLower(p)
	p_cln := RemoveDiacritics(p_lwr)

	fmt.Printf("raw request: %s | cleaned: %s\n", p, p_cln)
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
	playerQ := app.PlayerFromQ(r)
	seasonQ, sznErr := app.SeasonFromQ(r)
	if sznErr != nil {
		http.Error(w, sznErr.Error(), http.StatusUnprocessableEntity)
	}
	var wErr error

	if app.Store.Maps.PNameInSzn(playerQ, seasonQ) {
		_, wErr = fmt.Fprintf(w, "player %s exists in season %d\n", playerQ, seasonQ)
	} else {
		_, wErr = fmt.Fprintf(w, "player %s does not exist in season %d\n", playerQ, seasonQ)
	}

	if wErr != nil {
		http.Error(w,
			fmt.Sprintf("failed to write HTTP response\n**%s", wErr),
			http.StatusInternalServerError)
	}
}
