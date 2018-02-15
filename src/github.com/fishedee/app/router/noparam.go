package router

import (
	"net/http"
)

func NewNoParamMiddleware() RouterMiddleware {
	return func(handler []interface{}) interface{} {
		last := handler[len(handler)-1]
		netHandler, isNoParam := last.(http.HandlerFunc)
		if isNoParam == false {
			return last
		}
		return func(w http.ResponseWriter, r *http.Request, param map[string]string) {
			netHandler(w, r)
		}
	}
}
