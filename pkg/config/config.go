package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Can't find the file .env : ", err)
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)

	if err != nil {
		fmt.Println("Environment can't be loaded: ", err)
		return nil, err
	}
	return &config, nil
}
