// GENERATED
// DO NOT EDIT
// GENERATOR: scripts/gencode/gencode.go
// ARGUMENTS: --component server --config ../../api/api.yml --package apiserver --out-dir ./ --out ./apiserver.gen.go --model-package github.com/adamlouis/goq/pkg/goqmodel
package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adamlouis/goq/pkg/goqmodel"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type HTTPHandler interface {
	ListJobs(w http.ResponseWriter, req *http.Request)
	GetJob(w http.ResponseWriter, req *http.Request)
	DeleteJob(w http.ResponseWriter, req *http.Request)
	QueueJob(w http.ResponseWriter, req *http.Request)
	ClaimSomeJob(w http.ResponseWriter, req *http.Request)
	ClaimJob(w http.ResponseWriter, req *http.Request)
	ReleaseJob(w http.ResponseWriter, req *http.Request)
	SetJobSuccess(w http.ResponseWriter, req *http.Request)
	SetJobError(w http.ResponseWriter, req *http.Request)
}
type APIHandler interface {
	ListJobs(ctx context.Context, queryParams *goqmodel.ListJobsQueryParams) (*goqmodel.ListJobsResponse, error)
	GetJob(ctx context.Context, pathParams *goqmodel.GetJobPathParams) (*goqmodel.Job, error)
	DeleteJob(ctx context.Context, pathParams *goqmodel.DeleteJobPathParams) error
	QueueJob(ctx context.Context, body *goqmodel.Job) (*goqmodel.Job, error)
	ClaimSomeJob(ctx context.Context, body *goqmodel.ClaimSomeJobRequest) (*goqmodel.Job, error)
	ClaimJob(ctx context.Context, pathParams *goqmodel.ClaimJobPathParams) (*goqmodel.Job, error)
	ReleaseJob(ctx context.Context, pathParams *goqmodel.ReleaseJobPathParams) (*goqmodel.Job, error)
	SetJobSuccess(ctx context.Context, pathParams *goqmodel.SetJobSuccessPathParams) (*goqmodel.Job, error)
	SetJobError(ctx context.Context, pathParams *goqmodel.SetJobErrorPathParams) (*goqmodel.Job, error)
}

func RegisterRouter(apiHandler APIHandler, r *mux.Router, c ErrorCoder) {
	h := apiHandlerToHTTPHandler(apiHandler, c)
	r.Handle("/jobs", http.HandlerFunc(h.ListJobs)).Methods(http.MethodGet)
	r.Handle("/jobs/{jobID}", http.HandlerFunc(h.GetJob)).Methods(http.MethodGet)
	r.Handle("/jobs/{jobID}", http.HandlerFunc(h.DeleteJob)).Methods(http.MethodDelete)
	r.Handle("/jobs:queue", http.HandlerFunc(h.QueueJob)).Methods(http.MethodPost)
	r.Handle("/jobs:claim", http.HandlerFunc(h.ClaimSomeJob)).Methods(http.MethodPost)
	r.Handle("/jobs/{jobID}:claim", http.HandlerFunc(h.ClaimJob)).Methods(http.MethodPost)
	r.Handle("/jobs/{jobID}:release", http.HandlerFunc(h.ReleaseJob)).Methods(http.MethodPost)
	r.Handle("/jobs/{jobID}:success", http.HandlerFunc(h.SetJobSuccess)).Methods(http.MethodPost)
	r.Handle("/jobs/{jobID}:error", http.HandlerFunc(h.SetJobError)).Methods(http.MethodPost)
}

func apiHandlerToHTTPHandler(apiHandler APIHandler, errorCoder ErrorCoder) HTTPHandler {
	return &httpHandler{
		apiHandler: apiHandler,
		errorCoder: errorCoder,
	}
}

type httpHandler struct {
	apiHandler APIHandler
	errorCoder ErrorCoder
}

type ErrorCoder func(e error) int

// sendError sends an error response
func (h *httpHandler) sendError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(h.errorCoder(err))
	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	e.Encode(&errorResponse{
		Message: err.Error(),
	})
}

func sendErrorWithCode(w http.ResponseWriter, code int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	e.Encode(&errorResponse{
		Message: err.Error(),
	})
}

// sendOK sends an success response
func sendOK(w http.ResponseWriter, body interface{}) {
	w.Header().Add("Content-Type", "application/json")
	code := http.StatusOK
	if body == nil {
		code = http.StatusNoContent
	}
	w.WriteHeader(code)
	if body != nil {
		e := json.NewEncoder(w)
		e.SetEscapeHTML(false)
		e.Encode(body)
	}
}

type errorResponse struct {
	Message string `json:"message"`
}

