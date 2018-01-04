package main

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/graphia/particle"
	"github.com/stretchr/testify/assert"
	"gopkg.in/libgit2/git2go.v25"
)

func TestGetFilesInDocumentsDir(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	// using the small test repo, there are three documents
	filesExpected := 3

	files, err := getFilesInDir("documents")
	if err != nil {
		t.Error("error", err)
	}
	filesCount := len(files)
	assert.Equal(t, filesExpected, filesCount)
}

func TestGetFilesInAppendicesDir(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	// using the small test repo, there are two appendices
	filesExpected := 2

	// ensure a third non-Markdown file exists in the repo but won't be returned
	_, err := os.Stat(filepath.Join(repoPath, "appendices", "another_file.txt"))
	assert.False(t, os.IsNotExist(err))

	files, err := getFilesInDir("appendices")
	if err != nil {
		t.Error("error", err)
	}
	filesCount := len(files)
	assert.Equal(t, filesExpected, filesCount)

}

func TestGetFilesInDirContents(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	files, _ := getFilesInDir("appendices")

	file := files[0]

	// file attributes
	assert.Equal(t, "index.md", file.Filename)
	assert.Equal(t, "appendix_1", file.Document)
	assert.Equal(t, "appendices", file.Path)

	// frontmattter metadata
	assert.Equal(t, "Appendix 1", file.FrontMatter.Title)
	assert.Equal(t, "1.1", file.FrontMatter.Version)
	assert.Equal(t, "Arnold Pye", file.FrontMatter.Author)
	assert.Equal(t, []string{"Traffic News", "KBBL TV"}, file.FrontMatter.Tags)
	assert.Equal(t, "The first appendix is the best appendix", file.FrontMatter.Synopsis)
}

func TestGetFilesInNonExistantDir(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	_, err := getFilesInDir("fanfic")
	assert.Contains(t, err.Error(), "directory not found")
}

func TestGetConvertedFile(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getConvertedFile("documents", "document_2", "index.md")
	if err != nil {
		t.Error("error", err)
	}

	assert.Equal(t, file.Filename, "index.md")
	assert.Equal(t, file.Document, "document_2")
	assert.Equal(t, file.Path, "documents")

	// just look for the title rather than rerender md in here
	assert.Contains(t, *file.HTML, "<h1>Document 2</h1>")

	// markdown should be nil
	assert.Nil(t, file.Markdown)
}

func TestGetRawFile(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getRawFile("documents", "document_2", "index.md")
	if err != nil {
		t.Error("error", err)
	}

	raw, _ := ioutil.ReadFile(filepath.Join(
		config.Repository,
		"documents",
		"document_2",
		"index.md",
	))

	contents, err := particle.YAMLEncoding.DecodeString(string(raw), &FrontMatter{})

	assert.Equal(t, file.Filename, "index.md")
	assert.Equal(t, file.Document, "document_2")
	assert.Equal(t, file.Path, "documents")
	assert.Equal(t, *file.Markdown, string(contents))

	// this time, HTML should be nil
	assert.Nil(t, file.HTML)
}

func TestGetFileBothMarkdownAndHTML(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getFile("documents", "document_2", "index.md", true, true)
	if err != nil {
		t.Error("error", err)
	}

	raw, _ := ioutil.ReadFile(filepath.Join(
		config.Repository,
		"documents",
		"document_2",
		"index.md",
	))

	contents, err := particle.YAMLEncoding.DecodeString(string(raw), &FrontMatter{})

	repo, _ := repository(config)
	hc, _ := headCommit(repo)
	repoInfo := RepositoryInfo{LatestRevision: hc.Id().String()}

	assert.Equal(t, file.Filename, "index.md")
	assert.Equal(t, file.Document, "document_2")
	assert.Equal(t, file.Path, "documents")
	assert.Equal(t, *file.Markdown, string(contents))
	assert.Contains(t, *file.HTML, "<h1>Document 2</h1")
	assert.Equal(t, *file.RepositoryInfo, repoInfo)

}

