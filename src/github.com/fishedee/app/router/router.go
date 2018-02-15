package router

import (
	. "github.com/fishedee/container"
	. "github.com/fishedee/language"
	"net/http"
	"reflect"
	"strings"
)

type Router struct {
	trie      *TrieArray
	methodMap map[string]int
}

type routerHandlerFunc func(w http.ResponseWriter, r *http.Request, param map[string]string) int

type routerUrlPrefixHandler struct {
	param   map[int]string
	handler routerHandlerFunc
}

type routerHandler struct {
	urlExactHandler       routerHandlerFunc
	urlPrefixHandler      map[int]routerUrlPrefixHandler
	staticPrefixHandler   routerHandlerFunc
	notFoundPrefixHandler routerHandlerFunc
}

type routerPathInfo []routerHandler

func newRouter(trieTree *TrieTree) *Router {
	router := &Router{}
	router.methodMap = map[string]int{}
	entrys := routerMethod.Entrys()
	for i := routerMethod.HEAD; i <= routerMethod.PATCH; i++ {
		router.methodMap[entrys[i]] = i
	}
	router.trie = router.build(trieTree)
	return router
}

func (this *Router) combineParent(current routerFactoryPathInfo, parent routerFactoryPathInfo) {
	for method, parentMethodInfo := range parent {
		currentMethodInfo, isExist := current[method]
		if isExist == false {
			current[method] = currentMethodInfo
			continue
		}
		//合并urlPrefixHandler
		for seg, parentUrlPrefixHandler := range parentMethodInfo.urlPrefixHandler {
			_, isExist := currentMethodInfo.urlPrefixHandler[seg]
			if isExist == false {
				currentMethodInfo.urlPrefixHandler[seg] = parentUrlPrefixHandler
			}
		}
		//合并staticPrefixHandler
		if currentMethodInfo.staticPrefixHandler == nil {
			currentMethodInfo.staticPrefixHandler = parentMethodInfo.staticPrefixHandler
		}
		//合并notFoundPrefixHandler
		if currentMethodInfo.notFoundPrefixHandler == nil {
			currentMethodInfo.notFoundPrefixHandler = parentMethodInfo.notFoundPrefixHandler
		}
	}
}

type routerResponseWriter struct {
	writer http.ResponseWriter
	omit   bool
	status int
}

func newRouterResponseWriter(writer http.ResponseWriter) *routerResponseWriter {
	result := &routerResponseWriter{}
	result.writer = writer
	result.omit = false
	result.status = 200
	return result
}

func (this *routerResponseWriter) Header() http.Header {
	return this.writer.Header()
}

func (this *routerResponseWriter) Write(data []byte) (int, error) {
	if this.omit {
		return len(data), nil
	}
	return this.writer.Write(data)
}

func (this *routerResponseWriter) WriteHeader(status int) {
	this.status = status
	if status == 404 {
		this.omit = true
		return
	}
	this.writer.WriteHeader(status)
}

func (this *routerResponseWriter) GetStatus() int {
	return this.status
}

func (this *Router) catchNotFound(in interface{}, isCatch bool) routerHandlerFunc {
	origin := in.(routerFactoryHandlerFunc)
	if isCatch == false {
		return func(w http.ResponseWriter, r *http.Request, param map[string]string) int {
			origin(w, r, param)
			return 200
		}
	} else {
		return func(w http.ResponseWriter, r *http.Request, param map[string]string) int {
			fakeWriter := newRouterResponseWriter(w)
			origin(w, r, param)
			return fakeWriter.GetStatus()
		}
	}
}

func (this *Router) changeMethod(origin *routerFactoryHandler) routerHandler {
	result := routerHandler{
		urlExactHandler:       nil,
		urlPrefixHandler:      map[int]routerUrlPrefixHandler{},
		staticPrefixHandler:   nil,
		notFoundPrefixHandler: nil,
	}
	if origin.urlExactHandler != nil {
		result.urlExactHandler = this.catchNotFound(origin.urlExactHandler, false)
	}
	if origin.urlPrefixHandler != nil {
		for seg, singleUrlPrefixHandler := range origin.urlPrefixHandler {
			result.urlPrefixHandler[seg] = routerUrlPrefixHandler{
				param:   singleUrlPrefixHandler.param,
				handler: this.catchNotFound(singleUrlPrefixHandler.handler, false),
			}
		}
	}
	if origin.staticPrefixHandler != nil {
		result.staticPrefixHandler = this.catchNotFound(origin.staticPrefixHandler, true)
	}
	if origin.notFoundPrefixHandler != nil {
		result.notFoundPrefixHandler = this.catchNotFound(origin.notFoundPrefixHandler, false)
	}
	return result
}

