package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"time"

	"github.com/graphia/particle"
	"github.com/stretchr/testify/assert"
	"gopkg.in/libgit2/git2go.v25"
)

// repository

func TestCorrectRepo(t *testing.T) {

	testConfig, err := loadConfig(&smallRepoPath)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, "../tests/backend/repositories/small", testConfig.Repository)
}

func TestInvalidRepo(t *testing.T) {
	wd, _ := os.Getwd()

	invalidConfig, _ := loadConfig(&invalidRepoPath)
	_, err := repository(invalidConfig)

	msg := fmt.Sprintf(
		"Failed to resolve path '%s': No such file or directory",
		filepath.Join(wd, invalidConfig.Repository),
	)

	assert.Contains(t, err.Error(), msg)
}

// headCommit

func TestNoHeadCommit(t *testing.T) {
	var repo *git.Repository

	wd, _ := os.Getwd()
	path := filepath.Join(wd, "../tests/tmp/repositories/empty")

	// initialise the empty repo or open it if it's there
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
		repo, _ = git.InitRepository(path, false)
	} else {
		repo, _ = git.OpenRepository(path)
	}

	_, err := headCommit(repo)

	msg := "Cannot find repository head"

	assert.Contains(t, err.Error(), msg)

}

func TestHeadCommit(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_file"
	oid, _ := setupSmallTestRepo(repoPath)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)

	assert.Equal(t, hc.Id(), oid)
}

func createRandomFile(repo *git.Repository, filename, msg string) error {
	nc := NewCommit{
		Message: msg,
		Name:    "Barney Gumble",
		Email:   "barney.gumble@hotmail.com",
		Files: []NewCommitFile{
			NewCommitFile{
				Body:     "# The quick brown fox\n\njumped over the lazy dog",
				Filename: filename,
				Path:     "documents",
				FrontMatter: FrontMatter{
					Title:  "Document Twelve",
					Author: "Bernard Arnold Gumble",
				},
			},
		},
	}
	_, err := createFiles(nc)
	return err
}

func TestAllCommits(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_commits"
	_, _ = setupSmallTestRepo(repoPath)

	repo, _ := repository(config)

	msg := "Well well, if it isn't Mr Plow"

	err := createRandomFile(repo, "document_12.md", msg)
	if err != nil {
		panic(err)
	}

	commits, _ := allCommits(repo, 10)

	var messages []string

	for _, commit := range commits {
		messages = append(messages, commit.Message)
	}

	// should return both commit messages
	assert.Equal(t, []string{msg, "Quick, honk at that broad!"}, messages)

}

func TestAllCommitsWithLimitOf3(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_commits"
	_, _ = setupSmallTestRepo(repoPath)

	repo, _ := repository(config)

	for i := 12; i <= 16; i++ {
		err := createRandomFile(
			repo,
			fmt.Sprintf("document_%d.md", i),
			fmt.Sprintf("Commit Message %d", i),
		)

		if err != nil {
			panic(err)
		}
	}

	commits, _ := allCommits(repo, 3)

	var messages []string

	for _, commit := range commits {
		messages = append(messages, commit.Message)
	}

	// should return both commit messages
	assert.Equal(t, []string{"Commit Message 16", "Commit Message 15", "Commit Message 14"}, messages)

}

func TestDiffForCommit(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_commits"
	oid, _ := setupSmallTestRepo(repoPath)
	cs, _ := diffForCommit(oid.String())

	// Make sure only the correct five markdown files were added
	assert.Equal(t, 5, cs.NumDeltas)
	assert.Equal(t, 0, cs.NumDeleted)
	assert.Equal(t, 5, cs.NumAdded)
	assert.Contains(t, cs.FullDiff, "+++ b/documents/document_1.md")
	assert.Contains(t, cs.FullDiff, "+++ b/documents/document_2.md")
	assert.Contains(t, cs.FullDiff, "+++ b/documents/document_3.md")
	assert.Contains(t, cs.FullDiff, "+++ b/appendices/appendix_1.md")
	assert.Contains(t, cs.FullDiff, "+++ b/appendices/appendix_2.md")

	// Ensure the file contents are included
	assert.Contains(t, cs.FullDiff, "+Lorem ipsum dolor sit")

	var allFilesInRepo = []string{
		"documents/document_1.md",
		"documents/document_2.md",
		"documents/document_3.md",
		"appendices/appendix_1.md",
		"appendices/appendix_2.md",
	}

	for _, path := range allFilesInRepo {
		// Ensure 'old' files are empty
		assert.Empty(t, cs.Files[path].Old)

		// Ensure 'new' files contain the correct info
		contents, _ := ioutil.ReadFile(filepath.Join(repoPath, path))
		assert.Equal(t, cs.Files[path].New, string(contents))
	}

}

