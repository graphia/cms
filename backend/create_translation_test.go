package main

import (
	"fmt"
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
	config.Repository = repoPath

	type args struct {
		nt   NewTranslation
		user User
	}
	tests := []struct {
		name            string
		args            args
		wantOid         *git.Oid
		wantErr         bool
		errMsg          string
		commitMsg       string
		createDuplicate bool
	}{
		{
			name: "Non-enabled language code",
			args: args{
				nt: NewTranslation{
					Path:           "documents",
					SourceFilename: "index.md",
					SourceDocument: "document_1",
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
					SourceFilename: "index.md",
					SourceDocument: "document_1",
					LanguageCode:   "fi",
					RepositoryInfo: RepositoryInfo{LatestRevision: lr.String()},
				},
				user: mh,
			},
			commitMsg: "Finnish translation initiated",
		},
		{
			name: "Translation already exists",
			args: args{
				nt: NewTranslation{
					Path:           "documents",
					SourceFilename: "index.md",
					SourceDocument: "document_1",
					LanguageCode:   "fi",
					RepositoryInfo: RepositoryInfo{LatestRevision: lr.String()},
				},
				user: mh,
			},
			createDuplicate: true,
			wantErr:         true,
			errMsg:          "file already exists",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.createDuplicate {
				repo, _ := repository(config)
				_, _ = createRandomFile(repo, "document_1", "fi", "Whoosh")
			}

			oid, fn, err := createTranslation(tt.args.nt, tt.args.user)

			if tt.wantErr {
				Debug.Println("returning!")
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}

			// check target filename contains language code immediately prior to ".md"
			tfn := tt.args.nt.TargetFilename()
			assert.Contains(t, tfn, fmt.Sprintf("%s.md", tt.args.nt.LanguageCode))
			assert.Equal(t, fn, tfn)

			// check target file exists
			_, err = os.Stat(filepath.Join(repoPath, tt.args.nt.Path, tt.args.nt.SourceDocument, tfn))
			assert.False(t, os.IsNotExist(err))

			source, _ := getFile(tt.args.nt.Path, tt.args.nt.SourceDocument, tt.args.nt.SourceFilename, true, false)
			target, _ := getFile(tt.args.nt.Path, tt.args.nt.SourceDocument, tfn, true, false)

			// check source and target markdown hasn't changed and source was not in Draft mode
			assert.Equal(t, source.Markdown, target.Markdown)
			assert.False(t, source.FrontMatter.Draft, false)

			// check that the target has been set to draft
			assert.NotEqual(t, source, target)
			assert.True(t, target.FrontMatter.Draft)

			// check the commit message is correct
			repo, _ := repository(config)
			lastCommit, _ := repo.LookupCommit(oid)
			assert.Equal(t, tt.commitMsg, lastCommit.Message())

		})
	}
}
