package app

import (
	"Backend/internal/platform/config"
	"Backend/internal/platform/db/ent"
	"go.uber.org/zap"
)

type Application struct {
	Config    *config.Config
	EntClient *ent.Client
	Logger    *zap.Logger
}
