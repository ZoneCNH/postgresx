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
	runtime, err := exampleconfig.FromEnv("postgresx-health-example")
	if err != nil {
		log.Fatal(err)
	}
	if !runtime.Live {
		fmt.Println("postgresx health example dry-run: explicit config loaded")
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

	status := client.Check(ctx)
	if status.Status != postgresx.HealthHealthy {
		log.Fatal(status.Message)
	}
}
