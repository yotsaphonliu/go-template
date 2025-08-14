package azure_ad

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	ClientID      string
	ClientSecret  string
	TenantID      string
	GraphEndpoint string
}

func InitConfig() (*Config, error) {
	azClientID := viper.GetString("AZ_CLIENT_ID")
	if azClientID == "" {
		azClientID = viper.GetString("AzureAD.ClientID")
	}

	azClientSecret := viper.GetString("AZ_CLIENT_SECRET")
	if azClientSecret == "" {
		azClientSecret = viper.GetString("AzureAD.ClientSecret")
	}

	azTenantID := viper.GetString("AZ_TENANT_ID")
	if azTenantID == "" {
		azTenantID = viper.GetString("AzureAD.TenantID")
	}

	azGraphEndpoint := viper.GetString("AZ_GRAPH_ENDPOINT")
	if azGraphEndpoint == "" {
		azGraphEndpoint = viper.GetString("AzureAD.GraphEndpoint")
	}

	config := &Config{
		ClientID:      azClientID,
		ClientSecret:  azClientSecret,
		TenantID:      azTenantID,
		GraphEndpoint: azGraphEndpoint,
	}

	if config.ClientID == "" {
		return nil, errors.New("ClientID Not found")
	}

	if config.ClientSecret == "" {
		return nil, errors.New("ClientSecret Not found")
	}

	if config.TenantID == "" {
		return nil, errors.New("TenantID Not found")
	}

	if config.GraphEndpoint == "" {
		return nil, errors.New("GraphEndpoint Not found")
	}

	return config, nil
}
