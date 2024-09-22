package school

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
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
	Create(ctx gin.Context, school School) error
	Update(ctx gin.Context, school School) error
	Delete(ctx gin.Context, school School) error
	GetSchoolByName(ctx gin.Context, school School) (School, error)
	GetSchoolByDistrictID(ctx gin.Context, districtID int) ([]School, error)
	GetAllProvinces(ctx *gin.Context) ([]sqlc.Province, error)
	GetDistrictsByProvince(ctx *gin.Context) ([]sqlc.District, error)
}

func (c *Core) Create(ctx gin.Context, school School) error {

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
