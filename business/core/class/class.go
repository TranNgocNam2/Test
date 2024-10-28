package class

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/middleware"
	"Backend/internal/weekday"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

var (
	ErrProgramOrSubjectNotFound = errors.New("Không tìm thấy chương trình học hoặc môn học!")
	ErrInvalidClassStartTime    = errors.New("Thời gian bắt đầu lớp học không hợp lệ!")
	ErrSessionNotFound          = errors.New("Không có buổi học nào trong môn học này!")
	ErrInvalidWeekDay           = errors.New("Số ngày học trong tuần không khớp với số buổi học trong môn học!")
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

func (c *Core) Create(ctx *gin.Context, newClass NewClass) (uuid.UUID, error) {
	staffID, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return uuid.Nil, err
	}

	programSubject, err := c.queries.GetProgramSubject(ctx, sqlc.GetProgramSubjectParams{
		ProgramID: newClass.ProgramID,
		SubjectID: newClass.SubjectID,
	})
	if err != nil {
		return uuid.Nil, ErrProgramOrSubjectNotFound
	}

	if newClass.StartDate.Before(programSubject.StartDate) || newClass.StartDate.After(programSubject.EndDate) {
		return uuid.Nil, ErrInvalidClassStartTime
	}

	if int(programSubject.TimePerSession) != len(newClass.Slots.WeekDays) {
		return uuid.Nil, ErrInvalidWeekDay
	}

	sessions, err := c.queries.GetSessionBySubject(ctx, newClass.SubjectID)
	if err != nil {
		return uuid.Nil, ErrSessionNotFound
	}

	slots := generateSlots(newClass, sessions, programSubject.TimePerSession, programSubject.EndDate)

	// Check if the last slot's end time is after the programs end time
	var endDateClass *time.Time
	lastSlot := slots[len(slots)-1:][0].EndTime
	if lastSlot != nil && lastSlot.Before(programSubject.EndDate) {
		endDate := time.Date(lastSlot.Year(), lastSlot.Month(), lastSlot.Day(),
			0, 0, 0, 0, time.Local)
		endDateClass = &endDate
	}

	dbClass := sqlc.CreateClassParams{
		ID:               newClass.ID,
		ProgramSubjectID: programSubject.ID,
		Name:             newClass.Name,
		Code:             newClass.Code,
		Link:             newClass.Link,
		StartDate:        newClass.StartDate,
		EndDate:          endDateClass,
		Password:         newClass.Password,
		CreatedBy:        staffID,
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)
	if err = qtx.CreateClass(ctx, dbClass); err != nil {
		return uuid.Nil, err
	}

	if _, err = qtx.CreateSlots(ctx, slots); err != nil {
		return uuid.Nil, err
	}
	tx.Commit(ctx)
	return dbClass.ID, nil
}

func (c *Core) GetByID(ctx *gin.Context, id uuid.UUID) (Details, error) {
	dbClass, err := c.queries.GetClassByID(ctx, id)
	if err != nil {
		return Details{}, err
	}
	class := Details{
		ID:        dbClass.ID,
		Name:      dbClass.Name,
		Link:      *dbClass.Link,
		StartDate: dbClass.StartDate,
		EndDate:   dbClass.EndDate,
		Teachers:  nil,
		Slots:     nil,
	}

	programSubject, err := c.queries.GetProgramSubjectByID(ctx, dbClass.ProgramSubjectID)
	if err != nil {
		return Details{}, err
	}
	class.ProgramID = programSubject.ProgramID

	dbSubject, err := c.queries.GetSubjectById(ctx, programSubject.SubjectID)
	if err != nil {
		return Details{}, err
	}
	class.Subject = toCoreSubject(dbSubject)

	dbTeachers, err := c.queries.GetTeachersByClassID(ctx, dbClass.ID)
	if err != nil {
		return Details{}, nil
	}

	class.Teachers = toCoreTeacherSlice(dbTeachers)

	var slots []Slot
	dbSlots, err := c.queries.GetSlotsByClassID(ctx, dbClass.ID)
	if err != nil {
		return Details{}, err
	}
	for _, dbSlot := range dbSlots {
		dbSession, err := c.queries.GetSessionByID(ctx, dbSlot.SessionID)
		if err != nil {
			return Details{}, err
		}

		session := toCoreSession(dbSession)
		dbTeacher, err := c.queries.GetTeacherBySessionID(ctx, dbSession.ID)
		if err != nil {
			return Details{}, err
		}
		slot := Slot{
			ID:        dbSlot.ID,
			StartTime: *dbSlot.StartTime,
			EndTime:   *dbSlot.EndTime,
			Session:   session,
		}
		slots = append(slots, slot)
	}

	return class, nil
}

func generateSlots(newClass NewClass, sessions []sqlc.Session, duration int16, endDate time.Time) []sqlc.CreateSlotsParams {
	var slots []sqlc.CreateSlotsParams

	currentDate := newClass.Slots.StartDate

	for i, session := range sessions {
		if newClass.Slots.StartDate == nil || newClass.Slots.StartTime == nil ||
			len(newClass.Slots.WeekDays) == 0 {
			slot := sqlc.CreateSlotsParams{
				ID:        uuid.New(),
				SessionID: session.ID,
				ClassID:   newClass.ID,
				StartTime: nil,
				EndTime:   nil,
				Index:     session.Index,
			}
			slots = append(slots, slot)
			return slots
		}
		weekDay := newClass.Slots.WeekDays[i%len(newClass.Slots.WeekDays)]
		weeksToAdd := i / len(newClass.Slots.WeekDays)

		slotDate := weekday.Next(currentDate.AddDate(0, 0, 7*weeksToAdd), weekDay)

		slotStartTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(),
			newClass.Slots.StartTime.Hour(), newClass.Slots.StartTime.Minute(), 0, 0, time.Local)
		slotEndTime := slotStartTime.Add(time.Hour * time.Duration(duration))

		startTime := &slotStartTime
		endTime := &slotEndTime

		if slotStartTime.After(endDate) || slotEndTime.After(endDate) {
			startTime = nil
			endTime = nil
		}

		slot := sqlc.CreateSlotsParams{
			ID:        uuid.New(),
			SessionID: session.ID,
			ClassID:   newClass.ID,
			StartTime: startTime,
			EndTime:   endTime,
			Index:     session.Index,
		}

		slots = append(slots, slot)
	}

	return slots
}
