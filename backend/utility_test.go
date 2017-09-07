package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"github.com/asdine/storm"
	"github.com/dgrijalva/jwt-go"
	"github.com/urfave/negroni"

	"gopkg.in/libgit2/git2go.v25"
)

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func setupSmallTestRepo(dest string) (oid *git.Oid, err error) {
	src := "../tests/backend/repositories/small"
	oid, err = setupTestRepo(src, dest)
	return
}

func setupSubdirsTestRepo(dest string) (oid *git.Oid, err error) {
	src := "../tests/backend/repositories/subdirs"
	oid, err = setupTestRepo(src, dest)
	return
}

func setupMultipleFiletypesTestRepo(dest string) (oid *git.Oid, err error) {
	src := "../tests/backend/repositories/multiple_filetypes"
	oid, err = setupTestRepo(src, dest)
	return
}

func setupTestRepo(src, dest string) (oid *git.Oid, err error) {

	// copy the small repo skeleton to specified path
	err = os.RemoveAll(dest)
	if err != nil {
		return nil, err
	}

	err = CopyDir(src, dest)
	if err != nil {
		return nil, err
	}

	// now initialise the git repository with a blanket commit
	sig := &git.Signature{
		Name:  "Joe Quimby, Jr",
		Email: "diamond.joe@springfield.nt.gov",
		When:  time.Now(),
	}

	repo, err := git.InitRepository(dest, false)
	if err != nil {
		panic(err)
	}

	idx, err := repo.Index()
	if err != nil {
		panic(err)
	}

	err = idx.AddAll(
		[]string{
			filepath.Join("appendices", "*.md"),
			filepath.Join("documents", "*.md"),
			filepath.Join("documents", "*.json"),
			filepath.Join("appendices", "*.json"),
			filepath.Join("documents", "document_1", "*.*"),
			filepath.Join("appendices", "appendix_1", "*.*"),
		},
		git.IndexAddForce,
		nil,
	)
	if err != nil {
		panic(err)
	}

	err = idx.Write()
	if err != nil {
		panic(err)
	}

	treeID, err := idx.WriteTree()
	if err != nil {
		panic(err)
	}

	message := "Quick, honk at that broad!\n"
	tree, err := repo.LookupTree(treeID)
	if err != nil {
		panic(err)
	}

	oid, err = repo.CreateCommit("HEAD", sig, sig, message, tree)
	if err != nil {
		panic(err)
	}

	// initialise a config obj so createFile looks in the right place

	testConfigPath := "../config/test.yml"
	config, err = loadConfig(&testConfigPath)
	config.Repository = filepath.Join(dest)

	return oid, err
}

func flushDB(path string) storm.DB {

	os.RemoveAll(path)

	stormDB, err := storm.Open(path)
	if err != nil {
		panic(fmt.Sprintf("Database cannot be openend %s", err.Error()))
	}
	return *stormDB
}

func setupTestKeys() {
	signBytes, err := ioutil.ReadFile(testPrivKeyPath)
	if err != nil {
		panic(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}

	verifyBytes, err := ioutil.ReadFile(testPubKeyPath)
	if err != nil {
		panic(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}
}

func createTestServerWithContext() (server *httptest.Server) {
	n := negroni.New()

	r := unprotectedRouter()
	pr := protectedRouter()

	r.Handle("/api/*", negroni.New(
		negroni.HandlerFunc(apiTestMiddleware),
		negroni.Wrap(pr),
	))

	n.UseHandler(r)

	server = httptest.NewServer(n)

	return server

}

func apiTestUser() (user User) {
	return User{Name: "Selma Bouvier", Email: "selma.b@aol.com"}
}

func apiTestMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := apiTestUser()
	ctx := newContextWithCurrentUser(r.Context(), r, user)
	next(w, r.WithContext(ctx))
}
