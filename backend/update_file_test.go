package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateFile(t *testing.T) {

	// copy the small repo tmp/repositories/update_file
	repoPath := "../tests/tmp/repositories/update_file"

	oid, _ := setupSmallTestRepo(repoPath)

	// initialise a config obj so updateFile looks in the right place
	config = Config{
		Repository: filepath.Join(repoPath),
	}

	rw := RepoWrite{
		Body:     "# The quick brown fox\n\njumped over the lazy dog",
		Filename: "document_2.md",
		Path:     "documents",
		Message:  "Update document 2",
		Name:     "Milhouse van Houten",
		Email:    "milhouse@springfield.gov",
		FrontMatter: FrontMatter{
			Title:  "Document Two",
			Author: "Hans Moleman",
		},
	}

	oid, err := updateFile(rw)

	if err != nil {
		panic(err)
	}
	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// our commit hash should now equal the repo's head
	assert.Equal(t, oid, hc.Id())

	// ensure the file exists and has the right content
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, rw.Path, rw.Filename))
	assert.Contains(t, string(contents), rw.Body)
	assert.Contains(t, string(contents), rw.FrontMatter.Author)
	assert.Contains(t, string(contents), rw.FrontMatter.Title)

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, rw.Name)
	assert.Equal(t, lastCommit.Committer().Email, rw.Email)
	assert.Equal(t, lastCommit.Message(), rw.Message)

	// finally clean up by removing the tmp repo
	_ = os.RemoveAll(repoPath)

}

func TestUpdateFiles(t *testing.T) {

	// copy the small repo tmp/repositories/update_file
	repoPath := "../tests/tmp/repositories/update_files"

	oid, _ := setupSmallTestRepo(repoPath)

	// initialise a config obj so updateFile looks in the right place
	config = Config{
		Repository: filepath.Join(repoPath),
	}

	fw1 := FileWrite{
		Filename:  "document_1.md",
		Path:      "documents",
		Extension: "md",
	}

	fw2 := FileWrite{
		Filename:  "document_2.md",
		Path:      "documents",
		Extension: "md",
	}

	rw := NewRepoWrite{
		Message: "Update document 2",
		Name:    "Milhouse van Houten",
		Email:   "milhouse@springfield.gov",
		Files:   []FileWrite{fw1, fw2},
	}

	oid, err := updateFiles(rw)

	if err != nil {
		panic(err)
	}
	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// our commit hash should now equal the repo's head
	assert.Equal(t, oid, hc.Id())

	// ensure both files exists and have the right content

	file1Contents, _ := ioutil.ReadFile(filepath.Join(repoPath, fw1.Path, fw1.Filename))
	assert.Contains(t, string(file1Contents), fw1.Body)
	assert.Contains(t, string(file1Contents), fw1.FrontMatter.Author)
	assert.Contains(t, string(file1Contents), fw1.FrontMatter.Title)

	file2Contents, _ := ioutil.ReadFile(filepath.Join(repoPath, fw2.Path, fw2.Filename))
	assert.Contains(t, string(file2Contents), fw2.Body)
	assert.Contains(t, string(file2Contents), fw2.FrontMatter.Author)
	assert.Contains(t, string(file2Contents), fw2.FrontMatter.Title)

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, rw.Name)
	assert.Equal(t, lastCommit.Committer().Email, rw.Email)
	assert.Equal(t, lastCommit.Message(), rw.Message)

	// finally clean up by removing the tmp repo
	_ = os.RemoveAll(repoPath)

}

func TestUpdateFileWhenNotExists(t *testing.T) {

	repoPath := "../tests/tmp/repositories/update_file"

	setupSmallTestRepo(repoPath)

	rw := RepoWrite{
		Body:     "# The quick brown fox\n\njumped over the lazy dog",
		Filename: "document_9.md",
		Path:     "documents",
		Message:  "Add document 9",
		Name:     "Ned Flanders",
		Email:    "nedward.flanders@leftorium.com",
	}

	repo, _ := repository(config)
	hcBefore, _ := headCommit(repo)

	_, err := updateFile(rw)

	// check error message is correct
	assert.Contains(t, err.Error(), "file does not exist")

	hcAfter, _ := headCommit(repo)

	// nothing should have been committed
	assert.Equal(t, hcBefore, hcAfter)

}
