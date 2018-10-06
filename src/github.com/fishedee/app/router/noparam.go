package router

import (
	. "github.com/fishedee/language"
	"net/http"
	"reflect"
	"runtime"
)

func NewNoParamMiddleware() RouterMiddleware {
	return func(prev RouterMiddlewareContext) RouterMiddlewareContext {
		if _, isExist := prev.Data["name"]; isExist == false {
			name := runtime.FuncForPC(reflect.ValueOf(prev.Handler).Pointer()).Name()
			nameInfo := Explode(name, ".")
			lastName := nameInfo[len(nameInfo)-1]
			if len(lastName) >= 3 && lastName[len(lastName)-3:] == "-fm" {
				lastName = lastName[0 : len(lastName)-3]
			}
			if lastName[0] == '(' {
				lastName = lastName[1:]
			}
			if lastName[len(lastName)-1] == ')' {
				lastName = lastName[:len(lastName)-1]
			}
			prev.Data["name"] = lastName
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
