package api

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// Config holds all values required to log into databases and some extra settings like timeouts.
type Config struct {
	Redis    *redis.Options
	Postgres *pg.Options
	Timeout  *struct {
		Read     time.Duration
		Write    time.Duration
		Idle     time.Duration
		Shutdown time.Duration
	}
}

// NewConfig creates a new config instance by reading the config.yml file
// at the root of the project. Values inside of this config file will be overwritten
// by environment variables.
func NewConfig() (*Config, error) {
	c := new(Config)

	// Use environment variables as well
	viper.AutomaticEnv()
	viper.SetConfigFile("./config.yml")

	// Read the file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Tries to deserialize the yaml-file
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}

	return c, nil
}
