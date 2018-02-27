package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/graphia/particle"
	"gopkg.in/libgit2/git2go.v25"
	yaml "gopkg.in/yaml.v2"
)

var (
	// ErrMetadataNotFound is returned when the _index.md file is missing
	// this isn't a catestrophic problem, but should be logged
	ErrMetadataNotFound = errors.New("_index.md not found")

	// ErrDirectoryNotFound occurs when a directory can't be found in the
	// git repository
	ErrDirectoryNotFound = errors.New("directory not found")

	// ErrRepoOutOfSync occurs when changes are made to the repository
	// between starting to edit and submitting the changes
	ErrRepoOutOfSync = errors.New("repository out of sync")

	// ErrFileAlreadyExists prevents the accidental overwriting
	// of files
	ErrFileAlreadyExists = errors.New("file already exists")
)

// getFilesInDir returns a list of FileItems for listing
func getFilesInDir(directory string) (files []FileItem, err error) {

	// Initialising the slice so json.Marshal returns an empty
	// array instead of `null`
	files = []FileItem{}

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	ht, err := headTree(repo)
	if err != nil {
		return nil, err
	}

	// ensure that the directory exists
	entry, _ := ht.EntryByPath(directory)
	if entry == nil {
		return nil, ErrDirectoryNotFound
	}

	if entry.Type != git.ObjectTree {
		return nil, fmt.Errorf("%s is not a directory", directory)
	}

	tree, err := repo.LookupTree(entry.Id)
	if err != nil {
		return nil, fmt.Errorf("couldn't find tree for entry %s", entry.Id)
	}

	defer tree.Free()

	walkIterator := func(currentDir string, te *git.TreeEntry) int {
		var fm FrontMatter
		var blob *git.Blob

		var ext string

		if te.Type == git.ObjectBlob {

			ext = filepath.Ext(te.Name)

			// skip unless it's a Markdown file
			if ext != ".md" {
				Warning.Println("not a markdown file, skipping:", te.Name)
				return 0
			}

			// skip if it's a directory metadata file `_index.md`
			if te.Name == "_index.md" {
				Warning.Println("is a metadata file, skipping:", te.Name)
				return 0
			}

			blob, err = repo.LookupBlob(te.Id)
			if err != nil {
				Warning.Println("Failed to find blob", te.Id)
				return -1
			}

			fm, err = getMetadataFromBlob(blob)
			if err != nil {
				Warning.Println("Failed to read frontmatter", te.Id)
				return -1
			}

			if err != nil {
				Warning.Println("Failed to decode file", string(blob.Contents()))
				Warning.Println("Frontmatter:", fm)
				return -1
			}

			fi := FileItem{
				Filename:    te.Name,
				Document:    filepath.Clean(currentDir),
				Path:        directory,
				FrontMatter: fm,
			}

			files = append(files, fi)

		}

		return 0
	}

	err = tree.Walk(walkIterator)

	return files, err

}

func createFiles(nc NewCommit, user User) (oid *git.Oid, err error) {

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	ht, err := headTree(repo)
	if err != nil {
		return nil, err
	}

	// check none of the files already exist
	for _, ncf := range nc.Files {

		target := filepath.Join(ncf.Path, ncf.Document, ncf.Filename)

		//FIXME would it make more sense to switch this
		//around and check for the error instead?
		entry, _ := ht.EntryByPath(target)
		if entry != nil {
			return nil, ErrFileAlreadyExists
		}

	}

	oid, err = writeFiles(repo, nc, user)

	return oid, err
}

func updateFiles(nc NewCommit, user User) (oid *git.Oid, err error) {

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	oid, err = writeFiles(repo, nc, user)

	return oid, err
}

