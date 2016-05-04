package web

import (
	"bytes"
	_ "github.com/a"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/fishedee/language"
	"io/ioutil"
	"reflect"
	"strings"
)

var (
	routerControllerMethod = map[string]map[string]methodInfo{}
)

func firstLowerName(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func firstUpperName(name string) string {
	return strings.ToUpper(name[0:1]) + name[1:]
}

func isPublic(name string) bool {
	fisrtStr := name[0:1]
	if fisrtStr >= "A" && fisrtStr <= "Z" {
		return true
	} else {
		return false
	}
}

func runBeegoRequest(controller reflect.Value, method methodInfo, ctx *context.Context) {
	urlMethod := ctx.Input.Method()
	target := controller.Interface().(ControllerInterface)
	defer language.CatchCarsh(func(exception language.Exception) {
		target.Log.Error("Buiness Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
		result = []reflect.Value{reflect.ValueOf(exception)}
	})
	target.Init(controller, ctx.Request, ctx.ResponseWriter, nil)
	var controllerResult interface{}
	if urlMethod == "GET" || urlMethod == "POST" ||
		urlMethod == "DELETE" || urlMethod == "PUT" {
		result := runBeegoBusinessRequest(target, method.methodType, []reflect.Value{controller})
		if len(result) != 1 {
			panic("url controller should has return value " + url)
		}
		controllerResult = result[0].Interface()
	} else {
		controllerResult = nil
	}
	target.AutoRender(controllerResult, method.viewName)
}

func runBeegoBusinessRequest(target ControllerInterface, method reflect.Value, arguments []reflect.Value) (result []reflect.Value) {
	defer language.Catch(func(exception language.Exception) {
		target.Log.Error("Buiness Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
		result = []reflect.Value{reflect.ValueOf(exception)}
	})
	result = method.Call(arguments)
	if len(result) == 0 {
		result = []reflect.Value{reflect.Zero(reflect.TypeOf(ControllerInterface))}
	}
	return
}

func handleBeegoRequest(ctx *context.Context) {
	//查找路由
	url := ctx.Input.URL()
	urlArray := language.Explode(url, "/")
	if len(urlArray) < 2 {
		ctx.Abort("404", "File Not Found")
		return
	}
	method, isExist := routerControllerMethod[strings.ToLower(urlArray[0])][strings.ToLower(urlArray[1])]
	if isExist == false {
		ctx.Abort("404", "File Not Found")
		return
	}

	//设置http的body
	if len(ctx.Input.RequestBody) != 0 {
		ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(ctx.Input.RequestBody))
	}

	//执行路由
	controller := reflect.New(method.controllerType)
	runBeegoRequest(controller, method, ctx)
}

type methodInfo struct {
	controllerName string
	methodName     string
	viewName       string
	controllerType reflect.Type
	methodType     reflect.Method
}

func InitRoute(namespace string, target interface{}) {
	controllerType := reflect.TypeOf(target)
	routerControllerMethod[controllerType] = map[string]methodInfo{}

	numMethod := controllerType.NumMethod()
	for i := 0; i != numMethod; i++ {
		singleMethod := controllerType.Method(i)
		singleMethodName := singleMethod.Name
		if isPublic(singleMethodName) == false {
			continue
		}
		methodNameInfo := language.Explode(singleMethodName, "_")
		if len(methodNameInfo) < 2 {
			continue
		}
		routerControllerMethod[singleMethod.Name][singleMethodInfo.name] = &methodInfo{
			controllerName: strings.ToLower(controllerType.Name()),
			methodName:     strings.ToLower(methodNameInfo[0]),
			viewName:       firstLowerName(methodNameInfo[1]),
			controllerType: controllerType,
			methodType:     singleMethod,
		}
	}
}

func Run() {
	beego.Any("*.*", handleBeegoRequest)
	beego.Run()
}
