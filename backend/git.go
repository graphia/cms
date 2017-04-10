package main

import (
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
