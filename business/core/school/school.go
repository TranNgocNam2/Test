package school

import (
	"Backend/internal/order"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

var (
	ErrInvalidID = errors.New("ID không hợp lệ!")
)

type Storer interface {
	Create(ctx *gin.Context, school School) error
	Update(ctx *gin.Context, school School) error
	Delete(ctx *gin.Context, school School) error
	GetByID(ctx *gin.Context, id uuid.UUID) (School, error)
	Query(ctx *gin.Context, filter QueryFilter, orderBy order.By, pageNumber, rowsPerPage int) ([]School, error)
	Count(ctx *gin.Context, filter QueryFilter) (int, error)
	GetByDistrict(ctx *gin.Context, districtID int32) ([]School, error)
	GetAllProvinces(ctx *gin.Context) ([]Province, error)
	GetDistrictsByProvince(ctx *gin.Context, provinceID int32) ([]District, error)
}

type Core struct {
	storer Storer
}

func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

func (c *Core) Create(ctx *gin.Context, newSchool NewSchool) (error, int) {
	var school = School{
		ID:         uuid.New(),
		Name:       newSchool.Name,
		Address:    newSchool.Address,
		DistrictID: newSchool.DistrictID,
	}
	if err := c.storer.Create(ctx, school); err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (c *Core) Update(ctx *gin.Context, updateSchool UpdateSchool) (error, int) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ErrInvalidID, http.StatusBadRequest
	}

	school, err := c.storer.GetByID(ctx, id)
	if err != nil {
		return err, http.StatusNotFound
	}

	if updateSchool.Name != nil {
		school.Name = *updateSchool.Name
	}

	if updateSchool.Address != nil {
		school.Address = *updateSchool.Address
	}

	if updateSchool.DistrictID != nil {
		school.DistrictID = *updateSchool.DistrictID
	}

	if err = c.storer.Update(ctx, school); err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (c *Core) Delete(ctx *gin.Context) (error, int) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ErrInvalidID, http.StatusBadRequest
	}

	school, err := c.storer.GetByID(ctx, id)
	if err != nil {
		return err, http.StatusNotFound
	}

	if err = c.storer.Delete(ctx, school); err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (c *Core) GetSchoolByID(ctx *gin.Context, id uuid.UUID) (School, error, int) {
	school, err := c.storer.GetByID(ctx, id)
	if err != nil {
		return School{}, err, http.StatusNotFound
	}

	return school, nil, http.StatusOK
}

func (c *Core) GetSchoolsPaginated(ctx *gin.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]School, error, int) {
	if err := filter.Validate(); err != nil {
		return nil, err, http.StatusBadRequest
	}

	schools, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		schools = []School{}
	}

	return schools, nil, http.StatusOK
}

func (c *Core) Count(ctx *gin.Context, filter QueryFilter) (int, error, int) {
	if err := filter.Validate(); err != nil {
		return 0, err, http.StatusBadRequest
	}

	count, err := c.storer.Count(ctx, filter)
	if err != nil {
		count = 0
	}

	return count, err, http.StatusOK
}

func (c *Core) GetSchoolsByDistrictID(ctx *gin.Context) ([]School, error, int) {
	districtID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, ErrInvalidID, http.StatusBadRequest
	}

	schools, err := c.storer.GetByDistrict(ctx, int32(districtID))
	if err != nil {
		return nil, err, http.StatusNotFound
	}

	return schools, nil, http.StatusOK
}

func (c *Core) GetAllProvinces(ctx *gin.Context) ([]Province, error, int) {
	provinces, err := c.storer.GetAllProvinces(ctx)
	if err != nil {
		return nil, err, http.StatusNotFound
	}
	return provinces, nil, http.StatusOK
}

func (c *Core) GetDistrictsByProvinceID(ctx *gin.Context) ([]District, error, int) {
	provinceID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return nil, ErrInvalidID, http.StatusBadRequest
	}

	districts, err := c.storer.GetDistrictsByProvince(ctx, int32(provinceID))
	if err != nil {
		return nil, err, http.StatusNotFound
	}
	return districts, nil, http.StatusOK
}
