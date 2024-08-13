package db

import (
	db "Backend/api/internal/platform/db/sqlc"
	"Backend/api/internal/platform/logger"
	"Backend/kit/enum"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var DB *pgxpool.Pool
var Queries *db.Queries

func ConnectDB(ctx context.Context, dsn string) {
	var err error
	DB, err = pgxpool.New(ctx, dsn)
	Queries = db.New(DB)
	if err != nil {
		logger.Log.Fatal(enum.ErrorConnectDB, zap.Error(err))
	}
	return
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
