package getenv

import "testing"

func TestLoadDotEnv(t *testing.T) {
	if err := LoadDotEnv(); err != nil {
		t.Errorf(`LoadDotEnv failed: %e`, err)
	}

}
