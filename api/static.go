package main

import (
	"net/http"

	"github.com/jdetok/go-api-jdeko.me/applog"
)

// testing 7/20/25 with frontend in separate directory

const fsPath string = "/app/static"
const bballPath string = "/app/static/bball/bball.html"
const abtPath string = "/app/static/about/about.html"
const bballAbtPath string = "/app/static/about/bball_about.html"
const brontoPath string = "/app/static/bronto/bronto.html"

// const fsPath string = "/static"
// const bballPath string = "/static/bball/bball.html"
// const abtPath string = "/static/about/about.html"
// const brontoPath string = "/static/bronto/bronto.html"

func (app *application) rootHandler(w http.ResponseWriter, r *http.Request) {
	applog.LogHTTP(r)
	w.Header().Set("Cache-Control", "no-store")
	http.FileServer(http.Dir(fsPath)).ServeHTTP(w, r)
}

func (app *application) abtHandler(w http.ResponseWriter, r *http.Request) {
	applog.LogHTTP(r)
	http.ServeFile(w, r, abtPath)
}

func (app *application) brontoHandler(w http.ResponseWriter, r *http.Request) {
	applog.LogHTTP(r)
	http.ServeFile(w, r, brontoPath)
}

func (app *application) bballHandler(w http.ResponseWriter, r *http.Request) {
	applog.LogHTTP(r)
	http.ServeFile(w, r, bballPath)
}

func (app *application) bballAbtHandler(w http.ResponseWriter, r *http.Request) {
	applog.LogHTTP(r)
	http.ServeFile(w, r, bballAbtPath)
}

func (app *application) cssNoCache(w http.ResponseWriter, r *http.Request) {
	applog.LogHTTP(r)
	w.Header().Set("Cache-Control", "no-store")
	http.StripPrefix("/css/", http.FileServer(http.Dir(fsPath+"/css"))).ServeHTTP(w, r)
}

func (app *application) jsNoCache(w http.ResponseWriter, r *http.Request) {
	applog.LogHTTP(r)
	w.Header().Set("Cache-Control", "no-store")
	http.StripPrefix("/js/", http.FileServer(http.Dir(fsPath+"**/js"))).ServeHTTP(w, r)
}
