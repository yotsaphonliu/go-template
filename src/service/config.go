package service

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	AdminUsername string
	AdminPassword string
	AdminEmail    string
}

func InitConfig() (*Config, error) {
	adminUsername := viper.GetString("ADMIN_USERNAME")
	adminPassword := viper.GetString("ADMIN_PASSWORD")
	adminEmail := viper.GetString("ADMIN_EMAIL")

	if adminUsername == "" {
		adminUsername = viper.GetString("Admin.Username")
	}

	if adminPassword == "" {
		adminPassword = viper.GetString("Admin.Password")
	}

	if adminEmail == "" {
		adminEmail = viper.GetString("Admin.Email")
	}

	config := &Config{
		AdminUsername: adminUsername,
		AdminPassword: adminPassword,
		AdminEmail:    adminEmail,
	}

	if config.AdminUsername == "" {
		return nil, errors.New("ADMIN_USERNAME or Admin.Username config is not set")
	}

	if config.AdminPassword == "" {
		return nil, errors.New("ADMIN_PASSWORD or Admin.Password config is not set")
	}

	if config.AdminEmail == "" {
		return nil, errors.New("ADMIN_EMAIL or Admin.Email config is not set")
	}

	return config, nil
}
