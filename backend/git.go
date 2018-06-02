package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/libgit2/git2go.v25"
)

func repository(c Config) (repo *git.Repository, err error) {

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(wd, c.Repository)

	repo, err = git.OpenRepository(path)
	if err != nil {
		return nil, err
	}

	return
}

// headCommit returns the commit object at the repository's head
func headCommit(repo *git.Repository) (commit *git.Commit, err error) {

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("Cannot find repository head (%s)", err)
	}

	commit, err = repo.LookupCommit(head.Target())
	if err != nil {
		return nil, fmt.Errorf("Cannot find commit: %s (%s)", head.Target(), err)
	}

	return
}

// headTree returns the tree belonging to the headCommit
func headTree(repo *git.Repository) (tree *git.Tree, err error) {

	var hc *git.Commit

	hc, err = headCommit(repo)
	if err != nil {
		return nil, err
	}

	tree, err = hc.Tree()
	if err != nil {
		return nil, err
	}

	return
}

func getRepositoryInfo() (ri RepositoryInfo, err error) {

	var lr *git.Oid

	repo, err := repository(config)
	if err != nil {
		return ri, err
	}

	lr, err = getLatestRevision(repo)

	return RepositoryInfo{LatestRevision: lr.String()}, err
}

func getLatestRevision(repo *git.Repository) (oid *git.Oid, err error) {

	hc, err := headCommit(repo)

	if err != nil && err != ErrMetadataNotFound {
		Error.Println("Could not retrieve headCommit", err.Error())
		return nil, err
	}

	return hc.Id(), err

}

func checkLatestRevision(repo *git.Repository, hash string) error {
	var lr *git.Oid
	var err error

	if hash == "" {
		return fmt.Errorf("No hash provided")
	}

	lr, err = getLatestRevision(repo)
	if err != nil {
		return err
	}

	if lr.String() != hash {
		return ErrRepoOutOfSync
	}

	return nil
}

func getCommits(qty int) (commits []Commit, err error) {
	repo, err := repository(config)
	if err != nil {
		return commits, err
	}

	commits, err = allCommits(repo, qty)

	return commits, err
}

// allCommits returns all commits for the current branch (master)
func allCommits(repo *git.Repository, qty int) (commits []Commit, err error) {

	var commit Commit

	hc, err := headCommit(repo)
	if err != nil {
		return commits, err
	}

	revWalk, err := repo.Walk()
	if err != nil {
		return commits, err
	}

	err = revWalk.Push(hc.Id())

	revWalkIterator := func(c *git.Commit) bool {

		Debug.Println("Commited at:", c.Author().When)
		commit = Commit{
			Summary:    c.Summary(),
			Message:    c.Message(),
			ID:         c.Id().String(),
			ObjectType: c.Type().String(),
			Author:     c.Author(),
			Time:       c.Author().When,
		}
		commits = append(commits, commit)
		return true
	}

	err = revWalk.Iterate(revWalkIterator)
	if err != nil {
		return commits, err
	}

	if qty > 0 && qty <= len(commits) {
		return commits[0:qty], nil
	}

	return commits, nil
}

func countCommits() (qty int, err error) {
	repo, err := repository(config)
	if err != nil {
		return qty, err
	}
	qty, err = getCommitsCount(repo)
	return qty, err
}

func getCommitsCount(repo *git.Repository) (qty int, err error) {
	hc, err := headCommit(repo)
	if err != nil {
		return qty, err
	}

	revWalk, err := repo.Walk()
	if err != nil {
		return qty, err
	}
	err = revWalk.Push(hc.Id())

	revWalkIterator := func(c *git.Commit) bool {
		qty++
		return true
	}

	err = revWalk.Iterate(revWalkIterator)

	return qty, err
}

func diffForCommit(hash string) (cs Changeset, err error) {
	repo, err := repository(config)

	commitOid, err := git.NewOid(hash)
	if err != nil {
		return cs, err
	}
	commit, err := repo.LookupCommit(commitOid)
	if err != nil {
		return cs, err
	}
	defer commit.Free()

	commitTree, err := commit.Tree()
	if err != nil {
		return cs, err
	}
	defer commitTree.Free()

	options, err := git.DefaultDiffOptions()
	if err != nil {
		return cs, err
	}
	options.IdAbbrev = 40

	var parentTree *git.Tree
	if commit.ParentCount() > 0 {
		parentTree, err = commit.Parent(0).Tree()
		if err != nil {
			return cs, err
		}
	}

	gitDiff, err := repo.DiffTreeToTree(parentTree, commitTree, &options)
	if err != nil {
		return cs, err
	}

	// Show all file patch diffs in a commit.
	numDeltas, err := gitDiff.NumDeltas()

	numDiffs := 0
	numAdded := 0
	numDeleted := 0

	var buffer bytes.Buffer

	var files map[string]ChangesetFiles

	files = make(map[string]ChangesetFiles)
	var csf ChangesetFiles

	err = gitDiff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {

		patch, err := gitDiff.Patch(numDiffs)
		// increment counter as we loop through
		numDiffs++

		if err != nil {
			return nil, err
		}

		patchString, err := patch.String()
		if err != nil {
			return nil, err
		}

		buffer.WriteString(fmt.Sprintf("\n%s", patchString))
		patch.Free()

		switch file.Status {
		case git.DeltaAdded:
			numAdded++
		case git.DeltaDeleted:
			numDeleted++
		}

		var old, new []byte

		old, err = getFileContentsByOid(repo, file.OldFile.Oid)
		if err != nil {
			return nil, err
		}
		new, err = getFileContentsByOid(repo, file.NewFile.Oid)
		if err != nil {
			return nil, err
		}

		csf = ChangesetFiles{}

		if hasImageExt(file.OldFile.Path) {
			csf.Old = base64.StdEncoding.EncodeToString(old)
		} else {
			csf.Old = string(old)
		}

		if hasImageExt(file.NewFile.Path) {
			csf.New = base64.StdEncoding.EncodeToString(new)
		} else {
			csf.New = string(new)
		}

		files[file.NewFile.Path] = csf

		return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
			return func(line git.DiffLine) error {
				return nil
			}, nil
		}, nil

	}, git.DiffDetailLines)

	if err != nil {
		return cs, err
	}

	cs = Changeset{
		NumDeltas:  numDeltas,
		NumAdded:   numAdded,
		NumDeleted: numDeleted,
		FullDiff:   buffer.String(),
		Files:      files,
		Message:    commit.Message(),
		Author:     commit.Author(),
		Hash:       commit.Id().String(),
		Time:       commit.Committer().When,
	}

	return cs, nil

}

