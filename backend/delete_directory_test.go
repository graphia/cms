package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteDirectories(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_dir"

	setupSmallTestRepo(repoPath)

	user := User{
		Name:  "Moe Szyslak",
		Email: "moe@moes.com",
	}

	ncd := NewCommitDirectory{Path: "appendices"}

	nc := NewCommit{
		Message:     "Deleted directories",
		Directories: []NewCommitDirectory{ncd},
	}

	repo, _ := repository(config)

	oid, err := deleteDirectories(nc, user)
	if err != nil {
		panic(err)
	}

	hc, _ := headCommit(repo)

	// our commit hash should now equal the repo's head
	assert.Equal(t, oid, hc.Id())

	// ensure the file isn't present on the filesystem
	_, err = os.Stat(filepath.Join(repoPath, ncd.Path))
	assert.True(t, os.IsNotExist(err))

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)
	assert.Equal(t, lastCommit.Message(), nc.Message)

}

func TestDeleteDirectoriesNotExists(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_file"

	setupSmallTestRepo(repoPath)

	user := User{
		Name:  "Barney Gumble",
		Email: "barney@plow-king.com",
	}

	ncd := NewCommitDirectory{Path: "menu"}

	nc := NewCommit{

		Directories: []NewCommitDirectory{ncd},
	}

	_, err := deleteDirectories(nc, user)

	assert.Equal(t, err.Error(), "directory does not exist: menu")

}
