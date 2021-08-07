package webserver

import (
	"net/http"

	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/internal/session"
)

func (wh *webHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		jsonlog.Log("error", err.Error())
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if !wh.checker.Check(username, password) {
		http.Error(w, "failed to log in", http.StatusBadRequest)
		return
	}

	if err := wh.sessionManger.Delete(w, r); err != nil {
		jsonlog.Log("error", err.Error())
	}

	if err := wh.sessionManger.Create(w, r, &session.Payload{
		Username:      username,
		Authenticated: true,
	}); err != nil {
		jsonlog.Log("error", err.Error())
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
