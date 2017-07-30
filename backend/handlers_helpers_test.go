package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isImageURI(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// imagey things
		{name: ".jpg", args: args{uri: "https://server/files/file.jpg"}, want: true},
		{name: ".gif", args: args{uri: "https://server/files/file.png"}, want: true},
		{name: ".png", args: args{uri: "https://server/files/file.gif"}, want: true},

		// non-imagey things
		{name: ".html", args: args{uri: "https://server/files/file.html"}, want: false},
		{name: ".md", args: args{uri: "https://server/files/file.md"}, want: false},
		{name: ".pdf", args: args{uri: "https://server/files/file.pdf"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isImageURI(tt.args.uri))
		})
	}
}

func Test_extractImagePath(t *testing.T) {

	originalPath := "/cms/documents/students.md/images/jimbo.jpg"
	extractedPath := extractImagePath(originalPath)

	// ensure 'cms' segment is removed
	assert.NotContains(t, extractedPath, "cms")

	// ensure the markdown extension is removed
	assert.NotContains(t, extractedPath, ".md")

	// ensure it contains the repository prefix and it's at the beginning
	assert.Contains(t, extractedPath, config.Repository)
	assert.True(t, strings.HasPrefix(extractedPath, config.Repository))

}
