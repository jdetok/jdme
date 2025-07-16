package logs

import (
	"fmt"
	"net/http"
	"time"
)

// func LogNew()

func LogHTTP(r *http.Request) {
	fmt.Printf("===REQUEST RECEIVED - %v===\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("- Remote Addr: %v\n", r.RemoteAddr)
	fmt.Printf("- Referrer: %v\n", r.Referer())
	fmt.Printf("- User Agent: %v\n", r.UserAgent())
	fmt.Printf("- %v %v\n\n", r.Method, r.RequestURI)
}

func LogDebug(msg string) {
	fmt.Println(msg)
}
