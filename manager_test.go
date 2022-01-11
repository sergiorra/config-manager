package config_manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	type Config struct {
		LogLevel string
		Rpchost  string
	}

	cfg := new(Config)
	mgr := NewManager(WithDefault("env", "LOCAL"), WithDefault("loglevel", "TEST"))

	err := mgr.Load(cfg)
	if err != nil {
		t.Error(err)
	}

	loglevel := mgr.viper.GetString("loglevel")
	env := mgr.viper.GetString("env")

	assert.Equal(t, loglevel, "TEST")
	assert.Equal(t, env, "LOCAL")

}
