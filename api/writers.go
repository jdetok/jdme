package api

import (
	"net/http"

	"github.com/jdetok/go-api-jdeko.me/pkg/resp"
)

func (app *App) WriteJson(w http.ResponseWriter, jr resp.JsonResp) {
	w.Header().Set("Content-Type", "application/json")
	if err := jr.WriteResp(w); err != nil {
		app.Lg.Errorf("JSON write failed: %v", err)
		http.Error(w, "Failed to write JSON response", http.StatusInternalServerError)
	}
}

func (app *App) WriteJsonErr(w http.ResponseWriter, jr resp.JsonResp) error {
	w.Header().Set("Content-Type", "application/json")
	if err := jr.WriteResp(w); err != nil {
		return err
	}
	return nil
}

// accept slice of bytes in JSON structure and write to response writers
func (app *App) JSONWriter(w http.ResponseWriter, js []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (app *App) ErrHTTP(w http.ResponseWriter, err error, rm *resp.RespMeta, code int) {
	rm.CountErr(err)
	app.Lg.Errorf("** ERROR OCCURED: %v", err.Error())
	if wErr := app.WriteJsonErr(w, rm); wErr != nil {
		http.Error(w, err.Error(), code)
	}
}
