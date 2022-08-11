package grpc

import (
	"context"
	"fmt"
	"github.com/aaronland/go-artisanal-integers/server"
	"github.com/aaronland/go-artisanal-integers/service"
	"google.golang.org/grpc"
	"net"
	"net/url"
)

func init() {
	ctx := context.Background()
	server.RegisterServer(ctx, "grpc", NewGRPCServer)
}

type GRPCServer struct {
	server.Server
	grpc_server *grpc.Server
	address     string
}

func NewGRPCServer(ctx context.Context, uri string) (server.Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	service_uri := q.Get("service")

	if service_uri == "" {
		return nil, fmt.Errorf("Missing ?server= parameter")
	}

	svc, err := service.NewService(ctx, service_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new service, %w", err)
	}

	wr := &ServiceWrapper{
		service: svc,
	}

	grpc_server := grpc.NewServer()

	RegisterArtisanalIntegerServiceServer(grpc_server, wr)

	s := &GRPCServer{
		address:     u.Host,
		grpc_server: grpc_server,
	}

	return s, nil
}

func (s *GRPCServer) Address() string {
	return s.address
}

func (s *GRPCServer) ListenAndServe(ctx context.Context, args ...interface{}) error {

	lis, err := net.Listen("tcp", s.address)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return s.grpc_server.Serve(lis)
}
