package main

import (
	"gopkg.in/libgit2/git2go.v25"
)

func getCommits() (commits []Commit, err error) {
	repo, err := repository(config)
	if err != nil {
		return commits, err
	}

	commits, err = allCommits(repo)

	return commits, err
}

// allCommits returns all commits for the current branch (master)
func allCommits(repo *git.Repository) (commits []Commit, err error) {

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

	return commits, nil
}
