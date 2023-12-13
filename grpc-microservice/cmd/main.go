package main

import (
	pb "OpenTelemetryGRPC/grpc-microservice/gen"
	"OpenTelemetryGRPC/grpc-microservice/internal/constant"
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"
)

// server is used to implement greeting.GreetingServiceServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements greeting.GreetingServiceServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	ctx, span := constant.Tracer.Start(ctx, "Hello From Grpc")
	time.Sleep(time.Duration(rand.Intn(100)+200) * time.Millisecond)
	defer span.End()
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

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

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
