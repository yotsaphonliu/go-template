package azure_ad

import (
	"context"
	"fmt"

	"github.com/microsoftgraph/msgraph-sdk-go/groups"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type GetListGroupRequest struct {
	DisplayName       string `json:"display_name"`
	MailNickname      string `json:"mail_nickname"`
	UserPrincipalName string `json:"user_principal_name"`
	Password          string `json:"password"`
}

type GetGroupResponse struct {
	GroupID          string
	GroupDisplayName string
	GroupDescription string
}

func (AzureADServiceClient *AzureADServiceClient) GetListGroup() ([]GetGroupResponse, error) {
	//Create options to set $top=999
	top := int32(999)
	requestParams := &groups.GroupsRequestBuilderGetQueryParameters{
		Top: &top,
	}
	requestConfig := &groups.GroupsRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParams,
	}

	request := AzureADServiceClient.graphService.Groups()
	var resp []GetGroupResponse

	responsible, err := request.Get(context.Background(), requestConfig)
	if err != nil {
		return nil, err
	}

	for _, group := range responsible.GetValue() {
		var out GetGroupResponse

		id := group.GetId()
		if id != nil {
			out.GroupID = *id
		}
		name := group.GetDisplayName()
		if name != nil {
			out.GroupDisplayName = *name
		}

		description := group.GetDescription()
		if description != nil {
			out.GroupDescription = *description
		}

		resp = append(resp, out)
	}
	// Handle pagination
	nextLink := responsible.GetOdataNextLink()
	for nextLink != nil {
		// Create a new request builder for the next link
		nextPageRequest := groups.NewGroupsRequestBuilder(*nextLink, AzureADServiceClient.graphService.RequestAdapter)

		// Fetch the next page
		nextResponse, err := nextPageRequest.Get(context.Background(), nil)
		if err != nil {
			return nil, err
		}

		for _, group := range nextResponse.GetValue() {
			var out GetGroupResponse

			id := group.GetId()
			if id != nil {
				out.GroupID = *id
			}
			name := group.GetDisplayName()
			if name != nil {
				out.GroupDisplayName = *name
			}

			description := group.GetDescription()
			if description != nil {
				out.GroupDescription = *description
			}

			resp = append(resp, out)
		}

		// Update the nextLink for the loop
		nextLink = nextResponse.GetOdataNextLink()
	}

	return resp, nil
}

func (AzureADServiceClient *AzureADServiceClient) GetUserListGroup(azureUserID string) ([]GetGroupResponse, error) {

	responsible, err := AzureADServiceClient.graphService.Users().ByUserId(azureUserID).MemberOf().Get(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	var resp []GetGroupResponse
	for _, group := range responsible.GetValue() {
		var out GetGroupResponse

		id := group.GetId()
		if id != nil {
			out.GroupID = *id
		}
		resp = append(resp, out)
	}

	return resp, nil
}

func (AzureADServiceClient *AzureADServiceClient) AddUserToGroup(azureUserID string, azureGroupID string) error {
	requestBody := graphmodels.NewReferenceCreate()
	odataId := fmt.Sprintf("https://graph.microsoft.com/v1.0/directoryObjects/%s", azureUserID)
	requestBody.SetOdataId(&odataId)

	err := AzureADServiceClient.graphService.Groups().ByGroupId(azureGroupID).Members().Ref().Post(context.Background(), requestBody, nil)
	if err != nil {
		AzureADServiceClient.logger.Errorf("AddUserToGroup Error : %s", err)
		return err
	}

	return nil
}

func (AzureADServiceClient *AzureADServiceClient) RemoveUserFromGroup(azureUserID string, azureGroupID string) error {
	err := AzureADServiceClient.graphService.Groups().ByGroupId(azureGroupID).Members().ByDirectoryObjectId(azureUserID).Ref().Delete(context.Background(), nil)
	if err != nil {
		return err
	}

	return nil
}
