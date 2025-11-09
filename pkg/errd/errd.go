package errd

import (
	"fmt"
	"runtime"
)

// error types to use throughout repo

type ErrMeta struct {
	Caller string
}

func (m *ErrMeta) GetCaller() {
	pc, _, _, _ := runtime.Caller(1)
	m.Caller = runtime.FuncForPC(pc).Name()
}

type ValidationError struct {
	ErrMeta
	Val any
}

func (e *ValidationError) Error() string {
	e.GetCaller()
	return fmt.Sprintf("%s ERROR | %v is invalid", e.Caller, e.Val)
}

type RequestError struct {
	ErrMeta
	// resp.RespMeta
	Msg string
	Err error
}
