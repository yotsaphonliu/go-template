package endpoint

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go-template/src/core/handlers/render"
	"go-template/src/service"
	"go-template/src/version"
)

type HealthCheckEndpoint interface {
	HealthCheck(c *fiber.Ctx) error
}

type healthCheckEndpoint struct {
	Service   *service.Service
	startTime time.Time
}

func NewHealthCheckEndpoint(sv *service.Service) HealthCheckEndpoint {
	return &healthCheckEndpoint{
		Service:   sv,
		startTime: time.Now(),
	}
}

type HealthCheckServiceDetail struct {
	ServiceName string      `json:"service_name"`
	Status      string      `json:"status"`
	StartTime   string      `json:"start_time,omitempty"`
	UpTime      string      `json:"up_time,omitempty"`
	Version     string      `json:"version,omitempty"`
	Commit      string      `json:"commit,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

const Version = "0.0.1"

func (ep *healthCheckEndpoint) HealthCheck(c *fiber.Ctx) error {
	return render.JSON(c, HealthCheckServiceDetail{
		ServiceName: "go-template",
		Status:      "Online",
		StartTime:   ep.startTime.String(),
		UpTime:      time.Since(ep.startTime).String(),
		Version:     Version,
		Commit:      version.GitCommit,
	}, nil)
}
