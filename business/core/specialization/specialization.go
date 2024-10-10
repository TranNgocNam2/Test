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
	"time"
)

var (
	ErrSkillNotFound        = errors.New("Kỹ năng không có trong hệ thống!")
	ErrSubjectNotFound      = errors.New("Môn học không có trong hệ thống!")
	ErrSpecCodeAlreadyExist = errors.New("Mã chuyên ngành đã tồn tại!")
	ErrSpecNotFound         = errors.New("Chuyên ngành không tồn tại!")
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
	staffID, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	if _, err = c.queries.GetSpecializationByCode(ctx, newSpecialization.Code); err == nil {
		return ErrSpecCodeAlreadyExist
	}

	var dbSpecialization = sqlc.CreateSpecializationParams{
		ID:          newSpecialization.ID,
		Name:        newSpecialization.Name,
		Code:        newSpecialization.Code,
		TimeAmount:  newSpecialization.TimeAmount,
		ImageLink:   newSpecialization.Image,
		Description: newSpecialization.Description,
		CreatedBy:   staffID,
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	if err = qtx.CreateSpecialization(ctx, dbSpecialization); err != nil {
		return err
	}

	if err = processSpecSkills(ctx, qtx, newSpecialization); err != nil {
		return err
	}

	if err = processSpecSubjects(ctx, qtx, newSpecialization, staffID); err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}

func (c *Core) GetByID(ctx *gin.Context, id uuid.UUID) (Specialization, error) {
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return Specialization{}, err
	}

	dbSpec, err := c.queries.GetSpecializationByID(ctx, id)
	if err != nil {
		return Specialization{}, ErrSpecNotFound
	}

	spec := toCoreSpecialization(dbSpec)
	dbSpecSkills, err := c.queries.GetSkillsBySpecialization(ctx, dbSpec.ID)
	if err != nil {
		return Specialization{}, ErrSkillNotFound
	}
	if dbSpecSkills != nil {
		for _, skill := range dbSpecSkills {
			spec.Skills = append(spec.Skills, &struct {
				ID   *uuid.UUID
				Name *string
			}{
				ID:   &skill.ID,
				Name: &skill.Name,
			})
		}
	}

	dbSpecSubjects, err := c.queries.GetSubjectsBySpecialization(ctx, dbSpec.ID)
	if err != nil {
		return Specialization{}, ErrSubjectNotFound
	}
	if dbSpecSubjects != nil {
		for _, subject := range dbSpecSubjects {
			spec.Subjects = append(spec.Subjects, &struct {
				ID          *uuid.UUID
				Name        *string
				Image       *string
				Code        *string
				LastUpdated time.Time
			}{
				ID:          &subject.ID,
				Name:        &subject.Name,
				Image:       &subject.ImageLink,
				Code:        &subject.Code,
				LastUpdated: subject.CreatedAt,
			})
		}
	}

	return spec, nil
}
func processSpecSkills(ctx *gin.Context, qtx *sqlc.Queries, specialization Specialization) error {
	if specialization.Skills != nil {
		var skillIDs []uuid.UUID
		for _, skill := range specialization.Skills {
			skillIDs = append(skillIDs, *skill.ID)
		}
		dbSkill, err := qtx.GetSkillsByIDs(ctx, skillIDs)
		if err != nil || (len(dbSkill) != len(skillIDs)) {
			return ErrSkillNotFound
		}

		specSkills := sqlc.CreateSpecializationSkillsParams{
			SpecializationID: specialization.ID,
			SkillIds:         skillIDs,
		}
		err = qtx.CreateSpecializationSkills(ctx, specSkills)
		if err != nil {
			return err
		}
	}
	return nil
}

func processSpecSubjects(ctx *gin.Context, qtx *sqlc.Queries, specialization Specialization, staffID string) error {
	if specialization.Subjects != nil {
		var subjectIDs []uuid.UUID
		for _, subject := range specialization.Subjects {
			subjectIDs = append(subjectIDs, *subject.ID)
		}
		dbSubject, err := qtx.GetSubjectsByIDs(ctx, subjectIDs)
		if err != nil || (len(dbSubject) != len(subjectIDs)) {
			return ErrSubjectNotFound
		}

		specSubjects := sqlc.CreateSpecializationSubjectsParams{
			SpecializationID: specialization.ID,
			SubjectIds:       subjectIDs,
			CreatedBy:        staffID,
		}
		err = qtx.CreateSpecializationSubjects(ctx, specSubjects)
		if err != nil {
			return err
		}
	}
	return nil
}