func TestGetFileNeitherMarkdownOrHTML(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getFile("documents", "document_2", "index.md", false, false)
	if err != nil {
		t.Error("error", err)
	}

	repo, _ := repository(config)
	hc, _ := headCommit(repo)
	repoInfo := RepositoryInfo{LatestRevision: hc.Id().String()}

	assert.Equal(t, file.Filename, "index.md")
	assert.Equal(t, file.Document, "document_2")
	assert.Equal(t, file.Path, "documents")
	assert.Nil(t, file.HTML)
	assert.Nil(t, file.Markdown)
	assert.Equal(t, *file.RepositoryInfo, repoInfo)

}

func TestGetFileNoRepoMetadata(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getFile("appendices", "appendix_1", "index.md", false, false)

	repo, _ := repository(config)
	hc, _ := headCommit(repo)
	repoInfo := RepositoryInfo{LatestRevision: hc.Id().String()}

	assert.Nil(t, err) // make sure getFile doesn't return an error
	assert.Equal(t, file.Filename, "index.md")
	assert.Equal(t, file.Document, "appendix_1")
	assert.Equal(t, file.Path, "appendices")
	assert.Nil(t, file.DirectoryInfo)
	assert.Equal(t, *file.RepositoryInfo, repoInfo)

}

func TestListRootDirectories(t *testing.T) {

	// test with subdirectories to ensure we're only returning root
	// directories
	repoPath := "../tests/tmp/repositories/list_directories_subdirs"
	setupSubdirsTestRepo(repoPath)

	directoriesExpected := []Directory{
		Directory{
			Path:          "appendices",
			DirectoryInfo: DirectoryInfo{},
		},
		Directory{
			Path: "documents",
			DirectoryInfo: DirectoryInfo{
				Title:       "Documents",
				Description: "Documents go here",
			},
		},
	}

	directories, err := listRootDirectories()
	if err != nil {
		t.Error("error", err)
	}

	assert.Equal(t, 2, len(directories))
	assert.Equal(t, directoriesExpected, directories)
}

func TestRootDirectorySummary(t *testing.T) {

	// test with subdirectories to ensure we're only returning root
	// directories
	repoPath := "../tests/tmp/repositories/directories_summary"
	setupSubdirsTestRepo(repoPath)

	documentFiles, _ := getFilesInDir("documents")
	appendicesFiles, _ := getFilesInDir("appendices")

	expectedSummary := []DirectorySummary{
		DirectorySummary{
			Path:          "appendices",
			DirectoryInfo: DirectoryInfo{},
			Contents:      appendicesFiles,
		},
		DirectorySummary{
			Path: "documents",
			DirectoryInfo: DirectoryInfo{
				Title:       "Documents",
				Description: "Documents go here",
			},
			Contents: documentFiles,
		},
	}

	summary, err := listRootDirectorySummary()
	if err != nil {
		t.Error("error", err)
	}

	assert.Equal(t, 2, len(summary))
	assert.Equal(t, expectedSummary, summary)
}

func TestCountFiles(t *testing.T) {
	repoPath := "../tests/tmp/repositories/count_files"
	setupMultipleFiletypesTestRepo(repoPath)
	cf, _ := countFiles()

	expectedCounts := map[string]int{
		"images":          3,
		"documents":       6,
		"structured data": 2,
		"tabular data":    1,
		"other":           1,
	}

	assert.Equal(t, expectedCounts, cf)
}

func TestGetMediaType(t *testing.T) {
	type args struct {
		extension string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Joint Photographic Experts Group", args: args{extension: ".jpg"}, want: "image/jpeg"},
		{name: "Portable Network Graphics", args: args{extension: ".png"}, want: "image/png"},
		{name: "Graphics Interchange Format", args: args{extension: ".gif"}, want: "image/gif"},
		{name: "Tagged Image File Format", args: args{extension: ".tiff"}, want: "image/tiff"},

		{name: "Extensible Markup Language", args: args{extension: ".xml"}, want: "text/xml"},
		{name: "Comma Seperated Values", args: args{extension: ".csv"}, want: "text/csv"},
		{name: "JavaScript Object Notation", args: args{extension: ".json"}, want: "text/json"},

		// BMP isn't included in test config, should return 'unknown/[ext]'
		{name: "Bitmap", args: args{extension: ".bmp"}, want: "unknown/bmp"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getMediaType(tt.args.extension))
		})
	}
}

