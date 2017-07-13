package main

import (
	"fmt"
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

	fw := FileWrite{
		Body:     "# The quick brown fox\n\njumped over the lazy dog",
		Filename: "document_2.md",
		Path:     "documents",
		FrontMatter: FrontMatter{
			Title:  "Document Two",
			Author: "Hans Moleman",
		},
	}

	rw := NewRepoWrite{

		Message: "Update document 2",
		Name:    "Milhouse van Houten",
		Email:   "milhouse@springfield.gov",
		Files:   []FileWrite{fw},
	}

	oid, err := updateFiles(rw)

	if err != nil {
		panic(err)
	}
	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	// our commit hash should now equal the repo's head
	assert.Equal(t, oid, hc.Id())

	// ensure the file exists and has the right content
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, fw.Path, fw.Filename))
	assert.Contains(t, string(contents), fw.Body)
	assert.Contains(t, string(contents), fw.FrontMatter.Author)
	assert.Contains(t, string(contents), fw.FrontMatter.Title)

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
		Body:      "Cows don't look like cows on film. You gotta use horses.",
		FrontMatter: FrontMatter{
			Author: "Lindsay Naegle",
			Title:  "I'm an alcoholic",
		},
	}

	fw2 := FileWrite{
		Filename:  "document_2.md",
		Path:      "documents",
		Extension: "md",
		Body:      "You don't win friends with salad.",
		FrontMatter: FrontMatter{
			Author: "Lindsay Naegle",
			Title:  "Children are the future, today belongs to me!",
		},
	}

	rw := NewRepoWrite{
		Message: "Update document 1 and 2",
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
	assert.Contains(t, string(file1Contents), fmt.Sprintf("author: %s", fw1.FrontMatter.Author))
	assert.Contains(t, string(file1Contents), fmt.Sprintf("title: %s", fw1.FrontMatter.Title))

	file2Contents, _ := ioutil.ReadFile(filepath.Join(repoPath, fw2.Path, fw2.Filename))
	assert.Contains(t, string(file2Contents), fw2.Body)
	assert.Contains(t, string(file2Contents), fmt.Sprintf("author: %s", fw2.FrontMatter.Author))
	assert.Contains(t, string(file2Contents), fmt.Sprintf("title: %s", fw2.FrontMatter.Title))

	// ensure the most recent commit has the right name and email
	lastCommit, _ := repo.LookupCommit(oid)
	assert.Equal(t, lastCommit.Committer().Name, rw.Name)
	assert.Equal(t, lastCommit.Committer().Email, rw.Email)
	assert.Equal(t, lastCommit.Message(), rw.Message)

}

func TestUpdateFileWhenNotExists(t *testing.T) {

	repoPath := "../tests/tmp/repositories/update_file"

	setupSmallTestRepo(repoPath)

	rw := NewRepoWrite{
		Message: "Add document 9",
		Name:    "Ned Flanders",
		Email:   "nedward.flanders@leftorium.com",
		Files: []FileWrite{
			FileWrite{
				Body:     "# The quick brown fox\n\njumped over the lazy dog",
				Filename: "document_9.md",
				Path:     "documents",
			},
		},
	}

	repo, _ := repository(config)
	hcBefore, _ := headCommit(repo)

	_, err := updateFiles(rw)

	// check error message is correct
	assert.Contains(t, err.Error(), "file not found: documents/document_9.md")

	hcAfter, _ := headCommit(repo)

	// nothing should have been committed
	assert.Equal(t, hcBefore, hcAfter)

}
