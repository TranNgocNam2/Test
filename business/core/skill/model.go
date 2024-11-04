package skill

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
)

type Skill struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type NewSkill struct {
	ID   uuid.UUID
	Name string
}

type UpdateSkill struct {
	Name string
}

func toCoreSkill(dbSkill sqlc.Skill) Skill {
	return Skill{
		ID:   dbSkill.ID,
		Name: dbSkill.Name,
	}
}
