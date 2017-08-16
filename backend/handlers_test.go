package main

import (
	"bytes"
	"encoding/base64"
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

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
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/create_directory"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories")

	ncd := NewCommitDirectory{Path: "bobbins"}

	nc := &NewCommit{
		Message:     "Forty whacks with a wet noodle",
		Directories: []NewCommitDirectory{ncd},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	b := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, b)

	resp, err := client.Do(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var receiver SuccessResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// ensure returned commit hash is hte same as the repo's head
	assert.Equal(t, receiver.Oid, hc.Id().String())

	// ensure the file exists and has the right content
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, ncd.Path, ".keep"))
	assert.Equal(t, "", string(contents))

	// ensure the most recent commit has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	user := apiTestUser()
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)

	// ensure that the commit message has been set correctly
	msg := fmt.Sprintf("Added directories: %s", ncd.Path)
	assert.Equal(t, lastCommit.Message(), msg)

}

func TestApiCreateFileInDirectory(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/create_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files")

	ncf := NewCommitFile{
		Path:     "documents",
		Filename: "document_6.md",
		Body:     "# The quick brown fox",
		FrontMatter: FrontMatter{
			Title:  "Document Six",
			Author: "Kent Brockman & Troy McClure",
		},
	}

	nc := &NewCommit{
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

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

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

	user := apiTestUser()
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)

}

func TestApiCreateImageFileInDirectory(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/create_image_file"
	setupMultipleFiletypesTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files")

	pngImage, _ := ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "images", "image_1.png"))

	ncf := NewCommitFile{
		Path:          "documents/document_1/images",
		Filename:      "image_4.png",
		Base64Encoded: true,
		Body:          base64.StdEncoding.EncodeToString(pngImage),
	}

	nc := &NewCommit{
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

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var receiver SuccessResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// ensure returned commit hash is hte same as the repo's head
	assert.Equal(t, receiver.Oid, hc.Id().String())

	// ensure the file exists and has the right content
	_, err = os.Stat(filepath.Join(repoPath, ncf.Path, ncf.Filename))
	assert.False(t, os.IsNotExist(err))

	file, _ := ioutil.ReadFile(filepath.Join(repoPath, ncf.Path, ncf.Filename))
	assert.Equal(t, pngImage, file)

	// ensure the most recent commit has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)

	user := apiTestUser()
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)

}

func TestApiUpdateFileInDirectory(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/update_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_3.md")

	ncf := NewCommitFile{
		Path:     "documents",
		Filename: "document_3.md",
		Body:     "# The quick brown fox",
		FrontMatter: FrontMatter{
			Title:  "Document Three",
			Author: "Timothy Lovejoy",
		},
	}

	nc := &NewCommit{
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

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

	user := apiTestUser()
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)

}

