package webserver

import (
	"net/http"

	"github.com/adamlouis/goq/internal/pkg/jsonlog"
)

func (wh *webHandler) GetLogout(w http.ResponseWriter, r *http.Request) {
	if err := wh.sessionManger.Delete(w, r); err != nil {
		jsonlog.Log("error", err.Error())
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
