package postgresql

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string

	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	MaxOpenConns int32
	SSLMode      string
}

func InitConfig() (*Config, error) {
	dbHost := viper.GetString("PG_HOST")
	if dbHost == "" {
		dbHost = viper.GetString("Database.PostgreSQL.Host")
	}

	dbPort := viper.GetString("PG_PORT")
	if dbPort == "" {
		dbPort = viper.GetString("Database.PostgreSQL.Port")
	}

	dbUsername := viper.GetString("PG_USERNAME")
	if dbUsername == "" {
		dbUsername = viper.GetString("Database.PostgreSQL.Username")
	}

	dbPassword := viper.GetString("PG_PASSWORD")
	if dbPassword == "" {
		dbPassword = viper.GetString("Database.PostgreSQL.Password")
	}

	dbDBName := viper.GetString("PG_DB_NAME")
	if dbDBName == "" {
		dbDBName = viper.GetString("Database.PostgreSQL.DBName")
	}

	config := &Config{
		LogLevel: viper.GetString("Database.Log.Level"),

		Host:         dbHost,
		Port:         dbPort,
		Username:     dbUsername,
		Password:     dbPassword,
		DatabaseName: dbDBName,
		MaxOpenConns: viper.GetInt32("Database.PostgreSQL.MaxOpenConns"),
		SSLMode:      viper.GetString("Database.PostgreSQL.SSLMode"),
	}

	if config.LogLevel == "" {
		config.LogLevel = viper.GetString("Log.Level")
	}

	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == "" {
		config.Port = "5432"
	}
	if config.Username == "" {
		config.Username = "postgres"
	}
	if config.Password == "" {
		config.Password = "postgres"
	}
	if config.DatabaseName == "" {
		config.DatabaseName = "evaluation"
	}
	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}

	return config, nil
}

func bulkParamsString(paramPerInsert int, values []interface{}) string {

	sqlParams := ""

	innerPlaceholders := make([]string, 0)
	for i := paramPerInsert; i > 0; i-- {
		if sqlParams == "" {
			sqlParams += "("
		}

		innerPlaceholders = append(innerPlaceholders, fmt.Sprintf("$%d", len(values)-i+1))

	}

	if sqlParams != "" {
		sqlParams += strings.Join(innerPlaceholders, ",")
		sqlParams += ")"
	}

	return sqlParams

}
