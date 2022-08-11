package main

import (
	"context"
	"flag"
	_ "github.com/aaronland/go-artisanal-integers-grpc"
	"github.com/aaronland/go-artisanal-integers/server"
	"log"
)

func main() {

	server_uri := flag.String("server-uri", "grpc://localhost:8080?service=memory://", "")

	flag.Parse()

	ctx := context.Background()

	s, err := server.NewArtisanalIntegerServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Failed to create new server, %v", err)
	}

	log.Printf("Listening on %s\n", s.Address())

	err = s.ListenAndServe(ctx)

	if err != nil {
		log.Fatalf("Unable to serve requests, %v", err)
	}
}
