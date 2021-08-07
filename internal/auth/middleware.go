package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/adamlouis/goq/internal/session"
	"github.com/gorilla/mux"
)

type ctxKeyType string

type AuthTypeValue string

const (
	AuthType ctxKeyType = "AUTH_TYPE"
	Username ctxKeyType = "USERNAME"

	AuthTypeWebSession AuthTypeValue = "WEB_SESSION"
	AuthTypeBearer     AuthTypeValue = "BEARER"
	_bearerPrefix                    = "Bearer "
)

func GetMiddleware(sessionManager session.Manager, apiKeyChecker KChecker) []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{
		// auth for web session
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if p, _ := sessionManager.Get(w, r); p != nil {
					if p.Authenticated {
						ctx := context.WithValue(
							context.WithValue(
								r.Context(),
								AuthType,
								AuthTypeWebSession),
							Username,
							p.Username,
						)
						r = r.WithContext(ctx)
					}
				}
				next.ServeHTTP(w, r)
			})
		},
		// auth for api bearer token
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				az := r.Header["Authorization"]
				if len(az) == 1 && strings.HasPrefix(az[0], _bearerPrefix) {
					hv := az[0]
					token := hv[len(_bearerPrefix):]
					if apiKeyChecker.Check(token) {
						ctx := context.WithValue(r.Context(), AuthType, AuthTypeBearer)
						r = r.WithContext(ctx)
					}
				}
				next.ServeHTTP(w, r)
			})
		},
		// authz for all routes
		// could make this more declarative
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// by default, DO NOT allow ... unless some condition is met
				allow := false
				authType, ok := r.Context().Value(AuthType).(AuthTypeValue)

				// all can access login & logout
				if r.URL.Path == "/login" || r.URL.Path == "/logout" {
					allow = true
				}

				// api bearer token can access api
				if strings.HasPrefix(r.URL.Path, "/api") && ok && authType == AuthTypeBearer {
					allow = true
				}

				// web token can access all non-api routes
				if !strings.HasPrefix(r.URL.Path, "/api") && ok && authType == AuthTypeWebSession {
					allow = true
				}

				// allow if some condition is met
				if allow {
					next.ServeHTTP(w, r)
					return
				}

				// return a terse error for api
				if strings.HasPrefix(r.URL.Path, "/api") {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}

				// return a styled error for web
				w.Header().Add("Content-Type", "text/html")
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `
				<html>
					<head>
					</head>
						<style>
							body {
								font-family: 'Courier New', monospace;
							}

							h1 {
								font-size: 1.2em;
							}
							h2 {
								font-size: 1em;
							}
						</style>
					<body>
						<h1>forbidden</h1>
						<a href="/login">login</a>
					</body>
				</html>
				`)
			})
		},
	}

}
