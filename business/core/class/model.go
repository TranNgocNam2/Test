package class

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type NewClass struct {
	ID        uuid.UUID
	ProgramId uuid.UUID
	SubjectId uuid.UUID
	Name      string
	Code      string
	Link      *string
	Slots     struct {
		WeekDays  []time.Weekday
		StartTime *time.Time
		StartDate *time.Time
	}
	Password string
}

type UpdateSlot struct {
	ID        uuid.UUID
	StartTime time.Time
	EndTime   time.Time
	TeacherId string
	Index     int32
}

type UpdateClass struct {
	Name     string
	Code     string
	Password *string
}

type UpdateMeeting struct {
	Link string
}

type Details struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Code          string     `json:"code"`
	Link          string     `json:"link"`
	Password      *string    `json:"password,omitempty"`
	StartDate     *time.Time `json:"startDate"`
	EndDate       *time.Time `json:"endDate"`
	Program       Program    `json:"program"`
	Subject       Subject    `json:"subject"`
	Slots         []Slot     `json:"slots"`
	Skills        []Skill    `json:"skills"`
	TotalLearners int64      `json:"totalLearners,omitempty"`
}

type Class struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Code          string     `json:"code"`
	Program       Program    `json:"program"`
	StartDate     *time.Time `json:"startDate"`
	EndDate       *time.Time `json:"endDate"`
	Status        *int16     `json:"status,omitempty"`
	Subject       Subject    `json:"subject"`
	Skills        []Skill    `json:"skills"`
	TotalLearners int64      `json:"totalLearners"`
	TotalSlots    int64      `json:"totalSlots,omitempty"`
	CurrentSlot   int64      `json:"currentSlot,omitempty"`
}

type Skill struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Program struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Subject struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}

type Teacher struct {
	ID       string  `json:"id"`
	FullName string  `json:"fullName"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
}

type Slot struct {
	ID        uuid.UUID `json:"id"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Index     int32     `json:"index"`
	Session   Session   `json:"session"`
	Teacher   Teacher   `json:"teacher"`
}

type Session struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CheckTeacherTime struct {
	TeacherId string
	SlotId    uuid.UUID
	StartTime time.Time
	EndTime   time.Time
}

func toCoreSubject(dbSubject sqlc.Subject) Subject {
	return Subject{
		ID:   dbSubject.ID,
		Name: dbSubject.Name,
		Code: dbSubject.Code,
	}
}

func toCoreProgram(dbProgram sqlc.Program) Program {
	return Program{
		ID:   dbProgram.ID,
		Name: dbProgram.Name,
	}
}

func toCoreSession(dbSession sqlc.Session) Session {
	return Session{
		ID:   dbSession.ID,
		Name: dbSession.Name,
	}
}

func toCoreTeacher(dbTeacher sqlc.User) Teacher {
	return Teacher{
		ID:       dbTeacher.ID,
		FullName: *dbTeacher.FullName,
		Email:    dbTeacher.Email,
		Phone:    dbTeacher.Phone,
	}
}

func toCoreSkillSlice(dbSkills []sqlc.Skill) []Skill {
	skills := make([]Skill, len(dbSkills))
	for i, dbSkill := range dbSkills {
		skills[i] = Skill{
			ID:   dbSkill.ID,
			Name: dbSkill.Name,
		}
	}
	return skills
}
