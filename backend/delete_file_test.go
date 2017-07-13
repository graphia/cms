package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteFiles(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_file"

	setupSmallTestRepo(repoPath)

	fw1 := FileWrite{
		Filename: "document_1.md",
		Path:     "documents",
	}

	fw2 := FileWrite{
		Filename: "document_2.md",
		Path:     "documents",
	}

	rw := NewRepoWrite{
		Message: "Delete documents 1 and 2",
		Name:    "Milhouse van Houten",
		Email:   "milhouse@springfield.gov",
		Files:   []FileWrite{fw1, fw2},
	}

	repo, _ := repository(config)

	oid, err := deleteFiles(rw)
	if err != nil {
		panic(err)
	}

	hc, _ := headCommit(repo)

	// our commit hash should now equal the repo's head
	assert.Equal(t, oid, hc.Id())

	// ensure the file isn't present on the filesystem
	_, err = os.Stat(filepath.Join(repoPath, fw1.Path, fw1.Filename))
	_, err = os.Stat(filepath.Join(repoPath, fw2.Path, fw2.Filename))

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

	fw := FileWrite{
		Filename: "document_5.md",
		Path:     "documents",
	}

	rw := NewRepoWrite{
		Message: "Delete document 5",
		Name:    "Milhouse van Houten",
		Email:   "milhouse@springfield.gov",
		Files:   []FileWrite{fw},
	}

	_, err := deleteFiles(rw)

	assert.Contains(t, err.Error(), "file does not exist")

}
