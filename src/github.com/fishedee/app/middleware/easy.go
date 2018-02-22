package middleware

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/app/render"
	. "github.com/fishedee/app/router"
	. "github.com/fishedee/app/session"
	. "github.com/fishedee/app/validator"
	. "github.com/fishedee/language"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

func NewEasyMiddleware(log Log, validatorFactory ValidatorFactory, sessionFactory SessionFactory, renderFactory RenderFactory) RouterMiddleware {
	return func(handler []interface{}) interface{} {
		last := handler[len(handler)-1]
		lastHandler, isOk := last.(func(v Validator, s Session) interface{})
		if isOk == false {
			return last
		}
		first := handler[0]
		name := runtime.FuncForPC(reflect.ValueOf(first).Pointer()).Name()
		nameInfo := Explode(name, "_")
		if len(nameInfo) == 0 {
			return last
		}
		renderName := strings.ToLower(nameInfo[len(nameInfo)-1])
		renderChange := func(err Exception, result interface{}) interface{} {
			if renderName == "raw" {
				if err.GetCode() != 0 {
					return []byte(err.GetMessage())
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
						"data": result,
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
		return func(w http.ResponseWriter, r *http.Request, p RouterParam) {
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
				defer Catch(func(e Exception) {
					exception = e
					log.Error("Buiness Error Code:[%d] Message:[%s]\nStackTrace:[%s]", e.GetCode(), e.GetMessage(), e.GetStackTrace())
				})
				result = lastHandler(validator, session)
			}()
			data := renderChange(exception, result)
			err := render.Format(renderName, data)
			if err != nil {
				panic(err)
			}
		}

	}
}
