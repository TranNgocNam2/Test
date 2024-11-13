package learner

import (
	"Backend/business/core/learner/certificate"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
	"time"
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

func (c *Core) JoinClass(ctx *gin.Context, classAccess ClassAccess) error {
	learner, err := middleware.AuthorizeVerifiedLearner(ctx, c.queries)
	if err != nil {
		return err
	}

	dbClass, err := c.queries.GetClassCompletedByCode(ctx, classAccess.Code)
	if err != nil {
		return model.ErrClassNotFound
	}
	if dbClass.StartDate.Before(time.Now()) {
		return model.ErrClassStarted
	}

	if strings.Compare(dbClass.Password, classAccess.Password) != 0 {
		return model.ErrWrongPassword
	}

	dbSlots, _ := c.queries.GetSlotsByClassId(ctx, dbClass.ID)

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	classLearner := sqlc.AddLearnerToClassParams{
		ID:        uuid.New(),
		ClassID:   dbClass.ID,
		LearnerID: learner.ID,
	}

	err = qtx.AddLearnerToClass(ctx, classLearner)
	if err != nil {
		return err
	}

	for _, dbSlot := range dbSlots {
		err = qtx.GenerateLearnerAttendance(ctx, sqlc.GenerateLearnerAttendanceParams{
			ClassLearnerID: classLearner.ID,
			SlotID:         dbSlot.ID,
		})
		if err != nil {
			return err
		}
	}
	tx.Commit(ctx)
	return nil
}

func (c *Core) JoinSpecialization(ctx *gin.Context, specializationId uuid.UUID) error {
	learner, err := middleware.AuthorizeVerifiedLearner(ctx, c.queries)
	if err != nil {
		return err
	}

	learnerSpec, _ := c.queries.CountLearnerInSpecialization(ctx,
		sqlc.CountLearnerInSpecializationParams{
			LearnerID:        learner.ID,
			SpecializationID: specializationId,
		})

	if learnerSpec > 0 {
		return model.ErrAlreadyJoinedSpecialization
	}

	specialization, err := c.queries.GetPublishedSpecializationById(ctx, specializationId)
	if err != nil {
		return model.ErrSpecNotFound
	}

	err = c.queries.AddLearnerToSpecialization(ctx, sqlc.AddLearnerToSpecializationParams{
		LearnerID:        learner.ID,
		SpecializationID: specialization.ID,
	})
	if err != nil {
		return err
	}

	subjectIds, err := c.queries.GetSubjectIdsBySpecialization(ctx, specializationId)
	if err != nil {
		c.logger.Error(err.Error())
		return nil
	}

	learnerCertParams := sqlc.GetCertificationsByLearnerAndSubjectsParams{
		LearnerID:  learner.ID,
		SubjectIds: subjectIds,
		Status:     certificate.Valid,
	}

	subjectCerts, err := c.queries.GetCertificationsByLearnerAndSubjects(ctx, learnerCertParams)
	if err != nil {
		c.logger.Error(err.Error())
		return nil
	}
	if len(subjectCerts) == len(subjectIds) {
		specCert := sqlc.CreateSpecializationCertificateParams{
			LearnerID:        learner.ID,
			SpecializationID: &specialization.ID,
			Name:             specialization.Name,
			Status:           certificate.Valid,
		}

		err = c.queries.CreateSpecializationCertificate(ctx, specCert)
		if err != nil {
			c.logger.Error(err.Error())
			return nil
		}
	}
	return nil
}

func (c *Core) SubmitAttendance(ctx *gin.Context, classId uuid.UUID, attendanceSubmission AttendanceSubmission) error {
	learner, err := middleware.AuthorizeVerifiedLearner(ctx, c.queries)
	if err != nil {
		return err
	}

	class, err := c.queries.GetClassById(ctx, classId)
	if err != nil {
		return model.ErrClassNotFound
	}

	classLearner, err := c.queries.GetLearnerByClassId(ctx,
		sqlc.GetLearnerByClassIdParams{
			ClassID:   class.ID,
			LearnerID: learner.ID,
		})
	if err != nil {
		return model.LearnerNotInClass
	}

	slot, err := c.queries.GetSlotByClassIdAndIndex(ctx,
		sqlc.GetSlotByClassIdAndIndexParams{
			ClassID: class.ID,
			Index:   attendanceSubmission.Index,
		})
	if err != nil {
		return model.ErrSlotNotFound
	}

	if strings.Compare(*slot.AttendanceCode, attendanceSubmission.AttendanceCode) != 0 {
		return model.ErrInvalidAttendanceCode
	}

	now := time.Now().Format(time.DateTime)
	currentTime, _ := time.Parse(time.DateTime, now)

	if slot.EndTime.Before(currentTime) {
		return model.ErrSlotEnded
	}

	if currentTime.Before(*slot.StartTime) {
		return model.ErrSlotNotStarted
	}

	learnerAttendance, _ := c.queries.GetLearnerAttendanceByClassLearnerAndSlot(ctx,
		sqlc.GetLearnerAttendanceByClassLearnerAndSlotParams{
			ClassLearnerID: classLearner.ID,
			SlotID:         slot.ID,
		})

	err = c.queries.SubmitLearnerAttendance(ctx,
		sqlc.SubmitLearnerAttendanceParams{
			Status: Attended,
			ID:     learnerAttendance.ID,
		})
	if err != nil {
		return err
	}

	return nil
}