func writeFiles(repo *git.Repository, nc NewCommit, user User) (oid *git.Oid, err error) {

	index, err := repo.Index()
	if err != nil {
		Error.Println("Failed to get repo index", err.Error())
		return nil, err
	}
	defer index.Free()

	err = checkLatestRevision(repo, nc.RepositoryInfo.LatestRevision)
	if err != nil {
		return oid, err
	}

	var contents []byte

	for _, ncf := range nc.Files {

		var ie git.IndexEntry

		// get the file contents in the correct format
		contents, err = extractContents(ncf)
		if err != nil {
			Error.Println("Failed to extract contents", err.Error())
			return oid, err
		}

		oid, err = repo.CreateBlobFromBuffer(contents)
		if err != nil {
			Error.Println("Failed to create blob from buffer", err.Error())
			return oid, err
		}

		// build the git index entry and add it to the index
		ie = buildIndexEntry(oid, ncf)

		err = index.Add(&ie)
		if err != nil {
			return oid, err
		}

	}

	oid, err = writeTreeAndCommit(repo, index, nc.Message, user)

	return oid, err

}

func writeMetadataFiles(repo *git.Repository, nc NewCommit, user User) (oid *git.Oid, err error) {

	index, err := repo.Index()
	if err != nil {
		return nil, err
	}
	defer index.Free()

	var meta []byte

	for _, ncd := range nc.Directories {

		var ie git.IndexEntry

		body := []byte(ncd.DirectoryInfo.Body)

		if (ncd.DirectoryInfo != DirectoryInfo{}) {
			meta = make([]byte, particle.YAMLEncoding.EncodeLen(body, &ncd.DirectoryInfo))
			particle.YAMLEncoding.Encode(meta, body, &ncd.DirectoryInfo)
		}

		oid, err = repo.CreateBlobFromBuffer(meta)
		if err != nil {
			return nil, err
		}

		// build the git index entry and add it to the index
		ie = buildIndexEntryDirectory(oid, ncd)

		err = index.Add(&ie)
		if err != nil {
			return nil, err
		}

	}

	oid, err = writeTreeAndCommit(repo, index, nc.Message, user)
	if err != nil {
		return oid, err
	}

	return oid, err

}

func listRootDirectories() (directories []Directory, err error) {

	// Initialising the slice so json.Marshal returns an empty
	// array instead of `null`
	directories = []Directory{}

	repo, err := repository(config)
	if err != nil {
		return directories, err
	}
	defer repo.Free()

	ht, err := headTree(repo)
	if err != nil {
		return directories, err
	}

	defer ht.Free()

	walkIterator := func(_ string, te *git.TreeEntry) int {

		if te.Type == git.ObjectTree {

			// check for _index.md file
			tree, err := repo.LookupTree(te.Id)
			if err != nil {
				return 0
			}

			di, err := getMetadata(repo, tree)

			// if there is any kind of error except ErrMetadataNotFound,
			// something's wrong, quit
			if err != ErrMetadataNotFound && err != nil {
				Error.Println("Metadata found but not retrievable", err.Error())
				return 0
			}

			directories = append(directories, Directory{
				Path:          te.Name,
				DirectoryInfo: di,
			})

			return 1

		}

		return 0
	}

	err = ht.Walk(walkIterator)

	return directories, err

}

func listRootDirectorySummary() (summary []DirectorySummary, err error) {

	// Initialise the slice so we get an empty array instead of null
	summary = []DirectorySummary{}
	var contents []FileItem

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	ht, err := headTree(repo)
	if err != nil {
		return nil, err
	}

	defer ht.Free()

	walkIterator := func(_ string, te *git.TreeEntry) int {

		if te.Type == git.ObjectTree {

			contents, err = getFilesInDir(te.Name)
			if err != nil {
				Error.Println("Failed to retrieve files when generating summary", te.Name, err)
				return 0
			}

			// check for _index.md file
			tree, err := repo.LookupTree(te.Id)
			if err != nil {
				return 0
			}

			di, err := getMetadata(repo, tree)

			summary = append(summary, DirectorySummary{
				Path:          te.Name,
				DirectoryInfo: di,
				Contents:      contents,
			})

			return 1

		}

		return 0
	}

	err = ht.Walk(walkIterator)

	return summary, err
}

