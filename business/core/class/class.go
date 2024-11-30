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
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"gitlab.com/innovia69420/kit/enum/role"
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

	dbProgram, err := c.queries.GetProgramById(ctx, newClass.ProgramId)
	if err != nil {
		return uuid.Nil, model.ErrProgramNotFound
	}

	dbSubject, err := c.queries.GetPublishedSubjectById(ctx, newClass.SubjectId)
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
			0, 0, 0, 0, time.UTC)
		startDateClass = &startDate
	}

	var endDateClass *time.Time
	lastSlot := slots[len(slots)-1:][0].EndTime
	if lastSlot != nil && lastSlot.Before(dbProgram.EndDate) {
		endDate := time.Date(lastSlot.Year(), lastSlot.Month(), lastSlot.Day(),
			0, 0, 0, 0, time.UTC)
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

func (c *Core) QueryByManager(ctx *gin.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Class, error) {
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return nil, err
	}

	if err := filter.Validate(); err != nil {
		return nil, nil
	}

	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT
						id, name, code, subject_id, program_id, link, start_date, end_date, status
			FROM classes`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, false)

	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var dbClasses []sqlc.Class
	err = pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbClasses)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, nil
	}

	if dbClasses == nil {
		return nil, nil
	}

	var classes []Class

	for _, dbClass := range dbClasses {
		class := Class{
			ID:        dbClass.ID,
			Name:      dbClass.Name,
			Code:      dbClass.Code,
			StartDate: dbClass.StartDate,
			EndDate:   dbClass.EndDate,
			Status:    &dbClass.Status,
		}

		dbProgram, _ := c.queries.GetProgramById(ctx, dbClass.ProgramID)
		class.Program = toCoreProgram(dbProgram)

		dbSubject, _ := c.queries.GetSubjectById(ctx, dbClass.SubjectID)
		class.Subject = toCoreSubject(dbSubject)

		dbSkills, err := c.queries.GetSkillsBySubjectId(ctx, dbSubject.ID)
		if err != nil {
			return nil, nil
		}
		class.Skills = toCoreSkillSlice(dbSkills)

		totalLearners, err := c.queries.CountLearnersByClassId(ctx, dbClass.ID)
		if err != nil {
			return nil, nil
		}
		class.TotalLearners = totalLearners
		class.TotalSlots, _ = c.queries.CountSlotsByClassId(ctx, dbClass.ID)
		class.CurrentSlot, _ = c.queries.CountCompletedSlotsByClassId(ctx, dbClass.ID)

		classes = append(classes, class)
	}

	return classes, nil
}

func (c *Core) QueryByTeacher(ctx *gin.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Class, error) {
	teacherId, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		return nil, err
	}

	if err := filter.Validate(); err != nil {
		return nil, nil
	}

	data := map[string]interface{}{
		"teacher_id":    teacherId,
		"status":        COMPLETED,
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT
					c.id, c.name, c.code, c.subject_id, c.program_id, c.link, c.start_date, c.end_date
			FROM classes c
				JOIN slots s ON s.class_id = c.id
					WHERE s.teacher_id = :teacher_id AND c.status = :status`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf, true)
	buf.WriteString(" GROUP BY c.id")
	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var dbClasses []sqlc.Class
	err = pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbClasses)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, nil
	}

	if dbClasses == nil {
		return nil, nil
	}

	var classes []Class

	for _, dbClass := range dbClasses {
		class := Class{
			ID:        dbClass.ID,
			Name:      dbClass.Name,
			Code:      dbClass.Code,
			StartDate: dbClass.StartDate,
			EndDate:   dbClass.EndDate,
			Status:    nil,
		}

		dbProgram, _ := c.queries.GetProgramById(ctx, dbClass.ProgramID)
		class.Program = toCoreProgram(dbProgram)

		dbSubject, _ := c.queries.GetSubjectById(ctx, dbClass.SubjectID)
		class.Subject = toCoreSubject(dbSubject)

		dbSkills, err := c.queries.GetSkillsBySubjectId(ctx, dbSubject.ID)
		if err != nil {
			return nil, nil
		}
		class.Skills = toCoreSkillSlice(dbSkills)

		totalLearners, err := c.queries.CountLearnersByClassId(ctx, dbClass.ID)
		if err != nil {
			return nil, nil
		}
		class.TotalLearners = totalLearners

		class.TotalSlots, _ = c.queries.CountSlotsByClassId(ctx, dbClass.ID)
		class.CurrentSlot, _ = c.queries.CountCompletedSlotsByClassId(ctx, dbClass.ID)

		classes = append(classes, class)
	}

	return classes, nil
}

