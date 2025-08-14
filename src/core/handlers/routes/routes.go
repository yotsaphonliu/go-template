package routes

import (
	ctx "context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go-template/src/core/handlers/middlewares"
	"go-template/src/core/handlers/routes/endpoint"
	"go-template/src/core/log"
	"go-template/src/otel"
	"go-template/src/service"
)

func NewRouter(config *Config, logger log.Logger, sv *service.Service) {

	app := fiber.New()

	app.Use(
		cors.New(),
		otel.Middleware(),
		middlewares.CorrelationMiddleware(sv),
		middlewares.LoggingMiddleware(sv),
		middlewares.WrapError(),
		middlewares.ServiceCodeMiddleware(),
	)

	// Required Auth
	requiredAuth := middlewares.RequiredAuth(sv)

	// Required Role
	requiredOS := middlewares.RequiredRoles(sv, "aaa", "bbb")
	//requiredRoot := middlewares.RequiredRoles(sv, "aaa")

	// Endpoint
	healthCheckEndpoint := endpoint.NewHealthCheckEndpoint(sv)
	//paramsEndpoint := endpoint.NewParameterEndpoint(sv)
	userEndpoint := endpoint.NewUserEndpoint(sv)
	loginEndpoint := endpoint.NewLoginEndpoint(sv)

	api := app.Group("/api")

	api.Get("/health-check", healthCheckEndpoint.HealthCheck)

	api.Post("/root-login", loginEndpoint.LoginRoot)

	api.Get("/me", requiredAuth, loginEndpoint.GetMe)
	api.Post("/logout", requiredAuth, loginEndpoint.Logout)

	// Public api but req azure AD token
	//api.Get("/user-permission", requiredAzureAuth, roleInformationEndpoint.GetUserAllRoleWithPermission)
	//api.Post("/user/active", requiredAzureAuth, userEndpoint.ActiveUserPublic)
	//api.Post("/user/freeze", requiredAzureAuth, userEndpoint.FreezeUserPublic)

	//params := api.Group("/params")
	//{
	//
	//}

	user := api.Group("user", requiredAuth, requiredOS)
	{
		user.Post("/create", userEndpoint.CreateUser).Name("UM02001")
		user.Post("/list", userEndpoint.InquiryUserList).Name("UM02004")
	}

	// Waiting os signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_, cancel := ctx.WithTimeout(ctx.Background(), 5*time.Second)
		defer cancel()

		logger.Infof("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	logger.Infof("Serving HTTP API at http://127.0.0.1:%d", config.Port)
	err := app.Listen(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		logger.Panicf(err.Error())
	}
}
