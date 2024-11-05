package learner

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/password"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Core struct {
	db      *sqlx.DB
	queries *sqlc.Queries
	logger  *zap.Logger
	pool    *pgxpool.Pool
}

func NewCore(app *app.Application) *Core {
	return &Core{
		db:      app.DB,
		queries: app.Queries,
		logger:  app.Logger,
		pool:    app.Pool,
	}
}

func (c *Core) JoinClass(ctx *gin.Context, classAccess ClassAccess) error {
	learner, err := middleware.AuthorizeLearner(ctx, c.queries)
	if err != nil {
		return err
	}

	if learner.Status != Verified {
		return model.ErrUnauthorizedFeatureAccess
	}

	dbClass, err := c.queries.GetClassCompletedByCode(ctx, classAccess.Code)
	if err != nil {
		return model.ErrClassNotFound
	}
	if dbClass.StartDate.Before(time.Now()) {
		return model.ErrClassStarted
	}

	if !password.Verify(dbClass.Password, classAccess.Password) {
		return model.ErrWrongPassword
	}

	err = c.queries.AddLearnerToClass(ctx, sqlc.AddLearnerToClassParams{
		ClassID:   dbClass.ID,
		LearnerID: learner.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