func TestLookupFileHistory(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_history_test"
	oid, _ := setupSmallTestRepo(repoPath)
	repo, _ := repository(config)

	assert.NotNil(t, oid)

	dir := "documents"
	filename := "history_test.md"
	path := fmt.Sprintf("%s/%s", dir, filename)

	template := NewCommit{

		Name:  "Hyman Krustofski",
		Email: "hyman@springfieldsynagogue.org",
		Files: []NewCommitFile{
			NewCommitFile{
				Filename: filename,
				Path:     dir,
				FrontMatter: FrontMatter{
					Title:  "History Tests",
					Author: "Helen Lovejoy",
				},
			},
		},
	}

	var r1, r2, r3 NewCommit

	// Revision 1
	r1 = template
	r1.Files[0].Body = "# r1"
	r1.Message = "r1"
	oid, _ = createFiles(r1)
	assert.NotNil(t, oid)

	// Revision 2
	r2 = template
	r2.Files[0].Body = "# r2"
	r2.Message = "r2"
	oid, _ = updateFiles(r2)
	assert.NotNil(t, oid)

	// Revision 3
	r3 = template
	r3.Files[0].Body = "# r3"
	r3.Message = "r3"
	oid, _ = updateFiles(r3)
	assert.NotNil(t, oid)

	var history []HistoricCommit

	history, _ = lookupFileHistory(repo, path, 3)

	assert.Equal(t, 3, len(history))

	var messages []string
	for _, commit := range history {
		messages = append(messages, commit.Message)
	}

	sort.Strings(messages)

	assert.Equal(
		t,
		[]string{
			"r1",
			"r2",
			"r3",
		},
		messages,
	)

	// Check that retrieving a subset of the history also works
	history, _ = lookupFileHistory(repo, path, 2)
	assert.Equal(t, 2, len(history))
}

func TestLookupFileHistorySortsByTime(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_history_test"
	oid, _ := setupSmallTestRepo(repoPath)
	repo, _ := repository(config)

	template := NewCommit{

		Name:  "Hyman Krustofski",
		Email: "hyman@springfieldsynagogue.org",
		Files: []NewCommitFile{
			NewCommitFile{
				Filename: "sort_test.md",
				Path:     "documents",
				FrontMatter: FrontMatter{
					Title:  "History Tests",
					Author: "Helen Lovejoy",
				},
			},
		},
	}

	// Revisions ordered numerically but entered in the wrong order
	var r1, r2, r3 NewCommit

	// Revision 2
	r2 = template
	r2.Files[0].Body = "# r2"
	r2.Message = "r2"
	oid, _ = writeHistoricFiles(repo, r2, time.Date(2016, 1, 1, 14, 0, 0, 0, time.UTC))
	assert.NotNil(t, oid)

	// Revision 3
	r3 = template
	r3.Files[0].Body = "# r3"
	r3.Message = "r3"
	oid, _ = writeHistoricFiles(repo, r3, time.Date(2016, 1, 1, 15, 0, 0, 0, time.UTC))
	assert.NotNil(t, oid)

	// Revision 1
	r1 = template
	r1.Files[0].Body = "# r1"
	r1.Message = "r1"
	oid, _ = writeHistoricFiles(repo, r1, time.Date(2016, 1, 1, 13, 0, 0, 0, time.UTC))
	assert.NotNil(t, oid)

	var history []HistoricCommit

	history, _ = lookupFileHistory(repo, "documents/sort_test.md", 3)

	assert.Equal(t, 3, len(history))

	var messages []string
	for _, commit := range history {
		messages = append(messages, commit.Message)
	}

	assert.Equal(
		t,
		[]string{
			"r3",
			"r2",
			"r1",
		},
		messages,
	)

	// Check that retrieving a subset of the history also works

	var subsetMessages []string

	history, _ = lookupFileHistory(repo, "documents/sort_test.md", 2)
	assert.Equal(t, 2, len(history))

	for _, commit := range history {
		subsetMessages = append(subsetMessages, commit.Message)
	}

	assert.Equal(
		t,
		[]string{
			"r3",
			"r2",
		},
		subsetMessages,
	)

}

