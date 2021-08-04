package webserver

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/adamlouis/goq/internal/apiserver"
	"github.com/adamlouis/goq/internal/job"
	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/pkg/goqmodel"
	"github.com/gorilla/mux"
)

//go:embed templates
var templatesFS embed.FS

type WebHandler interface {
	Home(w http.ResponseWriter, req *http.Request)
	Jobs(w http.ResponseWriter, req *http.Request)
	Job(w http.ResponseWriter, req *http.Request)
}

func NewWebHandler(apiHandler apiserver.APIHandler, reporter job.Reporter) WebHandler {
	return &whs{
		apiHandler: apiHandler,
		reporter:   reporter,
	}
}

type whs struct {
	apiHandler apiserver.APIHandler
	reporter   job.Reporter
}

type pageData struct {
	Title            string
	JobStr           string
	Job              *goqmodel.Job
	Jobs             []*goqmodel.Job
	JobCountByName   map[string]int64
	JobCountByStatus map[job.JobStatus]int64
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
	r.Handle("/jobs/{jobID}", http.HandlerFunc(wh.Job)).Methods(http.MethodGet)
}

func (w *whs) Home(wrt http.ResponseWriter, req *http.Request) {
	w.Jobs(wrt, req)
}

func (w *whs) Jobs(wrt http.ResponseWriter, req *http.Request) {
	bn, err := w.reporter.GetCountByName(req.Context())
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}
	bs, err := w.reporter.GetCountByStatus(req.Context())
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}

	jobs := []*goqmodel.Job{}
	r, err := w.apiHandler.ListJobs(req.Context(), &goqmodel.ListJobsQueryParams{
		PageSize: 100,
	})

	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	} else {
		jobs = r.Jobs
	}

	newTemplate("jobs.go.html", []string{"templates/jobs.go.html"}).Execute(wrt, pageData{
		Title:            "GOQ",
		Jobs:             jobs,
		JobCountByName:   bn,
		JobCountByStatus: bs,
	})
}

func (w *whs) Job(wrt http.ResponseWriter, req *http.Request) {
	jobID := mux.Vars(req)["jobID"]

	j, err := w.apiHandler.GetJob(req.Context(), &goqmodel.GetJobPathParams{
		JobID: jobID,
	})

	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}

	s, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}

	newTemplate("job.go.html", []string{"templates/job.go.html"}).Execute(wrt, pageData{
		Title:  "GOQ",
		Job:    j,
		JobStr: string(s),
	})
}