func createTranslation(nt NewTranslation, user User) (oid *git.Oid, target string, err error) {

	repo, err := repository(config)
	target = nt.TargetFilename()

	exists, err := fileExists(repo, nt.Path, nt.SourceDocument, target)
	if exists {
		return oid, target, ErrFileAlreadyExists
	}

	err = checkLatestRevision(repo, nt.RepositoryInfo.LatestRevision)
	if err != nil {
		return oid, target, err
	}

	language, err := getLanguage(nt.LanguageCode)
	if err != nil {
		return oid, target, err
	}

	sf, err := getFile(nt.Path, nt.SourceDocument, nt.SourceFilename, true, true)
	if err != nil {
		return oid, target, err
	}

	index, err := repo.Index()
	if err != nil {
		Error.Println("Failed to get repo index", err.Error())
		return oid, target, err
	}
	defer index.Free()

	// set the new translation to draft
	sf.FrontMatter.Draft = true
	contents := sf.ToMarkdown()

	boid, err := repo.CreateBlobFromBuffer(contents)
	if err != nil {
		return oid, target, err
	}

	ie := buildIndexEntryTranslation(boid, nt, len(contents))

	err = index.Add(&ie)
	if err != nil {
		return oid, target, err
	}

	msg := fmt.Sprintf("%s translation initiated", language.Name)

	oid, err = writeTreeAndCommit(repo, index, msg, user)

	return oid, target, err
}

func createDirectories(nc NewCommit, user User) (oid *git.Oid, err error) {

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}

	if len(nc.Directories) == 0 {
		return nil, fmt.Errorf("at least one new directory must be specified")
	}

	var addedDirs []string

	// git can't track empty directories, so, Rails-style, we'll add an
	// empty file called .keep for each directory and ensure the body is blank
	for _, ncd := range nc.Directories {
		addedDirs = append(addedDirs, ncd.Path)
	}

	// And set the commit message sensibly so we don't need to prompt
	// the user
	nc.Message = fmt.Sprintf("Added directories: %s", strings.Join(addedDirs, ","))

	oid, err = writeDirectories(repo, nc, user)
	if err != nil {
		return nil, err
	}

	return oid, err
}

func updateDirectories(nc NewCommit, user User) (oid *git.Oid, err error) {

	// loop through nc files modifying the metadata
	// if the _index.md file does not exist, it should create it!
	repo, err := repository(config)
	if err != nil {
		return oid, err
	}

	// make sure that the dirs included in the nc are in the commit

	if len(nc.Directories) == 0 {
		return oid, fmt.Errorf("at least one directory must be specified")
	}

	oid, err = writeMetadataFiles(repo, nc, user)
	if err != nil {
		Error.Printf("Failed to write metadata files %s", err.Error())
		return oid, err
	}

	return oid, err
}

func writeDirectories(repo *git.Repository, nc NewCommit, user User) (oid *git.Oid, err error) {
	index, err := repo.Index()
	if err != nil {
		return nil, err
	}
	defer index.Free()

	for _, ncd := range nc.Directories {

		target := filepath.Join(ncd.Path)
		absoluteTarget := filepath.Join(config.Repository, target)

		_, err = os.Stat(absoluteTarget)
		if err == nil {
			return nil, fmt.Errorf("directory already exists %s", target)
		}

		if ncd.Path == "" {
			return nil, fmt.Errorf("path must be specified when creating a directory: %s", ncd)
		}

		var meta = []byte("")
		body := []byte(ncd.DirectoryInfo.Body)

		var ie git.IndexEntry

		// if we have some DirectoryInfo metadata, overwrite meta with it in the
		// usual FrontMatter manner

		if (ncd.DirectoryInfo != DirectoryInfo{}) {
			meta = make([]byte, particle.YAMLEncoding.EncodeLen(body, &ncd.DirectoryInfo))
			particle.YAMLEncoding.Encode(meta, body, &ncd.DirectoryInfo)
		}

		oid, err = repo.CreateBlobFromBuffer(meta)
		if err != nil {
			return nil, err
		}

		// build the git index entry and add it to the index
		ie = buildIndexEntryDirectory(oid, ncd)

		err = index.Add(&ie)
		if err != nil {
			return nil, err
		}

	}

	oid, err = writeTreeAndCommit(repo, index, nc.Message, user)

	if err != nil {
		return oid, err
	}

	return oid, err

}

