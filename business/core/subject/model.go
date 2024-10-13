package subject

import "github.com/google/uuid"

type Subject struct {
	Name            string
	Description     string
	Image           string
	TimePerSession  int
	SessionPerWeeke int
	Skill           []uuid.UUID
}
