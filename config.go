package api

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

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

func NewConfig() (*Config, error) {
	c := new(Config)

	viper.AutomaticEnv()
	viper.SetConfigFile("./config.yml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}

	return c, nil
}
