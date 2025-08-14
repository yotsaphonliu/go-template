package service

import (
	"net/http"
	"time"

	"go-template/src/core/model"
	"go-template/src/core/utils"
	"go-template/src/custom_error"
)

type LoginAzureParams struct {
	Code        string `json:"code" validate:"required"`
	RedirectURI string `json:"redirect_uri" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token" example:"T1lSVDdFOGJSb0Q0R2Y3UjhKVlJFeTdHdkNSSm9ZdlRUYUlLVUM4MXpBVHpVNnZC"`
}

type LoginAzureWithAccessTokenParams struct {
	AccessToken string `json:"access_token" validate:"required"`
	RedirectURI string `json:"redirect_uri" validate:"required"`
}

func removeDuplicates(input []string) []string {
	unique := make(map[string]bool)
	result := []string{}

	for _, element := range input {
		if !unique[element] {
			unique[element] = true
			result = append(result, element)
		}
	}

	return result
}

type LoginRootParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ctx *Context) LoginRoot(params LoginRootParams) (*LoginResponse, error) {
	logger := ctx.getLogger("LoginRoot")
	logger.Infof("Begin")
	logger.Debugf("params: %v", params)
	defer logger.Infof("End")

	if err := ValidateInput(params); err != nil {
		logger.Errorf("ValidateInput error : %s", err)
		return nil, err
	}

	apiKey := ""
	if params.Username == ctx.Config.AdminUsername && params.Password == ctx.Config.AdminPassword {
		apiKey = utils.GenerateApiKey()
		userApiKeys := make([]model.ApiKey, 0)
		userInternalRole := make([]string, 0)
		//userInternalRole = append(userInternalRole, string(model.ROLE_ADMIN_ROOT))
		for _, role := range userInternalRole {
			userApiKeys = append(userApiKeys, model.ApiKey{
				Key:          apiKey,
				AzureUserID:  "ADMIN-CONFIG-" + ctx.Config.AdminUsername,
				UserID:       0,
				EmailAddress: "root@mail.com",
				UserRoleName: role,
				ExpireTime:   time.Now().Add(time.Minute * time.Duration(10000)), // FIXME: add config
			})
		}

		err := ctx.DB.InsertApiKeys(userApiKeys, true)
		if err != nil {
			ctx.Logger.Errorf("Login error : %s", err)
			return nil, &custom_error.InternalError{
				Code:    custom_error.DBError,
				Message: err.Error(),
			}
		}
	} else {
		ctx.Logger.Errorf("Inactive user")
		return nil, &custom_error.UserError{
			Code:           custom_error.InvalidUsernameOrPassword,
			Message:        "Invalid Username or Password",
			HTTPStatusCode: http.StatusBadRequest,
		}
	}

	return &LoginResponse{
		Token: apiKey,
	}, nil
}

type GetMeResponse struct {
	EmailAddress string   `json:"email_address"`
	Username     string   `json:"username"`
	FullName     string   `json:"full_name"`
	Department   string   `json:"department"`
	ProfilePic   string   `json:"profile_pic"`
	Role         []string `json:"role"`
}

func (ctx *Context) GetMe() (*GetMeResponse, error) {
	logger := ctx.getLogger("GetMe")
	logger.Infof("Begin")
	defer logger.Infof("End")

	username := "ผู้ดูแลระบบ"
	fullName := "ผู้ดูแลระบบ"
	department := "ผู้ดูแลระบบ"
	if ctx.UserID != 0 {
		//user, err := ctx.DB.GetUserByUserID(ctx.UserID)
		//if err != nil {
		//	return nil, &custom_error.InternalError{
		//		Code:    custom_error.DBError,
		//		Message: err.Error(),
		//	}
		//}
		//
		//username = user.TitleName + " " + user.FirstName + " " + user.LastName
		//fullName = user.TitleName + " " + user.FirstName + " " + user.LastName
		//department = user.DepartmentName
	}

	// return data
	return &GetMeResponse{
		EmailAddress: ctx.EmailAddress,
		Username:     username,
		FullName:     fullName,
		Role:         ctx.Role,
		Department:   department,
		ProfilePic:   ctx.ProfilePic,
	}, nil
}

func (ctx *Context) Logout(token string) error {
	logger := ctx.getLogger("Logout")
	logger.Infof("Begin")
	defer logger.Infof("End")

	err := ctx.DB.DeleteApiKey(token)
	if err != nil {
		return &custom_error.InternalError{
			Code:    custom_error.DBError,
			Message: err.Error(),
		}
	}

	return nil
}

func (ctx *Context) RemoveExpireApiKey() {
	logger := ctx.getLogger("RemoveExpireApiKey")
	logger.Infof("Begin")
	defer logger.Infof("End")

	err := ctx.DB.DeleteExpireApiKey()
	if err != nil {
		logger.Errorf("DeleteExpireApiKey error: %+v", err)
	}
}
