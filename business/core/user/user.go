package user

import (
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/common/status"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/page"
	"Backend/internal/slice"
	"bytes"
	"github.com/google/uuid"
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

	var isVerified = false
	if newUser.Role != role.LEARNER {
		_, err := middleware.AuthorizeAdmin(ctx, c.queries)
		if err != nil {
			return err
		}
		isVerified = true
	}

	var dbUser = sqlc.CreateUserParams{
		ID:         newUser.ID,
		Email:      newUser.Email.Address,
		AuthRole:   newUser.Role,
		FullName:   &newUser.FullName,
		IsVerified: isVerified,
	}

	if err := c.queries.CreateUser(ctx, dbUser); err != nil {
		return err
	}

	return nil
}

func (c *Core) Verify(ctx *gin.Context, verificationId uuid.UUID, verifyLearner VerifyLearner) error {
	admin, err := middleware.AuthorizeAdmin(ctx, c.queries)
	if err != nil {
		return err
	}

	verifyUser, err := c.queries.GetLearnerVerificationById(ctx, verificationId)
	if err != nil || (verifyUser.ImageLink == nil &&
		status.Verification(verifyLearner.Status) == status.Verified) {
		return model.ErrInvalidVerificationInfo
	}

	dbUser, err := c.queries.GetUserById(ctx, verifyUser.LearnerID)
	if err != nil {
		return model.ErrUserNotFound
	}

	if dbUser.AuthRole != role.LEARNER {
		return model.ErrUserCannotBeVerified
	}
	if dbUser.IsVerified {
		return model.ErrLearnerAlreadyVerified
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := c.queries.WithTx(tx)

	if err = qtx.VerifyLearner(ctx, sqlc.VerifyLearnerParams{
		VerifiedBy: &admin.ID,
		Status:     verifyLearner.Status,
		LearnerID:  dbUser.ID,
		Note:       verifyLearner.Note,
		ID:         verifyUser.ID,
	}); err != nil {
		return err
	}

	if status.Verification(verifyLearner.Status) == status.Verified {
		if err = qtx.UpdateVerification(ctx, sqlc.UpdateVerificationParams{
			IsVerified: true,
			SchoolID:   &verifyUser.SchoolID,
			Type:       &verifyUser.Type,
			ID:         dbUser.ID,
		}); err != nil {
			return err
		}
	}

	tx.Commit(ctx)
	return nil
}

func (c *Core) GetByID(ctx *gin.Context, id string) (Details, error) {
	dbUser, err := c.queries.GetUserById(ctx, id)
	if err != nil {
		return Details{}, model.ErrUserNotFound
	}
	user := toCoreUser(dbUser)

	if dbUser.AuthRole == role.LEARNER && dbUser.IsVerified {
		dbSchool, _ := c.queries.GetSchoolById(ctx, *dbUser.SchoolID)
		user.School = &School{
			ID:   dbSchool.ID,
			Name: dbSchool.Name,
		}
	}
	return user, nil
}

func (c *Core) GetCurrent(ctx *gin.Context) (Details, error) {
	dbUser, err := middleware.AuthorizeUser(ctx, c.queries)
	if err != nil {
		return Details{}, model.ErrUserNotFound
	}
	user := toCoreUser(*dbUser)

	if dbUser.AuthRole == role.LEARNER && dbUser.IsVerified && dbUser.SchoolID != nil {
		dbSchool, _ := c.queries.GetSchoolById(ctx, *dbUser.SchoolID)
		user.School = &School{
			ID:   dbSchool.ID,
			Name: dbSchool.Name,
		}
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
		ProfilePhoto: &updatedUser.Photo,
		ID:           dbUser.ID,
	}); err != nil {
		return err
	}

	return nil
}

func (c *Core) Handle(ctx *gin.Context, id string) (string, error) {
	_, err := middleware.AuthorizeAdmin(ctx, c.queries)
	if err != nil {
		return "", err
	}

	user, err := c.queries.GetUserById(ctx, id)
	if err != nil {
		return "", model.ErrUserNotFound
	}

	if status.User(user.Status) == status.Valid {
		user.Status = int32(status.Invalid)
	} else {
		user.Status = int32(status.Valid)
	}

	if err = c.queries.HandleUserStatus(ctx, sqlc.HandleUserStatusParams{
		Status: user.Status,
		ID:     user.ID,
	}); err != nil {
		return "", err
	}
	return status.GetUserStatus(user.Status), nil
}

func (c *Core) GetUsers(ctx *gin.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]User, error) {
	if *filter.Role == role.TEACHER || *filter.Role == role.LEARNER {
		_, err := middleware.AuthorizeStaff(ctx, c.queries)
		if err != nil {
			return nil, err
		}
	}
	if *filter.Role == role.MANAGER {
		_, err := middleware.AuthorizeAdmin(ctx, c.queries)
		if err != nil {
			return nil, err
		}
	}
	if err := filter.Validate(); err != nil {
		return nil, nil
	}
	data := map[string]interface{}{
		"offset":        (page.Number - 1) * page.Size,
		"rows_per_page": page.Size,
	}
	const q = `SELECT u.id, u.full_name, u.email, u.auth_role, u.status, u.profile_photo, u.is_verified, u.phone, u.status, u.school_id, u.type 
				FROM users u`
	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, false, UserStatus)
	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())
	var dbUsers []sqlc.User
	err := pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbUsers)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, nil
	}
	if dbUsers == nil {
		return nil, nil
	}

	var users []User
	for _, dbUser := range dbUsers {
		user := User{
			ID:       dbUser.ID,
			FullName: *dbUser.FullName,
			Email:    dbUser.Email,
			Phone:    dbUser.Phone,
			Photo:    dbUser.ProfilePhoto,
			Role:     &dbUser.AuthRole,
			Status:   &dbUser.Status,
			Type:     dbUser.Type,
		}
		if *filter.Role == role.LEARNER && dbUser.IsVerified && dbUser.SchoolID != nil {
			dbSchool, _ := c.queries.GetSchoolById(ctx, *dbUser.SchoolID)
			user.School = &School{
				ID:   dbSchool.ID,
				Name: dbSchool.Name,
			}
		}
		users = append(users, user)
	}
	return users, nil
}

