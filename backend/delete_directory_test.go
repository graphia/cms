package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteDir(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_dir"

	setupSmallTestRepo(repoPath)

	rw := RepoWrite{
		Path:  "appendices",
		Name:  "Moe Szyslak",
		Email: "moe@moes.com",
	}

	repo, _ := repository(config)

	oid, err := deleteDirectory(rw)
	if err != nil {
		panic(err)
	}

	hc, _ := headCommit(repo)

	// our commit hash should now equal the repo's head
	assert.Equal(t, oid, hc.Id())

	// ensure the file isn't present on the filesystem
	_, err = os.Stat(filepath.Join(repoPath, rw.Path))
	assert.True(t, os.IsNotExist(err))

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, rw.Name)
	assert.Equal(t, lastCommit.Committer().Email, rw.Email)
	assert.Equal(t, lastCommit.Message(), rw.Message)

}

func TestDeleteDirNotExists(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_file"

	setupSmallTestRepo(repoPath)

	rw := RepoWrite{
		Path:  "menu",
		Name:  "Barney Gumble",
		Email: "barney@plow-king.com",
	}

	_, err := deleteDirectory(rw)

	assert.Contains(t, err.Error(), "directory does not exist")

}
