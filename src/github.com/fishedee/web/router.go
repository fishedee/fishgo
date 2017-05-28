package web

import (
	"github.com/fishedee/language"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type ControllerInterface interface {
	AutoRender(interface{}, string)
}

type methodInfo struct {
	viewName       string
	controllerType reflect.Type
	methodType     reflect.Method
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

func (this *handlerType) addRoute(namespace string, target ControllerInterface) {
	if this.routerControllerMethod == nil {
		this.routerControllerMethod = map[string]methodInfo{}
	}
	controllerType := reflect.TypeOf(target)
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
		var url string
		if namespace != "" {
			url = strings.ToLower(namespace + "/" + methodName[0])
		} else {
			url = strings.ToLower(methodName[0])
		}
		this.routerControllerMethod[url] = methodInfo{
			viewName:       this.firstLowerName(methodName[1]),
			controllerType: controllerType.Elem(),
			methodType:     singleMethod,
		}
	}
	//预热ioc
	injectIoc(reflect.ValueOf(target), nil)
}

func (this *handlerType) handleRequest(request *http.Request, response http.ResponseWriter) {
	//查找路由
	url := request.URL.Path
	url = strings.ToLower(strings.Trim(url, "/"))
	method, isExist := this.routerControllerMethod[url]
	if isExist == false {
		globalBasic.Log.Error("file not found : %s", url)
		response.WriteHeader(404)
		response.Write([]byte("file not found"))
		return
	}

	//执行路由
	controller := reflect.New(method.controllerType)
	this.runRequest(controller, method, request, response)
}

func (this *handlerType) runRequest(controller reflect.Value, method methodInfo, request *http.Request, response http.ResponseWriter) {
	urlMethod := request.Method
	basic := initBasic(request, response, nil)
	target := controller.Interface().(ControllerInterface)
	injectIoc(controller, basic)
	defer language.CatchCrash(func(exception language.Exception) {
		basic.Log.Critical("Buiness Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
		response.WriteHeader(500)
		response.Write([]byte("server internal error"))
	})
	var controllerResult interface{}
	if urlMethod == "GET" || urlMethod == "POST" ||
		urlMethod == "DELETE" || urlMethod == "PUT" {
		result := this.runRequestBusiness(target, method.methodType.Func, []reflect.Value{controller}, basic)
		if len(result) >= 1 {
			controllerResult = result[0].Interface()
		} else {
			controllerResult = nil
		}
	} else {
		controllerResult = nil
	}
	target.AutoRender(controllerResult, method.viewName)
}

func (this *handlerType) runRequestBusiness(target ControllerInterface, method reflect.Value, arguments []reflect.Value, basic *Basic) (result []reflect.Value) {
	defer language.Catch(func(exception language.Exception) {
		basic.Log.Error("Buiness Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
		result = []reflect.Value{reflect.ValueOf(exception)}
	})
	result = method.Call(arguments)
	return
}

var handler handlerType

func InitRoute(namespace string, target ControllerInterface) {
	handler.addRoute(namespace, target)
}

func Run() error {
	//启动服务器
	httpPort := globalBasic.Config.GetInt("httpport")
	if httpPort == 0 {
		httpPort = 8080
	}
	globalBasic.Log.Debug("Server is Running :%v", httpPort)
	err := globalBasic.Grace.ListenAndServe(httpPort, &handler)
	if err != nil {
		globalBasic.Log.Error("Listen fail! " + err.Error())
		return err
	}

	//删除收尾的资源
	destroyBasic()
	return nil
}
