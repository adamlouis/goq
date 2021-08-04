package jobsqlite3

import (
	"context"

	"github.com/adamlouis/goq/internal/job"
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
