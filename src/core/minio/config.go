package minio

import (
	"github.com/spf13/viper"
)

type Config struct {
	EndpointUrl string `mapstructure:"EndpointUrl"`
	Username    string `mapstructure:"User"`
	Password    string `mapstructure:"Password"`
	Bucket      string `mapstructure:"BucketName"`
	UseSSL      bool   `mapstructure:"UseSSL"`
}

func InitConfig() (*Config, error) {

	baseConf := struct {
		Conf Config `mapstructure:"Minio"`
	}{}
	err := viper.Unmarshal(&baseConf)
	if err != nil {
		return nil, err
	}

	conf := baseConf.Conf

	return &conf, nil
}
