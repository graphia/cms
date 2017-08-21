package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewCommitValidation(t *testing.T) {

	tests := []struct {
		name    string
		message string
		want    string
	}{
		{name: "Too short", message: "a", want: "must be at least 5 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc := NewCommit{Message: tt.message}
			err := validate.Struct(nc)
			json := validationErrorsToJSON(err)
			assert.Contains(t, json["Message"], tt.want)
		})
	}

}
