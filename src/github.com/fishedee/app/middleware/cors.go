package middleware

import (
	. "github.com/fishedee/app/cors"
	. "github.com/fishedee/app/router"
	"net/http"
)

func NewCorsMiddleware(cors Cors) RouterMiddleware {
	return func(handler []interface{}) interface{} {
		last := handler[len(handler)-1].(func(w http.ResponseWriter, r *http.Request, param RouterParam))
		return func(w http.ResponseWriter, r *http.Request, param RouterParam) {
			cors.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
				last(w, r, param)
			})
		}
	}
}
