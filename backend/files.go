package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/graphia/particle"
	"gopkg.in/libgit2/git2go.v25"
)

// getFilesInDir returns a list of FileItems for listing
func getFilesInDir(directory string) (files []FileItem, err error) {

	var pathFound = false

	repo, err := repository(config)
	if err != nil {
		return
	}
	defer repo.Free()

	tree, err := headTree(repo)
	if err != nil {
		return
	}

	walkIterator := func(td string, te *git.TreeEntry) int {

		td = strings.TrimRight(td, "/")

		// we've found the directory we're looking for
		if td == directory {
			Debug.Println("found directory", directory)
			pathFound = true
		}

		if te.Type == git.ObjectBlob && td == directory {

			fi := FileItem{
				AbsoluteFilename: fmt.Sprintf("%s%s", td, te.Name),
				Filename:         te.Name,
				Path:             td,
				Author:           "Joey Joe",
			}

			Debug.Println("found file", fi)

			files = append(files, fi)

		}

		return 0
	}

	err = tree.Walk(walkIterator)

	if !pathFound {
		return nil, fmt.Errorf("directory '%s' not found", directory)
	}

	return files, err
}

func createFile(rw RepoWrite) (oid *git.Oid, err error) {

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	tree, err := headTree(repo)
	if err != nil {
		return nil, err
	}

	target := filepath.Join(rw.Path, rw.Filename)
	Debug.Println("looking for file", target)

	file, _ := tree.EntryByPath(target)
	if file != nil {
		return nil, fmt.Errorf("file already exists %s", target)
	}

	oid, err = writeFile(repo, rw)
	return oid, err

}

func createEmptyFile(rw RepoWrite) (oid *git.Oid, err error) {

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	tree, err := headTree(repo)
	if err != nil {
		return nil, err
	}

	target := filepath.Join(rw.Path, rw.Filename)
	Debug.Println("looking for file", target)

	file, _ := tree.EntryByPath(target)
	if file != nil {
		return nil, fmt.Errorf("file already exists %s", target)
	}

	index, err := repo.Index()
	if err != nil {
		return nil, err
	}
	defer index.Free()

	// add an empty blob to the repo
	oid, err = repo.CreateBlobFromBuffer([]byte{})
	if err != nil {
		return nil, err
	}

	// build the git index entry and add it to the index
	ie := buildIndexEntry(oid, rw)

	Debug.Println("IndexEntry", ie)
	err = index.Add(&ie)
	if err != nil {
		return nil, err
	}

	// write the tree, persisting our addition to the git repo
	treeID, err := index.WriteTree()
	if err != nil {
		return nil, err
	}

	// and use the tree's id to find the actual updated tree
	tree, err = repo.LookupTree(treeID)
	if err != nil {
		return nil, err
	}

	// find the repository's tip, where we're committing to
	tip, err := headCommit(repo)
	if err != nil {
		return nil, err
	}

	// git signatures
	author := sign(rw)
	committer := sign(rw)

	// now commit our updated tree to the tip (parent)
	oid, err = repo.CreateCommit("HEAD", author, committer, rw.Message, tree, tip)
	if err != nil {
		return nil, err
	}

	// checkout to keep file system in sync with git
	err = repo.CheckoutHead(
		&git.CheckoutOpts{Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutForce},
	)

	if err != nil {
		Error.Println("Could not checkout head:", err.Error())
	}

	return oid, err

}

// TODO how do we know that a commit hasn't been made in
// between us serving and receiving the update?
func updateFile(rw RepoWrite) (oid *git.Oid, err error) {

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	tree, err := headTree(repo)
	if err != nil {
		return nil, err
	}

	target := filepath.Join(rw.Path, rw.Filename)

	file, _ := tree.EntryByPath(target)
	if file == nil {
		return nil, fmt.Errorf("file does not exist %s", target)
	}

	oid, err = writeFile(repo, rw)
	if err != nil {
		return nil, fmt.Errorf("file modification failed %s", err.Error())
	}

	return oid, err

}

func writeFile(repo *git.Repository, rw RepoWrite) (oid *git.Oid, err error) {

	index, err := repo.Index()
	if err != nil {
		return nil, err
	}
	defer index.Free()

	// add frontmatter to the file contents, followed by the document body
	contents := particle.YAMLEncoding.EncodeToString([]byte(rw.Body), &rw.FrontMatter)

	// and add the combined contents to a blob in the repo
	oid, err = repo.CreateBlobFromBuffer([]byte(contents))
	if err != nil {
		return nil, err
	}

	// build the git index entry and add it to the index
	ie := buildIndexEntry(oid, rw)

	Debug.Println("IndexEntry", ie)
	err = index.Add(&ie)
	if err != nil {
		return nil, err
	}

	// write the tree, persisting our addition to the git repo
	treeID, err := index.WriteTree()
	if err != nil {
		return nil, err
	}

	// and use the tree's id to find the actual updated tree
	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return nil, err
	}

	// find the repository's tip, where we're committing to
	tip, err := headCommit(repo)
	if err != nil {
		return nil, err
	}

	// git signatures
	author := sign(rw)
	committer := sign(rw)

	// now commit our updated tree to the tip (parent)
	oid, err = repo.CreateCommit("HEAD", author, committer, rw.Message, tree, tip)
	if err != nil {
		return nil, err
	}

	// checkout to keep file system in sync with git
	err = repo.CheckoutHead(
		&git.CheckoutOpts{Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutForce},
	)

	if err != nil {
		Error.Println("Could not checkout head:", err.Error())
	}

	return oid, err

}

