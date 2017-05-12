package main

import (
	"time"
)

// RepoWrite contains all info we need to perform a git commit
// TODO rename to FileWrite fw
type RepoWrite struct {
	Filename    string
	Path        string
	Body        string
	Message     string
	Name        string
	Email       string
	FrontMatter FrontMatter
}

// Response is a general response containing arbitrary data
type Response struct {
	Data string `json:"data"`
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

// FrontMatter contains the document's metadata
type FrontMatter struct {
	Title  string
	Author string
}

// Directory contains the directory's metadata
// FIXME eventually it will, currently just the name, need to
// work out how best to store it
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

// UserCredentials is the subset of User required for auth
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User holds all information specific to a user
type User struct {
	ID          int    `json:"id" storm:"id,increment"`
	Name        string `json:"name" validate:"required,min=3,max=64"`
	Username    string `json:"username" storm:"unique" validate:"required,min=3,max=32"`
	Password    []byte `json:"password" validate:"required,min=6"`
	Email       string `json:"email" storm:"unique" validate:"email,required"`
	Active      bool   `json:"active"`
	TokenString string `json:"token_string" storm:"unique"`
}

// LimitedUser is a 'safe' subset of user data that we can
// send out via the API. Password is omitted
type LimitedUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Token holds a JSON Web Token
type Token struct {
	Token string `json:"token"`
}
