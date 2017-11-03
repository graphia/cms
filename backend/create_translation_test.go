package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/libgit2/git2go.v25"
)

func Test_createTranslation(t *testing.T) {

	repoPath := "../tests/tmp/repositories/create_translation"
	lr, _ := setupSmallTestRepo(repoPath)

	multilingualConfig := "../config/test/multilingual.yml"
	config, _ = loadConfig(&multilingualConfig)

	// loadConfig overwrites repo config written by
	// setupSmallTestRepo, so manually set it back to
	// repoPath
	config.Repository = filepath.Join(repoPath)

	type args struct {
		nt   NewTranslation
		user User
	}
	tests := []struct {
		name      string
		args      args
		wantOid   *git.Oid
		wantErr   bool
		errMsg    string
		commitMsg string
	}{
		{
			name: "Non-enabled language code",
			args: args{
				nt: NewTranslation{
					Path:           "documents",
					SourceFilename: "document_1.md",
					LanguageCode:   "no",
					RepositoryInfo: RepositoryInfo{LatestRevision: lr.String()},
				},
				user: mh,
			},
			wantErr: true,
			errMsg:  "Language 'no' not enabled",
		},
		{
			name: "Translation file created properly",
			args: args{
				nt: NewTranslation{
					Path:           "documents",
					SourceFilename: "document_1.md",
					LanguageCode:   "fi",
					RepositoryInfo: RepositoryInfo{LatestRevision: lr.String()},
				},
				user: mh,
			},
			commitMsg: "Finnish translation initiated",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oid, fn, err := createTranslation(tt.args.nt, tt.args.user)

			if tt.wantErr {
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}

			// check target filename contains language code immediately prior to ".md"
			tfn := tt.args.nt.TargetFilename()
			assert.Contains(t, tfn, fmt.Sprintf("%s.md", tt.args.nt.LanguageCode))
			assert.Equal(t, fn, tfn)

			// check target file exists
			_, err = os.Stat(filepath.Join(repoPath, tt.args.nt.Path, tfn))
			assert.False(t, os.IsNotExist(err))

			// check source and target are equal
			source, _ := ioutil.ReadFile(filepath.Join(repoPath, tt.args.nt.Path, tt.args.nt.SourceFilename))
			target, _ := ioutil.ReadFile(filepath.Join(repoPath, tt.args.nt.Path, tfn))
			assert.Equal(t, source, target)

			// check the commit message is correct
			repo, _ := repository(config)
			lastCommit, _ := repo.LookupCommit(oid)
			assert.Equal(t, tt.commitMsg, lastCommit.Message())

		})
	}
}
