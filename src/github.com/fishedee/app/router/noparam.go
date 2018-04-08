package router

import (
	. "github.com/fishedee/language"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

func NewNoParamMiddleware() RouterMiddleware {
	return func(prev RouterMiddlewareContext) RouterMiddlewareContext {
		if _, isExist := prev.Data["name"]; isExist == false {
			name := runtime.FuncForPC(reflect.ValueOf(prev.Handler).Pointer()).Name()
			if rbc := strings.LastIndexByte(name, ')'); rbc != -1 {
				name = name[0:rbc]
			}
			nameInfo := Explode(name, ".")
			prev.Data["name"] = nameInfo[len(nameInfo)-1]
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
