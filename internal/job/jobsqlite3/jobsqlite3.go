package jobsqlite3

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/adamlouis/goq/internal/job"
	"github.com/adamlouis/goq/internal/jsonlogic"
	"github.com/adamlouis/goq/internal/jsonlogic/jsonlogicsqlite3"
	"github.com/adamlouis/goq/internal/pkg/crudutil"
	"github.com/adamlouis/goq/internal/pkg/errtype"
	"github.com/adamlouis/goq/internal/pkg/sqlite3util"
	"github.com/adamlouis/goq/pkg/goqmodel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewJobRepository(db sqlx.Ext) job.Repository {
	return &jobRepo{
		db: db,
	}
}

type jobRepo struct {
	db sqlx.Ext
}

var (
	searchColumns = []string{
		"id",
		"name",
		"status",
		"input",
		"output",
		"succeed_at",
		"claimed_at",
		"scheduled_for",
		"errored_at",
		"created_at",
		"updated_at",
	}
)

//go:embed migration/*.sql
var MigrationFS embed.FS

type jobRow struct {
	ID           string  `db:"id"`
	Name         string  `db:"name"`
	Status       string  `db:"status"`
	Input        []byte  `db:"input"`
	Output       []byte  `db:"output"`
	SucceededAt  *string `db:"succeed_at"`
	ClaimedAt    *string `db:"claimed_at"`
	ScheduledFor *string `db:"scheduled_for"`
	ErroredAt    *string `db:"errored_at"`
	CreatedAt    string  `db:"created_at"`
	UpdatedAt    string  `db:"updated_at"`
}

func (jr *jobRepo) Queue(ctx context.Context, j *goqmodel.Job) (*goqmodel.Job, error) {
	j.ID = uuid.New().String()
	j.Status = string(job.JobStatusQueued)

	var scheduledForStr *string
	if j.ScheduledFor != nil {
		scheduledFor, err := toInternalTime(*j.ScheduledFor)
		if err != nil {
			return nil, err
		}
		s := scheduledFor.Format(sqlite3util.DatetimeFormat)
		scheduledForStr = &s
	}

	input, err := json.Marshal(j.Input)
	if err != nil {
		return nil, err
	}

	_, err = jr.db.Exec(`
				INSERT INTO
					job
						(id, name, status, input, scheduled_for)
					VALUES
						(?, ?, ?, ?, ?)`,
		j.ID, j.Name, j.Status, input, scheduledForStr)
	if err != nil {
		return nil, err
	}

	return jr.Get(ctx, j.ID)
}

func (jr *jobRepo) Delete(ctx context.Context, id string) error {
	return crudutil.Delete(jr.db, `DELETE FROM job WHERE id = ?`, id)
}

func (jr *jobRepo) Get(ctx context.Context, id string) (*goqmodel.Job, error) {
	row := jr.db.QueryRowx(`
			SELECT
				id, name, status, input, output, succeed_at, errored_at, claimed_at, created_at, updated_at, scheduled_for
			FROM job
			WHERE id = ?`,
		id,
	)

	var r jobRow
	err := row.StructScan(&r)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errtype.NotFoundError{Err: err}
		}
		return nil, err
	}

	j, err := jobRowToJob(&r)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func min(m, n int) int {
	if n < m {
		return n
	}
	return m
}

