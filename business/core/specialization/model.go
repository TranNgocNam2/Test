package specialization

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type Details struct {
	ID          uuid.UUID
	Name        string
	Code        string
	Status      int16
	Description *string
	TimeAmount  *float64
	Image       *string
	CreatedAt   time.Time
	Skills      []*struct {
		ID   uuid.UUID
		Name string
	}
	Subjects []*struct {
		ID           uuid.UUID
		Name         string
		Image        string
		Code         string
		LastUpdated  time.Time
		TotalSession int64
	}
}

type NewSpecialization struct {
	ID          uuid.UUID
	Name        string
	Code        string
	Status      int16
	Description *string
	TimeAmount  *float64
	Image       *string
	Skills      []uuid.UUID
	Subjects    []uuid.UUID
}

type UpdateSpecialization struct {
	Name        string
	Code        string
	Status      int16
	Description string
	TimeAmount  float64
	Image       string
	Skills      []uuid.UUID
	Subjects    []uuid.UUID
}

type Specialization struct {
	ID           uuid.UUID
	Name         string
	Code         string
	Status       int16
	Image        *string
	TotalSubject int64
	Skills       []*struct {
		ID   uuid.UUID
		Name string
	}
}

func toCoreSpecialization(dbSpec sqlc.Specialization) Specialization {
	spec := Specialization{
		ID:     dbSpec.ID,
		Name:   dbSpec.Name,
		Code:   dbSpec.Code,
		Status: dbSpec.Status,
		Image:  dbSpec.ImageLink,
	}

	return spec
}
func toCoreSpecializationDetails(dbSpec sqlc.Specialization, dbSpecSkills []sqlc.Skill) Details {
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

	for _, skill := range dbSpecSkills {
		specDetails.Skills = append(specDetails.Skills, &struct {
			ID   uuid.UUID
			Name string
		}{
			ID:   skill.ID,
			Name: skill.Name,
		})
	}

	return specDetails
}
