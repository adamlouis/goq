// GENERATED
// DO NOT EDIT
// GENERATOR: scripts/gencode/gencode.go
// ARGUMENTS: --component client --config ../../api/api.yml --package goqclient --out ./client.gen.go --model-package github.com/adamlouis/goq/pkg/goqmodel
package goqclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/adamlouis/goq/pkg/goqmodel"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Client interface {
	ListJobs(ctx context.Context, queryParams *goqmodel.ListJobsQueryParams) (*goqmodel.ListJobsResponse, int, error)
	GetJob(ctx context.Context, pathParams *goqmodel.GetJobPathParams) (*goqmodel.Job, int, error)
	DeleteJob(ctx context.Context, pathParams *goqmodel.DeleteJobPathParams) (int, error)
	QueueJob(ctx context.Context, body *goqmodel.Job) (*goqmodel.Job, int, error)
	ClaimSomeJob(ctx context.Context, body *goqmodel.ClaimSomeJobRequest) (*goqmodel.Job, int, error)
	ClaimJob(ctx context.Context, pathParams *goqmodel.ClaimJobPathParams) (*goqmodel.Job, int, error)
	ReleaseJob(ctx context.Context, pathParams *goqmodel.ReleaseJobPathParams) (*goqmodel.Job, int, error)
	SetJobSuccess(ctx context.Context, pathParams *goqmodel.SetJobSuccessPathParams, body *goqmodel.Job) (*goqmodel.Job, int, error)
	SetJobError(ctx context.Context, pathParams *goqmodel.SetJobErrorPathParams, body *goqmodel.Job) (*goqmodel.Job, int, error)
	ListSchedulers(ctx context.Context, queryParams *goqmodel.ListSchedulersRequest) (*goqmodel.ListSchedulersResponse, int, error)
	PostScheduler(ctx context.Context, body *goqmodel.Scheduler) (*goqmodel.Scheduler, int, error)
	GetScheduler(ctx context.Context, pathParams *goqmodel.GetSchedulerPathParams) (*goqmodel.Scheduler, int, error)
	PutScheduler(ctx context.Context, pathParams *goqmodel.PutSchedulerPathParams, body *goqmodel.Scheduler) (*goqmodel.Scheduler, int, error)
	DeleteScheduler(ctx context.Context, pathParams *goqmodel.DeleteSchedulerPathParams) (int, error)
}

func NewHTTPClient(baseURL string) Client {
	return &client{
		baseURL: baseURL,
	}
}

type client struct {
	baseURL string
}

