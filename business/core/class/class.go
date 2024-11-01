package class

import (
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/weekday"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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

func (c *Core) Create(ctx *gin.Context, newClass NewClass) (uuid.UUID, error) {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = c.queries.GetClassCompletedByCode(ctx, newClass.Code)
	if err == nil {
		return uuid.Nil, model.ErrClassCodeAlreadyExist
	}

	dbProgram, err := c.queries.GetProgramByID(ctx, newClass.ProgramId)
	if err != nil {
		return uuid.Nil, model.ErrProgramNotFound
	}

	dbSubject, err := c.queries.GetPublishedSubjectByID(ctx, newClass.SubjectId)
	if err != nil {
		return uuid.Nil, model.ErrSubjectNotFound
	}

	sessions, err := c.queries.GetSessionsBySubjectId(ctx, newClass.SubjectId)
	if err != nil {
		return uuid.Nil, model.ErrSessionNotFound
	}

	slots := generateSlots(newClass, sessions, dbSubject.TimePerSession, dbProgram.EndDate)

	// Check if the last slot's end time is after the programs end time
	var startDateClass *time.Time
	firstSlot := slots[0].StartTime
	if firstSlot != nil && firstSlot.Before(dbProgram.StartDate) {
		return uuid.Nil, model.ErrInvalidClassStartTime
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
		return Details{}, model.ErrClassNotFound
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
		return Details{}, model.ErrSubjectNotFound
	}
	class.Subject = toCoreSubject(dbSubject)

	dbProgram, err := c.queries.GetProgramByID(ctx, dbClass.ProgramID)
	if err != nil {
		return Details{}, model.ErrProgramNotFound
	}
	class.Program = toCoreProgram(dbProgram)

	dbTeachers, err := c.queries.GetTeachersByClassId(ctx, dbClass.ID)
	if err == nil {
		class.Teachers = toCoreTeacherSlice(dbTeachers)
	}

	var slots []Slot
	dbSlots, _ := c.queries.GetSlotsByClassId(ctx, dbClass.ID)
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
		return model.ErrClassNotFound
	}

	for _, teacherId := range teacherIds {
		_, err = c.queries.GetTeacherByID(ctx, teacherId)
		if err != nil {
			return model.ErrTeacherNotFound
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
		return model.ErrClassNotFound
	}

	dbProgram, _ := c.queries.GetProgramByID(ctx, dbClass.ProgramID)
	if err = validateSlotTimes(dbClass, dbProgram, updateSlots); err != nil {
		return err
	}

	dbSlots, _ := c.queries.GetSlotsByClassId(ctx, dbClass.ID)
	if len(dbSlots) != len(updateSlots) {
		return model.ErrInvalidSlotCount
	}

	if hasOverlappingSlots(updateSlots) {
		return model.ErrInvalidSlotTime
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	for _, updateSlot := range updateSlots {
		dbSlot, err := c.queries.GetSlotById(ctx, updateSlot.ID)
		if err != nil {
			return model.ErrSlotNotFound
		}

		slot := sqlc.UpdateSlotParams{
			ID:        dbSlot.ID,
			StartTime: &updateSlot.StartTime,
			EndTime:   &updateSlot.EndTime,
		}

		if updateSlot.TeacherId != "" {
			err = c.validateTeacherInClass(ctx, dbClass.ID, updateSlot.TeacherId)
			if err != nil {
				return err
			}
			slot.TeacherID = &updateSlot.TeacherId
		}

		if err = qtx.UpdateSlot(ctx, slot); err != nil {
			return err
		}
	}
	classStatus := INCOMPLETE
	totalSlots, _ := qtx.CountSlotsHaveTeacherByClassId(ctx, dbClass.ID)
	if int(totalSlots) == len(updateSlots) {
		classStatus = COMPLETED
	}

	updateClass := sqlc.UpdateActiveClassParams{
		ID:        dbClass.ID,
		StartDate: dbClass.StartDate,
		EndDate:   dbClass.EndDate,
		Status:    int16(classStatus),
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
		return model.ErrClassNotFound
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
		return model.ErrClassNotFound
	}
	if updateClass.Code != dbClass.Code {
		_, err = c.queries.GetClassCompletedByCode(ctx, updateClass.Code)
		if err == nil {
			return model.ErrClassCodeAlreadyExist
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

func (c *Core) IsTeacherAvailable(ctx *gin.Context, teacherTime CheckTeacherTime) (map[uuid.UUID]bool, error) {
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return nil, err
	}

	dbClass, err := c.queries.GetClassById(ctx, teacherTime.ClassId)
	if err != nil {
		return nil, model.ErrClassNotFound
	}
	dbSlots, _ := c.queries.GetSlotsByClassId(ctx, dbClass.ID)
	availabilityMap := make(map[uuid.UUID]bool)

	for _, dbSlot := range dbSlots {
		checkCondition := sqlc.CheckTeacherTimeOverlapParams{
			TeacherID: teacherTime.TeacherId,
			EndTime:   dbSlot.EndTime,
			StartTime: dbSlot.StartTime,
		}
		status, err := c.queries.CheckTeacherTimeOverlap(ctx, checkCondition)
		if err != nil {
			return nil, err
		}
		availabilityMap[dbSlot.ID] = !status
	}

	return availabilityMap, nil
}

func (c *Core) validateTeacherInClass(ctx *gin.Context, classID uuid.UUID, teacherID string) error {
	classTeacher := sqlc.CheckTeacherInClassParams{
		TeacherID: teacherID,
		ClassID:   classID,
	}
	isTeacherInClass, err := c.queries.CheckTeacherInClass(ctx, classTeacher)
	if err != nil || !isTeacherInClass {
		return model.ErrTeacherIsNotInClass
	}
	return nil
}

func validateSlotTimes(dbClass sqlc.Class, dbProgram sqlc.Program, updateSlots []UpdateSlot) error {
	firstSlot := updateSlots[0].StartTime
	if firstSlot.Before(dbProgram.StartDate) {
		return model.ErrInvalidSlotStartTime
	}

	firstSlot = time.Date(firstSlot.Year(), firstSlot.Month(), firstSlot.Day(), 0, 0, 0, 0, time.Local)
	dbClass.StartDate = &firstSlot

	lastSlot := updateSlots[len(updateSlots)-1].EndTime
	if lastSlot.After(dbProgram.EndDate) {
		return model.ErrInvalidSlotEndTime
	}

	if lastSlot.Hour() != 0 || lastSlot.Minute() != 0 {
		lastSlot = time.Date(lastSlot.Year(), lastSlot.Month(), lastSlot.Day()+1, 0, 0, 0, 0, time.Local)
	} else {
		lastSlot = time.Date(lastSlot.Year(), lastSlot.Month(), lastSlot.Day(), 0, 0, 0, 0, time.Local)
	}
	dbClass.EndDate = &lastSlot

	return nil
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