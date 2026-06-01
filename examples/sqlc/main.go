package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ZoneCNH/postgresx/examples/internal/exampleconfig"
	"github.com/ZoneCNH/postgresx/pkg/postgresx"
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
	runtime, err := exampleconfig.FromEnv("postgresx-sqlc-example")
	if err != nil {
		log.Fatal(err)
	}
	if !runtime.Live {
		fmt.Println("postgresx sqlc example dry-run: DBTX boundary available")
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

	q := newQueries(client.Queryer())
	if err := q.Ping(ctx); err != nil {
		log.Fatal(err)
	}
}
