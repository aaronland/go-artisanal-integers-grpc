package grpc

import (
	"context"
	"github.com/aaronland/go-artisanal-integers/client"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type GRPCClient struct {
	client.Client
	client  ArtisanalIntegerServiceClient
	conn    *grpc.ClientConn
	address string
	mu      *sync.RWMutex
	ttl     time.Duration
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

	q := u.Query()

	ttl_secs := 60

	ttl_str := q.Get("ttl")

	if ttl_str != "" {

		t, err := strconv.Atoi(ttl_str)

		if err != nil {
			return nil, err
		}

		ttl_secs = t
	}

	ttl := time.Second * time.Duration(ttl_secs)

	mu := new(sync.RWMutex)

	cl := &GRPCClient{
		address: u.Host,
		mu:      mu,
		ttl:     ttl,
	}

	err = cl.ensureClient(ctx)

	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (c *GRPCClient) NextInt(ctx context.Context) (int64, error) {

	err := c.ensureClient(ctx)

	if err != nil {
		return -1, err
	}

	e := &emptypb.Empty{}

	i, err := c.client.NextInt(ctx, e)

	if err != nil {
		return -1, err
	}

	return i.Integer, nil
}

func (c *GRPCClient) ensureClient(ctx context.Context) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	conn, err := c.connect(ctx)

	if err != nil {
		return err
	}

	client := NewArtisanalIntegerServiceClient(conn)

	c.conn = conn
	c.client = client

	now := time.Now()
	then := now.Add(c.ttl)

	go func(ttl time.Time) {

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:

				if t.After(ttl) {

					c.mu.Lock()
					defer c.mu.Unlock()

					c.conn.Close()

					c.conn = nil
					c.client = nil
				}
			}

			if c.conn == nil {
				break
			}
		}

	}(then)

	return nil
}

func (c *GRPCClient) connect(ctx context.Context) (*grpc.ClientConn, error) {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.DialContext(ctx, c.address, opts...)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
