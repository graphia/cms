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

	type rec struct {
		Files         []FileItem    `json:"files"`
		DirectoryInfo DirectoryInfo `json:"info"`
	}

	var receiver rec

	json.NewDecoder(resp.Body).Decode(&receiver)

	// ensure directory info is correct
	assert.Equal(
		t,
		DirectoryInfo{
			Title:       "Documents",
			Description: "Documents go here",
		},
		receiver.DirectoryInfo,
	)

	// ensure we get 3 files back
	assert.Equal(t, 3, len(receiver.Files))

	// ensure all files are returned
	var expectedFilenames, actualFilenames []string

	expectedFilenames = []string{"document_1.md", "document_2.md", "document_3.md"}
	for _, fi := range receiver.Files {
		actualFilenames = append(actualFilenames, fi.Filename)
	}
	assert.Equal(t, expectedFilenames, actualFilenames)

	// ensure frontmatter is returned correctly
	var expectedTitles, actualTitles []string

	expectedTitles = []string{"document 1", "document 2", "document 3"}
	for _, fi := range receiver.Files {
		actualTitles = append(actualTitles, fi.FrontMatter.Title)
	}
	assert.Equal(t, expectedTitles, actualTitles)

}

func TestApiListDirectorySummaryHandler(t *testing.T) {

	server = httptest.NewServer(protectedRouter())

	// setup and send the request
	repoPath := "../tests/tmp/repositories/list_directory_summary"
	setupSubdirsTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/summary")

	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)

	resp, _ := client.Do(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var receiver []DirectorySummary
	json.NewDecoder(resp.Body).Decode(&receiver)

	// prepare the expected output
	documentFiles, _ := getFilesInDir("documents")
	appendicesFiles, _ := getFilesInDir("appendices")

	expectedSummary := []DirectorySummary{
		DirectorySummary{
			Path:          "appendices",
			DirectoryInfo: DirectoryInfo{},
			Contents:      appendicesFiles,
		},
		DirectorySummary{
			Path: "documents",
			DirectoryInfo: DirectoryInfo{
				Title:       "Documents",
				Description: "Documents go here",
			},
			Contents: documentFiles,
		},
	}

	assert.Equal(t, 2, len(receiver))
	assert.Equal(t, expectedSummary, receiver)

}

func Test_apiGetDirectoryMetadata(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/create_directory"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents")

	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)

	resp, _ := client.Do(req)

	var receiver DirectoryInfo

	json.NewDecoder(resp.Body).Decode(&receiver)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(
		t,
		DirectoryInfo{
			Title:       "Documents",
			Description: "Documents go here",
		},
		receiver,
	)
}

