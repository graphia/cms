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

	assert.Equal(t, "../tests/data/repositories/small", testConfig.Repository)
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
	os.RemoveAll(path)

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

func createRandomFile(repo *git.Repository, document, lc, msg string) (*git.Oid, error) {

	user := User{
		Name:  "Barney Gumble",
		Email: "barney.gumble@hotmail.com",
	}

	ri, _ := getRepositoryInfo()

	fn := "index.md"

	if lc != "en" {
		fn = fmt.Sprintf("index.%s.md", lc)
	}

	nc := NewCommit{
		Message:        msg,
		RepositoryInfo: ri,

		Files: []NewCommitFile{
			NewCommitFile{
				Body:     "# The quick brown fox\n\njumped over the lazy dog",
				Filename: fn,
				Document: document,
				Path:     "documents",
				FrontMatter: FrontMatter{
					Title:  "Document Twelve",
					Author: "Bernard Arnold Gumble",
				},
			},
		},
	}
	oid, err := createFiles(nc, user)
	return oid, err
}

func TestAllCommits(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_commits"
	_, _ = setupSmallTestRepo(repoPath)

	repo, _ := repository(config)

	msg := "Well well, if it isn't Mr Plow"

	_, err := createRandomFile(repo, "document_12.md", "en", msg)
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
		_, err := createRandomFile(
			repo,
			fmt.Sprintf("document_%d.md", i),
			"en",
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
	assert.Equal(t, 6, cs.NumDeltas)
	assert.Equal(t, 0, cs.NumDeleted)
	assert.Equal(t, 6, cs.NumAdded)
	assert.Contains(t, cs.FullDiff, "+++ b/documents/_index.md")
	assert.Contains(t, cs.FullDiff, "+++ b/documents/document_1/index.md")
	assert.Contains(t, cs.FullDiff, "+++ b/documents/document_2/index.md")
	assert.Contains(t, cs.FullDiff, "+++ b/documents/document_3/index.md")
	assert.Contains(t, cs.FullDiff, "+++ b/appendices/appendix_1/index.md")
	assert.Contains(t, cs.FullDiff, "+++ b/appendices/appendix_2/index.md")

	// Ensure the file contents are included
	assert.Contains(t, cs.FullDiff, "+Lorem ipsum dolor sit")

	var allFilesInRepo = []string{
		"documents/_index.md",
		"documents/document_1/index.md",
		"documents/document_2/index.md",
		"documents/document_3/index.md",
		"appendices/appendix_1/index.md",
		"appendices/appendix_2/index.md",
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

	user := User{
		Name:  "Hyman Krustofski",
		Email: "hyman@springfieldsynagogue.org",
	}

	template := NewCommit{
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
	r1.RepositoryInfo = RepositoryInfo{LatestRevision: oid.String()}
	oid, _ = createFiles(r1, user)
	assert.NotNil(t, oid)

	// Revision 2
	r2 = template
	r2.Files[0].Body = "# r2"
	r2.Message = "r2"
	r2.RepositoryInfo = RepositoryInfo{LatestRevision: oid.String()}
	oid, _ = updateFiles(r2, user)
	assert.NotNil(t, oid)

	// Revision 3
	r3 = template
	r3.Files[0].Body = "# r3"
	r3.Message = "r3"
	r3.RepositoryInfo = RepositoryInfo{LatestRevision: oid.String()}
	oid, _ = updateFiles(r3, user)
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

	user := User{
		Name:  "Hyman Krustofski",
		Email: "hyman@springfieldsynagogue.org",
	}

	template := NewCommit{
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
	r2.RepositoryInfo = RepositoryInfo{LatestRevision: oid.String()}
	oid, _ = writeHistoricFiles(repo, r2, user, time.Date(2016, 1, 1, 14, 0, 0, 0, time.UTC))
	assert.NotNil(t, oid)

	// Revision 3
	r3 = template
	r3.Files[0].Body = "# r3"
	r3.Message = "r3"
	r3.RepositoryInfo = RepositoryInfo{LatestRevision: oid.String()}
	oid, _ = writeHistoricFiles(repo, r3, user, time.Date(2016, 1, 1, 15, 0, 0, 0, time.UTC))
	assert.NotNil(t, oid)

	// Revision 1
	r1 = template
	r1.Files[0].Body = "# r1"
	r1.Message = "r1"
	r1.RepositoryInfo = RepositoryInfo{LatestRevision: oid.String()}
	oid, _ = writeHistoricFiles(repo, r1, user, time.Date(2016, 1, 1, 13, 0, 0, 0, time.UTC))
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

	user := User{
		Name:  "Hyman Krustofski",
		Email: "hyman@springfieldsynagogue.org",
	}

	template := NewCommit{
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
	r1.RepositoryInfo = RepositoryInfo{LatestRevision: oid.String()}
	oid, _ = createFiles(r1, user)
	assert.NotNil(t, oid)

	// Revision 2 (unrelated)
	r2 = template
	r2.Files[0].Filename = "document_12.md"
	r2.Files[0].Body = "# r2"
	r2.Message = "r2"
	r2.RepositoryInfo = RepositoryInfo{LatestRevision: oid.String()}
	oid, _ = createFiles(r2, user)
	assert.NotNil(t, oid)

	// Revision 3
	r3 = template
	r3.Files[0].Filename = "document_11.md"
	r3.Files[0].Body = "# r3"
	r3.Message = "r3"
	r3.RepositoryInfo = RepositoryInfo{LatestRevision: oid.String()}
	oid, _ = updateFiles(r3, user)
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

func writeHistoricFiles(repo *git.Repository, nc NewCommit, user User, time time.Time) (oid *git.Oid, err error) {

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
	author := historicSign(user, time)
	committer := historicSign(user, time)

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

func historicSign(user User, time time.Time) *git.Signature {
	return &git.Signature{
		Name:  user.Name,
		Email: user.Email,
		When:  time,
	}
}

func Test_canInitializeGitRepository(t *testing.T) {

	// setup actual git repo
	gitRepoPath := "../tests/tmp/repositories/initialize_repo"
	_, _ = setupSmallTestRepo(gitRepoPath)

	// setup empty dir
	emptyDirPath := "../tests/tmp/repositories/empty"
	os.RemoveAll(emptyDirPath)
	os.Mkdir(emptyDirPath, 0777)

	// setup full dir
	fullDirPath := "../tests/tmp/repositories/full"
	os.RemoveAll(fullDirPath)
	CopyDir("../tests/data/repositories/small", fullDirPath)

	// setup file
	filePath := "../tests/tmp/repositories/file"
	os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)

	// setup non-existant target
	nonExistantPath := "../tests/tmp/repositories/non_existant_directory"

	type args struct {
		path string
	}
	type want struct {
		err     bool
		message string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Actual Git Repository",
			args: args{path: gitRepoPath},
			want: want{err: true, message: "git repo already exists at"},
		},
		{
			name: "Empty Directory",
			args: args{path: emptyDirPath},
			want: want{err: false},
		},
		{
			name: "Directory with files",
			args: args{path: fullDirPath},
			want: want{err: false},
		},
		{
			name: "File",
			args: args{path: filePath},
			want: want{err: true, message: "file exists at"},
		},
		{
			name: "Non-existant directory",
			args: args{path: nonExistantPath},
			want: want{err: true, message: "directory does not exist"},
		},
	}

	var err error

	for _, tt := range tests {

		err = canInitializeGitRepository(tt.args.path)

		//assert.Equal(t, tt.want, canInitializeGitRepository(tt.args.path))
		if tt.want.err {
			assert.NotNil(t, err)
		}
	}
}

func Test_initializeGitRepository(t *testing.T) {

	_ = createUser(ck)
	ck, _ := getUserByUsername("cookie.kwan")

	// setup actual git repo
	gitRepoPath := "../tests/tmp/repositories/initialize_repo"
	_, _ = setupSmallTestRepo(gitRepoPath)

	// setup empty dir
	emptyDirPath := "../tests/tmp/repositories/empty"
	os.RemoveAll(emptyDirPath)
	os.Mkdir(emptyDirPath, 0777)

	// setup full dir
	fullDirPath := "../tests/tmp/repositories/full"
	os.RemoveAll(fullDirPath)
	CopyDir("../tests/data/repositories/small", fullDirPath)

	// setup file
	filePath := "../tests/tmp/repositories/file"
	os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)

	// setup non-existant target
	nonExistantPath := "../tests/tmp/repositories/non_existant_directory"

	type args struct {
		path string
	}
	type want struct {
		oid     bool
		err     bool
		message string
	}
	tests := []struct {
		name string
		args args
		want want
	}{

		{
			name: "Actual Git Repository",
			args: args{path: gitRepoPath},
			want: want{oid: false, err: true, message: "cannot initialise repo"},
		},
		{
			name: "Empty Directory",
			args: args{path: emptyDirPath},
			want: want{oid: true, err: false},
		},

		{
			name: "Directory with files",
			args: args{path: fullDirPath},
		},

		{
			name: "File",
			args: args{path: filePath},
			want: want{oid: false, err: true, message: "cannot initialise repo"},
		},
		{
			name: "Non-existant directory",
			args: args{path: nonExistantPath},
			want: want{oid: false, err: true, message: "cannot initialise repo"},
		},
	}

	var oid *git.Oid
	var err error

	for _, tt := range tests {
		//assert.Equal(t, tt.want, canInitializeGitRepository(tt.args.path))
		oid, err = initializeGitRepository(ck, tt.args.path)

		// ensure a commit has been made to the repo
		if tt.want.oid {
			assert.NotNil(t, oid)
		}

		if tt.want.err {
			assert.NotNil(t, err)
			assert.Equal(t, tt.want.message, err.Error())
		}

	}
}

func Test_getRepositoryInfo(t *testing.T) {

	gitRepoPath := "../tests/tmp/repositories/repo_info"
	_, _ = setupSmallTestRepo(gitRepoPath)
	repo, _ := repository(config)
	lr, _ := getLatestRevision(repo)
	expected := RepositoryInfo{LatestRevision: lr.String()}
	actual, _ := getRepositoryInfo()

	assert.Equal(t, expected, actual)
}

func Test_getLatestRevision(t *testing.T) {

	gitRepoPath := "../tests/tmp/repositories/repo_info"
	expected, _ := setupSmallTestRepo(gitRepoPath)
	repo, _ := repository(config)
	actual, _ := getLatestRevision(repo)

	assert.Equal(t, expected, actual)
}

func Test_checkLatestRevision(t *testing.T) {

	gitRepoPath := "../tests/tmp/repositories/repo_info"
	firstOid, _ := setupSmallTestRepo(gitRepoPath)
	repo, _ := repository(config)

	secondOid, _ := createRandomFile(repo, "document_6", "en", "Second commit!")

	type args struct {
		repo *git.Repository
		hash string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "Old Commit",
			args: args{repo: repo, hash: firstOid.String()},
			want: ErrRepoOutOfSync,
		},
		{
			name: "Most recent Commit",
			args: args{repo: repo, hash: secondOid.String()},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkLatestRevision(repo, tt.args.hash)
			assert.Equal(t, tt.want, err)
		})
	}
}

func Test_getCommitsCount(t *testing.T) {

	gitRepoPath := "../tests/tmp/repositories/repo_info"
	_, _ = setupSmallTestRepo(gitRepoPath)

	tests := []struct {
		name    string
		wantQty int
		wantErr bool
	}{
		{
			name:    "One extra commit",
			wantQty: 2,
		},
		{
			name:    "Three extra commits",
			wantQty: 4,
		},
		{
			name:    "Five extra commits",
			wantQty: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo, _ := repository(config)

			for i := 1; i < tt.wantQty; i++ {
				createRandomFile(
					repo,
					fmt.Sprintf("%d.md", i),
					"en",
					fmt.Sprintf("Added file %d", i),
				)
			}

			actualQty, _ := getCommitsCount(repo)
			assert.Equal(t, tt.wantQty, actualQty)

		})
	}
}
