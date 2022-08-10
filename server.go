package grpc

import (
	"context"
	"google.golang.org/grpc"
	"fmt"
	"net/url"
	"github.com/aaronland/go-artisanal-integers/server"
	"github.com/aaronland/go-artisanal-integers/service"	
)

type GRPCArtisanalIntegerServer struct {
	server.ArtisanalIntegerServer
	grpc_server
	address string
	service service.Service
}

func init() {
	ctx := context.Background()
	server.RegisterArtisanalIntegerServer(ctx, "gprc", NewGRPCArtisanalIntegerServer)
}

func NewGRPCArtisanalIntegerServer(ctx context.Context, uri string) (server.ArtisanalIntegerServer, error) {

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
	
	grpc_server := grpc.NewServer()

	// register something here...
	
	s := &GRPCArtisanalIntegerServer{
		address: u.Host,
		service: service,
	}

	return s, nil
}

func (s *GRPCArtisanalIntegerServer) Address() string {
	return s.address
}

func (s *GRPCArtisanalIntegerServer) ListenAndServer(ctx context.Context, args ...interface{}) error {

	lis, err := net.Listen("tcp", s.address)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return grpc_server.Serve(lis)
}

