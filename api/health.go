package api

import (
	"net/http"
)

func (app *App) HndlHealth(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	if err := app.DB.Ping(); err != nil {
		http.Error(w, "db down", 500)
		return
	}
	var recGameDate string
	row := app.DB.QueryRow(`select gdate from stats.pbox order by gdate desc limit 1`)
	if err := row.Scan(&recGameDate); err != nil {
		http.Error(w, "db ping successful and query returned, failed to scan row to string", 500)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("jdeko.me is healthy :)\nmost recent game date in db: " + recGameDate))
}
