package app

import (
	"Backend/internal/config"
	"Backend/internal/db/ent"
	"Backend/internal/db/sqlc"
	"go.uber.org/zap"
)

type Application struct {
	Config    *config.Config
	EntClient *ent.Client
	Logger    *zap.Logger
	Queries   *sqlc.Queries
}
