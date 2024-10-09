package subject

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Core struct {
	db      *sqlx.DB
	queries *sqlc.Queries
	logger  *zap.Logger
}

func NewCore(app *app.Application) *Core {
	return &Core{
		db:      app.Db,
		queries: app.Queries,
		logger:  app.Logger,
	}
}
