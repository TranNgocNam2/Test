package specialization

import (
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"bytes"
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

func (c *Core) Create(ctx *gin.Context, newSpec NewSpecialization) (uuid.UUID, error) {
	staffID, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return uuid.Nil, err
	}

	if _, err = c.queries.GetSpecializationByCode(ctx, newSpec.Code); err == nil {
		return uuid.Nil, model.ErrSpecCodeAlreadyExist
	}

	var dbSpec = sqlc.CreateSpecializationParams{
		ID:          newSpec.ID,
		Name:        newSpec.Name,
		Code:        newSpec.Code,
		ImageLink:   newSpec.Image,
		Description: newSpec.Description,
		CreatedBy:   staffID,
	}

	if err = c.queries.CreateSpecialization(ctx, dbSpec); err != nil {
		return uuid.Nil, err
	}

	return dbSpec.ID, nil
}

func (c *Core) GetById(ctx *gin.Context, id uuid.UUID) (Details, error) {
	dbSpec, err := c.queries.GetSpecializationById(ctx, id)
	if err != nil {
		return Details{}, model.ErrSpecNotFound
	}

	//if dbSpec.Status == Draft || dbSpec.Status == Deleted {
	//	if _, err = middleware.AuthorizeStaff(ctx, c.queries); err != nil {
	//		return Details{}, err
	//	}
	//}

	spec := toCoreSpecializationDetails(dbSpec)

	dbSubjects, err := c.queries.GetSubjectsBySpecialization(ctx, dbSpec.ID)
	if err != nil {
		return Details{}, model.ErrSubjectNotFound
	}
	if dbSubjects != nil {
		for _, dbSubject := range dbSubjects {
			totalSessions, err := c.queries.CountSessionsBySubjectId(ctx, dbSubject.ID)
			if err != nil {
				return Details{}, err
			}

			subject := Subject{
				ID:            dbSubject.ID,
				Name:          dbSubject.Name,
				Image:         *dbSubject.ImageLink,
				Code:          dbSubject.Code,
				LastUpdated:   dbSubject.UpdatedAt,
				Skills:        nil,
				Index:         dbSubject.Index,
				MinPassGrade:  *dbSubject.MinPassGrade,
				MinAttendance: *dbSubject.MinAttendance,
				CreatedBy:     dbSubject.CreatedBy,
				CreatedAt:     dbSubject.CreatedAt,
				UpdatedBy:     dbSubject.UpdatedBy,
				TotalSessions: 0,
			}

			subject.TotalSessions = totalSessions
			dbSkills, err := c.queries.GetSkillsBySubjectId(ctx, dbSubject.ID)
			if err != nil {
				return Details{}, model.ErrSkillNotFound
			}

			for _, dbSkill := range dbSkills {
				skill := toCoreSkill(dbSkill)
				subject.Skills = append(subject.Skills, skill)
			}
			spec.Subjects = append(spec.Subjects, subject)
		}
	}

	return spec, nil
}

