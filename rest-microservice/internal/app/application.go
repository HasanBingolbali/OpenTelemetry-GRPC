package app

import (
	pb "OpenTelemetryGRPC/rest-microservice/gen"
	"OpenTelemetryGRPC/rest-microservice/internal/constant"
	"OpenTelemetryGRPC/rest-microservice/internal/middleware"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"time"
)

type Application struct {
	*fiber.App
}

func (a *Application) Register() {
	a.Use(middleware.OtelMiddleware())
	a.Use(middleware.AuthMiddleware())
	a.Use(cors.New())
	a.Use(recover.New())
	a.RegisterHTTPRoutes()
}

func (a *Application) RegisterHTTPRoutes() {
	a.Get("/get", func(c *fiber.Ctx) error {
		ctx, _ := c.Locals(constant.TracerKey).(context.Context)
		ctx, span := constant.Tracer.Start(ctx, "ENTERING TO THE GET METHOD")
		defer span.End()
		time.Sleep(time.Duration(rand.Intn(250)+100) * time.Millisecond)
		dbFunc := func(ctx context.Context) {
			_, span := constant.Tracer.Start(ctx, "SQL SELECT")
			defer span.End()
			time.Sleep(time.Duration(rand.Intn(400)) * time.Millisecond)
		}
		dbFunc(ctx)
		grpcCall := func(ctx context.Context) {
			conn, err := grpc.Dial("otel-grpc-go.staging.svc.cluster.local:80",
				grpc.WithInsecure(),
				grpc.WithBlock(),
				grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewGreeterClient(conn)
			time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
			ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			name := "world"
			defer cancel()
			r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			log.Println(r.GetMessage())
		}
		grpcCall(ctx)
		return c.Status(200).JSON(fiber.Map{"follow": "me"})
	})
}