func deleteDirectories(nc NewCommit, user User) (oid *git.Oid, err error) {

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	ht, err := headTree(repo)
	if err != nil {
		return nil, err
	}

	// first grab the repo's index
	index, err := repo.Index()
	if err != nil {
		return nil, err
	}

	err = checkLatestRevision(repo, nc.RepositoryInfo.LatestRevision)
	if err != nil {
		return oid, err
	}

	for _, ncd := range nc.Directories {

		// ensure that the directory exists before we try to delete it
		d, _ := ht.EntryByPath(ncd.Path)
		if d == nil {
			return nil, fmt.Errorf("directory does not exist: %s", ncd.Path)
		}

		// and remove the target by path and everything beneath it
		Debug.Println("Removing directory:", ncd.Path)
		err = index.RemoveDirectory(ncd.Path, 0)
		if err != nil {
			return nil, err
		}

	}

	oid, err = writeTreeAndCommit(repo, index, nc.Message, user)
	if err != nil {
		return oid, err
	}

	return oid, err

}

func deleteFiles(nc NewCommit, user User) (oid *git.Oid, err error) {

	repo, err := repository(config)
	if err != nil {
		return oid, err
	}
	defer repo.Free()

	err = checkLatestRevision(repo, nc.RepositoryInfo.LatestRevision)
	if err != nil {
		return oid, err
	}

	ht, err := headTree(repo)
	if err != nil {
		return oid, err
	}

	// first grab the repo's index
	index, err := repo.Index()
	if err != nil {
		return oid, err
	}
	defer index.Free()

	for _, ncf := range nc.Files {

		target := filepath.Join(ncf.Path, ncf.Document, ncf.Filename)

		// ensure that the file exists before we try to delete it
		file, err := ht.EntryByPath(target)

		if file == nil {
			return oid, fmt.Errorf("file does not exist %s", target)
		}

		if err != nil {
			return oid, err
		}

		// and remove the target by path
		err = index.RemoveByPath(target)
		if err != nil {
			return oid, err
		}

	}

	// FIXME moving to Bundles this will *also* delete the translations
	// this is the correct behaviour but need to make it clear in the UI
	//
	// if we're deleting the accompanying attachment directories too
	for _, ncd := range nc.Directories {

		// ensure that the directory exists before we try to delete it
		// if it doesn't, continue rather than the usual error
		d, _ := ht.EntryByPath(ncd.Path)
		if d == nil {
			Warning.Println("directory does not exist", ncd.Path)
			continue
		}

		// and remove the target by path and everything beneath it
		err = index.RemoveDirectory(ncd.Path, 0)
		if err != nil {
			Error.Println("cannot delete directory", ncd.Path, err.Error())
			return nil, err
		}

	}

	// final check, if no commit message supplied use a generic one
	if nc.Message == "" {
		nc.Message = "File deleted"
	}

	oid, err = writeTreeAndCommit(repo, index, nc.Message, user)
	if err != nil {
		return oid, err
	}

	return oid, err

}

func sign(user User) *git.Signature {
	return &git.Signature{
		Name:  user.Name,
		Email: user.Email,
		When:  time.Now(),
	}
}

func buildIndexEntry(oid *git.Oid, ncf NewCommitFile) git.IndexEntry {
	return git.IndexEntry{
		Id:   oid,
		Path: filepath.Join(ncf.Path, ncf.Document, ncf.Filename),
		Size: uint32(len(ncf.Body)),

		Ctime: git.IndexTime{},
		Gid:   uint32(os.Getgid()),
		Uid:   uint32(os.Getuid()),
		Mode:  git.FilemodeBlob,
		Mtime: git.IndexTime{},
	}
}

func buildIndexEntryDirectory(oid *git.Oid, ncd NewCommitDirectory) git.IndexEntry {
	return git.IndexEntry{
		Id:   oid,
		Path: filepath.Join(ncd.Path, "_index.md"),
		Size: uint32(0),

		Ctime: git.IndexTime{},
		Gid:   uint32(os.Getgid()),
		Uid:   uint32(os.Getuid()),
		Mode:  git.FilemodeBlob,
		Mtime: git.IndexTime{},
	}
}