func (c *Core) CountByTeacher(ctx *gin.Context, filter QueryFilter) int {
	teacherId, _ := middleware.AuthorizeTeacher(ctx, c.queries)

	data := map[string]interface{}{
		"teacher_id": teacherId,
		"status":     COMPLETED,
	}
	if err := filter.Validate(); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	const q = `SELECT COUNT (DISTINCT (c.id)) AS count
			FROM classes c
				JOIN slots s ON s.class_id = c.id
					WHERE s.teacher_id = :teacher_id AND c.status = :status`

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

func (c *Core) QueryByLearner(ctx *gin.Context) ([]Class, error) {
	learner, err := middleware.AuthorizeVerifiedLearner(ctx, c.queries)
	if err != nil {
		return nil, err
	}

	dbClasses, err := c.queries.GetClassesByLearnerId(ctx, learner.ID)
	if err != nil {
		return nil, nil
	}

	if dbClasses == nil {
		return nil, nil
	}

	var classes []Class

	for _, dbClass := range dbClasses {
		class := Class{
			ID:        dbClass.ID,
			Name:      dbClass.Name,
			Code:      dbClass.Code,
			StartDate: dbClass.StartDate,
			EndDate:   dbClass.EndDate,
			Status:    nil,
		}

		dbProgram, _ := c.queries.GetProgramById(ctx, dbClass.ProgramID)
		class.Program = toCoreProgram(dbProgram)

		dbSubject, _ := c.queries.GetSubjectById(ctx, dbClass.SubjectID)
		class.Subject = toCoreSubject(dbSubject)

		dbSkills, err := c.queries.GetSkillsBySubjectId(ctx, dbSubject.ID)
		if err != nil {
			return nil, nil
		}
		class.Skills = toCoreSkillSlice(dbSkills)

		totalLearners, err := c.queries.CountLearnersByClassId(ctx, dbClass.ID)
		if err != nil {
			return nil, nil
		}
		class.TotalLearners = totalLearners
		class.TotalSlots, _ = c.queries.CountSlotsByClassId(ctx, dbClass.ID)
		class.CurrentSlot, _ = c.queries.CountCompletedSlotsByClassId(ctx, dbClass.ID)

		classes = append(classes, class)
	}
	return classes, nil
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
	applyFilter(filter, data, buf, false)

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
	user, err := middleware.AuthorizeUser(ctx, c.queries)
	if err != nil {
		return Details{}, err
	}

	dbClass, err := c.queries.GetClassById(ctx, id)
	if err != nil {
		return Details{}, model.ErrClassNotFound
	}

	class := Details{
		ID:        dbClass.ID,
		Name:      dbClass.Name,
		Code:      dbClass.Code,
		Link:      *dbClass.Link,
		StartDate: dbClass.StartDate,
		EndDate:   dbClass.EndDate,
		Password:  &dbClass.Password,
	}

	if user.AuthRole == role.LEARNER {
		class.Password = nil
	}

	dbSubject, err := c.queries.GetSubjectById(ctx, dbClass.SubjectID)
	if err != nil {
		return Details{}, model.ErrSubjectNotFound
	}
	class.Subject = toCoreSubject(dbSubject)

	dbProgram, err := c.queries.GetProgramById(ctx, dbClass.ProgramID)
	if err != nil {
		return Details{}, model.ErrProgramNotFound
	}
	class.Program = toCoreProgram(dbProgram)

	dbSkills, err := c.queries.GetSkillsBySubjectId(ctx, dbSubject.ID)
	if err != nil {
		return Details{}, model.ErrSkillNotFound
	}
	class.Skills = toCoreSkillSlice(dbSkills)

	var slots []Slot
	dbSlots, _ := c.queries.GetSlotsByClassId(ctx, dbClass.ID)
	for _, dbSlot := range dbSlots {
		dbSession, _ := c.queries.GetSessionById(ctx, dbSlot.SessionID)
		session := toCoreSession(dbSession)
		startTime := *dbSlot.StartTime
		endTime := *dbSlot.EndTime
		slot := Slot{
			ID:         dbSlot.ID,
			StartTime:  startTime,
			EndTime:    endTime,
			Index:      dbSlot.Index,
			Session:    session,
			RecordLink: dbSlot.RecordLink,
		}

		if dbSlot.TeacherID != nil {
			dbTeacher, _ := c.queries.GetTeacherById(ctx, *dbSlot.TeacherID)
			slot.Teacher = toCoreTeacher(dbTeacher)
		}

		slots = append(slots, slot)
	}
	class.Slots = slots

	return class, nil
}

func (c *Core) UpdateSlot(ctx *gin.Context, id uuid.UUID, updateSlots []UpdateSlot, status int) error {
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbClass, err := c.queries.GetClassById(ctx, id)
	if err != nil {
		return model.ErrClassNotFound
	}

	currentTime := time.Now().UTC()

	dbProgram, _ := c.queries.GetProgramById(ctx, dbClass.ProgramID)
	if err = validateSlotTimes(dbClass, dbProgram, updateSlots); err != nil {
		return err
	}

	slotCount, _ := c.queries.CountSlotsByClassId(ctx, dbClass.ID)
	if int(slotCount) != len(updateSlots) {
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
		dbSlot, err := c.queries.GetSlotByIdAndIndex(ctx,
			sqlc.GetSlotByIdAndIndexParams{
				ID:    updateSlot.ID,
				Index: updateSlot.Index,
			})
		if err != nil {
			return model.ErrSlotNotFound
		}

		if dbSlot.StartTime != nil && dbSlot.StartTime.Before(currentTime) {
			updateSlot.StartTime = *dbSlot.StartTime
			updateSlot.EndTime = *dbSlot.EndTime
			updateSlot.TeacherId = *dbSlot.TeacherID
		}

		if dbSlot.StartTime.After(currentTime) {
			isTeacherAvailable, err := c.IsTeacherAvailable(ctx, CheckTeacherTime{
				TeacherId: updateSlot.TeacherId,
				StartTime: updateSlot.StartTime,
				EndTime:   updateSlot.EndTime,
			})
			if err != nil || !isTeacherAvailable {
				return model.ErrTeacherNotAvailable
			}
		}

		slot := sqlc.UpdateSlotParams{
			ID:        dbSlot.ID,
			StartTime: &updateSlot.StartTime,
			EndTime:   &updateSlot.EndTime,
			TeacherID: &updateSlot.TeacherId,
		}

		if err = qtx.UpdateSlot(ctx, slot); err != nil {
			return err
		}
	}

	if dbClass.Status == COMPLETED {
		status = COMPLETED
	}

	updateClass := sqlc.UpdateClassStatusAndDateParams{
		ID:        dbClass.ID,
		StartDate: dbClass.StartDate,
		EndDate:   dbClass.EndDate,
		Status:    int16(status),
	}

	err = qtx.UpdateClassStatusAndDate(ctx, updateClass)
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

	if dbClass.StartDate.After(time.Now().UTC()) {
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

func (c *Core) UpdateMeetingLink(ctx *gin.Context, id uuid.UUID, updateMeeting UpdateMeeting) error {
	teacherId, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		return err
	}

	class, err := c.queries.GetClassById(ctx, id)
	if err != nil {
		return model.ErrClassNotFound
	}

	if class.Status != COMPLETED {
		return model.ErrClassNotCompleted
	}

	if class.EndDate != nil && class.EndDate.UTC().Before(time.Now().UTC()) {
		return model.ErrClassIsEnded
	}

	isTeacherInClass, err := c.queries.CheckTeacherInClass(ctx, sqlc.CheckTeacherInClassParams{
		TeacherID: &teacherId,
		ClassID:   id,
	})
	if err != nil || !isTeacherInClass {
		return model.ErrTeacherIsNotInClass
	}

	err = c.queries.UpdateMeetingLink(ctx, sqlc.UpdateMeetingLinkParams{
		Link:      &updateMeeting.Link,
		UpdatedBy: &teacherId,
		ID:        class.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Core) IsTeacherAvailable(ctx *gin.Context, teacherTime CheckTeacherTime) (bool, error) {
	teacher, err := c.queries.GetTeacherById(ctx, teacherTime.TeacherId)
	if err != nil {
		return false, model.ErrTeacherNotFound
	}
	checkCondition := sqlc.CheckTeacherTimeOverlapParams{
		TeacherID: &teacher.ID,
		EndTime:   &teacherTime.EndTime,
		StartTime: &teacherTime.StartTime,
		SlotID:    teacherTime.SlotId,
	}

	status, err := c.queries.CheckTeacherTimeOverlap(ctx, checkCondition)
	if err != nil {
		return false, err
	}

	return !status, nil
}

func validateSlotTimes(dbClass sqlc.Class, dbProgram sqlc.Program, updateSlots []UpdateSlot) error {
	firstSlot := updateSlots[0].StartTime
	if firstSlot.Before(dbProgram.StartDate) {
		return model.ErrInvalidSlotStartTime
	}

	firstSlot = time.Date(firstSlot.Year(), firstSlot.Month(), firstSlot.Day(), 0, 0, 0, 0, time.UTC)
	dbClass.StartDate = &firstSlot

	lastSlot := updateSlots[len(updateSlots)-1].EndTime
	if lastSlot.After(dbProgram.EndDate) {
		return model.ErrInvalidSlotEndTime
	}

	if lastSlot.Hour() != 0 || lastSlot.Minute() != 0 {
		lastSlot = time.Date(lastSlot.Year(), lastSlot.Month(), lastSlot.Day()+1, 0, 0, 0, 0, time.UTC)
	} else {
		lastSlot = time.Date(lastSlot.Year(), lastSlot.Month(), lastSlot.Day(), 0, 0, 0, 0, time.UTC)
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

func generateSlots(newClass NewClass, sessions []sqlc.Session, duration float32, endDate time.Time) []sqlc.CreateSlotsParams {
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
		} else {
			weekDay := newClass.Slots.WeekDays[i%len(newClass.Slots.WeekDays)]
			weeksToAdd := i / len(newClass.Slots.WeekDays)

			slotDate := weekday.Next(currentDate.AddDate(0, 0, weeksToAdd), weekDay)
			currentDate = &slotDate

			slotStartTime := time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(),
				newClass.Slots.StartTime.Hour(), newClass.Slots.StartTime.Minute(), 0, 0, time.UTC)
			slotEndTime := slotStartTime.Add(time.Duration(duration * float32(time.Hour)))

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
	}

	return slots
}
