package main

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/aaronland/go-artisanal-integers-grpc"
	"github.com/aaronland/go-artisanal-integers/client"
	"log"
)

func main() {

	client_uri := flag.String("client-uri", "grpc://localhost:8080", "")

	flag.Parse()

	ctx := context.Background()

	cl, err := client.NewClient(ctx, *client_uri)

	if err != nil {
		log.Fatalf("Failed to create new client, %v", err)
	}

	i, err := cl.NextInt(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(i)
}