func TestLookupFileHistoryOnlyReturnsRelevantCommits(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_history_test"
	oid, _ := setupSmallTestRepo(repoPath)
	repo, _ := repository(config)

	template := NewCommit{

		Name:  "Hyman Krustofski",
		Email: "hyman@springfieldsynagogue.org",
		Files: []NewCommitFile{
			NewCommitFile{
				Path: "documents",
				FrontMatter: FrontMatter{
					Title:  "History Tests",
					Author: "Helen Lovejoy",
				},
			},
		},
	}

	// Revisions ordered numerically but entered in the wrong order
	var r1, r2, r3 NewCommit

	// Revision 1
	r1 = template
	r1.Files[0].Filename = "document_11.md"
	r1.Files[0].Body = "# r1"
	r1.Message = "r1"
	oid, _ = createFiles(r1)
	assert.NotNil(t, oid)

	// Revision 2 (unrelated)
	r2 = template
	r2.Files[0].Filename = "document_12.md"
	r2.Files[0].Body = "# r2"
	r2.Message = "r2"
	oid, _ = createFiles(r2)
	assert.NotNil(t, oid)

	// Revision 1
	r3 = template
	r3.Files[0].Filename = "document_11.md"
	r3.Files[0].Body = "# r3"
	r3.Message = "r3"
	oid, _ = updateFiles(r3)
	assert.NotNil(t, oid)

	var history []HistoricCommit

	history, _ = lookupFileHistory(repo, "documents/document_11.md", 2)

	assert.Equal(t, 2, len(history))

	var messages []string
	for _, commit := range history {
		messages = append(messages, commit.Message)
	}

	// r2 shouldn't be present because it refers to 'documents/document_12.md'
	assert.Equal(
		t,
		[]string{
			"r3",
			"r1",
		},
		messages,
	)

}

// Utility functions

func writeHistoricFiles(repo *git.Repository, nc NewCommit, time time.Time) (oid *git.Oid, err error) {

	index, err := repo.Index()
	if err != nil {
		return nil, err
	}
	defer index.Free()

	var contents string

	for _, ncf := range nc.Files {

		var ie git.IndexEntry

		contents = particle.YAMLEncoding.EncodeToString([]byte(ncf.Body), &ncf.FrontMatter)

		oid, err = repo.CreateBlobFromBuffer([]byte(contents))
		if err != nil {
			return nil, err
		}

		// build the git index entry and add it to the index
		ie = buildIndexEntry(oid, ncf)

		err = index.Add(&ie)
		if err != nil {
			return nil, err
		}

	}

	// write the tree, persisting our addition to the git repo
	treeID, err := index.WriteTree()
	if err != nil {
		return nil, err
	}

	// and use the tree's id to find the actual updated tree
	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return nil, err
	}

	// find the repository's tip, where we're committing to
	tip, err := headCommit(repo)
	if err != nil {
		return nil, err
	}

	// git signatures
	author := historicSign(nc, time)
	committer := historicSign(nc, time)

	// now commit our updated tree to the tip (parent)
	oid, err = repo.CreateCommit("HEAD", author, committer, nc.Message, tree, tip)
	if err != nil {
		return nil, err
	}

	// checkout to keep file system in sync with git
	err = repo.CheckoutHead(
		&git.CheckoutOpts{Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutForce},
	)

	return oid, err

}

func historicSign(nc NewCommit, time time.Time) *git.Signature {
	return &git.Signature{
		Name:  nc.Name,
		Email: nc.Email,
		When:  time,
	}
}
