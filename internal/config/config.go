package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbConnString string `mapstructure:"DB_CONN_STRING"`
}

func ReadEnvConfig() *Config {
	var config *Config = &Config{}

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return config
}
