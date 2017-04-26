package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/husobee/vestigo"
)

// Authentication functionality 🔑

func authLoginHandler(w http.ResponseWriter, r *http.Request) {
	// start session, generate a JWT if the credentials are ok

	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	if strings.ToLower(user.Username) != "someone" {
		if user.Password != "p@ssword" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}
	}

	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		panic(err)
	}

	tokenString, err := token.SignedString(signKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		panic(err)
	}

	response := Token{tokenString}
	JSONResponse(response, w)

}

func authLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// end session
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{"Gained access to protected resource"}
	JSONResponse(response, w)
}

// cmsGeneralHandler takes care of serving the frontend CMS portion
// of the app. Paths that match '/cms' are routed here; if the file exists
// on the filesystem it's served. If not, cms/index.html is served instead
//
// GET /cms                        -> public/cms/index.html
// GET /cms/javascripts/app.js     -> public/cms/javascripts/app.js
// GET /cms/something/nonexistant  -> public/cms/index.html
func cmsGeneralHandler(w http.ResponseWriter, r *http.Request) {
	Debug.Println(r.URL.Path)

	// TODO make this path configurable
	base := "./frontend/public"

	path := filepath.Join(base, r.URL.Path)

	// if we can't find a file (asset) to serve but the path begins
	// with /cms, serve the CMS's index. It is likely someone navigating
	// directly to a resource or refreshing a page. If Vue's router
	// *still* can't match it to a page, it'll show an appropriate error
	index := filepath.Join(base, "cms", "index.html")

	if _, err := os.Stat(path); err == nil {
		http.ServeFile(w, r, path)
	} else {
		http.ServeFile(w, r, index)
	}
}

func apiRootHandler(w http.ResponseWriter, r *http.Request) {}

// Directory level functionality 🗃

// apiListDirectoriesHandler returns a JSON object representing
// the Repository's contents.
//
// GET /api/directories/
//
// eg. when the documents directory contains Documents 1 and 2:
//
// {
//   documents: [
//	   {Filename: 'Document 1', ...},
//	   {Filename: 'Document 2', ...}
//	 ],
//	 appendices: [...]
// }
func apiListDirectoriesHandler(w http.ResponseWriter, r *http.Request) {

	fi1 := FileItem{Filename: "abc123"}
	fi2 := FileItem{Filename: "def234"}

	filesByDirectory := map[string][]FileItem{
		"documents": []FileItem{fi1, fi2},
	}

	output, err := json.Marshal(filesByDirectory)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)

	w.WriteHeader(200)
}

// apiCreateDirectoryHandler creates an 'empty' directory. It actually
// contains a hidden `.keep` file, so it's trackable by git
//
// POST /api/directories
// {
//	 "path": "recipes",
//	 "name": "Martin Prince",
//	 "email": "mp@springfield-elementary.gov",
//	 "message": "Added new directory called Bobbins"
// }
func apiCreateDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var rw RepoWrite
	var sr SuccessResponse

	rw = RepoWrite{}
	json.NewDecoder(r.Body).Decode(&rw)

	oid, err := createDirectory(rw)

	if err != nil {
		fr := FailureResponse{Message: err.Error()}
		output, err := json.Marshal(fr)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(400)
		w.Write(output)
	}

	sr = SuccessResponse{
		Message: "Directory created",
		Oid:     oid.String(),
	}

	output, err := json.Marshal(sr)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	w.Write(output)

}

func apiRenameDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(200)
}

// apiDeleteDirectoryHandler deletes a directory and all of its contents.
// Note, this is rather dangerous and, on the front end, should be guarded from
// accidental clicking 😾
//
// DELETE /api/directories/:directory
// {
//	 "path": "recipes",
//	 "name": "Martin Prince",
//	 "email": "mp@springfield-elementary.gov",
// }
//
// returns a SuccessResponse containing the git commit hash or a FailureResponse
// containing an error message
func apiDeleteDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var rw RepoWrite
	var fr FailureResponse
	var sr SuccessResponse

	// set up the RepoWrite with git params, an appropraiate message and then
	// specify the directory based on the path
	rw = RepoWrite{}
	json.NewDecoder(r.Body).Decode(&rw)
	rw.Message = fmt.Sprintf("Deleted %s directory", rw.Path)
	rw.Path = vestigo.Param(r, "directory")

	oid, err := deleteDirectory(rw)

	if err != nil {

		msg := "Failed to delete directory"

		Error.Println(msg, err.Error())

		fr = FailureResponse{
			Message: fmt.Sprintf("%s: %s", msg, err.Error()),
		}

		output, err := json.Marshal(fr)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(400)
		w.Write(output)

		return
	}

	sr = SuccessResponse{
		Message: "Directory deleted",
		Oid:     oid.String(),
	}

	output, err := json.Marshal(sr)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	w.Write(output)
}

// Inside a directory functionality 🗂

// apiListFilesInDirectoryHandler returns a JSON array containing
// the all files belonging to :directory
//
// GET /api/directories/:directory/files
//
// eg. when the documents directory contains Documents 1 and 2:
//
// [
//	 {"filename": "Document 1", ...},
//   {"filename": "Document 2", ...}
// ]
func apiListFilesInDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	directory := vestigo.Param(r, "directory")
	files, err := getFilesInDir(directory)
	if err != nil {
		Error.Println("Could not get list of files in directory", directory, err.Error())
	}

	output, err := json.Marshal(files)
	if err != nil {
		Error.Println("Could not create JSON", directory, err.Error())
		Error.Println("Files", files)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(200)
	w.Write(output)
}

