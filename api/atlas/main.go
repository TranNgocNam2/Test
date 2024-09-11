//go:build ignore

package main

import (
	"Backend/internal/config"
	"Backend/internal/db/ent/migrate"
	"context"
	"fmt"
	"gitlab.com/innovia69420/kit/file"
	"os"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	path, err := file.WorkingDirectory()
	if err != nil {
		fmt.Println("Failed getting working directory")
	}
	cfg, err := config.LoadAllAppConfig(path)
	if err != nil {
		fmt.Println("Failed loading app config")
	}

	// Create a local migration directory able to understand Atlas migration file format for replay.
	dir, err := atlas.NewLocalDir("internal/db/ent/migrate/migrations")
	if err != nil {
		fmt.Println("Failed creating migraton file")
		fmt.Println(err.Error())
	}
	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
		schema.WithFormatter(atlas.DefaultFormatter),
	}
	if len(os.Args) != 2 {
		fmt.Println("Arguments exceed expectation")
	}
	// Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
	err = migrate.NamedDiff(ctx, cfg.DatabaseUrl, os.Args[1], opts...)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed creating migration file")
		os.Exit(1)
	}
}
