package db

import (
	"Backend/api/internal/platform/db/ent"
	"Backend/kit/enum"
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func ConnectDB(ctx context.Context, databaseUrl string, logger *zap.Logger) *ent.Client {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		logger.Fatal(enum.ErrorConnectDB, zap.Error(err))
	}

	if err = db.PingContext(ctx); err != nil {
		logger.Fatal(enum.ErrorConnectDB, zap.Error(err))
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}
