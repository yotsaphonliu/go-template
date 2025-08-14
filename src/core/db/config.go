package db

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBType string
}

func InitConfig() (*Config, error) {
	dbType := viper.GetString("DB_TYPE")
	if dbType == "" {
		dbType = viper.GetString("Database.Type")
	}

	config := &Config{
		DBType: dbType,
	}

	if config.DBType == "" {
		panic("Database.Type not set")
	} else {
		switch config.DBType {
		case "postgres":
		// case "tidb":
		default:
			panic("Unsupported database type")
		}
	}
	return config, nil
}
