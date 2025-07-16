package main

import (
	"net/http"

	"github.com/jdetok/go-api-jdeko.me/internal/env"
	"github.com/jdetok/go-api-jdeko.me/internal/logs"
)

var fsPath string = env.GetString("STATIC_PATH")

func (app *application) rootHandler(w http.ResponseWriter, r *http.Request) {
	logs.LogHTTP(r)
	w.Header().Set("Cache-Control", "no-store")
	http.FileServer(http.Dir(fsPath)).ServeHTTP(w, r)
}

func (app *application) cssNoCache(w http.ResponseWriter, r *http.Request) {
	logs.LogHTTP(r)
	w.Header().Set("Cache-Control", "no-store")
	http.StripPrefix("/css/", http.FileServer(http.Dir(fsPath+"/css"))).ServeHTTP(w, r)
}

func (app *application) jsNoCache(w http.ResponseWriter, r *http.Request) {
	logs.LogHTTP(r)
	w.Header().Set("Cache-Control", "no-store")
	http.StripPrefix("/js/", http.FileServer(http.Dir(fsPath+"**/js"))).ServeHTTP(w, r)
}
