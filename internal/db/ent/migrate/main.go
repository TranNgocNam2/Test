//go:build ignore

package main

import (
	"Backend/internal/db/ent/migrate"
	"context"
	"fmt"
	"os"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
	"gitlab.com/innovia69420/kit/logger"
)

func main() {
	ctx := context.Background()
	lg := logger.FromCtx(ctx)
	// Create a local migration directory able to understand Atlas migration file format for replay.
	dir, err := atlas.NewLocalDir("internal/db/ent/migrate/migrations")
	if err != nil {
		fmt.Println("stupid 1")
		lg.Sugar().Fatalf("Failed creating migraton file")
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
		lg.Sugar().Fatalf("Arguments exceed expectation")
	}
	// Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
	err = migrate.NamedDiff(ctx, "postgres://vzart:pass@localhost:5432/test?sslmode=disable", os.Args[1], opts...)
	if err != nil {
		fmt.Println(os.Getenv("DB_DSN"))
		fmt.Println(err)
		lg.Sugar().Fatalf("Failed creating migration file")
	}
}
