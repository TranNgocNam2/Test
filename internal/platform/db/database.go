package db

import (
	"Backend/internal/platform/db/ent"
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gitlab.com/innovia69420/kit/enum/message"
	"go.uber.org/zap"
)

func ConnectDB(ctx context.Context, databaseUrl string, logger *zap.Logger) *ent.Client {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		logger.Fatal(message.FailedConnectDatabase, zap.Error(err))
	}

	if err = db.PingContext(ctx); err != nil {
		logger.Fatal(message.FailedConnectDatabase, zap.Error(err))
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}
