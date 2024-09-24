package school

import (
	"github.com/google/uuid"
)

type QueryFilter struct {
	ID         *uuid.UUID
	Name       *string `validate:"omitempty,min=3"`
	DistrictID *int
}
