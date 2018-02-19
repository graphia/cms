package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteFiles(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_file"

	oid, _ := setupSmallTestRepo(repoPath)

	user := User{
		Name:  "Milhouse van Houten",
		Email: "milhouse@springfield.gov",
	}

	ncf1 := NewCommitFile{
		Filename: "index.md",
		Document: "document_1",
		Path:     "documents",
	}

	ncf2 := NewCommitFile{
		Filename: "index.md",
		Document: "document_2",
		Path:     "documents",
	}

	nc := NewCommit{
		Message:        "Delete documents 1 and 2",
		Files:          []NewCommitFile{ncf1, ncf2},
		RepositoryInfo: RepositoryInfo{LatestRevision: oid.String()},
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

	oid, _ := setupSmallTestRepo(repoPath)

	user := User{
		Name:  "Milhouse van Houten",
		Email: "milhouse@springfield.gov",
	}

	ncf := NewCommitFile{
		Filename: "index.md",
		Document: "document_5",
		Path:     "documents",
	}

	nc := NewCommit{
		Message:        "Delete document 5",
		Files:          []NewCommitFile{ncf},
		RepositoryInfo: RepositoryInfo{LatestRevision: oid.String()},
	}

	_, err := deleteFiles(nc, user)

	assert.Contains(t, err.Error(), "file does not exist")

}

func TestDeleteFilesNoMessage(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_file_no_message"

	oid, _ := setupSmallTestRepo(repoPath)

	user := User{
		Name:  "Milhouse van Houten",
		Email: "milhouse@springfield.gov",
	}

	ncf := NewCommitFile{
		Filename: "index.md",
		Document: "document_1",
		Path:     "documents",
	}

	nc := NewCommit{
		//Message: "Delete something..",
		Files:          []NewCommitFile{ncf},
		RepositoryInfo: RepositoryInfo{LatestRevision: oid.String()},
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
	_, err = os.Stat(filepath.Join(repoPath, ncf.Path, ncf.Filename))
	assert.True(t, os.IsNotExist(err))

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)

	// and finally, check that a suitable message has been added
	assert.Equal(t, "File deleted", lastCommit.Message())

}

func TestDeleteFilesRepoOutOfDate(t *testing.T) {

	repoPath := "../tests/tmp/repositories/delete_file_no_message"

	oid, _ := setupSmallTestRepo(repoPath)

	user := User{
		Name:  "Milhouse van Houten",
		Email: "milhouse@springfield.gov",
	}

	ncf := NewCommitFile{
		Filename: "index.md",
		Document: "document_1",
		Path:     "documents",
	}

	nc := NewCommit{
		//Message: "Delete something..",
		Files:          []NewCommitFile{ncf},
		RepositoryInfo: RepositoryInfo{LatestRevision: oid.String()},
	}

	repo, _ := repository(config)

	// importantly, add another file to outdate the repo
	_, _ = createRandomFile(repo, "document_5", "en", "whoosh")

	_, err := deleteFiles(nc, user)

	assert.Equal(t, err, ErrRepoOutOfSync)

	// ensure the file hasn't been deleted
	_, err = os.Stat(filepath.Join(repoPath, ncf.Path, ncf.Document, ncf.Filename))
	assert.False(t, os.IsNotExist(err))

}

func TestDeleteFileAndAttachmentsDirectory(t *testing.T) {
	var err error

	repoPath := "../tests/tmp/repositories/delete_file_and_dir"

	oid, _ := setupMultipleFiletypesTestRepo(repoPath)

	user := User{
		Name:  "Milhouse van Houten",
		Email: "milhouse@springfield.gov",
	}

	ncf1 := NewCommitFile{
		Filename: "index.md",
		Document: "document_1",
		Path:     "documents",
	}

	ncd1 := NewCommitDirectory{
		Path: "documents/document_1",
	}

	nc := NewCommit{
		Message:        "Delete documents 1 and 2",
		Files:          []NewCommitFile{ncf1},
		Directories:    []NewCommitDirectory{ncd1},
		RepositoryInfo: RepositoryInfo{LatestRevision: oid.String()},
	}

	repo, _ := repository(config)

	// ensure the file is present on the filesystem
	_, err = os.Stat(filepath.Join(repoPath, ncf1.Path, ncf1.Document, ncf1.Filename))
	assert.False(t, os.IsNotExist(err))

	// ensure the directory is present on the filesystem
	_, err = os.Stat(filepath.Join(repoPath, ncd1.Path))
	assert.False(t, os.IsNotExist(err))

	// actually delete the files
	oid, err = deleteFiles(nc, user)
	if err != nil {
		panic(err)
	}

	// our commit hash should now equal the repo's head
	hc, _ := headCommit(repo)
	assert.Equal(t, oid, hc.Id())

	// ensure the file isn't present on the filesystem
	_, err = os.Stat(filepath.Join(repoPath, ncf1.Path, ncf1.Filename))
	assert.True(t, os.IsNotExist(err))

	// ensure the directory isn't present on the filesystem
	_, err = os.Stat(filepath.Join(repoPath, ncd1.Path))
	assert.True(t, os.IsNotExist(err))

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)
	assert.Equal(t, lastCommit.Message(), nc.Message)

}
