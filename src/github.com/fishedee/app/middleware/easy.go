package middleware

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/app/metric"
	. "github.com/fishedee/app/quicktag"
	. "github.com/fishedee/app/render"
	. "github.com/fishedee/app/router"
	. "github.com/fishedee/app/session"
	. "github.com/fishedee/app/validator"
	. "github.com/fishedee/language"
	"net/http"
	"strings"
)

func NewEasyMiddleware(log Log, validatorFactory ValidatorFactory, sessionFactory SessionFactory, renderFactory RenderFactory, metric Metric) RouterMiddleware {
	var serverError MetricCounter
	if metric != nil {
		serverError = metric.GetCounter("server.error")
	}
	jsonQuickTag := NewQuickTag("json")

	return func(prev RouterMiddlewareContext) RouterMiddlewareContext {
		lastHandler, isOk := prev.Handler.(func(v Validator, s Session) interface{})
		if isOk == false {
			return prev
		}
		name := prev.Data["name"].(string)
		nameInfo := Explode(name, "_")
		if len(nameInfo) == 0 {
			return prev
		}
		renderName := strings.ToLower(nameInfo[len(nameInfo)-1])

		renderChange := func(err Exception, result interface{}) interface{} {
			if renderName == "raw" {
				if err.GetCode() != 0 {
					return []byte(err.GetMessage())
				} else {
					return result
				}
			} else if renderName == "file" {
				if err.GetCode() != 0 {
					return "error.html"
				} else {
					return result
				}
			} else if renderName == "text" {
				if err.GetCode() != 0 {
					return err.GetMessage()
				} else {
					return result
				}
			} else if renderName == "redirect" {
				if err.GetCode() != 0 {
					return []interface{}{302, "/"}
				} else {
					return result
				}
			} else if renderName == "json" {
				if err.GetCode() != 0 {
					return map[string]interface{}{
						"code": err.GetCode(),
						"msg":  err.GetMessage(),
						"data": nil,
					}
				} else {
					return map[string]interface{}{
						"code": 0,
						"msg":  "",
						"data": jsonQuickTag.GetTagInstance(result),
					}
				}
			} else if renderName == "html" {
				if err.GetCode() != 0 {
					return []interface{}{"error.html", map[string]interface{}{
						"code":  err.GetCode(),
						"msg":   err.GetMessage(),
						"stack": err.GetStackTrace(),
					}}
				} else {
					return result
				}
			} else {
				return result
			}
		}
		return RouterMiddlewareContext{
			Data: prev.Data,
			Handler: func(w http.ResponseWriter, r *http.Request, p RouterParam) {
				param := map[string]string{}
				if p != nil {
					for _, singleP := range p {
						param[singleP.Key] = singleP.Value
					}
				}
				validator := validatorFactory.Create(r, param)
				session := sessionFactory.Create(w, r)
				render := renderFactory.Create(w, r)
				var result interface{}
				var exception Exception
				func() {
					doException := func(e Exception) {
						if serverError != nil {
							serverError.Inc(1)
						}
						exception = e
						log.Error("Buiness Error Code:[%d] Message:[%s]\nStackTrace:[%s]", e.GetCode(), e.GetMessage(), e.GetStackTrace())
					}
					defer Catch(doException)
					result = lastHandler(validator, session)
					if resultException, isOk := result.(*Exception); isOk {
						result = nil
						doException(*resultException)
					} else if resultError, isOk := result.(error); isOk {
						result = nil
						doException(*NewException(1, resultError.Error()))
					}
				}()
				data := renderChange(exception, result)
				err := render.Format(renderName, data)
				if err != nil {
					panic(err)
				}
			},
		}

	}
}
