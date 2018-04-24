package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/husobee/vestigo"
)

// SuccessResponse contains information about a successful
// update to the repository
type SuccessResponse struct {
	Message string `json:"message"`
	Oid     string `json:"oid,omitempty"`
	Meta    string `json:"meta,omitempty"`
}

// FailureResponse accompanies the HTTP status code with
// some more information as to why the update failed
type FailureResponse struct {
	Message string `json:"message"`
	Meta    string `json:"meta,omitempty"`
}

// HTTPS Redirect 👉
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {

	if config.Host == "" {
		Error.Println("Host not configured in config, cannot redirect ot HTTPS")
	}

	target := fmt.Sprintf(
		"https://%s:%s%s",
		config.Host,
		config.HTTPSListenPort,
		r.RequestURI,
	)

	Debug.Println("redirecting to HTTPS", target)
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}

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
	var fr FailureResponse
	var uc UserCredentials

	err := json.NewDecoder(r.Body).Decode(&uc)
	if err != nil {
		fr = FailureResponse{Message: "Forbidden"}
		JSONResponse(fr, http.StatusForbidden, w)
		return
	}

	user, err := getUserByUsername(uc.Username)
	if err != nil {
		fr = FailureResponse{Message: fmt.Sprintf("User not found: %s", uc.Username)}
		JSONResponse(fr, http.StatusNotFound, w)
		return
	}

	err = user.checkPassword(uc.Password)
	if err != nil {
		fr = FailureResponse{
			Message: "Invalid credentials",
		}
		JSONResponse(fr, http.StatusUnauthorized, w)
		return
	}

	if !user.Active {
		fr = FailureResponse{
			Message: "User is inactive",
		}
		JSONResponse(fr, http.StatusForbidden, w)
		return
	}

	token, err := newToken(user)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintln("Failed to create new token", err.Error()),
		}
		JSONResponse(fr, http.StatusInternalServerError, w)
		return
	}

	tokenString, err := newTokenString(token)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintln("Failed to create new token string", err.Error()),
		}
		JSONResponse(fr, http.StatusInternalServerError, w)
		return
	}

	Debug.Println("Setting user token", tokenString)

	err = user.setToken(tokenString)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintln("Failed to set the user token", err.Error()),
		}
		JSONResponse(fr, http.StatusInternalServerError, w)
		return
	}

	type output struct {
		Token       `json:"jwt"`
		LimitedUser `json:"user"`
	}

	response := output{
		Token:       Token{tokenString},
		LimitedUser: user.limitedUser(),
	}

	JSONResponse(response, http.StatusOK, w)

}

func apiLogoutHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	var err error
	var fr FailureResponse
	var sr SuccessResponse

	user = getCurrentUser(r.Context())

	err = user.unsetToken()
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Not a git repository"),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	// token is unset
	sr = SuccessResponse{Message: "Logged out"}
	JSONResponse(sr, http.StatusOK, w)

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

	ot, err := request.ParseFromRequest(
		r,
		request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		},
	)

	if err != nil {
		response := FailureResponse{Message: "Could not authorise user"}
		JSONResponse(response, http.StatusUnauthorized, w)
		return
	}

	if !ot.Valid {
		response := FailureResponse{Message: "Token invalid"}
		JSONResponse(response, http.StatusUnauthorized, w)
		return
	}

	user := getCurrentUser(r.Context())

	nt, err := newToken(user)
	if err != nil {
		response := FailureResponse{Message: "Could not create new token"}
		JSONResponse(response, http.StatusBadRequest, w)
		return
	}

	tokenString, err := newTokenString(nt)
	if err != nil {
		response := FailureResponse{Message: "Could not extract string from new token"}
		JSONResponse(response, http.StatusBadRequest, w)
		return
	}

	user.setToken(tokenString)

	response := Token{tokenString}
	JSONResponse(response, http.StatusOK, w)

}

