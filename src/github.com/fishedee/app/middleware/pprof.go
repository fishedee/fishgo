package middleware

import (
	. "github.com/fishedee/app/router"
	"net/http"
	"net/http/pprof"
)

func NewPProfMiddleware() RouterMiddleware {
	path := "/debug/pprof/"
	pathLen := len(path)
	return func(prev RouterMiddlewareContext) RouterMiddlewareContext {
		last := prev.Handler.(func(w http.ResponseWriter, r *http.Request, param RouterParam))
		return RouterMiddlewareContext{
			Data: prev.Data,
			Handler: func(w http.ResponseWriter, r *http.Request, param RouterParam) {
				if r.URL.Path == "/debug/pprof" {
					http.Redirect(w, r, "/debug/pprof/", 301)
					return
				}
				if len(r.URL.Path) >= pathLen &&
					r.URL.Path[0:pathLen] == path {
					url := r.URL.Path[pathLen:]
					if url == "/cmdline" {
						pprof.Cmdline(w, r)
					} else if url == "/profile" {
						pprof.Profile(w, r)
					} else if url == "/symbol" {
						pprof.Symbol(w, r)
					} else if url == "/trace" {
						pprof.Trace(w, r)
					} else {
						pprof.Index(w, r)
					}
					return
				}
				last(w, r, param)
			},
		}
	}
}
