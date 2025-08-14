package otel

import (
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
)

func Middleware(opts ...otelfiber.Option) fiber.Handler {
	return otelfiber.Middleware(opts...)
}
