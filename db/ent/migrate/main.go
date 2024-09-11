//go:build ignore

package main

import (
	"Backend/db/ent/migrate"
	"context"
	"fmt"
	"os"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	// Create a local migration directory able to understand Atlas migration file format for replay.
	dir, err := atlas.NewLocalDir("ent/migrate/migrations")

	if err != nil {
		fmt.Println("stupid 1")
	}
	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
		schema.WithFormatter(atlas.DefaultFormatter),
	}
	if len(os.Args) != 2 {
		fmt.Println("stupid 2")
	}
	// Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
	err = migrate.NamedDiff(ctx, "postgres://vzart:pass@localhost:5432/test?sslmode=disable", os.Args[1], opts...)
	if err != nil {
		fmt.Println(err)
	}
}
