package validate

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var validate *validator.Validate

func init() {

	// Instantiate a validator.
	validate = validator.New()

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Check(val any) error {
	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		var validationErrors validator.ValidationErrors
		ok := errors.As(err, &validationErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, validationError := range validationErrors {
			field := FieldError{
				Field: validationError.Field(),
				Err:   validationError.Error(),
			}
			fields = append(fields, field)
		}

		return fields
	}

	return nil
}
