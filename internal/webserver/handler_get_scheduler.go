package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/pkg/goqmodel"
	"github.com/gorilla/mux"
)

func (wh *webHandler) GetScheduler(w http.ResponseWriter, r *http.Request) {
	schedulerID := mux.Vars(r)["schedulerID"]

	s, err := wh.apiHandler.GetScheduler(r.Context(), &goqmodel.GetSchedulerPathParams{
		SchedulerID: schedulerID,
	})

	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}

	sb, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		jsonlog.Log("error", err) // TODO: handle error
	}

	newTemplate("scheduler.go.html", []string{"templates/common.go.html", "templates/scheduler.go.html"}).Execute(w, pageData{
		Title:        "GOQ",
		SchedulerStr: string(sb),
	})
}
