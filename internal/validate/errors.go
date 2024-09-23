package validate

import "encoding/json"

type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
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

//func (fe FieldErrors) Fields() map[string]string {
//	m := make(map[string]string)
//	for _, fld := range fe {
//		m[fld.Field] = fld.Err
//	}
//	return m
//}
