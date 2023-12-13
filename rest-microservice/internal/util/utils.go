package util

import (
	"OpenTelemetryGRPC/rest-microservice/internal/constant"
	"context"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

func SpanFromCtx(c *fiber.Ctx) trace.Span {
	ctx, ok := c.Locals(constant.TracerKey).(context.Context)
	if !ok {
		return trace.SpanFromContext(c.Context())
	}
	return trace.SpanFromContext(ctx)
}
