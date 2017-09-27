package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_updateDirectories(t *testing.T) {

	user := User{
		Name:  "Bernice Hibbert",
		Email: "bernice.hibbert@gmail.com",
	}

	type args struct {
		nc   NewCommit
		user User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errMsg  string
	}{
		{
			name: "When a metadata file exists",
			args: args{
				nc: NewCommit{
					Message: "Added beverage information",
					Directories: []NewCommitDirectory{NewCommitDirectory{
						Path: "documents", // there is already a _index.md in documents
						DirectoryInfo: DirectoryInfo{
							Title:       "Buzz Cola",
							Description: "Twice the sugar, twice the caffeine",
							Body:        "# Buzz Cola\nThe taste you'll kill for!",
						},
					}},
				},
				user: user,
			},
		},
		{
			name: "When a metadata file does not exist",
			args: args{
				nc: NewCommit{
					Message: "Added beverage information",
					Directories: []NewCommitDirectory{NewCommitDirectory{
						Path: "appendices", // there is no _index.md in appendices
						DirectoryInfo: DirectoryInfo{
							Title:       "Buzz Cola",
							Description: "Twice the sugar, twice the caffeine",
							Body:        "# Buzz Cola\nThe taste you'll kill for!",
						},
					}},
				},
				user: user,
			},
		},
		{
			name: "When no directories are specified",
			args: args{
				nc: NewCommit{
					Message:     "There's nothing here",
					Directories: []NewCommitDirectory{},
				},
			},
			wantErr: true,
			errMsg:  "at least one directory must be specified",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repoPath := "../tests/tmp/repositories/update_directories"
			setupSmallTestRepo(repoPath)

			oid, err := updateDirectories(tt.args.nc, user)

			if !tt.wantErr {
				repo, _ := repository(config)
				lc, _ := repo.LookupCommit(oid)

				// make sure the committer info is correct
				assert.Equal(t, lc.Committer().Name, user.Name)
				assert.Equal(t, lc.Committer().Email, user.Email)
				assert.Equal(t, lc.Message(), tt.args.nc.Message)

				// make sure the file contents are correct
				for _, f := range tt.args.nc.Directories {
					contents, _ := ioutil.ReadFile(filepath.Join(repoPath, f.Path, "_index.md"))
					assert.Contains(t, string(contents), "---\ntitle: Buzz Cola\ndescription: Twice the sugar, twice the caffeine\n---")
					assert.Contains(t, string(contents), "kill")

				}
			} else {
				assert.Contains(t, err.Error(), tt.errMsg)
			}

		})
	}
}
