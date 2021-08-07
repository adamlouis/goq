package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/adamlouis/goq/internal/apiserver"
	"github.com/adamlouis/goq/internal/auth"
	"github.com/adamlouis/goq/internal/job/jobsqlite3"
	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/internal/pkg/sqlite3util"
	"github.com/adamlouis/goq/internal/scheduler/schedulersqlite3"
	"github.com/adamlouis/goq/internal/session"
	"github.com/adamlouis/goq/internal/webserver"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	fDotenv = flag.String("dotenv", "", "a .env file from which to read environment variables. useful for local development.")
)

const (
	_envServerPort       = "GOQ_SERVER_PORT"
	_envRootUsername     = "GOQ_ROOT_USERNAME"
	_envRootPassword     = "GOQ_ROOT_PASSWORD"
	_envSessionStorePath = "GOQ_SESSION_STORE_PATH"
	_envSessionKey       = "GOQ_SESSION_KEY"
	_envAPIKey           = "GOQ_API_KEY"

	defaultServerPort = 9944
)

type config struct {
	ServerPort       int
	RootUsername     string
	RootPassword     string
	SessionStorePath string
	SessionKey       string
	APIKey           string
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
		ServerPort:       serverPort,
		RootUsername:     os.Getenv(_envRootUsername),
		RootPassword:     os.Getenv(_envRootPassword),
		SessionStorePath: os.Getenv(_envSessionStorePath),
		SessionKey:       os.Getenv(_envSessionKey),
		APIKey:           os.Getenv(_envAPIKey),
	}, nil
}

func newDB(c *config, path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // TODO: use RW lock or WAL rather than 1 max conn
	return db, nil
}

func main() {
	ctx := context.Background()

	c, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if len(c.APIKey) == 0 {
		log.Fatalf("API key required")
	}
	if len(c.RootPassword) == 0 {
		log.Fatalf("root password required")
	}
	if len(c.SessionKey) == 0 {
		log.Fatalf("session key required")
	}

	jobDB, err := newDB(c, "./db/job.db")
	if err != nil {
		log.Fatal(err)
	}
	defer jobDB.Close()

	if err := sqlite3util.NewMigrator(jobDB, jobsqlite3.MigrationFS).Up(); err != nil {
		log.Fatal(err)
	}

	schedulerDB, err := newDB(c, "./db/scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer schedulerDB.Close()

	if err := sqlite3util.NewMigrator(schedulerDB, schedulersqlite3.MigrationFS).Up(); err != nil {
		log.Fatal(err)
	}

	sessionManager := session.NewFSManager("web-session", sessions.NewFilesystemStore(c.SessionStorePath, []byte(c.SessionKey)))
	upChecker := auth.NewConstUPChecker(c.RootUsername, c.RootPassword)
	apiKeyChecker := auth.NewConstKChecker(c.APIKey)

	rootRouter := mux.NewRouter()
	rootRouter.Use(loggerMiddleware)
	rootRouter.Use(auth.GetMiddleware(sessionManager, apiKeyChecker)...)

	apiHdl := apiserver.NewAPIHandler(ctx, jobDB, schedulerDB)
	apiRouter := rootRouter.PathPrefix("/api").Subrouter()
	apiserver.RegisterRouter(apiHdl, apiRouter, apiserver.GetErrorCode)

	webHdl := webserver.NewWebHandler(apiHdl, jobsqlite3.NewJobReporter(jobDB), sessionManager, upChecker)
	webserver.RegisterRouter(webHdl, rootRouter)

	addr := fmt.Sprintf(":%d", c.ServerPort)
	srv := &http.Server{
		Handler:      rootRouter,
		Addr:         addr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	jsonlog.Log("type", "SERVER_STARTED", "port", c.ServerPort)
	log.Fatal(srv.ListenAndServe())

}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next.ServeHTTP(w, r)
		// TODO: do w/ tracing / opentelemetry, start a span, pass down context, etc
		// TODO: produce just 1 structured event per req w/ all metadata
		jsonlog.Log(
			"name", fmt.Sprintf("%s:%s", r.Method, r.URL.Path),
			"type", "REQUEST",
			"method", r.Method,
			"duration_ms", time.Since(now)/time.Millisecond,
			"path", r.URL.Path,
			"time", time.Now().Format(time.RFC3339),
		)
	})
}
