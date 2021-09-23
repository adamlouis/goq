package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (wh *webHandler) GetSchedulerCreate(w http.ResponseWriter, r *http.Request) {
	newTemplate("create-scheduler.go.html", []string{"templates/common.go.html", "templates/create-scheduler.go.html"}).Execute(w, pageData{
		Title: "GOQ",
	})
}
func (wh *webHandler) PostSchedulerCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		jsonlog.Log("error", err.Error())
	}

	s := &goqmodel.Scheduler{}
	err := json.Unmarshal([]byte(r.Form.Get("body")), s)
	if err != nil {
		jsonlog.Log("error", err.Error())
	} else {
		_, err = wh.apiHandler.PostScheduler(r.Context(), s)
		if err != nil {
			jsonlog.Log("error", err.Error())
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
