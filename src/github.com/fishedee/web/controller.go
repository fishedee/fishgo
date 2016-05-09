package web

import (
	"bytes"
	"net/http"
)

type ControllerInterface interface {
	ModelInterface
	initEmpty(ControllerInterface, interface{})
	init(ControllerInterface, interface{}, interface{}, interface{})
	AutoRender(interface{}, string)
}

type Controller struct {
	*Basic
	appController ControllerInterface
}

type memoryResponseWriter struct {
	header     http.Header
	headerCode int
	data       []byte
}

func (this *memoryResponseWriter) Header() http.Header {
	if this.header == nil {
		this.header = http.Header{}
	}
	return this.header
}

func (this *memoryResponseWriter) Write(in []byte) (int, error) {
	this.data = append(this.data, in...)
	return len(this.data), nil
}

func (this *memoryResponseWriter) WriteHeader(headerCode int) {
	this.headerCode = headerCode
}

func (this *Controller) initEmpty(target ControllerInterface, t interface{}) {
	request, err := http.NewRequest("get", "/", bytes.NewReader([]byte("")))
	if err != nil {
		panic(err)
	}
	response := &memoryResponseWriter{}
	this.init(target, request, response, t)
}

func (this *Controller) init(target ControllerInterface, request interface{}, response interface{}, t interface{}) {
	this.appController = target
	this.Basic = initBasic(request, response, t)
	initModel(this.appController)
}

func (this *Controller) SetAppController(controller ControllerInterface) {
	if this.appController != nil {
		return
	}
	this.appController = controller
	this.Basic = controller.GetBasic().(*Basic)
}

func (this *Controller) GetBasic() interface{} {
	return this.Basic
}

func (this *Controller) AutoRender(result interface{}, view string) {

}

func (this *Controller) Check(requireStruct interface{}) {
	this.Ctx.GetParamToStruct(requireStruct)
}

func (this *Controller) CheckGet(requireStruct interface{}) {
	this.Ctx.GetUrlParamToStruct(requireStruct)
}

func (this *Controller) CheckPost(requireStruct interface{}) {
	this.Ctx.GetFormParamToStruct(requireStruct)
}

func (this *Controller) Write(data []byte) {
	this.Ctx.Write(data)
}

func (this *Controller) WriteHeader(key string, value string) {
	this.Ctx.WriteHeader(key, value)
}

func (this *Controller) WriteMimeHeader(mime string, title string) {
	this.Ctx.WriteMimeHeader(mime, title)
}

func InitController(controller interface{}) {
}
