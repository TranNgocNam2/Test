package user

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	ErrEmailAlreadyExists = errors.New("Email đã tồn tại trong hệ thống!")
	ErrPhoneAlreadyExists = errors.New("Số điện thoại đã tồn tại trong hệ thống!")
	ErrUserAlreadyExist   = errors.New("Người dùng đã tồn tại trong hệ thống!")
	ErrUserNotFound       = errors.New("Người dùng không tồn tại trong hệ thống!")
)

type Core struct {
	db      *sqlx.DB
	queries *sqlc.Queries
	logger  *zap.Logger
}

func NewCore(app *app.Application) *Core {
	return &Core{
		db:      app.Db,
		queries: app.Queries,
		logger:  app.Logger,
	}
}

func (c *Core) Create(ctx *gin.Context, newUser User) error {
	if _, err := c.queries.GetUserByEmail(ctx, newUser.Email.Address); err == nil {
		return ErrEmailAlreadyExists
	}

	if _, err := c.queries.GetUserByID(ctx, newUser.ID); err == nil {
		return ErrUserAlreadyExist
	}

	var dbUser = sqlc.CreateUserParams{
		ID:       newUser.ID,
		Email:    newUser.Email.Address,
		AuthRole: newUser.Role,
	}

	if err := c.queries.CreateUser(ctx, dbUser); err != nil {
		return err
	}

	return nil
}

func (c *Core) GetByID(ctx *gin.Context, id string) (User, error) {
	dbUser, err := c.queries.GetUserByID(ctx, id)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	user := toCoreUser(dbUser)
	if dbUser.SchoolID != nil {
		dbSchool, _ := c.queries.GetSchoolByID(ctx, *dbUser.SchoolID)
		user.School = &struct {
			ID   *uuid.UUID
			Name *string
		}{
			ID:   &dbSchool.ID,
			Name: &dbSchool.Name,
		}
	}

	return user, nil
}

func (c *Core) Update(ctx *gin.Context, updatedUser User) error {
	dbUser, err := c.queries.GetUserByID(ctx, updatedUser.ID)
	if err != nil {
		return ErrUserNotFound
	}

	if updatedUser.Email.Address != dbUser.Email {
		if _, err = c.queries.GetUserByEmail(ctx, updatedUser.Email.Address); err == nil {
			return ErrEmailAlreadyExists
		}
	}

	if updatedUser.Phone != nil && *updatedUser.Phone != *dbUser.Phone {
		if _, err = c.queries.GetUserByPhone(ctx, updatedUser.Phone); err == nil {
			return ErrPhoneAlreadyExists
		}
	}

	if updatedUser.School == nil {
		updatedUser.School.ID = dbUser.SchoolID
	}

	var dbUserUpdate = sqlc.UpdateUserParams{
		FullName:     updatedUser.FullName,
		Email:        updatedUser.Email.Address,
		Phone:        updatedUser.Phone,
		Gender:       updatedUser.Gender,
		SchoolID:     updatedUser.School.ID,
		ProfilePhoto: updatedUser.Photo,
		ID:           updatedUser.ID,
	}

	if err = c.queries.UpdateUser(ctx, dbUserUpdate); err != nil {
		return err
	}

	return nil
}
