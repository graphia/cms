package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/asdine/storm"
	"github.com/dgrijalva/jwt-go"
	"github.com/gliderlabs/ssh"
	"github.com/husobee/vestigo"
	"github.com/urfave/negroni"
	gossh "golang.org/x/crypto/ssh"
	"gopkg.in/go-playground/validator.v9"
)

var (
	config Config

	// This was set to default to /etc/ but VSCode's Go debugger config isn't working properly
	// see, https://github.com/Microsoft/vscode-go/issues/1134 so for ease now set it to the
	// location of the test config
	//
	argConfigPath = flag.String("config", "/etc/graphia.yml", "the config file")
	logEnabled    = flag.Bool("log-to-file", false, "enable logging")
	verifyKey     *rsa.PublicKey
	signKey       *rsa.PrivateKey
	db            storm.DB
	validate      *validator.Validate
)

// init loads config and sets up logging without requiring
// main to be run, allowing us to log while testing too
func init() {

	var p, envConfigPath string
	var err error

	flag.Parse()

	// If a CONFIG env var exists, use its value
	// to retrieve the config file, otherwise fall
	// back to the command line arg
	envConfigPath = os.Getenv("CONFIG")

	if envConfigPath != "" {
		p = envConfigPath
	} else {
		p = *argConfigPath
	}

	config, err = loadConfig(&p)
	if err != nil {
		panic(err)
	}

	validate = validator.New()

	setupLogging(*logEnabled)
}

func main() {

	Debug.Println("Initialised with config:", config)

	var r, pr *vestigo.Router
	var n *negroni.Negroni

	setupKeys()
	r = unprotectedRouter()
	pr = protectedRouter()
	n = setupMiddleware(r, pr)
	db = setupDB(config.Database)

	Debug.Println("Router and Middleware set up")

	if config.SSHEnabled {
		setupSSH()
	}

	n.Run(fmt.Sprintf(":%s", config.Port))
}

func setupKeys() {
	signBytes, err := ioutil.ReadFile(config.PrivateKeyPath)
	if err != nil {
		panic(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}

	verifyBytes, err := ioutil.ReadFile(config.PublicKeyPath)
	if err != nil {
		panic(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}
}

func setupMiddleware(r, pr *vestigo.Router) (n *negroni.Negroni) {
	n = negroni.New()
	n.UseHandler(r)

	r.Handle("/api/*", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(pr),
	))

	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	return
}

// These routes are proteced
func unprotectedRouter() (r *vestigo.Router) {
	r = vestigo.NewRouter()

	// authentication endpoints
	r.Post("/auth/login", authLoginHandler)

	// setup endpoints
	r.Get("/setup/create_initial_user", setupAllowCreateInitialUserHandler)
	r.Post("/setup/create_initial_user", setupCreateInitialUserHandler)

	// rather than above rule, do a check to see if the file exists and serve it
	// if it doesn't, serve index.html :>
	r.HandleFunc("/cms", cmsGeneralHandler)
	r.HandleFunc("/cms/*", cmsGeneralHandler)

	// serve everything in build by default
	r.Handle("/*", http.FileServer(http.Dir(config.Static)))

	if config.CORSEnabled {

		Warning.Println("CORS:", config.CORSEnabled)
		Warning.Println("CORS origin:", config.CORSOrigin)

		r.SetGlobalCors(&vestigo.CorsAccessControl{
			AllowOrigin:      []string{"*", config.CORSOrigin},
			AllowHeaders:     []string{"Authorization"},
			AllowCredentials: true,
			MaxAge:           3600 * time.Second,
		})
	}

	return r
}

// These routes are proteced 👮
// A JWT is required
func protectedRouter() (r *vestigo.Router) {
	r = vestigo.NewRouter()

	r.Get("/api", apiRootHandler)

	// initial setup endpoints
	r.Get("/api/setup/initialize_repository", apiSetupAllowInitializeRepositoryHandler)
	r.Post("/api/setup/initialize_repository", apiSetupInitializeRepositoryHandler)

	// auth-related endpoints
	r.Post("/api/renew", authRenewTokenHandler)

	// directory endpoints
	r.Get("/api/summary", apiDirectorySummary)
	r.Get("/api/directories", apiListDirectoriesHandler)
	r.Get("/api/directories/:directory", apiGetDirectoryMetadataHandler)
	r.Patch("/api/directories/:directory", apiUpdateDirectoriesHandler)
	r.Post("/api/directories", apiCreateDirectoryHandler)
	r.Delete("/api/directories/:directory", apiDeleteDirectoryHandler)

	// file endpoints
	r.Get("/api/directories/:directory/files", apiListFilesInDirectoryHandler)
	r.Post("/api/directories/:directory/files", apiCreateFileInDirectoryHandler)
	r.Get("/api/directories/:directory/files/:file", apiGetFileInDirectoryHandler)
	r.Get("/api/directories/:directory/files/:file/edit", apiEditFileInDirectoryHandler)
	r.Patch("/api/directories/:directory/files/:file", apiUpdateFileInDirectoryHandler)
	r.Delete("/api/directories/:directory/files/:file", apiDeleteFileFromDirectoryHandler)
	r.Post("/api/directories/:directory/files/:file/translate", apiTranslateFileHandler)

	r.Get("/api/directories/:directory/files/:file/history", apiGetFileHistoryHandler)

	// attachment endpoint
	// note filename used rather than :file because we're not using the extension
	r.Get("/api/directories/:directory/files/:filename/attachments", apiGetFileAttachmentsHandler)
	r.Get("/api/directories/:directory/files/:filename/attachments/:file", apiGetFileAttachmentHandler)

	r.Get("/api/users", apiListUsersHandler)
	r.Post("/api/users", apiCreateUserHandler)
	r.Get("/api/users/:username", apiGetUserHandler)
	r.Post("/api/users/:username", apiUpdateUserHandler)
	r.Delete("/api/users/:username", apiDeleteUserHandler)

	// repo endpoints
	r.Get("/api/repository_info", apiGetRepositoryInformationHandler)
	r.Get("/api/recent_commits", apiGetCommitsHandler)
	r.Get("/api/commits/:hash", apiGetCommitHandler)

	// cms endpoints
	r.Post("/api/publish", apiPublishHandler)
	r.Get("/api/translation_info", apiGetLanguageInformationHandler)

	// missing operations:
	// how should file and directory moves/copies be represented?
	// auth..

	return r
}

func setupDB(path string) storm.DB {
	stormDB, err := storm.Open(path)
	if err != nil {
		panic(fmt.Sprintf("Database cannot be openend %s", err.Error()))
	}
	return *stormDB
}

func setupSSH() {

	ssh.Handle(func(s ssh.Session) {
		authorizedKey := gossh.MarshalAuthorizedKey(s.PublicKey())
		io.WriteString(s, fmt.Sprintf("public key used by %s:\n", s.User()))
		s.Write(authorizedKey)
	})

	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return false
	})

	Debug.Println("starting ssh server on port 2222...")
	ssh.ListenAndServe(":2222", nil, publicKeyOption)
}
