package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ZoneCNH/postgresx/examples/internal/exampleconfig"
	"github.com/ZoneCNH/postgresx/pkg/postgresx"
)

func main() {
	ctx := context.Background()
	runtime, err := exampleconfig.FromEnv("postgresx-basic-example")
	if err != nil {
		log.Fatal(err)
	}
	if !runtime.Live {
		fmt.Println("postgresx basic example dry-run: explicit config loaded")
		return
	}
	client, err := postgresx.Open(ctx, runtime.Config)
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
