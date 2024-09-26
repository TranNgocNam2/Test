package school

import (
	"Backend/internal/validate"
	"fmt"
)

type QueryFilter struct {
	Name *string `validate:"omitempty,min=3"`
}

func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func (qf *QueryFilter) WithName(name string) {
	qf.Name = &name
}
