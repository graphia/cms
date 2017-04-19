package main

import (
	"time"
)

// RepoWrite contains all info we need to perform a git commit
// TODO rename to FileWrite fw
type RepoWrite struct {
	Filename string
	Path     string
	Body     string
	Message  string
	Name     string
	Email    string
}

// SuccessResponse contains information about a successful
// update to the repository
type SuccessResponse struct {
	Message string `json:"message"`
	Oid     string `json:"oid"`
}

// FailureResponse accompanies the HTTP status code with
// some more information as to why the update failed
type FailureResponse struct {
	Message string `json:"message"`
}

type FrontMatter struct {
	Title  string
	Author string
}

type Directory struct {
	Name string
}

// FileItem contains enough file information for listing
// HTML and raw Markdown content is omitted
type FileItem struct {
	AbsoluteFilename string    `json:"absolute_filename"`
	Filename         string    `json:"filename"`
	Path             string    `json:"path"`
	Author           string    `json:"author"`
	Date             time.Time `json:"updated_at"`
}

// File represents a Markdown file and can be returned with
// HTML or Markdown contents (or both if required)
type File struct {
	AbsoluteFilename string  `json:"absolute_filename"`
	Filename         string  `json:"filename"`
	Path             string  `json:"path"`
	HTML             *string `json:"html"`
	Markdown         *string `json:"markdown"`
	Author           string  `json:"author"`
	Title            string  `json:"title"`
}
