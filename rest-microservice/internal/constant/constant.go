package constant

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"os"
)

var (
	Tracer    trace.Tracer
	TracerKey = "gofiber"
)

func init() {
	Tracer = otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))
}
