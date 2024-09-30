package schooldb

import (
	"Backend/business/core/school"
	"Backend/business/db/pgx"
	_ "Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/order"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	_ "go.uber.org/zap"
)

type Store struct {
	db      *sqlx.DB
	queries *sqlc.Queries
	logger  *zap.Logger
}

func NewStore(db *sqlx.DB, queries *sqlc.Queries, logger *zap.Logger) *Store {
	return &Store{
		db:      db,
		queries: queries,
		logger:  logger,
	}
}

func (s *Store) Create(ctx *gin.Context, school school.School) error {
	newSchoolDB := sqlc.CreateSchoolParams{
		ID:         uuid.New(),
		Name:       school.Name,
		Address:    school.Address,
		DistrictID: school.DistrictID,
	}
	if err := s.queries.CreateSchool(ctx, newSchoolDB); err != nil {
		return err
	}

	return nil
}

func (s *Store) Update(ctx *gin.Context, school school.School) error {
	updateSchoolDB := sqlc.UpdateSchoolParams{
		Name:       school.Name,
		Address:    school.Address,
		DistrictID: school.DistrictID,
		ID:         school.ID,
	}
	if err := s.queries.UpdateSchool(ctx, updateSchoolDB); err != nil {
		return err
	}

	return nil
}

func (s *Store) Delete(ctx *gin.Context, school school.School) error {
	if err := s.queries.DeleteSchool(ctx, school.ID); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetByID(ctx *gin.Context, id uuid.UUID) (school.School, error) {
	schoolDB, err := s.queries.GetSchoolByID(ctx, id)
	if err != nil {
		return school.School{}, err
	}

	return toCoreSchool(schoolDB), nil
}

func (s *Store) Query(ctx *gin.Context, filter school.QueryFilter, orderBy order.By, pageNumber, rowsPerPage int) ([]school.School, error) {
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

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	s.logger.Info(buf.String())

	var schoolsDB []sqlc.School

	if err := pgx.NamedQuerySlice(ctx, s.logger, s.db, buf.String(), data, &schoolsDB); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toCoreSchoolSlice(schoolsDB), nil
}

func (s *Store) Count(ctx *gin.Context, filter school.QueryFilter) (int, error) {
	data := map[string]interface{}{}

	const q = `SELECT
                      count (1)
               FROM
                      schools`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}

	if err := pgx.NamedQueryStruct(ctx, s.logger, s.db, buf.String(), data, &count); err != nil {
		return 0, err
	}

	return count.Count, nil
}

func (s *Store) GetByDistrict(ctx *gin.Context, districtID int32) ([]school.School, error) {
	schools, err := s.queries.GetSchoolsByDistrictID(ctx, districtID)
	if err != nil {
		return nil, err
	}

	return toCoreSchoolSlice(schools), nil
}

func (s *Store) GetAllProvinces(ctx *gin.Context) ([]school.Province, error) {
	provinces, err := s.queries.GetAllProvince(ctx)
	if err != nil {
		return nil, err
	}

	return toCoreProvinceSlice(provinces), nil
}

func (s *Store) GetDistrictsByProvince(ctx *gin.Context, provinceID int32) ([]school.District, error) {
	districts, err := s.queries.GetDistrictsByProvince(ctx, provinceID)
	if err != nil {
		return nil, err
	}

	return toCoreDistrictSlice(districts), nil
}
