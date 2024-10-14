package subject

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Subject struct {
	Name           string
	Code           string
	Description    string
	Image          string
	TimePerSession int
	SessionPerWeek int
	Skills         []uuid.UUID
}

type SubjectDraft struct {
	ID          uuid.UUID
	Name        string
	Code        string
	Description string
	Image       string
	Status      int
	Skills      []uuid.UUID
	Sessions    []Session
}

type SubjectPulished struct {
	ID          uuid.UUID
	Name        string
	Code        string
	Description string
	Image       string
	Skills      []uuid.UUID
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
