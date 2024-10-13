package app

import (
	"Backend/business/db/sqlc"
	"Backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Application struct {
	Config  *config.Config
	Logger  *zap.Logger
	Db      *sqlx.DB
	Queries *sqlc.Queries
	Pool    *pgxpool.Pool
}
