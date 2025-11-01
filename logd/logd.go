package logd

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	INFO    string = "INFO"
	WARNING string = "* WARNING"
	ERROR   string = "** ERROR"
	FATAL   string = "*** FATAL ERROR"
)

type LogdW struct {
	mu  sync.Mutex
	out io.Writer
}
type Logd struct {
	lw *LogdW
	lg *log.Logger
}

func NewLogd(out io.Writer) *Logd {
	lw := LogdWInit(out)
	return &Logd{lw: lw,
		lg: log.New(lw, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logd) log(level, msg string, args ...any) {
	prefix := fmt.Sprintf("[%s]", level)
	l.lg.SetPrefix(prefix)
	msgf := fmt.Sprintf(msg, args...)

	// calldepth set to 3 to catch original caller
	if err := l.lg.Output(3, msgf); err != nil {
		l.lg.Printf("failed to output log msg %s", msgf)
	}
}

// setup LogdW
func LogdWInit(w io.Writer) *LogdW {
	return &LogdW{out: w}
}

// concurent safe write
func (lw *LogdW) Write(p []byte) (n int, err error) {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	return lw.out.Write(p)
}

// create and return the log file
func SetupLogdF(pathfile string) (*os.File, error) {
	ts := time.Now().Format("01022006_150405")
	fname := fmt.Sprintf("%s_%s.log", pathfile, ts)
	f, err := os.Create(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to create file at %s\n**%w", fname, err)
	}
	return f, nil
}

// HIGH LEVEL FUNCS TO CALL IN SOURCE
func (l *Logd) Infof(msg string, args ...any)  { l.log(INFO, msg, args...) }
func (l *Logd) Warnf(msg string, args ...any)  { l.log(WARNING, msg, args...) }
func (l *Logd) Errorf(msg string, args ...any) { l.log(ERROR, msg, args...) }
func (l *Logd) Fatalf(msg string, args ...any) { l.log(ERROR, msg, args...) }
