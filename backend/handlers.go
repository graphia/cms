package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/husobee/vestigo"
	"golang.org/x/crypto/bcrypt"
)

// Authentication functionality 🔑

// authLoginHandler checks the supplied UserCredentials and, if a user
// exists and the password matches, returns a JWT
//
// POST /auth/login
//
// {username: "jimbo.jones", password: "psssstsecret"}
//
// If successful, the response:
//
// {token: "xxxxx.yyyyy.zzzzz"}
//
func authLoginHandler(w http.ResponseWriter, r *http.Request) {
	// start session, generate a JWT if the credentials are ok

	var uc UserCredentials

	err := json.NewDecoder(r.Body).Decode(&uc)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	user, err := getUserByUsername(uc.Username)

	if err != nil {
		response := FailureResponse{Message: fmt.Sprintf("User not found: %s", uc.Username)}
		JSONResponse(response, http.StatusBadRequest, w)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(uc.Password))

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := FailureResponse{Message: "Invalid credentials"}
		json, err := json.Marshal(response)
		if err != nil {
			Debug.Println("Failed", err)
			panic(err)
		}
		w.Write(json)
		return
	}

	token, err := newToken(user)
	if err != nil {
		panic(err)
	}
	tokenString, err := newTokenString(token)
	if err != nil {
		panic(err)
	}

	Debug.Println("Setting user token", tokenString)
	err = setToken(user, tokenString)
	if err != nil {
		panic(err)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		panic(err)
	}

	response := Token{tokenString}
	JSONResponse(response, http.StatusOK, w)

}

// authRenewTokenHandler issues a new token providing a valid one is
// provided. This is intended to be called periodically by the client so
// user sessions don't expire while they're still using the site
//
// POST /api/renew
//
// {token: "xxxxx.yyyyy.zzzzz"}
//
func authRenewTokenHandler(w http.ResponseWriter, r *http.Request) {

	token, err := request.ParseFromRequest(
		r,
		request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		},
	)

	if err == nil {
		if token.Valid {

			Debug.Println("Token valid, issuing a new one")
			// TODO print the claims?

			claims := token.Claims.(jwt.MapClaims)
			username := claims["sub"]
			user, err := getUserByUsername("joey")
			if err != nil {
				panic(fmt.Errorf("Cannot find user %s, %s", username, err.Error()))
			}

			token, err := newToken(user)
			if err != nil {
				panic(err)
			}

			tokenString, err := newTokenString(token)
			if err != nil {
				panic(err)
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Error extracting the key")
				panic(err)
			}

			response := Token{tokenString}
			JSONResponse(response, http.StatusOK, w)

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			response := FailureResponse{Message: "Invalid credentials"}
			json, err := json.Marshal(response)
			if err != nil {
				panic(err)
			}
			w.Write(json)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}

}

