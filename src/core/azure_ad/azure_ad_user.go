package azure_ad

import (
	"context"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type CreateUserRequest struct {
	DisplayName       string `json:"display_name"`
	MailNickname      string `json:"mail_nickname"`
	UserPrincipalName string `json:"user_principal_name"`
	Password          string `json:"password"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

func (AzureADServiceClient *AzureADServiceClient) CreateUserAD(params CreateUserRequest) (string, error) {
	requestBody := graphmodels.NewUser()
	accountEnabled := true
	requestBody.SetAccountEnabled(&accountEnabled)
	displayName := params.DisplayName
	requestBody.SetDisplayName(&displayName)
	mailNickname := params.MailNickname
	requestBody.SetMailNickname(&mailNickname)
	userPrincipalName := params.UserPrincipalName // "TestCreateFromAPI_1@blockfintthailand.onmicrosoft.com"
	requestBody.SetUserPrincipalName(&userPrincipalName)
	passwordProfile := graphmodels.NewPasswordProfile()
	forceChangePasswordNextSignIn := true
	passwordProfile.SetForceChangePasswordNextSignIn(&forceChangePasswordNextSignIn)
	password := params.Password
	passwordProfile.SetPassword(&password)
	requestBody.SetPasswordProfile(passwordProfile)

	users, err := AzureADServiceClient.graphService.Users().Post(context.Background(), requestBody, nil)
	if err != nil {
		return "", err
	}

	return *users.GetId(), nil
}

func (AzureADServiceClient *AzureADServiceClient) EnableUserToAzureAD(azureUserID string, enable bool) error {
	requestBody := graphmodels.NewUser()
	accountEnabled := enable
	requestBody.SetAccountEnabled(&accountEnabled)

	_, err := AzureADServiceClient.graphService.Users().ByUserId(azureUserID).Patch(context.Background(), requestBody, nil)
	if err != nil {
		return err
	}

	return nil
}

func (AzureADServiceClient *AzureADServiceClient) DeleteUserToAzureAD(azureUserID string) error {
	err := AzureADServiceClient.graphService.Users().ByUserId(azureUserID).Delete(context.Background(), nil)
	if err != nil {
		return err
	}

	return nil
}
