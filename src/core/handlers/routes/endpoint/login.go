package endpoint

import (
	"github.com/gofiber/fiber/v2"
	"go-template/src/core/handlers/render"
	"go-template/src/core/utils"
	"go-template/src/custom_error"
	"go-template/src/service"
)

type LoginEndpoint interface {
	LoginRoot(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	GetMe(c *fiber.Ctx) error
}

type loginEndpoint struct {
	Service *service.Service
}

func NewLoginEndpoint(sv *service.Service) LoginEndpoint {
	return &loginEndpoint{
		Service: sv,
	}
}

func (ep *loginEndpoint) LoginRoot(c *fiber.Ctx) error {
	ctx := ep.Service.NewContext(c)

	params := &service.LoginRootParams{}
	if err := c.BodyParser(params); err != nil {
		return &custom_error.ValidationError{
			Code:    custom_error.InvalidJSONString,
			Message: "Invalid JSON string",
		}
	}

	result, err := ctx.LoginRoot(*params)
	if err != nil {
		return err
	}

	return render.JSON(c, result, nil)
}

func (ep *loginEndpoint) GetMe(c *fiber.Ctx) error {
	ctx := ep.Service.NewContext(c)

	result, err := ctx.GetMe()
	if err != nil {
		return err
	}

	return render.JSON(c, result, nil)
}

func (ep *loginEndpoint) Logout(c *fiber.Ctx) error {
	ctx := ep.Service.NewContext(c)

	authorizationHeader := c.Get("Authorization")
	bearerToken := utils.ExtractBearerToken(authorizationHeader)
	if bearerToken == "" {
		return render.JSON(c, nil, nil)
	}

	err := ctx.Logout(bearerToken)
	if err != nil {
		return err
	}

	return render.JSON(c, nil, nil)
}