func (this *Router) change(pathInfo routerFactoryPathInfo) routerPathInfo {
	methodLen := routerMethod.PATCH - routerMethod.HEAD + 2
	var result routerPathInfo
	result = make([]routerHandler, methodLen, methodLen)
	for i := 0; i != len(result); i++ {
		methodPathInfo := pathInfo[i]
		result[i] = this.changeMethod(methodPathInfo)
	}
	return result
}

func (this *Router) build(trieTree *TrieTree) *TrieArray {
	myTrieTree := NewTrieTree()
	trieTree.Walk(func(key string, value interface{}, parentKey string, parentValue interface{}) {
		currentPathInfo := value.(routerFactoryPathInfo)
		parentPathInfo := parentValue.(routerFactoryPathInfo)
		this.combineParent(currentPathInfo, parentPathInfo)
		newPathInfo := this.change(currentPathInfo)
		myTrieTree.Set(key, newPathInfo)
	})
	return myTrieTree.ToTrieArray()
}

func (this *Router) findHandler(url string, method int) (routerHandlerFunc, map[string]string, routerHandlerFunc) {
	handlerKey, handlerValue := this.trie.LongestPrefixMatch(strings.ToLower(url))
	handler := handlerValue.(routerPathInfo)[method]
	if handler.urlExactHandler != nil && len(handlerKey) == len(url) {
		return handler.urlExactHandler, nil, handler.notFoundPrefixHandler
	}
	if len(handler.urlPrefixHandler) != 0 {
		urlSegment := Explode(url, "/")
		urlPrefixHandler, isExist := handler.urlPrefixHandler[len(urlSegment)]
		if isExist {
			urlParam := map[string]string{}
			for index, key := range urlPrefixHandler.param {
				urlParam[key] = urlSegment[index]
			}
			return urlPrefixHandler.handler, urlParam, handler.notFoundPrefixHandler
		}
	}
	if handler.staticPrefixHandler != nil {
		return handler.staticPrefixHandler, nil, handler.notFoundPrefixHandler
	}
	return handler.notFoundPrefixHandler, nil, handler.notFoundPrefixHandler
}

func (this *Router) ServeHttp(w http.ResponseWriter, r *http.Request) {
	url := strings.ToLower(r.URL.Path)
	method, isExist := this.methodMap[r.Method]
	if isExist == false {
		panic("unsupport method " + r.Method)
	}
	handler, param, notFoundHandler := this.findHandler(url, method)
	status := handler(w, r, param)
	if status == 404 {
		notFoundHandler(w, r, param)
	}
}

type RouterFactory struct {
	basePath   string
	middleware []RouterMiddleware
	tree       map[string]routerFactoryPathInfo
	group      []*RouterFactory
}

type RouterMiddleware func([]interface{}) interface{}

var routerMethod struct {
	EnumStruct
	HEAD    int `enum:"1,HEAD"`
	OPTIONS int `enum:"2,OPTIONS"`
	GET     int `enum:"3,GET"`
	POST    int `enum:"4,POST"`
	DELETE  int `enum:"5,DELETE"`
	PUT     int `enum:"6,PUT"`
	PATCH   int `enum:"7,PATCH"`
	ANY     int `enum:"8,ANY"`
}

type routerFactoryUrlPrefixHandler struct {
	param   map[int]string
	handler interface{}
}

type routerFactoryHandlerFunc func(w http.ResponseWriter, r *http.Request, param map[string]string)

type routerFactoryHandler struct {
	urlExactHandler       interface{}
	urlPrefixHandler      map[int]*routerFactoryUrlPrefixHandler
	staticPrefixHandler   interface{}
	notFoundPrefixHandler interface{}
}

type routerFactoryPathInfo map[int]*routerFactoryHandler