func TestGetAttachments(t *testing.T) {
	repoPath := "../tests/tmp/repositories/count_files"
	setupMultipleFiletypesTestRepo(repoPath)

	attachments, _ := getAttachments("appendices/appendix_1")

	assert.Equal(t, 4, len(attachments))

	//assert.Contains(t assert.TestingT, s interface{}, contains interface{}, msgAndArgs ...interface{})

	// JSON doc
	jsonAttachmentContents, _ := ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "data", "data.json"))
	xmlAttachmentContents, _ := ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "data", "data.xml"))
	pngAttachmentContents, _ := ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "images", "image_1.png"))
	jpegAttachmentContents, _ := ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "images", "image_2.jpg"))

	expectedAttachments := []Attachment{
		Attachment{
			Path:      "appendices/appendix_1/data",
			Extension: ".json",
			MediaType: "text/json",
			Data:      base64.StdEncoding.EncodeToString(jsonAttachmentContents),
			Filename:  "data.json",
		},
		Attachment{
			Path:      "appendices/appendix_1/data",
			Extension: ".xml",
			MediaType: "text/xml",
			Data:      base64.StdEncoding.EncodeToString(xmlAttachmentContents),
			Filename:  "data.xml",
		},
		Attachment{
			Path:      "appendices/appendix_1/images",
			Extension: ".png",
			MediaType: "image/png",
			Data:      base64.StdEncoding.EncodeToString(pngAttachmentContents),
			Filename:  "image_1.png",
		},
		Attachment{
			Path:      "appendices/appendix_1/images",
			Extension: ".jpg",
			MediaType: "image/jpeg",
			Data:      base64.StdEncoding.EncodeToString(jpegAttachmentContents),
			Filename:  "image_2.jpg",
		},
	}

	for _, a := range expectedAttachments {
		t.Run(a.Filename, func(t *testing.T) {
			assert.Contains(t, attachments, a)
		})
	}

}

func TestExtractContents(t *testing.T) {
	repoPath := "../tests/tmp/repositories/extract_contents"
	setupMultipleFiletypesTestRepo(repoPath)

	var gifImage, jpegImage, pngImage []byte

	gifImage, _ = ioutil.ReadFile(filepath.Join(repoPath, "documents", "document_1", "images", "image_1.gif"))
	pngImage, _ = ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "images", "image_1.png"))
	jpegImage, _ = ioutil.ReadFile(filepath.Join(repoPath, "appendices", "appendix_1", "images", "image_2.jpg"))

	type args struct {
		ncf NewCommitFile
	}
	tests := []struct {
		name         string
		args         args
		wantContents []byte
		wantErr      bool
	}{

		{
			name: "Markdown",
			args: args{
				ncf: NewCommitFile{
					FrontMatter: FrontMatter{
						Author:   "Bernice Hibbert",
						Date:     "2016-04-05",
						Draft:    true,
						Slug:     "pangram",
						Synopsis: "Use all of the characters",
						Tags:     nil,
						Title:    "Pangram",
						Version:  "1.0",
					},
					Body:     "the quick *brown* fox jumped over the **lazy** dog",
					Filename: "pangram.md",
				},
			},
			// Multiline string so any leading whitespace remains 😒
			wantContents: []byte(`---
author: Bernice Hibbert
date: 2016-04-05
draft: true
slug: pangram
synopsis: Use all of the characters
tags: []
title: Pangram
version: "1.0"
---

the quick *brown* fox jumped over the **lazy** dog`,
			),
			wantErr: false,
		},

		{
			name: "JavaScript Object Notation",
			args: args{
				ncf: NewCommitFile{
					Body:     `{"hello": "world"}`,
					Filename: "hello-world.json",
				},
			},
			wantContents: []byte("{\"hello\": \"world\"}"),
			wantErr:      false,
		},
		{
			name: "Extensible Markup Language",
			args: args{
				ncf: NewCommitFile{
					Body:     `<?xml version="1.0" encoding="UTF-8"?><hello>world</hello>`,
					Filename: "hello-world.xml",
				},
			},
			wantContents: []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><hello>world</hello>"),
			wantErr:      false,
		},
		{
			name: "Graphics Interchange Format",
			args: args{
				ncf: NewCommitFile{
					Body:          base64.StdEncoding.EncodeToString(gifImage),
					Filename:      "hello-world.gif",
					Base64Encoded: true,
				},
			},
			wantContents: gifImage,
			wantErr:      false,
		},
		{
			name: "Join Photographic Experts Group",
			args: args{
				ncf: NewCommitFile{
					Body:          base64.StdEncoding.EncodeToString(jpegImage),
					Filename:      "hello-world.jpeg",
					Base64Encoded: true,
				},
			},
			wantContents: jpegImage,
			wantErr:      false,
		},
		{
			name: "Portable Network Graphics",
			args: args{
				ncf: NewCommitFile{
					Body:          base64.StdEncoding.EncodeToString(pngImage),
					Filename:      "hello-world.png",
					Base64Encoded: true,
				},
			},
			wantContents: pngImage,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contents, _ := extractContents(tt.args.ncf)
			assert.Equal(t, tt.wantContents, contents)
		})
	}
}

