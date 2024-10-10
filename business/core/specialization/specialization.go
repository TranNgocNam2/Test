package specialization

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	ErrSkillNotFound   = errors.New("Kỹ năng không có trong hệ thống!")
	ErrSubjectNotFound = errors.New("Môn học không có trong hệ thống!")
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

func (c *Core) Create(ctx *gin.Context, newSpecialization Specialization) error {
	userID, err := middleware.AuthorizeUser(ctx, c.queries)
	if err != nil {
		return err
	}

	var dbSpecialization = sqlc.CreateSpecializationParams{
		ID:          newSpecialization.ID,
		Name:        newSpecialization.Name,
		Code:        newSpecialization.Code,
		TimeAmount:  newSpecialization.TimeAmount,
		ImageLink:   newSpecialization.Image,
		Description: newSpecialization.Description,
		CreatedBy:   userID,
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	err = qtx.CreateSpecialization(ctx, dbSpecialization)
	if err != nil {
		return err
	}

	if newSpecialization.Skills != nil {
		var skillIDs []uuid.UUID
		for _, skill := range newSpecialization.Skills {
			skillIDs = append(skillIDs, *skill.ID)
		}
		dbSkill, err := qtx.GetSkillsByIDs(ctx, skillIDs)
		if err != nil || (len(dbSkill) != len(skillIDs)) {
			return ErrSkillNotFound
		}

		specSkills := sqlc.CreateSpecializationSkillsParams{
			SpecializationID: dbSpecialization.ID,
			SkillIds:         skillIDs,
		}
		err = qtx.CreateSpecializationSkills(ctx, specSkills)
		if err != nil {
			return err
		}
	}

	if newSpecialization.Subjects != nil {
		var subjectIDs []uuid.UUID
		for _, subject := range newSpecialization.Subjects {
			subjectIDs = append(subjectIDs, *subject.ID)
		}
		dbSubject, err := c.queries.GetSubjectsByIDs(ctx, subjectIDs)
		if err != nil || (len(dbSubject) != len(subjectIDs)) {
			return ErrSubjectNotFound
		}

		specSubjects := sqlc.CreateSpecializationSubjectsParams{
			SpecializationID: dbSpecialization.ID,
			SubjectIds:       subjectIDs,
			CreatedBy:        userID,
		}

		err = qtx.CreateSpecializationSubjects(ctx, specSubjects)
		if err != nil {
			return err
		}
	}
	tx.Commit(ctx)
	return nil
}
