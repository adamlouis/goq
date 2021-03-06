package webserver

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/adamlouis/goq/internal/apiserver"
	"github.com/adamlouis/goq/internal/auth"
	"github.com/adamlouis/goq/internal/session"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

//go:embed templates
var templatesFS embed.FS

type WebHandler interface {
	GetHome(w http.ResponseWriter, req *http.Request)
	GetJob(w http.ResponseWriter, req *http.Request)
	GetScheduler(w http.ResponseWriter, req *http.Request)
	GetSearch(w http.ResponseWriter, req *http.Request)
	GetLogin(w http.ResponseWriter, req *http.Request)
	PostLogin(w http.ResponseWriter, req *http.Request)
	GetLogout(w http.ResponseWriter, req *http.Request)
	GetJobCreate(w http.ResponseWriter, req *http.Request)
	PostJobCreate(w http.ResponseWriter, req *http.Request)
	GetSchedulerCreate(w http.ResponseWriter, req *http.Request)
	PostSchedulerCreate(w http.ResponseWriter, req *http.Request)
}

func NewWebHandler(
	apiHandler apiserver.APIHandler,
	jobDB *sqlx.DB,
	sessionManger session.Manager,
	checker auth.UPChecker,
) WebHandler {
	return &webHandler{
		apiHandler:    apiHandler,
		jobDB:         jobDB,
		sessionManger: sessionManger,
		checker:       checker,
	}
}

type webHandler struct {
	apiHandler    apiserver.APIHandler
	jobDB         *sqlx.DB
	sessionManger session.Manager
	checker       auth.UPChecker
}

func RegisterRouter(wh WebHandler, r *mux.Router) {
	r.Handle("/", http.HandlerFunc(wh.GetHome)).Methods(http.MethodGet)
	r.Handle("/jobs/{jobID}", http.HandlerFunc(wh.GetJob)).Methods(http.MethodGet)
	r.Handle("/schedulers/{schedulerID}", http.HandlerFunc(wh.GetScheduler)).Methods(http.MethodGet)
	r.Handle("/jobs:create", http.HandlerFunc(wh.GetJobCreate)).Methods(http.MethodGet)
	r.Handle("/jobs:create", http.HandlerFunc(wh.PostJobCreate)).Methods(http.MethodPost)
	r.Handle("/schedulers:create", http.HandlerFunc(wh.GetSchedulerCreate)).Methods(http.MethodGet)
	r.Handle("/schedulers:create", http.HandlerFunc(wh.PostSchedulerCreate)).Methods(http.MethodPost)
	r.Handle("/search", http.HandlerFunc(wh.GetSearch)).Methods(http.MethodGet)
	r.Handle("/login", http.HandlerFunc(wh.GetLogin)).Methods(http.MethodGet)
	r.Handle("/login", http.HandlerFunc(wh.PostLogin)).Methods(http.MethodPost)
	r.Handle("/logout", http.HandlerFunc(wh.GetLogout)).Methods(http.MethodGet)
}

var tmplFuncs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
}

// in development, load templates from local filesystem
// in production, use embeded FS
//
// in development, this allows template files to update without re-building
// in production, this allows for a self-contained executable
func newTemplate(name string, patterns []string) *template.Template {
	if os.Getenv("GOQ_MODE") == "DEVELOPMENT" {
		resolved := make([]string, len(patterns))
		for i := range patterns {
			resolved[i] = fmt.Sprintf("internal/webserver/%s", patterns[i])
		}
		return template.Must(template.New(name).Funcs(tmplFuncs).ParseFiles(resolved...))
	}
	return template.Must(template.New(name).Funcs(tmplFuncs).ParseFS(templatesFS, patterns...))
}