func getFileContentsByOid(repo *git.Repository, oid *git.Oid) (contents []byte, err error) {

	// A hash of 40 zeroes means no file is expected, return nil
	if oid.String() == "0000000000000000000000000000000000000000" {
		return nil, err
	}

	blob, err := repo.LookupBlob(oid)
	if err != nil {
		return contents, err
	}
	contents = blob.Contents()

	return contents, err
}

func getFileHistory(path string, size int) (history []HistoricCommit, err error) {
	repo, err := repository(config)
	if err != nil {
		return history, err
	}

	history, err = lookupFileHistory(repo, path, size)

	return history, err
}

func lookupFileHistory(repo *git.Repository, path string, size int) ([]HistoricCommit, error) {

	if len(path) == 0 {
		return nil, nil
	}
	var err error

	revwalk, err := repo.Walk()
	if err != nil {
		return nil, err
	}
	defer revwalk.Free()

	err = revwalk.PushHead()
	if err != nil {
		return nil, err
	}

	revwalk.Sorting(git.SortTime)

	var entry *git.TreeEntry
	var fh []HistoricCommit

	err = revwalk.Iterate(func(commit *git.Commit) bool {
		defer commit.Free()

		tree, err := commit.Tree()
		if err != nil {
			Error.Println("Failed to extract tree from commit", tree.Id())
			return false
		}
		defer tree.Free()

		entry, err = tree.EntryByPath(path)
		if err != nil {
			Warning.Printf("Cannot find entry '%s' in tree '%s'", path, tree.Id())
		}

		var hc HistoricCommit

		if entry != nil && err == nil {

			// make sure the file (entry) we're looking for is present in this commit
			if len(fh) > 0 && fh[len(fh)-1].EntryID == entry.Id.String() {
				fh = fh[:len(fh)-1]
			}

			changeset, err := diffForCommit(commit.Id().String())
			if err != nil {
				Warning.Println("couldn't get diff for commmit", commit.Id().String())
				return false
			}

			hc = HistoricCommit{
				ID:      commit.Id().String(),
				EntryID: entry.Id.String(),
				Message: commit.Message(),
				Author:  commit.Author(),
				Time:    commit.Author().When,
				Old:     changeset.Files[path].Old,
				New:     changeset.Files[path].New,
			}

			fh = append(fh, hc)

			if size > 0 && len(fh) >= size {
				Info.Println("History limit reached, exiting")
				return false
			}
		}

		return true
	})

	return fh, nil

}

func canInitializeGitRepository(path string) error {

	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("directory does not exist '%s'", path)
	}

	if !stat.IsDir() {
		return fmt.Errorf("file exists at '%s'", path)
	}

	_, err = git.OpenRepository(path)
	if err == nil {
		return fmt.Errorf("git repo already exists at '%s'", path)
	}

	// TODO
	// Q. what do we do if it's a subdirectory of a git repo?
	// A. nothing - for the time being! 😎

	Warning.Printf("nothing obstructing git repo initialisation at '%s'", path)
	return nil
}

func initializeGitRepository(user User, path string) (oid *git.Oid, err error) {

	err = canInitializeGitRepository(path)

	if err != nil {
		return oid, fmt.Errorf("cannot initialise repo")
	}

	repo, err := git.InitRepository(path, false)
	if err != nil {
		return oid, err
	}

	idx, err := repo.Index()
	if err != nil {
		return oid, err
	}

	err = idx.AddAll([]string{}, git.IndexAddForce, nil)
	if err != nil {
		return oid, err
	}

	err = idx.Write()
	if err != nil {
		return oid, err
	}

	treeID, err := idx.WriteTree()
	if err != nil {
		return oid, err
	}

	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return oid, err
	}

	// git signatures
	author := sign(user)
	committer := sign(user)
	message := "Initialising repository"

	oid, err = repo.CreateCommit("HEAD", author, committer, message, tree)
	if err != nil {
		return oid, err
	}

	// checkout to keep file system in sync with git
	err = repo.CheckoutHead(
		&git.CheckoutOpts{Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutForce},
	)

	return oid, err
}
