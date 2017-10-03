package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildStaticSite(t *testing.T) {

	configPath := "../config/hugo.test.yml"
	repoPath := "../tests/tmp/repositories/publish"
	buildPath := "../tests/tmp/builds/test_build_static_site"

	// make sure the build path is clear
	os.RemoveAll(buildPath)

	_, _ = setupSmallTestRepo(repoPath)

	config.HugoBin = "hugo"
	config.HugoConfigFile = configPath

	output, err := buildStaticSite()
	assert.Nil(t, err)
	assert.Contains(t, string(output), "Started building sites")

	thingsThatShouldExist := []string{"index.html", "index.xml", "documents", "appendices", "sitemap.xml"}

	for _, item := range thingsThatShouldExist {
		t.Run(fmt.Sprintf("publishing creates %s", item), func(t *testing.T) {
			_, err = os.Stat(filepath.Join(buildPath, item))
			assert.False(t, os.IsNotExist(err), fmt.Sprintf("object does not exist: %s", item))
		})
	}

}
