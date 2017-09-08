package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDirectory(t *testing.T) {
	var err error

	repoPath := "../tests/tmp/repositories/create_directory"

	setupSmallTestRepo(repoPath)

	newDir := "recipes"
	commitMessage := fmt.Sprintf("Added directories: %s", newDir)

	user := User{
		Name:  "Luigi Risotto",
		Email: "luigi@luigis-restaurant.com",
	}

	nc := NewCommit{
		Directories: []NewCommitDirectory{
			NewCommitDirectory{
				Path: newDir,
			},
		},
	}

	repo, _ := repository(config)
	oid, _ := createDirectories(nc, user)
	hc, _ := headCommit(repo)

	// our commit hash should now equal the repo's head
	assert.Equal(t, oid, hc.Id())

	// ensure the file exists and has the right content
	_, err = os.Stat(filepath.Join(repoPath, newDir, "_index.md"))
	assert.False(t, os.IsNotExist(err))

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)
	assert.Equal(t, lastCommit.Message(), commitMessage)

	// finally clean up by removing the tmp repo
	_ = os.RemoveAll(repoPath)

}

func TestCreateDirectoryWhenExists(t *testing.T) {

	repoPath := "../tests/tmp/repositories/create_file"
	setupSmallTestRepo(repoPath)

	user := User{
		Name:  "Luigi Risotto",
		Email: "luigi@luigis-restaurant.com",
	}

	nc := NewCommit{
		Directories: []NewCommitDirectory{
			NewCommitDirectory{Path: "appendices"},
		},
	}
	repo, _ := repository(config)
	hcBefore, _ := headCommit(repo)

	_, err := createDirectories(nc, user)

	// check error message is correct
	assert.Contains(t, err.Error(), "directory already exists")

	hcAfter, _ := headCommit(repo)

	// nothing should have been committed
	assert.Equal(t, hcBefore, hcAfter)

}

func TestCreateDirectoryNonSpecified(t *testing.T) {

	repoPath := "../tests/tmp/repositories/create_file"
	setupSmallTestRepo(repoPath)

	user := User{
		Name:  "Luigi Risotto",
		Email: "luigi@luigis-restaurant.com",
	}

	nc := NewCommit{
		Directories: []NewCommitDirectory{}, // empty
	}
	repo, _ := repository(config)
	hcBefore, _ := headCommit(repo)

	_, err := createDirectories(nc, user)

	// check error message is correct
	assert.Contains(t, err.Error(), "at least one new directory must be specified")

	hcAfter, _ := headCommit(repo)

	// nothing should have been committed
	assert.Equal(t, hcBefore, hcAfter)

}
