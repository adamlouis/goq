package webserver

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/adamlouis/goq/internal/job"
	"github.com/adamlouis/goq/internal/job/jobsqlite3"
	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (wh *webHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	username := ""
	if p, _ := wh.sessionManger.Get(w, r); p != nil {
		username = p.Username
	}

	pivot, err := wh.getJobStatusTable(r.Context())
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}

	jns := []string{}
	for _, r := range pivot {
		for j, c := range r {
			if j == 0 {
				jns = append(jns, c)
			}

		}
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
		Pivot:      pivot,
		Jobs:       toJobTmpls(jobs),
		JobNames:   jns,
		Schedulers: toSchedulerTmpls(schedulers),
	})
}

func (wh *webHandler) getJobStatusTable(ctx context.Context) ([][]string, error) {
	tx, err := wh.jobDB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	report, err := jobsqlite3.NewReporter(tx).GetCountByNameByStatus(ctx)
	if err != nil {
		return nil, err
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
		row := make([]string, len(statusColumns)+1)
		row[0] = name
		for i, s := range statusColumns {
			row[i+1] = fmt.Sprintf("%d", report[name][s])
		}
		pivot = append(pivot, row)
	}

	return pivot, nil
}
