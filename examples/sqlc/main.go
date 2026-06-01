package main

import (
	"context"
	"log"

	"github.com/bytechainx/postgresx"
	"github.com/bytechainx/postgresx/examples/internal/exampleconfig"
)

type queries struct {
	db postgresx.Queryer
}

func newQueries(db postgresx.Queryer) *queries {
	return &queries{db: db}
}

func (q *queries) Ping(ctx context.Context) error {
	_, err := q.db.Exec(ctx, "SELECT 1")
	return err
}

func main() {
	ctx := context.Background()
	cfg, err := exampleconfig.FromEnv("postgresx-sqlc-example")
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

	q := newQueries(client.Queryer())
	if err := q.Ping(ctx); err != nil {
		log.Fatal(err)
	}
}
