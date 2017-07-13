package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/libgit2/git2go.v25"
)

var (
	server *httptest.Server
)

// API Tests

func TestApiListDirectoriesHandler(t *testing.T) {
	server = httptest.NewServer(protectedRouter())

	repoPath := "../tests/tmp/repositories/list_directories"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories")

	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)

	resp, _ := client.Do(req)

	var receiver []Directory

	json.NewDecoder(resp.Body).Decode(&receiver)

	directoriesExpected := []string{"appendices", "documents"}

	var directoryNames []string
	for _, directory := range receiver {
		directoryNames = append(directoryNames, directory.Name)
	}

	assert.Equal(t, directoryNames, directoriesExpected)

}

func TestApiListDirectorySummaryHandler(t *testing.T) {

	server = httptest.NewServer(protectedRouter())

	// setup and send the request
	repoPath := "../tests/tmp/repositories/list_directory_summary"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/summary")

	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)

	resp, _ := client.Do(req)

	var receiver map[string][]FileItem
	json.NewDecoder(resp.Body).Decode(&receiver)

	// prepare the expected output
	documentFiles, _ := getFilesInDir("documents")
	appendicesFiles, _ := getFilesInDir("appendices")

	expectedSummary := map[string][]FileItem{
		"documents":  documentFiles,
		"appendices": appendicesFiles,
	}

	assert.Equal(t, 2, len(receiver))
	assert.Equal(t, expectedSummary, receiver)

}

func TestApiCreateDirectory(t *testing.T) {
	server = httptest.NewServer(protectedRouter())

	repoPath := "../tests/tmp/repositories/create_directory"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories")

	ncf := NewCommitFile{Path: "bobbins"}

	nc := &NewCommit{
		Email:   "martin.prince@springfield.k12.us",
		Name:    "Martin Prince",
		Message: "Forty whacks with a wet noodle",
		Files:   []NewCommitFile{ncf},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	b := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, b)

	resp, err := client.Do(req)

	var receiver SuccessResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// ensure returned commit hash is hte same as the repo's head
	assert.Equal(t, receiver.Oid, hc.Id().String())

	// ensure the file exists and has the right content
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, ncf.Path, ".keep"))
	assert.Equal(t, "", string(contents))

	// ensure the most recent commit has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, nc.Name)
	assert.Equal(t, lastCommit.Committer().Email, nc.Email)

	// ensure that the commit message has been set correctly
	msg := fmt.Sprintf("Added directories: %s", ncf.Path)
	assert.Equal(t, lastCommit.Message(), msg)

}

func TestApiCreateFileInDirectory(t *testing.T) {
	server = httptest.NewServer(protectedRouter())

	repoPath := "../tests/tmp/repositories/create_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files")

	ncf := NewCommitFile{
		Path:      "documents",
		Filename:  "document_6.md",
		Extension: "md",
		Body:      "# The quick brown fox",
		FrontMatter: FrontMatter{
			Title:  "Document Six",
			Author: "Kent Brockman & Troy McClure",
		},
	}

	nc := &NewCommit{
		Email:   "martin.prince@springfield.k12.us",
		Name:    "Martin Prince",
		Message: "Forty whacks with a wet noodle",
		Files:   []NewCommitFile{ncf},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	b := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, b)

	resp, err := client.Do(req)

	var receiver SuccessResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// ensure returned commit hash is hte same as the repo's head
	assert.Equal(t, receiver.Oid, hc.Id().String())

	// ensure the file exists and has the right content
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, ncf.Path, ncf.Filename))
	assert.Contains(t, string(contents), ncf.Body)
	assert.Contains(t, string(contents), ncf.FrontMatter.Author)
	assert.Contains(t, string(contents), ncf.FrontMatter.Title)

	// ensure the most recent commit has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, nc.Name)
	assert.Equal(t, lastCommit.Committer().Email, nc.Email)

}

func TestApiUpdateFileInDirectory(t *testing.T) {
	server = httptest.NewServer(protectedRouter())

	repoPath := "../tests/tmp/repositories/update_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_3.md")

	ncf := NewCommitFile{
		Path:      "documents",
		Filename:  "document_3.md",
		Extension: "md",
		Body:      "# The quick brown fox",
		FrontMatter: FrontMatter{
			Title:  "Document Three",
			Author: "Timothy Lovejoy",
		},
	}

	nc := &NewCommit{
		Email:   "martin.prince@springfield.k12.us",
		Name:    "Martin Prince",
		Message: "Forty whacks with a wet noodle",
		Files:   []NewCommitFile{ncf},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	buff := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("PATCH", target, buff)

	resp, err := client.Do(req)

	var receiver SuccessResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// ensure returned commit hash is hte same as the repo's head
	assert.Equal(t, receiver.Oid, hc.Id().String())

	// ensure the file exists and has the right content
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, ncf.Path, ncf.Filename))
	assert.Contains(t, string(contents), ncf.Body)
	assert.Contains(t, string(contents), ncf.FrontMatter.Author)
	assert.Contains(t, string(contents), ncf.FrontMatter.Title)

	// ensure the most recent commit has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, nc.Name)
	assert.Equal(t, lastCommit.Committer().Email, nc.Email)

}

