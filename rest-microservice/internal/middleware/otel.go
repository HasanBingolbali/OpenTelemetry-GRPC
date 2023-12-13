package middleware

import (
	"OpenTelemetryGRPC/rest-microservice/internal/constant"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"time"
)

func OtelMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		spanOptions := []trace.SpanStartOption{
			trace.WithAttributes(semconv.HTTPMethodKey.String(c.Method())),
			trace.WithAttributes(semconv.HTTPURLKey.String(c.OriginalURL())),
			trace.WithAttributes(semconv.NetHostIPKey.String(c.Get("x-forwarded-for"))),
		}
		spanName := fmt.Sprintf("%s %s", c.Method(), c.OriginalURL())
		ctx, span := constant.Tracer.Start(
			c.Context(),
			spanName,
			spanOptions...,
		)
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*10))
		defer cancel()
		c.Locals(constant.TracerKey, ctx)
		defer span.End()

		err := c.Next()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return nil
	}
}
