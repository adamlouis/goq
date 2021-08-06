package webserver

import (
	"io/ioutil"
	"net/http"
	"os"
)

func (wh *webHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	// parse form
	body, _ := ioutil.ReadAll(r.Body)
	form := dumbDecode(body)
	username := form["username"]
	password := form["password"]

	// check credentials
	ru := os.Getenv("GOQ_ROOT_USERNAME")
	rp := os.Getenv("GOQ_ROOT_PASSWORD")

	matches := len(ru) > 0 && len(rp) > 0 && ru == username && rp == password

	if !matches {
		http.Error(w, "failed to log in", http.StatusBadRequest)
		return
	}

	// delete old session
	old, _ := store.Get(r, "web-session")
	old.Options.MaxAge = 0
	old.Save(r, w)

	// create new session
	session, _ := store.New(r, "web-session")
	session.Values["username"] = username
	session.Values["authenticated"] = true
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}
