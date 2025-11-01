package api

import (
	"fmt"
	"os"
	"time"
)

func SetupLogf(pathfile string) (*os.File, error) {
	ts := time.Now().Format("01022006_0506")
	fname := fmt.Sprintf("%s_%s.log", pathfile, ts)
	f, err := os.Create(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to create file at %s\n**%w", fname, err)
	}
	return f, nil
}

func WriteLogf(f *os.File, msg string) {
	defer f.Close()
	if _, err := f.Write([]byte(msg)); err != nil {
		fmt.Println("failed to write log")
	}
}
