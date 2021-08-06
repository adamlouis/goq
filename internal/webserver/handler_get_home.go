package webserver

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/adamlouis/goq/internal/job"
	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (wh *webHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "web-session")
	username := fmt.Sprintf("%v", session.Values["username"])

	report, err := wh.reporter.GetCountByNameByStatus(r.Context())
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
	result, err := wh.apiHandler.ListJobs(r.Context(), &goqmodel.ListJobsQueryParams{
		PageSize: 100,
	})
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	} else {
		jobs = result.Jobs
	}

	schedulers := []*goqmodel.Scheduler{}
	rs, err := wh.apiHandler.ListSchedulers(r.Context(), &goqmodel.ListSchedulersRequest{
		PageSize: 100,
	})
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	} else {
		schedulers = rs.Schedulers
	}

	newTemplate("home.go.html", []string{"templates/common.go.html", "templates/home.go.html"}).Execute(w, pageData{
		Username:   username,
		Title:      "GOQ",
		Jobs:       toJobTmpls(jobs),
		Pivot:      pivot,
		Schedulers: toSchedulerTmpls(schedulers),
	})
}
