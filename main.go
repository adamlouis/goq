package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/adamlouis/goq/internal/apiserver"
	"github.com/adamlouis/goq/internal/webserver"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	fDotenv = flag.String("dotenv", "", "a .env file from which to read environment variables. useful for local development.")
)

// TODO: revisit config-based init
const (
	_envServerPort              = "GOQ_SERVER_PORT"
	_envSQLite3ConnectionString = "GOQ_SQLITE3_CONNECTION_STRING"
	// _envStaticDir              = "SQUIRRELBYTE_STATIC_DIR"
	// _envAllowedHTTPMethods     = "SQUIRRELBYTE_ALLOWED_HTTP_METHODS"
	// _envAllowedHTTPPaths       = "SQUIRRELBYTE_ALLOWED_HTTP_PATHS"

	defaultServerPort = 9944
)

type config struct {
	ServerPort              int
	SQLite3ConnectionString string
	// StaticDir          string
	// AllowedHTTPMethods map[string]bool
	// AllowedHTTPPaths   map[string]bool
}

func loadConfig() (*config, error) {
	if fDotenv != nil && *fDotenv != "" {
		err := godotenv.Load(*fDotenv)
		if err != nil {
			return nil, err
		}
	}

	serverPort := defaultServerPort
	pstr := os.Getenv(_envServerPort)
	if pstr != "" {
		p, err := strconv.Atoi(pstr)
		if err != nil {
			return nil, err
		}
		serverPort = p
	}

	return &config{
		ServerPort:              serverPort,
		SQLite3ConnectionString: os.Getenv(_envSQLite3ConnectionString),
	}, nil
}

func main() {

	c, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	rootRouter := mux.NewRouter()

	apiHdl := apiserver.NewAPIHandler()
	apiRouter := rootRouter.PathPrefix("/api").Subrouter()
	apiserver.RegisterRouter(apiHdl, apiRouter, apiserver.GetErrorCode)

	webHdl := webserver.NewWebHandler()
	webserver.RegisterRouter(webHdl, rootRouter)

	addr := fmt.Sprintf(":%d", c.ServerPort)
	srv := &http.Server{
		Handler:      rootRouter,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("running on %s\n", addr)
	srv.ListenAndServe()
}
