package webserver

import (
	"embed"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

//go:embed templates
var templatesFS embed.FS

type WebHandler interface {
	Home(w http.ResponseWriter, req *http.Request)
}

func NewWebHandler() WebHandler {
	return &whs{}
}

type whs struct{}

type pageData struct {
	Title string
}

func RegisterRouter(wh WebHandler, r *mux.Router) {
	r.Handle("/", http.HandlerFunc(wh.Home)).Methods(http.MethodGet)
}

func (w *whs) Home(wrt http.ResponseWriter, req *http.Request) {
	t := template.Must(template.New("home.go.html").ParseFS(templatesFS, "templates/home.go.html"))
	t.Execute(wrt, pageData{
		Title: "GOQ",
	})
}
