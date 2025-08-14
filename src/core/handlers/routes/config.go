package routes

import "github.com/spf13/viper"

type Config struct {
	Port int
}

func InitConfig() (*Config, error) {
	config := &Config{
		Port: viper.GetInt("API.HTTPServerPort"),
	}

	if config.Port == 0 {
		config.Port = 9092
	}

	return config, nil
}
