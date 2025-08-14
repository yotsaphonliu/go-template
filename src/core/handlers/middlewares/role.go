package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go-template/src/core/handlers/render"
	"go-template/src/service"
)

// RequiredRoles required roles
// TODO: NEW role MODEL
func RequiredRoles(sv *service.Service, roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := sv.NewContext(c)
		userRole := ctx.Role

		for _, role := range roles {
			for _, ur := range userRole {
				if string(role) == ur {
					return c.Next()
				}
			}
		}
		return render.Error(c, fiber.ErrUnauthorized)
	}
}
