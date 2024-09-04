package app

import (
	"Backend/api/internal/platform/config"
	"Backend/api/internal/platform/db/ent"
	"go.uber.org/zap"
)

type Application struct {
	Config    *config.Config
	EntClient *ent.Client
	Logger    *zap.Logger
}
