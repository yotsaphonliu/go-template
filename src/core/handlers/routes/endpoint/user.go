package endpoint

import (
	"github.com/gofiber/fiber/v2"
	"go-template/src/core/handlers/render"
	"go-template/src/service"
)

type UserEndpoint interface {
	CreateUser(c *fiber.Ctx) error
	InquiryUserList(c *fiber.Ctx) error
}

type userEndpoint struct {
	Service *service.Service
}

func NewUserEndpoint(sv *service.Service) UserEndpoint {
	return &userEndpoint{
		Service: sv,
	}
}

func (e *userEndpoint) CreateUser(c *fiber.Ctx) error {
	//ctx := e.Service.NewContext(c)

	//params := &model.CreateUserInformationParams{}
	//if err := c.BodyParser(params); err != nil {
	//	return &custom_error.ValidationError{
	//		Code:    custom_error.InvalidJSONString,
	//		Message: "Invalid JSON string",
	//	}
	//}
	//
	//result, err := ctx.CreateUserInformation(params)
	//if err != nil {
	//	return err
	//}

	return render.JSON(c, nil, nil)
}

func (e *userEndpoint) InquiryUserList(c *fiber.Ctx) error {
	//ctx := e.Service.NewContext(c)

	//params := &model.InquiryUserListParams{}
	//if err := c.BodyParser(params); err != nil {
	//	return &custom_error.ValidationError{
	//		Code:    custom_error.InvalidJSONString,
	//		Message: "Invalid JSON string",
	//	}
	//}
	//
	//result, pagination, err := ctx.InquiryUserList(params)
	//if err != nil {
	//	return err
	//}

	return render.JSON(c, nil, nil)
}
