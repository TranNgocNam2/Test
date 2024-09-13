package app

import (
	"Backend/internal/config"
	"Backend/internal/db/sqlc"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Application struct {
	Config  *config.Config
	Logger  *zap.Logger
	Db      *sqlx.DB
	Queries *sqlc.Queries
}
