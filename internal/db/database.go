package db

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"gitlab.com/innovia69420/kit/enum/message"
	"go.uber.org/zap"
)

func ConnectDB(ctx context.Context, dsn string, logger *zap.Logger) *sqlx.DB {
	db, err := sqlx.Connect("pgx", dsn)

	if err != nil {
		logger.Fatal(message.FailedConnectDatabase)
	}

	return db
}
