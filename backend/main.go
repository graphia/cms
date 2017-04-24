package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/husobee/vestigo"
	"github.com/urfave/negroni"
)

var (
	config         Config
	configFilePath = flag.String("config", "/etc/graphia.yml", "the config file")
	logEnabled     = flag.Bool("log-to-file", false, "enable logging")
)

// init loads config and sets up logging without requiring
// main to be run, allowing us to log while testing too
func init() {
	var err error

	flag.Parse()

	config, err = loadConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	setupLogging(*logEnabled)
}

func main() {

	Debug.Println("Initialised with config:", config)

	var r *vestigo.Router
	var n *negroni.Negroni

	r = setupRouter()
	n = setupMiddleware(r)

	Debug.Println("Router and Middleware set up")

	n.Run(fmt.Sprintf(":%s", config.Port))
}

func setupMiddleware(r *vestigo.Router) (n *negroni.Negroni) {
	n = negroni.New()
	n.UseHandler(r)
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	return
}

func setupRouter() (r *vestigo.Router) {
	r = vestigo.NewRouter()
	r.Get("/api", apiRootHandler)

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

	if config.CORSEnabled {

		Warning.Println("CORS:", config.CORSEnabled)
		Warning.Println("CORS origin:", config.CORSOrigin)

		r.SetGlobalCors(&vestigo.CorsAccessControl{
			AllowOrigin:      []string{"*", config.CORSOrigin},
			AllowCredentials: true,
			MaxAge:           3600 * time.Second,
		})
	}

	// rather than above rule, do a check to see if the file exists and serve it
	// if it doesn't, serve index.html :>
	r.HandleFunc("/cms", cmsGeneralHandler)
	r.HandleFunc("/cms/*", cmsGeneralHandler)

	return
}
