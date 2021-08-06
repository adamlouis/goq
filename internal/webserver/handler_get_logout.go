package webserver

import (
	"net/http"
)

func (wh *webHandler) GetLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "web-session")
	session.Options.MaxAge = 0
	store.Save(r, w, session)
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusFound)
}
