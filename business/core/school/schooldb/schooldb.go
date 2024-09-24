package schooldb

import (
	"Backend/business/core/school"
	"Backend/business/db/sqlc"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrSchoolNotFound   = errors.New("school not found")
	ErrDistrictNotFound = errors.New("district not found")
	ErrProvinceNotFound = errors.New("province/city not found")
	FailedToDelete      = errors.New("failed to delete school")
	FailedToCreate      = errors.New("failed to create school")
	FailedToUpdate      = errors.New("failed to update school")
)

type Store struct {
	db      *sqlx.DB
	queries *sqlc.Queries
}

func NewStore(db *sqlx.DB, queries *sqlc.Queries) *Store {
	return &Store{
		db:      db,
		queries: queries,
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
		return FailedToCreate
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
		return FailedToUpdate
	}

	return nil
}

func (s *Store) Delete(ctx *gin.Context, school school.School) error {
	if err := s.queries.DeleteSchool(ctx, school.ID); err != nil {
		return FailedToDelete
	}

	return nil
}

func (s *Store) GetByID(ctx *gin.Context, id uuid.UUID) (school.School, error) {
	schoolDB, err := s.queries.GetSchoolByID(ctx, id)
	if err != nil {
		return school.School{}, ErrSchoolNotFound
	}

	return toCoreSchool(schoolDB), nil
}

func (s *Store) GetByDistrict(ctx *gin.Context, districtID int32) ([]school.School, error) {
	schools, err := s.queries.GetSchoolsByDistrictID(ctx, districtID)
	if err != nil {
		return nil, ErrSchoolNotFound
	}

	return toCoreSchoolSlice(schools), nil
}

func (s *Store) GetAllProvinces(ctx *gin.Context) ([]school.Province, error) {
	provinces, err := s.queries.GetAllProvince(ctx)
	if err != nil {
		return nil, ErrProvinceNotFound
	}

	return toCoreProvinceSlice(provinces), nil
}

func (s *Store) GetDistrictsByProvince(ctx *gin.Context, provinceID int32) ([]school.District, error) {
	districts, err := s.queries.GetDistrictsByProvince(ctx, provinceID)
	if err != nil {
		return nil, ErrDistrictNotFound
	}

	return toCoreDistrictSlice(districts), nil
}
