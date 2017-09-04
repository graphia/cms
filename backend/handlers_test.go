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

	"github.com/graphia/particle"
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

	var directoryPaths []string
	for _, directory := range receiver {
		directoryPaths = append(directoryPaths, directory.Path)
	}

	assert.Equal(t, directoryPaths, directoriesExpected)

}

func TestApiListDirectoriesHandlerNoDirectory(t *testing.T) {
	server = httptest.NewServer(protectedRouter())

	testConfigPath := "../config/test.yml"
	config, _ = loadConfig(&testConfigPath)
	config.Repository = "../tests/tmp/repositories/non_existant_repo"

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories")

	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)

	resp, _ := client.Do(req)

	var fr FailureResponse

	json.NewDecoder(resp.Body).Decode(&fr)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "No repository found", fr.Message)

}

func TestApiListDirectoriesHandlerNoGit(t *testing.T) {

	// remove any existing repo
	repoPath := "../tests/tmp/repositories/uninitialized"
	_ = os.RemoveAll(repoPath)

	// copy the repo template to the expected location
	// but don't initialise the git repo
	template := "../tests/backend/repositories/small"

	_ = CopyDir(template, repoPath)

	testConfigPath := "../config/test.yml"
	config, _ = loadConfig(&testConfigPath)
	config.Repository = repoPath

	server = httptest.NewServer(protectedRouter())

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories")

	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)

	resp, _ := client.Do(req)

	var fr FailureResponse

	json.NewDecoder(resp.Body).Decode(&fr)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Not a git repository", fr.Message)

}

func Test_apiListFilesInDirectoryHandler(t *testing.T) {
	server = httptest.NewServer(protectedRouter())
	repoPath := "../tests/tmp/repositories/list_files_in_directory"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files")

	req, _ := http.NewRequest("GET", target, nil)
	client := &http.Client{}
	resp, _ := client.Do(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var receiver []FileItem
	json.NewDecoder(resp.Body).Decode(&receiver)

	// ensure we get 3 files back
	assert.Equal(t, 3, len(receiver))

	// ensure all files are returned
	var expectedFilenames, actualFilenames []string

	expectedFilenames = []string{"document_1.md", "document_2.md", "document_3.md"}
	for _, fi := range receiver {
		actualFilenames = append(actualFilenames, fi.Filename)
	}
	assert.Equal(t, expectedFilenames, actualFilenames)

	// ensure frontmatter is returned correctly
	var expectedTitles, actualTitles []string

	expectedTitles = []string{"document 1", "document 2", "document 3"}
	for _, fi := range receiver {
		actualTitles = append(actualTitles, fi.FrontMatter.Title)
	}
	assert.Equal(t, expectedTitles, actualTitles)

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

	ncd := NewCommitDirectory{
		Path: "leftorium",
		DirectoryInfo: DirectoryInfo{
			Name:        "The Leftorium",
			Description: "Left-handed goods for all!",
		},
	}

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
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, ncd.Path, ".info"))
	assert.Contains(t, string(contents), fmt.Sprintf("%s: %s", "name", ncd.DirectoryInfo.Name))
	assert.Contains(t, string(contents), fmt.Sprintf("%s: %s", "description", ncd.DirectoryInfo.Description))

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

func TestApiCreateFileInDirectoryWithErrors(t *testing.T) {
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
		// No Message
		Files: []NewCommitFile{ncf},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	b := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, b)

	resp, err := client.Do(req)

	errors := make(map[string]string)

	json.NewDecoder(resp.Body).Decode(&errors)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, "is a required field", errors["message"])

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

func TestApiUpdateFileInDirectoryWithErrors(t *testing.T) {
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
		// No Message
		Files: []NewCommitFile{ncf},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	buff := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("PATCH", target, buff)

	resp, err := client.Do(req)

	errors := make(map[string]string)

	json.NewDecoder(resp.Body).Decode(&errors)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, "is a required field", errors["message"])

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

	raw, _ := ioutil.ReadFile(filepath.Join(
		config.Repository,
		"documents",
		"document_3.md",
	))

	contents, err := particle.YAMLEncoding.DecodeString(string(raw), &FrontMatter{})

	assert.Equal(t, *file.Markdown, string(contents))

}

