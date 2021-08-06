package webserver

import "net/http"

func (wh *webHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	newTemplate("login.go.html", []string{"templates/common.go.html", "templates/login.go.html"}).Execute(w, pageData{
		Title: "GOQ",
	})
}
