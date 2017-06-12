package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/libgit2/git2go.v25"
)

// repository

func TestCorrectRepo(t *testing.T) {

	testConfig, err := loadConfig(&smallRepoPath)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, "../tests/backend/repositories/small", testConfig.Repository)
}

func TestInvalidRepo(t *testing.T) {
	wd, _ := os.Getwd()

	invalidConfig, _ := loadConfig(&invalidRepoPath)
	_, err := repository(invalidConfig)

	msg := fmt.Sprintf(
		"Failed to resolve path '%s': No such file or directory",
		filepath.Join(wd, invalidConfig.Repository),
	)

	assert.Contains(t, err.Error(), msg)
}

// headCommit

func TestNoHeadCommit(t *testing.T) {
	var repo *git.Repository

	wd, _ := os.Getwd()
	path := filepath.Join(wd, "../tests/tmp/repositories/empty")

	// initialise the empty repo or open it if it's there
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
		repo, _ = git.InitRepository(path, false)
	} else {
		repo, _ = git.OpenRepository(path)
	}

	_, err := headCommit(repo)

	msg := "Cannot find repository head"

	assert.Contains(t, err.Error(), msg)

}

func TestHeadCommit(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_file"
	oid, _ := setupSmallTestRepo(repoPath)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	assert.Equal(t, hc.Id(), oid)
}

func createRandomFile(repo *git.Repository, filename, msg string) error {
	rw := RepoWrite{
		Body:     "# The quick brown fox\n\njumped over the lazy dog",
		Filename: filename,
		Path:     "documents",
		Message:  msg,
		Name:     "Barney Gumble",
		Email:    "barney.gumble@hotmail.com",
		FrontMatter: FrontMatter{
			Title:  "Document Twelve",
			Author: "Bernard Arnold Gumble",
		},
	}
	_, err := createFile(rw)
	return err
}

func TestAllCommits(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_commits"
	_, _ = setupSmallTestRepo(repoPath)

	repo, _ := repository(config)

	msg := "Well well, if it isn't Mr Plow"

	err := createRandomFile(repo, "document_12.md", msg)
	if err != nil {
		panic(err)
	}

	commits, _ := allCommits(repo, 10)

	var messages []string

	for _, commit := range commits {
		messages = append(messages, commit.Message)
	}

	// should return both commit messages
	assert.Equal(t, []string{msg, "Quick, honk at that broad!"}, messages)

}

func TestAllCommitsWithLimitOf3(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_commits"
	_, _ = setupSmallTestRepo(repoPath)

	repo, _ := repository(config)

	for i := 12; i <= 16; i++ {
		err := createRandomFile(
			repo,
			fmt.Sprintf("document_%d.md", i),
			fmt.Sprintf("Commit Message %d", i),
		)

		if err != nil {
			panic(err)
		}
	}

	commits, _ := allCommits(repo, 3)

	var messages []string

	for _, commit := range commits {
		messages = append(messages, commit.Message)
	}

	// should return both commit messages
	assert.Equal(t, []string{"Commit Message 16", "Commit Message 15", "Commit Message 14"}, messages)

}

func TestDiffForCommit(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_commits"
	oid, _ := setupSmallTestRepo(repoPath)
	count, diff, _ := diffForCommit(oid.String())

	// Make sure only the correct five markdown files were added
	assert.Equal(t, 5, count)
	assert.Contains(t, diff, "+++ b/documents/document_1.md")
	assert.Contains(t, diff, "+++ b/documents/document_2.md")
	assert.Contains(t, diff, "+++ b/documents/document_3.md")
	assert.Contains(t, diff, "+++ b/appendices/appendix_1.md")
	assert.Contains(t, diff, "+++ b/appendices/appendix_2.md")

	// Ensure the file contents are included
	assert.Contains(t, diff, "+Lorem ipsum dolor sit")

}
