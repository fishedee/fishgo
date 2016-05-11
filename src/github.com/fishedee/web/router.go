package web

import (
	"github.com/fishedee/language"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type methodInfo struct {
	viewName       string
	controllerType reflect.Type
	methodIndex    int
}

type handlerType struct {
	routerControllerMethod map[string]methodInfo
}

func (this *handlerType) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	beginTime := time.Now().UnixNano()
	this.handleRequest(request, response)
	endTime := time.Now().UnixNano()
	globalBasic.Log.Debug("%s %s : %s", request.Method, request.URL.String(), time.Duration(endTime-beginTime).String())
}

func (this *handlerType) firstLowerName(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func (this *handlerType) firstUpperName(name string) string {
	return strings.ToUpper(name[0:1]) + name[1:]
}

func (this *handlerType) isPublic(name string) bool {
	fisrtStr := name[0:1]
	if fisrtStr >= "A" && fisrtStr <= "Z" {
		return true
	} else {
		return false
	}
}

func (this *handlerType) addRoute(namespace string, target interface{}) {
	if this.routerControllerMethod == nil {
		this.routerControllerMethod = map[string]methodInfo{}
	}
	controllerValue := getIocRealTarget(target)
	controllerType := controllerValue.Type()
	numMethod := controllerType.NumMethod()
	for i := 0; i != numMethod; i++ {
		singleMethod := controllerType.Method(i)
		singleMethodName := singleMethod.Name
		if this.isPublic(singleMethodName) == false {
			continue
		}
		methodName := language.Explode(singleMethodName, "_")
		if len(methodName) < 2 {
			continue
		}
		namespace := strings.Trim(namespace, "/")
		methodName[0] = strings.Trim(methodName[0], "/")
		url := strings.ToLower(namespace + "/" + methodName[0])
		this.routerControllerMethod[url] = methodInfo{
			viewName:       this.firstLowerName(methodName[1]),
			controllerType: controllerType,
			methodIndex:    i,
		}
	}
}

func (this *handlerType) handleRequest(request *http.Request, response http.ResponseWriter) {
	//查找路由
	url := request.URL.Path
	url = strings.ToLower(strings.Trim(url, "/"))
	method, isExist := this.routerControllerMethod[url]
	if isExist == false {
		response.WriteHeader(404)
		response.Write([]byte("file not found"))
		return
	}

	//执行路由
	this.runRequest(method, request, response)
}

func (this *handlerType) runRequest(method methodInfo, request *http.Request, response http.ResponseWriter) {
	urlMethod := request.Method
	basic := initBasic(request, response, nil)
	defer language.CatchCrash(func(exception language.Exception) {
		basic.Log.Critical("Buiness Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
		response.WriteHeader(500)
		response.Write([]byte("server internal error"))
	})
	target, err := newIocInstanse(method.controllerType, basic)
	if err != nil {
		panic(err)
	}

	var controllerResult interface{}
	if urlMethod == "GET" || urlMethod == "POST" ||
		urlMethod == "DELETE" || urlMethod == "PUT" {
		result := this.runRequestBusiness(basic, target, method.methodIndex)
		if len(result) >= 1 {
			controllerResult = result[0].Interface()
		} else {
			controllerResult = nil
		}
	} else {
		controllerResult = nil
	}
	target.Interface().(ControllerInterface).AutoRender(controllerResult, method.viewName)
}

func (this *handlerType) runRequestBusiness(basic *Basic, target reflect.Value, methodIndex int) (result []reflect.Value) {
	defer language.Catch(func(exception language.Exception) {
		basic.Log.Error("Buiness Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
		result = []reflect.Value{reflect.ValueOf(exception)}
	})
	result = target.Method(methodIndex).Call(nil)
	return
}

var handler handlerType

func InitRoute(namespace string, target interface{}) {
	handler.addRoute(namespace, target)
}

func Run() error {
	httpPort := globalBasic.Config.GetInt("httpport")
	if httpPort == 0 {
		httpPort = 8080
	}
	globalBasic.Log.Debug("Server is Running :%v", httpPort)
	err := http.ListenAndServe(":"+strconv.Itoa(httpPort), &handler)
	if err != nil {
		globalBasic.Log.Error("Listen fail! " + err.Error())
		return err
	}
	return nil
}