func (c *Core) Update(ctx *gin.Context, id uuid.UUID, updateSpec UpdateSpecialization) error {
	staffID, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbSpec, err := c.queries.GetSpecializationById(ctx, id)
	if err != nil {
		return model.ErrSpecNotFound
	}

	if dbSpec.Code != updateSpec.Code {
		_, err = c.queries.GetSpecializationByCode(ctx, updateSpec.Code)
		if err == nil {
			return model.ErrSpecCodeAlreadyExist
		}
	}

	var dbUpdateSpecialization sqlc.UpdateSpecializationParams

	if dbSpec.Status == Published {
		if updateSpec.Status == 0 {
			return model.ErrSpecStatusCannotBeDraft
		}
		dbUpdateSpecialization = sqlc.UpdateSpecializationParams{
			ID:          id,
			ImageLink:   &updateSpec.Image,
			Description: &updateSpec.Description,
			UpdatedBy:   &staffID,
			Name:        dbSpec.Name,
			Code:        dbSpec.Code,
			Status:      Published,
		}
	}

	if dbSpec.Status == Draft {
		dbUpdateSpecialization = sqlc.UpdateSpecializationParams{
			ID:     id,
			Name:   updateSpec.Name,
			Code:   updateSpec.Code,
			Status: updateSpec.Status,
			ImageLink:   &updateSpec.Image,
			Description: &updateSpec.Description,
			UpdatedBy:   &staffID,
		}
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	timeAmount, err := c.processSpecSubjects(ctx, qtx, dbSpec.ID, updateSpec.Subjects, staffID)
	if err != nil {
		return err
	}

	dbUpdateSpecialization.TimeAmount = timeAmount

	if err = qtx.UpdateSpecialization(ctx, dbUpdateSpecialization); err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil

}

func (c *Core) Query(ctx *gin.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) []Specialization {
	if err := filter.Validate(); err != nil {
		return nil
	}

	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT
						id, name, code, time_amount, image_link, description, status
			FROM specializations`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var dbSpecializations []sqlc.Specialization
	err := pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbSpecializations)
	if err != nil {
		c.logger.Error(err.Error())
		return nil
	}

	if dbSpecializations == nil {
		return nil
	}

	var specializations []Specialization

	for _, dbSpec := range dbSpecializations {
		spec := toCoreSpecialization(dbSpec)
		totalSubjects, err := c.queries.CountSubjectsBySpecializationId(ctx, dbSpec.ID)
		if err != nil {
			c.logger.Error(err.Error())
			return nil
		}

		spec.TotalSubjects = totalSubjects
		specializations = append(specializations, spec)
	}

	return specializations
}

func (c *Core) Delete(ctx *gin.Context, id uuid.UUID) error {
	staffID, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbSpec, err := c.queries.GetSpecializationById(ctx, id)
	if err != nil {
		return model.ErrSpecNotFound
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	if err = qtx.DeleteSpecializationSubjects(ctx, dbSpec.ID); err != nil {
		return err
	}

	if dbSpec.Status == Draft {
		if err = qtx.DeleteSpecialization(ctx, dbSpec.ID); err != nil {
			return err
		}
	}

	if dbSpec.Status == Published {
		if err = qtx.UpdateSpecializationStatus(ctx, sqlc.UpdateSpecializationStatusParams{
			UpdatedBy: &staffID,
			ID:        dbSpec.ID,
		}); err != nil {
			return err
		}
	}

	tx.Commit(ctx)
	return nil
}

func (c *Core) Count(ctx *gin.Context, filter QueryFilter) int {
	if err := filter.Validate(); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	data := map[string]interface{}{}

	const q = `SELECT
                        count(1)
               FROM
                        specializations`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}

	if err := pgx.NamedQueryStruct(ctx, c.logger, c.db, buf.String(), data, &count); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	return count.Count
}

func (c *Core) processSpecSubjects(ctx *gin.Context, qtx *sqlc.Queries, specializationId uuid.UUID, specSubjects []SpecSubject, staffID string) (*float32, error) {
	err := qtx.DeleteSpecializationSubjects(ctx, specializationId)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, err
	}
	var timeAmount float32
	for _, specSubject := range specSubjects {
		dbSubject, err := qtx.GetSubjectById(ctx, specSubject.ID)
		if err != nil {
			return nil, model.ErrSubjectNotFound
		}
		timeAmount += dbSubject.TimePerSession * float32(dbSubject.TotalSessions)

		dbSpecSubject := sqlc.CreateSpecializationSubjectParams{
			SpecializationID: specializationId,
			SubjectID:        specSubject.ID,
			Index:            specSubject.Index,
			CreatedBy:        staffID,
		}
		err = qtx.CreateSpecializationSubject(ctx, dbSpecSubject)
		if err != nil {
			c.logger.Error(err.Error())
			return nil, err
		}
	}

	return &timeAmount, nil
}
