package main

import (
	"OpenTelemetryGRPC/rest-microservice/internal/app"
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()
	application := app.Application{App: fiber.New()}
	application.Register()
	err = application.Listen(":8080")
	if err != nil {
		log.Fatal("failed to listen: %w", err)
	}

}
