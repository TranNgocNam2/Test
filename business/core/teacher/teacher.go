package teacher

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/code"
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

func (c *Core) GenerateAttendanceCode(ctx *gin.Context, slotId uuid.UUID) (string, error) {
	teacherId, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		return "", err
	}

	slot, err := c.queries.GetSlotById(ctx, slotId)
	if err != nil {
		return "", model.ErrSlotNotFound
	}

	if strings.Compare(*slot.TeacherID, teacherId) != 0 {
		return "", model.ErrTeacherIsNotInSlot
	}

	if slot.StartTime.After(time.Now()) {
		return "", model.ErrSlotNotStarted
	}

	if slot.EndTime.Before(time.Now()) {
		return "", model.ErrSlotEnded
	}
	if slot.AttendanceCode != nil {
		return *slot.AttendanceCode, nil
	}

	attendanceCode := code.GenerateAttendance(6)

	err = c.queries.UpdateAttendanceCode(ctx, sqlc.UpdateAttendanceCodeParams{
		AttendanceCode: &attendanceCode,
		ID:             slot.ID,
	})
	if err != nil {
		c.logger.Error(err.Error())
		return "", err
	}

	return attendanceCode, nil
}

func (c *Core) GetTeachersInClass(ctx *gin.Context, classId uuid.UUID) ([]Teacher, error) {
	_, err := middleware.AuthorizeUser(ctx, c.queries)
	if err != nil {
		return nil, err
	}

	dbClass, err := c.queries.GetClassById(ctx, classId)
	if err != nil {
		return nil, model.ErrClassNotFound
	}

	dbTeachers, err := c.queries.GetTeachersInClass(ctx, dbClass.ID)
	if err != nil {
		return nil, model.ErrTeacherNotFound
	}

	var teachers []Teacher
	for _, dbTeacher := range dbTeachers {
		teachers = append(teachers, toCoreTeacher(dbTeacher))
	}

	return teachers, nil
}

func (c *Core) UpdateRecordLink(ctx *gin.Context, slotId uuid.UUID, recordLink UpdateRecord) error {
	teacherID, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		return err
	}

	slot, err := c.queries.GetSlotById(ctx, slotId)
	if err != nil {
		return model.ErrSlotNotFound
	}

	if strings.Compare(*slot.TeacherID, teacherID) != 0 {
		return model.ErrTeacherIsNotInSlot
	}
	if slot.StartTime.After(time.Now()) {
		return model.ErrStartTimeNotStarted
	}

	err = c.queries.UpdateRecordLink(ctx, sqlc.UpdateRecordLinkParams{
		RecordLink: &recordLink.Link,
		ID:         slotId,
	})
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	return nil
}