func createDirectory(rw RepoWrite) (oid *git.Oid, err error) {

	// check that the directory does not already exist
	// FIXME this could be improved by performing the check
	// against the git repository rather than just checking
	// the filesystem
	_, err = os.Stat(filepath.Join(config.Repository, rw.Path))
	if err == nil {
		return nil, fmt.Errorf("directory already exists")
	}

	// git can't track empty directories, so, Rails-style, we'll add an
	// empty file called .keep
	rw.Filename = ".keep"
	rw.Body = ""

	// And set the commit message sensibly so we don't need to prompt
	// the user
	rw.Message = fmt.Sprintf("Added %s directory", rw.Path)

	oid, err = createEmptyFile(rw)
	if err != nil {
		return nil, err
	}

	return oid, err
}

func deleteDirectory(rw RepoWrite) (oid *git.Oid, err error) {

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	target := rw.Path

	ht, err := headTree(repo)
	if err != nil {
		return nil, err
	}

	// ensure that the directory exists before we try to delete it
	directory, _ := ht.EntryByPath(target)
	if directory == nil {
		return nil, fmt.Errorf("directory does not exist %s", target)
	}

	// now go head and remove it

	// first grab the repo's index
	index, err := repo.Index()
	if err != nil {
		return nil, err
	}

	Debug.Println("Removing directory:", target)

	// and remove the target by path and everything beneath it
	err = index.RemoveDirectory(target, 0)
	if err != nil {
		return nil, err
	}

	// write the tree, persisting our deletion to the git repo
	treeID, err := index.WriteTree()
	if err != nil {
		return nil, err
	}

	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return nil, err
	}

	// find the repository's tip, where we're committing to
	tip, err := headCommit(repo)
	if err != nil {
		return nil, err
	}

	// git signatures
	author := sign(rw)
	committer := sign(rw)

	// now commit our updated tree to the tip (parent)
	oid, err = repo.CreateCommit("HEAD", author, committer, rw.Message, tree, tip)
	if err != nil {
		return nil, err
	}

	// checkout to keep file system in sync with git
	err = repo.CheckoutHead(
		&git.CheckoutOpts{Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutForce},
	)

	if err != nil {
		Error.Println("Could not checkout head:", err.Error())
	}

	return oid, err

}

func deleteFile(rw RepoWrite) (oid *git.Oid, err error) {

	repo, err := repository(config)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	target := filepath.Join(rw.Path, rw.Filename)

	ht, err := headTree(repo)
	if err != nil {
		return nil, err
	}

	// ensure that the file exists before we try to delete it
	file, _ := ht.EntryByPath(target)
	if file == nil {
		return nil, fmt.Errorf("file does not exist %s", target)
	}

	// now go head and remove it

	// first grab the repo's index
	index, err := repo.Index()
	if err != nil {
		return nil, err
	}

	// and remove the target by path
	err = index.RemoveByPath(target)
	if err != nil {
		return nil, err
	}

	// write the tree, persisting our deletion to the git repo
	treeID, err := index.WriteTree()
	if err != nil {
		return nil, err
	}

	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return nil, err
	}

	// find the repository's tip, where we're committing to
	tip, err := headCommit(repo)
	if err != nil {
		return nil, err
	}

	// git signatures
	author := sign(rw)
	committer := sign(rw)

	// now commit our updated tree to the tip (parent)
	oid, err = repo.CreateCommit("HEAD", author, committer, rw.Message, tree, tip)
	if err != nil {
		return nil, err
	}

	// checkout to keep file system in sync with git
	err = repo.CheckoutHead(
		&git.CheckoutOpts{Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing | git.CheckoutForce},
	)

	if err != nil {
		Error.Println("Could not checkout head:", err.Error())
	}

	return oid, err

}

func sign(rw RepoWrite) *git.Signature {
	return &git.Signature{
		Name:  rw.Name,
		Email: rw.Email,
		When:  time.Now(),
	}
}

func buildIndexEntry(oid *git.Oid, rw RepoWrite) git.IndexEntry {
	return git.IndexEntry{
		Id:   oid,
		Path: filepath.Join(rw.Path, rw.Filename),
		Size: uint32(len(rw.Body)),

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
		Filename: filename,
		Path:     directory,
		HTML:     html,
		Markdown: markdown,

		// front matter derived attributes
		Title:  fm.Title,
		Author: fm.Author,
	}

	return file, err
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
