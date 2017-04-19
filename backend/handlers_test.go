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

func init() {
	server = httptest.NewServer(setupRouter())
}

func TestApiCreateDirectory(t *testing.T) {
	repoPath := "../tests/tmp/repositories/create_directory"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories")

	rw := &RepoWrite{
		Email:   "martin.prince@springfield.k12.us",
		Name:    "Martin Prince",
		Message: "Forty whacks with a wet noodle",
		Path:    "bobbins",
	}

	payload, err := json.Marshal(rw)
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
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, rw.Path, ".keep"))
	assert.Equal(t, string(contents), "")

	// ensure the most recent commit has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, rw.Name)
	assert.Equal(t, lastCommit.Committer().Email, rw.Email)

	// ensure that the commit message has been set correctly
	msg := fmt.Sprintf("Added %s directory", rw.Path)
	assert.Equal(t, lastCommit.Message(), msg)

}

func TestApiCreateFileInDirectory(t *testing.T) {

	repoPath := "../tests/tmp/repositories/create_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files")

	rw := &RepoWrite{
		Body:     "# The quick brown fox",
		Email:    "martin.prince@springfield.k12.us",
		Name:     "Martin Prince",
		Message:  "Forty whacks with a wet noodle",
		Path:     "documents",
		Filename: "document_6.md",
	}

	payload, err := json.Marshal(rw)
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
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, rw.Path, rw.Filename))
	assert.Equal(t, string(contents), rw.Body)

	// ensure the most recent commit has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, rw.Name)
	assert.Equal(t, lastCommit.Committer().Email, rw.Email)

}

func TestApiUpdateFileInDirectory(t *testing.T) {

	repoPath := "../tests/tmp/repositories/update_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_3.md")

	rw := &RepoWrite{
		Body:     "# The quick brown fox",
		Email:    "martin.prince@springfield.k12.us",
		Name:     "Martin Prince",
		Message:  "Forty whacks with a wet noodle",
		Path:     "documents",
		Filename: "document_3.md",
	}

	payload, err := json.Marshal(rw)
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
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, rw.Path, rw.Filename))
	assert.Equal(t, string(contents), rw.Body)

	// ensure the most recent commit has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, rw.Name)
	assert.Equal(t, lastCommit.Committer().Email, rw.Email)

}

func TestApiDeleteFileFromDirectory(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_2.md")

	rw := &RepoWrite{
		Email:    "clancy.wiggum@springfield.police.gov",
		Name:     "Clarence Wiggum",
		Message:  "Suspect is hatless. Repeat, hatless.",
		Filename: "document_2.md",
		Path:     "documents",
	}

	payload, err := json.Marshal(rw)
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

	// ensure the most recent rw has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, rw.Name)
	assert.Equal(t, lastCommit.Committer().Email, rw.Email)

	// TODO ensure file isn't present on the filesystem
	// ensure the file exists and has the right content
	_, err = os.Stat(filepath.Join(repoPath, rw.Path, rw.Filename))
	assert.True(t, os.IsNotExist(err))

}

func TestApiDeleteDirectory(t *testing.T) {

	var err error

	repoPath := "../tests/tmp/repositories/delete_dir"
	setupSmallTestRepo(repoPath)

	directory := "appendices"

	target := fmt.Sprintf("%s/%s/%s", server.URL, "api/directories", directory)

	rw := &RepoWrite{
		Email: "julius.hibbert@springfield-hospital.com",
		Name:  "Julius Hibbert",
	}

	// before deleting, make sure appendix files are present
	_, err = os.Stat(filepath.Join(repoPath, directory, "appendix_1.md"))
	assert.False(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(repoPath, directory, "appendix_2.md"))
	assert.False(t, os.IsNotExist(err))

	payload, err := json.Marshal(rw)
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

	// ensure the most recent rw has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, rw.Name)
	assert.Equal(t, lastCommit.Committer().Email, rw.Email)

	// ensure the file exists and has the right content
	_, err = os.Stat(filepath.Join(repoPath, directory, "appendix_1.md"))
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(repoPath, directory, "appendix_2.md"))
	assert.True(t, os.IsNotExist(err))

}

func TestApiDeleteDirectoryNotExists(t *testing.T) {

	var err error

	repoPath := "../tests/tmp/repositories/delete_dir"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s/%s", server.URL, "api/directories", "favourites")

	rw := &RepoWrite{
		Email: "julius.hibbert@springfield-hospital.com",
		Name:  "Julius Hibbert",
	}

	payload, err := json.Marshal(rw)
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
