package program

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

func (c *Core) Create(ctx *gin.Context, newProgram NewProgram) (uuid.UUID, error) {
	staffID, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return uuid.Nil, err
	}

	dbProgram := sqlc.CreateProgramParams{
		ID:          newProgram.ID,
		Name:        newProgram.Name,
		StartDate:   newProgram.StartDate,
		EndDate:     newProgram.EndDate,
		CreatedBy:   staffID,
		Description: newProgram.Description,
	}
	if err = c.queries.CreateProgram(ctx, dbProgram); err != nil {
		return uuid.Nil, err
	}
	return dbProgram.ID, nil
}

func (c *Core) Update(ctx *gin.Context, id uuid.UUID, updateProgram UpdateProgram) error {
	staffID, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbProgram, err := c.queries.GetProgramByID(ctx, id)
	if err != nil {
		return model.ErrProgramNotFound
	}

	if dbProgram.StartDate.Before(time.Now()) {
		return model.ErrCannotUpdateProgram
	}

	dbUpdateProgram := sqlc.UpdateProgramParams{
		Name:        updateProgram.Name,
		StartDate:   updateProgram.StartDate,
		EndDate:     updateProgram.EndDate,
		UpdatedBy:   &staffID,
		Description: updateProgram.Description,
		ID:          id,
	}

	if err = c.queries.UpdateProgram(ctx, dbUpdateProgram); err != nil {
		return err
	}

	return nil
}

func (c *Core) Query(ctx *gin.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) []Program {
	if err := filter.Validate(); err != nil {
		return nil
	}

	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT
						id, name, start_date, end_date
			FROM programs`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var dbPrograms []sqlc.Program
	err := pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbPrograms)
	if err != nil {
		c.logger.Error(err.Error())
		return nil
	}

	if dbPrograms == nil {
		return nil
	}

	var programs []Program
	for _, dbProgram := range dbPrograms {
		program := toCoreProgram(dbProgram)
		program.TotalClasses, err = c.queries.CountClassesByProgramId(ctx, dbProgram.ID)
		if err != nil {
			c.logger.Error(err.Error())
			return nil
		}
		programs = append(programs, program)
	}
	return programs
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
                        programs`

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

func (c *Core) Delete(ctx *gin.Context, id uuid.UUID) error {
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbProgram, err := c.queries.GetProgramByID(ctx, id)
	if err != nil {
		return model.ErrProgramNotFound
	}

	totalClasses, _ := c.queries.CountClassesByProgramId(ctx, id)
	if totalClasses > 0 || dbProgram.StartDate.Before(time.Now()) {
		return model.ErrCannotDeleteProgram
	}

	if err = c.queries.DeleteProgram(ctx, id); err != nil {
		return err
	}

	return nil
}
