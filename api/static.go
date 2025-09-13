package api

import (
	"net/http"
)

const fsPath string = "/app/static"
const bballPath string = "/app/static/bball/bball.html"
const abtPath string = "/app/static/about/about.html"
const bballAbtPath string = "/app/static/about/bball_about.html"
const brontoPath string = "/app/static/bronto/bronto.html"

// root URL, serve static directly
func (app *App) rootHndl(w http.ResponseWriter, r *http.Request) {
	LogHTTP(r)
	w.Header().Set("Cache-Control", "no-store")
	http.FileServer(http.Dir(fsPath)).ServeHTTP(w, r)
}

// /about handler, serves static files from the about directory
func (app *App) abtHndl(w http.ResponseWriter, r *http.Request) {
	LogHTTP(r)
	http.ServeFile(w, r, abtPath)
}

// /bronto handler, serves static files from the bronto directory
func (app *App) brontoHndl(w http.ResponseWriter, r *http.Request) {
	LogHTTP(r)
	http.ServeFile(w, r, brontoPath)
}

// /bball base handler, serves bball.html
func (app *App) bballHndl(w http.ResponseWriter, r *http.Request) {
	LogHTTP(r)
	http.ServeFile(w, r, bballPath)
}

// /bball/about handler, serves bball_about.html
func (app *App) bballAbtHndl(w http.ResponseWriter, r *http.Request) {
	LogHTTP(r)
	http.ServeFile(w, r, bballAbtPath)
}

// prevent css files from caching
func (app *App) cssNostore(w http.ResponseWriter, r *http.Request) {
	LogHTTP(r)
	w.Header().Set("Cache-Control", "no-store")
	http.StripPrefix("/css/", http.FileServer(http.Dir(fsPath+"/css"))).ServeHTTP(w, r)
}

// prevent js files from caching
func (app *App) jsNostore(w http.ResponseWriter, r *http.Request) {
	LogHTTP(r)
	w.Header().Set("Cache-Control", "no-store")
	http.StripPrefix("/js/", http.FileServer(http.Dir(fsPath+"**/js"))).ServeHTTP(w, r)
}