func (jr *jobRepo) Search(ctx context.Context, r *goqmodel.SearchJobsRequest) (*goqmodel.SearchJobsResponse, error) {
	sb := sq.
		StatementBuilder.
		Select(searchColumns...).
		From("job")

	sqlz := jsonlogicsqlite3.NewSQLizer()

	where, err := sqlz.ToSQL(r.Where)
	if err != nil {
		return nil, err
	}
	fmt.Println("WHERE", where)
	sb = sb.Where(where)

	orderBys, err := jsonlogic.AllToSQL(sqlz, r.OrderBy)
	if err != nil {
		return nil, err
	}
	orderBys = append(orderBys, "id ASC") // add id for stable order
	sb = sb.OrderBy(orderBys...)

	fmt.Println("ORDER BY", orderBys)
	sb = sb.Limit(uint64(min(r.PageSize, 100)))

	sql, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	fmt.Println(sql, args)
	rows, err := jr.db.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	jobs := []*goqmodel.Job{}
	for rows.Next() {
		var r jobRow
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}
		p, err := jobRowToJob(&r)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, p)
	}
	return &goqmodel.SearchJobsResponse{
		Jobs: jobs,
	}, nil

}
func (jr *jobRepo) List(ctx context.Context, args *goqmodel.ListJobsQueryParams) (*goqmodel.ListJobsResponse, error) {
	sz, err := crudutil.GetPageSize(args.PageSize, 500)
	if err != nil {
		return nil, err
	}

	sb := sq.
		StatementBuilder.
		Select("id, name, status, input, output, succeed_at, errored_at, claimed_at, created_at, updated_at, scheduled_for").
		From("job").
		OrderBy("created_at desc, id ASC").
		Limit(uint64(sz) + 1) // get n+1 so we know if there's a next page

	offset := uint64(0)
	if args.PageToken != "" {
		page := &listJobsPageData{}
		err := crudutil.DecodePageData(args.PageToken, page)
		if err != nil {
			return nil, err
		}
		offset = page.Offset
	}
	sb = sb.Offset(offset)

	sql, sqlArgs, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := jr.db.Queryx(sql, sqlArgs...)
	if err != nil {
		return nil, err
	}

	jobs := make([]*goqmodel.Job, 0, sz)

	for rows.Next() {
		var r jobRow
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}
		j, err := jobRowToJob(&r)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}

	nextPageToken := ""
	if len(jobs) > int(sz) {
		jobs = jobs[0 : len(jobs)-1]
		s, err := crudutil.EncodePageData(&listJobsPageData{
			Offset: offset + uint64(len(jobs)),
		})
		if err != nil {
			return nil, err
		}
		nextPageToken = s
	}

	return &goqmodel.ListJobsResponse{
		Jobs:          jobs,
		NextPageToken: nextPageToken,
	}, nil
}

type listJobsPageData struct {
	Offset uint64 `json:"offset"`
}

// // TODO - handle concurrency, locking, etc
func (jr *jobRepo) Claim(ctx context.Context, opts *job.ClaimOptions) (*goqmodel.Job, error) {
	sb := sq.
		StatementBuilder.
		Select("id, status").
		From("job").
		OrderBy("created_at ASC, id ASC").
		Where(sq.Eq{"status": job.JobStatusQueued}).
		Where(sq.Or{
			sq.Expr("scheduled_for IS NULL"),
			sq.LtOrEq{"scheduled_for": "CURRENT_TIMESTAMP"},
		}).
		Limit(1)

	if opts.JobID != "" {
		sb = sb.Where(sq.Eq{"id": opts.JobID})
	}

	if len(opts.Names) > 0 {
		ors := make(sq.Or, len(opts.Names))
		for i, n := range opts.Names {
			ors[i] = sq.Eq{"name": n}
		}
		sb = sb.Where(ors)
	}

	query, queryArgs, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	row := jr.db.QueryRowx(query, queryArgs...)

	var r jobRow
	err = row.StructScan(&r)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	r.Status = string(job.JobStatusClaimed)
	_, err = jr.db.Exec(`UPDATE job SET status = ?, claimed_at = CURRENT_TIMESTAMP WHERE id = ?`, r.Status, r.ID)
	if err != nil {
		return nil, err
	}

	return jr.Get(ctx, r.ID)
}

