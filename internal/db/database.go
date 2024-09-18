package db

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"gitlab.com/innovia69420/kit/enum/message"
	"go.uber.org/zap"
)

func ConnectDB(dsn string, logger *zap.Logger) (db *sqlx.DB, err error) {
	db, err = sqlx.Connect("pgx", dsn)

	if err != nil {
		logger.Error(message.FailedConnectDatabase)
		return nil, err
	}
	return db, nil
}
