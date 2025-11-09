package resp

// resp package defines & implements all data structures to be returned as
// json in http responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JsonResp interface {
	// New() *JsonResp
	WriteResp(http.ResponseWriter) error
}

type RespMeta struct {
	Request        string   `json:"request"`
	Url            string   `json:"requested_url"`
	ErrorsOccurred uint     `json:"errors_occured"`
	Errors         []string `json:"errors,omitempty"`
}

func NewRespMeta(r *http.Request) *RespMeta {
	return &RespMeta{
		Request:        fmt.Sprintf("%s %s%s", r.Method, r.Host, r.URL.String()),
		Url:            r.URL.String(),
		ErrorsOccurred: 0,
	}
}

func (rm *RespMeta) WriteResp(w http.ResponseWriter) error {
	if err := json.NewEncoder(w).Encode(rm); err != nil {
		return err
	}
	return nil
}

func (rm *RespMeta) CountErr(err error) {
	rm.ErrorsOccurred += 1
	rm.Errors = append(rm.Errors, err.Error())
}