func NewRouterFactory() *RouterFactory {
	routerFactory := &RouterFactory{}
	routerFactory.basePath = "/"
	routerFactory.middleware = []RouterMiddleware{}
	routerFactory.tree = map[string]routerFactoryPathInfo{}
	routerFactory.group = []*RouterFactory{}
	routerFactory.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404 page not found By Fish"))
		w.WriteHeader(404)
	})
	return routerFactory
}

func (this *RouterFactory) Use(middleware RouterMiddleware) *RouterFactory {
	this.middleware = append(this.middleware, middleware)
	return this
}

func (this *RouterFactory) isPublic(name string) bool {
	fisrtStr := name[0:1]
	if fisrtStr >= "A" && fisrtStr <= "Z" {
		return true
	} else {
		return false
	}
}

func (this *RouterFactory) changeUrlPrefix(priority int, path string, handler interface{}) (int, string, interface{}) {
	//过滤非url逻辑
	if priority != 1 {
		return priority, path, handler
	}

	//过滤前缀url逻辑
	pathInfo := Explode(strings.ToLower(path), "/")

	var singlePathIndex = 0
	for ; singlePathIndex != len(pathInfo); singlePathIndex++ {
		if pathInfo[singlePathIndex][0] == ':' {
			break
		}
	}
	if singlePathIndex == len(pathInfo) {
		return priority, path, handler
	}

	//处理前缀url逻辑
	urlPrefixHandler := &routerFactoryUrlPrefixHandler{
		param:   map[int]string{},
		handler: handler,
	}
	path = "/" + Implode(pathInfo[0:singlePathIndex], "/")
	for ; singlePathIndex != len(pathInfo); singlePathIndex++ {
		if pathInfo[singlePathIndex][0] != ':' {
			panic("invalid path : " + path)
		}
		urlPrefixHandler.param[singlePathIndex] = pathInfo[singlePathIndex][1:]
	}
	return 2, path, urlPrefixHandler
}

func (this *RouterFactory) addSingleRoute(method int, priority int, path string, handler interface{}) {
	//处理path
	pathInfo := Explode(strings.ToLower(path), "/")
	path = "/" + Implode(pathInfo, "/")

	//处理特殊的url前缀逻辑
	priority, path, handler = this.changeUrlPrefix(priority, path, handler)

	treeInfo, isExist := this.tree[path]
	if isExist == false {
		treeInfo = routerFactoryPathInfo{}
		this.tree[path] = treeInfo
	}
	methodInfo, isExist := treeInfo[method]
	if isExist == false {
		methodInfo = &routerFactoryHandler{
			urlExactHandler:       nil,
			urlPrefixHandler:      map[int]*routerFactoryUrlPrefixHandler{},
			staticPrefixHandler:   nil,
			notFoundPrefixHandler: nil,
		}
		treeInfo[method] = methodInfo
	}
	if priority == 1 {
		methodInfo.urlExactHandler = handler
	} else if priority == 2 {
		methodInfo.urlPrefixHandler[len(pathInfo)] = handler.(*routerFactoryUrlPrefixHandler)
	} else if priority == 3 {
		methodInfo.staticPrefixHandler = handler
	} else if priority == 4 {
		methodInfo.notFoundPrefixHandler = handler
	}
}

func (this *RouterFactory) addRoute(method int, priority int, path string, handler interface{}) {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()
	if handlerType.Kind() == reflect.Func {
		this.addSingleRoute(method, priority, path, handler)
	} else if handlerType.Kind() == reflect.Ptr &&
		handlerType.Elem().Kind() == reflect.Struct {
		methodNum := handlerValue.NumMethod()
		for i := 0; i != methodNum; i++ {
			methodHandler := handlerType.Method(i)
			methodName := methodHandler.Name
			if this.isPublic(methodName) == false {
				continue
			}
			methodNameArray := Explode(methodName, "_")
			this.addSingleRoute(method, priority, path+"/"+methodNameArray[0], handlerValue.Method(i).Interface())
		}
	} else {
		panic("invalid handler type " + handlerType.String())
	}
}

func (this *RouterFactory) HEAD(path string, handler interface{}) *RouterFactory {
	this.addRoute(routerMethod.HEAD, 1, path, handler)
	return this
}

func (this *RouterFactory) OPTIONS(path string, handler interface{}) *RouterFactory {
	this.addRoute(routerMethod.OPTIONS, 1, path, handler)
	return this
}

