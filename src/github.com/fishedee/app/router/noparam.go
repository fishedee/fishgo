package router

import (
	"net/http"
)

func NewNoParamMiddleware() RouterMiddleware {
	return func(handler []interface{}) interface{} {
		last := handler[len(handler)-1]
		var netHandler http.HandlerFunc
		if handler, isNoParam := last.(http.HandlerFunc); isNoParam == true {
			netHandler = handler
		} else if handler2, isNoParam := last.(func(w http.ResponseWriter, r *http.Request)); isNoParam == true {
			netHandler = handler2
		} else if handler3, isNoParam := last.(http.Handler); isNoParam == true {
			netHandler = handler3.ServeHTTP
		} else {
			return last
		}
		return func(w http.ResponseWriter, r *http.Request, param RouterParam) {
			netHandler(w, r)
		}
	}
}