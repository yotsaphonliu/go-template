package azure_ad

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type AzureToken struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

type AzureTokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorCodes       []int  `json:"error_codes"`
	Timestamp        string `json:"timestamp"`
	TraceId          string `json:"trace_id"`
	CorrelationId    string `json:"correlation_id"`
}

type AzureLoginParams struct {
	Code        string `json:"code" validate:"required"`
	RedirectURI string `json:"redirect_uri" validate:"required"`
	//AccessToken string
}

func (AzureADServiceClient *AzureADServiceClient) AzureLogin(params AzureLoginParams) (*ProfileMeResponse, string, error) {

	azureURL := "https://login.microsoftonline.com"
	tenantID := "/" + AzureADServiceClient.config.TenantID
	resource := "/oauth2/v2.0/token"
	data := url.Values{}
	data.Set("client_id", AzureADServiceClient.config.ClientID)
	data.Set("redirect_uri", params.RedirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("client_secret", AzureADServiceClient.config.ClientSecret)
	data.Set("code", params.Code)

	u, _ := url.ParseRequestURI(azureURL)
	u.Path = tenantID + resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		return nil, "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errResp := AzureTokenError{}
		_ = json.Unmarshal(body, &errResp)
		return nil, "", fmt.Errorf("invalid status, expect 200 but got %v, error : %v", resp.StatusCode, errResp.ErrorDescription)
	}

	token := AzureToken{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, "", err
	}
	if token.AccessToken == "" {
		return nil, "", errors.New("Invalid code")
	}

	profile, err := AzureADServiceClient.GetProfileMe(token.AccessToken)
	if err != nil {
		return nil, "", err
	}

	profilePic, err := AzureADServiceClient.GetMeProfilePic(token.AccessToken)
	if err != nil {
		return nil, "", err
	}

	return profile, profilePic, nil
}

type AzureLoginWithADAccessTokenParams struct {
	AccessToken string `json:"access_token" validate:"required"`
	RedirectURI string `json:"redirect_uri" validate:"required"`
	//AccessToken string
}

func (AzureADServiceClient *AzureADServiceClient) AzureLoginWithAccessToken(params AzureLoginWithADAccessTokenParams) (*ProfileMeResponse, string, error) {
	profile, err := AzureADServiceClient.GetProfileMe(params.AccessToken)
	if err != nil {
		return nil, "", err
	}

	profilePic, err := AzureADServiceClient.GetMeProfilePic(params.AccessToken)
	if err != nil {
		return nil, "", err
	}

	return profile, profilePic, nil
}
