package main

import (
	"bytes"
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
		commit = Commit{
			Message:    c.Summary(),
			ID:         c.Id().String(),
			ObjectType: c.Type().String(),
			Author:     c.Author(),
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

	commitTree, err := commit.Tree()
	if err != nil {
		return cs, err
	}
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
	if err != nil {
		return cs, err
	}

	var buffer bytes.Buffer

	for d := 0; d < numDeltas; d++ {

		patch, err := gitDiff.Patch(d)
		if err != nil {
			return cs, err
		}

		patchString, err := patch.String()
		if err != nil {
			return cs, err
		}

		buffer.WriteString(fmt.Sprintf("\n%s", patchString))
		patch.Free()
	}

	cs = Changeset{
		NumDeltas: numDeltas,
		Diff:      buffer.String(),
	}

	return cs, nil

}
