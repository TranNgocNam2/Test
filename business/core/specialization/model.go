package specialization

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type Details struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Status      int16     `json:"status"`
	Description *string   `json:"description"`
	TimeAmount  *float64  `json:"timeAmount"`
	Image       *string   `json:"image"`
	CreatedAt   time.Time `json:"createdAt"`
	Subjects    []Subject `json:"subjects"`
}

type NewSpecialization struct {
	ID          uuid.UUID
	Name        string
	Code        string
	Status      int16
	Description *string
	TimeAmount  *float64
	Image       *string
}

type UpdateSpecialization struct {
	Name        string
	Code        string
	Status      int16
	Description string
	TimeAmount  float64
	Image       string
	Subjects    []SpecSubject
}

type SpecSubject struct {
	ID    uuid.UUID
	Index int16
}

type Specialization struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	Status        int16     `json:"status"`
	Image         *string   `json:"image"`
	Description   *string   `json:"description"`
	TotalSubjects int64     `json:"totalSubjects"`
}

type Subject struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Image         string     `json:"image"`
	Code          string     `json:"code"`
	LastUpdated   *time.Time `json:"lastUpdated"`
	Skills        []Skill    `json:"skills"`
	Index         int16      `json:"index"`
	MinPassGrade  float32    `json:"minPassGrade"`
	MinAttendance float32    `json:"minAttendance"`
	CreatedBy     string     `json:"createdBy"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedBy     *string    `json:"updatedBy"`
	TotalSessions int64      `json:"totalSessions"`
}

type Skill struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func toCoreSkill(dbSkill sqlc.Skill) Skill {
	skill := Skill{
		ID:   dbSkill.ID,
		Name: dbSkill.Name,
	}

	return skill
}

func toCoreSpecialization(dbSpec sqlc.Specialization) Specialization {
	spec := Specialization{
		ID:          dbSpec.ID,
		Name:        dbSpec.Name,
		Code:        dbSpec.Code,
		Status:      dbSpec.Status,
		Image:       dbSpec.ImageLink,
		Description: dbSpec.Description,
	}

	return spec
}
func toCoreSpecializationDetails(dbSpec sqlc.Specialization) Details {
	specDetails := Details{
		ID:          dbSpec.ID,
		Name:        dbSpec.Name,
		Code:        dbSpec.Code,
		Status:      dbSpec.Status,
		Description: dbSpec.Description,
		TimeAmount:  dbSpec.TimeAmount,
		Image:       dbSpec.ImageLink,
		CreatedAt:   dbSpec.CreatedAt,
	}

	return specDetails
}
