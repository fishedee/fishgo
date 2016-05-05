package web

import (
	"bytes"
	"github.com/fishedee/encoding"
	"github.com/fishedee/language"
	"io/ioutil"
	"mime"
	"net/http"
	"testing"
)

type ControllerInterface interface {
	ModelInterface
	InitEmpty(ControllerInterface, *testing.T)
	Init(ControllerInterface, *http.Request, http.ResponseWriter, *testing.T)
	AutoRender(interface{}, string)
}

type Controller struct {
	*Basic
	appController ControllerInterface
	inputData     interface{}
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

func (this *Controller) InitEmpty(target ControllerInterface, t *testing.T) {
	request, err := http.NewRequest("get", "/", bytes.NewReader([]byte("")))
	if err != nil {
		panic(err)
	}
	response := &memoryResponseWriter{}
	this.Init(target, request, response, t)
}

func (this *Controller) Init(target ControllerInterface, request *http.Request, response http.ResponseWriter, t *testing.T) {
	this.appController = target
	this.Basic = initBasic(request, response, t)
	initModel(this.appController)
	this.parseInput()
}

func (this *Controller) SetAppController(controller ControllerInterface) {
	if this.appController != nil {
		return
	}
	this.appController = controller
	this.Basic = controller.GetBasic()
}

func (this *Controller) GetBasic() *Basic {
	return this.Basic
}

func (this *Controller) AutoRender(result interface{}, view string) {

}

func (this *Controller) parseInput() {
	//取出get数据
	request := this.Ctx.Request
	queryInput := request.URL.RawQuery

	//取出post数据
	postInput := ""
	ct := request.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	ct, _, err := mime.ParseMediaType(ct)
	if ct == "application/x-www-form-urlencoded" {
		byteArray, _ := ioutil.ReadAll(this.Ctx.Request.Body)
		this.Ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(byteArray))
		postInput = string(byteArray)
	}

	//解析数据
	input := queryInput + "&" + postInput
	this.inputData = nil
	err = encoding.DecodeUrlQuery([]byte(input), &this.inputData)
	if err != nil {
		language.Throw(1, err.Error())
	}
}

func (this *Controller) Check(requireStruct interface{}) {
	//导出到struct
	err := language.MapToArray(this.inputData, requireStruct, "url")
	if err != nil {
		language.Throw(1, err.Error())
	}
}

func (this *Controller) CheckGet(requireStruct interface{}) {
	if this.Ctx.Request.Method != "GET" {
		language.Throw(1, "请求Method不是Get方法")
	}
	this.Check(requireStruct)
}

func (this *Controller) CheckPost(requireStruct interface{}) {
	if this.Ctx.Request.Method != "POST" {
		language.Throw(1, "请求Method不是POST方法")
	}
	this.Check(requireStruct)
}

func (this *Controller) Write(data []byte) {
	writer := this.Ctx.ResponseWriter
	writer.Write(data)
}

func (this *Controller) WriteHeader(key string, value string) {
	writer := this.Ctx.ResponseWriter
	writer.Header().Set(key, value)
}

func (this *Controller) WriteMimeHeader(mime string, title string) {
	writer := this.Ctx.ResponseWriter
	writerHeader := writer.Header()
	if mime == "json" {
		writerHeader.Set("Content-Type", "application/x-javascript; charset=utf-8")
	} else if mime == "javascript" {
		writerHeader.Set("Content-Type", "application/x-javascript; charset=utf-8")
	} else if mime == "plain" {
		writerHeader.Set("Content-Type", "text/plain; charset=utf-8")
	} else if mime == "xlsx" {
		writerHeader.Set("Content-Type", "application/vnd.openxmlformats-officedocument; charset=UTF-8")
		writerHeader.Set("Pragma", "public")
		writerHeader.Set("Expires", "0")
		writerHeader.Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
		writerHeader.Set("Content-Type", "application/force-download")
		writerHeader.Set("Content-Type", "application/octet-stream")
		writerHeader.Set("Content-Type", "application/download")
		writerHeader.Set("Content-Disposition", "attachment;filename="+title+".xlsx")
		writerHeader.Set("Content-Transfer-Encoding", "binary")
	} else if mime == "csv" {
		writerHeader.Set("Content-Type", "application/vnd.ms-excel; charset=UTF-8")
		writerHeader.Set("Pragma", "public")
		writerHeader.Set("Expires", "0")
		writerHeader.Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
		writerHeader.Set("Content-Type", "application/force-download")
		writerHeader.Set("Content-Type", "application/octet-stream")
		writerHeader.Set("Content-Type", "application/download")
		writerHeader.Set("Content-Disposition", "attachment;filename="+title+".csv")
		writerHeader.Set("Content-Transfer-Encoding", "binary")
	} else {
		panic("invalid mime [" + mime + "]")
	}
}
