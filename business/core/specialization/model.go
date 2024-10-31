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
	Subjects    []*struct {
		ID            uuid.UUID `json:"id"`
		Name          string    `json:"name"`
		Image         string    `json:"image"`
		Code          string    `json:"code"`
		LastUpdated   time.Time `json:"lastUpdated"`
		TotalSessions int64     `json:"totalSessions"`
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
}

type UpdateSpecialization struct {
	Name        string
	Code        string
	Status      int16
	Description string
	TimeAmount  float64
	Image       string
	Subjects    []uuid.UUID
}

type Specialization struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	Status        int16     `json:"status"`
	Image         *string   `json:"image"`
	TotalSubjects int64     `json:"totalSubjects"`
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
