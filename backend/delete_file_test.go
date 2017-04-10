package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteFile(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_file"

	setupSmallTestRepo(repoPath)

	rw := RepoWrite{
		Filename: "document_2.md",
		Path:     "documents",
		Message:  "Delete document 2",
		Name:     "Milhouse van Houten",
		Email:    "millhouse@springfield.gov",
	}

	repo, _ := repository(config)

	oid, err := deleteFile(rw)
	if err != nil {
		panic(err)
	}

	hc, _ := headCommit(repo)

	// our commit hash should now equal the repo's head
	assert.Equal(t, oid, hc.Id())

	// ensure the file isn't present on the filesystem
	_, err = os.Stat(filepath.Join(repoPath, rw.Path, rw.Filename))
	assert.True(t, os.IsNotExist(err))

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, rw.Name)
	assert.Equal(t, lastCommit.Committer().Email, rw.Email)
	assert.Equal(t, lastCommit.Message(), rw.Message)

}

func TestDeleteFileNotExists(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_file"

	setupSmallTestRepo(repoPath)

	rw := RepoWrite{
		Filename: "document_8.md",
		Path:     "documents",
		Message:  "Delete document 8",
		Name:     "Nelson Muntz",
		Email:    "muntz@ha-ha.org",
	}

	_, err := deleteFile(rw)

	assert.Contains(t, err.Error(), "file does not exist")

}
