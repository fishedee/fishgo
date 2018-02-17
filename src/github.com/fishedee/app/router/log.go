package router

import (
	. "github.com/fishedee/language"
	"net/http"
	"time"
)

type RouterLogInterface interface {
	Critical(format string, v ...interface{})
	Error(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

func NewLogMiddleware(log RouterLogInterface) RouterMiddleware {
	return func(handler []interface{}) interface{} {
		last := handler[len(handler)-1].(func(w http.ResponseWriter, r *http.Request, param map[string]string))
		run := func(w http.ResponseWriter, r *http.Request, param map[string]string) {
			defer CatchCrash(func(exception Exception) {
				log.Critical("Buiness Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
				w.WriteHeader(500)
				w.Write([]byte("server internal error"))
			})
			last(w, r, param)
		}
		return func(w http.ResponseWriter, r *http.Request, param map[string]string) {
			beginTime := time.Now().UnixNano()
			run(w, r, param)
			endTime := time.Now().UnixNano()
			log.Debug("%s %s : %s", r.Method, r.URL.String(), time.Duration(endTime-beginTime).String())
		}

	}
}
