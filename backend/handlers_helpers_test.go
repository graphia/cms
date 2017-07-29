package main

import (
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
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Trim Markdown Extension",
			args: args{uri: "/cms/documents/students.md/images/jimbo.jpg"},
			want: "repo/documents/students/images/jimbo.jpg",
		},
		{
			name: "Omit 'CMS'",
			args: args{uri: "/cms/documents/employees/images/lenny.png"},
			want: "repo/documents/employees/images/lenny.png",
		},
		{
			name: "Prefix with repository path",
			args: args{uri: "/cms/documents/pensioners/images/abe.gif"},
			want: "repo/documents/pensioners/images/abe.gif",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, extractImagePath(tt.args.uri))
		})
	}
}