func TestApiDeleteFileFromDirectory(t *testing.T) {
	server = httptest.NewServer(protectedRouter())

	repoPath := "../tests/tmp/repositories/delete_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_2.md")

	ncf1 := NewCommitFile{
		Filename: "document_1.md",
		Path:     "documents",
	}

	ncf2 := NewCommitFile{
		Filename: "document_2.md",
		Path:     "documents",
	}

	nc := &NewCommit{
		Email:   "clancy.wiggum@springfield.police.gov",
		Name:    "Clarence Wiggum",
		Message: "Suspect is hatless. Repeat, hatless.",
		Files:   []NewCommitFile{ncf1, ncf2},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	buff := bytes.NewBuffer(payload)

	req, _ := http.NewRequest("DELETE", target, buff)

	resp, err := client.Do(req)

	var receiver SuccessResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// ensure returned commit hash is hte same as the repo's head
	assert.Equal(t, receiver.Oid, hc.Id().String())

	// ensure the most recent nc has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, nc.Name)
	assert.Equal(t, lastCommit.Committer().Email, nc.Email)

	// TODO ensure files aren't present on the filesystem
	// ensure the file exists and has the right content
	_, err = os.Stat(filepath.Join(repoPath, ncf1.Path, ncf1.Filename))
	_, err = os.Stat(filepath.Join(repoPath, ncf2.Path, ncf2.Filename))
	assert.True(t, os.IsNotExist(err))

}

func TestApiDeleteDirectory(t *testing.T) {
	server = httptest.NewServer(protectedRouter())

	var err error

	repoPath := "../tests/tmp/repositories/delete_dir"
	setupSmallTestRepo(repoPath)

	directory := "appendices"

	target := fmt.Sprintf("%s/%s/%s", server.URL, "api/directories", directory)

	nc := &NewCommit{
		Email: "julius.hibbert@springfield-hospital.com",
		Name:  "Julius Hibbert",
	}

	// before deleting, make sure appendix files are present
	_, err = os.Stat(filepath.Join(repoPath, directory, "appendix_1.md"))
	assert.False(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(repoPath, directory, "appendix_2.md"))
	assert.False(t, os.IsNotExist(err))

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	buff := bytes.NewBuffer(payload)

	req, _ := http.NewRequest("DELETE", target, buff)

	resp, err := client.Do(req)

	var receiver SuccessResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// ensure returned commit hash is hte same as the repo's head
	assert.Equal(t, receiver.Oid, hc.Id().String())

	// ensure the most recent nc has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, nc.Name)
	assert.Equal(t, lastCommit.Committer().Email, nc.Email)

	// ensure the file exists and has the right content
	_, err = os.Stat(filepath.Join(repoPath, directory, "appendix_1.md"))
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(repoPath, directory, "appendix_2.md"))
	assert.True(t, os.IsNotExist(err))

}

func TestApiDeleteDirectoryNotExists(t *testing.T) {
	server = httptest.NewServer(protectedRouter())

	var err error

	repoPath := "../tests/tmp/repositories/delete_dir"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s/%s", server.URL, "api/directories", "favourites")

	nc := &NewCommit{
		Email: "julius.hibbert@springfield-hospital.com",
		Name:  "Julius Hibbert",
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	buff := bytes.NewBuffer(payload)

	req, _ := http.NewRequest("DELETE", target, buff)

	resp, err := client.Do(req)

	var receiver FailureResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	assert.Contains(t, receiver.Message, "Failed to delete directory")

}

func TestApiGetFileInDirectory(t *testing.T) {
	server = httptest.NewServer(protectedRouter())

	repoPath := "../tests/tmp/repositories/create_directory"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/directories/documents/files/document_3.md",
	)

	resp, err := http.Get(target)
	if err != nil {
		panic(err)
	}

	var file File

	json.NewDecoder(resp.Body).Decode(&file)

	assert.Equal(t, file.Filename, "document_3.md")
	assert.Equal(t, file.Path, "documents")
	assert.Contains(t, *file.HTML, "<h1>Document 3</h1>")

}

func TestApiEditFileInDirectory(t *testing.T) {
	repoPath := "../tests/tmp/repositories/create_directory"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/directories/documents/files/document_3.md/edit",
	)

	resp, err := http.Get(target)
	if err != nil {
		panic(err)
	}

	var file File

	json.NewDecoder(resp.Body).Decode(&file)

	assert.Equal(t, file.Filename, "document_3.md")
	assert.Equal(t, file.Path, "documents")

	contents, _ := ioutil.ReadFile(filepath.Join(
		config.Repository,
		"documents",
		"document_3.md",
	))

	assert.Equal(t, *file.Markdown, string(contents))

}
