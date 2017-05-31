package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFilesInDocumentsDir(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	// using the small test repo, there are three documents
	filesExpected := 3

	files, err := getFilesInDir("documents")
	if err != nil {
		t.Error("error", err)
	}
	filesCount := len(files)
	assert.Equal(t, filesExpected, filesCount)
}

func TestGetFilesInAppendicesDir(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	// using the small test repo, there are two appendices
	filesExpected := 2

	// ensure a third non-Markdown file exists in the repo but won't be returned
	_, err := os.Stat(filepath.Join(repoPath, "appendices", "another_file.txt"))
	assert.False(t, os.IsNotExist(err))

	files, err := getFilesInDir("appendices")
	if err != nil {
		t.Error("error", err)
	}
	filesCount := len(files)
	assert.Equal(t, filesExpected, filesCount)

}

func TestGetFilesInDirContents(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	files, _ := getFilesInDir("appendices")

	file := files[0]

	assert.Equal(t, "appendix_1.md", file.Filename)
	assert.Equal(t, "Appendix 1", file.Title)
	assert.Equal(t, "appendices", file.Path)
	assert.Equal(t, "appendices/appendix_1.md", file.AbsoluteFilename)
	assert.Equal(t, "1.1", file.Version)
	assert.Equal(t, "Arnold Pye", file.Author)
	assert.Equal(t, []string{"Traffic News", "KBBL TV"}, file.Tags)
}

func TestGetFilesInNonExistantDir(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	// there isn't a directory called fanfic, so should
	// raise not found error
	_, err := getFilesInDir("fanfic")
	assert.Contains(t, err.Error(), "directory 'fanfic' not found")
}

func TestGetConvertedFile(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getConvertedFile("documents", "document_2.md")
	if err != nil {
		t.Error("error", err)
	}

	assert.Equal(t, file.Filename, "document_2.md")
	assert.Equal(t, file.Path, "documents")

	// just look for the title rather than rerender md in  here
	assert.Contains(t, *file.HTML, "<h1>Document 2</h1>")

	// markdown should be nil
	assert.Nil(t, file.Markdown)
}

func TestGetRawFile(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getRawFile("documents", "document_2.md")
	if err != nil {
		t.Error("error", err)
	}

	contents, _ := ioutil.ReadFile(filepath.Join(
		config.Repository,
		"documents",
		"document_2.md",
	))

	assert.Equal(t, file.Filename, "document_2.md")
	assert.Equal(t, file.Path, "documents")
	assert.Equal(t, *file.Markdown, string(contents))

	// this time, HTML should be nil
	assert.Nil(t, file.HTML)
}

func TestGetFileBothMarkdownAndHTML(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getFile("documents", "document_2.md", true, true)
	if err != nil {
		t.Error("error", err)
	}

	contents, _ := ioutil.ReadFile(filepath.Join(
		config.Repository,
		"documents",
		"document_2.md",
	))

	assert.Equal(t, file.Filename, "document_2.md")
	assert.Equal(t, file.Path, "documents")
	assert.Equal(t, *file.Markdown, string(contents))
	assert.Contains(t, *file.HTML, "<h1>Document 2</h1")

}

func TestGetFileNeitherMarkdownOrHTML(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getFile("documents", "document_2.md", false, false)
	if err != nil {
		t.Error("error", err)
	}

	assert.Equal(t, file.Filename, "document_2.md")
	assert.Equal(t, file.Path, "documents")
	assert.Nil(t, file.HTML)
	assert.Nil(t, file.Markdown)

}
