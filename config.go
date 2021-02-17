package api

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// Config holds all values required to log into databases and some extra settings like timeouts.
type Config struct {
	Port      uint
	SharedKey string
	Redis     *redis.Options
	Postgres  *pg.Options
	Timeout   *struct {
		Read     time.Duration
		Write    time.Duration
		Idle     time.Duration
		Shutdown time.Duration
	}
}

// NewConfig creates a new config instance by reading the config.yml file
// at the root of the project. Values inside of this config file will be overwritten
// by environment variables.
func NewConfig() *Config {
	c := new(Config)

	// Use environment variables as well

	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// FIXME: required for now. Change lauch.json or restructure project
	viper.AddConfigPath("..")

	// Read the file
	Must(viper.ReadInConfig())
	Must(viper.Unmarshal(c))

	return c
}
