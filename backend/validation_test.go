package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/go-playground/validator.v9"
)

type testStruct struct {
	Name  string `validate:"required,min=5,max=20"`
	Age   int    `validate:"omitempty,min=18,max=99"`
	Email string `validate:"omitempty,email"`
}

// should return a map[string][]string with attribute
// name and a slice of error messages
func TestValidationErrorsToJSON(t *testing.T) {
	validate := validator.New()

	// two failures; missing name and low age
	fail := testStruct{Age: 14}

	err := validate.Struct(fail)

	// there should be errors
	assert.NotNil(t, err)

	errors := validationErrorsToJSON(err)

	expectedOut := make(map[string]string)
	expectedOut["Name"] = "is a required field"
	expectedOut["Age"] = "must be at least 18"

	assert.Equal(t, expectedOut, errors)
}

func TestValidationErrorMessageRequiredTag(t *testing.T) {
	validate := validator.New()

	fail := testStruct{}

	err := validate.Struct(fail)

	// there should be errors
	assert.NotNil(t, err)

	errors := validationErrorsToJSON(err)

	expectedOut := make(map[string]string)
	expectedOut["Name"] = "is a required field"

	assert.Equal(t, expectedOut, errors)
}

func TestValidationErrorMessageMinLengthTag(t *testing.T) {
	validate := validator.New()

	fail := testStruct{Name: "xy"}

	err := validate.Struct(fail)

	// there should be errors
	assert.NotNil(t, err)

	errors := validationErrorsToJSON(err)

	expectedOut := make(map[string]string)
	expectedOut["Name"] = "must be at least 5 characters"

	assert.Equal(t, expectedOut, errors)
}

func TestValidationErrorMessageMaxLengthTag(t *testing.T) {
	validate := validator.New()

	fail := testStruct{Name: "abcdefghijklmnopqrstuvwxyz"}

	err := validate.Struct(fail)

	// there should be errors
	assert.NotNil(t, err)

	errors := validationErrorsToJSON(err)

	expectedOut := make(map[string]string)
	expectedOut["Name"] = "must be no more than 20 characters"

	assert.Equal(t, expectedOut, errors)
}

func TestValidationErrorMessageMinIntTag(t *testing.T) {
	validate := validator.New()

	fail := testStruct{Name: "Nelson Muntz", Age: 14}

	err := validate.Struct(fail)

	// there should be errors
	assert.NotNil(t, err)

	errors := validationErrorsToJSON(err)

	expectedOut := make(map[string]string)
	expectedOut["Age"] = "must be at least 18"

	assert.Equal(t, expectedOut, errors)
}

func TestValidationErrorMessageEmail(t *testing.T) {
	validate := validator.New()

	fail := testStruct{Email: "something@"}

	err := validate.Struct(fail)

	// there should be errors
	assert.NotNil(t, err)

	errors := validationErrorsToJSON(err)

	expectedOut := make(map[string]string)
	expectedOut["Name"] = "is a required field"
	expectedOut["Email"] = "is not a valid email address"

	assert.Equal(t, expectedOut, errors)
}

func TestValidationErrorMessageMaxIntTag(t *testing.T) {
	validate := validator.New()

	fail := testStruct{Name: "Kearney Zzyzwicz", Age: 130}

	err := validate.Struct(fail)

	// there should be errors
	assert.NotNil(t, err)

	errors := validationErrorsToJSON(err)

	expectedOut := make(map[string]string)
	expectedOut["Age"] = "must be smaller than 99"

	assert.Equal(t, expectedOut, errors)
}
