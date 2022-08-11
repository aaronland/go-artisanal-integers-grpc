package grpc

import (
	"context"
	"fmt"
	"github.com/aaronland/go-artisanal-integers/server"
	"github.com/aaronland/go-artisanal-integers/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"net/url"
)

type GRPCArtisanalIntegerServer struct {
	server.ArtisanalIntegerServer
	grpc_server *grpc.Server
	address     string
}

type ServiceWrapper struct {
	ArtisanalIntegerServiceServer
	service service.Service
}

func (wr *ServiceWrapper) NextInt(ctx context.Context, e *emptypb.Empty) (*ArtisanalInteger, error) {

	v, err := wr.service.NextInt(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to generate next int, %w", err)
	}

	i := &ArtisanalInteger{
		Integer: v,
	}

	return i, nil
}

func init() {
	ctx := context.Background()
	server.RegisterArtisanalIntegerServer(ctx, "grpc", NewGRPCArtisanalIntegerServer)
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

	wr := &ServiceWrapper{
		service: svc,
	}

	grpc_server := grpc.NewServer()

	RegisterArtisanalIntegerServiceServer(grpc_server, wr)

	s := &GRPCArtisanalIntegerServer{
		address:     u.Host,
		grpc_server: grpc_server,
	}

	return s, nil
}

func (s *GRPCArtisanalIntegerServer) Address() string {
	return s.address
}

func (s *GRPCArtisanalIntegerServer) ListenAndServe(ctx context.Context, args ...interface{}) error {

	lis, err := net.Listen("tcp", s.address)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return s.grpc_server.Serve(lis)
}
