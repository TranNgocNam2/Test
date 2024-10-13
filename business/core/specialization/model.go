package specialization

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type Specialization struct {
	ID          uuid.UUID
	Name        string
	Code        string
	Status      int16
	Description *string
	TimeAmount  *float64
	Image       *string
	CreatedAt   time.Time
	Skills      []*struct {
		ID   *uuid.UUID
		Name *string
	}
	Subjects []*struct {
		ID           *uuid.UUID
		Name         *string
		Image        *string
		Code         *string
		LastUpdated  time.Time
		TotalSession *int64
	}
}

func toCoreSpecialization(dbSpec sqlc.Specialization) Specialization {
	return Specialization{
		ID:          dbSpec.ID,
		Name:        dbSpec.Name,
		Code:        dbSpec.Code,
		Status:      dbSpec.Status,
		Description: dbSpec.Description,
		TimeAmount:  dbSpec.TimeAmount,
		Image:       dbSpec.ImageLink,
		CreatedAt:   dbSpec.CreatedAt,
	}
}
