package main

import (
	"context"
	"log"

	"github.com/ZoneCNH/postgresx/pkg/postgresx"
	"github.com/ZoneCNH/postgresx/examples/internal/exampleconfig"
)

type sliceSource []postgresx.Migration

func (s sliceSource) List(context.Context) ([]postgresx.Migration, error) {
	return []postgresx.Migration(s), nil
}

func main() {
	ctx := context.Background()
	cfg, err := exampleconfig.FromEnv("postgresx-migration-example")
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

	runner := postgresx.NewMigrationRunner(client)
	err = runner.Up(ctx, sliceSource{
		{
			Version: 1,
			Name:    "create_example_table",
			UpSQL:   "CREATE TABLE IF NOT EXISTS postgresx_example (id BIGSERIAL PRIMARY KEY)",
			DownSQL: "DROP TABLE IF EXISTS postgresx_example",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
