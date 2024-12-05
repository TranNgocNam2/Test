package transcript

import (
	"Backend/business/core/learner"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/web/payload"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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

func (c *Core) ChangeScore(ctx *gin.Context, classId uuid.UUID, req payload.UpdateLearnerTranscript) error {
	_, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		c.logger.Error(err.Error())
		return middleware.ErrInvalidUser
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	for _, transcript := range req.Learners {
		classLearner, err := qtx.GetClassLearnerByClassAndLearner(ctx, sqlc.GetClassLearnerByClassAndLearnerParams{
			LearnerID: transcript.LearnerId,
			ClassID:   classId,
		})

		if err != nil {
			c.logger.Error("learner with id: % is not in class")
			return model.LearnerNotInClass
		}

		learnerTranscript, err := qtx.GetLearnerTranscript(ctx, sqlc.GetLearnerTranscriptParams{
			ClassLearnerID: classLearner.ID,
			TranscriptID:   transcript.TranscriptId,
		})

		if err != nil {
			c.logger.Error("learner with id: %s does not have this transcript")
			return err
		}

		err = qtx.UpdateLearnerTranscriptGrade(ctx, sqlc.UpdateLearnerTranscriptGradeParams{
			ID:    learnerTranscript.ID,
			Grade: &transcript.Grade,
		})

		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
