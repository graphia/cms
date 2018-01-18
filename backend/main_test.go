package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// should raise an appropriate error
func TestMissingConfig(t *testing.T) {
	_, err := loadConfig(&invalidPath)
	msg := "no such file or directory"
	assert.Contains(t, err.Error(), msg)
}

// should load the specified file
func TestCustomConfig(t *testing.T) {
	c, err := loadConfig(&smallRepoPath)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, c.HTTPListenPort, "8081")
	assert.Equal(t, c.Repository, "../tests/backend/repositories/small")
}
