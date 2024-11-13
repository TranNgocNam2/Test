package learner

import (
	"Backend/business/core/learner/certificate"
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"bytes"
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
	if dbClass.StartDate.Before(time.Now().UTC()) {
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

	now := time.Now().UTC().Format(time.DateTime)
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

func (c *Core) GetLearnersInClass(ctx *gin.Context, classId uuid.UUID, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Learner, error) {
	_, err := middleware.AuthorizeUser(ctx, c.queries)
	if err != nil {
		return nil, err
	}
	if err := filter.Validate(); err != nil {
		return nil, nil
	}

	data := map[string]interface{}{
		"class_id":      classId,
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT
						u.id, u.full_name AS full_name, u.email, u.phone, u.gender, u.profile_photo, u.status AS status, 
						s.id AS school_id, s.name AS school_name, cl.id AS class_learner_id
			FROM users u
				JOIN class_learners cl ON u.id = cl.learner_id
				JOIN classes c ON cl.class_id = c.id
				JOIN schools s ON s.id = u.school_id
					WHERE c.id = :class_id`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, true)
	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())
	var dbLearners []sqlc.GetLearnersByClassIdRow
	err = pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbLearners)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, nil
	}

	if dbLearners == nil {
		return nil, nil
	}
	var learners []Learner

	for _, dbLearner := range dbLearners {
		learner := Learner{
			ID:       dbLearner.ID,
			FullName: *dbLearner.FullName,
			Email:    dbLearner.Email,
			Phone:    *dbLearner.Phone,
			Gender:   dbLearner.Gender,
			Photo:    *dbLearner.ProfilePhoto,
			School: School{
				ID:   *dbLearner.SchoolID,
				Name: dbLearner.SchoolName,
			},
			ImageLink: dbLearner.Image,
		}

		dbAttendances, _ := c.queries.GetAttendanceByClassLearner(ctx, dbLearner.ClassLearnerID)
		learner.Attendances = toCoreAttendanceSlice(dbAttendances)

		dbAssignments, err := c.queries.GetAssignmentsByClassLearner(ctx, dbLearner.ClassLearnerID)
		if err != nil {
			return nil, nil
		}
		learner.Assignments = toCoreAssignmentSlice(dbAssignments)
		learners = append(learners, learner)
	}
	return learners, nil
}

func (c *Core) CountLearnersInClass(ctx *gin.Context, classId uuid.UUID, filter QueryFilter) int {
	if err := filter.Validate(); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	data := map[string]interface{}{
		"class_id": classId,
	}

	const q = `SELECT
                         COUNT(u.id) AS count
               FROM
                         users u
							JOIN class_learners cl ON u.id = cl.learner_id
 							JOIN classes c ON cl.class_id = c.id
							JOIN schools s ON s.id = u.school_id
								WHERE c.id = :class_id`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, true)

	var count struct {
		Count int `db:"count"`
	}

	if err := pgx.NamedQueryStruct(ctx, c.logger, c.db, buf.String(), data, &count); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	return count.Count
}

func (c *Core) GetLearnersAttendance(ctx *gin.Context, slotId uuid.UUID, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]AttendanceRecord, error) {
	_, err := middleware.AuthorizeWithoutLearner(ctx, c.queries)
	if err != nil {
		return nil, err
	}

	if err := filter.Validate(); err != nil {
		return nil, nil
	}

	data := map[string]interface{}{
		"slot_id":       slotId,
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT 
					u.id, u.full_name AS full_name, s.id AS school_id, s.name AS school_name, la.status
			FROM users u
    			JOIN class_learners cl ON u.id = cl.learner_id
    			JOIN schools s ON s.id = u.school_id
    			JOIN learner_attendances la ON la.class_learner_id = cl.id 
					WHERE la.slot_id = :slot_id`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, true)
	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())
	var dbAttendanceRecords []sqlc.GetLearnerAttendanceBySlotRow
	err = pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbAttendanceRecords)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, nil
	}

	if dbAttendanceRecords == nil {
		return nil, nil
	}
	var attendanceRecords []AttendanceRecord
	for _, dbAttendanceRecord := range dbAttendanceRecords {
		attendanceRecord := AttendanceRecord{
			ID:       dbAttendanceRecord.ID,
			FullName: *dbAttendanceRecord.FullName,
			School: School{
				ID:   dbAttendanceRecord.SchoolID,
				Name: dbAttendanceRecord.SchoolName,
			},
			Status: dbAttendanceRecord.Status,
		}
		attendanceRecords = append(attendanceRecords, attendanceRecord)
	}
	return attendanceRecords, nil
}

func (c *Core) CountLearnersAttendance(ctx *gin.Context, slotId uuid.UUID, filter QueryFilter) int {
	if err := filter.Validate(); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	data := map[string]interface{}{
		"slot_id": slotId,
	}

	const q = `SELECT
                         COUNT(u.id) AS count
			FROM users u
    			JOIN class_learners cl ON u.id = cl.learner_id
    			JOIN schools s ON s.id = u.school_id
    			JOIN learner_attendances la ON la.class_learner_id = cl.id
					WHERE la.slot_id = :slot_id`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, true)

	var count struct {
		Count int `db:"count"`
	}

	if err := pgx.NamedQueryStruct(ctx, c.logger, c.db, buf.String(), data, &count); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	return count.Count
}
