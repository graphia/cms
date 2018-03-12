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

	tag := e.Tag()

	switch tag {
	case "required":
		return "is a required field"
	case "email":
		return "is not a valid email address"
	case "min":

		// min is available to both string and int attributes, and we probably
		// want a different message for each
		switch e.Type().String() {
		case "int":
			return fmt.Sprintf("must be at least %s", e.Param())
		case "string":
			return fmt.Sprintf("must be at least %s characters", e.Param())
		}

	case "max":

		// min is available to both string and int attributes, and we probably
		// want a different message for each
		switch e.Type().String() {
		case "int":
			return fmt.Sprintf("must be smaller than %s", e.Param())
		case "string":
			return fmt.Sprintf("must be no more than %s characters", e.Param())
		}

	case "alphanumunicode":
		return "only alphanumeric unicode characters are permitted"
	}

	Warning.Println("No error message format specified", tag)
	return "there was an unspecified error"

}
