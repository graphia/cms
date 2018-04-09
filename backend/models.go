package main

import (
	"path/filepath"
	"time"

	"github.com/graphia/particle"
	"gopkg.in/libgit2/git2go.v25"
)

// NewCommit will replace RepoWrite and allow multiple files
type NewCommit struct {
	Message        string               `json:"message" validate:"required,min=5"`
	Files          []NewCommitFile      `json:"files"`
	Directories    []NewCommitDirectory `json:"directories"`
	RepositoryInfo `json:"repository_info"`
}

// NewCommitDirectory holds directory info for creating new dirs
type NewCommitDirectory struct {
	Path          string        `json:"name"`
	DirectoryInfo DirectoryInfo `json:"info"`
}

// NewCommitFile will replace RepoWrite's file attributes
type NewCommitFile struct {
	Filename      string      `json:"filename" validate:"required"`
	Document      string      `json:"document" validate:"required"`
	Path          string      `json:"path" validate:"required"`
	Body          string      `json:"body"`
	FrontMatter   FrontMatter `json:"frontmatter"`
	Base64Encoded bool        `json:"base_64_encoded"`
}

// ToMarkdown returns the file's actual contents, frontmatter and document body
func (ncf NewCommitFile) ToMarkdown() (b []byte) {
	md := []byte(ncf.Body)
	b = make([]byte, particle.YAMLEncoding.EncodeLen(md, ncf.FrontMatter))
	particle.YAMLEncoding.Encode(b, md, ncf.FrontMatter)
	return b
}

// NewTranslation creates a new copy of a file ready for translation
type NewTranslation struct {
	SourceFilename string `json:"source_filename" validate:"required"`
	SourceDocument string `json:"source_document" validate:"required"`
	Path           string `json:"path" validate:"required"`
	LanguageCode   string `json:"language_code" validate:"required"`
	RepositoryInfo `json:"repository_info"`
}

// TargetFilename provides the new filename with the
// language code inserted
// FIXME this should be strengthened so it will work with already-translated files too
func (nt NewTranslation) TargetFilename() string {
	return translationFilename(nt.SourceFilename, nt.LanguageCode)
}

// FrontMatter contains the document's metadata
//
// date is a string because we never actually *do* anything with it
// other than read/write it, and Go has no sane 'Date' type
// https://github.com/golang/go/issues/21365
type FrontMatter struct {
	Author   string   `json:"author"         yaml:"author"`
	Date     string   `json:"date,omitempty" yaml:"date"`
	Draft    bool     `json:"draft"          yaml:"draft"`
	Synopsis string   `json:"synopsis"       yaml:"synopsis"`
	Tags     []string `json:"tags"           yaml:"tags"`
	Title    string   `json:"title"          yaml:"title"`
	Version  string   `json:"version"        yaml:"version"`
}

// Directory contains the directory's metadata
// FIXME eventually it will, currently just the name, need to
// work out how best to store it
type Directory struct {
	Path           string `json:"path" yaml:"path"`
	DirectoryInfo  `json:"info"`
	RepositoryInfo `json:"repository_info,omitempty"`
}

// DirectorySummary contains the directory's metadata plus
// an array of its contents
type DirectorySummary struct {
	Path          string `json:"path"`
	DirectoryInfo `json:"info"`
	Contents      []FileItem `json:"contents"`
}

// DirectoryInfo contains the fields that will be written to
// a directory's .info file
type DirectoryInfo struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Body        string `json:"body" yaml:"-"`
	HTML        string `json:"html" yaml:"-"`
}

// RepositoryInfo provides some data about the repo, such as
// the latest revision
// validated as required because when sent with a NewCommit
// we need to ensure that we're working from an up-to-date tree
type RepositoryInfo struct {
	LatestRevision string `json:"latest_revision" validate:"required"`
}

// ServerInfo holds basic details/counts displayed on the dashboard
type ServerInfo struct {
	Title   string         `json:"title"`
	Users   int            `json:"users"`
	Commits int            `json:"commits"`
	Counts  map[string]int `json:"files"`
}

// Language contains a language's name and code for localisation
type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}

// FileItem contains enough file information for listing
// HTML and raw Markdown content is omitted
type FileItem struct {
	Filename    string      `json:"filename"`
	Path        string      `json:"path"`
	Document    string      `json:"document"`
	Date        time.Time   `json:"updated_at"`
	FrontMatter FrontMatter `json:"frontmatter"`
}

// File represents a Markdown file and can be returned with
// HTML or Markdown contents (or both if required)
type File struct {
	Filename       string          `json:"filename"`
	Path           string          `json:"path"`
	Document       string          `json:"document"`
	Language       string          `json:"language"`
	HTML           *string         `json:"html"`
	Markdown       *string         `json:"markdown"`
	FrontMatter    FrontMatter     `json:"frontmatter"`
	DirectoryInfo  *DirectoryInfo  `json:"directory_info,omitempty"`
	RepositoryInfo *RepositoryInfo `json:"repository_info,omitempty"`
	Translations   []string        `json:"translations"`
}

// FullPath constructs the absolute path using the path, document and filename
func (f File) FullPath() string {
	return filepath.Join(f.Path, f.Document, f.Filename)
}

// ToMarkdown returns the file's actual contents, frontmatter and document body
func (f File) ToMarkdown() (b []byte) {
	md := []byte(*f.Markdown)
	b = make([]byte, particle.YAMLEncoding.EncodeLen(md, f.FrontMatter))
	particle.YAMLEncoding.Encode(b, md, f.FrontMatter)
	return b
}

// Attachment belongs to a File, usually an image
type Attachment struct {
	Path      string `json:"path"`
	Filename  string `json:"filename"`
	Extension string `json:"extension"`
	MediaType string `json:"filetype"`
	Data      string `json:"data"`
}

// Token holds a JSON Web Token
type Token struct {
	Token string `json:"token"`
}

// SetupOption indicates whether or not to display initial setup screen
type SetupOption struct {
	Enabled bool   `json:"enabled"`
	Meta    string `json:"meta,omitempty"`
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
	Old        string         `json:"old"`
	New        string         `json:"new"`
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
