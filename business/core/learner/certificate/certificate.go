package certificate

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"fmt"
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

func (c *Core) GetById(ctx *gin.Context, id uuid.UUID) (*Certificate, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Commit(ctx)

	qtx := c.queries.WithTx(tx)

	dbCert, err := qtx.GetCertificateById(ctx, id)
	if err != nil || dbCert.Status == Invalid {
		fmt.Println(err)
		return nil, model.ErrCertificateNotFound
	}

	certificate := Certificate{
		ID:             dbCert.ID,
		Name:           dbCert.Name,
		CreatedAt:      dbCert.CreatedAt,
		Specialization: nil,
		Subject:        nil,
	}

	if dbCert.SpecializationID != nil {
		dbSpecialization, err := qtx.GetSpecializationById(ctx, *dbCert.SpecializationID)
		if err != nil {
			c.logger.Error(err.Error())
			return nil, nil
		}
		certificate.Specialization = &Specialization{
			ID:          dbSpecialization.ID,
			Name:        dbSpecialization.Name,
			Code:        dbSpecialization.Code,
			TimeAmount:  *dbSpecialization.TimeAmount,
			ImageLink:   *dbSpecialization.ImageLink,
			Description: *dbSpecialization.Description,
		}
	}

	if dbCert.SubjectID != nil {
		dbSubject, err := qtx.GetSubjectById(ctx, *dbCert.SubjectID)
		if err != nil {
			c.logger.Error(err.Error())
			return nil, nil
		}
		certificate.Subject = &Subject{
			ID:          dbSubject.ID,
			Name:        dbSubject.Name,
			Code:        dbSubject.Code,
			Description: *dbSubject.Description,
			ImageLink:   *dbSubject.ImageLink,
		}

		if dbCert.ClassID == nil {
			return &certificate, nil
		}

		program, err := qtx.GetProgramByClassId(ctx, *dbCert.ClassID)
		if err != nil {
			c.logger.Error(err.Error())
			return &certificate, nil
		}
		certificate.Program = &Program{
			ID:        program.ID,
			Name:      program.Name,
			StartDate: program.StartDate,
			EndDate:   program.EndDate,
		}
	}

	return &certificate, nil
}

//func (c *Core) Query(ctx *gin.Context, learnerId string) ([]Certificate, error) {
//	tx, err := c.pool.Begin(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	qtx := c.queries.WithTx(tx)
//	learner, err := qtx.GetUserById(ctx, learnerId)
//	learnerId = learner.ID
//	if err != nil || !learner.IsVerified ||
//		learner.AuthRole != role.LEARNER ||
//		learner.Status == Invalid {
//		currentUser, err := qtx.GetUserById(ctx, learnerId)
//		if err != nil {
//			learnerId = ""
//		}
//		learnerId = currentUser.ID
//	}
//
//	return nil, nil
//}
