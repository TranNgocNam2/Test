package validate

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func Check(val any) error {
	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Err:   verror.Translate(translator),
			}
			fields = append(fields, field)
		}

		return fields
	}

	return nil
}
