package app

import (
	"Backend/business/db/sqlc"
	"Backend/internal/config"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Application struct {
	Config  *config.Config
	Logger  *zap.Logger
	Db      *sqlx.DB
	Queries *sqlc.Queries
}
