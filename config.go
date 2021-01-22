package api

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type Config struct {
	Redis *redis.Options
}

func (c *Config) ReadConfig() error {
	viper.AutomaticEnv()
	viper.SetConfigFile("./config.yml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}
