package router

import (
	. "github.com/fishedee/container"
	. "github.com/fishedee/language"
	"net/http"
	"strings"
	"sync"
)

type Router struct {
	trie *TrieArray
	pool sync.Pool
}

type RouterSingleParam struct {
	Key   string
	Value string
}
type RouterParam []RouterSingleParam

type routerContext struct {
	param RouterParam
}

type routerHandlerFunc func(w http.ResponseWriter, r *http.Request, param RouterParam) int

type routerUrlPrefixHandlerParam struct {
	index int
	name  string
}
type routerUrlPrefixHandler struct {
	segment int
	param   []routerUrlPrefixHandlerParam
	handler routerHandlerFunc
}

type routerHandler struct {
	prefix                []string
	prefixLength          int
	urlExactHandler       routerHandlerFunc
	urlPrefixHandler      []routerUrlPrefixHandler
	staticPrefixHandler   routerHandlerFunc
	notFoundPrefixHandler routerHandlerFunc
}

type routerPathInfo []routerHandler

func newRouter(trieTree *TrieTree, maxSegment int) *Router {
	router := &Router{}
	router.trie = router.build(trieTree)
	router.pool = sync.Pool{
		New: func() interface{} {
			return &routerContext{
				param: make([]RouterSingleParam, maxSegment, maxSegment),
			}
		},
	}
	return router
}

