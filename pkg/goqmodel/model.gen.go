// GENERATED
// DO NOT EDIT
// GENERATOR: scripts/gencode/gencode.go
// ARGUMENTS: --component model --config ../../api/api.yml --package goqmodel --out ./model.gen.go
package goqmodel

type JSONObject map[string]interface{}
type Job struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Status       string     `json:"status"`
	Input        JSONObject `json:"input"`
	Output       JSONObject `json:"output"`
	ScheduledFor *string    `json:"scheduled_for"`
	SucceededAt  *string    `json:"succeeded_at"`
	ErroredAt    *string    `json:"errored_at"`
	ClaimedAt    *string    `json:"claimed_at"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
}
type ListJobsQueryParams struct {
	Where     string `json:"where"`
	OrderBy   string `json:"order_by"`
	PageSize  int    `json:"page_size"`
	PageToken string `json:"page_token"`
}
type ListJobsResponse struct {
	Jobs          []*Job `json:"jobs"`
	NextPageToken string `json:"next_page_token"`
}
type ClaimSomeJobRequest struct {
	Names []string `json:"names"`
}
type Scheduler struct {
	ID        string     `json:"id"`
	Schedule  string     `json:"schedule"`
	JobName   string     `json:"job_name"`
	Input     JSONObject `json:"input"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}
type ListSchedulersRequest struct {
	PageToken string `json:"page_token"`
	PageSize  int    `json:"page_size"`
}
type ListSchedulersResponse struct {
	Schedulers    []*Scheduler `json:"schedulers"`
	NextPageToken string       `json:"next_page_token"`
}
type GetJobPathParams struct {
	JobID string
}
type DeleteJobPathParams struct {
	JobID string
}
type ClaimJobPathParams struct {
	JobID string
}
type ReleaseJobPathParams struct {
	JobID string
}
type SetJobSuccessPathParams struct {
	JobID string
}
type SetJobErrorPathParams struct {
	JobID string
}
type GetSchedulerPathParams struct {
	SchedulerID string
}
type PutSchedulerPathParams struct {
	SchedulerID string
}
type DeleteSchedulerPathParams struct {
	SchedulerID string
}