// authAllowCreateInitialUser will simply return true if there are no users
// and false if there are some
//
// GET /auth/show_initial_setup
//
// {"enabled": false}
func authAllowCreateInitialUser(w http.ResponseWriter, r *http.Request) {
	var zeroUsers bool

	count, err := countUsers()
	if err != nil {
		panic(err)
	}

	zeroUsers = (count == 0)

	output, err := json.Marshal(InitialSetup{Enabled: zeroUsers})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// authCreateInitialUser allows for the creation of the system's first user and
// unlike typical user creation, doesn't require the instigator to be logged in
//
// POST /auth/create_initial_user
//
// {"username": "lhutz", "name": "Lionel Hutz" ...}
//
// If successful, the response should be a token:
//
// {token: "xxxxx.yyyyy.zzzzz"}
func authCreateInitialUser(w http.ResponseWriter, r *http.Request) {

	var sr SuccessResponse
	var fr FailureResponse

	// check that there are no users
	qty, err := countUsers()
	if err != nil {
		// handle error counting users
		Debug.Println(err)
	}

	if qty > 0 {
		// users exist, don't allow a new one to be created
		fr = FailureResponse{
			Message: "Users already exist. Cannot create initial user",
		}
		output, err := json.Marshal(fr)
		if err != nil {
			Error.Println(err.Error())
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	// get details from body and set active
	user := User{}
	json.NewDecoder(r.Body).Decode(&user)
	user.Active = true

	// create the user
	err = createUser(user)

	if err != nil {

		errors := validationErrorsToJSON(err)

		output, err := json.Marshal(errors)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)
		return
	}

	sr = SuccessResponse{
		Message: "User created",
	}

	output, err := json.Marshal(sr)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(201)
	w.Write(output)

}

// cmsGeneralHandler takes care of serving the frontend's CMS portion
// of the app. Paths that match '/cms' are routed here; if the file exists
// on the filesystem (in the location specified by `config.Static`) it's
// served. If not, cms/index.html is served instead
//
// GET /cms                        -> public/cms/index.html
// GET /cms/javascripts/app.js     -> public/cms/javascripts/app.js
// GET /cms/something/nonexistant  -> public/cms/index.html
func cmsGeneralHandler(w http.ResponseWriter, r *http.Request) {

	var index, path, uri string

	uri = r.RequestURI

	// if we're dealing with an image we *don't* necessarily know from
	// where we'll be serving it; it could be from the preview or the
	// editor but as far as the actual document is concerned the path
	// is relative
	if isImageURI(uri) {
		path = extractImagePath(uri)
	} else {
		path = filepath.Join(config.Static, r.URL.Path)
	}

	// look for the file on the filesystem
	_, err := os.Stat(path)

	if err == nil {
		// found it, serve it
		http.ServeFile(w, r, path)
	} else {
		// if we can't find a file (asset) to serve but the path begins
		// with /cms, serve the CMS's index. It is likely someone
		// navigating directly to a resource or refreshing a page. If
		// Vue's route *still* can't match it to a page, it'll show an
		// appropriate error
		index = filepath.Join(config.Static, "cms", "index.html")
		http.ServeFile(w, r, index)
	}
}

func apiRootHandler(w http.ResponseWriter, r *http.Request) {}

// Directory level functionality 🗃

// apiListDirectoriesHandler returns an array of directories
//
// GET /api/directories/
//
// eg. when the documents directory contains Documents 1 and 2:
//
// {
//	  {name: 'documents'},
//	  {name: 'appendices'}
// }
func apiListDirectoriesHandler(w http.ResponseWriter, r *http.Request) {

	var directories []Directory
	var err error

	directories, err = listRootDirectories()

	output, err := json.Marshal(directories)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// apiListDirectoriesHandler returns a JSON object representing
// the Repository's contents.
//
// GET /api/directory_summary
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
func apiDirectorySummary(w http.ResponseWriter, r *http.Request) {

	var summary map[string][]FileItem
	var err error

	summary = make(map[string][]FileItem)

	summary, err = listRootDirectorySummary()

	output, err := json.Marshal(summary)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)

}

// apiCreateDirectoryHandler creates an 'empty' directory. It actually
// contains a hidden `.keep` file, so it's trackable by git
//
// POST /api/directories
// {
//	  "name": "Martin Prince",
//	  "email": "mp@springfield-elementary.gov",
//	  "message": "Added new directory called Bobbins",
//	  "directories": [
//	    {"path": "documents"},
//	    {"path": "appendices"}
//	  ]
// }
func apiCreateDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var nc NewCommit
	var sr SuccessResponse
	var fr FailureResponse

	json.NewDecoder(r.Body).Decode(&nc)

	user := getCurrentUser(r.Context())

	oid, err := createDirectories(nc, user)

	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to create directory: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	sr = SuccessResponse{
		Message: "Directory created",
		Oid:     oid.String(),
	}

	output, err := json.Marshal(sr)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(output)

}

func apiRenameDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusOK)
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
// DELETE /api/directories
// {
//	  "name": "Martin Prince",
//	  "email": "mp@springfield-elementary.gov",
//	  "message": "Added new directory called Bobbins",
//	  "directories": [
//	    {"path": "documents"}
//	  ]
// }
//
// returns a SuccessResponse containing the git commit hash or a FailureResponse
// containing an error message
func apiDeleteDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var directory string
	var nc NewCommit
	var fr FailureResponse
	var sr SuccessResponse

	directory = vestigo.Param(r, "directory")

	// set up the RepoWrite with git params, an appropraiate message and then
	// specify the directory based on the path
	nc = NewCommit{}

	json.NewDecoder(r.Body).Decode(&nc)

	if len(nc.Directories) == 0 {
		fr = FailureResponse{
			Message: "No directories specified for deletion",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	if !pathInDirectories(directory, &nc.Directories) {
		fr = FailureResponse{
			Message: "No specified directory matches path",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	user := getCurrentUser(r.Context())

	oid, err := deleteDirectories(nc, user)

	if err != nil {

		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to delete directory: %s", err.Error()),
		}

		JSONResponse(fr, http.StatusBadRequest, w)

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

	w.WriteHeader(http.StatusOK)
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
	var fr FailureResponse

	directory := vestigo.Param(r, "directory")
	files, err := getFilesInDir(directory)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Could not get list of files in directory %s: %s", directory, err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)

	}

	output, err := json.Marshal(files)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Could not create JSON: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// apiCreateFileInDirectory creates a file the specified directory
//
// POST /api/directories/:directory/files
// {
//	  "message": "Added document six"
//	  "files": [
//	    {"filename": "document_six.md", "path": "documents", "extension": "md", body: "blah blah"}
//	  ]
// }
func apiCreateFileInDirectoryHandler(w http.ResponseWriter, r *http.Request) {

	var nc NewCommit
	var fr FailureResponse
	var sr SuccessResponse

	// FIXME should check that params match at least once file in nc.Files
	json.NewDecoder(r.Body).Decode(&nc)

	user := getCurrentUser(r.Context())

	oid, err := createFiles(nc, user)
	if err != nil {
		fr = FailureResponse{Message: fmt.Sprintf("Failed to create files: %s", err.Error())}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	sr = SuccessResponse{
		Message: "File(s) created",
		Oid:     oid.String(),
	}

	output, err := json.Marshal(sr)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(201)
	w.Write(output)

	Debug.Println("File(s) created", oid)

}

// apiUpdateFileInDirectory updates an existing file the specified
// directory
//
// PATCH /api/directories/:directory/files/:filename
// {
//	  "message": "Added document six"
//	  "files": [
//	    {"filename": "document_six.md", "path": "documents", "extension": "md", "body": "what a nice doc"}
//	  ]
// }
func apiUpdateFileInDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var filename, directory string
	var nc NewCommit
	var fr FailureResponse
	var sr SuccessResponse

	filename = vestigo.Param(r, "file")
	directory = vestigo.Param(r, "directory")

	//  FIXME should check that params match at least once file in nc.Files
	json.NewDecoder(r.Body).Decode(&nc)

	if len(nc.Files) == 0 {
		fr = FailureResponse{
			Message: "No files specified for update",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	if !pathInFiles(directory, filename, &nc.Files) {
		fr = FailureResponse{
			Message: "No supplied file matches path",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	user := getCurrentUser(r.Context())

	oid, err := updateFiles(nc, user)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to update files: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
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
//	  "message": "Deleted document 6 as it's no longer required",
//	  "files": [
//	    {"filename": "document_six.md", "path": "documents"}
//	  ]
// }
func apiDeleteFileFromDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var filename, directory string
	var nc NewCommit
	var fr FailureResponse
	var sr SuccessResponse

	filename = vestigo.Param(r, "file")
	directory = vestigo.Param(r, "directory")

	json.NewDecoder(r.Body).Decode(&nc)

	if len(nc.Files) == 0 {
		fr = FailureResponse{
			Message: "No files specified for deletion",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	if !pathInFiles(directory, filename, &nc.Files) {
		fr = FailureResponse{
			Message: "No supplied file matches path",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	Debug.Println("nc", nc)

	if len(nc.Files) == 0 {
		response := FailureResponse{Message: "No files specified for deletion"}
		JSONResponse(response, http.StatusBadRequest, w)
		return
	}

	user := getCurrentUser(r.Context())

	oid, err := deleteFiles(nc, user)

	if err != nil {
		fr = FailureResponse{
			Message: err.Error(),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
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
	var fr FailureResponse

	directory := vestigo.Param(r, "directory")
	filename := vestigo.Param(r, "file")

	file, err := getConvertedFile(directory, filename)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to get converted file: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	output, err := json.Marshal(file)
	if err != nil {
		Error.Println("Failed to convert file to JSON", file)
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to create JSON from file: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func apiGetFileAttachmentsHandler(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse

	directory := vestigo.Param(r, "directory")
	filename := vestigo.Param(r, "filename")

	path := fmt.Sprintf("%s/%s", directory, filename)

	files, err := getAttachments(path)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to get converted file: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	output, err := json.Marshal(files)
	if err != nil {
		Error.Println("Failed to convert file to JSON", files)
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to create JSON from file: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func apiGetFileAttachmentHandler(w http.ResponseWriter, r *http.Request) {

	Debug.Println("***** HANDLING ATTACHMENT!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
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

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// user functionality 👩🏽‍💻

// apiListUsers
func apiListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := allUsers()
	if err != nil {
		Error.Println("Could not get list of users", err.Error())
	}

	output, err := json.Marshal(users)
	if err != nil {
		Error.Println("Could not create JSON", err.Error())
		Error.Println("Users", users)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(output)

}

// apiGetUser
func apiGetUser(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse

	username := vestigo.Param(r, "username")
	Debug.Println("Retrieving user", username)

	user, err := getLimitedUserByUsername(username)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to find restricted user %s: %s", username, err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	output, err := json.Marshal(user)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to find restricted user %s: %s", username, err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

// apiCreateUser
func apiCreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	var sr SuccessResponse
	var fr FailureResponse

	json.NewDecoder(r.Body).Decode(&user)

	user.Active = true

	err := createUser(user)
	if err != nil {

		errors := validationErrorsToJSON(err)

		output, err := json.Marshal(errors)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(output)

		return
	}

	Debug.Println("User was created successfully")

	sr = SuccessResponse{
		Message: "User created",
	}

	output, err := json.Marshal(sr)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to generate JSON %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	w.WriteHeader(201)
	w.Write(output)

}

// apiUpdateUser
func apiUpdateUser(w http.ResponseWriter, r *http.Request) {}

// apiDeleteUser
func apiDeleteUser(w http.ResponseWriter, r *http.Request) {}

func apiPublish(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse

	output, err := buildStaticSite()
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to publish site: %s", err.Error()),
			Meta:    string(output),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	Info.Println("Site published!")

	sr := SuccessResponse{
		Message: "Published successfully",
	}

	repsonse, err := json.Marshal(sr)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to generate response: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(repsonse)

}

// Repository data 💁

// GET /api/recent_commits
//
// returns the most recent commits made to the repository. Currently hard-coded
// to 20, but would make sense to accept a param
//
// [
//	  {
//	    "message": "Changed some stuff",
//	    "id": "e2da99aa078c",
//	    "object_type": "blob",
//	    "author": "Peter Yates"
//	    "time": "Fri Jul 14 12:34:45 2017 +0100"
//	  },
// ]
func apiGetCommits(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse
	var commits []Commit
	var err error

	// TODO manage quantity param

	commits, err = getCommits(20)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to retrieve recent commits: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	response, err := json.Marshal(commits)

	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to convert recent commits to JSON: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GET /api/commits/:commit_hash
//
// returns a single commit containing relevant info plus the list of files
// and diff information, plus the file's contents before and after the change
//
//	{
//	  "num_deltas": 2,
//	  "num_added": 1,
//	  "num_deleted": 1,
//	  "full_diff": "raw diff text",
//	  "files": {
//	    "documents/document_1.md": {
//	      "old": "the quick brown fox",
//	      "new": "the thick brown fox"
//	    }
//	  },
//	  "message": "changed quick to thick",
//	  "author": "Cletus Spuckler",
//	  "hash": "e2da99aa078c",
//	  "timestamp": "Fri Jul 14 12:34:45 2017 +0100"
//	},
func apiGetCommit(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse
	var hash string

	hash = vestigo.Param(r, "hash")

	cs, err := diffForCommit(hash)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to generate diff for commit %s: %s", hash, err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	output, err := json.Marshal(cs)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to convert diff to JSON: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(output)

}

// GET /api/directories/:directory/files/:filename/history
//
// returns the basic commit information for every commit that has
// affected the specified file
//
// [
//	  {
//	    "message": "Changed some stuff",
//	    "id": "e2da99aa078c",
//	    "object_type": "blob",
//	    "author": "Peter Yates"
//	    "time": "Fri Jul 14 12:34:45 2017 +0100"
//	  },
// ]
func apiGetFileHistory(w http.ResponseWriter, r *http.Request) {
	directory := vestigo.Param(r, "directory")
	filename := vestigo.Param(r, "file")

	path := fmt.Sprintf("%s/%s", directory, filename)

	history, err := getFileHistory(path, 10)
	if err != nil {
		panic(err)
	}

	output, err := json.Marshal(history)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)

}