func (c *client) ListJobs(ctx context.Context, queryParams *goqmodel.ListJobsQueryParams) (*goqmodel.ListJobsResponse, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/jobs", c.baseURL))
	if err != nil {
		return nil, -1, err
	}
	u.Query().Add("where", queryParams.Where)
	u.Query().Add("order_by", queryParams.OrderBy)
	u.Query().Add("page_size", strconv.Itoa(queryParams.PageSize))
	u.Query().Add("page_token", queryParams.PageToken)
	var requestBody io.Reader
	req, err := http.NewRequest(http.MethodGet, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.ListJobsResponse{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) GetJob(ctx context.Context, pathParams *goqmodel.GetJobPathParams) (*goqmodel.Job, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/jobs/%v", c.baseURL, pathParams.JobID))
	if err != nil {
		return nil, -1, err
	}
	var requestBody io.Reader
	req, err := http.NewRequest(http.MethodGet, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.Job{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) DeleteJob(ctx context.Context, pathParams *goqmodel.DeleteJobPathParams) (int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/jobs/%v", c.baseURL, pathParams.JobID))
	if err != nil {
		return -1, err
	}
	var requestBody io.Reader
	req, err := http.NewRequest(http.MethodDelete, u.String(), requestBody)
	if err != nil {
		return -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	return resp.StatusCode, nil
}
func (c *client) QueueJob(ctx context.Context, body *goqmodel.Job) (*goqmodel.Job, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/jobs:queue", c.baseURL))
	if err != nil {
		return nil, -1, err
	}
	var requestBody io.Reader
	if jsonBytes, err := json.Marshal(body); err != nil {
		return nil, -1, err
	} else {
		requestBody = bytes.NewBuffer(jsonBytes)
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.Job{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) ClaimSomeJob(ctx context.Context, body *goqmodel.ClaimSomeJobRequest) (*goqmodel.Job, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/jobs:claim", c.baseURL))
	if err != nil {
		return nil, -1, err
	}
	var requestBody io.Reader
	if jsonBytes, err := json.Marshal(body); err != nil {
		return nil, -1, err
	} else {
		requestBody = bytes.NewBuffer(jsonBytes)
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.Job{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) ClaimJob(ctx context.Context, pathParams *goqmodel.ClaimJobPathParams) (*goqmodel.Job, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/jobs/%v:claim", c.baseURL, pathParams.JobID))
	if err != nil {
		return nil, -1, err
	}
	var requestBody io.Reader
	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.Job{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) ReleaseJob(ctx context.Context, pathParams *goqmodel.ReleaseJobPathParams) (*goqmodel.Job, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/jobs/%v:release", c.baseURL, pathParams.JobID))
	if err != nil {
		return nil, -1, err
	}
	var requestBody io.Reader
	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.Job{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) SetJobSuccess(ctx context.Context, pathParams *goqmodel.SetJobSuccessPathParams, body *goqmodel.Job) (*goqmodel.Job, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/jobs/%v:success", c.baseURL, pathParams.JobID))
	if err != nil {
		return nil, -1, err
	}
	var requestBody io.Reader
	if jsonBytes, err := json.Marshal(body); err != nil {
		return nil, -1, err
	} else {
		requestBody = bytes.NewBuffer(jsonBytes)
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.Job{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) SetJobError(ctx context.Context, pathParams *goqmodel.SetJobErrorPathParams, body *goqmodel.Job) (*goqmodel.Job, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/jobs/%v:error", c.baseURL, pathParams.JobID))
	if err != nil {
		return nil, -1, err
	}
	var requestBody io.Reader
	if jsonBytes, err := json.Marshal(body); err != nil {
		return nil, -1, err
	} else {
		requestBody = bytes.NewBuffer(jsonBytes)
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.Job{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) ListSchedulers(ctx context.Context, queryParams *goqmodel.ListSchedulersRequest) (*goqmodel.ListSchedulersResponse, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/schedulers", c.baseURL))
	if err != nil {
		return nil, -1, err
	}
	u.Query().Add("page_token", queryParams.PageToken)
	u.Query().Add("page_size", strconv.Itoa(queryParams.PageSize))
	var requestBody io.Reader
	req, err := http.NewRequest(http.MethodGet, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.ListSchedulersResponse{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) PostScheduler(ctx context.Context, body *goqmodel.Scheduler) (*goqmodel.Scheduler, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/schedulers", c.baseURL))
	if err != nil {
		return nil, -1, err
	}
	var requestBody io.Reader
	if jsonBytes, err := json.Marshal(body); err != nil {
		return nil, -1, err
	} else {
		requestBody = bytes.NewBuffer(jsonBytes)
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.Scheduler{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) GetScheduler(ctx context.Context, pathParams *goqmodel.GetSchedulerPathParams) (*goqmodel.Scheduler, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/schedulers/%v", c.baseURL, pathParams.SchedulerID))
	if err != nil {
		return nil, -1, err
	}
	var requestBody io.Reader
	req, err := http.NewRequest(http.MethodGet, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.Scheduler{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) PutScheduler(ctx context.Context, pathParams *goqmodel.PutSchedulerPathParams, body *goqmodel.Scheduler) (*goqmodel.Scheduler, int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/schedulers/%v", c.baseURL, pathParams.SchedulerID))
	if err != nil {
		return nil, -1, err
	}
	var requestBody io.Reader
	if jsonBytes, err := json.Marshal(body); err != nil {
		return nil, -1, err
	} else {
		requestBody = bytes.NewBuffer(jsonBytes)
	}
	req, err := http.NewRequest(http.MethodPut, u.String(), requestBody)
	if err != nil {
		return nil, -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	respBody := goqmodel.Scheduler{}
	if len(respBytes) == 0 {
		return nil, resp.StatusCode, nil
	}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return &respBody, resp.StatusCode, nil
}
func (c *client) DeleteScheduler(ctx context.Context, pathParams *goqmodel.DeleteSchedulerPathParams) (int, error) {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("%s/schedulers/%v", c.baseURL, pathParams.SchedulerID))
	if err != nil {
		return -1, err
	}
	var requestBody io.Reader
	req, err := http.NewRequest(http.MethodDelete, u.String(), requestBody)
	if err != nil {
		return -1, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(respBytes))
	}
	return resp.StatusCode, nil
}
