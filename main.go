package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adamlouis/goq/internal/apiserver"
	"github.com/adamlouis/goq/internal/job/jobsqlite3"
	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/internal/pkg/sqlite3util"
	"github.com/adamlouis/goq/internal/scheduler/schedulersqlite3"
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

// TODO: revisit config-based init
const (
	_envServerPort = "GOQ_SERVER_PORT"
	// _envSQLite3ConnectionString = "GOQ_SQLITE3_CONNECTION_STRING"
	// _envStaticDir              = "SQUIRRELBYTE_STATIC_DIR"
	// _envAllowedHTTPMethods     = "SQUIRRELBYTE_ALLOWED_HTTP_METHODS"
	// _envAllowedHTTPPaths       = "SQUIRRELBYTE_ALLOWED_HTTP_PATHS"

	defaultServerPort = 9944
)

type config struct {
	ServerPort int
	// SQLite3ConnectionString string
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
		ServerPort: serverPort,
		// SQLite3ConnectionString: os.Getenv(_envSQLite3ConnectionString),
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

	rootRouter := mux.NewRouter()
	rootRouter.Use(loggerMiddleware, webSessionAuthorizationProcessor, bearerAuthorizationProcessor, aclMiddleware)

	apiHdl := apiserver.NewAPIHandler(ctx, jobDB, schedulerDB)
	apiRouter := rootRouter.PathPrefix("/api").Subrouter()
	apiserver.RegisterRouter(apiHdl, apiRouter, apiserver.GetErrorCode)

	webHdl := webserver.NewWebHandler(apiHdl, jobsqlite3.NewJobReporter(jobDB))
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

// TODO: pass authz/authn pkg, test, etc
var (
	key   = []byte(os.Getenv("GOQ_SESSION_KEY"))
	store = sessions.NewFilesystemStore(os.Getenv("GOQ_SESSION_STORE_PATH"), key)
)

type AuthTypeName string

const AuthType AuthTypeName = "AUTH_TYPE"

type AuthTypeValue string

const AuthTypeWeb AuthTypeValue = "WEB_SESSION"
const AuthTypeAPI AuthTypeValue = "BEARER"

func webSessionAuthorizationProcessor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "web-session")
		if auth, ok := session.Values["authenticated"].(bool); ok && auth {
			ctx := context.WithValue(r.Context(), AuthType, AuthTypeWeb)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func bearerAuthorizationProcessor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		az := r.Header["Authorization"]
		if len(az) == 1 {
			apiKey := os.Getenv("GOQ_API_KEY")
			if len(apiKey) > 0 {
				if az[0] == fmt.Sprintf("Bearer %s", apiKey) {
					ctx := context.WithValue(r.Context(), AuthType, AuthTypeAPI)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func aclMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow := false
		val, ok := r.Context().Value(AuthType).(AuthTypeValue)

		if r.URL.Path == "/login" || r.URL.Path == "/logout" {
			allow = true
		}

		if strings.HasPrefix(r.URL.Path, "/api") && ok && val == AuthTypeAPI {
			allow = true
		}

		if !strings.HasPrefix(r.URL.Path, "/api") && ok && val == AuthTypeWeb {
			allow = true
		}

		if allow {
			next.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api") {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, `forbidden - <a href="/login">login</a>`)
	})
}
