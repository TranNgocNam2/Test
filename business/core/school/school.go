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

func (c *Core) Create(ctx *gin.Context, newSchool NewSchool) (uuid.UUID, error) {
	var dbSchool = sqlc.CreateSchoolParams{
		ID:         uuid.New(),
		Name:       newSchool.Name,
		Address:    newSchool.Address,
		DistrictID: newSchool.DistrictId,
	}

	if err := c.queries.CreateSchool(ctx, dbSchool); err != nil {
		return uuid.Nil, err
	}
	return dbSchool.ID, nil
}

func (c *Core) Update(ctx *gin.Context, id uuid.UUID, updatedSchool UpdateSchool) error {

	dbSchool, err := c.queries.GetSchoolById(ctx, id)
	if err != nil {
		return model.ErrSchoolNotFound
	}

	updateSchool := sqlc.UpdateSchoolParams{
		Name:       updatedSchool.Name,
		Address:    updatedSchool.Address,
		DistrictID: updatedSchool.DistrictId,
		ID:         dbSchool.ID,
	}

	if err = c.queries.UpdateSchool(ctx, updateSchool); err != nil {
		return err
	}
	return nil
}

func (c *Core) Delete(ctx *gin.Context, id uuid.UUID) error {
	dbSchool, err := c.queries.GetSchoolById(ctx, id)
	if err != nil {
		return model.ErrSchoolNotFound
	}

	if err = c.queries.DeleteSchool(ctx, dbSchool.ID); err != nil {
		return err
	}
	return nil
}

func (c *Core) GetByID(ctx *gin.Context, id uuid.UUID) (School, error) {
	dbSchool, err := c.queries.GetAllSchoolInformationById(ctx, id)
	if err != nil {
		return School{}, model.ErrSchoolNotFound
	}

	school := School{
		ID:      dbSchool.ID,
		Name:    dbSchool.Name,
		Address: dbSchool.Address,
		District: District{
			ID:   dbSchool.DistrictID,
			Name: dbSchool.DistrictName,
		},
		Province: Province{
			ID:   dbSchool.ProvinceID,
			Name: dbSchool.ProvinceName,
		},
	}

	return school, nil
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
				s.id, s.name, s.address, d.id AS district_id, d.name AS district_name, p.id AS province_id, p.name AS province_name 
					FROM schools s 
						JOIN districts d ON s.district_id = d.id
						JOIN provinces p ON p.id = d.province_id`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)
	orderByClause := orderByClause(orderBy)

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var dbSchools []sqlc.GetAllSchoolsRow

	if err := pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbSchools); err != nil {
		c.logger.Error(err.Error())
		return nil
	}
	if len(dbSchools) == 0 || dbSchools == nil {
		return nil
	}

	var result []School
	for _, dbSchool := range dbSchools {
		school := School{
			ID:      dbSchool.ID,
			Name:    dbSchool.Name,
			Address: dbSchool.Address,
			District: District{
				ID:   dbSchool.DistrictID,
				Name: dbSchool.DistrictName,
			},
			Province: Province{
				ID:   dbSchool.ProvinceID,
				Name: dbSchool.ProvinceName,
			},
		}

		result = append(result, school)
	}

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
               FROM schools s 
						JOIN districts d ON s.district_id = d.id
						JOIN provinces p ON p.id = d.province_id`

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

func (c *Core) GetSchoolsByDistrictId(ctx *gin.Context, id int) ([]School, error) {
	district, err := c.queries.GetDistrictById(ctx, int32(id))
	if err != nil {
		return nil, model.ErrDistrictNotFound
	}

	province, _ := c.queries.GetProvinceById(ctx, district.ProvinceID)

	dbSchools, err := c.queries.GetSchoolsByDistrictId(ctx, district.ID)
	if err != nil {
		return nil, model.ErrSchoolNotFound
	}

	var schools []School
	for _, dbSchool := range dbSchools {
		school := toCoreSchool(dbSchool)
		school.District = toCoreDistrict(district)
		school.Province = toCoreProvince(province)
		schools = append(schools, school)
	}

	return schools, nil
}

func (c *Core) GetAllProvinces(ctx *gin.Context) ([]Province, error) {
	provinces, err := c.queries.GetAllProvince(ctx)
	if err != nil {
		return nil, err
	}
	return toCoreProvinceSlice(provinces), nil
}

func (c *Core) GetDistrictsByProvinceId(ctx *gin.Context, id int) ([]District, error) {
	province, err := c.queries.GetProvinceById(ctx, int32(id))
	if err != nil {
		return nil, model.ErrProvinceNotFound
	}
	districts, _ := c.queries.GetDistrictsByProvince(ctx, province.ID)

	return toCoreDistrictSlice(districts), nil
}
