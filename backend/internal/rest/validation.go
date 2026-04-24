package rest

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func mapValidationErrors(err error, request any) (map[string]string, bool) {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil, false
	}

	fields := make(map[string]string, len(validationErrors))
	for _, fieldErr := range validationErrors {
		fields[fieldErr.Field()] = validationErrorMessage(fieldErr)
	}

	return fields, true
}

func validationErrorMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "is required"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fieldErr.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fieldErr.Param())
	case "alphanum":
		return "must contain only letters and numbers"
	case "url":
		return "must be a valid URL"
	default:
		return fieldErr.Error()
	}
}
