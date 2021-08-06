package webserver

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/adamlouis/goq/pkg/goqmodel"
	"github.com/gorilla/sessions"
)

type SchedulerTmpl struct {
	*goqmodel.Scheduler
	InputStr string
}

type JobTmpl struct {
	*goqmodel.Job
	InputStr, OutputStr string
}

type pageData struct {
	Username     string
	Title        string
	JobStr       string
	SchedulerStr string
	Pivot        [][]string
	Jobs         []*JobTmpl
	Schedulers   []*SchedulerTmpl
}

func newTemplate(name string, patterns []string) *template.Template {
	if os.Getenv("GOQ_MODE") == "DEVELOPMENT" {
		resolved := make([]string, len(patterns))
		for i := range patterns {
			resolved[i] = fmt.Sprintf("internal/webserver/%s", patterns[i])
		}
		return template.Must(template.New(name).ParseFiles(resolved...))
	}
	return template.Must(template.New(name).ParseFS(templatesFS, patterns...))
}

func toStrMap(m map[interface{}]interface{}) map[string]interface{} {
	r := map[string]interface{}{}
	for key, value := range m {
		r[fmt.Sprintf("%v", key)] = value
	}
	return r
}

//go:embed templates
var templatesFS embed.FS

// TODO: pass down, etc.
var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte(os.Getenv("GOQ_SESSION_KEY"))
	store = sessions.NewFilesystemStore(os.Getenv("GOQ_SESSION_STORE_PATH"), key)
)

// use gorilla/sessions
// autogen
func dumbDecode(b []byte) map[string]string {
	result := map[string]string{}
	fields := strings.Split(string(b), "&")
	for _, field := range fields {
		parts := strings.Split(field, "=")
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
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
