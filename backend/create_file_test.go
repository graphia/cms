package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFiles(t *testing.T) {

	repoPath := "../tests/tmp/repositories/create_file"

	setupSmallTestRepo(repoPath)

	user := User{
		Name:  "Milhouse van Houten",
		Email: "milhouse@springfield.gov",
	}

	ncf1 := NewCommitFile{
		Filename:  "document_14.md",
		Path:      "documents",
		Extension: "md",
		Body:      "Cows don't look like cows on film. You gotta use horses.",
		FrontMatter: FrontMatter{
			Author: "Lindsay Naegle",
			Title:  "I'm an alcoholic",
		},
	}

	ncf2 := NewCommitFile{
		Filename:  "document_15.md",
		Path:      "documents",
		Extension: "md",
		Body:      "You don't win friends with salad.",
		FrontMatter: FrontMatter{
			Author: "Lindsay Naegle",
			Title:  "Children are the future, today belongs to me!",
		},
	}

	nc := NewCommit{
		Message: "Update document 2",
		Files:   []NewCommitFile{ncf1, ncf2},
	}

	oid, err := createFiles(nc, user)

	if err != nil {
		panic(err)
	}
	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// our commit hash should now equal the repo's head
	assert.Equal(t, oid, hc.Id())

	// ensure both files exists and have the right content

	file1Contents, _ := ioutil.ReadFile(filepath.Join(repoPath, ncf1.Path, ncf1.Filename))
	assert.Contains(t, string(file1Contents), ncf1.Body)
	assert.Contains(t, string(file1Contents), fmt.Sprintf("author: %s", ncf1.FrontMatter.Author))
	assert.Contains(t, string(file1Contents), fmt.Sprintf("title: %s", ncf1.FrontMatter.Title))

	file2Contents, _ := ioutil.ReadFile(filepath.Join(repoPath, ncf2.Path, ncf2.Filename))
	assert.Contains(t, string(file2Contents), ncf2.Body)
	assert.Contains(t, string(file2Contents), fmt.Sprintf("author: %s", ncf2.FrontMatter.Author))
	assert.Contains(t, string(file2Contents), fmt.Sprintf("title: %s", ncf2.FrontMatter.Title))

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, user.Name)
	assert.Equal(t, lastCommit.Committer().Email, user.Email)
	assert.Equal(t, lastCommit.Message(), nc.Message)

}

func TestCreateFileWhenExists(t *testing.T) {

	repoPath := "../tests/tmp/repositories/create_file"
	setupSmallTestRepo(repoPath)
	repo, _ := repository(config)

	user := User{
		Name:  "Ned Flanders",
		Email: "nedward.flanders@leftorium.com",
	}

	nc := NewCommit{
		Message: "Add document 2",
		Files: []NewCommitFile{
			NewCommitFile{
				Body:     "# The quick brown fox\n\njumped over the lazy dog",
				Filename: "document_2.md",
				Path:     "documents",
			},
		},
	}

	hcBefore, _ := headCommit(repo)

	_, err := createFiles(nc, user)

	// check error message is correct
	assert.Contains(t, err.Error(), "file already exists")

	hcAfter, _ := headCommit(repo)

	// nothing should have been committed
	assert.Equal(t, hcBefore, hcAfter)

}
