package router

import (
	"net/http"
	"reflect"
	"runtime"
)

func NewNoParamMiddleware() RouterMiddleware {
	return func(prev RouterMiddlewareContext) RouterMiddlewareContext {
		if _, isExist := prev.Data["name"]; isExist == false {
			prev.Data["name"] = runtime.FuncForPC(reflect.ValueOf(prev.Handler).Pointer()).Name()
		}
		last := prev.Handler
		var netHandler http.HandlerFunc
		if handler, isNoParam := last.(http.HandlerFunc); isNoParam == true {
			netHandler = handler
		} else if handler2, isNoParam := last.(func(w http.ResponseWriter, r *http.Request)); isNoParam == true {
			netHandler = handler2
		} else if handler3, isNoParam := last.(http.Handler); isNoParam == true {
			netHandler = handler3.ServeHTTP
		}
		if netHandler != nil {
			last = func(w http.ResponseWriter, r *http.Request, param RouterParam) {
				netHandler(w, r)
			}
		}
		return RouterMiddlewareContext{
			Data:    prev.Data,
			Handler: last,
		}
	}
}
