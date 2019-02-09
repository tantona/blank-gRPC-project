package grpctools

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/reflection"
)

// Options provides an abstraction for GRPC + GRPC Gateway services
type Options struct {
	HTTPPort               string
	GRPCPort               string
	GRPCServerTimeout      time.Duration
	RegisterHandlerFunc    func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) (err error)
	RegisterGRPCServerFunc func(*grpc.Server)
	CustomServer           serverInterface
}

var defaultGRPCServerTimeout = time.Duration(10 * time.Second)

type serverInterface interface {
	ServeGRPC()
	ServeHTTP() error
	CheckGRPCConnectivity(*Options) error
}

type server struct {
	Context  context.Context
	Listener net.Listener
	Options  *Options
}

func (srv *server) ServeGRPC() {
	s := grpc.NewServer()
	srv.Options.RegisterGRPCServerFunc(s)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Infof("Running gRPC server on localhost%s", srv.Options.GRPCPort)

	if err := s.Serve(srv.Listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (srv *server) ServeHTTP() error {
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := srv.Options.RegisterHandlerFunc(srv.Context, mux, fmt.Sprintf("localhost%s", srv.Options.GRPCPort), opts)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Infof("Running http server on localhost%s", srv.Options.HTTPPort)
	return http.ListenAndServe(srv.Options.HTTPPort, mux)

}

func (srv *server) CheckGRPCConnectivity(options *Options) error {
	if options.GRPCServerTimeout == time.Duration(0) {
		options.GRPCServerTimeout = defaultGRPCServerTimeout
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}
	log.Debug("Checking gRPC Server connectivity...")
	conn, err := grpc.Dial(fmt.Sprintf("localhost%s", options.GRPCPort), opts...)
	defer conn.Close()
	if err != nil {
		return err
	}

	startTime := time.Now()
	for {
		state := conn.GetState()
		log.Debugf("gRPC ClientConn State: %s", state)
		if state == connectivity.Ready {
			log.Debug("gRPC Server Ready")
			break
		}
		since := time.Since(startTime)
		if since >= options.GRPCServerTimeout {
			return fmt.Errorf("starting gRPC server timed out after %s", options.GRPCServerTimeout)
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

// Run Registers the grpc + http services specified in Options
func Run(options *Options) {
	var s serverInterface
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if options.CustomServer != nil {
		s = options.CustomServer
	} else {
		lis, err := net.Listen("tcp", options.GRPCPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s = &server{
			Listener: lis,
			Options:  options,
			Context:  ctx,
		}
	}

	// If HTTP Options are defined then spin off gRPC server on a separate thread
	if options.HTTPPort != "" && options.RegisterHandlerFunc != nil {
		log.Debug("HTTP_PORT and RegisterHandlerFunc are defined - starting gRPC with HTTP Proxy")
		go s.ServeGRPC()

		err := s.CheckGRPCConnectivity(options)
		if err != nil {
			log.Fatal(err)
		}

		if err := s.ServeHTTP(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// If HTTP server not needed - just serve GRPC
	s.ServeGRPC()
}
