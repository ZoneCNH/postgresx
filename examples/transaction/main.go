package main

import (
	"context"
	"log"

	"github.com/bytechainx/postgresx"
	"github.com/bytechainx/postgresx/examples/internal/exampleconfig"
)

func main() {
	ctx := context.Background()
	cfg, err := exampleconfig.FromEnv("postgresx-transaction-example")
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

	err = client.WithTx(ctx, func(ctx context.Context, tx postgresx.Tx) error {
		_, err := tx.Exec(ctx, "SELECT 1")
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}
