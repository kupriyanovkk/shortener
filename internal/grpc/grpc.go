package grpc

import (
	"context"
	"net"

	"github.com/kupriyanovkk/shortener/internal/config"
	pb "github.com/kupriyanovkk/shortener/internal/grpc/proto"
	"google.golang.org/grpc"
)

// ShortenerServer is the server API for Shortener service.
type ShortenerServer struct {
	pb.UnimplementedShortenerServer
	app *config.App
}

// NewShortenerGRPCServer creates a new ShortenerGRPCServer.
//
// It takes an app parameter of type *config.App and returns a ShortenerServer and an error.
func NewShortenerGRPCServer(app *config.App) (s *ShortenerServer, err error) {
	s = &ShortenerServer{
		app: app,
	}

	return s, err
}

// Run runs the ShortenerServer.
//
// It takes a context and a config flags as parameters.
// Returns an error.
func (s *ShortenerServer) Run(ctx context.Context) error {
	listen, err := net.Listen("tcp", s.app.Flags.ServerAddress)
	if err != nil {
		return err
	}

	gRPCServer := grpc.NewServer()

	pb.RegisterShortenerServer(gRPCServer, s)

	return gRPCServer.Serve(listen)
}
