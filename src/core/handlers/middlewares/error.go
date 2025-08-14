package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go-template/src/core/handlers/render"
)

// WrapError wrap error
func WrapError() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return render.Error(c, err)
		}
		return nil
	}
}
