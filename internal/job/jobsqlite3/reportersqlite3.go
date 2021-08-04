package jobsqlite3

import (
	"context"

	"github.com/adamlouis/goq/internal/job"
	"github.com/jmoiron/sqlx"
)

func NewJobReporter(db *sqlx.DB) job.Reporter {
	return &jobReporter{
		db: db,
	}
}

type jobReporter struct {
	db *sqlx.DB
}

type nameStatusCountRow struct {
	Name   string `db:"name"`
	Status string `db:"status"`
	Count  int64  `db:"c"`
}

func (jr *jobReporter) GetCountByNameByStatus(ctx context.Context) (map[string]map[job.JobStatus]int64, error) {
	tx, err := jr.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Queryx("SELECT name, status, count(1) as c FROM job GROUP BY name, status ORDER BY c")
	if err != nil {
		return nil, err
	}

	result := map[string]map[job.JobStatus]int64{}
	for rows.Next() {
		var r nameStatusCountRow
		err = rows.StructScan(&r)
		if err != nil {
			return nil, err
		}

		if _, ok := result[r.Name]; !ok {
			result[r.Name] = map[job.JobStatus]int64{}
		}
		result[r.Name][job.JobStatus(r.Status)] = r.Count
	}
	return result, nil
}
