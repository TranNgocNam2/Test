package validate

import (
	"encoding/json"
	"github.com/pkg/errors"
)

var (
	ErrValidation = errors.New("Yêu cầu không đúng định dạng: %w!")
)

type FieldError struct {
	Field string `json:"field,omitempty"`
	Err   string `json:"error,omitempty"`
}

func NewFieldsError(field string, err error) error {
	return FieldErrors{
		{
			Field: field,
			Err:   err.Error(),
		},
	}
}

// FieldErrors represents a collection of field errors.
type FieldErrors []FieldError

func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}
