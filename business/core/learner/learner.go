package learner

import (
	"Backend/business/core/learner/certificate"
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/common/status"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/slice"
	"bytes"
	"fmt"
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

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		c.logger.Error(err.Error())
		return model.ErrFailedToAddLearnerToClass
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	dbClass, err := qtx.GetClassCompletedByCode(ctx, classAccess.Code)
	if err != nil {
		return model.ErrClassNotFound
	}

	_, err = c.queries.GetClassLearnerByClassAndLearner(ctx,
		sqlc.GetClassLearnerByClassAndLearnerParams{
			ClassID:   dbClass.ID,
			LearnerID: learner.ID,
		})
	if err == nil {
		return model.ErrLearnerAlreadyInClass
	}

	if dbClass.StartDate.Before(time.Now()) {
		return model.ErrClassStarted
	}

	if strings.Compare(dbClass.Password, classAccess.Password) != 0 {
		return model.ErrWrongPassword
	}

	dbSlots, _ := qtx.GetSlotsByClassId(ctx, dbClass.ID)
	var slotIds []uuid.UUID
	for _, dbSlot := range dbSlots {
		scheduleConflict, _ := c.queries.CheckLearnerTimeOverlap(ctx,
			sqlc.CheckLearnerTimeOverlapParams{
				LearnerID: learner.ID,
				EndTime:   dbSlot.EndTime,
				StartTime: dbSlot.StartTime,
			})
		if scheduleConflict {
			return fmt.Errorf(model.ErrScheduleConflict, dbSlot.StartTime.Format("15:04 02/01/2006"),
				dbSlot.EndTime.Format("15:04 02/01/2006"))
		}
		slotIds = append(slotIds, dbSlot.ID)
	}

	classLearner := sqlc.AddLearnerToClassParams{
		ID:        uuid.New(),
		ClassID:   dbClass.ID,
		LearnerID: learner.ID,
	}
	err = qtx.AddLearnerToClass(ctx, classLearner)
	if err != nil {
		c.logger.Error(err.Error())
		return model.ErrFailedToAddLearnerToClass
	}

	err = qtx.GenerateLearnerAttendance(ctx, sqlc.GenerateLearnerAttendanceParams{
		ClassLearnerID: classLearner.ID,
		SlotIds:        slotIds,
	})
	if err != nil {
		c.logger.Error(err.Error())
		return model.ErrFailedToAddLearnerToClass
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

	classLearner, err := c.queries.GetClassLearnerByClassAndLearner(ctx,
		sqlc.GetClassLearnerByClassAndLearnerParams{
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

	currentTime := time.Now()

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
			Status: int32(status.Attended),
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
		"status":        int16(status.Valid),
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT
				u.id, u.full_name AS full_name, u.email, u.phone, u.profile_photo, u.status AS status, u.type,
				s.id AS school_id, s.name AS school_name, cl.id AS class_learner_id
			FROM users u
				JOIN class_learners cl ON u.id = cl.learner_id
				JOIN classes c ON cl.class_id = c.id
        		JOIN schools s ON s.id = u.school_id
					WHERE c.id = :class_id AND u.is_verified = true AND u.status = :status`

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
			Phone:    dbLearner.Phone,
			Photo:    dbLearner.ProfilePhoto,
			Type:     dbLearner.Type,
			School: School{
				ID:   dbLearner.SchoolID,
				Name: dbLearner.SchoolName,
			},
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
		"status":   int16(status.Valid),
	}

	const q = `SELECT
                         COUNT(u.id) AS count
               FROM
                         users u
							JOIN class_learners cl ON u.id = cl.learner_id
 							JOIN classes c ON cl.class_id = c.id
							JOIN schools s ON s.id = u.school_id
								WHERE c.id = :class_id AND u.is_verified = true AND u.status = :status`

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
		"status":        int16(status.Valid),
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT 
					u.id, u.full_name AS full_name, s.id AS school_id, s.name AS school_name, la.status
			FROM users u
    			JOIN class_learners cl ON u.id = cl.learner_id
    			JOIN schools s ON s.id = u.school_id
    			JOIN learner_attendances la ON la.class_learner_id = cl.id 
					WHERE la.slot_id = :slot_id AND u.is_verified = true AND u.status = :status`

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
		"status":  int16(status.Valid),
	}

	const q = `SELECT
                         COUNT(u.id) AS count
			FROM users u
    			JOIN class_learners cl ON u.id = cl.learner_id
    			JOIN schools s ON s.id = u.school_id
    			JOIN learner_attendances la ON la.class_learner_id = cl.id
					WHERE la.slot_id = :slot_id AND u.is_verified = true AND u.status = :status`

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

func (c *Core) CreateVerificationInfo(ctx *gin.Context, updateLearner UpdateLearner) (uuid.UUID, error) {
	learner, err := middleware.AuthorizeLearner(ctx, c.queries)
	if err != nil {
		return uuid.Nil, err
	}

	if learner.IsVerified {
		return uuid.Nil, model.ErrLearnerAlreadyVerified
	}

	school, err := c.queries.GetSchoolById(ctx, updateLearner.SchoolId)
	if err != nil {
		return uuid.Nil, model.ErrSchoolNotFound
	}

	verification, err := c.queries.GetLearnerVerificationByLearnerId(ctx,
		sqlc.GetLearnerVerificationByLearnerIdParams{
			LearnerID: learner.ID,
			Status:    int16(status.Pending),
		})
	if err == nil && status.Verification(verification.Status) == status.Pending {
		return uuid.Nil, model.ErrVerificationPending
	}

	verificationId, err := c.queries.CreateVerificationRequest(ctx, sqlc.CreateVerificationRequestParams{
		ImageLink: updateLearner.ImageLinks,
		Type:      updateLearner.Type,
		SchoolID:  school.ID,
		LearnerID: learner.ID,
		Status:    int16(status.Pending),
		ID:        uuid.New(),
	})
	if err != nil {
		return uuid.Nil, err
	}
	return verificationId, nil
}

func (c *Core) CancelVerification(ctx *gin.Context, verificationId uuid.UUID) error {
	learner, err := middleware.AuthorizeLearner(ctx, c.queries)
	if err != nil {
		return err
	}
	if learner.IsVerified {
		return model.ErrLearnerAlreadyVerified
	}

	verification, err := c.queries.GetLearnerVerificationById(ctx, verificationId)
	if verification.LearnerID != learner.ID {
		return model.ErrUnauthorizedFeatureAccess
	}
	if err != nil && verification.Status != int16(status.Pending) {
		return model.ErrVerificationNotFound
	}

	err = c.queries.VerifyLearner(ctx, sqlc.VerifyLearnerParams{
		VerifiedBy: nil,
		Status:     int16(status.Cancelled),
		Note:       "Học viên đã hủy yêu cầu xác thực!",
		LearnerID:  learner.ID,
		ID:         verification.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) GetVerificationsInformation(ctx *gin.Context) (*VerifyLearnerInfo, error) {
	learner, err := middleware.AuthorizeLearner(ctx, c.queries)
	if err != nil {
		return nil, err
	}

	dbVerifications, err := c.queries.GetVerificationLearners(ctx, learner.ID)
	if err != nil || dbVerifications == nil {
		return nil, nil
	}

	learnerVerification := VerifyLearnerInfo{
		ID:            learner.ID,
		FullName:      *learner.FullName,
		Email:         learner.Email,
		Verifications: nil,
	}
	for _, dbVerification := range dbVerifications {
		verification := struct {
			ID        uuid.UUID `json:"id"`
			Status    int16     `json:"status"`
			Note      *string   `json:"note"`
			ImageLink []string  `json:"imageLink"`
			Type      int16     `json:"type"`
			School    School    `json:"school"`
			CreatedAt time.Time `json:"createdAt"`
		}{
			ID:        dbVerification.ID,
			Status:    dbVerification.Status,
			Note:      dbVerification.Note,
			ImageLink: slice.ParseFromString(dbVerification.ImageLink),
			Type:      dbVerification.Type,
			School: School{
				ID:   dbVerification.SchoolID,
				Name: dbVerification.SchoolName,
			},
			CreatedAt: dbVerification.CreatedAt,
		}
		learnerVerification.Verifications = append(learnerVerification.Verifications, verification)
	}

	return &learnerVerification, nil
}
