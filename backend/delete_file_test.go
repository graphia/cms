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

	user := User{
		Name:  "Milhouse van Houten",
		Email: "milhouse@springfield.gov",
	}

	ncf1 := NewCommitFile{
		Filename: "document_1.md",
		Path:     "documents",
	}

	ncf2 := NewCommitFile{
		Filename: "document_2.md",
		Path:     "documents",
	}

	nc := NewCommit{
		Message: "Delete documents 1 and 2",
		Files:   []NewCommitFile{ncf1, ncf2},
	}

	repo, _ := repository(config)

	oid, err := deleteFiles(nc, user)
	if err != nil {
		panic(err)
	}

	hc, _ := headCommit(repo)

	// our commit hash should now equal the repo's head
	assert.Equal(t, oid, hc.Id())

	// ensure the file isn't present on the filesystem
	_, err = os.Stat(filepath.Join(repoPath, ncf1.Path, ncf1.Filename))
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(filepath.Join(repoPath, ncf2.Path, ncf2.Filename))
	assert.True(t, os.IsNotExist(err))

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)
	assert.Equal(t, lastCommit.Message(), nc.Message)

}

func TestDeleteFileNotExists(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_file"

	setupSmallTestRepo(repoPath)

	user := User{
		Name:  "Milhouse van Houten",
		Email: "milhouse@springfield.gov",
	}

	ncf := NewCommitFile{
		Filename: "document_5.md",
		Path:     "documents",
	}

	nc := NewCommit{
		Message: "Delete document 5",
		Files:   []NewCommitFile{ncf},
	}

	_, err := deleteFiles(nc, user)

	assert.Contains(t, err.Error(), "file does not exist")

}
