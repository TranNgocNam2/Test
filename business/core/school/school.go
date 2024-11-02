package school

import (
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/order"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		db:      app.DB,
		queries: app.Queries,
		logger:  app.Logger,
	}
}

func (c *Core) Create(ctx *gin.Context, newSchool School) (uuid.UUID, error) {
	var dbSchool = sqlc.CreateSchoolParams{
		ID:         uuid.New(),
		Name:       newSchool.Name,
		Address:    newSchool.Address,
		DistrictID: newSchool.DistrictID,
	}

	if err := c.queries.CreateSchool(ctx, dbSchool); err != nil {
		return uuid.Nil, err
	}
	return dbSchool.ID, nil
}

func (c *Core) Update(ctx *gin.Context, id uuid.UUID, updatedSchool School) error {

	dbSchool, err := c.queries.GetSchoolByID(ctx, id)
	if err != nil {
		return model.ErrSchoolNotFound
	}

	updateSchool := sqlc.UpdateSchoolParams{
		Name:       updatedSchool.Name,
		Address:    updatedSchool.Address,
		DistrictID: updatedSchool.DistrictID,
		ID:         dbSchool.ID,
	}

	if err = c.queries.UpdateSchool(ctx, updateSchool); err != nil {
		return err
	}
	return nil
}

func (c *Core) Delete(ctx *gin.Context, id uuid.UUID) error {
	dbSchool, err := c.queries.GetSchoolByID(ctx, id)
	if err != nil {
		return model.ErrSchoolNotFound
	}

	if err = c.queries.DeleteSchool(ctx, dbSchool.ID); err != nil {
		return err
	}
	return nil
}

func (c *Core) GetByID(ctx *gin.Context, id uuid.UUID) (School, error) {
	school, err := c.queries.GetSchoolByID(ctx, id)

	if err != nil {
		return School{}, model.ErrSchoolNotFound
	}

	return toCoreSchool(school), nil
}

func (c *Core) Query(ctx *gin.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) []School {
	if err := filter.Validate(); err != nil {
		c.logger.Error(err.Error())
		return nil
	}

	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT
                        id, name, address, district_id
               FROM
                        schools`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)
	orderByClause := orderByClause(orderBy)

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var schools []sqlc.School

	if err := pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &schools); err != nil {
		c.logger.Error(err.Error())
		return nil
	}

	result := toCoreSchoolSlice(schools)

	return result
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
                        schools`

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

func (c *Core) GetSchoolsByDistrictID(ctx *gin.Context, id int) ([]School, error) {
	schools, err := c.queries.GetSchoolsByDistrictID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return toCoreSchoolSlice(schools), nil
}

func (c *Core) GetAllProvinces(ctx *gin.Context) ([]Province, error) {
	provinces, err := c.queries.GetAllProvince(ctx)
	if err != nil {
		return nil, err
	}
	return toCoreProvinceSlice(provinces), nil
}

func (c *Core) GetDistrictsByProvinceID(ctx *gin.Context, id int) ([]District, error) {
	districts, err := c.queries.GetDistrictsByProvince(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return toCoreDistrictSlice(districts), nil
}