func buildIndexEntryTranslation(oid *git.Oid, nt NewTranslation, size int) git.IndexEntry {
	return git.IndexEntry{
		Id:   oid,
		Path: filepath.Join(nt.Path, nt.SourceDocument, nt.TargetFilename()),
		Size: uint32(size),

		Ctime: git.IndexTime{},
		Gid:   uint32(os.Getgid()),
		Uid:   uint32(os.Getuid()),
		Mode:  git.FilemodeBlob,
		Mtime: git.IndexTime{},
	}
}

func getFile(directory, document, filename string, includeMd, includeHTML bool) (file *File, err error) {
	var html, markdown *string = nil, nil
	var fm FrontMatter

	repo, err := repository(config)

	tree, err := headTree(repo)
	if err != nil {
		return nil, err
	}
	target := filepath.Join(directory, document, filename)

	entry, err := tree.EntryByPath(target)
	if err != nil {
		return nil, err
	}

	blob, err := repo.LookupBlob(entry.Id)
	if err != nil {
		return nil, err
	}
	defer blob.Free()

	md, err := particle.YAMLEncoding.DecodeString(string(blob.Contents()), &fm)

	if includeMd {
		str := string(md)
		markdown = &str
	}

	if includeHTML {
		str := string(renderMarkdown(md))
		html = &str
	}

	di, err := getMetadataFromDirectory(directory)
	// if we get any error other than ErrMetadataNotFound,
	// return it, otherwise it's ok and we can continue
	if err != nil && err != ErrMetadataNotFound {
		return file, err
	}

	hc, err := headCommit(repo)
	if err != nil && err != ErrMetadataNotFound {
		return file, err
	}

	ri := RepositoryInfo{LatestRevision: hc.Id().String()}
	if err != nil {
		return nil, err
	}

	translations, err := getTranslations(repo, directory, document, filename)

	file = &File{
		Filename:       filename,
		Document:       document,
		Path:           directory,
		HTML:           html,
		Markdown:       markdown,
		FrontMatter:    fm,
		DirectoryInfo:  di,
		RepositoryInfo: &ri,
		Translations:   translations,
	}

	return file, nil
}

func getTranslations(repo *git.Repository, directory, document, filename string) (langs []string, err error) {

	langs = []string{}
	tree, err := headTree(repo)

	if !config.TranslationEnabled {
		return langs, fmt.Errorf("translation is not enabled")
	}

	for _, lc := range config.EnabledLanguages {

		target := filepath.Join(directory, document, translationFilename(filename, lc))

		Debug.Println("checking for translation", target)

		// ignore not found error and add to output if present
		entry, _ := tree.EntryByPath(target)

		if entry != nil {
			Debug.Println("found translation", target)
			langs = append(langs, lc)
		}
	}

	if len(langs) == 0 {
		return langs, fmt.Errorf("No translations found")
	}

	return langs, err
}

func getAttachments(directory string) (files []Attachment, err error) {

	// Initialise the slice so [] is marshalled instead of null
	files = make([]Attachment, 0)

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	ht, err := headTree(repo)
	if err != nil {
		return nil, err
	}

	// ensure that the directory exists
	entry, _ := ht.EntryByPath(directory)
	if entry == nil {
		return nil, fmt.Errorf("directory '%s' not found", directory)
	}

	if entry.Type != git.ObjectTree {
		return nil, fmt.Errorf("%s is not a directory", directory)
	}

	tree, err := repo.LookupTree(entry.Id)
	if err != nil {
		return nil, fmt.Errorf("couldn't find tree for entry %s", entry.Id)
	}

	defer tree.Free()

	walkIterator := func(path string, te *git.TreeEntry) int {
		var blob *git.Blob
		var ext string
		var attachment Attachment

		if te.Type == git.ObjectBlob {

			ext = filepath.Ext(te.Name)

			// skip if a Markdown file
			if ext == ".md" {
				Debug.Println("markdown file, skipping:", te.Name)
				return 0
			}

			blob, err = repo.LookupBlob(te.Id)

			if err != nil {
				Warning.Println("Failed to find blob", te.Id)
				return -1
			}

			data := blob.Contents()

			attachment = Attachment{
				Filename:  te.Name,
				Extension: ext,
				Data:      base64.StdEncoding.EncodeToString(data),
				Path:      filepath.Join(directory, path),
				MediaType: getMediaType(ext),
			}

			files = append(files, attachment)

		}

		return 0
	}

	err = tree.Walk(walkIterator)

	return files, err

}

