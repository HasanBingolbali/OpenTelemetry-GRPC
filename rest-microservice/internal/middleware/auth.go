package middleware

import (
	"OpenTelemetryGRPC/rest-microservice/internal/util"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Query("user")
		if userID == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		util.SpanFromCtx(c).SetAttributes(attribute.String("user-id", userID))
		return c.Next()
	}
}
