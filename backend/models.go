package main

import (
	"time"

	"gopkg.in/libgit2/git2go.v25"
)

// NewCommit will replace RepoWrite and allow multiple files
// TODO rename to NewCommit
type NewCommit struct {
	Message     string               `json:"message"`
	Files       []NewCommitFile      `json:"files"`
	Directories []NewCommitDirectory `json:"directories"`
}

// NewCommitDirectory holds directory info for creating new dirs
type NewCommitDirectory struct {
	Path string `json:"name"`
}

// NewCommitFile will replace RepoWrite's file attributes
type NewCommitFile struct {
	Filename      string      `json:"filename"`
	Path          string      `json:"path"`
	Body          string      `json:"body"`
	FrontMatter   FrontMatter `json:"frontmatter"`
	Base64Encoded bool        `json:"base_64_encoded"`
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
	Meta    string `json:"meta,omitempty"`
}

// FrontMatter contains the document's metadata
type FrontMatter struct {
	Author   string   `yaml:"author"`
	Slug     string   `yaml:"slug"`
	Synopsis string   `yaml:"synopsis"`
	Tags     []string `yaml:"tags"`
	Title    string   `yaml:"title"`
	Version  string   `yaml:"version"`
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
	Synopsis         string    `json:"synopsis"`
	Version          string    `json:"version"`
	Tags             []string  `json:"tags"`
	Title            string    `json:"title"`
}

// File represents a Markdown file and can be returned with
// HTML or Markdown contents (or both if required)
type File struct {
	AbsoluteFilename     string   `json:"absolute_filename"`
	Filename             string   `json:"filename"`
	AttachmentsDirectory string   `json:"attachments_directory"`
	Path                 string   `json:"path"`
	HTML                 *string  `json:"html"`
	Markdown             *string  `json:"markdown"`
	Author               string   `json:"author"`
	Title                string   `json:"title"`
	Synopsis             string   `json:"synopsis"`
	Version              string   `json:"version"`
	Tags                 []string `json:"tags"`
}

// Attachment belongs to a File, usually an image
type Attachment struct {
	Path             string `json:"path"`
	Filename         string `json:"filename"`
	AbsoluteFilename string `json:"absolute_filename"`
	Extension        string `json:"extension"`
	MediaType        string `json:"filetype"`
	Data             string `json:"data"`
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
	Password    string `json:"password" validate:"required,min=6"`
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

// InitialSetup indicates whether or not to display initial setup screen
type InitialSetup struct {
	Enabled bool `json:"enabled"`
}

// Commit holds metadata for a Git Commit
type Commit struct {
	Message    string         `json:"message"`
	ID         string         `json:"id"`
	ObjectType string         `json:"object_type"`
	Author     *git.Signature `json:"author"`
	Time       time.Time      `json:"timestamp"`
}

// HistoricCommit is a commit used as part of a log
type HistoricCommit struct {
	EntryID    string         `json:"entry"`
	Message    string         `json:"message"`
	ID         string         `json:"id"`
	ObjectType string         `json:"object_type"`
	Author     *git.Signature `json:"author"`
	Time       time.Time      `json:"timestamp"`
}

// Changeset holds data about a previous commit, including the full delta
type Changeset struct {
	NumDeltas  int                       `json:"num_deltas"`
	NumAdded   int                       `json:"num_added"`
	NumDeleted int                       `json:"num_deleted"`
	FullDiff   string                    `json:"full_diff"`
	Files      map[string]ChangesetFiles `json:"files"`
	Message    string                    `json:"message"`
	Author     *git.Signature            `json:"author"`
	Hash       string                    `json:"hash"`
	Time       time.Time                 `json:"timestamp"`
}

// ChangesetFiles holds a copy of the file before and after the change
type ChangesetFiles struct {
	Old string `json:"old,omitempty"`
	New string `json:"new,omitempty"`
}
