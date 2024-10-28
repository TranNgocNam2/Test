package class

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type NewClass struct {
	ID        uuid.UUID
	ProgramID uuid.UUID
	SubjectID uuid.UUID
	Name      string
	Code      string
	Link      *string
	StartDate time.Time
	Slots     struct {
		WeekDays  []time.Weekday
		StartTime *time.Time
		StartDate *time.Time
	}
	Password string
}

type Details struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Link      string     `json:"link"`
	StartDate time.Time  `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	ProgramID uuid.UUID  `json:"programID"`
	Subject   Subject    `json:"subject"`
	Teachers  []Teacher  `json:"teachers"`
	Slots     []Slot     `json:"slots"`
}

type Subject struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

type Teacher struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Gender   *int16 `json:"gender"`
}

type Slot struct {
	ID        uuid.UUID `json:"id"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Session   Session   `json:"session"`
	Teacher   Teacher   `json:"teacher"`
}

type Session struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Index int32     `json:"index"`
}

func toCoreSubject(dbSubject sqlc.Subject) Subject {
	return Subject{
		ID:   dbSubject.ID,
		Name: dbSubject.Name,
		Code: dbSubject.Code,
	}
}

func toCoreSession(dbSession sqlc.Session) Session {
	return Session{
		ID:    dbSession.ID,
		Name:  dbSession.Name,
		Index: dbSession.Index,
	}
}

func toCoreTeacher(dbTeacher sqlc.User) Teacher {
	return Teacher{
		ID:       dbTeacher.ID,
		FullName: *dbTeacher.FullName,
		Email:    dbTeacher.Email,
		Phone:    *dbTeacher.Phone,
		Gender:   dbTeacher.Gender,
	}
}

func toCoreTeacherSlice(dbTeachers []sqlc.User) []Teacher {
	var teachers []Teacher
	for i, dbTeacher := range dbTeachers {
		teachers[i] = toCoreTeacher(dbTeacher)
	}
	return teachers
}
