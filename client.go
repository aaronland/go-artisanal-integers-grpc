package grpc

import (
	"context"
	"github.com/aaronland/go-artisanal-integers/client"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/url"
)

type GRPCClient struct {
	client.Client
	client ArtisanalIntegerServiceClient
}

func init() {
	ctx := context.Background()
	client.RegisterClient(ctx, "grpc", NewGRPCClient)
}

func NewGRPCClient(ctx context.Context, uri string) (client.Client, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(u.Host, opts...)

	if err != nil {
		return nil, err
	}

	client := NewArtisanalIntegerServiceClient(conn)

	cl := &GRPCClient{
		client: client,
	}

	return cl, nil
}

func (c *GRPCClient) NextInt(ctx context.Context) (int64, error) {

	e := &emptypb.Empty{}

	i, err := c.client.NextInt(ctx, e)

	if err != nil {
		return -1, err
	}

	return i.Integer, nil
}
