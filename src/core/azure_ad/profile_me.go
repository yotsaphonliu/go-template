package azure_ad

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ProfileMeResponse struct {
	OdataContext      string   `json:"@odata.context"`
	BusinessPhones    []string `json:"businessPhones"`
	DisplayName       string   `json:"displayName"`
	GivenName         string   `json:"givenName"`
	JobTitle          string   `json:"jobTitle"`
	Mail              string   `json:"mail"`
	MobilePhone       string   `json:"mobilePhone"`
	OfficeLocation    string   `json:"officeLocation"`
	PreferredLanguage string   `json:"preferredLanguage"`
	Surname           string   `json:"surname"`
	UserPrincipalName string   `json:"userPrincipalName"`
	ID                string   `json:"id"`
}

type ProfileMeError struct {
	Error struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		InnerError struct {
			Timestamp       time.Time `json:"timestamp"`
			RequestID       string    `json:"request-id"`
			Date            string    `json:"date"`
			ClientRequestID string    `json:"client-request-id"`
		} `json:"innerError"`
	} `json:"error"`
}

func (AzureADServiceClient *AzureADServiceClient) GetProfileMe(accessToken string) (*ProfileMeResponse, error) {
	graphUrl := "https://graph.microsoft.com/v1.0/me"
	req, _ := http.NewRequest(http.MethodGet, graphUrl, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		errResp := ProfileMeError{}
		_ = json.Unmarshal(data, &errResp)
		return nil, fmt.Errorf("invalid status, expect 200 but got %v error : code = %s, msg = %s", resp.StatusCode, errResp.Error.Code, errResp.Error.Message)
	}

	profile := ProfileMeResponse{}
	err = json.Unmarshal(data, &profile)
	if err != nil {
		return &profile, err
	}

	// AzureADServiceClient.logger.Debugf("Profile : %+v", profile)

	return &profile, nil
}

func (AzureADServiceClient *AzureADServiceClient) GetMeProfilePic(accessToken string) (string, error) {
	graphUrl := "https://graph.microsoft.com/v1.0/me/photo/$value"
	req, _ := http.NewRequest(http.MethodGet, graphUrl, nil)
	req.Header.Set("Content-Type", "image/jpg")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		errResp := ProfileMeError{}
		_ = json.Unmarshal(data, &errResp)
		AzureADServiceClient.logger.Errorf("invalid status, expect 200 but got %v error : code = %s, msg = %s", resp.StatusCode, errResp.Error.Code, errResp.Error.Message)
		return "", nil
	}

	base64String := base64.StdEncoding.EncodeToString(data)

	// Print or use the Base64-encoded string as needed
	// AzureADServiceClient.logger.Debugf("Base64-encoded profile picture:", base64String)

	return base64String, nil
}
