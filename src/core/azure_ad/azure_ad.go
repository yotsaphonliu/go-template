package azure_ad

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	"go-template/src/core/log"
)

type AzureADService interface {
	CreateUserAD(params CreateUserRequest) (string, error)
	GetListGroup() ([]GetGroupResponse, error)
	GetUserListGroup(azureUserID string) ([]GetGroupResponse, error)
	AzureLogin(params AzureLoginParams) (*ProfileMeResponse, string, error)
	GetProfileMe(accessToken string) (*ProfileMeResponse, error)
	GetMeProfilePic(accessToken string) (string, error)
	AddUserToGroup(azureUserID string, azureGroupID string) error
	RemoveUserFromGroup(azureUserID string, azureGroupID string) error
	EnableUserToAzureAD(azureUserID string, enable bool) error
	DeleteUserToAzureAD(azureUserID string) error
	AzureLoginWithAccessToken(params AzureLoginWithADAccessTokenParams) (*ProfileMeResponse, string, error)
}

type AzureADServiceClient struct {
	logger       log.Logger
	app          *azidentity.ClientSecretCredential
	graphService *graph.GraphServiceClient
	config       *Config
}

func New(config *Config, logger log.Logger) (azureADServiceClient *AzureADServiceClient, err error) {
	azureADServiceClient = &AzureADServiceClient{}
	azureADServiceClient.logger = logger.WithFields(log.Fields{
		"module": "azure_ad_service",
	})

	azureADServiceClient.config = config

	// logger.Debugf("\n AD Connection ClientID : %s \n", config.ClientID)
	// logger.Debugf("\n AD Connection ClientSecret : %s \n", config.ClientSecret)
	// logger.Debugf("\n AD Connection TenantID : %s \n", config.TenantID)
	// logger.Debugf("\n AD Connection GraphEndpoint : %s \n", config.GraphEndpoint)

	azureADServiceClient.app = initMSALApp(config)
	azureADServiceClient.graphService = initGraphClient(azureADServiceClient.app)

	return azureADServiceClient, nil
}

func initMSALApp(config *Config) *azidentity.ClientSecretCredential {
	cred, _ := azidentity.NewClientSecretCredential(
		config.TenantID,
		config.ClientID,
		config.ClientSecret,
		nil,
	)

	return cred
}

func initGraphClient(app *azidentity.ClientSecretCredential) *graph.GraphServiceClient {
	graphClient, _ := graph.NewGraphServiceClientWithCredentials(
		app, nil)

	return graphClient
}
