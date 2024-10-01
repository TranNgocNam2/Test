package school

import (
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/order"
	"bytes"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gitlab.com/innovia69420/kit/web/request"
	"go.uber.org/zap"
)

var (
	ErrInvalidID          = errors.New("ID không hợp lệ!")
	ErrSchoolNotFound     = errors.New("Không tìm thấy trường học!")
	ErrCreateSchoolFailed = errors.New("Có lỗi trong quá trình tạo trường học!")
	ErrUpdateSchoolFailed = errors.New("Có lỗi trong quá trình cập nhật trường học!")
	ErrDeleteSchoolFailed = errors.New("Có lỗi trong quá trình xóa trường học!")
)

type Core struct {
	db      *sqlx.DB
	queries *sqlc.Queries
	logger  *zap.Logger
}

func NewCore(db *sqlx.DB, queries *sqlc.Queries, logger *zap.Logger) *Core {
	return &Core{
		db:      db,
		queries: queries,
		logger:  logger,
	}
}

func (c *Core) Create(ctx *gin.Context, request request.NewSchool) error {
	var school = sqlc.CreateSchoolParams{
		ID:         uuid.New(),
		Name:       request.Name,
		Address:    request.Address,
		DistrictID: request.DistrictId,
	}

	if err := c.queries.CreateSchool(ctx, school); err != nil {
		return ErrCreateSchoolFailed
	}
	return nil
}

func (c *Core) Update(ctx *gin.Context, request request.UpdateSchool) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ErrInvalidID
	}

	school, err := c.queries.GetSchoolByID(ctx, id)
	if err != nil {
		return ErrSchoolNotFound
	}

	school.Name = request.Name
	school.Address = request.Address
	school.DistrictID = request.DistrictId

	params := sqlc.UpdateSchoolParams{
		Name:       school.Name,
		Address:    school.Address,
		DistrictID: school.DistrictID,
		ID:         school.ID,
	}

	if err = c.queries.UpdateSchool(ctx, params); err != nil {
		return ErrUpdateSchoolFailed
	}
	return nil
}

func (c *Core) Delete(ctx *gin.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ErrInvalidID
	}

	school, err := c.queries.GetSchoolByID(ctx, id)
	if err != nil {
		return ErrSchoolNotFound
	}

	if err = c.queries.DeleteSchool(ctx, school.ID); err != nil {
		return ErrDeleteSchoolFailed
	}
	return nil
}

func (c *Core) GetSchoolByID(ctx *gin.Context, id uuid.UUID) (*School, error) {
	school, err := c.queries.GetSchoolByID(ctx, id)

	if err != nil {
		return nil, ErrSchoolNotFound
	}

	result := School{
		ID:         school.ID,
		Name:       school.Name,
		Address:    school.Address,
		DistrictID: school.DistrictID,
	}

	return &result, nil
}

func (c *Core) GetSchoolsPaginated(ctx *gin.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) []School {
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

func (c *Core) GetSchoolsByDistrictID(ctx *gin.Context) ([]School, error) {
	districtID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, ErrInvalidID
	}

	schools, err := c.queries.GetSchoolsByDistrictID(ctx, int32(districtID))
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

func (c *Core) GetDistrictsByProvinceID(ctx *gin.Context) ([]District, error) {
	provinceID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, ErrInvalidID
	}

	districts, err := c.queries.GetDistrictsByProvince(ctx, int32(provinceID))
	if err != nil {
		return nil, err
	}
	return toCoreDistrictSlice(districts), nil
}
