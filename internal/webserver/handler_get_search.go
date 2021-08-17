package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (wh *webHandler) GetSearch(w http.ResponseWriter, r *http.Request) {
	username := ""
	if p, _ := wh.sessionManger.Get(w, r); p != nil {
		username = p.Username
	}

	req := &goqmodel.SearchJobsRequest{
		Where:    map[string]interface{}{},
		OrderBy:  []interface{}{},
		PageSize: 100,
	}

	jq := r.URL.Query().Get("jq")
	parsed := &goqmodel.SearchJobsRequest{}
	if err := json.Unmarshal([]byte(jq), parsed); err == nil {
		req = parsed
	}

	req.PageSize = 100

	jobs := []*goqmodel.Job{}
	result, err := wh.apiHandler.SearchJobs(r.Context(), req)
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	} else {
		jobs = result.Jobs
	}

	newTemplate("search.go.html", []string{"templates/common.go.html", "templates/search.go.html"}).Execute(w, pageData{
		Username: username,
		Title:    "GOQ",
		JQ:       jq,
		Jobs:     toJobTmpls(jobs),
	})
}
