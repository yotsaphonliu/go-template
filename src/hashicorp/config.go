package hashicorp

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	URL         string
	TokenMethod struct {
		Token string
	}
	UserPassMethod struct {
		User     string
		Password string
	}
	AppRoleMethod struct {
		Path     string
		RoleID   string
		SecretID string
	}
	AuthMethod AuthMethod //possible value: token, user_password (default: token)
}

type AuthMethod string

const (
	Token        AuthMethod = "token"
	UserPassword AuthMethod = "user_password"
	AppRole      AuthMethod = "app_role"
)

func InitConfig() (*Config, error) {
	config := &Config{
		URL:        viper.GetString("HashiCorp.URL"),
		AuthMethod: AuthMethod(viper.GetString("HashiCorp.AuthMethod")),
	}

	if config.URL == "" {
		return nil, fmt.Errorf("HashiCorp.URL not set")
	}
	if config.AuthMethod == "" {
		config.AuthMethod = Token
	}

	switch config.AuthMethod {
	case Token, "":
		config.TokenMethod.Token = viper.GetString("HashiCorp.TokenMethod.Token")
		if config.TokenMethod.Token == "" {
			return nil, fmt.Errorf("HashiCorp.TokenMethod.Token not set")
		}

	case UserPassword:
		config.UserPassMethod.User = viper.GetString("HashiCorp.UserPassMethod.User")
		config.UserPassMethod.Password = viper.GetString("HashiCorp.UserPassMethod.Password")

		if config.UserPassMethod.User == "" {
			return nil, fmt.Errorf("HashiCorp.UserPassMethod.User not set")
		}
		if config.UserPassMethod.Password == "" {
			return nil, fmt.Errorf("HashiCorp.UserPassMethod.Password not set")
		}

	case AppRole:
		config.AppRoleMethod.Path = viper.GetString("HashiCorp.AppRoleMethod.Path")
		config.AppRoleMethod.RoleID = viper.GetString("HashiCorp.AppRoleMethod.RoleID")
		config.AppRoleMethod.SecretID = viper.GetString("HashiCorp.AppRoleMethod.SecretID")

		if config.AppRoleMethod.Path == "" {
			return nil, fmt.Errorf("HashiCorp.AppRoleMethod.Path not set")
		}
		if config.AppRoleMethod.RoleID == "" {
			return nil, fmt.Errorf("HashiCorp.AppRoleMethod.RoleID not set")
		}
		if config.AppRoleMethod.SecretID == "" {
			return nil, fmt.Errorf("HashiCorp.AppRoleMethod.SecretID not set")
		}

	default:
		return nil, fmt.Errorf("unknown auth method: %s", config.AuthMethod)
	}

	return config, nil
}
