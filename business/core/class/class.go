package class

import (
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/weekday"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

var (
	ErrProgramNotFound       = errors.New("Không tìm thấy chương trình học!")
	ErrSubjectNotFound       = errors.New("Không tìm thấy môn học!")
	ErrInvalidClassStartTime = errors.New("Thời gian bắt đầu lớp học không hợp lệ!")
	ErrInvalidSlotStartTime  = errors.New("Thời gian bắt đầu buổi học không hợp lệ!")
	ErrInvalidSlotEndTime    = errors.New("Thời gian kết thúc buổi học không hợp lệ!")
	ErrSessionNotFound       = errors.New("Không có buổi học nào trong môn học này!")
	ErrInvalidWeekDay        = errors.New("Số ngày học trong tuần không khớp với số buổi học trong môn học!")
	ErrClassNotFound         = errors.New("Không tìm thấy lớp học!")
	ErrClassCodeAlreadyExist = errors.New("Mã của lớp học đã tồn tại!")
	ErrTeacherNotFound       = errors.New("Không tìm thấy giáo viên!")
	ErrInvalidSlotCount      = errors.New("Số lượng buổi học không hợp lệ!")
	ErrInvalidSlotTime       = errors.New("Thời gian buổi học không hợp lệ!")
	ErrSlotNotFound          = errors.New("Không tìm thấy buổi học!")
	ErrTeacherNotAvailable   = errors.New("Giáo viên không thể dạy vào thời gian này!")
	ErrTeacherIsNotInClass   = errors.New("Giáo viên không thuộc lớp học này!")
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
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = c.queries.GetClassByCode(ctx, newClass.Code)
	if err == nil {
		return uuid.Nil, ErrClassCodeAlreadyExist
	}

	dbProgram, err := c.queries.GetProgramByID(ctx, newClass.ProgramId)
	if err != nil {
		return uuid.Nil, ErrProgramNotFound
	}

	dbSubject, err := c.queries.GetPublishedSubjectByID(ctx, newClass.SubjectId)
	if err != nil {
		return uuid.Nil, ErrSubjectNotFound
	}

	sessions, err := c.queries.GetSessionsBySubjectId(ctx, newClass.SubjectId)
	if err != nil {
		return uuid.Nil, ErrSessionNotFound
	}

	slots := generateSlots(newClass, sessions, dbSubject.TimePerSession, dbProgram.EndDate)

	// Check if the last slot's end time is after the programs end time
	var startDateClass *time.Time
	firstSlot := slots[0].StartTime
	if firstSlot != nil && firstSlot.Before(dbProgram.StartDate) {
		return uuid.Nil, ErrInvalidClassStartTime
	}

	if firstSlot != nil && firstSlot.After(dbProgram.StartDate) {
		startDate := time.Date(firstSlot.Year(), firstSlot.Month(), firstSlot.Day(),
			0, 0, 0, 0, time.Local)
		startDateClass = &startDate
	}

	var endDateClass *time.Time
	lastSlot := slots[len(slots)-1:][0].EndTime
	if lastSlot != nil && lastSlot.Before(dbProgram.EndDate) {
		endDate := time.Date(lastSlot.Year(), lastSlot.Month(), lastSlot.Day(),
			0, 0, 0, 0, time.Local)
		endDateClass = &endDate
	}

	dbClass := sqlc.CreateClassParams{
		ID:        newClass.ID,
		SubjectID: dbSubject.ID,
		ProgramID: dbProgram.ID,
		Name:      newClass.Name,
		Code:      newClass.Code,
		Link:      newClass.Link,
		StartDate: startDateClass,
		EndDate:   endDateClass,
		Password:  newClass.Password,
		CreatedBy: staffId,
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

func (c *Core) QueryByManager(ctx *gin.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) []Class {
	if err := filter.Validate(); err != nil {
		return nil
	}

	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT
						id, name, code, subject_id, program_id, link, start_date, end_date, status
			FROM classes`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var dbClasses []sqlc.Class
	err := pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbClasses)
	if err != nil {
		c.logger.Error(err.Error())
		return nil
	}

	if dbClasses == nil {
		return nil
	}

	var classes []Class

	for _, dbClass := range dbClasses {
		class := Class{
			ID:   dbClass.ID,
			Name: dbClass.Name,
			Code: dbClass.Code,
		}

		dbProgram, _ := c.queries.GetProgramByID(ctx, dbClass.ProgramID)
		class.Program = toCoreProgram(dbProgram)

		dbSubject, _ := c.queries.GetSubjectById(ctx, dbClass.SubjectID)
		class.Subject = toCoreSubject(dbSubject)

		dbTeachers, err := c.queries.GetTeachersByClassId(ctx, dbClass.ID)
		if err != nil {
			return nil
		}
		class.Teachers = toCoreTeacherSlice(dbTeachers)

		dbSkills, err := c.queries.GetSkillsBySubjectId(ctx, dbSubject.ID)
		if err != nil {
			return nil
		}
		class.Skills = toCoreSkillSlice(dbSkills)

		totalLearners, err := c.queries.CountLearnersByClassId(ctx, dbClass.ID)
		if err != nil {
			return nil
		}
		class.TotalLearners = totalLearners

		classes = append(classes, class)
	}

	return classes
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
                        classes`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}

	if err := pgx.NamedQueryStruct(ctx, c.logger, c.db, buf.String(), data, &count); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	return count.Count
}

func (c *Core) GetByID(ctx *gin.Context, id uuid.UUID) (Details, error) {
	dbClass, err := c.queries.GetClassById(ctx, id)
	if err != nil {
		return Details{}, ErrClassNotFound
	}

	class := Details{
		ID:        dbClass.ID,
		Name:      dbClass.Name,
		Link:      *dbClass.Link,
		StartDate: dbClass.StartDate,
		EndDate:   dbClass.EndDate,
	}

	dbSubject, err := c.queries.GetSubjectById(ctx, dbClass.SubjectID)
	if err != nil {
		return Details{}, ErrSubjectNotFound
	}
	class.Subject = toCoreSubject(dbSubject)

	dbProgram, err := c.queries.GetProgramByID(ctx, dbClass.ProgramID)
	if err != nil {
		return Details{}, ErrProgramNotFound
	}
	class.Program = toCoreProgram(dbProgram)

	dbTeachers, err := c.queries.GetTeachersByClassId(ctx, dbClass.ID)
	if err == nil {
		class.Teachers = toCoreTeacherSlice(dbTeachers)
	}

	var slots []Slot
	dbSlots, _ := c.queries.GetSlotsByClassID(ctx, dbClass.ID)
	for _, dbSlot := range dbSlots {
		dbSession, _ := c.queries.GetSessionById(ctx, dbSlot.SessionID)
		session := toCoreSession(dbSession)
		startTime := *dbSlot.StartTime
		endTime := *dbSlot.EndTime
		slot := Slot{
			ID:        dbSlot.ID,
			StartTime: startTime.Format(time.DateTime),
			EndTime:   endTime.Format(time.DateTime),
			Index:     dbSlot.Index,
			Session:   session,
		}

		if dbSlot.TeacherID != nil {
			dbTeacher, _ := c.queries.GetTeacherByID(ctx, *dbSlot.TeacherID)
			slot.Teacher = toCoreTeacher(dbTeacher)
		}

		slots = append(slots, slot)
	}
	class.Slots = slots

	return class, nil
}

func (c *Core) UpdateClassTeacher(ctx *gin.Context, id uuid.UUID, teacherIds []string) error {
	staffID, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbClass, err := c.queries.GetClassById(ctx, id)
	if err != nil {
		return ErrClassNotFound
	}

	for _, teacherId := range teacherIds {
		_, err = c.queries.GetTeacherByID(ctx, teacherId)
		if err != nil {
			return ErrTeacherNotFound
		}
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := c.queries.WithTx(tx)

	classTeacher := sqlc.AddTeacherToClassParams{
		TeacherIds: teacherIds,
		ClassID:    dbClass.ID,
		CreatedBy:  staffID,
	}

	if err = qtx.RemoveTeacherFromClass(ctx, dbClass.ID); err != nil {
		return err
	}

	if err = qtx.AddTeacherToClass(ctx, classTeacher); err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}

func (c *Core) UpdateSlot(ctx *gin.Context, id uuid.UUID, updateSlots []UpdateSlot) error {
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbClass, err := c.queries.GetClassById(ctx, id)
	if err != nil {
		return ErrClassNotFound
	}

	dbProgram, _ := c.queries.GetProgramByID(ctx, dbClass.ProgramID)

	firstSlot := updateSlots[0].StartTime
	if firstSlot.Before(dbProgram.StartDate) {
		return ErrInvalidSlotStartTime
	}

	firstSlot = time.Date(firstSlot.Year(), firstSlot.Month(), firstSlot.Day(),
		0, 0, 0, 0, time.Local)
	dbClass.StartDate = &firstSlot

	lastSlot := updateSlots[len(updateSlots)-1:][0].EndTime
	if lastSlot.After(dbProgram.EndDate) {
		return ErrInvalidSlotEndTime
	}

	lastSlot = time.Date(lastSlot.Year(), lastSlot.Month(), lastSlot.Day(),
		0, 0, 0, 0, time.Local)
	dbClass.EndDate = &lastSlot

	dbSlots, _ := c.queries.GetSlotsByClassID(ctx, dbClass.ID)
	if len(dbSlots) != len(updateSlots) {
		return ErrInvalidSlotCount
	}

	if hasOverlappingSlots(updateSlots) {
		return ErrInvalidSlotTime
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	for _, updateSlot := range updateSlots {
		classTeacher := sqlc.CheckTeacherInClassParams{
			TeacherID: updateSlot.TeacherId,
			ClassID:   dbClass.ID,
		}

		isTeacherInClass, err := c.queries.CheckTeacherInClass(ctx, classTeacher)
		if err != nil || !isTeacherInClass {
			return ErrTeacherIsNotInClass
		}

		teacherTime := CheckTeacherTime{
			TeacherId: &updateSlot.TeacherId,
			StartTime: &updateSlot.StartTime,
			EndTime:   &updateSlot.EndTime,
		}

		if c.IsTeacherAvailable(ctx, teacherTime) {
			return ErrTeacherNotAvailable
		}

		slot := sqlc.UpdateSlotParams{
			ID:        updateSlot.ID,
			StartTime: &updateSlot.StartTime,
			EndTime:   &updateSlot.EndTime,
			TeacherID: &updateSlot.TeacherId,
		}

		if err = qtx.UpdateSlot(ctx, slot); err != nil {
			return err
		}
	}

	updateClass := sqlc.UpdateActiveClassParams{
		ID:        dbClass.ID,
		StartDate: dbClass.StartDate,
		EndDate:   dbClass.EndDate,
	}

	err = qtx.UpdateActiveClass(ctx, updateClass)
	if err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}

func (c *Core) Delete(ctx *gin.Context, id uuid.UUID) error {
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbClass, err := c.queries.GetClassById(ctx, id)
	if err != nil {
		return ErrClassNotFound
	}
	fmt.Println(dbClass.StartDate.Before(time.Now()))
	if dbClass.StartDate.After(time.Now()) {
		err = c.queries.DeleteClass(ctx, dbClass.ID)
		if err != nil {
			return err
		}
	} else {
		err = c.queries.SoftDeleteClass(ctx, dbClass.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Core) Update(ctx *gin.Context, id uuid.UUID, updateClass UpdateClass) error {
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbClass, err := c.queries.GetClassById(ctx, id)
	if err != nil {
		return ErrClassNotFound
	}
	if updateClass.Code != dbClass.Code {
		_, err = c.queries.GetClassByCode(ctx, updateClass.Code)
		if err == nil {
			return ErrClassCodeAlreadyExist
		}
	}

	dbUpdateClass := sqlc.UpdateClassParams{
		Name: updateClass.Name,
		Code: updateClass.Code,
		ID:   dbClass.ID,
	}

	if updateClass.Password != nil {
		dbUpdateClass.Password = *updateClass.Password
	}

	err = c.queries.UpdateClass(ctx, dbUpdateClass)
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) IsTeacherAvailable(ctx *gin.Context, teacherTime CheckTeacherTime) bool {
	checkCondition := sqlc.CheckTeacherTimeOverlapParams{
		TeacherID: teacherTime.TeacherId,
		EndTime:   teacherTime.EndTime,
		StartTime: teacherTime.StartTime,
	}

	status, err := c.queries.CheckTeacherTimeOverlap(ctx, checkCondition)
	if err != nil {
		return false
	}

	return status
}

func hasOverlappingSlots(updateSlots []UpdateSlot) bool {
	for i := 1; i < len(updateSlots); i++ {
		if !updateSlots[i].StartTime.After(updateSlots[i-1].EndTime) {
			return true
		}

		for j := 0; j < i; j++ {
			if updateSlots[i].StartTime.Before(updateSlots[j].EndTime) && updateSlots[i].EndTime.After(updateSlots[j].StartTime) {
				return true
			}
		}
	}
	return false
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
		}
		weekDay := newClass.Slots.WeekDays[i%len(newClass.Slots.WeekDays)]
		weeksToAdd := i / len(newClass.Slots.WeekDays)

		slotDate := weekday.Next(currentDate.AddDate(0, 0, weeksToAdd), weekDay)
		currentDate = &slotDate

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
