package school

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type Core struct {
	storer Storer
	app    *app.Application
}

func NewCore(app *app.Application) *Core {
	return &Core{
		app: app,
	}
}

type Storer interface {
	Create(ctx *gin.Context, school sqlc.School) error
	Update(ctx *gin.Context, school School) error
	Delete(ctx *gin.Context) error
	GetSchoolByID(ctx *gin.Context, id uuid.UUID) (School, error)
	GetSchoolByName(ctx gin.Context, school School) (School, error)
	GetSchoolByDistrictID(ctx gin.Context, districtID int) ([]School, error)
	GetAllProvinces(ctx *gin.Context) ([]sqlc.Province, error)
	GetDistrictsByProvince(ctx *gin.Context) ([]sqlc.District, error)
}

func (c *Core) Create(ctx *gin.Context, school sqlc.CreateSchoolParams) error {
	if err := c.app.Queries.CreateSchool(ctx, school); err != nil {
		return err
	}
	return nil
}

func (c *Core) Delete(ctx *gin.Context) (error, int) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return fmt.Errorf("invalid school id"), http.StatusBadRequest
	}

	_, err = c.GetSchoolByID(ctx, id)
	if err != nil {
		return fmt.Errorf("school not found"), http.StatusNotFound
	}

	if err := c.app.Queries.DeleteSchool(ctx, id); err != nil {
		return fmt.Errorf("failed to delete school"), http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (c *Core) GetSchoolByID(ctx *gin.Context, id uuid.UUID) (sqlc.School, error) {
	school, err := c.app.Queries.GetSchoolByID(ctx, id)
	if err != nil {
		return sqlc.School{}, err
	}
	return school, nil
}

func (c *Core) GetAllProvinces(ctx *gin.Context) ([]sqlc.Province, error) {
	provinces, err := c.app.Queries.GetAllProvince(ctx)
	if err != nil {
		return nil, err
	}
	return provinces, nil
}

func (c *Core) GetDistrictsByProvince(ctx *gin.Context) ([]sqlc.District, error) {
	provinceID, err := strconv.Atoi(ctx.Param("province_id"))
	if err != nil {
		return nil, err
	}

	districts, err := c.app.Queries.GetDistrictsByProvince(ctx, int32(provinceID))
	if err != nil {
		return nil, err
	}
	return districts, nil
}
