package user

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	Update(ctx *gin.Context, school sqlc.School) error
	Delete(ctx *gin.Context) error
	GetSchoolByID(ctx *gin.Context, id uuid.UUID) (sqlc.School, error)
	GetSchoolByName(ctx gin.Context, school sqlc.School) (sqlc.School, error)
	GetSchoolByDistrictID(ctx gin.Context, districtID int) ([]sqlc.School, error)
	GetAllProvinces(ctx *gin.Context) ([]sqlc.Province, error)
	GetDistrictsByProvince(ctx *gin.Context) ([]sqlc.District, error)
}
