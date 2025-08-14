package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go-template/src/core/log"
	"go-template/src/service"
)

type contextKey string

const contextKeyTraceID contextKey = "TraceID"
const contextKeySpanID contextKey = "SpanID"

func CorrelationMiddleware(sv *service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Get("x-request-no")
		c.Locals(contextKeyTraceID, traceID)
		c.Locals(contextKeySpanID, uuid.NewString())
		sv.SpanID = GetSpanID(c)
		sv.TraceID = GetTraceID(c)

		routeName := c.Route().Name
		if routeName != "" {
			c.Set("Service-Code", routeName)
		}

		return c.Next()
	}
}

func ServiceCodeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Proceed to the next handler in the chain
		err := c.Next()

		// Set "Service-Code" header based on the route name, if available
		if routeName := c.Route().Name; routeName != "" {
			c.Set("Service-Code", routeName)
		}

		return err
	}
}

func GetSpanID(c *fiber.Ctx) string {
	spanID := c.Locals(contextKeySpanID)
	if ret, ok := spanID.(string); ok {
		return ret
	}
	return ""
}

func GetTraceID(c *fiber.Ctx) string {
	traceID := c.Locals(contextKeyTraceID)
	if ret, ok := traceID.(string); ok {
		return ret
	}
	return ""
}

func LoggingMiddleware(sv *service.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.OriginalURL() != "/api/health-check" {
			startTime := time.Now()
			appCtx := sv.NewContext(c)
			logger := appCtx.Logger.WithFields(log.Fields{
				"package":   "http_api",
				"remote_ip": c.Context().RemoteIP().String(),
				"method":    c.Method(),
				"path":      c.OriginalURL(),
				"span_id":   GetSpanID(c),
				"trace_id":  GetTraceID(c),
			})

			c.Next()

			duration := time.Since(startTime)
			statusCode := c.Response().StatusCode()
			logger = logger.WithFields(log.Fields{
				"duration":    duration.String(),
				"status_code": statusCode,
			})
			reqBody := masking(CompactJSON(c.Request().Body()))
			resBody := masking(c.Response().Body())

			logger.Debugf("JSON Request: %s", CompactJSON(c.Request().Body()))
			logger.Debugf("JSON Response %s", c.Response().Body())

			err := appCtx.CreateActivityLog(GetTraceID(c), reqBody, resBody)
			if err != nil {
				logger.Errorf("CreateActivityLog error : %v", err)
			}

			if statusCode != http.StatusOK && statusCode != http.StatusCreated && statusCode != http.StatusAccepted {
				logger.Errorf("%s", c.Response().Body())
			}
			logger.Infof("%s %s", c.Method(), c.OriginalURL())
		} else {
			c.Next()
		}

		return nil
	}
}

func CompactJSON(src []byte) []byte {
	var dst bytes.Buffer
	if err := json.Compact(&dst, src); err != nil {
		return nil
	}
	return dst.Bytes()
}

func masking(b []byte) []byte {
	if len(b) > 5000 {
		return []byte(fmt.Sprintf("length is %v bytes", len(b)))
	}
	return b
}
