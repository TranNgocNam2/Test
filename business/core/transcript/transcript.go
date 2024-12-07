package transcript

import (
	"Backend/business/core/learner/certificate"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/web/payload"
	"math"
	"time"

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

func (c *Core) ChangeScore(ctx *gin.Context, classId uuid.UUID, req []payload.LearnerTranscript) error {
	_, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		return middleware.ErrInvalidUser
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	for _, transcript := range req {
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

func (c *Core) SubmitScore(ctx *gin.Context, classId uuid.UUID) error {
	_, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		return middleware.ErrInvalidUser
	}

	class, err := c.queries.GetClassById(ctx, classId)
	if err != nil {
		return model.ErrClassNotFound
	}

	subject, err := c.queries.GetSubjectById(ctx, class.SubjectID)
	if err != nil {
		return model.ErrSubjectNotFound
	}

	classLearners, err := c.queries.GetLearnersByClassId(ctx, class.ID)
	if err != nil {
		return model.CannotGetAllLearners
	}

	for _, learner := range classLearners {
		transcripts, err := c.queries.GetLearnerTranscriptByClassLearnerId(learner.ClassLearnerID)
		if err != nil {
			return err
		}
		var totalGrade float64
		totalGrade = 0
		for _, transcript := range transcripts {
			if float64(*transcript.Grade) < transcript.MinGrade {
				// Update transcript status
				if err = c.queries.UpdateTranscriptStatus(ctx, sqlc.UpdateTranscriptStatusParams{
					ClassLearnerID: learner.ClassLearnerID,
					TranscriptID:   transcript.TranscriptID,
					Status:         0,
				}); err != nil {
					return err
				}
			} else {
				if err = c.queries.UpdateTranscriptStatus(ctx, sqlc.UpdateTranscriptStatusParams{
					ClassLearnerID: learner.ClassLearnerID,
					TranscriptID:   transcript.TranscriptID,
					Status:         1,
				}); err != nil {
					return err
				}
				totalGrade = float64(*transcript.Grade) * transcript.Weight
			}
		}

		attendaces, err := c.queries.CountAttendace(ctx, learner.ClassLearnerID)
		if err != nil {
			return err
		}

		slots, err := c.queries.CountSlotsByClassId(ctx, classId)
		if err != nil {
			return err
		}

		if totalGrade < float64(*subject.MinPassGrade) || math.Ceil(float64(attendaces)/float64(slots)*100) < float64(*subject.MinAttendance) {
			if err = c.queries.UpdateClassStatus(ctx, sqlc.UpdateClassStatusParams{
				ID:     learner.ClassLearnerID,
				Status: 0,
			}); err != nil {
				return err
			}
		} else {
			if err = c.queries.UpdateClassStatus(ctx, sqlc.UpdateClassStatusParams{
				ID:     learner.ClassLearnerID,
				Status: 1,
			}); err != nil {
				return err
			}
		}

	}
	return nil
}
