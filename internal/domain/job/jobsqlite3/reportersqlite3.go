package jobsqlite3

import (
	"context"

	"github.com/adamlouis/goq/internal/domain/job"
	"github.com/jmoiron/sqlx"
)

func NewJobReporter(db sqlx.Ext) job.Reporter {
	return &jobReporter{
		db: db,
	}
}

type jobReporter struct {
	db sqlx.Ext
}

type statusCountRow struct {
	Status string `db:"status"`
	Count  int64  `db:"c"`
}
type nameCountRow struct {
	Name  string `db:"name"`
	Count int64  `db:"c"`
}

func (jr *jobReporter) GetCountByStatus(ctx context.Context) (map[job.JobStatus]int64, error) {
	rows, err := jr.db.Queryx("SELECT status, count(1) as c FROM job GROUP BY status ORDER BY c")
	if err != nil {
		return nil, err
	}

	result := map[job.JobStatus]int64{}
	for rows.Next() {
		var r statusCountRow
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}
		result[job.JobStatus(r.Status)] = r.Count
	}
	return result, nil
}
func (jr *jobReporter) GetCountByName(ctx context.Context) (map[string]int64, error) {
	rows, err := jr.db.Queryx("SELECT name, count(1) as c FROM job GROUP BY name ORDER BY c")
	if err != nil {
		return nil, err
	}

	result := map[string]int64{}
	for rows.Next() {
		var r nameCountRow
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}
		result[r.Name] = r.Count
	}
	return result, nil
}

// 	j.ID = uuid.New().String()
// 	j.Status = string(job.JobStatusQueued)

// 	var scheduledForStr *string
// 	if j.ScheduledFor != nil {
// 		scheduledFor, err := toInternalTime(*j.ScheduledFor)
// 		if err != nil {
// 			return nil, err
// 		}
// 		s := scheduledFor.Format(sqlite3util.DatetimeFormat)
// 		scheduledForStr = &s
// 	}

// 	input, err := json.Marshal(j.Input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = jr.db.Exec(`
// 				INSERT INTO
// 					job
// 						(id, name, status, input, scheduled_for)
// 					VALUES
// 						(?, ?, ?, ?, ?)`,
// 		j.ID, j.Name, j.Status, input, scheduledForStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jr.Get(ctx, j.ID)
// }

// func (jr *jobRepo) Delete(ctx context.Context, id string) error {
// 	return crudutil.Delete(jr.db, `DELETE FROM job WHERE id = ?`, id)
// }

// func (jr *jobRepo) Get(ctx context.Context, id string) (*goqmodel.Job, error) {
// 	row := jr.db.QueryRowx(`
// 			SELECT
// 				id, name, status, input, output, succeed_at, errored_at, claimed_at, created_at, updated_at, scheduled_for
// 			FROM job
// 			WHERE id = ?`,
// 		id,
// 	)

// 	var r jobRow
// 	err := row.StructScan(&r)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, errtype.NotFoundError{Err: err}
// 		}
// 		return nil, err
// 	}

// 	j, err := jobRowToJob(&r)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return j, nil
// }

// func (jr *jobRepo) List(ctx context.Context, args *goqmodel.ListJobsQueryParams) (*goqmodel.ListJobsResponse, error) {
// 	sz, err := crudutil.GetPageSize(args.PageSize, 500)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sb := sq.
// 		StatementBuilder.
// 		Select("id, name, status, input, output, succeed_at, errored_at, claimed_at, created_at, updated_at, scheduled_for").
// 		From("job").
// 		OrderBy("created_at ASC, id ASC").
// 		Limit(uint64(sz) + 1) // get n+1 so we know if there's a next page

// 	offset := uint64(0)
// 	if args.PageToken != "" {
// 		page := &listJobsPageData{}
// 		err := crudutil.DecodePageData(args.PageToken, page)
// 		if err != nil {
// 			return nil, err
// 		}
// 		offset = page.Offset
// 	}
// 	sb = sb.Offset(offset)

// 	sql, sqlArgs, err := sb.ToSql()
// 	if err != nil {
// 		return nil, err
// 	}

// 	rows, err := jr.db.Queryx(sql, sqlArgs...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	jobs := make([]*goqmodel.Job, 0, sz)

// 	for rows.Next() {
// 		var r jobRow
// 		err = rows.StructScan(&r)
// 		if err != nil {
// 			return nil, err
// 		}
// 		j, err := jobRowToJob(&r)
// 		if err != nil {
// 			return nil, err
// 		}
// 		jobs = append(jobs, j)
// 	}

// 	nextPageToken := ""
// 	if len(jobs) > int(sz) {
// 		jobs = jobs[0 : len(jobs)-1]
// 		s, err := crudutil.EncodePageData(&listJobsPageData{
// 			Offset: offset + uint64(len(jobs)),
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 		nextPageToken = s
// 	}

// 	return &goqmodel.ListJobsResponse{
// 		Jobs:          jobs,
// 		NextPageToken: nextPageToken,
// 	}, nil
// }

// type listJobsPageData struct {
// 	Offset uint64 `json:"offset"`
// }

