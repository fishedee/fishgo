package middleware

import (
	. "github.com/fishedee/app/router"
	. "github.com/fishedee/language"
	"reflect"
	"strings"
)

func isPublic(name string) bool {
	fisrtStr := name[0:1]
	if fisrtStr >= "A" && fisrtStr <= "Z" {
		return true
	} else {
		return false
	}
}

func analyseObjectRouter(factory *RouterFactory, name string) (func(path string, handler interface{}) *RouterFactory, string) {
	methodFunction := map[string]func(path string, handler interface{}) *RouterFactory{
		"head":    factory.HEAD,
		"options": factory.OPTIONS,
		"get":     factory.GET,
		"post":    factory.POST,
		"put":     factory.PUT,
		"delete":  factory.DELETE,
		"patch":   factory.PATCH,
		"any":     factory.Any,
		"none":    nil,
	}
	method := factory.Any
	functionName := "/"

	nameArray := Explode(strings.ToLower(name), "_")
	router, isExist := methodFunction[nameArray[0]]
	if isExist == true {
		method = router
		if len(nameArray) == 1 {
			return method, functionName
		} else {
			return method, nameArray[1]
		}
	} else {
		functionName = nameArray[0]
		return method, functionName
	}
}

func ObjectRouter(factory *RouterFactory, basePath string, handler interface{}) {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()
	methodNum := handlerValue.NumMethod()
	for i := 0; i != methodNum; i++ {
		methodHandler := handlerType.Method(i)
		methodName := methodHandler.Name
		if isPublic(methodName) == false {
			continue
		}
		addRouter, path := analyseObjectRouter(factory, methodName)
		if addRouter == nil {
			continue
		}
		addRouter(basePath+"/"+path, RouterMiddlewareContext{
			Data: map[string]interface{}{
				"name": methodName,
			},
			Handler: handlerValue.Method(i).Interface(),
		})
	}
}
