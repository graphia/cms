package main

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

func validationErrorsToJSON(errs error) map[string]string {

	messages := make(map[string]string)

	for _, e := range errs.(validator.ValidationErrors) {
		messages[e.Field()] = validationErrorMessage(e)
	}

	return messages
}

func validationErrorMessage(e validator.FieldError) string {
	// depending on the type of validation display an appropriate message

	switch e.Tag() {
	case "required":
		return "is a required field"
	case "min":

		switch e.Type().String() {
		case "int":
			return fmt.Sprintf("must be at least %s", e.Param())
		case "string":
			return fmt.Sprintf("must be at least %s characters", e.Param())
		}

	case "max":

		switch e.Type().String() {
		case "int":
			return fmt.Sprintf("must be smaller than %s", e.Param())
		case "string":
			return fmt.Sprintf("must be no more than %s characters", e.Param())
		}

	case "alphanumunicode":
		return "only alphanumeric unicode characters are permitted"
	}

	return "there was an unspecified error"

}
