package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jdetok/golib/logd"
)

// basic log for http requests
func LogHTTP(r *http.Request) {
	logd.Logc(fmt.Sprintf(`
+++ REQUEST RECEIVED - %v
- Request URL: %v
- Method: %v | Request URI: %v
- Referrer: %v
- Remote Addr: %v
- User Agent: %v`,
		time.Now().Format("2006-01-02 15:04:05"),
		r.URL,
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
		r.Referer(),
		r.UserAgent()),
	)
}
