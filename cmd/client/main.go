package main

import (
	"context"
	"flag"
	"fmt"
	aa_grpc "github.com/aaronland/go-artisanal-integers-grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func main() {

	address := flag.String("address", "localhost:8080", "")

	flag.Parse()

	ctx := context.Background()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(*address, opts...)

	if err != nil {
		log.Fatalf("fail to dial '%s', %v", *address, err)
	}

	defer conn.Close()

	client := aa_grpc.NewArtisanalIntegerServiceClient(conn)

	e := &emptypb.Empty{}

	i, err := client.NextInt(ctx, e)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(i.Integer)
}