func getConvertedFile(directory, document, filename string) (file *File, err error) {
	Debug.Println("Getting converted file", directory, filename)
	file, err = getFile(directory, document, filename, false, true)
	if err != nil {
		return nil, err
	}
	return
}

func getRawFile(directory, document, filename string) (file *File, err error) {
	Debug.Println("Getting raw file", directory, filename)
	file, err = getFile(directory, document, filename, true, false)
	if err != nil {
		return nil, err
	}
	return
}

func countFiles() (counts map[string]int, err error) {
	repo, err := repository(config)
	if err != nil {
		return counts, err
	}

	counts, err = getFileCounts(repo)
	return counts, err
}

func getFileCounts(repo *git.Repository) (counter map[string]int, err error) {

	ht, err := headTree(repo)
	if err != nil {
		return counter, err
	}

	defer ht.Free()

	counter = make(map[string]int)
	fileCategoryLookup := make(map[string]string)

	// initialise a 'catch all' counter called other at 0
	counter["other"] = 0

	// loop through file categories and build
	// counter and lookup map
	for category, filetypes := range config.FileCategories {
		counter[category] = 0

		for _, filetype := range filetypes {
			fileCategoryLookup[filetype] = category
		}
	}

	walkIterator := func(_ string, te *git.TreeEntry) int {

		var ext, fc string

		if te.Type == git.ObjectBlob {
			ext = filepath.Ext(te.Name)
			fc = fileCategoryLookup[ext]

			if fc != "" {
				counter[fc]++
			} else {
				counter["other"]++
			}
		}

		return 0
	}

	err = ht.Walk(walkIterator)

	return counter, err
}

func writeTreeAndCommit(repo *git.Repository, index *git.Index, message string, user User) (oid *git.Oid, err error) {

	// write the tree, persisting our addition to the git repo
	treeID, err := index.WriteTree()
	if err != nil {
		return oid, err
	}

	// and use the tree's id to find the actual updated tree
	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return oid, err
	}

	// find the repository's tip, where we're committing to
	tip, err := headCommit(repo)
	if err != nil {
		return oid, err
	}

	// git signatures
	author := sign(user)
	committer := sign(user)

	// now commit our updated tree to the tip (parent)
	oid, err = repo.CreateCommit("HEAD", author, committer, message, tree, tip)
	if err != nil {
		return oid, err
	}

	// checkout to keep file system in sync with git
	err = repo.CheckoutHead(
		&git.CheckoutOpts{Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutForce},
	)

	return oid, err

}

func pathInFiles(directory, document, filename string, files *[]NewCommitFile) bool {

	// check that at least one file in files matches the directory and filename
	for _, file := range *files {
		if file.Path == directory &&
			file.Document == document &&
			file.Filename == filename {
			return true
		}
	}

	return false
}

func pathInDirectories(directory string, directories *[]NewCommitDirectory) bool {

	// check that at least one directory files matches the path's directory
	for _, d := range *directories {
		if d.Path == directory {
			return true
		}
	}

	return false
}

// getMediaType returns the correct media when passed a file extension
//
//https://en.wikipedia.org/wiki/Data_URI_scheme#Syntax
func getMediaType(extension string) string {

	var extensionNoDot string

	if len(config.MediaTypes) == 0 {
		Error.Println("No media types configured")
		return "none"
	}

	extensionNoDot = strings.Replace(extension, ".", "", 1)

	mt := config.MediaTypes[extensionNoDot]

	if mt == "" {
		fallback := fmt.Sprintf("unknown/%s", extensionNoDot)
		Warning.Printf("No media type found for '%s', returning '%s'", extension, fallback)
		return fallback
	}

	return mt
}

