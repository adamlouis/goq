package session

import (
	"net/http"

	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/gorilla/sessions"
)

type Payload struct {
	Username      string
	Authenticated bool
}

type Manager interface {
	Delete(w http.ResponseWriter, r *http.Request) error
	Create(w http.ResponseWriter, r *http.Request, payload *Payload) error
	Get(w http.ResponseWriter, r *http.Request) (*Payload, error)
}

func NewFSManager(cookieName string, store *sessions.FilesystemStore) Manager {
	return &fssm{
		cookieName: cookieName,
		store:      store,
	}
}

type fssm struct {
	cookieName string
	store      *sessions.FilesystemStore
}

func (f *fssm) Delete(w http.ResponseWriter, r *http.Request) error {
	s, _ := f.store.Get(r, f.cookieName)
	if !s.IsNew {
		s.Options.MaxAge = 0
		return s.Save(r, w)
	}
	return nil
}
func (f *fssm) Create(w http.ResponseWriter, r *http.Request, payload *Payload) error {
	s, err := f.store.New(r, f.cookieName)
	if err != nil {
		// TODO: why "the value is not valid"?
		// see source:
		// ErrMacInvalid = cookieError{typ: decodeError, msg: "the value is not valid"}
		jsonlog.Log("error", err.Error())
	}
	s.Values["username"] = payload.Username
	s.Values["authenticated"] = payload.Authenticated
	return s.Save(r, w)
}
func (f *fssm) Get(w http.ResponseWriter, r *http.Request) (*Payload, error) {
	session, err := f.store.Get(r, f.cookieName)
	if err != nil {
		return nil, err
	}
	p := &Payload{}
	if a, ok := session.Values["authenticated"].(bool); ok {
		p.Authenticated = a
	}
	if u, ok := session.Values["username"].(string); ok {
		p.Username = u
	}
	return p, nil
}
