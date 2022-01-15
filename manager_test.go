package config_manager

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestOptionalFunctions(t *testing.T) {
	type Config struct {
		LogLevel string
		RpcPort  string
	}

	if err := os.Setenv("APP_RPCPORT", "8080"); err != nil {
		t.Error(err)
	}

	cfg := &Config{}
	mgr := NewManager(WithEnvPrefix("APP"), WithDefault("env", "LOCAL"), WithDefault("loglevel", "TEST"))

	if err := mgr.Load(cfg); err != nil {
		t.Error(err)
	}

	if err := os.Unsetenv("APP_RPCPORT"); err != nil {
		t.Error(err)
	}

	assert.Equal(t, cfg.LogLevel, "TEST")
	assert.Equal(t, cfg.RpcPort, "8080")
}
