package user

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"database/sql"
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
	if dbUser.SchoolID.Valid {
		dbSchool, _ := c.queries.GetSchoolByID(ctx, dbUser.SchoolID.UUID)
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

	phoneNumber := sql.NullString{
		String: *updatedUser.Phone,
		Valid:  true,
	}
	if updatedUser.Phone != nil && *updatedUser.Phone != *dbUser.Phone {
		if _, err = c.queries.GetUserByPhone(ctx, phoneNumber); err == nil {
			return ErrPhoneAlreadyExists
		}
	}

	schoolID := uuid.NullUUID{
		UUID:  uuid.Nil,
		Valid: false,
	}

	if updatedUser.School != nil {
		schoolID = uuid.NullUUID{
			UUID:  *updatedUser.School.ID,
			Valid: true,
		}
	} else {
		schoolID = dbUser.SchoolID
	}

	var dbUserUpdate = sqlc.UpdateUserParams{
		FullName: sql.NullString{
			String: *updatedUser.FullName,
			Valid:  true,
		},
		Email: updatedUser.Email.Address,
		Phone: phoneNumber,
		Gender: sql.NullInt16{
			Int16: updatedUser.Gender,
			Valid: true,
		},
		SchoolID: schoolID,
		ProfilePhoto: sql.NullString{
			String: *updatedUser.Photo,
			Valid:  true,
		},
		ID: updatedUser.ID,
	}

	if err = c.queries.UpdateUser(ctx, dbUserUpdate); err != nil {
		return err
	}

	return nil
}
