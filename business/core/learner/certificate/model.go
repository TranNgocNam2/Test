package certificate

import (
	"github.com/google/uuid"
	"time"
)

type Certificate struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	CreatedAt      time.Time       `json:"createdAt"`
	Specialization *Specialization `json:"specialization,omitempty"`
	Subject        *Subject        `json:"subject,omitempty"`
	Program        *Program        `json:"program,omitempty"`
}

type Specialization struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	TimeAmount  float32   `json:"timeAmount"`
	ImageLink   string    `json:"imageLink"`
	Description string    `json:"description"`
}
type Subject struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	ImageLink   string    `json:"imageLink"`
}

type Program struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