// apiCreateFileInDirectory creates a file the specified directory
//
// POST /api/directories/:directory/files
// {
//   "filename": "document_6.md",
//   "body": "The *quick* brown fox [jumped](appendices/jumping.md) ...",
//	 "message": "Added document six"
//   ...
// }
//
// returns a SuccessResponse containing the git commit hash or a FailureResponse
// containing an error message
func apiCreateFileInDirectoryHandler(w http.ResponseWriter, r *http.Request) {

	var rw RepoWrite
	var fr FailureResponse
	var sr SuccessResponse

	directory := vestigo.Param(r, "directory")

	rw = RepoWrite{Path: directory}
	json.NewDecoder(r.Body).Decode(&rw)

	oid, err := createFile(rw)
	if err != nil {
		fr = FailureResponse{
			Message: err.Error(),
		}
		output, err := json.Marshal(fr)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(400)
		w.Write(output)
	}

	sr = SuccessResponse{
		Message: "File created",
		Oid:     oid.String(),
	}

	output, err := json.Marshal(sr)
	if err != nil {
		panic(err)
	}

	w.Write(output)

	Debug.Println("File created", oid)

}

// apiUpdateFileInDirectory updates an existing file the specified
// directory
//
// PATCH /api/directories/:directory/files/:filename
// {
//   "body": "The *quick* brown dog [jumped](appendices/jumping.md) ...",
//	 "message": "Fixed grammar in document six"
//   ...
// }
//
// returns a SuccessResponse containing the git commit hash or a
// FailureResponse containing an error message
func apiUpdateFileInDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var rw RepoWrite
	var fr FailureResponse
	var sr SuccessResponse

	directory := vestigo.Param(r, "directory")
	filename := vestigo.Param(r, "file")

	rw = RepoWrite{Path: directory, Filename: filename}
	json.NewDecoder(r.Body).Decode(&rw)

	oid, err := updateFile(rw)
	if err != nil {
		Error.Println("Failed to update file", err.Error())
		fr = FailureResponse{
			Message: err.Error(),
		}
		output, err := json.Marshal(fr)
		if err != nil {
			Error.Println(err.Error())
		}

		w.WriteHeader(400)
		w.Write(output)
	}

	sr = SuccessResponse{
		Message: "File updated",
		Oid:     oid.String(),
	}

	output, err := json.Marshal(sr)
	if err != nil {
		Error.Println(err.Error())
	}

	w.Write(output)

	Debug.Println("File updated", oid)

}

// apiDeleteFileFromDirectoryHandler deletes a file from the specified
// directory
//
// DELETE /api/directories/:directory/files/:filename
// {
//   "message": "Deleted document 6 as it's no longer required"
// }
//
// returns a SuccessResponse containing the git commit hash or a
// FailureResponse containing an error message
func apiDeleteFileFromDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var rw RepoWrite
	var fr FailureResponse
	var sr SuccessResponse

	directory := vestigo.Param(r, "directory")
	filename := vestigo.Param(r, "file")

	rw = RepoWrite{Path: directory, Filename: filename}
	json.NewDecoder(r.Body).Decode(&rw)

	oid, err := deleteFile(rw)

	if err != nil {
		Error.Println("Failed to delete file", err.Error())
		fr = FailureResponse{
			Message: err.Error(),
		}
		output, err := json.Marshal(fr)
		if err != nil {
			Error.Println(err.Error())
		}

		w.WriteHeader(400)
		w.Write(output)
	}

	sr = SuccessResponse{
		Message: "File deleted",
		Oid:     oid.String(),
	}

	output, err := json.Marshal(sr)
	if err != nil {
		Error.Println(err.Error())
	}

	w.Write(output)

	Debug.Println("File updated", oid)

}

// apiGetFileInDirectoryHandler returns a File object representing the
// specified file to be displayed in a web page; the raw Markdown is not
// sent but the compiled HTML is.
//
// GET /api/directories/:directory/files/:filename
//
// returns
//
// {
//   "filename": "document_3.md",
//   "path": "documents",
//   "author": "Carl Carlson",
//	 "markdown": nil,
//   "html": "<h1>The quick brown fox</h1><p>Jumped over <em>the</em> lazy dog</p>"
//	 ...
// }
func apiGetFileInDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	directory := vestigo.Param(r, "directory")
	filename := vestigo.Param(r, "file")

	file, err := getConvertedFile(directory, filename)
	if err != nil {
		panic(fmt.Errorf("failed to get file %s", err.Error()))
	}

	output, err := json.Marshal(file)
	if err != nil {
		Error.Println(err.Error())
	}

	w.WriteHeader(200)
	w.Write(output)
}

// apiEditFileInDirectoryHandler returns a File object representing the
// specified file to be used on the editor page of the application. A
// server-renedered preview isn't shown, so we don't generate HTML but
// just send the raw Markdown
//
// GET /api/directories/:directory/files/:filename/edit
//
// returns
//
// {
//   "filename": "document_3.md",
//   "path": "documents",
//   "author": "Carl Carlson",
//	 "markdown": "# The quick brown fox\nJumped over *the* lazy dog",
//	 "html": nil
//	 ...
// }
func apiEditFileInDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	directory := vestigo.Param(r, "directory")
	filename := vestigo.Param(r, "file")

	file, err := getRawFile(directory, filename)
	if err != nil {
		panic(fmt.Errorf("failed to get file %s", err.Error()))
	}

	output, err := json.Marshal(file)
	if err != nil {
		Error.Println(err.Error())
	}

	w.WriteHeader(200)
	w.Write(output)
}
