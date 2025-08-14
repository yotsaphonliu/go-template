package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go-template/src/core/handlers/render"
	"go-template/src/core/utils"
	"go-template/src/service"
)

func RequiredAuth(sv *service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := sv.NewContext(c)
		authorizationHeader := c.Get("Authorization")
		bearerToken := utils.ExtractBearerToken(authorizationHeader)
		if bearerToken == "" {
			return render.Error(c, fiber.ErrUnauthorized)
		}

		var azureUserID string
		var userID int64
		var role []string
		var emailAddress string
		var profilePic string

		// FIXME: Change to config (Key expire time)
		apiKeyData, err := ctx.DB.VerifyApiKey(bearerToken, time.Now().Add(time.Minute*time.Duration(10000)))
		if err != nil {
			return render.Error(c, fiber.ErrUnauthorized)
		}
		if apiKeyData == nil || len(apiKeyData) == 0 {
			return render.Error(c, fiber.ErrUnauthorized)
		} else {
			// Set user_id and username
			azureUserID = apiKeyData[0].AzureUserID
			userID = apiKeyData[0].UserID
			emailAddress = apiKeyData[0].EmailAddress
			profilePic = apiKeyData[0].UserProfilePic

			for _, r := range apiKeyData {
				role = append(role, r.UserRoleName)
			}
		}

		if azureUserID != "" {
			sv.AzureUserID = azureUserID
			sv.UserID = userID
			sv.Role = role
			sv.EmailAddress = emailAddress
			sv.ProfilePic = profilePic

			return c.Next()
		}

		return render.Error(c, fiber.ErrUnauthorized)
	}
}
