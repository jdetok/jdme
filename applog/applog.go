package applog

import (
	"fmt"
	"net/http"
	"time"
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

func LogHTTP(r *http.Request) {
	fmt.Printf("===REQUEST RECEIVED - %v===\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("- Remote Addr: %v\n", r.RemoteAddr)
	fmt.Printf("- Referrer: %v\n", r.Referer())
	fmt.Printf("- User Agent: %v\n", r.UserAgent())
	fmt.Printf("- %v %v\n\n", r.Method, r.RequestURI)
}