func (h *httpHandler) ListJobs(w http.ResponseWriter, req *http.Request) {
	whereQueryParam := req.URL.Query().Get("where")
	orderByQueryParam := req.URL.Query().Get("order_by")
	pageSizeQueryParam := 0
	if req.URL.Query().Get("page_size") != "" {
		q, err := strconv.Atoi(req.URL.Query().Get("page_size"))
		if err != nil {
			sendErrorWithCode(w, http.StatusBadRequest, err)
			return
		}
		pageSizeQueryParam = q
	}
	pageTokenQueryParam := req.URL.Query().Get("page_token")
	queryParams := goqmodel.ListJobsQueryParams{
		Where:     whereQueryParam,
		OrderBy:   orderByQueryParam,
		PageSize:  pageSizeQueryParam,
		PageToken: pageTokenQueryParam,
	}
	r, err := h.apiHandler.ListJobs(req.Context(), &queryParams)
	if err != nil {
		h.sendError(w, err)
		return
	}
	sendOK(w, r)
}
func (h *httpHandler) GetJob(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	jobID, ok := vars["jobID"]
	if !ok {
		sendErrorWithCode(w, http.StatusBadRequest, fmt.Errorf("invalid jobID path parameter"))
		return
	}
	pathParams := goqmodel.GetJobPathParams{
		JobID: jobID,
	}
	r, err := h.apiHandler.GetJob(req.Context(), &pathParams)
	if err != nil {
		h.sendError(w, err)
		return
	}
	sendOK(w, r)
}
func (h *httpHandler) DeleteJob(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	jobID, ok := vars["jobID"]
	if !ok {
		sendErrorWithCode(w, http.StatusBadRequest, fmt.Errorf("invalid jobID path parameter"))
		return
	}
	pathParams := goqmodel.DeleteJobPathParams{
		JobID: jobID,
	}
	err := h.apiHandler.DeleteJob(req.Context(), &pathParams)
	if err != nil {
		h.sendError(w, err)
		return
	}
	sendOK(w, struct{}{})
}
func (h *httpHandler) QueueJob(w http.ResponseWriter, req *http.Request) {
	var requestBody goqmodel.Job
	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		sendErrorWithCode(w, http.StatusBadRequest, err)
		return
	}
	r, err := h.apiHandler.QueueJob(req.Context(), &requestBody)
	if err != nil {
		h.sendError(w, err)
		return
	}
	sendOK(w, r)
}
func (h *httpHandler) ClaimSomeJob(w http.ResponseWriter, req *http.Request) {
	var requestBody goqmodel.ClaimSomeJobRequest
	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		sendErrorWithCode(w, http.StatusBadRequest, err)
		return
	}
	r, err := h.apiHandler.ClaimSomeJob(req.Context(), &requestBody)
	if err != nil {
		h.sendError(w, err)
		return
	}
	sendOK(w, r)
}
func (h *httpHandler) ClaimJob(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	jobID, ok := vars["jobID"]
	if !ok {
		sendErrorWithCode(w, http.StatusBadRequest, fmt.Errorf("invalid jobID path parameter"))
		return
	}
	pathParams := goqmodel.ClaimJobPathParams{
		JobID: jobID,
	}
	r, err := h.apiHandler.ClaimJob(req.Context(), &pathParams)
	if err != nil {
		h.sendError(w, err)
		return
	}
	sendOK(w, r)
}
func (h *httpHandler) ReleaseJob(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	jobID, ok := vars["jobID"]
	if !ok {
		sendErrorWithCode(w, http.StatusBadRequest, fmt.Errorf("invalid jobID path parameter"))
		return
	}
	pathParams := goqmodel.ReleaseJobPathParams{
		JobID: jobID,
	}
	r, err := h.apiHandler.ReleaseJob(req.Context(), &pathParams)
	if err != nil {
		h.sendError(w, err)
		return
	}
	sendOK(w, r)
}
func (h *httpHandler) SetJobSuccess(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	jobID, ok := vars["jobID"]
	if !ok {
		sendErrorWithCode(w, http.StatusBadRequest, fmt.Errorf("invalid jobID path parameter"))
		return
	}
	pathParams := goqmodel.SetJobSuccessPathParams{
		JobID: jobID,
	}
	r, err := h.apiHandler.SetJobSuccess(req.Context(), &pathParams)
	if err != nil {
		h.sendError(w, err)
		return
	}
	sendOK(w, r)
}
func (h *httpHandler) SetJobError(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	jobID, ok := vars["jobID"]
	if !ok {
		sendErrorWithCode(w, http.StatusBadRequest, fmt.Errorf("invalid jobID path parameter"))
		return
	}
	pathParams := goqmodel.SetJobErrorPathParams{
		JobID: jobID,
	}
	r, err := h.apiHandler.SetJobError(req.Context(), &pathParams)
	if err != nil {
		h.sendError(w, err)
		return
	}
	sendOK(w, r)
}
