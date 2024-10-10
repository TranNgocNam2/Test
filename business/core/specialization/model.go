package specialization

import (
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
		ID            *uuid.UUID
		Name          *string
		Image         *string
		Code          *string
		LastUpdated   time.Time
		TotalSessions *int16
	} `json:"subjects,omitempty"`
}
