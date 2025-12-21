package api

import (
	"net/http"
)

const fsPath string = "/app/static"
const wikiPath string = "/app/wiki/public/docs"
const bballPath string = "/app/static/bball/bball.html"
const abtPath string = "/app/static/about/about.html"
const bballAbtPath string = "/app/static/about/bball_about.html"
const brontoPath string = "/app/static/bronto/bronto.html"

// root URL, serve static directly
func (app *App) HndlRoot(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	w.Header().Set("Cache-Control", "no-store")
	http.FileServer(http.Dir(fsPath)).ServeHTTP(w, r)
}

func (app *App) ServeDocs(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	w.Header().Set("Cache-Control", "no-store")
	fs := http.FileServer(http.Dir(wikiPath))
	http.StripPrefix("/docs/", fs).ServeHTTP(w, r)
}

// /about handler, serves static files from the about directory
func (app *App) HndlAbt(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	http.ServeFile(w, r, abtPath)
}

// /bronto handler, serves static files from the bronto directory
func (app *App) HndlBronto(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	http.ServeFile(w, r, brontoPath)
}

// /bball base handler, serves bball.html
func (app *App) HndlBBall(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	http.ServeFile(w, r, bballPath)
}

// /bball/about handler, serves bball_about.html
func (app *App) HndlBBallAbt(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	http.ServeFile(w, r, bballAbtPath)
}

// prevent css files from caching
func (app *App) CSSNostore(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	w.Header().Set("Cache-Control", "no-store")
	http.StripPrefix("/css/", http.FileServer(http.Dir(fsPath+"/css"))).ServeHTTP(w, r)
}

// prevent js files from caching
func (app *App) JSNostore(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	w.Header().Set("Cache-Control", "no-store")
	http.StripPrefix("/js/", http.FileServer(http.Dir(fsPath+"**/js"))).ServeHTTP(w, r)
}

func (app *App) HndlHealth(w http.ResponseWriter, r *http.Request) {
	app.Lg.HTTPf(r)
	w.WriteHeader(200)
	w.Write([]byte("jdeko.me http server is healthy :)\n"))
}

func (app *App) HndlDBHealth(w http.ResponseWriter, r *http.Request) {
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