func TestApiGetAttachmentsHandler(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/get_attachments_handler"
	setupMultipleFiletypesTestRepo(repoPath)

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/directories/appendices/files/appendix_1/attachments",
	)

	resp, err := http.Get(target)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	if err != nil {
		panic(err)
	}

	var attachments []Attachment

	json.NewDecoder(resp.Body).Decode(&attachments)

	jsonAttachmentContents, _ := ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "data", "data.json"))
	xmlAttachmentContents, _ := ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "data", "data.xml"))
	pngAttachmentContents, _ := ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "images", "image_1.png"))
	jpegAttachmentContents, _ := ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "images", "image_2.jpg"))

	expectedAttachments := []Attachment{
		Attachment{
			Path:             "appendices/appendix_1/data",
			AbsoluteFilename: "appendices/appendix_1/data/data.json",
			Extension:        ".json",
			MediaType:        "text/json",
			Data:             base64.StdEncoding.EncodeToString(jsonAttachmentContents),
			Filename:         "data.json",
		},
		Attachment{
			Path:             "appendices/appendix_1/data",
			AbsoluteFilename: "appendices/appendix_1/data/data.xml",
			Extension:        ".xml",
			MediaType:        "text/xml",
			Data:             base64.StdEncoding.EncodeToString(xmlAttachmentContents),
			Filename:         "data.xml",
		},
		Attachment{
			Path:             "appendices/appendix_1/images",
			AbsoluteFilename: "appendices/appendix_1/images/image_1.png",
			Extension:        ".png",
			MediaType:        "image/png",
			Data:             base64.StdEncoding.EncodeToString(pngAttachmentContents),
			Filename:         "image_1.png",
		},
		Attachment{
			Path:             "appendices/appendix_1/images",
			AbsoluteFilename: "appendices/appendix_1/images/image_2.jpg",
			Extension:        ".jpg",
			MediaType:        "image/jpeg",
			Data:             base64.StdEncoding.EncodeToString(jpegAttachmentContents),
			Filename:         "image_2.jpg",
		},
	}

	for _, a := range expectedAttachments {
		t.Run(a.Filename, func(t *testing.T) {
			assert.Contains(t, attachments, a)
		})
	}

}

func TestApiGetAttachmentsNoDirectoryHandler(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/get_attachments_handler"
	setupMultipleFiletypesTestRepo(repoPath)

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/directories/documents/files/document_3/attachments",
	)

	resp, _ := http.Get(target)

	var fr FailureResponse

	json.NewDecoder(resp.Body).Decode(&fr)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "No attachments", fr.Message)

}

// Setup tests

func Test_setupAllowInitializeRepository_Success(t *testing.T) {
	server = createTestServerWithContext()

	fullDirPath := "../tests/tmp/repositories/full"
	os.RemoveAll(fullDirPath)
	CopyDir("../tests/backend/repositories/small", fullDirPath)

	config.Repository = fullDirPath

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/setup/initialize_repository",
	)

	resp, _ := http.Get(target)
	var so SetupOption
	json.NewDecoder(resp.Body).Decode(&so)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, so.Enabled)

}

func Test_setupAllowInitializeRepository_Fail(t *testing.T) {
	server = createTestServerWithContext()

	gitRepoPath := "../tests/tmp/repositories/allow_initialize"
	_, _ = setupSmallTestRepo(gitRepoPath)

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/setup/initialize_repository",
	)

	resp, _ := http.Get(target)
	var so SetupOption
	json.NewDecoder(resp.Body).Decode(&so)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, so.Enabled)
	assert.Contains(t, so.Meta, "git repo already exists at")
}

func Test_setupInitializeRepository_Success(t *testing.T) {
	server = createTestServerWithContext()

	newDir := "../tests/tmp/repositories/full"
	os.RemoveAll(newDir)
	CopyDir("../tests/backend/repositories/small", newDir)

	config.Repository = newDir

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/setup/initialize_repository",
	)

	resp, _ := http.Post(target, "application/json", nil)

	var sr SuccessResponse
	json.NewDecoder(resp.Body).Decode(&sr)

	// make sure the returned oid is at the head of the repo
	repo, _ := repository(config)
	hc, _ := headCommit(repo)
	assert.Equal(t, hc.Id().String(), sr.Oid)

	// and that the status and message are correct
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, sr.Message, "Repository initialised")

}

func Test_setupInitializeRepository_Failure(t *testing.T) {
	server = createTestServerWithContext()

	gitRepoPath := "../tests/tmp/repositories/allow_initialize"
	_, _ = setupSmallTestRepo(gitRepoPath)

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/setup/initialize_repository",
	)

	resp, _ := http.Post(target, "application/json", nil)
	Debug.Println(resp)

	var fr FailureResponse
	json.NewDecoder(resp.Body).Decode(&fr)

	Debug.Println(resp.Body)

	// and that the status and message are correct
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, fr.Message, "Cannot initialize repository, see log")
}
