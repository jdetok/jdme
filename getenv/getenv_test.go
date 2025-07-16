package getenv

import (
	"fmt"
	"testing"
)

func TestLoadDotEnv(t *testing.T) {
	if err := LoadDotEnv(); err != nil {
		t.Errorf(`LoadDotEnv failed: %e`, err)
	}

	t1, err := GetEnvStr("DB_CONN_STR")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(t1)
}
