package router

import (
	"net/http"
)

type RouterCorsInterface interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

func NewCorsMiddleware(cors RouterCorsInterface) RouterMiddleware {
	return func(handler []interface{}) interface{} {
		last := handler[len(handler)-1].(func(w http.ResponseWriter, r *http.Request, param RouterParam))
		return func(w http.ResponseWriter, r *http.Request, param RouterParam) {
			cors.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
				last(w, r, param)
			})
		}
	}
}
