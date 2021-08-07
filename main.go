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
	_envServerPort            = "GOQ_SERVER_PORT"
	_envRootUsername          = "GOQ_ROOT_USERNAME"
	_envRootPassword          = "GOQ_ROOT_PASSWORD"
	_envSessionStorePath      = "GOQ_SESSION_STORE_PATH"
	_envSessionKey            = "GOQ_SESSION_KEY"
	_envAPIKey                = "GOQ_API_KEY"
	_envSQLiteJobDBPath       = "GOQ_SQLITE_JOB_DB_PATH"
	_envSQLiteSchedulerDBPath = "GOQ_SQLITE_SCHEDULER_DB_PATH"

	defaultServerPort            = 9944
	defaultSQLiteJobDBPath       = "./db/job.db"
	defaultSQLiteSchedulerDBPath = "./db/scheduler.db"
	webSessionCookieName         = "web-session"
)

type config struct {
	ServerPort            int
	RootUsername          string
	RootPassword          string
	SessionStorePath      string
	SessionKey            string
	APIKey                string
	SQLiteJobDBPath       string
	SQLiteSchedulerDBPath string
}

type dbs struct {
	jobDB       *sqlx.DB
	schedulerDB *sqlx.DB
}

func main() {
	ctx := context.Background()

	// load config from environment
	c, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if len(c.APIKey) == 0 {
		log.Fatalf("API key required")
	}
	if len(c.RootUsername) == 0 {
		log.Fatalf("root username required")
	}
	if len(c.RootPassword) == 0 {
		log.Fatalf("root password required")
	}
	if len(c.SessionKey) == 0 {
		log.Fatalf("session key required")
	}

	// set up dbs
	dbs, err := initDBS(c)
	if err != nil {
		log.Fatal(err)
	}
	defer dbs.jobDB.Close()
	defer dbs.schedulerDB.Close()

	// set up auth
	sessionManager := session.NewFSManager(webSessionCookieName, sessions.NewFilesystemStore(c.SessionStorePath, []byte(c.SessionKey)))
	// TODO: accept bcyrpt hash & compare that rather than plain passwd + api key. i am in a national forest and have no internet to get bcrypt
	upChecker := auth.NewConstUPChecker(c.RootUsername, c.RootPassword)
	apiKeyChecker := auth.NewConstKChecker(c.APIKey)

	// set up routers
	rootRouter := mux.NewRouter()
	apiRouter := rootRouter.PathPrefix("/api").Subrouter()

	// add middleware
	rootRouter.Use(loggerMiddleware)
	rootRouter.Use(auth.GetMiddleware(sessionManager, apiKeyChecker)...)

	// register api handlers
	apiHdl := apiserver.NewAPIHandler(ctx, dbs.jobDB, dbs.schedulerDB)
	apiserver.RegisterRouter(apiHdl, apiRouter, apiserver.GetErrorCode)

	// register web handlers
	webHdl := webserver.NewWebHandler(apiHdl, dbs.jobDB, sessionManager, upChecker)
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

// load config from the environment
func loadConfig() (*config, error) {
	if fDotenv != nil && *fDotenv != "" {
		err := godotenv.Load(*fDotenv)
		if err != nil {
			return nil, err
		}
	}

	return &config{
		ServerPort:            stoi(os.Getenv(_envServerPort), defaultServerPort),
		RootUsername:          os.Getenv(_envRootUsername),
		RootPassword:          os.Getenv(_envRootPassword),
		SessionStorePath:      os.Getenv(_envSessionStorePath),
		SessionKey:            os.Getenv(_envSessionKey),
		APIKey:                os.Getenv(_envAPIKey),
		SQLiteJobDBPath:       sdefault(os.Getenv(_envSQLiteJobDBPath), defaultSQLiteJobDBPath),
		SQLiteSchedulerDBPath: sdefault(os.Getenv(_envSQLiteSchedulerDBPath), defaultSQLiteSchedulerDBPath),
	}, nil
}

// init the dbs
func initDBS(c *config) (*dbs, error) {
	// open dbs
	jobDB, err := sqlx.Open("sqlite3", "./db/job.db")
	if err != nil {
		return nil, err
	}
	jobDB.SetMaxOpenConns(1) // TODO: use RW lock or WAL rather than 1 max conn

	schedulerDB, err := sqlx.Open("sqlite3", "./db/scheduler.db")
	if err != nil {
		return nil, err
	}
	schedulerDB.SetMaxOpenConns(1) // TODO: use RW lock or WAL rather than 1 max conn

	// run migrations
	if err := sqlite3util.NewMigrator(jobDB, jobsqlite3.MigrationFS).Up(); err != nil {
		return nil, err
	}

	if err := sqlite3util.NewMigrator(schedulerDB, schedulersqlite3.MigrationFS).Up(); err != nil {
		return nil, err
	}

	return &dbs{
		jobDB:       jobDB,
		schedulerDB: schedulerDB,
	}, nil
}

func sdefault(s1, s2 string) string {
	if s1 != "" {
		return s1
	}
	return s2
}

// return int representation of s1
// if s1 is not an int, return def
func stoi(s1 string, def int) int {
	i, err := strconv.Atoi(s1)
	if err == nil {
		return i
	}
	return def
}
