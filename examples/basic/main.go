package main

import (
	"context"
	"log"

	"github.com/bytechainx/postgresx"
	"github.com/bytechainx/postgresx/examples/internal/exampleconfig"
)

func main() {
	ctx := context.Background()
	cfg, err := exampleconfig.FromEnv("postgresx-basic-example")
	if err != nil {
		log.Fatal(err)
	}
	client, err := postgresx.Open(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Close(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	if err := client.Ping(ctx); err != nil {
		log.Fatal(err)
	}
}