func Test_getMetadata(t *testing.T) {
	repoPath := "../tests/tmp/repositories/get_metadata"
	setupSubdirsTestRepo(repoPath)
	repo, _ := repository(config)

	ht, _ := headTree(repo)

	docs, _ := ht.EntryByPath("documents")
	docsTree, _ := repo.LookupTree(docs.Id)

	appendices, _ := ht.EntryByPath("appendices")
	appendicesTree, _ := repo.LookupTree(appendices.Id)

	type args struct {
		repo *git.Repository
		tree *git.Tree
	}
	tests := []struct {
		name    string
		args    args
		wantDi  DirectoryInfo
		wantErr bool
	}{
		{
			name: "Directory with _index.md",
			args: args{
				repo: repo,
				tree: docsTree,
			},
			wantDi: DirectoryInfo{
				Title:       "Documents",
				Description: "Documents go here",
			},
			wantErr: false,
		},
		{
			name: "Directory without _index.md",
			args: args{
				repo: repo,
				tree: appendicesTree,
			},
			wantDi:  DirectoryInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		if !tt.wantErr {
			md, _ := getMetadata(tt.args.repo, tt.args.tree)
			assert.Equal(t, tt.wantDi, md)
		} else {
			_, err := getMetadata(tt.args.repo, tt.args.tree)
			assert.Equal(t, ErrMetadataNotFound, err)
		}
	}
}

func Test_getMetadataFromBlob(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_metadata"
	setupMissingFrontmatterTestRepo(repoPath)
	repo, _ := repository(config)
	ht, _ := headTree(repo)

	entryFullFM, _ := ht.EntryByPath("documents/full-frontmatter/index.md")
	blobFullFM, _ := repo.LookupBlob(entryFullFM.Id)

	entryNoFM, _ := ht.EntryByPath("documents/no-frontmatter/index.md")
	blobNoFM, _ := repo.LookupBlob(entryNoFM.Id)

	entryBrokenFM, _ := ht.EntryByPath("documents/broken-frontmatter/index.md")
	blobBrokenFM, _ := repo.LookupBlob(entryBrokenFM.Id)

	type args struct {
		blob *git.Blob
	}
	tests := []struct {
		name    string
		args    args
		wantFm  FrontMatter
		wantErr bool
	}{
		{
			name: "documents/full-frontmatter/index.md",
			args: args{
				blob: blobFullFM,
			},
			wantFm: FrontMatter{
				Title:    "Document 1",
				Author:   "Gil Gunderson",
				Slug:     "document-1",
				Synopsis: "I brought that wall from home",
				Tags:     []string{"ol'", "gil", "ol' gil"},
			},
			wantErr: false,
		},
		{
			name: "documents/no-frontmatter/index.md",
			args: args{
				blob: blobNoFM,
			},
			wantFm:  FrontMatter{},
			wantErr: false,
		},
		{
			name: "documents/broken-frontmatter/index.md",
			args: args{
				blob: blobBrokenFM,
			},
			wantFm:  FrontMatter{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFm, _ := getMetadataFromBlob(tt.args.blob)
			assert.Equal(t, gotFm, tt.wantFm)
		})
	}
}

func Test_translationFilename(t *testing.T) {
	config.DefaultLanguage = "en"

	type args struct {
		fn   string
		code string
	}
	tests := []struct {
		name    string
		args    args
		wantTfn string
		wantErr bool
	}{
		{
			name:    "No language code",
			args:    args{fn: "test.md", code: "sv"},
			wantTfn: "test.sv.md",
		},
		{
			name:    "Enabled language code",
			args:    args{fn: "test.fi.md", code: "sv"},
			wantTfn: "test.sv.md",
		},
		{
			name:    "Enabled language code",
			args:    args{fn: "test.fi.md", code: "sv"},
			wantTfn: "test.sv.md",
		},
		{
			name:    "Default language code",
			args:    args{fn: "test.fi.md", code: "en"},
			wantTfn: "test.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTfn := translationFilename(tt.args.fn, tt.args.code)

			assert.Equal(t, tt.wantTfn, gotTfn)
		})
	}
}

