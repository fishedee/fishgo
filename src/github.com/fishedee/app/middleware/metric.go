package middleware

import (
	. "github.com/fishedee/app/metric"
	. "github.com/fishedee/app/router"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	"net/http"
	"time"
)

func getTaggedName(name string, tags map[string]string) string {
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
	if len(tagStr) != 0 {
		name = name + "?" + tagStr
	}
	return name
}

func NewPathMetricMiddleware(metric Metric, tags map[string]string) RouterMiddleware {
	return func(prev RouterMiddlewareContext) RouterMiddlewareContext {
		newTags := map[string]string{}
		if tags != nil {
			for k, v := range tags {
				newTags[k] = v
			}
		}
		newTags["path"] = prev.Data["path"].(string)
		pathRequest := metric.GetTimer(getTaggedName("path.request", newTags))

		last := prev.Handler.(func(w http.ResponseWriter, r *http.Request, param RouterParam))
		return RouterMiddlewareContext{
			Data: prev.Data,
			Handler: func(w http.ResponseWriter, r *http.Request, param RouterParam) {
				begin := time.Now()
				last(w, r, param)
				end := time.Now()
				duration := end.Sub(begin)
				pathRequest.Update(duration)
			},
		}
	}
}
