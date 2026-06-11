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
	runtime, err := exampleconfig.FromEnv("postgresx-transaction-example")
	if err != nil {
		log.Fatal(err)
	}
	if !runtime.Live {
		fmt.Println("postgresx transaction example dry-run: explicit transaction boundary prepared")
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

	err = client.WithTx(ctx, func(ctx context.Context, tx postgresx.Tx) error {
		_, err := tx.Exec(ctx, "SELECT 1")
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}