func (this *Router) combineParent(current routerFactoryPathInfo, parent routerFactoryPathInfo) {
	for method, parentMethodInfo := range parent {
		currentMethodInfo, isExist := current[method]
		if isExist == false {
			current[method] = parentMethodInfo
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
		return func(w http.ResponseWriter, r *http.Request, param RouterParam) int {
			origin(w, r, param)
			return 200
		}
	} else {
		return func(w http.ResponseWriter, r *http.Request, param RouterParam) int {
			fakeWriter := newRouterResponseWriter(w)
			origin(fakeWriter, r, param)
			return fakeWriter.GetStatus()
		}
	}
}

func (this *Router) changeMethod(origin *routerFactoryHandler) routerHandler {
	result := routerHandler{
		urlExactHandler:       nil,
		urlPrefixHandler:      []routerUrlPrefixHandler{},
		staticPrefixHandler:   nil,
		notFoundPrefixHandler: nil,
	}
	if origin.urlExactHandler != nil {
		result.urlExactHandler = this.catchNotFound(origin.urlExactHandler, false)
	}
	if origin.urlPrefixHandler != nil {
		resultUrlPrefixHandler := []routerUrlPrefixHandler{}
		for seg, singleUrlPrefixHandler := range origin.urlPrefixHandler {
			param := []routerUrlPrefixHandlerParam{}
			for key, value := range singleUrlPrefixHandler.param {
				param = append(param, routerUrlPrefixHandlerParam{
					index: key,
					name:  value,
				})
			}
			param = QuerySort(param, "index asc").([]routerUrlPrefixHandlerParam)
			resultUrlPrefixHandler = append(resultUrlPrefixHandler, routerUrlPrefixHandler{
				segment: seg,
				param:   param,
				handler: this.catchNotFound(singleUrlPrefixHandler.handler, false),
			})
		}
		result.urlPrefixHandler = QuerySort(resultUrlPrefixHandler, "segment asc").([]routerUrlPrefixHandler)
	}
	if origin.staticPrefixHandler != nil {
		result.staticPrefixHandler = this.catchNotFound(origin.staticPrefixHandler, true)
	}
	if origin.notFoundPrefixHandler != nil {
		result.notFoundPrefixHandler = this.catchNotFound(origin.notFoundPrefixHandler, false)
	}
	return result
}

func (this *Router) changeUrl(url string) ([]string, int) {
	length := strings.LastIndexByte(url, '/')
	if length == -1 {
		return nil, 0
	}
	url = url[:length]
	return Explode(url, "/"), length
}

func (this *Router) change(url string, pathInfo routerFactoryPathInfo) routerPathInfo {
	methodLen := RouterMethod.PATCH - RouterMethod.HEAD + 2
	var result routerPathInfo
	result = make([]routerHandler, methodLen, methodLen)
	for i := RouterMethod.HEAD; i <= RouterMethod.PATCH; i++ {
		methodPathInfo := pathInfo[i]
		newMethodPathInfo := this.changeMethod(methodPathInfo)
		newMethodPathInfo.prefix, newMethodPathInfo.prefixLength = this.changeUrl(url)
		result[i] = newMethodPathInfo
	}
	return result
}

func (this *Router) build(trieTree *TrieTree) *TrieArray {
	myTrieTree := NewTrieTree()
	trieTree.Walk(func(key string, value interface{}, parentKey string, parentValue interface{}) {
		if value == nil {
			value = routerFactoryPathInfo{}
		}
		if parentValue == nil {
			parentValue = routerFactoryPathInfo{}
		}
		currentPathInfo := value.(routerFactoryPathInfo)
		parentPathInfo := parentValue.(routerFactoryPathInfo)
		this.combineParent(currentPathInfo, parentPathInfo)
		trieTree.Set(key, currentPathInfo)
		newPathInfo := this.change(key, currentPathInfo)
		myTrieTree.Set(key, newPathInfo)
	})
	return myTrieTree.ToTrieArray()
}

func (this *Router) normalUrl(url string, a int) string {
	urlLen := len(url)
	begin := 0
	end := urlLen
	if urlLen > 0 && url[0] == '/' {
		begin++
	}
	if urlLen >= 2 && url[end-1] == '/' {
		end--
	}
	return url[begin:end]
}

func (this *Router) parseParam(prefix []string, suffix string, maxParam int) (RouterParam, *routerContext, bool) {
	context := this.pool.Get().(*routerContext)
	param := context.param
	k := 0
	for i := 0; i != len(prefix); i++ {
		param[k].Value = prefix[i]
		k++
	}
	lastIndex := -1
	for lastIndex < len(suffix) {
		index := lastIndex + 1
		for index < len(suffix) && suffix[index] != '/' {
			index++
		}
		if index-lastIndex > 1 {
			if k >= maxParam {
				this.pool.Put(context)
				return nil, nil, false
			}
			param[k].Value = suffix[lastIndex+1 : index]
			k++
		}
		lastIndex = index
	}
	return param[0:k], context, true
}

func (this *Router) findHandler(url string, method int) (routerHandlerFunc, RouterParam, routerHandlerFunc, *routerContext) {
	searchUrl := this.normalUrl(url, 1)
	var isExact bool
	var handlerValue interface{}
	if len(searchUrl) != 0 {
		handlerValue, isExact = this.trie.LongestPrefixMatchWithChar(searchUrl, '/')
	} else {
		var handlerKey string
		handlerKey, handlerValue = this.trie.LongestPrefixMatch(searchUrl)
		isExact = len(handlerKey) == len(searchUrl)
	}
	handler := handlerValue.(routerPathInfo)[method]
	if handler.urlExactHandler != nil && isExact {
		return handler.urlExactHandler, nil, handler.notFoundPrefixHandler, nil
	}
	if len(handler.urlPrefixHandler) != 0 {
		maxParam := handler.urlPrefixHandler[len(handler.urlPrefixHandler)-1].segment
		suffixUrl := searchUrl[handler.prefixLength:]
		param, context, isValid := this.parseParam(handler.prefix, suffixUrl, maxParam)
		if isValid {
			for _, urlPrefixHandler := range handler.urlPrefixHandler {
				if urlPrefixHandler.segment == len(param) {
					for _, singleParam := range urlPrefixHandler.param {
						param[singleParam.index].Key = singleParam.name
					}
					begin := urlPrefixHandler.param[0].index
					end := len(param)
					return urlPrefixHandler.handler, param[begin:end], handler.notFoundPrefixHandler, context
				}
			}
		}
	}
	if handler.staticPrefixHandler != nil {
		return handler.staticPrefixHandler, nil, handler.notFoundPrefixHandler, nil
	}
	return handler.notFoundPrefixHandler, nil, handler.notFoundPrefixHandler, nil
}

func (this *Router) ServeHttp(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	var methodInt int
	switch r.Method {
	case "HEAD":
		methodInt = RouterMethod.HEAD
		break
	case "OPTIONS":
		methodInt = RouterMethod.OPTIONS
		break
	case "GET":
		methodInt = RouterMethod.GET
		break
	case "POST":
		methodInt = RouterMethod.POST
		break
	case "DELETE":
		methodInt = RouterMethod.DELETE
		break
	case "PUT":
		methodInt = RouterMethod.PUT
		break
	case "PATCH":
		methodInt = RouterMethod.PATCH
	default:
		panic("unsupport method " + r.Method)
	}
	handler, param, notFoundHandler, context := this.findHandler(url, methodInt)
	status := handler(w, r, param)
	if status == 404 {
		notFoundHandler(w, r, param)
	}
	if context != nil {
		this.pool.Put(context)
	}
}

type RouterFactory struct {
	basePath   string
	middleware []RouterMiddleware
	tree       map[string]routerFactoryPathInfo
	group      []*RouterFactory
	maxSegment int
}

type RouterMiddleware func([]interface{}) interface{}

var RouterMethod struct {
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

type routerFactoryHandlerFunc func(w http.ResponseWriter, r *http.Request, param RouterParam)

type routerFactoryHandler struct {
	urlExactHandler       interface{}
	urlPrefixHandler      map[int]*routerFactoryUrlPrefixHandler
	staticPrefixHandler   interface{}
	notFoundPrefixHandler interface{}
}

type routerFactoryPathInfo map[int]*routerFactoryHandler

func NewRouterFactory() *RouterFactory {
	routerFactory := newRouterFactory("")
	routerFactory.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("404 page not found By Fish"))
	})
	return routerFactory
}

func newRouterFactory(basePath string) *RouterFactory {
	routerFactory := &RouterFactory{}
	routerFactory.basePath = basePath
	routerFactory.middleware = []RouterMiddleware{}
	routerFactory.tree = map[string]routerFactoryPathInfo{}
	routerFactory.group = []*RouterFactory{}
	routerFactory.maxSegment = 0
	return routerFactory
}

