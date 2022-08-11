package grpc

import (
	"context"
	"fmt"
	"github.com/aaronland/go-artisanal-integers/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
