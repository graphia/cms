package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/asdine/storm"
	"github.com/dgrijalva/jwt-go"
	"github.com/husobee/vestigo"
	"github.com/urfave/negroni"
	"gopkg.in/go-playground/validator.v9"
)

var (
	config        Config
	argConfigPath = flag.String("config", "/etc/graphia/cms.yml", "the config file")
	logEnabled    = flag.Bool("log-to-file", false, "enable logging")
	verifyKey     *rsa.PublicKey
	signKey       *rsa.PrivateKey
	db            storm.DB
	validate      *validator.Validate
	mailer        Mailer
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

	mailer = setupMailer()

	setupLogging(*logEnabled)
}

func main() {

	Debug.Println("Initialised with config:", config)

	var r, pr, ar *vestigo.Router
	var n *negroni.Negroni

	setupKeys()
	r = unprotectedRouter()
	pr = protectedRouter()
	ar = adminRouter()
	n = setupMiddleware(r, pr, ar)
	db = setupDB(config.Database)

	Debug.Println("Router and Middleware set up")

	if config.SSHEnabled {
		Debug.Println("SSH is enabled, listening on port", config.SSHListenPort)
		setupSSH()
	} else {
		Debug.Println("SSH not enabled :(")
	}

	var err error

	if config.HTTPSEnabled {

		errors := make(chan error, 0)

		go func() {

			Info.Println("Redirecting HTTP to HTTPS", config.HTTPListenPort)

			err = http.ListenAndServe(
				config.HTTPListenPortWithColon(),
				http.HandlerFunc(redirectToHTTPS),
			)

			if err != nil {
				errors <- err
			}

		}()

		go func() {

			Info.Println("Listening for SSL connections on", config.HTTPSListenPort)

			err = http.ListenAndServeTLS(
				config.HTTPSListenPortWithColon(),
				config.HTTPSCert,
				config.HTTPSKey,
				n,
			)

			if err != nil {
				errors <- err
			}

		}()

		err = <-errors

		if err != nil {
			panic(err)
		}

	} else {

		// Only listening on HTTP, aka 'dev mode'
		Info.Println("HTTPS is not enabled, listening on", config.HTTPListenPort)
		http.ListenAndServe(config.HTTPListenPortWithColon(), n)
	}

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

func setupMailer() Mailer {
	return Mailer{send: DefaultSender}
}

func setupMiddleware(r, pr, ar *vestigo.Router) (n *negroni.Negroni) {
	n = negroni.New()
	n.UseHandler(r)

	r.Handle("/api/admin/*", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.HandlerFunc(ValidateAdminMiddleware),
		negroni.Wrap(ar),
	))

	r.Handle("/api/*", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(pr),
	))

	//n.Use(negroni.NewLogger())
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

	r.Get("/setup/activate/:confirmation_key", setupGetUserByConfirmationKeyHandler)
	r.Patch("/setup/activate/:confirmation_key", setupActivateUserHandler)

	// rather than above rule, do a check to see if the file exists and serve it
	// if it doesn't, serve index.html :>
	r.HandleFunc("/cms", cmsGeneralHandler)
	r.HandleFunc("/cms/*", cmsGeneralHandler)

	// serve everything in build by default
	r.Handle("/*", http.FileServer(http.Dir(config.Static)))

	return
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
	r.Get("/api/directories/:directory/documents", apiListFilesInDirectoryHandler)
	r.Post("/api/directories/:directory/documents", apiCreateFileInDirectoryHandler)

	r.Get("/api/directories/:directory/documents/:document/files/:file", apiGetFileInDirectoryHandler)
	r.Get("/api/directories/:directory/documents/:document/files/:file/edit", apiEditFileInDirectoryHandler)

	r.Patch("/api/directories/:directory/documents/:document/files/:file", apiUpdateFileInDirectoryHandler)
	r.Delete("/api/directories/:directory/documents/:document/files/:file", apiDeleteFileFromDirectoryHandler)
	r.Post("/api/directories/:directory/documents/:document/files/:file/translate", apiTranslateFileHandler)

	r.Get("/api/directories/:directory/documents/:document/files/:file/history", apiGetFileHistoryHandler)

	// attachment endpoint
	// note filename used rather than :file because we're not using the extension
	r.Get("/api/directories/:directory/documents/:document/attachments", apiGetFileAttachmentsHandler)
	r.Get("/api/directories/:directory/documents/:document/attachments/:file", apiGetFileAttachmentHandler)

	// user retrieval endpoints
	r.Get("/api/users", apiListUsersHandler)
	r.Get("/api/users/:username", apiGetUserHandler)
	r.Get("/api/user_info", apiGetUserInfoHandler)

	r.Post("/api/logout", apiLogoutHandler)

	// user settings and ssh key management
	r.Patch("/api/settings/name", apiUpdateUserNameHandler)
	r.Patch("/api/settings/password", apiUpdatePasswordHandler)
	r.Get("/api/settings/ssh", apiUserListPublicKeysHandler)
	r.Post("/api/settings/ssh", apiUserAddPublicKeyHandler)
	r.Delete("/api/settings/ssh/:id", apiUserDeletePublicKeyHandler)

	// stats endpoints
	r.Get("/api/server_info", apiGetServerInformationHandler)

	// repo endpoints
	r.Get("/api/repository_info", apiGetRepositoryInformationHandler)
	r.Get("/api/recent_commits", apiGetCommitsHandler)
	r.Get("/api/commits/:hash", apiGetCommitHandler)
	r.Get("/api/history", apiGetHistoryHandler)

	// cms endpoints
	r.Post("/api/publish", apiPublishHandler)
	r.Get("/api/translation_info", apiGetLanguageInformationHandler)

	// missing operations:
	// how should file and directory moves/copies be represented?

	return r
}

// Endpoints only available to users who are Administrators
func adminRouter() (r *vestigo.Router) {

	r = vestigo.NewRouter()

	r.Post("/api/admin/users", apiCreateUserHandler)
	r.Get("/api/admin/users/:username", apiGetFullUserHandler)
	r.Patch("/api/admin/users/:username", apiUpdateUserHandler)
	r.Delete("/api/admin/users/:username", apiDeleteUserHandler)

	return r
}

func setupDB(path string) storm.DB {
	stormDB, err := storm.Open(path)
	if err != nil {
		panic(fmt.Sprintf("Database cannot be openend %s", err.Error()))
	}
	return *stormDB
}
