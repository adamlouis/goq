package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/pkg/goqmodel"
	"github.com/gorilla/mux"
)

func (wh *webHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	jobID := mux.Vars(r)["jobID"]

	j, err := wh.apiHandler.GetJob(r.Context(), &goqmodel.GetJobPathParams{
		JobID: jobID,
	})

	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}

	s, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}

	newTemplate("job.go.html", []string{"templates/common.go.html", "templates/job.go.html"}).Execute(w, pageData{
		Title:  "GOQ",
		JobStr: string(s),
	})
}
