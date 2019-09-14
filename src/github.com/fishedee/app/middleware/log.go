package middleware

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/app/metric"
	. "github.com/fishedee/app/router"
	. "github.com/fishedee/language"
	"net/http"
	"time"
)

func NewLogMiddleware(log Log, metric Metric) RouterMiddleware {
	var serverCrash MetricCounter
	var serverRequest MetricTimer
	if metric != nil {
		serverCrash = metric.GetCounter("server.crash")
		serverRequest = metric.GetTimer("server.request")
	}
	return func(prev RouterMiddlewareContext) RouterMiddlewareContext {
		last := prev.Handler.(func(w http.ResponseWriter, r *http.Request, param RouterParam))
		run := func(w http.ResponseWriter, r *http.Request, param RouterParam) {
			defer CatchCrash(func(exception Exception) {
				if serverCrash != nil {
					serverCrash.Inc(1)
				}
				log.Critical("Buiness Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
				w.WriteHeader(500)
				w.Write([]byte("server internal error"))
			})
			last(w, r, param)
		}
		return RouterMiddlewareContext{
			Data: prev.Data,
			Handler: func(w http.ResponseWriter, r *http.Request, param RouterParam) {
				beginTime := time.Now().UnixNano()
				run(w, r, param)
				endTime := time.Now().UnixNano()
				duration := time.Duration(endTime - beginTime)
				if serverRequest != nil {
					serverRequest.Update(duration)
				}
				log.Debug("%s %s : %s", r.Method, r.URL.String(), duration.String())
			},
		}

	}
}
