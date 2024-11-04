package user

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"gitlab.com/innovia69420/kit/enum/role"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Core struct {
	db      *sqlx.DB
	queries *sqlc.Queries
	logger  *zap.Logger
	pool    *pgxpool.Pool
}

func NewCore(app *app.Application) *Core {
	return &Core{
		db:      app.DB,
		queries: app.Queries,
		logger:  app.Logger,
		pool:    app.Pool,
	}
}

func (c *Core) Create(ctx *gin.Context, newUser NewUser) error {
	if _, err := c.queries.GetUserByEmail(ctx, newUser.Email.Address); err == nil {
		return model.ErrEmailAlreadyExists
	}

	if _, err := c.queries.GetUserByID(ctx, newUser.ID); err == nil {
		return model.ErrUserAlreadyExist
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

func (c *Core) Verify(ctx *gin.Context, id string, verifyUser VerifyUser) error {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbUser, err := c.queries.GetUserByID(ctx, id)
	if err != nil {
		return model.ErrUserNotFound
	}

	if dbUser.AuthRole != role.LEARNER {
		return model.ErrUserCannotBeVerified
	}

	if dbUser.Image == nil && verifyUser.Status == Verified {
		return model.ErrInvalidVerificationInfo
	}

	dbVerifyUser := sqlc.VerifyUserParams{
		ID:         dbUser.ID,
		Status:     verifyUser.Status,
		VerifiedBy: &staffId,
	}
	if err = c.queries.VerifyUser(ctx, dbVerifyUser); err != nil {
		return err
	}

	return nil
}

func (c *Core) GetByID(ctx *gin.Context, id string) (User, error) {
	dbUser, err := c.queries.GetUserByID(ctx, id)
	if err != nil {
		return User{}, model.ErrUserNotFound
	}
	user := toCoreUser(dbUser)
	if dbUser.SchoolID != nil {
		dbSchool, _ := c.queries.GetSchoolByID(ctx, *dbUser.SchoolID)
		user.School = &struct {
			ID   *uuid.UUID `json:"id"`
			Name *string    `json:"name"`
		}{
			ID:   &dbSchool.ID,
			Name: &dbSchool.Name,
		}
	}

	return user, nil
}

func (c *Core) Update(ctx *gin.Context, id string, updatedUser UpdateUser) error {
	dbUser, err := c.queries.GetUserByID(ctx, id)
	if err != nil {
		return model.ErrUserNotFound
	}

	if updatedUser.Email.Address != dbUser.Email {
		if _, err = c.queries.GetUserByEmail(ctx, updatedUser.Email.Address); err == nil {
			return model.ErrEmailAlreadyExists
		}
	}

	if updatedUser.Phone != "" && updatedUser.Phone != *dbUser.Phone {
		if _, err = c.queries.GetUserByPhone(ctx, &updatedUser.Phone); err == nil {
			return model.ErrPhoneAlreadyExists
		}
	}

	if updatedUser.SchoolID == nil {
		updatedUser.SchoolID = dbUser.SchoolID
	}

	var dbUserUpdate = sqlc.UpdateUserParams{
		FullName:     &updatedUser.FullName,
		Email:        updatedUser.Email.Address,
		Phone:        &updatedUser.Phone,
		Gender:       &updatedUser.Gender,
		SchoolID:     updatedUser.SchoolID,
		ProfilePhoto: &updatedUser.Photo,
		Status:       Pending,
		Image:        updatedUser.Image,
		ID:           dbUser.ID,
	}

	if err = c.queries.UpdateUser(ctx, dbUserUpdate); err != nil {
		return err
	}

	return nil
}
