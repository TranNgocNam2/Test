package db

import (
	"Backend/db/ent"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gitlab.com/innovia69420/kit/enum/message"
	"go.uber.org/zap"
)

func ConnectDB(ctx context.Context, databaseUrl string, logger *zap.Logger) (*ent.Client, *pgxpool.Pool) {
	pool, err := pgxpool.New(context.Background(), databaseUrl)

	if err != nil {
		logger.Fatal(message.FailedConnectDatabase, zap.Error(err))
	}

	db := stdlib.OpenDBFromPool(pool)

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv)), pool
}
