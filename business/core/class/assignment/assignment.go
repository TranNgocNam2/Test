package assignment

import (
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/web/payload"
	"bytes"
	"encoding/json"
	"time"

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

func (c *Core) CreateAssignment(ctx *gin.Context, classId uuid.UUID, asm payload.Assignment) (string, error) {
	_, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {

		c.logger.Error(err.Error())
		return "", middleware.ErrInvalidUser
	}

	class, err := c.queries.GetClassById(ctx, classId)
	if err != nil {
		c.logger.Error(err.Error())
		return "", model.ErrClassNotFound
	}

	dbProgram, _ := c.queries.GetProgramById(ctx, class.ProgramID)

	deadline, err := time.Parse(time.DateTime, asm.Deadline)
	if err != nil {
		c.logger.Error(err.Error())
		return "", model.ErrTimeFormat
	}

	if deadline.Before(dbProgram.StartDate) {
		return "", model.ErrInvalidDeadlineTime
	}

	if deadline.After(dbProgram.EndDate) {
		return "", model.ErrInvalidDeadlineTime
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	question, err := json.Marshal(asm.Question)
	if err != nil {
		return "", model.ErrDataConversion
	}

	asmId, err := qtx.InsertAssignment(ctx, sqlc.InsertAssignmentParams{
		ID:         uuid.New(),
		Classid:    classId,
		Question:   json.RawMessage(question),
		Status:     int16(*asm.Status),
		Type:       int16(*asm.Type),
		CanOverdue: &asm.CanOverdue,
		Deadline:   &deadline,
	})

	if err != nil {
		return "", err
	}

	var params []sqlc.InsertLearnerAssignmentParams
	classLearners, err := qtx.GetLearnersByClassId(ctx, classId)
	if err != nil {
		return "", err
	}

	for _, learner := range classLearners {
		param := sqlc.InsertLearnerAssignmentParams{
			ID:               uuid.New(),
			ClassLearnerID:   learner.ClassLearnerID,
			AssignmentID:     asmId,
			Grade:            0,
			GradingStatus:    NOT_GRADED,
			SubmissionStatus: NOT_SUBMITTED,
		}

		params = append(params, param)
	}

	if _, err = qtx.InsertLearnerAssignment(ctx, params); err != nil {
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	return asmId.String(), nil
}

func (c *Core) UpdateAssignment(ctx *gin.Context, classId uuid.UUID, asmId uuid.UUID, data payload.Assignment) error {
	_, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		c.logger.Error(err.Error())
		return middleware.ErrInvalidUser
	}

	class, err := c.queries.GetClassById(ctx, classId)
	if err != nil {
		return model.ErrClassNotFound
	}

	if res, err := c.queries.CheckAssignmentInClass(ctx, sqlc.CheckAssignmentInClassParams{
		ClassID: classId,
		ID:      asmId,
	}); err != nil || !res {
		c.logger.Error(err.Error())
		return model.InvalidClassAssignment
	}

	dbProgram, _ := c.queries.GetProgramById(ctx, class.ProgramID)

	deadline, err := time.Parse(time.DateTime, data.Deadline)
	if err != nil {
		return model.ErrTimeFormat
	}

	if deadline.Before(dbProgram.StartDate) {
		return model.ErrInvalidDeadlineTime
	}

	if deadline.After(dbProgram.EndDate) {
		return model.ErrInvalidDeadlineTime
	}

	question, err := json.Marshal(data.Question)
	if err != nil {
		return model.ErrDataConversion
	}

	if err := c.queries.UpdateAssignment(ctx, sqlc.UpdateAssignmentParams{
		Question:   json.RawMessage(question),
		Deadline:   &deadline,
		Status:     int16(*data.Status),
		Type:       int16(*data.Type),
		CanOverdue: &data.CanOverdue,
		ID:         asmId,
	}); err != nil {
		return err
	}

	return nil
}

func (c *Core) DeleteAssignment(ctx *gin.Context, asmId uuid.UUID) error {
	_, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		c.logger.Error(err.Error())
		return middleware.ErrInvalidUser
	}

	asm, err := c.queries.GetAssignmentById(ctx, asmId)
	if err != nil {
		c.logger.Error(err.Error())
		return model.ErrAssignmentNotFound
	}

	if asm.Status != VISIBLE {
		return model.ErrAssignmentDeletion
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	if err := qtx.DeleleLearnerAssignment(ctx, asmId); err != nil {
		return err
	}

	if err := qtx.DeleteAssignment(ctx, asmId); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Core) GetById(ctx *gin.Context, id uuid.UUID) (*Assignment, error) {
	asm, err := c.queries.GetAssignmentById(ctx, id)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, model.ErrAssignmentNotFound
	}

	result := Assignment{
		Id:         asm.ID,
		ClassId:    asm.ClassID,
		Question:   asm.Question,
		Deadline:   *asm.Deadline,
		Status:     int(asm.Status),
		Type:       int(asm.Type),
		CanOverdue: *asm.CanOverdue,
	}

	return &result, nil
}

func (c *Core) Query(ctx *gin.Context, classId uuid.UUID, orderBy order.By, pageNumber int, rowsPerPage int) []Assignment {

	data := map[string]interface{}{
		"class_id":      classId,
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT
                        id, class_id, question, deadline, status, type, can_overdue
               FROM assignments`

	buf := bytes.NewBufferString(q)
	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var dbAssignments []sqlc.Assignment
	if err := pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbAssignments); err != nil {
		c.logger.Error(err.Error())
		return nil
	}

	if dbAssignments == nil {
		return nil
	}

	var assignments []Assignment
	for _, asm := range dbAssignments {
		assignment := Assignment{
			Id:         asm.ID,
			ClassId:    asm.ClassID,
			Question:   asm.Question,
			Deadline:   *asm.Deadline,
			Status:     int(asm.Status),
			Type:       int(asm.Type),
			CanOverdue: *asm.CanOverdue,
		}

		assignments = append(assignments, assignment)
	}
	return assignments
}

func (c *Core) Count(ctx *gin.Context, classId uuid.UUID) int {
	data := map[string]interface{}{
		"class_id": classId,
	}

	const q = `SELECT COUNT(1) FROM assignments`
	buf := bytes.NewBufferString(q)
	var count struct {
		Count int `db:"count"`
	}

	if err := pgx.NamedQueryStruct(ctx, c.logger, c.db, buf.String(), data, &count); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	return count.Count
}

func (c *Core) GradeAssignment(ctx *gin.Context, learnerId uuid.UUID, asmId uuid.UUID, data payload.AssignmentGrade) error {
	_, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		c.logger.Error(err.Error())
		return middleware.ErrInvalidUser
	}

	asm, err := c.queries.GetAssignmentById(ctx, asmId)
	if err != nil {
		c.logger.Error(err.Error())
		return model.ErrAssignmentNotFound
	}

	if asm.Status == VISIBLE {
		return model.ErrCannotGradeVisibleAssignment
	}

	classLearner, err := c.queries.GetClassLearnerByClassAndLearner(ctx, sqlc.GetClassLearnerByClassAndLearnerParams{
		ClassID:   asm.ClassID,
		LearnerID: learnerId.String(),
	})

	if err != nil {
		c.logger.Error(err.Error())
		return model.LearnerNotInClass
	}

	learnerAsm, err := c.queries.GetLearnerAssignment(ctx, sqlc.GetLearnerAssignmentParams{
		ClassLearnerID: classLearner.ID,
		AssignmentID:   asm.ID,
	})

	if err != nil {
		c.logger.Error(err.Error())
		return model.ErrLearnerAssignmentNotFound
	}

	if learnerAsm.SubmissionStatus == NOT_SUBMITTED {
		return model.ErrGradingNotStartedAssignment
	}

	if err := c.queries.UpdateLearnerGrade(ctx, sqlc.UpdateLearnerGradeParams{
		ClassLearnerID: classLearner.ID,
		AssignmentID:   asm.ID,
		Grade:          data.Grade,
		GradingStatus:  GRADED,
	}); err != nil {
		return err
	}

	return nil
}

func (c *Core) GetLearnerAssignment(ctx *gin.Context, asmId uuid.UUID) (*LearnerAssignment, error) {
	learner, err := middleware.AuthorizeVerifiedLearner(ctx, c.queries)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, middleware.ErrInvalidUser
	}

	asm, err := c.queries.GetAssignmentById(ctx, asmId)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, model.ErrAssignmentNotFound
	}

	assignment := Assignment{
		Id:         asm.ID,
		ClassId:    asm.ClassID,
		Question:   asm.Question,
		Deadline:   *asm.Deadline,
		Status:     int(asm.Status),
		Type:       int(asm.Type),
		CanOverdue: *asm.CanOverdue,
	}

	classLearner, err := c.queries.GetClassLearnerByClassAndLearner(ctx, sqlc.GetClassLearnerByClassAndLearnerParams{
		ClassID:   asm.ClassID,
		LearnerID: learner.ID,
	})

	if err != nil {
		c.logger.Error(err.Error())
		return nil, model.LearnerNotInClass
	}

	learnerAsm, err := c.queries.GetLearnerAssignment(ctx, sqlc.GetLearnerAssignmentParams{
		ClassLearnerID: classLearner.ID,
		AssignmentID:   asm.ID,
	})

	if err != nil {
		c.logger.Error(err.Error())
		return nil, model.ErrLearnerAssignmentNotFound
	}

	learnerAssignment := LearnerAssignment{
		LearnerId:        learner.ID,
		Grade:            learnerAsm.Grade,
		SubmissionStatus: int(learnerAsm.SubmissionStatus),
		GradingStatus:    int(learnerAsm.GradingStatus),
		Data:             learnerAsm.Data,
		Assignment:       assignment,
	}

	return &learnerAssignment, nil
}

func (c *Core) SubmitAssignment(ctx *gin.Context, asmId uuid.UUID, req payload.LearnerSubmission) error {
	learner, err := middleware.AuthorizeVerifiedLearner(ctx, c.queries)
	if err != nil {
		c.logger.Error(err.Error())
		return middleware.ErrInvalidUser
	}

	asm, err := c.queries.GetAssignmentById(ctx, asmId)
	if err != nil {
		c.logger.Error(err.Error())
		return model.ErrAssignmentNotFound
	}

	if asm.Status == VISIBLE || asm.Status == CLOSED {
		return model.ErrInvalidAssignmentSubmision
	}

	now := time.Now().UTC()

	if now.After(*asm.Deadline) && !*asm.CanOverdue {
		return model.ErrSubmitOverdue
	}

	classLearner, err := c.queries.GetClassLearnerByClassAndLearner(ctx, sqlc.GetClassLearnerByClassAndLearnerParams{
		ClassID:   asm.ClassID,
		LearnerID: learner.ID,
	})

	if err != nil {
		c.logger.Error(err.Error())
		return model.LearnerNotInClass
	}

	_, err = c.queries.GetLearnerAssignment(ctx, sqlc.GetLearnerAssignmentParams{
		ClassLearnerID: classLearner.ID,
		AssignmentID:   asm.ID,
	})

	if err != nil {
		c.logger.Error(err.Error())
		return model.ErrLearnerAssignmentNotFound
	}

	data, err := json.Marshal(req.Data)
	if err != nil {
		return model.ErrDataConversion
	}

	submissionTime := time.Now().UTC()

	submissionStatus := SUBMITTED
	if submissionTime.After(*asm.Deadline) {
		submissionStatus = LATE
	}

	if err := c.queries.UpdateLearnerAssignment(ctx, sqlc.UpdateLearnerAssignmentParams{
		SubmissionStatus: int16(submissionStatus),
		Data:             json.RawMessage(data),
		AssignmentID:     asm.ID,
		ClassLearnerID:   classLearner.ID,
	}); err != nil {
		return err
	}

	return nil
}
