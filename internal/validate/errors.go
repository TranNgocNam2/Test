package validate

type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

// FieldErrors represents a collection of field errors.
type FieldErrors []FieldError
