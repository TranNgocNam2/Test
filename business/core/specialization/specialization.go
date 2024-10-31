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
		TimeAmount:  newSpec.TimeAmount,
		ImageLink:   newSpec.Image,
		Description: newSpec.Description,
		CreatedBy:   staffID,
	}

	if err = c.queries.CreateSpecialization(ctx, dbSpec); err != nil {
		return uuid.Nil, err
	}

	return dbSpec.ID, nil
}

func (c *Core) GetByID(ctx *gin.Context, id uuid.UUID) (Details, error) {
	dbSpec, err := c.queries.GetSpecializationByID(ctx, id)
	if err != nil {
		return Details{}, model.ErrSpecNotFound
	}

	if dbSpec.Status == Draft || dbSpec.Status == Delete {
		if _, err = middleware.AuthorizeStaff(ctx, c.queries); err != nil {
			return Details{}, err
		}
	}

	spec := toCoreSpecializationDetails(dbSpec)

	dbSpecSubjects, err := c.queries.GetSubjectsBySpecialization(ctx, dbSpec.ID)
	if err != nil {
		return Details{}, model.ErrSubjectNotFound
	}
	if dbSpecSubjects != nil {
		for _, subject := range dbSpecSubjects {
			totalSessions, err := c.queries.CountSessionsBySubjectID(ctx, subject.ID)
			if err != nil {
				return Details{}, err
			}
			spec.Subjects = append(spec.Subjects, &struct {
				ID            uuid.UUID `json:"id"`
				Name          string    `json:"name"`
				Image         string    `json:"image"`
				Code          string    `json:"code"`
				LastUpdated   time.Time `json:"lastUpdated"`
				TotalSessions int64     `json:"totalSessions"`
			}{
				ID:            subject.ID,
				Name:          subject.Name,
				Image:         *subject.ImageLink,
				Code:          subject.Code,
				LastUpdated:   subject.CreatedAt,
				TotalSessions: totalSessions,
			})
		}
	}

	return spec, nil
}

func (c *Core) Update(ctx *gin.Context, id uuid.UUID, updateSpec UpdateSpecialization) error {
	staffID, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbSpec, err := c.queries.GetSpecializationByID(ctx, id)
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

	if dbSpec.Status == Public {
		dbUpdateSpecialization = sqlc.UpdateSpecializationParams{
			ID:          id,
			TimeAmount:  &updateSpec.TimeAmount,
			ImageLink:   &updateSpec.Image,
			Description: &updateSpec.Description,
			UpdatedBy:   &staffID,
			Name:        dbSpec.Name,
			Code:        dbSpec.Code,
			Status:      1,
		}
	}

	if dbSpec.Status == Draft {
		dbUpdateSpecialization = sqlc.UpdateSpecializationParams{
			ID:          id,
			Name:        updateSpec.Name,
			Code:        updateSpec.Code,
			Status:      updateSpec.Status,
			TimeAmount:  &updateSpec.TimeAmount,
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

	if err = qtx.UpdateSpecialization(ctx, dbUpdateSpecialization); err != nil {
		return err
	}

	if err = processSpecSubjects(ctx, qtx, dbSpec.ID, updateSpec.Subjects, staffID); err != nil {
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
						id, name, code, time_amount, image_link
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
		totalSubjects, err := c.queries.CountSubjectsBySpecializationID(ctx, dbSpec.ID)
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

	dbSpec, err := c.queries.GetSpecializationByID(ctx, id)
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

	if dbSpec.Status == Public {
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

func processSpecSubjects(ctx *gin.Context, qtx *sqlc.Queries, specializationID uuid.UUID, subjectIDs []uuid.UUID, staffID string) error {
	if subjectIDs != nil {
		err := qtx.DeleteSpecializationSubjects(ctx, specializationID)
		if err != nil {
			return err
		}

		dbSubject, err := qtx.GetSubjectsByIDs(ctx, subjectIDs)
		if err != nil || (len(dbSubject) != len(subjectIDs)) {
			return model.ErrSubjectNotFound
		}

		specSubjects := sqlc.CreateSpecializationSubjectsParams{
			SpecializationID: specializationID,
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
