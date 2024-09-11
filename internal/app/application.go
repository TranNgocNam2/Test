package app

import (
	"Backend/db/ent"
	"Backend/db/sqlc"
	"Backend/internal/config"
	"go.uber.org/zap"
)

type Application struct {
	Config    *config.Config
	EntClient *ent.Client
	Logger    *zap.Logger
	Queries   *sqlc.Queries
}
