package certificate

import (
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/common/status"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/page"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"gitlab.com/innovia69420/kit/enum/role"
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

func (c *Core) GetById(ctx *gin.Context, id uuid.UUID) (*Certificate, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Commit(ctx)

	qtx := c.queries.WithTx(tx)

	dbCert, err := qtx.GetCertificateById(ctx, id)
	if err != nil || dbCert.Status == Invalid {
		return nil, model.ErrCertificateNotFound
	}

	certificate, err := handleCertificateData(ctx, qtx, dbCert)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, nil
	}

	return certificate, nil
}

func (c *Core) Query(ctx *gin.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Certificate, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Commit(ctx)

	qtx := c.queries.WithTx(tx)
	filter.LearnerId, err = getUserId(ctx, filter.LearnerId, qtx)
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

	const q = `SELECT c.id, c.name, c.created_at, c.specialization_id, c.subject_id, c.class_id
					FROM certificates c`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, false)

	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var dbCertificates []sqlc.Certificate
	err = pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbCertificates)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, nil
	}

	if dbCertificates == nil {
		return nil, nil
	}

	var certificates []Certificate

	for _, dbCert := range dbCertificates {
		certificate, err := handleCertificateData(ctx, qtx, dbCert)
		if err != nil {
			c.logger.Error(err.Error())
			return nil, nil
		}

		certificates = append(certificates, *certificate)
	}

	return certificates, nil
}

func (c *Core) Count(ctx *gin.Context, filter QueryFilter) int {
	if err := filter.Validate(); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	data := map[string]interface{}{}

	const q = `SELECT
                        count(1)
               FROM
                        certificates c`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, false)

	var count struct {
		Count int `db:"count"`
	}

	if err := pgx.NamedQueryStruct(ctx, c.logger, c.db, buf.String(), data, &count); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	return count.Count
}

//	func (c *Core) GetByLearnerAndSpecialization(ctx *gin.Context, specId uuid.UUID, learnerId *string) ([]Certificate, error) {
//		tx, err := c.pool.Begin(ctx)
//		if err != nil {
//			return nil, err
//		}
//		defer tx.Commit(ctx)
//
//		qtx := c.queries.WithTx(tx)
//		learnerId, err = getUserId(ctx, learnerId, qtx)
//		if err != nil || learnerId == nil {
//			return nil, err
//		}
//
//		subjectIds, err := qtx.GetSubjectIdsBySpecialization(ctx, specId)
//		if err != nil {
//			return nil, err
//		}
//
//		cer
//	}
func getUserId(ctx *gin.Context, userId *string, qtx *sqlc.Queries) (*string, error) {
	if userId != nil {
		learner, err := qtx.GetUserById(ctx, *userId)
		if err != nil || learner.AuthRole != role.LEARNER ||
			status.User(learner.Status) != status.Valid || !learner.IsVerified {
			return nil, model.ErrUserNotFound
		}
	} else {
		learner, err := middleware.AuthorizeVerifiedLearner(ctx, qtx)
		if err != nil {
			return nil, err
		}
		return &learner.ID, nil
	}
	return userId, nil
}

func handleCertificateData(ctx *gin.Context, qtx *sqlc.Queries, dbCert sqlc.Certificate) (*Certificate, error) {
	certificate := toCoreCertificate(dbCert)
	if dbCert.SpecializationID != nil {
		dbSpecialization, err := qtx.GetSpecializationById(ctx, *dbCert.SpecializationID)
		if err != nil {
			return nil, err
		}
		certificate.Specialization = toCoreSpecialization(dbSpecialization)
	}

	if dbCert.SubjectID != nil {
		dbProgram, err := qtx.GetProgramByClassId(ctx, *dbCert.ClassID)
		if err != nil {
			return nil, err
		}

		dbSubject, err := qtx.GetSubjectById(ctx, *dbCert.SubjectID)
		if err != nil {
			return nil, err
		}
		certificate.Subject = toCoreSubject(dbSubject, dbProgram)
	}
	return &certificate, nil
}
