package charmlogfiber

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Config struct {
	DefaultLevel     log.Level
	ClientErrorLevel log.Level
	ServerErrorLevel log.Level

	WithRequestID bool
}

// New returns a fiber.Handler (middleware) that logs requests using slog.
//
// Requests with errors are logged using slog.Error().
// Requests without errors are logged using slog.Info().
func New(logger *log.Logger) fiber.Handler {
	return NewWithConfig(logger, Config{
		DefaultLevel:     log.InfoLevel,
		ClientErrorLevel: log.WarnLevel,
		ServerErrorLevel: log.ErrorLevel,

		WithRequestID: true,
	})
}

// NewWithConfig returns a fiber.Handler (middleware) that logs requests using slog.
func NewWithConfig(logger *log.Logger, config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Path()
		start := time.Now()
		path := c.Path()

		requestID := uuid.New().String()
		if config.WithRequestID {
			c.Context().SetUserValue("request-id", requestID)
			c.Set("X-Request-ID", requestID)
		}

		err := c.Next()

		end := time.Now()
		latency := end.Sub(start)

		attributes := []interface{}{
			"status", c.Response().StatusCode(),
			"method", string(c.Context().Method()),
			"path", path,
			"ip", c.Context().RemoteIP().String(),
			"latency", latency,
			"user-agent", string(c.Context().UserAgent()),
			"time", end,
		}

		if config.WithRequestID {
			attributes = append(attributes, "request-id", requestID)
		}

		switch {
		case c.Response().StatusCode() >= http.StatusBadRequest:
			if err != nil {
				attributes = append(attributes, "error", err)
			}
			logger.Error("Incoming request with error",
				attributes...,
			)
		default:
			logger.Info("Incoming request", attributes...)
		}

		return err
	}
}

// GetRequestID returns the request identifier
func GetRequestID(c *fiber.Ctx) string {
	requestID, ok := c.Context().UserValue("request-id").(string)
	if !ok {
		return ""
	}

	return requestID
}