// extractContents retrieves the contents of the NewCommitFile and prepares
// them to be written to the repository.
//
// * Markdown files are combined with the FrontMatter
// * Base64 encoded files are decoded to a byte sequence
// * Plain text files are left untouched, simply converted to a byte slice
func extractContents(ncf NewCommitFile) (contents []byte, err error) {

	if ncf.Base64Encoded {

		contents, err := base64.StdEncoding.DecodeString(ncf.Body)
		if err != nil {
			Error.Println("Failed to decode file:", ncf)
			return contents, fmt.Errorf("Failed to decode file: %s", ncf.Filename)
		}
		return contents, err

	} else if filepath.Ext(ncf.Filename) == ".md" {

		return ncf.ToMarkdown(), err

	} else {
		return []byte(ncf.Body), err
	}

}

// A quicker, more-efficient way of extracting the frontmatter from
// a markdown file, this only reads the frontmatter from the top and
// skips the markdown beneath.
// This exists because particle slows down by reading the entire thing
func getMetadataFromBlob(blob *git.Blob) (fm FrontMatter, err error) {

	const fmBoundary = "---"

	var reader io.Reader
	var scanner *bufio.Scanner
	var fmText *bytes.Buffer
	var fmBoundaryCount = 0
	var textPresent bool

	reader = bytes.NewReader(blob.Contents())

	fmText = bytes.NewBuffer(nil)
	scanner = bufio.NewScanner(reader)

	// Read the blob contents line by line and write the
	// frontmatter to fmText

Scan:
	for scanner.Scan() {

		// When we hit the first yaml marker, '---'
		// increment the fmBoundaryCount, and while it's 1
		// read any lines into the fmText buffer
		if scanner.Text() == fmBoundary {
			fmBoundaryCount++
		}

		switch fmBoundaryCount {

		// look for 'normal' text before we've encountered any
		// YAML fm boundaries
		case 0:
			textPresent, err = regexp.MatchString("[A-z0-9]", scanner.Text())
			if err != nil {
				return fm, err
			}

			// and if we find any, exit.
			if textPresent {
				break Scan
			}

			continue Scan

		// once we hit the first boundary, write the lines to
		// the fmText buffer
		case 1:

			fmText.WriteString(scanner.Text())
			fmText.WriteString("\n")
			continue Scan

		// when we reach the second boundary, we're done
		case 2:
			break Scan

		default:
			Warning.Println("fmBoundaryCount exceeded two", fmBoundaryCount)
		}
	}

	err = yaml.Unmarshal(fmText.Bytes(), &fm)

	return fm, err
}

func getMetadataFromDirectory(directory string) (*DirectoryInfo, error) {

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	ht, err := headTree(repo)
	if err != nil {
		return nil, err
	}

	entry := ht.EntryByName(directory)

	tree, err := repo.LookupTree(entry.Id)
	if err != nil {
		return nil, err
	}

	md, err := getMetadata(repo, tree)
	if err == ErrMetadataNotFound {
		return nil, err
	}

	return &md, err
}

func getMetadata(repo *git.Repository, tree *git.Tree) (di DirectoryInfo, err error) {
	var reader io.Reader

	infoEntry, err := tree.EntryByPath("_index.md")
	if err != nil {
		Warning.Println("_index.md does not exist in the repository, skipping for", tree.Object.Id())
		return di, ErrMetadataNotFound
	}

	blob, err := repo.LookupBlob(infoEntry.Id)
	if err != nil {
		Warning.Println("_index.md cannot be retrieved, exiting", infoEntry.Id)
		return di, err
	}
	defer blob.Free()

	reader = bytes.NewReader(blob.Contents())

	_, err = particle.YAMLEncoding.DecodeReader(reader, &di)

	if err != nil {
		Warning.Println("_index.md cannot be decoded, exiting", blob.Contents())
		return di, err
	}

	return di, err
}

func translationFilename(fn, code string) (tfn string) {
	const ext = "md" // assuming we'll always be using .md for markdown
	const delim = "."
	var base string

	base = strings.Split(fn, delim)[0]

	if code == config.DefaultLanguage {
		return strings.Join([]string{base, ext}, delim)
	}

	return strings.Join([]string{base, code, ext}, delim)

}

func fileExists(repo *git.Repository, path, document, filename string) (exists bool, err error) {

	tree, err := headTree(repo)
	if err != nil {
		return false, err
	}

	_, err = tree.EntryByPath(filepath.Join(path, document, filename))
	if err != nil {
		return false, err
	}

	return true, err
}
