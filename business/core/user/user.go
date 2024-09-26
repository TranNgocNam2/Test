package user

import (
	"Backend/internal/app"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrPhoneAlreadyExists = errors.New("phone already exists")
	ErrorUserAlreadyExist = errors.New("user already exists")
)

type Core struct {
	storer Storer
	app    *app.Application
}

func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

type Storer interface {
	Create(ctx *gin.Context, user User) error
	//Update(ctx *gin.Context, school sqlc.School) error
	//Delete(ctx *gin.Context) error
	GetByID(ctx *gin.Context, id string) (User, error)
	GetByEmail(ctx *gin.Context, email mail.Address) (User, error)
	GetByPhone(ctx *gin.Context, phone string) (User, error)
	//GetSchoolByName(ctx gin.Context, school sqlc.School) (sqlc.School, error)
	//GetSchoolByDistrictID(ctx gin.Context, districtID int) ([]sqlc.School, error)
	//GetAllProvinces(ctx *gin.Context) ([]sqlc.Province, error)
	//GetDistrictsByProvince(ctx *gin.Context) ([]sqlc.District, error)
}

func (c *Core) CreateUser(ctx *gin.Context, newUser NewUser) (error, int) {
	if _, err := c.storer.GetByEmail(ctx, newUser.Email); err == nil {
		return ErrEmailAlreadyExists, http.StatusBadRequest
	}

	if _, err := c.storer.GetByPhone(ctx, newUser.Phone); err == nil {
		return ErrPhoneAlreadyExists, http.StatusBadRequest
	}

	if _, err := c.storer.GetByID(ctx, newUser.ID); err == nil {
		return ErrorUserAlreadyExist, http.StatusBadRequest
	}

	user := User{
		ID:       newUser.ID,
		FullName: newUser.FullName,
		Email:    newUser.Email.Address,
		Phone:    newUser.Phone,
		Gender:   int16(newUser.Gender),
		Role:     int16(newUser.Role),
		Photo:    newUser.Photo,
	}
	user.School.ID = *newUser.SchoolID

	if err := c.storer.Create(ctx, user); err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func (c *Core) GetUserByID(ctx *gin.Context) (User, error, int) {
	id := ctx.Param("id")
	user, err := c.storer.GetByID(ctx, id)
	if err != nil {
		return User{}, err, http.StatusNotFound
	}

	return user, nil, http.StatusOK
}
