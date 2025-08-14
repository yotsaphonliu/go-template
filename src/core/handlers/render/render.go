package render

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go-template/src/core/result"
	"go-template/src/custom_error"
)

// JSON render json to client
func JSON(c *fiber.Ctx, response interface{}, pagination interface{}) error {
	return c.
		Status(http.StatusOK).
		JSON(result.Result{
			Code:       200,
			Message:    "OK",
			Pagination: pagination,
			Data:       response,
		})
}

// Byte render byte to client
func Byte(c *fiber.Ctx, bytes []byte) error {
	_, err := c.Status(http.StatusOK).
		Write(bytes)

	return err
}

// Error render error to client
func Error(c *fiber.Ctx, err error) error {

	if locErr, ok := err.(result.Result); ok {
		return c.
			Status(locErr.HTTPStatusCode()).
			JSON(locErr)
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		return c.
			Status(fiberErr.Code).
			JSON(result.NewResultWithMessage(fiberErr.Error()))
	}

	if customErr, ok := err.(*custom_error.ValidationError); ok {
		return c.
			Status(http.StatusBadRequest).
			JSON(customErr)
	}

	if customErr, ok := err.(*custom_error.AuthorizationError); ok {

		if customErr.HTTPStatusCode != 0 {
			return c.
				Status(customErr.HTTPStatusCode).
				JSON(customErr)
		}

		return c.
			Status(http.StatusUnauthorized).
			JSON(customErr)
	}

	if customErr, ok := err.(*custom_error.UserError); ok {
		return c.
			Status(http.StatusBadRequest).
			JSON(customErr)
	}

	defaultErr := result.Result{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
	return c.
		Status(defaultErr.HTTPStatusCode()).
		JSON(defaultErr)
}
