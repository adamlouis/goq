package webserver

import (
	"encoding/json"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

type pageData struct {
	Username     string
	JQ           string
	Title        string
	JobStr       string
	SchedulerStr string
	JobNames     []string
	Pivot        [][]string
	PivotLinks   [][]string
	Jobs         []*JobTmpl
	Schedulers   []*SchedulerTmpl
}

type SchedulerTmpl struct {
	*goqmodel.Scheduler
	InputStr string
}

type JobTmpl struct {
	*goqmodel.Job
	InputStr, OutputStr string
}

func toSchedulerTmpls(ss []*goqmodel.Scheduler) []*SchedulerTmpl {
	r := make([]*SchedulerTmpl, len(ss))
	for i := range ss {
		ib, _ := json.Marshal(ss[i].Input)

		r[i] = &SchedulerTmpl{
			Scheduler: ss[i],
			InputStr:  string(ib),
		}
	}
	return r
}

func toJobTmpls(js []*goqmodel.Job) []*JobTmpl {
	r := make([]*JobTmpl, len(js))
	for i := range js {
		ib, _ := json.Marshal(js[i].Input)
		ob, _ := json.Marshal(js[i].Output)

		r[i] = &JobTmpl{
			Job:       js[i],
			InputStr:  string(ib),
			OutputStr: string(ob),
		}
	}
	return r
}
