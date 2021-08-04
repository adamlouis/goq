package webserver

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
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
	Job(w http.ResponseWriter, req *http.Request)
	Scheduler(w http.ResponseWriter, req *http.Request)
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

type SchedulerTmpl struct {
	*goqmodel.Scheduler
	InputStr string
}

type JobTmpl struct {
	*goqmodel.Job
	InputStr, OutputStr string
}

type pageData struct {
	Title        string
	JobStr       string
	SchedulerStr string
	Pivot        [][]string
	Jobs         []*JobTmpl
	Schedulers   []*SchedulerTmpl
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
	r.Handle("/jobs/{jobID}", http.HandlerFunc(wh.Job)).Methods(http.MethodGet)
	r.Handle("/schedulers/{schedulerID}", http.HandlerFunc(wh.Scheduler)).Methods(http.MethodGet)
}

func (w *whs) Home(wrt http.ResponseWriter, req *http.Request) {
	report, err := w.reporter.GetCountByNameByStatus(req.Context())
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}

	statusColumns := []job.JobStatus{
		job.JobStatusQueued,
		job.JobStatusClaimed,
		job.JobStatusSuccess,
		job.JobStatusError,
	}

	header := make([]string, 1+len(statusColumns))
	for i, s := range statusColumns {
		header[i+1] = string(s)
	}
	pivot := [][]string{header}

	names := make([]string, len(report))
	i := 0
	for name := range report {
		names[i] = name
		i += 1
	}
	sort.Strings(names)

	for _, name := range names {
		row := make([]string, 5)
		row[0] = name
		for i, s := range statusColumns {
			row[i+1] = fmt.Sprintf("%d", report[name][s])
		}
		pivot = append(pivot, row)
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

	schedulers := []*goqmodel.Scheduler{}
	rs, err := w.apiHandler.ListSchedulers(req.Context(), &goqmodel.ListSchedulersRequest{
		PageSize: 100,
	})
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	} else {
		schedulers = rs.Schedulers
	}

	newTemplate("home.go.html", []string{"templates/common.go.html", "templates/home.go.html"}).Execute(wrt, pageData{
		Title:      "GOQ",
		Jobs:       toJobTmpls(jobs),
		Pivot:      pivot,
		Schedulers: toSchedulerTmpls(schedulers),
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

	newTemplate("job.go.html", []string{"templates/common.go.html", "templates/job.go.html"}).Execute(wrt, pageData{
		Title:  "GOQ",
		JobStr: string(s),
	})
}

func (w *whs) Scheduler(wrt http.ResponseWriter, req *http.Request) {
	schedulerID := mux.Vars(req)["schedulerID"]

	s, err := w.apiHandler.GetScheduler(req.Context(), &goqmodel.GetSchedulerPathParams{
		SchedulerID: schedulerID,
	})

	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}

	sb, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}

	newTemplate("scheduler.go.html", []string{"templates/common.go.html", "templates/scheduler.go.html"}).Execute(wrt, pageData{
		Title:        "GOQ",
		SchedulerStr: string(sb),
	})
}

func toSchedulerTmpls(ss []*goqmodel.Scheduler) []*SchedulerTmpl {
	r := make([]*SchedulerTmpl, len(ss))
	for i := range ss {
		ib, _ := json.Marshal(ss[i].Input)

		r[i] = &SchedulerTmpl{
			Scheduler: ss[i],
			InputStr:  string(ib),
		}
	}
	return r
}

func toJobTmpls(js []*goqmodel.Job) []*JobTmpl {
	r := make([]*JobTmpl, len(js))
	for i := range js {
		ib, _ := json.Marshal(js[i].Input)
		ob, _ := json.Marshal(js[i].Output)

		r[i] = &JobTmpl{
			Job:       js[i],
			InputStr:  string(ib),
			OutputStr: string(ob),
		}
	}
	return r
}
