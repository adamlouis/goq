package webserver

import (
	"net/http"

	"github.com/adamlouis/goq/internal/apiserver"
	"github.com/adamlouis/goq/internal/job"
	"github.com/gorilla/mux"
)

type WebHandler interface {
	GetHome(w http.ResponseWriter, req *http.Request)
	GetJob(w http.ResponseWriter, req *http.Request)
	GetScheduler(w http.ResponseWriter, req *http.Request)
	GetLogin(w http.ResponseWriter, req *http.Request)
	PostLogin(w http.ResponseWriter, req *http.Request)
	GetLogout(w http.ResponseWriter, req *http.Request)
}

func NewWebHandler(apiHandler apiserver.APIHandler, reporter job.Reporter) WebHandler {
	return &webHandler{
		apiHandler: apiHandler,
		reporter:   reporter,
	}
}

type webHandler struct {
	apiHandler apiserver.APIHandler
	reporter   job.Reporter
}

func RegisterRouter(wh WebHandler, r *mux.Router) {
	r.Handle("/", http.HandlerFunc(wh.GetHome)).Methods(http.MethodGet)
	r.Handle("/jobs/{jobID}", http.HandlerFunc(wh.GetJob)).Methods(http.MethodGet)
	r.Handle("/schedulers/{schedulerID}", http.HandlerFunc(wh.GetScheduler)).Methods(http.MethodGet)
	r.Handle("/login", http.HandlerFunc(wh.GetLogin)).Methods(http.MethodGet)
	r.Handle("/login", http.HandlerFunc(wh.PostLogin)).Methods(http.MethodPost)
	r.Handle("/logout", http.HandlerFunc(wh.GetLogout)).Methods(http.MethodGet)
}
