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
	ErrorUserAlreadyExist = errors.New("Người dùng đã tồn tại trong hệ thống!")
	ErrorUserNotFound     = errors.New("Người dùng không tồn tại trong hệ thống!")
)

type Core struct {
	//storer Storer
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

//type Storer interface {
//	Create(ctx *gin.Context, user User) error
//	GetByID(ctx *gin.Context, id string) (User, error)
//	GetByEmail(ctx *gin.Context, email mail.Address) (User, error)
//	GetByPhone(ctx *gin.Context, phone string) (User, error)
//}

func (c *Core) Create(ctx *gin.Context, newUser User) error {

	if _, err := c.queries.GetUserByEmail(ctx, newUser.Email.Address); err == nil {
		return ErrEmailAlreadyExists
	}

	if _, err := c.queries.GetUserByPhone(ctx, newUser.Phone); err == nil {
		return ErrPhoneAlreadyExists
	}

	if _, err := c.queries.GetUserByID(ctx, newUser.ID); err == nil {
		return ErrorUserAlreadyExist
	}

	schoolID := uuid.NullUUID{
		UUID:  uuid.Nil,
		Valid: false,
	}

	if newUser.School != nil {
		schoolID = uuid.NullUUID{
			UUID:  *newUser.School.ID,
			Valid: true,
		}
	}

	var dbUser = sqlc.CreateUserParams{
		ID:           newUser.ID,
		FullName:     newUser.FullName,
		Email:        newUser.Email.Address,
		Phone:        newUser.Phone,
		Gender:       newUser.Gender,
		ProfilePhoto: newUser.Photo,
		Role:         newUser.Role,
		SchoolID:     schoolID,
	}

	if err := c.queries.CreateUser(ctx, dbUser); err != nil {
		return err
	}

	return nil
}

func (c *Core) GetUserByID(ctx *gin.Context) (User, error) {
	id := ctx.Param("id")
	dbUser, err := c.queries.GetUserByID(ctx, id)
	if err != nil {
		return User{}, ErrorUserNotFound
	}

	if !dbUser.SchoolID.Valid {
		return toCoreUser(dbUser), nil
	}

	dbSchool, err := c.queries.GetSchoolByID(ctx, dbUser.SchoolID.UUID)
	if err != nil {
		return User{}, err
	}

	user := toCoreUser(dbUser)
	user.School.ID = &dbSchool.ID
	user.School.Name = &dbSchool.Name

	return user, nil
}
