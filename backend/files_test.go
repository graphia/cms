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
	assert.Equal(t, "appendix_1.md", file.Filename)
	assert.Equal(t, "appendices", file.Path)
	assert.Equal(t, "appendices/appendix_1.md", file.AbsoluteFilename)

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

	// there isn't a directory called fanfic, so should
	// raise not found error
	_, err := getFilesInDir("fanfic")
	assert.Contains(t, err.Error(), "directory 'fanfic' not found")
}

func TestGetConvertedFile(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getConvertedFile("documents", "document_2.md")
	if err != nil {
		t.Error("error", err)
	}

	assert.Equal(t, file.Filename, "document_2.md")
	assert.Equal(t, file.Path, "documents")

	// just look for the title rather than rerender md in  here
	assert.Contains(t, *file.HTML, "<h1>Document 2</h1>")

	// markdown should be nil
	assert.Nil(t, file.Markdown)
}

func TestGetRawFile(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getRawFile("documents", "document_2.md")
	if err != nil {
		t.Error("error", err)
	}

	raw, _ := ioutil.ReadFile(filepath.Join(
		config.Repository,
		"documents",
		"document_2.md",
	))

	contents, err := particle.YAMLEncoding.DecodeString(string(raw), &FrontMatter{})

	assert.Equal(t, file.Filename, "document_2.md")
	assert.Equal(t, file.Path, "documents")
	assert.Equal(t, *file.Markdown, string(contents))

	// this time, HTML should be nil
	assert.Nil(t, file.HTML)
}

func TestGetFileBothMarkdownAndHTML(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getFile("documents", "document_2.md", true, true)
	if err != nil {
		t.Error("error", err)
	}

	raw, _ := ioutil.ReadFile(filepath.Join(
		config.Repository,
		"documents",
		"document_2.md",
	))

	contents, err := particle.YAMLEncoding.DecodeString(string(raw), &FrontMatter{})

	assert.Equal(t, file.Filename, "document_2.md")
	assert.Equal(t, file.Path, "documents")
	assert.Equal(t, *file.Markdown, string(contents))
	assert.Contains(t, *file.HTML, "<h1>Document 2</h1")

}

func TestGetFileNeitherMarkdownOrHTML(t *testing.T) {

	repoPath := "../tests/tmp/repositories/get_file"
	setupSmallTestRepo(repoPath)

	file, err := getFile("documents", "document_2.md", false, false)
	if err != nil {
		t.Error("error", err)
	}

	assert.Equal(t, file.Filename, "document_2.md")
	assert.Equal(t, file.Path, "documents")
	assert.Nil(t, file.HTML)
	assert.Nil(t, file.Markdown)

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

	expectedSummary := map[string][]FileItem{
		"appendices": appendicesFiles,
		"documents":  documentFiles,
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
		"documents":       5,
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
			Path:             "appendices/appendix_1/data",
			AbsoluteFilename: "appendices/appendix_1/data/data.json",
			Extension:        ".json",
			MediaType:        "text/json",
			Data:             base64.StdEncoding.EncodeToString(jsonAttachmentContents),
			Filename:         "data.json",
		},
		Attachment{
			Path:             "appendices/appendix_1/data",
			AbsoluteFilename: "appendices/appendix_1/data/data.xml",
			Extension:        ".xml",
			MediaType:        "text/xml",
			Data:             base64.StdEncoding.EncodeToString(xmlAttachmentContents),
			Filename:         "data.xml",
		},
		Attachment{
			Path:             "appendices/appendix_1/images",
			AbsoluteFilename: "appendices/appendix_1/images/image_1.png",
			Extension:        ".png",
			MediaType:        "image/png",
			Data:             base64.StdEncoding.EncodeToString(pngAttachmentContents),
			Filename:         "image_1.png",
		},
		Attachment{
			Path:             "appendices/appendix_1/images",
			AbsoluteFilename: "appendices/appendix_1/images/image_2.jpg",
			Extension:        ".jpg",
			MediaType:        "image/jpeg",
			Data:             base64.StdEncoding.EncodeToString(jpegAttachmentContents),
			Filename:         "image_2.jpg",
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
			assert.Equal(t, ErrMetaDataNotFound, err)
		}
	}
}