// setupAllowCreateInitialUser will simply return true if there are no users
// and false if there are some
//
// GET /setup/show_initial_setup
//
// {"enabled": false}
func setupAllowCreateInitialUserHandler(w http.ResponseWriter, r *http.Request) {
	var zeroUsers bool
	var so SetupOption

	count, err := countUsers()
	if err != nil {
		response := FailureResponse{Message: "Could not perform user count"}
		JSONResponse(response, http.StatusBadRequest, w)
		return
	}

	zeroUsers = (count == 0)

	so = SetupOption{Enabled: zeroUsers}
	JSONResponse(so, http.StatusOK, w)
}

// setupAllowInitializeRepository will return true if a Git repository does not
// exist at the location specified in `config.Repository`
//
// GET /setup/initialise_repository
//
// {"enabled": false}

func apiSetupAllowInitializeRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var response SetupOption

	err = canInitializeGitRepository(config.Repository)

	if err != nil {
		response = SetupOption{Enabled: false, Meta: err.Error()}
		JSONResponse(response, http.StatusOK, w)
		return
	}

	response = SetupOption{Enabled: true}
	JSONResponse(response, http.StatusOK, w)

}

// setupInitializeRepository will initialize an empty Git repository in the location
// specified by `config.Repository`, providing one doesn't already exist there
//
// POST /setup/create_repository
//
// {"oid": "a741330fec...", message: "Repository initialised"}
func apiSetupInitializeRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	path := config.Repository

	err = canInitializeGitRepository(path)
	if err != nil {
		fr := FailureResponse{Message: "Cannot initialize repository, see log"}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	user := getCurrentUser(r.Context())

	oid, err := initializeGitRepository(user, config.Repository)

	sr := SuccessResponse{Message: "Repository initialised", Oid: oid.String()}
	JSONResponse(sr, http.StatusOK, w)

}

// setupCreateInitialUser allows for the creation of the system's first user and
// unlike typical user creation, doesn't require the instigator to be logged in
//
// POST /setup/create_initial_user
//
// {"username": "lhutz", "name": "Lionel Hutz" ...}
//
// If successful, the response should be a token:
//
// {token: "xxxxx.yyyyy.zzzzz"}
func setupCreateInitialUserHandler(w http.ResponseWriter, r *http.Request) {

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

	// get details from body and ensure active
	user := User{}
	json.NewDecoder(r.Body).Decode(&user)

	user.Active = true
	user.Admin = true

	// create the user
	err = createUser(user)

	if err != nil {
		errors := validationErrorsToJSON(err)
		JSONResponse(errors, http.StatusBadRequest, w)
		return
	}

	sr = SuccessResponse{
		Message: "User created",
	}

	JSONResponse(sr, http.StatusCreated, w)

}

func setupGetUserByConfirmationKeyHandler(w http.ResponseWriter, r *http.Request) {
	confKey := vestigo.Param(r, "confirmation_key")
	lu, err := getLimitedUserByConfirmationKey(confKey)

	if err != nil {
		Error.Println("Failed to find user with key", confKey)
		errors := validationErrorsToJSON(err)
		JSONResponse(errors, http.StatusBadRequest, w)
		return
	}

	JSONResponse(lu, http.StatusOK, w)

}

func setupActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	var sr SuccessResponse
	var fr FailureResponse

	type newPassword struct {
		Password string
	}

	confKey := vestigo.Param(r, "confirmation_key")
	u, err := getUserByConfirmationKey(confKey)

	if err != nil {
		Error.Println("Failed to find user with key", confKey)
		fr = FailureResponse{Message: fmt.Sprintf("Failed to find user with key %s", confKey)}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	var np newPassword
	json.NewDecoder(r.Body).Decode(&np)

	err = activateUser(u, np.Password)
	if err != nil {
		errors := validationErrorsToJSON(err)
		JSONResponse(errors, http.StatusBadRequest, w)
		return
	}

	sr = SuccessResponse{Message: "User activated"}
	JSONResponse(sr, http.StatusOK, w)

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

	if hasImageExt(uri) {
		// if we're dealing with an image we *don't* necessarily know from
		// where we'll be serving it; it could be from the preview or the
		// editor but as far as the actual document is concerned the path
		// is relative
		path = extractImagePath(uri)
	} else {
		// for everything else, the path is the file's path from the static
		// directory
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
	var fr FailureResponse

	directories, err = listRootDirectories()

	if err != nil {

		var msg = strings.ToLower(err.Error())

		// no directory found 404
		if strings.HasPrefix(msg, "failed to resolve path") {
			fr = FailureResponse{
				Message: fmt.Sprintf("No repository found"),
			}
			JSONResponse(fr, http.StatusNotFound, w)
			return
		}

		// directory found but not git-controlled 400
		if strings.HasPrefix(msg, "could not find repository") {
			fr = FailureResponse{
				Message: fmt.Sprintf("Not a git repository"),
			}
			JSONResponse(fr, http.StatusBadRequest, w)
			return
		}

		// anything else
		fr = FailureResponse{
			Message: fmt.Sprintf("Could not retrieve directories: %s", msg),
		}

		JSONResponse(fr, http.StatusBadRequest, w)
		return

	}

	JSONResponse(directories, http.StatusOK, w)
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
	var fr FailureResponse
	var summary []DirectorySummary
	var err error

	summary, err = listRootDirectorySummary()

	if err != nil {
		fr = FailureResponse{
			Message: "No specified directory matches path",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(summary, http.StatusOK, w)

}

// apiCreateDirectoryHandler creates an 'empty' directory. It actually
// contains a _index.md file, so it's trackable by git
//
// POST /api/directories
// {
//	  "name": "Martin Prince",
//	  "email": "mp@springfield-elementary.gov",
//	  "message": "Added new directory called Bobbins",
//	  "directories": [
//	    {
//	      "path": "documents",
//	      "info": {"title": "Documents", "description": "Blah", "body": "Markdown"}
//	    },
//	  ]
// }
func apiCreateDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var nc NewCommit
	var sr SuccessResponse
	var fr FailureResponse

	json.NewDecoder(r.Body).Decode(&nc)

	user := getCurrentUser(r.Context())

	oid, err := createDirectories(nc, user)

	if err == ErrRepoOutOfSync {
		fr = FailureResponse{Message: "Repository out of sync with commit"}
		JSONResponse(fr, http.StatusConflict, w)
		return
	}

	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to create directory: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	sr = SuccessResponse{
		Message: "Directory created",
		Oid:     oid.String(),
	}

	JSONResponse(sr, http.StatusCreated, w)

}

func apiRenameDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusOK)
}

