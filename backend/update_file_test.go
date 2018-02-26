package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateFiles(t *testing.T) {

	// copy the small repo tmp/repositories/update_file
	repoPath := "../tests/tmp/repositories/update_files"

	oid, _ := setupSmallTestRepo(repoPath)

	// initialise a config obj so updateFile looks in the right place
	config = Config{
		Repository: filepath.Join(repoPath),
	}

	user := User{
		Name:  "Milhouse van Houten",
		Email: "milhouse@springfield.gov",
	}

	ncf1 := NewCommitFile{
		Filename: "document_1.md",
		Path:     "documents",
		Body:     "Cows don't look like cows on film. You gotta use horses.",
		FrontMatter: FrontMatter{
			Author: "Lindsay Naegle",
			Title:  "I'm an alcoholic",
		},
	}

	ncf2 := NewCommitFile{
		Filename: "document_2.md",
		Path:     "documents",
		Body:     "You don't win friends with salad.",
		FrontMatter: FrontMatter{
			Author: "Lindsay Naegle",
			Title:  "Children are the future, today belongs to me!",
		},
	}

	nc := NewCommit{
		Message:        "Update document 1 and 2",
		Files:          []NewCommitFile{ncf1, ncf2},
		RepositoryInfo: RepositoryInfo{LatestRevision: oid.String()},
	}

	oid, err := updateFiles(nc, user)

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

func TestUpdateFilesRepoOutOfDate(t *testing.T) {

	// copy the small repo tmp/repositories/update_file
	repoPath := "../tests/tmp/repositories/update_files"

	oid, _ := setupSmallTestRepo(repoPath)

	// initialise a config obj so updateFile looks in the right place
	config = Config{
		Repository: filepath.Join(repoPath),
	}

	user := User{
		Name:  "Milhouse van Houten",
		Email: "milhouse@springfield.gov",
	}

	nfc := NewCommitFile{
		Filename: "index.md",
		Document: "document_1",
		Path:     "documents",
		Body:     "Cows don't look like cows on film. You gotta use horses.",
		FrontMatter: FrontMatter{
			Author: "Lindsay Naegle",
			Title:  "I'm an alcoholic",
		},
	}

	nc := NewCommit{
		Message:        "Update document 1",
		Files:          []NewCommitFile{nfc},
		RepositoryInfo: RepositoryInfo{LatestRevision: oid.String()},
	}

	repo, _ := repository(config)
	_, _ = createRandomFile(repo, "document_5.md", "en", "whoosh")

	oid, err := updateFiles(nc, user)

	assert.Equal(t, err, ErrRepoOutOfSync)

	// make sure the file hasn't been updated
	contents, _ := ioutil.ReadFile(filepath.Join(repoPath, nfc.Path, nfc.Document, nfc.Filename))
	assert.Contains(t, string(contents), "Lorem ipsum dolor sit amet")

}