func Test_getTranslations(t *testing.T) {

	multilingualConfig := "../config/test/multilingual.yml"

	repoPath := "../tests/tmp/repositories/get_translations"
	setupTranslationsTestRepo(repoPath)

	config, _ = loadConfig(&multilingualConfig)
	config.Repository = repoPath
	repo, _ := repository(config)

	type args struct {
		repo      *git.Repository
		directory string
		document  string
		filename  string
	}
	tests := []struct {
		name      string
		args      args
		wantLangs []string
		wantErr   bool
		errMsg    string
	}{
		{
			name: "document_1",
			args: args{
				repo:      repo,
				directory: "documents",
				document:  "document_1",
				filename:  "index.md",
			},
			wantLangs: []string{"en", "sv"},
		},
		{
			name: "document_2",
			args: args{
				repo:      repo,
				directory: "documents",
				document:  "document_2",
				filename:  "index.md",
			},
			wantLangs: []string{"en", "fi", "sv"},
		},
		{
			name: "document_3",
			args: args{
				repo:      repo,
				directory: "documents",
				document:  "document_3",
				filename:  "index.md",
			},
			wantLangs: []string{"en", "fi"},
		},
		{
			name: "missing document",
			args: args{
				repo:      repo,
				directory: "documents",
				document:  "missing_document",
				filename:  "index.md",
			},
			wantErr: true,
			errMsg:  "No translations found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLangs, err := getTranslations(tt.args.repo, tt.args.directory, tt.args.document, tt.args.filename)
			if tt.wantErr {
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}
			assert.Equal(t, tt.wantLangs, gotLangs)
		})
	}
}

func Test_fileExists(t *testing.T) {

	repoPath := "../tests/tmp/repositories/file_exists"
	setupTranslationsTestRepo(repoPath)
	repo, _ := repository(config)

	type args struct {
		repo     *git.Repository
		path     string
		document string
		filename string
	}
	tests := []struct {
		name       string
		args       args
		wantExists bool
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "Existing file",
			args:       args{repo: repo, path: "documents", document: "document_1", filename: "index.md"},
			wantExists: true,
		},
		{
			name:       "Non-existing file",
			args:       args{repo: repo, path: "documents", document: "document_1", filename: "index.de.md"},
			wantExists: false,
			wantErr:    true,
			errMsg:     "the path 'index.de.md' does not exist in the given tree",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, err := fileExists(tt.args.repo, tt.args.path, tt.args.document, tt.args.filename)

			assert.Equal(t, tt.wantExists, gotExists)

			if tt.wantErr {
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}
