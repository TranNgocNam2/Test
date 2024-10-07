package user

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gitlab.com/innovia69420/kit/enum/role"
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

	if _, err := c.queries.GetUserByPhone(ctx, newUser.Phone); err == nil {
		return ErrPhoneAlreadyExists
	}

	if _, err := c.queries.GetUserByID(ctx, newUser.ID); err == nil {
		return ErrUserAlreadyExist
	}

	var dbUser = sqlc.CreateUserParams{
		ID:           newUser.ID,
		FullName:     newUser.FullName,
		Email:        newUser.Email.Address,
		Phone:        newUser.Phone,
		Gender:       newUser.Gender,
		ProfilePhoto: newUser.Photo,
		AuthRole:     newUser.Role,
	}

	if err := c.queries.CreateUser(ctx, dbUser); err != nil {
		return err
	}

	if dbUser.AuthRole == role.LEARNER {
		dbLearner := sqlc.CreateLeanerParams{
			ID:       dbUser.ID,
			SchoolID: *newUser.School.ID,
		}
		if err := c.queries.CreateLeaner(ctx, dbLearner); err != nil {
			return err
		}
	} else {
		createdBy := sql.NullString{
			String: *newUser.CreatedBy,
			Valid:  false,
		}
		if newUser.CreatedBy != nil {
			createdBy.Valid = true
		}

		dbStaff := sqlc.CreateStaffParams{
			ID:        dbUser.ID,
			Role:      dbUser.AuthRole,
			CreatedBy: createdBy,
		}
		if err := c.queries.CreateStaff(ctx, dbStaff); err != nil {
			return err
		}
	}

	return nil
}

func (c *Core) GetUserByID(ctx *gin.Context) (User, error) {
	id := ctx.Param("id")
	dbUser, err := c.queries.GetUserByID(ctx, id)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	user := toCoreUser(dbUser)
	if user.Role == role.LEARNER {
		dbLearner, _ := c.queries.GetLearnerByID(ctx, dbUser.ID)
		dbSchool, _ := c.queries.GetSchoolByID(ctx, dbLearner.SchoolID)
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