func (c *Core) CountUsers(ctx *gin.Context, filter QueryFilter) int {
	if err := filter.Validate(); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	data := map[string]interface{}{}

	const q = `SELECT
						 COUNT(u.id) AS count
				FROM users u`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, false, UserStatus)

	var count struct {
		Count int `db:"count"`
	}

	if err := pgx.NamedQueryStruct(ctx, c.logger, c.db, buf.String(), data, &count); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	return count.Count
}

func (c *Core) GetVerificationUsers(ctx *gin.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Verification, error) {
	_, err := middleware.AuthorizeAdmin(ctx, c.queries)
	if err != nil {
		return nil, err
	}
	if err := filter.Validate(); err != nil {
		return nil, nil
	}
	data := map[string]interface{}{
		"offset":        (page.Number - 1) * page.Size,
		"rows_per_page": page.Size,
	}
	const q = `SELECT u.id AS user_id, u.full_name, u.email, u.status,
       					vl.id, vl.image_link AS image_link, vl.type, vl.note, vl.created_at, vl.status,
       					s.id AS school_id, s.name AS school_name
				FROM users u
					JOIN verification_learners vl ON u.id = vl.learner_id
					JOIN schools s ON vl.school_id = s.id`
	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, false, VerifiedStatus)
	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var verificationUsers []sqlc.GetVerificationLearnersRow
	err = pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &verificationUsers)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, nil
	}
	if verificationUsers == nil {
		return nil, nil
	}

	var results []Verification
	for _, verificationUser := range verificationUsers {
		result := Verification{
			ID:        verificationUser.ID,
			Status:    verificationUser.Status,
			Note:      verificationUser.Note,
			ImageLink: slice.ParseFromString(verificationUser.ImageLink),
			Type:      verificationUser.Type,
			CreatedAt: verificationUser.CreatedAt,
			School: School{
				ID:   verificationUser.SchoolID,
				Name: verificationUser.SchoolName,
			},
			User: User{
				ID:       verificationUser.UserID,
				FullName: *verificationUser.FullName,
				Email:    verificationUser.Email,
			},
		}
		results = append(results, result)
	}
	return results, nil
}

func (c *Core) CountVerificationUsers(ctx *gin.Context, filter QueryFilter) int {
	if err := filter.Validate(); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	data := map[string]interface{}{}

	const q = `SELECT
                         COUNT(vl.id) AS count
				FROM users u
					JOIN verification_learners vl ON u.id = vl.learner_id
					JOIN schools s ON vl.school_id = s.id`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, false, VerifiedStatus)

	var count struct {
		Count int `db:"count"`
	}

	if err := pgx.NamedQueryStruct(ctx, c.logger, c.db, buf.String(), data, &count); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	return count.Count
}

func (c *Core) CreateLearner(ctx *gin.Context, learner NewLearner) error {
	_, err := middleware.AuthorizeAdmin(ctx, c.queries)
	if err != nil {
		return err
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := c.queries.WithTx(tx)

	if _, err := qtx.GetUserByEmail(ctx, learner.Email); err == nil {
		return model.ErrEmailAlreadyExists
	}

	if _, err := qtx.GetUserById(ctx, learner.ID); err == nil {
		return model.ErrUserAlreadyExist
	}

	if _, err := qtx.GetSchoolById(ctx, learner.SchoolID); err != nil {
		return model.ErrSchoolNotFound
	}

	err = qtx.CreateLearner(ctx, sqlc.CreateLearnerParams{
		ID:         learner.ID,
		Email:      learner.Email,
		AuthRole:   role.LEARNER,
		FullName:   &learner.FullName,
		IsVerified: true,
		SchoolID:   &learner.SchoolID,
	})
	if err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}

func (c *Core) UpdateLearner(ctx *gin.Context, id string, updatedLearner UpdateLearner) error {
	_, err := middleware.AuthorizeAdmin(ctx, c.queries)
	if err != nil {
		return err
	}
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := c.queries.WithTx(tx)

	dbUser, err := qtx.GetUserById(ctx, id)
	if err != nil {
		return model.ErrUserNotFound
	}

	if _, err := qtx.GetSchoolById(ctx, updatedLearner.SchoolID); err != nil {
		return model.ErrSchoolNotFound
	}

	if err = qtx.UpdateLearner(ctx, sqlc.UpdateLearnerParams{
		SchoolID: &updatedLearner.SchoolID,
		Type:     &updatedLearner.Type,
		ID:       dbUser.ID,
	}); err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}
