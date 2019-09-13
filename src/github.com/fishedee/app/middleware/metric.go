package middleware

import (
	. "github.com/fishedee/app/metric"
	. "github.com/fishedee/app/router"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	"net/http"
	"time"
)

//欠缺对status的错误和崩溃的上传
func NewMetricMiddleware(metric Metric, tags map[string]string) RouterMiddleware {
	tagList := []string{}
	if tags != nil {
		for k, v := range tags {
			vEncode, err := EncodeUrl(v)
			if err != nil {
				panic(err)
			}
			tagList = append(tagList, k+"="+vEncode)
		}
	}
	tagStr := Implode(tagList, "&")

	return func(prev RouterMiddlewareContext) RouterMiddlewareContext {
		last := prev.Handler.(func(w http.ResponseWriter, r *http.Request, param RouterParam))
		return RouterMiddlewareContext{
			Data: prev.Data,
			Handler: func(w http.ResponseWriter, r *http.Request, param RouterParam) {
				url := r.URL.Path
				begin := time.Now()
				defer func() {
					end := time.Now()
					duration := end.Sub(begin)
					name := "router_time?path=" + url
					if len(tagStr) != 0 {
						name = name + "&" + tagStr
					}
					metric.UpdateTimer(name, duration)
				}()
				last(w, r, param)
			},
		}
	}
}