// apiUpdateDirectoriesHandler can update one or more sets of directory metadata,
// this is stored in the `_index.md` file as regular Frontmatter, eg:
// ---
// title: My favourite document
// description: The greatest doc ever!
// ---
func apiUpdateDirectoriesHandler(w http.ResponseWriter, r *http.Request) {
	var directory string
	var nc NewCommit
	var sr SuccessResponse
	var fr FailureResponse

	directory = vestigo.Param(r, "directory")

	nc = NewCommit{}

	json.NewDecoder(r.Body).Decode(&nc)

	if len(nc.Directories) == 0 {
		fr = FailureResponse{
			Message: "No directories specified for updates",
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

	oid, err := updateDirectories(nc, user)

	if err == ErrRepoOutOfSync {
		fr = FailureResponse{Message: "Repository out of sync with commit"}
		JSONResponse(fr, http.StatusConflict, w)
		return
	}

	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to update directories: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	sr = SuccessResponse{
		Message: "Directories updated",
		Oid:     oid.String(),
	}

	JSONResponse(sr, http.StatusOK, w)

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

	if err == ErrRepoOutOfSync {
		fr = FailureResponse{Message: "Repository out of sync with commit"}
		JSONResponse(fr, http.StatusConflict, w)
		return
	}

	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintln("Failed to delete directory:", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	Warning.Println("Directory deleted", directory)

	sr = SuccessResponse{
		Message: "Directory deleted",
		Oid:     oid.String(),
	}

	JSONResponse(sr, http.StatusCreated, w)
}

// Inside a directory functionality 🗂

// apiListFilesInDirectoryHandler returns a JSON array containing
// the all files belonging to :directory
//
// GET /api/directories/:directory/documents
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

	if err == ErrDirectoryNotFound {
		fr = FailureResponse{
			Message: ErrDirectoryNotFound.Error(),
		}
		JSONResponse(fr, http.StatusNotFound, w)
		return

	} else if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Could not get list of files in directory %s: %s", directory, err.Error()),
		}
		JSONResponse(fr, http.StatusNotFound, w)
		return
	}

	metadata, err := getMetadataFromDirectory(directory)

	if err == ErrMetadataNotFound {
		Warning.Println("No metadata file found for", directory)
	} else if err != nil {

		fr = FailureResponse{
			Message: fmt.Sprintln("Something went wrong while retrieving metadata", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	type output struct {
		Files         []FileItem     `json:"files"`
		DirectoryInfo *DirectoryInfo `json:"info,omitempty"`
	}

	result := output{Files: files, DirectoryInfo: metadata}

	JSONResponse(result, http.StatusOK, w)
}

// apiGetDirectoryMetadata returns the `DirectoryInfo` for the
// given directory
//
// GET /api/directories/:directory
//
// {
//   title: "Krusty Burger",
//   description: "Krusty Burger, Ribwich and Breakfast Balls"
// }
func apiGetDirectoryMetadataHandler(w http.ResponseWriter, r *http.Request) {
	var di *DirectoryInfo
	var fr FailureResponse
	var err error

	directory := vestigo.Param(r, "directory")
	di = &DirectoryInfo{}

	di, err = getMetadataFromDirectory(directory)

	if err == ErrMetadataNotFound {

		fr = FailureResponse{
			Message: fmt.Sprintln("Could not find _index.md file in", directory, err.Error()),
		}
		JSONResponse(fr, http.StatusNotFound, w)
		return

	} else if err != nil {

		fr = FailureResponse{
			Message: fmt.Sprintln("Something went wrong while retrieving metadata", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(&di, http.StatusOK, w)

}

// apiCreateFileInDirectory creates a file the specified directory
//
// POST /api/directories/:directory/documents
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
	var err error

	// FIXME should check that params match at least once file in nc.Files
	json.NewDecoder(r.Body).Decode(&nc)

	err = validate.Struct(nc)
	if err != nil {
		errors := validationErrorsToJSON(err)
		JSONResponse(errors, http.StatusBadRequest, w)
		return
	}

	user := getCurrentUser(r.Context())

	oid, err := createFiles(nc, user)

	// If err is a ErrRepoOutOfSync, return a 409 (Edit Conflict) and appropriate message
	if err == ErrRepoOutOfSync {
		fr = FailureResponse{Message: "Repository out of sync with commit"}
		JSONResponse(fr, http.StatusConflict, w)
		return
	}

	if err != nil {
		fr = FailureResponse{Message: fmt.Sprintf("Failed to create files: %s", err.Error())}
		JSONResponse(fr, http.StatusOK, w)
		return
	}

	Debug.Println("File(s) created", oid)

	sr = SuccessResponse{
		Message: "File(s) created",
		Oid:     oid.String(),
	}

	JSONResponse(sr, http.StatusCreated, w)

}

// apiUpdateFileInDirectory updates an existing file the specified
// directory
//
// PATCH /api/directories/:directory/documents/:document/files/:filename
// {
//	  "message": "Added document six"
//	  "files": [
//	    {"filename": "document_six.md", "path": "documents", "extension": "md", "body": "what a nice doc"}
//	  ]
// }
func apiUpdateFileInDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var filename, document, directory string
	var nc NewCommit
	var fr FailureResponse
	var sr SuccessResponse
	var err error

	filename = vestigo.Param(r, "file")
	document = vestigo.Param(r, "document")
	directory = vestigo.Param(r, "directory")

	json.NewDecoder(r.Body).Decode(&nc)

	err = validate.Struct(nc)
	if err != nil {
		errors := validationErrorsToJSON(err)
		JSONResponse(errors, http.StatusBadRequest, w)
		return
	}

	if len(nc.Files) == 0 {
		fr = FailureResponse{
			Message: "No files specified for update",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	if !pathInFiles(directory, document, filename, &nc.Files) {
		fr = FailureResponse{
			Message: "No supplied file matches path",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	user := getCurrentUser(r.Context())

	oid, err := updateFiles(nc, user)

	// If err is a ErrRepoOutOfSync, return a 409 (Edit Conflict) and appropriate message
	if err == ErrRepoOutOfSync {
		fr = FailureResponse{Message: "Repository out of sync with commit"}
		JSONResponse(fr, http.StatusConflict, w)
		return
	}

	if err != nil {

		Error.Println("Failed to update files", nc.Files, err.Error())

		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to update files: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	Warning.Println("File updated", oid)

	sr = SuccessResponse{
		Message: "File updated",
		Oid:     oid.String(),
	}

	JSONResponse(sr, http.StatusCreated, w)

}

func apiTranslateFileHandler(w http.ResponseWriter, r *http.Request) {
	var filename, document, directory string
	var nt NewTranslation
	var fr FailureResponse
	var sr SuccessResponse
	var err error

	filename = vestigo.Param(r, "file")
	document = vestigo.Param(r, "document")
	directory = vestigo.Param(r, "directory")

	user := getCurrentUser(r.Context())

	json.NewDecoder(r.Body).Decode(&nt)

	if directory != nt.Path {
		fr = FailureResponse{Message: "Directory does not match payload"}
		Warning.Printf(
			"Directory does not match contents, param: %s, payload: %s",
			directory,
			nt.Path,
		)
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	if filename != nt.SourceFilename || document != nt.SourceDocument {
		fr = FailureResponse{Message: "Filename does not match payload"}
		Warning.Printf(
			"Filename does not match contents, param: %s, payload: %s",
			filename,
			nt.SourceFilename,
		)
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	oid, fn, err := createTranslation(nt, user)

	if err == ErrRepoOutOfSync {
		fr = FailureResponse{Message: "Repository out of sync with commit"}
		JSONResponse(fr, http.StatusConflict, w)
		return
	}

	if err != nil {
		Error.Println("Could not create translation:", err.Error())
		fr = FailureResponse{Message: err.Error()}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	sr = SuccessResponse{Message: "Translation created", Oid: oid.String(), Meta: fn}

	JSONResponse(sr, http.StatusCreated, w)

}

// apiDeleteFileFromDirectoryHandler deletes a file from the specified
// directory
//
// DELETE /api/directories/:directory/documents/:document/files/:filename
// {
//	  "message": "Deleted document 6 as it's no longer required",
//	  "files": [
//	    {"filename": "document_six.md", "path": "documents"}
//	  ]
// }
func apiDeleteFileFromDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var filename, document, directory string
	var nc NewCommit
	var fr FailureResponse
	var sr SuccessResponse

	filename = vestigo.Param(r, "file")
	document = vestigo.Param(r, "document")
	directory = vestigo.Param(r, "directory")

	json.NewDecoder(r.Body).Decode(&nc)

	if len(nc.Files) == 0 {
		fr = FailureResponse{
			Message: "No files specified for deletion",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	if !pathInFiles(directory, document, filename, &nc.Files) {
		fr = FailureResponse{
			Message: "No supplied file matches path",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	// if there's no commit message, add one
	// FIXME perhaps if more than one file is deleted we could list them
	// in the message rather than just the one specified by the URL
	// UI only offers one at a time so far so not vital.
	if nc.Message == "" {
		nc.Message = fmt.Sprintf("File deleted %s/%s/%s", directory, document, filename)
	}

	user := getCurrentUser(r.Context())

	oid, err := deleteFiles(nc, user)

	// If err is a ErrRepoOutOfSync, return a 409 (Edit Conflict) and appropriate message
	if err == ErrRepoOutOfSync {
		fr = FailureResponse{Message: "Repository out of sync with commit"}
		JSONResponse(fr, http.StatusConflict, w)
		return
	}

	if err != nil {
		Error.Println("Failed to delete directory:", err.Error())
		fr = FailureResponse{
			Message: err.Error(),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	Warning.Println("File updated with commit", oid)

	sr = SuccessResponse{
		Message: "File deleted",
		Oid:     oid.String(),
	}

	JSONResponse(sr, http.StatusCreated, w)

}

// apiGetFileInDirectoryHandler returns a File object representing the
// specified file to be displayed in a web page; the raw Markdown is not
// sent but the compiled HTML is.
//
// GET /api/directories/:directory/documents/:document/files/:filename
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
	document := vestigo.Param(r, "document")
	filename := vestigo.Param(r, "file")

	file, err := getConvertedFile(directory, document, filename)
	if err != nil {

		Error.Println("Could not find converted file", directory, document, filename, err.Error())

		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to get converted file: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(file, http.StatusOK, w)

}

func apiGetFileAttachmentsHandler(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse

	directory := vestigo.Param(r, "directory")
	document := vestigo.Param(r, "document")

	path := fmt.Sprintf("%s/%s", directory, document)

	files, err := getAttachments(path)

	if err != nil {

		Warning.Println("No attachments dir found for path", path)

		fr = FailureResponse{
			Message: "No attachments",
		}

		JSONResponse(fr, http.StatusNotFound, w)

		return
	}

	JSONResponse(files, http.StatusOK, w)
}

func apiGetFileAttachmentHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Not implemented yet"))
}

// apiEditFileInDirectoryHandler returns a File object representing the
// specified file to be used on the editor page of the application. A
// server-renedered preview isn't shown, so we don't generate HTML but
// just send the raw Markdown
//
// GET /api/directories/:directory/documents/:document/files/:filename/edit
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
	var fr FailureResponse

	directory := vestigo.Param(r, "directory")
	document := vestigo.Param(r, "document")
	filename := vestigo.Param(r, "file")

	file, err := getRawFile(directory, document, filename)

	if err != nil {
		Error.Println("Could not update file", err.Error())

		fr = FailureResponse{
			Message: fmt.Sprintln("Could not update file", file, err.Error()),
		}

		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(file, http.StatusOK, w)
}

// user functionality 👩🏽‍💻

// apiListUsers
func apiListUsersHandler(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse

	users, err := allUsers()
	if err != nil {

		Error.Println("Could not retrieve user list", err.Error())

		fr = FailureResponse{
			Message: fmt.Sprintf("Could not retrieve list of users"),
		}

		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(users, http.StatusOK, w)

}

// apiGetUser
func apiGetUserHandler(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse

	username := vestigo.Param(r, "username")
	Debug.Println("Retrieving user", username)

	user, err := getLimitedUserByUsername(username)
	if err != nil {

		Error.Println("Could not retrieve user", username, err.Error())

		fr = FailureResponse{
			Message: fmt.Sprintln("Failed to find restricted user", username, err.Error()),
		}
		JSONResponse(fr, http.StatusNotFound, w)
		return
	}

	JSONResponse(user, http.StatusOK, w)
}

// apiUpdateUser
func apiUpdateUserNameHandler(w http.ResponseWriter, r *http.Request) {

	// FIXME add validation!
	type payload struct {
		Name string `json:"name"`
	}

	var params payload
	var err error
	var fr FailureResponse

	user := getCurrentUser(r.Context())

	// only updating the name is supported at this time

	json.NewDecoder(r.Body).Decode(&params)

	err = db.UpdateField(&user, "Name", params.Name)
	if err != nil {
		fr = FailureResponse{
			Message: err.Error(),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(SuccessResponse{Message: "User updated"}, http.StatusOK, w)
}

// apiUpdatePasswordHandler updates the password after confirming that the
// user knows their current password
//
// PATCH /api/settings/password
//
// returns
//
// {
//   "current_password": "p455w0rd",
//   "new_password": "passwORDx123",
// }
func apiUpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {

	var fr FailureResponse
	var err error

	var pl PasswordUpdate

	user := getCurrentUser(r.Context())

	json.NewDecoder(r.Body).Decode(&pl)

	Debug.Println("password update payload", pl)

	err = validate.Struct(pl)
	if err != nil {
		JSONResponse(validationErrorsToJSON(err), http.StatusBadRequest, w)
		return
	}

	// check current matches old
	err = user.checkPassword(pl.CurrentPassword)
	if err != nil {
		fr = FailureResponse{
			Message: "Current password is not correct",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	// everything's ok, update
	err = user.setPassword(pl.NewPassword)
	if err != nil {
		fr = FailureResponse{
			Message: "Password could not be set",
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(SuccessResponse{Message: "Password updated"}, http.StatusOK, w)

}

func apiUserListPublicKeysHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	var fr FailureResponse

	user = getCurrentUser(r.Context())

	type rk struct {
		ID          int    `json:"id"`
		Fingerprint string `json:"fingerprint"`
		Raw         string `json:"raw"`
		Name        string `json:"name"`
	}

	upks, err := user.keys()
	if err != nil {
		Debug.Println("error", err.Error())
		fr = FailureResponse{
			Message: fmt.Sprintf("Cannot retrieve keys belonging to: %s", user.Username),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	// properly initialise the slice so if empty,
	// marshalled JSON is a empty array instead of null
	var keys []rk
	keys = make([]rk, 0)

	for _, upk := range upks {

		file, err := upk.File()
		if err != nil {
			Error.Println("Could not decode key", upk.Fingerprint)
			continue
		}

		keys = append(keys, rk{
			ID:          upk.ID,
			Name:        upk.Name,
			Fingerprint: upk.Fingerprint,
			Raw:         file,
		})
	}

	JSONResponse(keys, http.StatusOK, w)

}

func apiUserAddPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Name string `json:"name"`
		Key  string `json:"key"`
	}

	var user User
	var pl payload
	var err error
	var sr SuccessResponse
	var fr FailureResponse

	json.NewDecoder(r.Body).Decode(&pl)

	user = getCurrentUser(r.Context())

	Debug.Println("supplied key:", pl.Key)

	err = user.addPublicKey(pl.Name, pl.Key)

	if err != nil && err.Error() == "already exists" {
		Error.Println("Key already exists", user.Username, err.Error(), pl.Key)
		fr = FailureResponse{
			Message: fmt.Sprintf("Key already exists\n%s", pl.Key),
		}
		JSONResponse(fr, http.StatusConflict, w)
		return
	}

	if err != nil {
		Error.Println("Failed to create public key", user.Username, err.Error())
		fr = FailureResponse{
			Message: fmt.Sprintf("Cannot set public key for %s", user.Username),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	sr = SuccessResponse{
		Message: fmt.Sprintf("Public key created for %s", user.Username),
	}

	JSONResponse(sr, http.StatusOK, w)

}

func apiUserDeletePublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var pkr PublicKey
	var sr SuccessResponse
	var fr FailureResponse
	var user User

	sid := vestigo.Param(r, "id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Invalid id %d", id),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	user = getCurrentUser(r.Context())

	err = db.One("ID", id, &pkr)
	if err != nil {
		Warning.Println("Cannot find public key", id)
		fr = FailureResponse{
			Message: fmt.Sprintf("Cannot find public key %d", id),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}
	if user.ID != pkr.UserID {
		Warning.Println("Key does not belong to user", id, user)
		fr = FailureResponse{
			Message: fmt.Sprintf("Key %d does not belong to you", id),
		}
		JSONResponse(fr, http.StatusForbidden, w)
		return
	}

	err = db.DeleteStruct(&pkr)
	if user.ID != pkr.UserID {
		Error.Println("Could not delete public key", err.Error())
		fr = FailureResponse{
			Message: fmt.Sprintf("Could not delete key %d", id),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	// no errors, success!
	sr = SuccessResponse{
		Message: fmt.Sprintf("Key %d deleted", id),
	}

	JSONResponse(sr, http.StatusOK, w)

}
func apiPublishHandler(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse

	output, err := buildStaticSite()
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to publish site: %s", err.Error()),
			Meta:    string(output),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	Info.Println("Site published!")

	type publishResponse struct {
		Message string `json:"message"`
		Meta    string `json:"meta"`
	}

	pr := publishResponse{
		Message: "Published successfully",
		Meta:    string(output),
	}

	JSONResponse(pr, http.StatusOK, w)

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
func apiGetCommitsHandler(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse
	var commits []Commit
	var err error

	// TODO manage quantity param

	commits, err = getCommits(5)
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintf("Failed to retrieve recent commits: %s", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(commits, http.StatusOK, w)

}

// GET /api/repository_info
//
// returns the repositorys head commit's hash, used to ensure subsequent commits
// aren't applied to an out-of-sync repo
//
// [
//	  {
//	    "latest_revision": "abcde12345"
//	  },
// ]
func apiGetRepositoryInformationHandler(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse
	var ri RepositoryInfo
	var err error

	ri, err = getRepositoryInfo()
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintln("Failed to retrieve repository info", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(ri, http.StatusOK, w)

}

func apiGetServerInformationHandler(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse
	var si ServerInfo
	var err error

	fc, err := countFiles()
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintln("Failed to get file count", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	uc, err := countUsers()
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintln("Failed to get user count", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	cc, err := countCommits()
	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintln("Failed to get commits count", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	si = ServerInfo{
		Title:         config.SiteTitle,
		Users:         uc,
		Counts:        fc,
		Commits:       cc,
		Host:          config.Host,
		SSHListenPort: config.SSHListenPort,
	}

	if err != nil {
		fr = FailureResponse{
			Message: fmt.Sprintln("Failed to retrieve server info", err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(si, http.StatusOK, w)
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
func apiGetCommitHandler(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse
	var hash string

	hash = vestigo.Param(r, "hash")

	cs, err := diffForCommit(hash)
	if err != nil {
		Error.Println("Could not find commit", hash, err.Error())

		fr = FailureResponse{
			Message: fmt.Sprintln("Failed to generate diff for commit", hash, err.Error()),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(cs, http.StatusOK, w)
}

// GET /api/directories/:directory/documents/:document/files/:filename/history
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
func apiGetFileHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse

	directory := vestigo.Param(r, "directory")
	document := vestigo.Param(r, "document")
	filename := vestigo.Param(r, "file")

	path := filepath.Join(directory, document, filename)

	history, err := getFileHistory(path, 10)
	if err != nil {
		Error.Println("Could not get file history for file", path)

		fr = FailureResponse{
			Message: fmt.Sprintln("Could not get file history for file", path),
		}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	JSONResponse(history, http.StatusOK, w)

}

// User data 👩🏽‍💻

// GET /api/user_info
//
// returns the currently logged in User info
func apiGetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := getCurrentUser(r.Context())
	JSONResponse(user.limitedUser(), http.StatusOK, w)
}

// Translation data 💬

// GET /api/translation_info
//
// returns the CMS's translation/language settings
//
// [
//	  {
//	    translation_enabled: true,
//	    default_language: "en",
//	    languages: [
//	      {code: "en", name: "English", flag: "🇬🇧"},
//	      {code: "es", name: "Spanish", flag: "🇪🇸"}
//	    ]
//	  },
// ]
func apiGetLanguageInformationHandler(w http.ResponseWriter, r *http.Request) {

	// return enabled false if translation disabled in config
	if !config.TranslationEnabled {
		response := struct {
			TranslationEnabled bool `json:"translation_enabled"`
		}{
			false,
		}

		Debug.Println("Translation is disabled")
		JSONResponse(response, http.StatusOK, w)
		return
	}

	type language struct {
		Code string `json:"code"`
		Name string `json:"name"`
		Flag string `json:"flag"`
	}

	type languageInfo struct {
		TranslationEnabled bool       `json:"translation_enabled,omitempty"`
		DefaultLanguage    string     `json:"default_language"`
		Languages          []language `json:"languages,omitempty"`
	}

	var languages []language

	// build a list of just the enabled languages
	for _, lc := range config.EnabledLanguages {

		li := config.AllLanguages[lc]

		languages = append(languages, language{
			Code: lc,
			Name: li.Name,
			Flag: li.Flag,
		})

	}

	var li = languageInfo{
		TranslationEnabled: config.TranslationEnabled,
		DefaultLanguage:    config.DefaultLanguage,
		Languages:          languages,
	}

	Debug.Println("Translation is enabled", li)
	JSONResponse(li, http.StatusOK, w)

}
