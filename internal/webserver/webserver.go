package webserver

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/adamlouis/goq/internal/apiserver"
	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/pkg/goqmodel"
	"github.com/gorilla/mux"
)

//go:embed templates
var templatesFS embed.FS

type WebHandler interface {
	Home(w http.ResponseWriter, req *http.Request)
	Jobs(w http.ResponseWriter, req *http.Request)
}

func NewWebHandler(apiHandler apiserver.APIHandler) WebHandler {
	return &whs{
		apiHandler: apiHandler,
	}
}

type whs struct {
	apiHandler apiserver.APIHandler
}

type pageData struct {
	Title string
	Jobs  []*goqmodel.Job
}

func newTemplate(name string, patterns []string) *template.Template {
	if os.Getenv("GOQ_MODE") == "LOCAL" {
		resolved := make([]string, len(patterns))
		for i := range patterns {
			resolved[i] = fmt.Sprintf("internal/webserver/%s", patterns[i])
		}
		return template.Must(template.New(name).ParseFiles(resolved...))
	}
	return template.Must(template.New(name).ParseFS(templatesFS, patterns...))
}

func RegisterRouter(wh WebHandler, r *mux.Router) {
	r.Handle("/", http.HandlerFunc(wh.Home)).Methods(http.MethodGet)
	r.Handle("/jobs", http.HandlerFunc(wh.Jobs)).Methods(http.MethodGet)
}

func (w *whs) Home(wrt http.ResponseWriter, req *http.Request) {
	newTemplate("home.go.html", []string{"templates/home.go.html"}).Execute(wrt, pageData{
		Title: "GOQ",
	})
}
func (w *whs) Jobs(wrt http.ResponseWriter, req *http.Request) {
	jobs := []*goqmodel.Job{}
	r, err := w.apiHandler.ListJobs(req.Context(), &goqmodel.ListJobsQueryParams{
		PageSize: 100,
	})

	if err != nil {
		// TODO: handle error
		jsonlog.Log("error", err)
	} else {
		jobs = r.Jobs
	}

	newTemplate("jobs.go.html", []string{"templates/jobs.go.html"}).Execute(wrt, pageData{
		Title: "GOQ",
		Jobs:  jobs,
	})
}
