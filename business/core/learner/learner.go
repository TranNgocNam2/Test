package learner

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/password"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	dbSlots, _ := c.queries.GetSlotsByClassId(ctx, dbClass.ID)

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	classLearner := sqlc.AddLearnerToClassParams{
		ID:        uuid.New(),
		ClassID:   dbClass.ID,
		LearnerID: learner.ID,
	}

	err = qtx.AddLearnerToClass(ctx, classLearner)
	if err != nil {
		return err
	}

	for _, dbSlot := range dbSlots {
		err = qtx.GenerateLearnerAttendance(ctx, sqlc.GenerateLearnerAttendanceParams{
			ClassLearnerID: classLearner.ID,
			SlotID:         dbSlot.ID,
		})
		if err != nil {
			return err
		}
	}
	tx.Commit(ctx)
	return nil
}

func (c *Core) JoinSpecialization(ctx *gin.Context, specializationId uuid.UUID) error {
	learner, err := middleware.AuthorizeLearner(ctx, c.queries)
	if err != nil {
		return err
	}

	if learner.Status != Verified {
		return model.ErrUnauthorizedFeatureAccess
	}

	specialization, err := c.queries.GetPublishedSpecializationById(ctx, specializationId)
	if err != nil {
		return model.ErrSpecNotFound
	}

	err = c.queries.AddLearnerToSpecialization(ctx, sqlc.AddLearnerToSpecializationParams{
		LearnerID:        learner.ID,
		SpecializationID: specialization.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
