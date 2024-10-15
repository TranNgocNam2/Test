package subject

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Subject struct {
	ID             uuid.UUID
	Name           string
	Code           string
	Description    string
	Image          string
	TimePerSession int
	SessionPerWeek int
	Skills         []uuid.UUID
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
