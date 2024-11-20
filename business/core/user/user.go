package user

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/common/status"
	"Backend/internal/middleware"
	"gitlab.com/innovia69420/kit/enum/role"

	"github.com/gin-gonic/gin"
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

	if _, err := c.queries.GetUserById(ctx, newUser.ID); err == nil {
		return model.ErrUserAlreadyExist
	}

	var dbUser = sqlc.CreateUserParams{
		ID:       newUser.ID,
		Email:    newUser.Email.Address,
		AuthRole: newUser.Role,
		FullName: &newUser.FullName,
	}

	if err := c.queries.CreateUser(ctx, dbUser); err != nil {
		return err
	}

	return nil
}

func (c *Core) Verify(ctx *gin.Context, id string, verifyLearner VerifyLearner) error {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbUser, err := c.queries.GetUserById(ctx, id)
	if err != nil {
		return model.ErrUserNotFound
	}

	if dbUser.AuthRole != role.LEARNER {
		return model.ErrUserCannotBeVerified
	}

	verifyUser, err := c.queries.GetLearnerVerificationByUserId(ctx, dbUser.ID)
	if err != nil || (verifyUser.ImageLink == nil &&
		status.Verification(verifyLearner.Status) == status.Verified) {
		return model.ErrInvalidVerificationInfo
	}

	dbVerifyUser := sqlc.VerifyLearnerParams{
		VerifiedBy: &staffId,
		Status:     verifyLearner.Status,
		LearnerID:  dbUser.ID,
	}
	if err = c.queries.VerifyLearner(ctx, dbVerifyUser); err != nil {
		return err
	}

	return nil
}

func (c *Core) GetByID(ctx *gin.Context, id string) (Details, error) {
	dbUser, err := c.queries.GetUserById(ctx, id)
	if err != nil {
		return Details{}, model.ErrUserNotFound
	}
	user := toCoreUser(dbUser)

	if dbUser.AuthRole == role.LEARNER {
		learnerVerification, _ := c.queries.GetLearnerVerificationByUserId(ctx, dbUser.ID)
		dbSchool, _ := c.queries.GetSchoolById(ctx, learnerVerification.SchoolID)
		user.School = &School{
			ID:   dbSchool.ID,
			Name: dbSchool.Name,
		}
		user.Type = &learnerVerification.Type
		user.VerifiedStatus = &learnerVerification.Status
	}

	return user, nil
}

func (c *Core) GetCurrent(ctx *gin.Context) (Details, error) {
	dbUser, err := middleware.AuthorizeUser(ctx, c.queries)
	if err != nil {
		return Details{}, model.ErrUserNotFound
	}
	user := toCoreUser(*dbUser)

	if dbUser.AuthRole == role.LEARNER {
		learnerVerification, _ := c.queries.GetLearnerVerificationByUserId(ctx, dbUser.ID)
		dbSchool, _ := c.queries.GetSchoolById(ctx, learnerVerification.SchoolID)
		user.School = &School{
			ID:   dbSchool.ID,
			Name: dbSchool.Name,
		}
		user.Type = &learnerVerification.Type
		user.VerifiedStatus = &learnerVerification.Status
	}

	return user, nil
}

func (c *Core) Update(ctx *gin.Context, id string, updatedUser UpdateUser) error {
	dbUser, err := c.queries.GetUserById(ctx, id)
	if err != nil {
		return model.ErrUserNotFound
	}

	if err = c.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		FullName:     &updatedUser.FullName,
		Phone:        &updatedUser.Phone,
		ProfilePhoto: &updatedUser.Photo,
		ID:           dbUser.ID,
	}); err != nil {
		return err
	}

	return nil
}

func (c *Core) Handle(ctx *gin.Context, id string) error {
	_, err := middleware.AuthorizeAdmin(ctx, c.queries)
	if err != nil {
		return err
	}

	user, err := c.queries.GetUserById(ctx, id)
	if err != nil {
		return model.ErrUserNotFound
	}
	if status.User(user.Status) == status.Valid {
		user.Status = int32(status.Invalid)
	} else {
		user.Status = int32(status.Valid)
	}

	if err := c.queries.HandleUserStatus(ctx, sqlc.HandleUserStatusParams{
		ID:     user.ID,
		Status: user.Status,
	}); err != nil {
		return err
	}

	return nil
}