func (jr *jobRepo) Release(ctx context.Context, id string) (*goqmodel.Job, error) {
	j, err := jr.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if j.Status != string(job.JobStatusClaimed) {
		return nil, fmt.Errorf("only jobs with status CLAIMED can be released - %s has status %s", j.ID, j.Status)
	}

	j.Status = string(job.JobStatusQueued)
	_, err = jr.db.Exec(`UPDATE job SET status = ?, claimed_at = NULL WHERE id = ?`, job.JobStatusQueued, j.ID)
	if err != nil {
		return nil, err
	}

	return jr.Get(ctx, id)
}
func (jr *jobRepo) Success(ctx context.Context, id string, output goqmodel.JSONObject) (*goqmodel.Job, error) {
	j, err := jr.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if j.Status != string(job.JobStatusClaimed) {
		return nil, fmt.Errorf("only jobs with status CLAIMED can be updated with success - %s has status %s", j.ID, j.Status)
	}
	j.Status = string(job.JobStatusSuccess)

	var ob []byte
	if output != nil {
		b, err := json.Marshal(output)
		if err != nil {
			return nil, err
		}
		ob = b
	}

	_, err = jr.db.Exec(`UPDATE job SET status = ?, output = ?, succeed_at = CURRENT_TIMESTAMP WHERE id = ?`, job.JobStatusSuccess, ob, j.ID)
	if err != nil {
		return nil, err
	}

	return jr.Get(ctx, id)
}
func (jr *jobRepo) Error(ctx context.Context, id string, output goqmodel.JSONObject) (*goqmodel.Job, error) {
	j, err := jr.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if j.Status != string(job.JobStatusClaimed) {
		return nil, fmt.Errorf("only jobs with status CLAIMED can be update with error - %s has status %s", j.ID, j.Status)
	}
	j.Status = string(job.JobStatusError)

	var ob []byte
	if output != nil {
		b, err := json.Marshal(output)
		if err != nil {
			return nil, err
		}
		ob = b
	}

	_, err = jr.db.Exec(`UPDATE job SET status = ?, output = ?, errored_at = CURRENT_TIMESTAMP WHERE id = ?`, job.JobStatusError, ob, j.ID)
	if err != nil {
		return nil, err
	}

	return jr.Get(ctx, id)
}

func jobRowToJob(r *jobRow) (*goqmodel.Job, error) {
	c, err := time.Parse(sqlite3util.DatetimeFormat, r.CreatedAt)
	if err != nil {
		return nil, err
	}

	u, err := time.Parse(sqlite3util.DatetimeFormat, r.UpdatedAt)
	if err != nil {
		return nil, err
	}

	sa, err := toAPITimePtr(r.SucceededAt)
	if err != nil {
		return nil, err
	}

	er, err := toAPITimePtr(r.ErroredAt)
	if err != nil {
		return nil, err
	}

	ca, err := toAPITimePtr(r.ClaimedAt)
	if err != nil {
		return nil, err
	}

	sf, err := toAPITimePtr(r.ScheduledFor)
	if err != nil {
		return nil, err
	}

	var input map[string]interface{}
	err = json.Unmarshal(r.Input, &input)
	if err != nil {
		return nil, err
	}

	var output map[string]interface{}
	if r.Output != nil {
		err = json.Unmarshal(r.Output, &output)
		if err != nil {
			return nil, err
		}
	}

	return &goqmodel.Job{
		ID:           r.ID,
		Name:         r.Name,
		Status:       r.Status,
		Input:        input,
		Output:       output,
		SucceededAt:  sa,
		ErroredAt:    er,
		ClaimedAt:    ca,
		ScheduledFor: sf,
		CreatedAt:    toAPITime(c),
		UpdatedAt:    toAPITime(u),
	}, nil
}

func toAPITimePtr(s *string) (*string, error) {
	if s != nil {
		t, err := time.Parse(sqlite3util.DatetimeFormat, *s)
		if err != nil {
			return nil, err
		}
		apist := toAPITime(t)
		return &apist, nil
	}
	return nil, nil
}

func toAPITime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func toInternalTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
