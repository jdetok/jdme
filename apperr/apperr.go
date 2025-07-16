package apperr

import (
	"fmt"
	"net/http"
)

type AppErr struct {
	Process string
	Msg     string
	Err     error
	IsHTTP  bool
	MsgHTTP string
}

func (e *AppErr) BuildError(err error) error {
	return fmt.Errorf("** ERROR IN %s\n-- ***MSG: %s\n ****SOURCE FUNC ERR: %e",
		e.Process, e.Msg, err)
}

func (e *AppErr) HTTPErr(w http.ResponseWriter, err error) {
	e.MsgHTTP = fmt.Sprintf(`*Error occured within jdeko.me API --n%e`, err)
	http.Error(w, e.MsgHTTP, http.StatusInternalServerError)
}
