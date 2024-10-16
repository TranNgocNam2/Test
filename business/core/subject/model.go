package subject

import (
	"Backend/business/db/sqlc"
	"encoding/json"

	"github.com/google/uuid"
)

type Subject struct {
	ID             uuid.UUID
	Name           string
	Code           string
	Description    string
	Image          string
	TimePerSession int16
	SessionPerWeek int16
	Skills         []Skill
	TotalSessions  int
}

type SubjectDetail struct {
	ID            uuid.UUID
	Name          string
	Code          string
	Description   string
	Image         string
	Status        int
	Skills        []Skill
	Sessions      []Session
	TotalSessions int
}

type Skill struct {
	ID   uuid.UUID
	Name string
}

type Session struct {
	ID        uuid.UUID
	Name      string
	Index     int
	Materials []Material
}

type Material struct {
	ID       uuid.UUID
	Name     string
	Index    int
	IsShared bool
	Data     json.RawMessage
}

func toCoreSubject(dbSubject sqlc.Subject) Subject {
	return Subject{
		ID:             dbSubject.ID,
		Name:           dbSubject.Name,
		Code:           dbSubject.Code,
		Description:    dbSubject.Description,
		Image:          *dbSubject.ImageLink,
		TimePerSession: dbSubject.TimePerSession,
		SessionPerWeek: dbSubject.SessionsPerWeek,
	}
}