func (this *RouterFactory) Use(middleware RouterMiddleware) *RouterFactory {
	this.middleware = append(this.middleware, middleware)
	return this
}

func (this *RouterFactory) changeUrlPrefix(priority int, path string, handler interface{}) (int, string, interface{}) {
	//过滤非url逻辑
	if priority != 1 {
		return priority, path, handler
	}

	//过滤前缀url逻辑
	pathInfo := Explode(path, "/")

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
	path = Implode(pathInfo[0:singlePathIndex], "/")
	if len(path) != 0 {
		path += "/"
	}
	for ; singlePathIndex != len(pathInfo); singlePathIndex++ {
		if pathInfo[singlePathIndex][0] != ':' {
			panic("invalid path : " + path)
		}
		urlPrefixHandler.param[singlePathIndex] = pathInfo[singlePathIndex][1:]
	}

	if len(pathInfo) > this.maxSegment {
		this.maxSegment = len(pathInfo)
	}
	return 2, path, urlPrefixHandler
}

func (this *RouterFactory) addSingleRoute(method int, priority int, path string, handler interface{}) {
	//处理path
	pathInfo := Explode(this.basePath+"/"+path, "/")
	path = Implode(pathInfo, "/")
	if len(path) != 0 {
		path += "/"
	}

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
	this.addSingleRoute(method, priority, path, handler)
}

func (this *RouterFactory) HEAD(path string, handler interface{}) *RouterFactory {
	this.addRoute(RouterMethod.HEAD, 1, path, handler)
	return this
}

func (this *RouterFactory) OPTIONS(path string, handler interface{}) *RouterFactory {
	this.addRoute(RouterMethod.OPTIONS, 1, path, handler)
	return this
}

func (this *RouterFactory) GET(path string, handler interface{}) *RouterFactory {
	this.addRoute(RouterMethod.GET, 1, path, handler)
	return this
}

func (this *RouterFactory) POST(path string, handler interface{}) *RouterFactory {
	this.addRoute(RouterMethod.POST, 1, path, handler)
	return this
}

func (this *RouterFactory) DELETE(path string, handler interface{}) *RouterFactory {
	this.addRoute(RouterMethod.DELETE, 1, path, handler)
	return this
}

func (this *RouterFactory) PUT(path string, handler interface{}) *RouterFactory {
	this.addRoute(RouterMethod.PUT, 1, path, handler)
	return this
}

func (this *RouterFactory) PATCH(path string, handler interface{}) *RouterFactory {
	this.addRoute(RouterMethod.PATCH, 1, path, handler)
	return this
}

func (this *RouterFactory) Any(path string, handler interface{}) *RouterFactory {
	for i := RouterMethod.HEAD; i <= RouterMethod.PATCH; i++ {
		this.addRoute(i, 1, path, handler)
	}
	return this
}

func (this *RouterFactory) rejustPath(path string) string {
	pathInfo := Explode(path, "/")
	newPath := Implode(pathInfo, "/")
	return newPath
}

func (this *RouterFactory) Static(path string, dir string) *RouterFactory {
	absolutePath := this.rejustPath(this.basePath + "/" + path)
	handler := http.StripPrefix("/"+absolutePath, http.FileServer(http.Dir(dir)))
	this.addRoute(RouterMethod.HEAD, 3, path, handler)
	this.addRoute(RouterMethod.GET, 3, path, handler)
	return this
}

func (this *RouterFactory) NotFound(handler interface{}) *RouterFactory {
	for i := RouterMethod.HEAD; i <= RouterMethod.PATCH; i++ {
		this.addRoute(i, 4, "/", handler)
	}
	return this
}

func (this *RouterFactory) Group(basePath string, handler func(r *RouterFactory)) *RouterFactory {
	realBasePath := this.rejustPath(this.basePath + "/" + basePath)
	groupFactory := newRouterFactory(realBasePath)

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
	httpHandler, isOk := resultHandler.(func(w http.ResponseWriter, r *http.Request, param RouterParam))
	if isOk == false {
		panic("handler must be routerFactoryHandlerFunc type")
	}
	return routerFactoryHandlerFunc(httpHandler)
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
		trieTree.Set(path, mapper)
	}

	for _, singleGroup := range this.group {
		singleGroup.buildTrie(trieTree, middlewares)
	}
}

func (this *RouterFactory) getMaxSegment() int {
	maxSegment := this.maxSegment
	for _, singleGroup := range this.group {
		groupSegment := singleGroup.getMaxSegment()
		if groupSegment > maxSegment {
			maxSegment = groupSegment
		}
	}
	return maxSegment
}

func (this *RouterFactory) Create() *Router {
	trieTree := NewTrieTree()
	this.buildTrie(trieTree, []RouterMiddleware{})
	maxSegment := this.getMaxSegment()
	router := newRouter(trieTree, maxSegment)
	return router
}

func init() {
	InitEnumStruct(&RouterMethod)
}
