package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (wh *webHandler) GetJobCreate(w http.ResponseWriter, r *http.Request) {
	newTemplate("create-job.go.html", []string{"templates/common.go.html", "templates/create-job.go.html"}).Execute(w, pageData{
		Title: "GOQ",
	})

}
func (wh *webHandler) PostJobCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		jsonlog.Log("error", err.Error())
	}

	j := &goqmodel.Job{}
	err := json.Unmarshal([]byte(r.Form.Get("body")), j)
	if err != nil {
		jsonlog.Log("error", err.Error())
	} else {
		_, err = wh.apiHandler.QueueJob(r.Context(), j)
		if err != nil {
			jsonlog.Log("error", err.Error())
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
