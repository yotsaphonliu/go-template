package smtp_service

import (
	"errors"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	SMTPHost       string        `json:"smtp_host"`
	SMTPPort       int           `json:"smtp_port"`
	Username       string        `json:"username"`
	Password       string        `json:"password"`
	From           string        `json:"from"`
	FEEndPoint     string        `json:"fe_endpoint"`
	LinkExpireTime time.Duration `json:"link_expire_time"`
}

func InitConfig() (*Config, error) {
	SMTPHost := viper.GetString("SMTP_HOST")
	if SMTPHost == "" {
		SMTPHost = viper.GetString("SMTP.Host")
	}

	SMTPPort := viper.GetInt("SMTP_PORT")
	if SMTPPort == 0 {
		SMTPPort = viper.GetInt("SMTP.Port")
	}

	SMTPUsername := viper.GetString("SMTP_USERNAME")
	if SMTPUsername == "" {
		SMTPUsername = viper.GetString("SMTP.Username")
	}

	SMTPPassword := viper.GetString("SMTP_PASSWORD")
	if SMTPPassword == "" {
		SMTPPassword = viper.GetString("SMTP.Password")
	}

	SMTPFrom := viper.GetString("SMTP_FROM")
	if SMTPFrom == "" {
		SMTPFrom = viper.GetString("SMTP.From")
	}

	FEEndPoint := viper.GetString("FEEndPoint")
	if FEEndPoint == "" {
		FEEndPoint = viper.GetString("SMTP.FEEndPoint")
	}

	LinkExpireTime := viper.GetDuration("LinkExpireTime")
	if LinkExpireTime == 0 {
		LinkExpireTime = viper.GetDuration("SMTP.LinkExpireTime")
	}

	config := &Config{
		SMTPHost:       SMTPHost,
		SMTPPort:       SMTPPort,
		Username:       SMTPUsername,
		Password:       SMTPPassword,
		From:           SMTPFrom,
		FEEndPoint:     FEEndPoint,
		LinkExpireTime: LinkExpireTime,
	}

	if config.SMTPHost == "" {
		return nil, errors.New("SMTP Host Not found")
	}

	if config.SMTPPort == 0 {
		return nil, errors.New("SMTP Port Not found")
	}

	if config.Username == "" {
		return nil, errors.New("SMTP Username Not found")
	}

	if config.Password == "" {
		return nil, errors.New("SMTP Password Not found")
	}

	if config.From == "" {
		return nil, errors.New("SMTP From Not found")
	}

	return config, nil
}