func TestApiCreateDirectoryHandler(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/create_directory"
	setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories")

	ncd := NewCommitDirectory{
		Path: "leftorium",
		DirectoryInfo: DirectoryInfo{
			Title:       "The Leftorium",
			Description: "Left-handed goods for all!",
			Body:        "# Hi-didly-ho, neighbourinos!",
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
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, ncd.Path, "_index.md"))

	assert.Contains(t, string(contents), fmt.Sprintf("%s: %s", "title", ncd.DirectoryInfo.Title))
	assert.Contains(t, string(contents), fmt.Sprintf("%s: %s", "description", ncd.DirectoryInfo.Description))
	assert.Contains(t, string(contents), ncd.DirectoryInfo.Body)

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

func TestApiCreateFileInDirectoryHandler(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/create_file"
	lr, _ := setupSmallTestRepo(repoPath)

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
		Message:        "Forty whacks with a wet noodle",
		Files:          []NewCommitFile{ncf},
		RepositoryInfo: RepositoryInfo{LatestRevision: lr.String()},
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

func TestApiCreateFileInDirectoryNoRepoInfo(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/create_file"
	_, _ = setupSmallTestRepo(repoPath)

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

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	actual := make(map[string]string)

	json.NewDecoder(resp.Body).Decode(&actual)

	expected := make(map[string]string)
	expected["LatestRevision"] = "is a required field"

	assert.Equal(t, expected, actual)

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

func TestApiCreateFileInDirectoryRepoOutOfDate(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/create_file"
	firstCommit, _ := setupSmallTestRepo(repoPath)
	repo, _ := repository(config)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files")

	ncf := NewCommitFile{
		Path:     "documents",
		Filename: "document_9.md",
		Body:     "# The quick brown fox",
		FrontMatter: FrontMatter{
			Title:  "Document Six",
			Author: "Kent Brockman & Troy McClure",
		},
	}

	nc := &NewCommit{
		Message:        "First Commit",
		Files:          []NewCommitFile{ncf},
		RepositoryInfo: RepositoryInfo{LatestRevision: firstCommit.String()},
	}

	// Insert another commit so firstCommit is no longer current
	_, _ = createRandomFile(repo, "document_5.md", "whoosh")

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	b := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", target, b)

	resp, err := client.Do(req)

	var fr FailureResponse

	json.NewDecoder(resp.Body).Decode(&fr)

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	assert.Contains(t, fr.Message, "Repository out of sync with commit")

}

func TestApiCreateImageFileInDirectory(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/create_image_file"
	lr, _ := setupMultipleFiletypesTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files")

	pngImage, _ := ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "images", "image_1.png"))

	ncf := NewCommitFile{
		Path:          "documents/document_1/images",
		Filename:      "image_4.png",
		Base64Encoded: true,
		Body:          base64.StdEncoding.EncodeToString(pngImage),
	}

	nc := &NewCommit{
		Message:        "Forty whacks with a wet noodle",
		Files:          []NewCommitFile{ncf},
		RepositoryInfo: RepositoryInfo{LatestRevision: lr.String()},
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
	lr, _ := setupSmallTestRepo(repoPath)

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
		Message:        "Forty whacks with a wet noodle",
		Files:          []NewCommitFile{ncf},
		RepositoryInfo: RepositoryInfo{LatestRevision: lr.String()},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	buff := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("PATCH", target, buff)

	resp, err := client.Do(req)

	// We've updated a file by creating a commit, so it's a 201
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

func TestApiUpdateFileInDirectoryRepoOutOfDate(t *testing.T) {
	server = createTestServerWithContext()
	repo, _ := repository(config)

	repoPath := "../tests/tmp/repositories/update_file"
	firstCommit, _ := setupSmallTestRepo(repoPath)

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
		Message:        "First Commit",
		Files:          []NewCommitFile{ncf},
		RepositoryInfo: RepositoryInfo{LatestRevision: firstCommit.String()},
	}

	// Insert another commit so firstCommit is no longer current
	_, _ = createRandomFile(repo, "document_5.md", "whoosh")

	payload, _ := json.Marshal(nc)

	buff := bytes.NewBuffer(payload)

	client := &http.Client{}

	req, _ := http.NewRequest("PATCH", target, buff)

	resp, _ := client.Do(req)

	var fr FailureResponse

	json.NewDecoder(resp.Body).Decode(&fr)

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	assert.Contains(t, fr.Message, "Repository out of sync with commit")

}

// Make sure that the file specified in the URL is included in the payload
func TestApiUpdateOtherFileInDirectory(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/update_file"
	lr, _ := setupSmallTestRepo(repoPath)

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
		Message:        "Forty whacks with a wet noodle",
		Files:          []NewCommitFile{ncf},
		RepositoryInfo: RepositoryInfo{LatestRevision: lr.String()},
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
	lr, _ := setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/directories/documents/files/document_3.md")

	nc := &NewCommit{
		Message:        "Forty whacks with a wet noodle",
		Files:          []NewCommitFile{},
		RepositoryInfo: RepositoryInfo{LatestRevision: lr.String()},
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
	oid, _ := setupSmallTestRepo(repoPath)

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
		Files:          []NewCommitFile{ncf1, ncf2},
		RepositoryInfo: RepositoryInfo{LatestRevision: oid.String()},
	}

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	buff := bytes.NewBuffer(payload)

	req, _ := http.NewRequest("DELETE", target, buff)

	resp, err := client.Do(req)

	// We've created a commit so return 201, even though the
	// commit contains a deletion
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var receiver SuccessResponse

	json.NewDecoder(resp.Body).Decode(&receiver)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// ensure returned commit hash is hte same as the repo's head
	assert.Equal(t, receiver.Oid, hc.Id().String())

	// ensure the most recent nc has the right name and email
	oid, _ = git.NewOid(receiver.Oid)
	lastCommit, _ := repo.LookupCommit(oid)

	user := apiTestUser()
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)

	// ensure the files no longer exist
	_, err = os.Stat(filepath.Join(repoPath, ncf1.Path, ncf1.Filename))
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(repoPath, ncf2.Path, ncf2.Filename))
	assert.True(t, os.IsNotExist(err))

	// ensure a suitable commit message has been generated
	assert.Equal(t, "File deleted documents/document_2.md", lastCommit.Message())

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
	lr, _ := setupSmallTestRepo(repoPath)

	ncd := NewCommitDirectory{Path: "appendices"}
	nc := &NewCommit{
		Directories:    []NewCommitDirectory{ncd},
		RepositoryInfo: RepositoryInfo{LatestRevision: lr.String()},
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

	// Created because we're creating a commit (despite deleting
	// a directory)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

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

func TestApiDeleteDirectoryRepoOutOfDate(t *testing.T) {
	server = createTestServerWithContext()

	var err error

	repoPath := "../tests/tmp/repositories/delete_dir"
	firstCommit, _ := setupSmallTestRepo(repoPath)
	repo, _ := repository(config)

	ncd := NewCommitDirectory{Path: "appendices"}
	nc := &NewCommit{
		Directories:    []NewCommitDirectory{ncd},
		RepositoryInfo: RepositoryInfo{LatestRevision: firstCommit.String()},
	}

	_, _ = createRandomFile(repo, "document_5.md", "whoosh")

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

	// Created because we're creating a commit (despite deleting
	// a directory)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var fr FailureResponse

	json.NewDecoder(resp.Body).Decode(&fr)
	assert.Contains(t, fr.Message, "Repository out of sync with commit")

	// make sure dirs haven't actually been deleted
	_, err = os.Stat(filepath.Join(repoPath, ncd.Path, "appendix_1.md"))
	assert.False(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(repoPath, ncd.Path, "appendix_2.md"))
	assert.False(t, os.IsNotExist(err))

}

func TestApiDeleteDirectoryNoRepoInfo(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/delete_dir"
	setupSmallTestRepo(repoPath)

	ncd := NewCommitDirectory{Path: "appendices"}
	nc := &NewCommit{
		Directories: []NewCommitDirectory{ncd},
	}

	target := fmt.Sprintf("%s/%s/%s", server.URL, "api/directories", ncd.Path)

	payload, _ := json.Marshal(nc)

	client := &http.Client{}

	buff := bytes.NewBuffer(payload)

	req, _ := http.NewRequest("DELETE", target, buff)

	resp, _ := client.Do(req)

	// Created because we're creating a commit (despite deleting
	// a directory)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var fr FailureResponse
	json.NewDecoder(resp.Body).Decode(&fr)

	assert.Contains(t, fr.Message, "No hash provided")

}

// make sure error is returned when trying to delete a non-existant directory
func TestApiDeleteDirectoryNotExists(t *testing.T) {
	server = createTestServerWithContext()
	var err error

	repoPath := "../tests/tmp/repositories/delete_dir"
	lr, _ := setupSmallTestRepo(repoPath)

	target := fmt.Sprintf("%s/%s/%s", server.URL, "api/directories", "favourites")

	ncd := NewCommitDirectory{Path: "favourites"}
	nc := &NewCommit{
		Directories:    []NewCommitDirectory{ncd},
		RepositoryInfo: RepositoryInfo{LatestRevision: lr.String()},
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
	assert.Equal(t, &DirectoryInfo{Title: "Documents", Description: "Documents go here"}, file.DirectoryInfo)

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

func TestApiUpdateDirectoriesHandler(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/update_directories"
	setupSmallTestRepo(repoPath)

	nc := NewCommit{
		Message: "Added beverage information",
		Directories: []NewCommitDirectory{
			NewCommitDirectory{
				Path: "documents", // there is already a _index.md in documents
				DirectoryInfo: DirectoryInfo{
					Title:       "Buzz Cola",
					Description: "Twice the sugar, twice the caffeine",
					Body:        "# Buzz Cola\nThe taste you'll kill for!",
				},
			},
		},
	}

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/directories/documents",
	)

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

	// make sure the committer info is correct
	repo, _ := repository(config)
	oid, _ := git.NewOid(receiver.Oid)
	lc, _ := repo.LookupCommit(oid)

	user := apiTestUser()

	assert.Equal(t, lc.Committer().Name, user.Name)
	assert.Equal(t, lc.Committer().Email, user.Email)
	assert.Equal(t, lc.Message(), nc.Message)

	// finally ensure that the file has been written properly
	for _, f := range nc.Directories {
		contents, _ := ioutil.ReadFile(filepath.Join(repoPath, f.Path, "_index.md"))

		assert.Contains(t, string(contents), "---") // make sure yaml is delimited
		assert.Contains(t, string(contents), "title: Buzz Cola")
		assert.Contains(t, string(contents), "description: Twice the sugar, twice the caffeine")

		assert.Contains(t, string(contents), "# Buzz Cola")
		assert.Contains(t, string(contents), "The taste you'll kill for!")

	}

}

func TestApiUpdateDirectoriesHandlerNoDirectories(t *testing.T) {
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/update_directories"
	setupSmallTestRepo(repoPath)

	nc := NewCommit{
		Message:     "Added beverage information",
		Directories: []NewCommitDirectory{},
	}

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/directories/documents",
	)

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	buff := bytes.NewBuffer(payload)
	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", target, buff)
	resp, err := client.Do(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var receiver FailureResponse
	json.NewDecoder(resp.Body).Decode(&receiver)

	// make sure the error message is correct
	assert.Equal(t, "No directories specified for updates", receiver.Message)

}

func TestApiUpdateDirectoriesHandlerWrongDirectories(t *testing.T) {
	// in this example the path doesn't match the supplied dirs
	server = createTestServerWithContext()

	repoPath := "../tests/tmp/repositories/update_directories"
	setupSmallTestRepo(repoPath)

	nc := NewCommit{
		Message: "Added beverage information",
		Directories: []NewCommitDirectory{
			NewCommitDirectory{
				Path: "documents", // there is already a _index.md in documents
				DirectoryInfo: DirectoryInfo{
					Title:       "Buzz Cola",
					Description: "Twice the sugar, twice the caffeine",
				},
			},
		},
	}

	target := fmt.Sprintf(
		"%s/%s",
		server.URL,
		"api/directories/appendices", // note that in nc we're changing 'documents'
	)

	payload, err := json.Marshal(nc)
	if err != nil {
		panic(err)
	}

	buff := bytes.NewBuffer(payload)
	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", target, buff)
	resp, err := client.Do(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var receiver FailureResponse
	json.NewDecoder(resp.Body).Decode(&receiver)

	// make sure the error message is correct
	assert.Equal(t, "No specified directory matches path", receiver.Message)

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

func Test_apiGetRepositoryInformationHandler(t *testing.T) {
	server = createTestServerWithContext()

	gitRepoPath := "../tests/tmp/repositories/repo_info"
	oid, _ := setupSmallTestRepo(gitRepoPath)

	target := fmt.Sprintf("%s/%s", server.URL, "api/repository_info")

	client := &http.Client{}

	req, _ := http.NewRequest("GET", target, nil)

	resp, _ := client.Do(req)

	var actual RepositoryInfo

	json.NewDecoder(resp.Body).Decode(&actual)

	expected := RepositoryInfo{LatestRevision: oid.String()}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expected, actual)

}