// Make sure that the file specified in the URL is included in the payload
func TestApiUpdateOtherFileInDirectory(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/update_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_3.md")

	ncf := NewCommitFile{
		Path:     "documents",
		Filename: "document_2.md", // note, target contains document_2.md
		Body:     "# The quick brown fox",
		FrontMatter: FrontMatter{
			Title:  "Document Three",
			Author: "Timothy Lovejoy",
		},
	}

	nc := &NewCommit{
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

	var receiver FailureResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	assert.Equal(t, "No supplied file matches path", receiver.Message)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

// Make sure that at least one file is specified in the payload
func TestApiUpdateNoFilesInDirectory(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/update_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_3.md")

	nc := &NewCommit{
		Message: "Forty whacks with a wet noodle",
		Files:   []NewCommitFile{},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	buff := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("PATCH", target, buff)

	resp, err := client.Do(req)

	var receiver FailureResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	assert.Equal(t, "No files specified for update", receiver.Message)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func TestApiDeleteFileFromDirectory(t *testing.T) {
	server = createTestServerWithContext()

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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var receiver SuccessResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// ensure returned commit hash is hte same as the repo's head
	assert.Equal(t, receiver.Oid, hc.Id().String())

	// ensure the most recent nc has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)

	user := apiTestUser()
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)

	// ensure the files no longer exist
	_, err = os.Stat(filepath.Join(repoPath, ncf1.Path, ncf1.Filename))
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(repoPath, ncf2.Path, ncf2.Filename))
	assert.True(t, os.IsNotExist(err))

}

// Make sure that the file specified in the URL is included in the payload
func TestApiDeleteOtherFileFromDirectory(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/delete_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_3.md")

	ncf1 := NewCommitFile{
		Filename: "document_1.md",
		Path:     "documents",
	}

	ncf2 := NewCommitFile{
		Filename: "document_2.md",
		Path:     "documents",
	}

	nc := &NewCommit{
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

	var receiver FailureResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	assert.Equal(t, "No supplied file matches path", receiver.Message)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

// Make sure that at least one file is specified in the payload
func TestApiDeleteNoFilesFromDirectory(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/delete_file"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_3.md")

	nc := &NewCommit{
		Message: "Suspect is hatless. Repeat, hatless.",
		Files:   []NewCommitFile{},
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

	assert.Equal(t, "No files specified for deletion", receiver.Message)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func TestApiDeleteDirectory(t *testing.T) {
	server = createTestServerWithContext()

	var err error

	repoPath := "../tests/tmp/repositories/delete_dir"
	setupSmallTestRepo(repoPath)

	ncd := NewCommitDirectory{Path: "appendices"}
	nc := &NewCommit{
		Directories: []NewCommitDirectory{ncd},
	}

	target := fmt.Sprintf("%s/%s/%s", server.URL, "api/directories", ncd.Path)

	// before deleting, make sure appendix files are present
	_, err = os.Stat(filepath.Join(repoPath, ncd.Path, "appendix_1.md"))
	assert.False(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(repoPath, ncd.Path, "appendix_2.md"))
	assert.False(t, os.IsNotExist(err))

	payload, _ := json.Marshal(nc)

	client := &http.Client{}

	buff := bytes.NewBuffer(payload)

	req, _ := http.NewRequest("DELETE", target, buff)

	resp, err := client.Do(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var receiver SuccessResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// ensure returned commit hash is hte same as the repo's head
	assert.Equal(t, receiver.Oid, hc.Id().String())

	// ensure the most recent nc has the right name and email
	oid, _ := git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)
	user := apiTestUser()
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)

	// ensure the files no longer exist
	_, err = os.Stat(filepath.Join(repoPath, ncd.Path, "appendix_1.md"))
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(repoPath, ncd.Path, "appendix_2.md"))
	assert.True(t, os.IsNotExist(err))

}

// make sure error is returned when trying to delete a non-existant directory
func TestApiDeleteDirectoryNotExists(t *testing.T) {
	server = createTestServerWithContext()

	var err error

	repoPath := "../tests/tmp/repositories/delete_dir"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s/%s", server.URL, "api/directories", "favourites")

	ncd := NewCommitDirectory{Path: "favourites"}
	nc := &NewCommit{
		Directories: []NewCommitDirectory{ncd},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	buff := bytes.NewBuffer(payload)

	req, _ := http.NewRequest("DELETE", target, buff)

	resp, err := client.Do(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var receiver FailureResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	assert.Contains(t, receiver.Message, "Failed to delete directory")

}

// ensure correct error returned dir named in URL path isn't specified
// in the payload
func TestApiDeleteAnotherDirectory(t *testing.T) {
	server = createTestServerWithContext()

	var err error

	repoPath := "../tests/tmp/repositories/delete_dir"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s/%s", server.URL, "api/directories", "appendices")

	ncd := NewCommitDirectory{Path: "documents"}
	nc := &NewCommit{
		Directories: []NewCommitDirectory{ncd},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	buff := bytes.NewBuffer(payload)

	req, _ := http.NewRequest("DELETE", target, buff)

	resp, err := client.Do(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var receiver FailureResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	assert.Contains(t, receiver.Message, "No specified directory matches path")

}

func TestApiGetFileInDirectory(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/create_directory"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/directories/documents/files/document_3.md",
	)

	resp, _ := http.Get(target)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var file File

	json.NewDecoder(resp.Body).Decode(&file)

	assert.Equal(t, file.Filename, "document_3.md")
	assert.Equal(t, file.Path, "documents")
	assert.Contains(t, *file.HTML, "<h1>Document 3</h1>")

}

func TestApiEditFileInDirectory(t *testing.T) {

	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/create_directory"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/directories/documents/files/document_3.md/edit",
	)

	resp, err := http.Get(target)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

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
