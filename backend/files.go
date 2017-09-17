package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/graphia/particle"
	"gopkg.in/libgit2/git2go.v25"
)

var (
	// ErrMetaDataNotFound is returned when the _index.md file is missing
	// this isn't a catestrophic problem, but should be logged
	ErrMetaDataNotFound = errors.New("_index.md not found")
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

	walkIterator := func(_ string, te *git.TreeEntry) int {
		var fm FrontMatter
		var blob *git.Blob
		var reader io.Reader
		var ext string

		if te.Type == git.ObjectBlob {

			// skip unless it's a Markdown file

			ext = filepath.Ext(te.Name)

			if ext != ".md" {
				Warning.Println("not a markdown file, skipping:", te.Name)
				return 0
			}

			if te.Name == "_index.md" {
				Warning.Println("is a metadata file, skipping:", te.Name)
				return 0
			}

			blob, err = repo.LookupBlob(te.Id)

			if err != nil {
				Warning.Println("Failed to find blob", te.Id)
				return -1
			}

			reader = bytes.NewReader(blob.Contents())

			_, err := particle.YAMLEncoding.DecodeReader(reader, &fm)

			if err != nil {
				Warning.Println("Failed to decode file", string(blob.Contents()))
				Warning.Println("Frontmatter:", fm)
				return -1
			}

			fi := FileItem{
				AbsoluteFilename: fmt.Sprintf("%s/%s", directory, te.Name),
				Filename:         te.Name,
				Path:             directory,
				FrontMatter:      fm,
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

		// if no filename is specified it's likely that
		// we are simply creating a directory
		// TODO fix this
		if ncf.Filename == "" {
			continue
		}

		target := filepath.Join(ncf.Path, ncf.Filename)

		entry, _ := ht.EntryByPath(target)
		//if err != nil {
		//return nil, fmt.Errorf("file not found %s", target)
		//}

		if entry != nil {
			return nil, fmt.Errorf("file already exists")
		}

	}

	oid, err = writeFiles(repo, nc, user)

	return oid, err
}

// Replaces updateFile
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
		return nil, err
	}
	defer index.Free()

	var contents []byte

	for _, ncf := range nc.Files {

		var ie git.IndexEntry

		// get the file contents in the correct format
		contents, err = extractContents(ncf)
		if err != nil {
			return nil, err
		}

		oid, err = repo.CreateBlobFromBuffer(contents)
		if err != nil {
			return nil, err
		}

		// build the git index entry and add it to the index
		ie = buildIndexEntry(oid, ncf)

		err = index.Add(&ie)
		if err != nil {
			return nil, err
		}

	}

	oid, err = writeTreeAndCommit(repo, index, nc, user)
	if err != nil {
		return oid, err
	}

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

	oid, err = writeTreeAndCommit(repo, index, nc, user)
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

			// if there is any kind of error except ErrMetaDataNotFound,
			// something's wrong, quit
			if err != ErrMetaDataNotFound && err != nil {
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

	oid, err = writeTreeAndCommit(repo, index, nc, user)

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

	oid, err = writeTreeAndCommit(repo, index, nc, user)
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

	ht, err := headTree(repo)
	if err != nil {
		return oid, err
	}

	// first grab the repo's index
	index, err := repo.Index()
	if err != nil {
		return oid, err
	}

	for _, ncf := range nc.Files {

		target := filepath.Join(ncf.Path, ncf.Filename)

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

	oid, err = writeTreeAndCommit(repo, index, nc, user)
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
		Path: filepath.Join(ncf.Path, ncf.Filename),
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

func getFile(directory string, filename string, includeMd, includeHTML bool) (file *File, err error) {
	var html, markdown *string = nil, nil
	var fm FrontMatter

	repo, err := repository(config)

	tree, err := headTree(repo)
	if err != nil {
		return nil, err
	}
	target := filepath.Join(directory, filename)

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

	file = &File{
		Filename:    filename,
		Path:        directory,
		HTML:        html,
		Markdown:    markdown,
		FrontMatter: fm,
	}

	return file, err
}

func getAttachments(directory string) (files []Attachment, err error) {

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
				Filename:         te.Name,
				AbsoluteFilename: filepath.Join(directory, path, te.Name),
				Extension:        ext,
				Data:             base64.StdEncoding.EncodeToString(data),
				Path:             filepath.Join(directory, path),
				MediaType:        getMediaType(ext),
			}

			files = append(files, attachment)

		}

		return 0
	}

	err = tree.Walk(walkIterator)

	return files, err

}

func getConvertedFile(directory, filename string) (file *File, err error) {
	Debug.Println("Getting converted file", directory, filename)
	file, err = getFile(directory, filename, false, true)
	if err != nil {
		return nil, err
	}
	return
}

func getRawFile(directory, filename string) (file *File, err error) {
	Debug.Println("Getting raw file", directory, filename)
	file, err = getFile(directory, filename, true, false)
	if err != nil {
		return nil, err
	}
	return
}

func countFiles() (counter map[string]int, err error) {
	repo, err := repository(config)
	if err != nil {
		return counter, err
	}

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

func writeTreeAndCommit(repo *git.Repository, index *git.Index, nc NewCommit, user User) (oid *git.Oid, err error) {

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
	oid, err = repo.CreateCommit("HEAD", author, committer, nc.Message, tree, tip)
	if err != nil {
		return oid, err
	}

	// checkout to keep file system in sync with git
	err = repo.CheckoutHead(
		&git.CheckoutOpts{Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutForce},
	)

	return oid, err

}

func pathInFiles(directory, filename string, files *[]NewCommitFile) bool {

	// check that at least one file in files matches the directory and filename
	for _, file := range *files {
		if file.Path == directory && file.Filename == filename {
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

		return []byte(
			particle.YAMLEncoding.EncodeToString(
				[]byte(ncf.Body), &ncf.FrontMatter,
			),
		), err

	} else {
		return []byte(ncf.Body), err
	}

}

func getMetadataFromDirectory(directory string) (di DirectoryInfo, err error) {

	repo, err := repository(config)
	if err != nil {
		return di, err
	}
	defer repo.Free()

	ht, err := headTree(repo)
	if err != nil {
		return di, err
	}

	entry := ht.EntryByName(directory)

	tree, err := repo.LookupTree(entry.Id)
	if err != nil {
		return di, err
	}

	di, err = getMetadata(repo, tree)

	return di, err
}

func getMetadata(repo *git.Repository, tree *git.Tree) (di DirectoryInfo, err error) {
	var reader io.Reader

	infoEntry, err := tree.EntryByPath("_index.md")
	if err != nil {
		Warning.Println("_index.md does not exist in the repository, skipping for", tree.Object.Id())
		return di, ErrMetaDataNotFound
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