// // // TODO - handle concurrency, locking, etc
// func (jr *jobRepo) Claim(ctx context.Context, opts *job.ClaimOptions) (*goqmodel.Job, error) {
// 	sb := sq.
// 		StatementBuilder.
// 		Select("id, status").
// 		From("job").
// 		OrderBy("created_at ASC, id ASC").
// 		Where(sq.Eq{"status": job.JobStatusQueued}).
// 		Where(sq.Or{
// 			sq.Expr("scheduled_for IS NULL"),
// 			sq.LtOrEq{"scheduled_for": "CURRENT_TIMESTAMP"},
// 		}).
// 		Limit(1)

// 	if opts.JobID != "" {
// 		sb = sb.Where(sq.Eq{"id": opts.JobID})
// 	}

// 	if len(opts.Names) > 0 {
// 		ors := make(sq.Or, len(opts.Names))
// 		for i, n := range opts.Names {
// 			ors[i] = sq.Eq{"name": n}
// 		}
// 		sb = sb.Where(ors)
// 	}

// 	query, queryArgs, err := sb.ToSql()
// 	if err != nil {
// 		return nil, err
// 	}

// 	row := jr.db.QueryRowx(query, queryArgs...)

// 	var r jobRow
// 	err = row.StructScan(&r)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	r.Status = string(job.JobStatusClaimed)
// 	_, err = jr.db.Exec(`UPDATE job SET status = ?, claimed_at = CURRENT_TIMESTAMP WHERE id = ?`, r.Status, r.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jr.Get(ctx, r.ID)
// }

// func (jr *jobRepo) Release(ctx context.Context, id string) (*goqmodel.Job, error) {
// 	j, err := jr.Get(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if j.Status != string(job.JobStatusClaimed) {
// 		return nil, fmt.Errorf("only jobs with status CLAIMED can be released - %s has status %s", j.ID, j.Status)
// 	}

// 	j.Status = string(job.JobStatusQueued)
// 	_, err = jr.db.Exec(`UPDATE job SET status = ?, claimed_at = NULL WHERE id = ?`, job.JobStatusQueued, j.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jr.Get(ctx, id)
// }
// func (jr *jobRepo) Success(ctx context.Context, id string) (*goqmodel.Job, error) {
// 	j, err := jr.Get(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if j.Status != string(job.JobStatusClaimed) {
// 		return nil, fmt.Errorf("only jobs with status CLAIMED can be updated with success - %s has status %s", j.ID, j.Status)
// 	}
// 	j.Status = string(job.JobStatusSuccess)

// 	output, err := json.Marshal(j.Output)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = jr.db.Exec(`UPDATE job SET status = ?, output = ?, succeed_at = CURRENT_TIMESTAMP WHERE id = ?`, job.JobStatusSuccess, output, j.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jr.Get(ctx, id)
// }
// func (jr *jobRepo) Error(ctx context.Context, id string) (*goqmodel.Job, error) {
// 	j, err := jr.Get(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if j.Status != string(job.JobStatusClaimed) {
// 		return nil, fmt.Errorf("only jobs with status CLAIMED can be update with error - %s has status %s", j.ID, j.Status)
// 	}
// 	j.Status = string(job.JobStatusError)

// 	output, err := json.Marshal(j.Output)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = jr.db.Exec(`UPDATE job SET status = ?, output = ?, errored_at = CURRENT_TIMESTAMP WHERE id = ?`, job.JobStatusError, output, j.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jr.Get(ctx, id)
// }

// func jobRowToJob(r *jobRow) (*goqmodel.Job, error) {
// 	c, err := time.Parse(sqlite3util.DatetimeFormat, r.CreatedAt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	u, err := time.Parse(sqlite3util.DatetimeFormat, r.UpdatedAt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sa, err := toAPITimePtr(r.SucceededAt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	er, err := toAPITimePtr(r.ErroredAt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ca, err := toAPITimePtr(r.ClaimedAt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sf, err := toAPITimePtr(r.ScheduledFor)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var input map[string]interface{}
// 	err = json.Unmarshal(r.Input, &input)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var output map[string]interface{}
// 	err = json.Unmarshal(r.Output, &output)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &goqmodel.Job{
// 		ID:           r.ID,
// 		Name:         r.Name,
// 		Status:       r.Status,
// 		Input:        input,
// 		Output:       output,
// 		SucceededAt:  sa,
// 		ErroredAt:    er,
// 		ClaimedAt:    ca,
// 		ScheduledFor: sf,
// 		CreatedAt:    toAPITime(c),
// 		UpdatedAt:    toAPITime(u),
// 	}, nil
// }

// func toAPITimePtr(s *string) (*string, error) {
// 	if s != nil {
// 		t, err := time.Parse(sqlite3util.DatetimeFormat, *s)
// 		if err != nil {
// 			return nil, err
// 		}
// 		apist := toAPITime(t)
// 		return &apist, nil
// 	}
// 	return nil, nil
// }

// func toAPITime(t time.Time) string {
// 	return t.Format(time.RFC3339)
// }

// func toInternalTime(s string) (time.Time, error) {
// 	return time.Parse(time.RFC3339, s)
// }