func (this *RouterFactory) GET(path string, handler interface{}) *RouterFactory {
	this.addRoute(routerMethod.GET, 1, path, handler)
	return this
}

func (this *RouterFactory) POST(path string, handler interface{}) *RouterFactory {
	this.addRoute(routerMethod.POST, 1, path, handler)
	return this
}

func (this *RouterFactory) DELETE(path string, handler interface{}) *RouterFactory {
	this.addRoute(routerMethod.DELETE, 1, path, handler)
	return this
}

func (this *RouterFactory) PUT(path string, handler interface{}) *RouterFactory {
	this.addRoute(routerMethod.PUT, 1, path, handler)
	return this
}

func (this *RouterFactory) PATCH(path string, handler interface{}) *RouterFactory {
	this.addRoute(routerMethod.PATCH, 1, path, handler)
	return this
}

func (this *RouterFactory) Any(path string, handler interface{}) *RouterFactory {
	for i := routerMethod.HEAD; i <= routerMethod.PATCH; i++ {
		this.addRoute(i, 1, path, handler)
	}
	return this
}

func (this *RouterFactory) rejustPath(path string) string {
	pathInfo := Explode(strings.ToLower(path), "/")
	newPath := "/" + Implode(pathInfo, "/")
	return newPath
}

func (this *RouterFactory) Static(path string, dir string) *RouterFactory {
	absolutePath := this.rejustPath(this.basePath + "/" + path)
	handler := http.StripPrefix(absolutePath, http.FileServer(http.Dir(dir)))
	this.addRoute(routerMethod.HEAD, 3, path, handler)
	this.addRoute(routerMethod.GET, 3, path, handler)
	return this
}

func (this *RouterFactory) NotFound(handler interface{}) *RouterFactory {
	for i := routerMethod.HEAD; i <= routerMethod.PATCH; i++ {
		this.addRoute(i, 4, "/", handler)
	}
	return this
}

func (this *RouterFactory) Group(basePath string, handler func(r *RouterFactory)) *RouterFactory {
	groupFactory := NewRouterFactory()

	groupFactory.basePath = this.rejustPath(this.basePath + "/" + basePath)

	this.group = append(this.group, groupFactory)

	handler(groupFactory)

	return this
}

func (this *RouterFactory) createHandler(middlewares []RouterMiddleware, handler interface{}) routerFactoryHandlerFunc {
	middlewares = append(middlewares, NewNoParamMiddleware())
	allHandler := []interface{}{handler}
	for i := len(middlewares) - 1; i >= 0; i-- {
		curHandler := middlewares[i](allHandler)
		allHandler = append(allHandler, curHandler)
	}
	resultHandler := allHandler[len(allHandler)-1]
	httpHandler, isOk := resultHandler.(routerFactoryHandlerFunc)
	if isOk == false {
		panic("handler must be routerFactoryHandlerFunc type")
	}
	return httpHandler
}

func (this *RouterFactory) buildTrie(trieTree *TrieTree, rootMiddleware []RouterMiddleware) {
	middlewares := append(rootMiddleware, this.middleware...)

	for path, mapper := range this.tree {
		for _, methodMapper := range mapper {
			if methodMapper.urlExactHandler != nil {
				methodMapper.urlExactHandler = this.createHandler(middlewares, methodMapper.urlExactHandler)
			}
			for _, urlPrefixHandler := range methodMapper.urlPrefixHandler {
				urlPrefixHandler.handler = this.createHandler(middlewares, urlPrefixHandler.handler)
			}
			if methodMapper.staticPrefixHandler != nil {
				methodMapper.staticPrefixHandler = this.createHandler(middlewares, methodMapper.staticPrefixHandler)
			}
			if methodMapper.notFoundPrefixHandler != nil {
				methodMapper.notFoundPrefixHandler = this.createHandler(middlewares, methodMapper.notFoundPrefixHandler)
			}
		}
		absolutePath := this.rejustPath(this.basePath + "/" + path)
		trieTree.Set(absolutePath, mapper)
	}

	for _, singleGroup := range this.group {
		singleGroup.buildTrie(trieTree, middlewares)
	}
}

func (this *RouterFactory) Create() *Router {
	trieTree := NewTrieTree()
	this.buildTrie(trieTree, []RouterMiddleware{})
	router := newRouter(trieTree)
	return router
}

func init() {
	InitEnumStruct(&routerMethod)
}
