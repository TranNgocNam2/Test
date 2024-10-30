package subject

import (
	"Backend/business/db/sqlc"
	"encoding/json"

	"github.com/google/uuid"
)

type Subject struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Code           string    `json:"code"`
	Description    *string   `json:"description,omitempty"`
	Image          *string   `json:"image,omitempty"`
	TimePerSession int16     `json:"timePerSession"`
	SessionPerWeek int16     `json:"sessionPerWeek"`
	Skills         []Skill   `json:"skills,omitempty"`
	TotalSessions  int       `json:"totalSessions"`
}

type SubjectDetail struct {
	ID             uuid.UUID    `json:"id"`
	Name           string       `json:"name"`
	Code           string       `json:"code"`
	TimePerSession int          `json:"timePerSession"`
	MinPassGrade   float32      `json:"minPassGrade"`
	MinAttendance  float32      `json:"minAttendance"`
	Description    string       `json:"description"`
	Image          string       `json:"image"`
	Status         int          `json:"status"`
	Skills         []Skill      `json:"skills"`
	Sessions       []Session    `json:"sessions"`
	Transcripts    []Transcript `json:"transcripts"`
	TotalSessions  int          `json:"totalSessions"`
}

type Skill struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Session struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Index     int        `json:"index"`
	Materials []Material `json:"materials"`
}

type Material struct {
	ID       uuid.UUID       `json:"id"`
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Index    int             `json:"index"`
	IsShared bool            `json:"isShared"`
	Data     json.RawMessage `json:"data"`
}

type Transcript struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Index      int       `json:"index"`
	Percentage float32   `json:"percentage"`
	MinGrade   float32   `json:"minGrade"`
}

func toCoreSubject(dbSubject sqlc.Subject) Subject {
	return Subject{
		ID:             dbSubject.ID,
		Name:           dbSubject.Name,
		Code:           dbSubject.Code,
		Description:    dbSubject.Description,
		Image:          dbSubject.ImageLink,
		TimePerSession: dbSubject.TimePerSession,
	}
}
