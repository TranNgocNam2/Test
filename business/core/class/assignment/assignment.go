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
	"fmt"
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
		return "", err
	}

	class, err := c.queries.GetClassById(ctx, classId)
	if err != nil {
		return "", model.ErrClassNotFound
	}

	dbProgram, _ := c.queries.GetProgramById(ctx, class.ProgramID)

	deadline, err := time.Parse(time.DateOnly, asm.Deadline)
	if err != nil {
		return "", err
	}

	if deadline.Before(dbProgram.StartDate) {
		return "", model.ErrInvalidSlotStartTime
	}

	if deadline.After(dbProgram.EndDate) {
		return "", model.ErrInvalidSlotEndTime
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
	tx.Commit(ctx)

	return asmId.String(), nil
}

func (c *Core) UpdateAssignment(ctx *gin.Context, classId uuid.UUID, asmId uuid.UUID, data payload.Assignment) error {
	_, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		return err
	}

	if res, err := c.queries.CheckAssignmentInClass(ctx, sqlc.CheckAssignmentInClassParams{
		ClassID: classId,
		ID:      asmId,
	}); err != nil || !res {
		return err
	}

	class, err := c.queries.GetClassById(ctx, classId)
	if err != nil {
		return model.ErrClassNotFound
	}

	dbProgram, _ := c.queries.GetProgramById(ctx, class.ProgramID)

	deadline, err := time.Parse(time.DateOnly, data.Deadline)
	if err != nil {
		return err
	}

	if deadline.Before(dbProgram.StartDate) {
		return model.ErrInvalidSlotStartTime
	}

	if deadline.After(dbProgram.EndDate) {
		return model.ErrInvalidSlotEndTime
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
		return err
	}

	asm, err := c.queries.GetAssignmentById(ctx, asmId)
	if err != nil {
		return err
	}

	if asm.Status != VISIBLE {
		return fmt.Errorf("no no no")
	}

	if err := c.queries.DeleteAssignment(ctx, asmId); err != nil {
		return err
	}

	return nil
}

func (c *Core) GetById(ctx *gin.Context, id uuid.UUID) (*Assignment, error) {
	asm, err := c.queries.GetAssignmentById(ctx, id)
	if err != nil {
		return nil, err
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

func (c *Core) GradeAssignment(ctx *gin.Context, learnerId uuid.UUID, asmId uuid.UUID, data payload.AssignmentGrade) error {
	_, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		return err
	}

	asm, err := c.queries.GetAssignmentById(ctx, asmId)
	if err != nil {
		return err
	}

	classLearner, err := c.queries.GetLearnerByClassId(ctx, sqlc.GetLearnerByClassIdParams{
		ClassID:   asm.ClassID,
		LearnerID: learnerId.String(),
	})

	if err != nil {
		return err
	}

	return nil
}
