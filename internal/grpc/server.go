package grpc

import (
	"context"
	"log"
	"net"
	"sync"

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

// Run starts the gRPC server and listens for incoming connections.
func (s *ShortenerServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	listener, err := net.Listen("tcp", s.app.Flags.GRPCServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterShortenerServer(server, s)

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()
		server.GracefulStop()
	}()

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
