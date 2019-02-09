package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"PROJECT_ROOT/lib/grpctools"
	proto "PROJECT_ROOT/proto/services/example"
	"PROJECT_ROOT/services/example/server"
	"google.golang.org/grpc"
)

func main() {
	log.SetLevel(log.InfoLevel)
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Starting Server")
	grpctools.Run(&grpctools.Options{
		HTTPPort:            os.Getenv("HTTP_PORT"),
		GRPCPort:            os.Getenv("GRPC_PORT"),
		RegisterHandlerFunc: proto.RegisterExampleServiceHandlerFromEndpoint,
		RegisterGRPCServerFunc: func(s *grpc.Server) {
			proto.RegisterExampleServiceServer(s, &server.Server{})
		},
	})
}
