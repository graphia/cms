package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEmptyFile(t *testing.T) {

	repoPath := "../tests/tmp/repositories/create_file"

	setupSmallTestRepo(repoPath)

	rw := RepoWrite{
		Body:     "# Ensure this is discarded",
		Filename: "document_9.md",
		Path:     "documents",
		Message:  "Add document 9",
		Name:     "Milhouse van Houten",
		Email:    "millhouse@springfield.gov",
		FrontMatter: FrontMatter{
			Title:  "Ensure this is discarded",
			Author: "Ensure this is discarded",
		},
	}

	repo, _ := repository(config)
	oid, _ := createFile(rw)
	hc, _ := headCommit(repo)

	// our commit hash should now equal the repo's head
	assert.Equal(t, oid, hc.Id())

	// ensure the file exists and **is blank**
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, rw.Path, rw.Filename))
	assert.Contains(t, string(contents), "")

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, rw.Name)
	assert.Equal(t, lastCommit.Committer().Email, rw.Email)
	assert.Equal(t, lastCommit.Message(), rw.Message)

	// finally clean up by removing the tmp repo
	_ = os.RemoveAll(repoPath)

}

func TestCreateEmptyFileWhenExists(t *testing.T) {

	repoPath := "../tests/tmp/repositories/create_file"

	setupSmallTestRepo(repoPath)

	rw := RepoWrite{
		Body:     "# Ensure this is discarded",
		Filename: "document_2.md",
		Path:     "documents",
		Message:  "Add document 2",
		Name:     "Ned Flanders",
		Email:    "nedward.flanders@leftorium.com",
		FrontMatter: FrontMatter{
			Title:  "Ensure this is discarded",
			Author: "Ensure this is discarded",
		},
	}

	repo, _ := repository(config)
	hcBefore, _ := headCommit(repo)

	_, err := createFile(rw)

	// check error message is correct
	assert.Contains(t, err.Error(), "file already exists")

	hcAfter, _ := headCommit(repo)

	// nothing should have been committed
	assert.Equal(t, hcBefore, hcAfter)

}
