package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/husobee/vestigo"
	"github.com/urfave/negroni"
)

const (
	// FIXME load these from config
	privKeyPath = "./keys/app.rsa"
	pubKeyPath  = "./keys/app.rsa.pub"
)

var (
	config         Config
	configFilePath = flag.String("config", "/etc/graphia.yml", "the config file")
	logEnabled     = flag.Bool("log-to-file", false, "enable logging")
	verifyKey      *rsa.PublicKey
	signKey        *rsa.PrivateKey
)

// init loads config and sets up logging without requiring
// main to be run, allowing us to log while testing too
func init() {
	var err error

	flag.Parse()

	config, err = loadConfig(configFilePath)
	if err != nil {
		if err != nil {
			panic(err)
		}
	}

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

	Debug.Println("Router and Middleware set up")

	n.Run(fmt.Sprintf(":%s", config.Port))
}

func setupKeys() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		panic(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
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
	r.Post("/auth/logout", authLogoutHandler)

	// rather than above rule, do a check to see if the file exists and serve it
	// if it doesn't, serve index.html :>
	r.HandleFunc("/cms", cmsGeneralHandler)
	r.HandleFunc("/cms/*", cmsGeneralHandler)

	// TODO duplicated below, tidy up
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

	// auth-related endpoints
	r.Post("/api/renew", authRenewTokenHandler)

	// directory endpoints
	r.Get("/api/directories", apiListDirectoriesHandler)
	r.Post("/api/directories", apiCreateDirectoryHandler)
	r.Delete("/api/directories/:directory", apiDeleteDirectoryHandler)

	// file endpoints
	r.Get("/api/directories/:directory/files", apiListFilesInDirectoryHandler)
	r.Post("/api/directories/:directory/files", apiCreateFileInDirectoryHandler)
	r.Get("/api/directories/:directory/files/:file", apiGetFileInDirectoryHandler)
	r.Get("/api/directories/:directory/files/:file/edit", apiEditFileInDirectoryHandler)
	r.Patch("/api/directories/:directory/files/:file", apiUpdateFileInDirectoryHandler)
	r.Delete("/api/directories/:directory/files/:file", apiDeleteFileFromDirectoryHandler)

	// missing operations:
	// how should file and directory moves/copies be represented?
	// auth..

	return r
}
