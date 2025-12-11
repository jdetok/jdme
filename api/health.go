package api

import (
	"net/http"
	"time"
)

func (app *App) HndlHealth(w http.ResponseWriter, r *http.Request) {
	if time.Since(app.LastUpdate) > 3*time.Hour {
		http.Error(w, "store not updated", 500)
		return
	}
	if err := app.DB.Ping(); err != nil {
		http.Error(w, "db down", 500)
		return
	}
	w.WriteHeader(200)
}
