package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/fishedee/encoding"
	"github.com/fishedee/language"
)

type ContextSerializeRequest struct {
	Method string
	Url    string
	Header map[string][]string
}

type Context interface {
	//输入参数
	GetUrl() *url.URL
	GetBody() ([]byte, error)
	GetParam(key string) string
	SetParam(key string, data string)
	GetParamToStruct(requireStruct interface{})
	GetUrlParamToStruct(requireStruct interface{})
	GetFormParamToStruct(requireStruct interface{})
	GetFile(key string) ([]byte, string, error)
	GetCookie(key string) string

	//元信息
	GetProxy() []string
	GetMethod() string
	GetScheme() string
	GetHost() string
	GetPort() int
	GetSite() string
	GetRemoteAddr() string
	GetRemoteIP() string
	GetRemotePort() int
	GetUserAgent() string
	SetUserAgent(data string)
	GetHeader(key string) string
	IsUpload() bool
	IsLocal() bool

	//输出数据
	Write(data []byte)
	WriteHeader(key string, value string)
	WriteStatus(code int)
	WriteMimeHeader(mime string, title string)

	//危险操作
	GetRawRequest() interface{}
	GetRawResponseWriter() interface{}
	GetRawTesting() interface{}

	//序列化与反序列化
	SerializeRequest() (ContextSerializeRequest, error)
	DeSerializeRequest(ContextSerializeRequest) error
}

type contextImplement struct {
	request          *http.Request
	responseWriter   http.ResponseWriter
	testing          *testing.T
	inputData        map[string]interface{}
	serializeRequest *ContextSerializeRequest
}

func NewContext(request interface{}, response interface{}, t interface{}) Context {
	if t == nil {
		t = (*testing.T)(nil)
	}
	result := contextImplement{
		request:        request.(*http.Request),
		responseWriter: response.(http.ResponseWriter),
		testing:        t.(*testing.T),
	}
	result.parseInput()
	return &result
}

func (this *contextImplement) parseInput() {
	//取出get数据
	request := this.request
	queryInput := request.URL.RawQuery

	//取出post数据
	postInput := ""
	ct := request.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	ct, _, err := mime.ParseMediaType(ct)
	if ct == "application/x-www-form-urlencoded" &&
		this.request.Body != nil {
		byteArray, err := ioutil.ReadAll(this.request.Body)
		if err != nil {
			panic(err)
		}
		this.request.Body = ioutil.NopCloser(bytes.NewReader(byteArray))
		postInput = string(byteArray)
	}

	//解析数据
	input := queryInput + "&" + postInput
	this.inputData = map[string]interface{}{}
	err = encoding.DecodeUrlQuery([]byte(input), &this.inputData)
	if err != nil {
		language.Throw(1, err.Error())
	}
}

func (this *contextImplement) GetUrl() *url.URL {
	return this.request.URL
}

func (this *contextImplement) GetBody() ([]byte, error) {
	byteArray, err := ioutil.ReadAll(this.request.Body)
	if err != nil {
		return nil, err
	}
	return byteArray, nil
}

func (this *contextImplement) GetParam(key string) string {
	result, isExist := this.inputData[key]
	if isExist {
		return fmt.Sprintf("%v", result)
	}

	return this.request.FormValue(key)
}

func (this *contextImplement) SetParam(key string, data string) {
	this.inputData[key] = data
}

func (this *contextImplement) GetParamToStruct(requireStruct interface{}) {
	//导出到struct
	err := language.MapToArray(this.inputData, requireStruct, "url")
	if err != nil {
		language.Throw(1, err.Error())
	}
}

func (this *contextImplement) GetUrlParamToStruct(requireStruct interface{}) {
	if this.GetMethod() != "GET" {
		language.Throw(1, "请求Method不是Get方法: "+this.GetMethod())
	}
	this.GetParamToStruct(requireStruct)
}

func (this *contextImplement) GetFormParamToStruct(requireStruct interface{}) {
	if this.GetMethod() != "POST" {
		language.Throw(1, "请求Method不是POST方法: "+this.GetMethod())
	}
	this.GetParamToStruct(requireStruct)
}

func (this *contextImplement) GetFile(key string) ([]byte, string, error) {
	file, fileHeader, err := this.request.FormFile(key)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", err
	}
	return data, fileHeader.Filename, nil
}

func (this *contextImplement) GetCookie(key string) string {
	ck, err := this.request.Cookie(key)
	if err != nil {
		return ""
	}
	return ck.Value
}

func (this *contextImplement) GetProxy() []string {
	if ips := this.request.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

func (this *contextImplement) GetMethod() string {
	return this.request.Method
}

func (this *contextImplement) GetScheme() string {
	if this.request.URL.Scheme != "" {
		return this.request.URL.Scheme
	}
	if this.request.TLS == nil {
		return "http"
	}
	return "https"
}

func (this *contextImplement) GetHost() string {
	if this.request.Host != "" {
		hostParts := strings.Split(this.request.Host, ":")
		if len(hostParts) > 0 {
			return hostParts[0]
		}
		return this.request.Host
	}
	return "localhost"
}

func (this *contextImplement) GetPort() int {
	if this.request.Host != "" {
		hostParts := strings.Split(this.request.Host, ":")
		if len(hostParts) > 1 {
			port, err := strconv.Atoi(hostParts[1])
			if err != nil {
				return 80
			}
			return port
		}
		return 80
	}
	return 80
}

func (this *contextImplement) GetSite() string {
	return this.GetScheme() + "://" + this.request.Host
}

func (this *contextImplement) GetRemoteAddr() string {
	ips := this.GetProxy()
	if len(ips) > 0 && ips[0] != "" {
		return ips[0]
	}
	return this.request.RemoteAddr
}

func (this *contextImplement) GetRemoteIP() string {
	addr := this.GetRemoteAddr()
	ip := strings.Split(addr, ":")
	if len(ip) > 0 {
		if ip[0] != "" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (this *contextImplement) GetRemotePort() int {
	addr := this.GetRemoteAddr()
	ip := strings.Split(addr, ":")
	if len(ip) > 1 {
		port, err := strconv.Atoi(ip[1])
		if err != nil {
			return 80
		}
		return port
	}
	return 80
}

func (this *contextImplement) GetUserAgent() string {
	return this.request.Header.Get("User-Agent")
}

func (this *contextImplement) SetUserAgent(data string) {
	this.request.Header.Set("User-Agent", data)
}

func (this *contextImplement) GetHeader(key string) string {
	return this.request.Header.Get(key)
}

func (this *contextImplement) IsUpload() bool {
	return strings.Contains(this.request.Header.Get("Content-Type"), "multipart/form-data")
}

func (this *contextImplement) IsLocal() bool {
	return this.GetRemoteIP() == "127.0.0.1"
}

func (this *contextImplement) Write(data []byte) {
	writer := this.responseWriter
	writer.Write(data)
}

func (this *contextImplement) WriteHeader(key string, value string) {
	writer := this.responseWriter
	writer.Header().Set(key, value)
}

func (this *contextImplement) WriteStatus(code int) {
	writer := this.responseWriter
	writer.WriteHeader(code)
}

func (this *contextImplement) WriteMimeHeader(mime string, title string) {
	writer := this.responseWriter
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

func (this *contextImplement) GetRawRequest() interface{} {
	return this.request
}

func (this *contextImplement) GetRawResponseWriter() interface{} {
	return this.responseWriter
}

func (this *contextImplement) GetRawTesting() interface{} {
	return this.testing
}

func (this *contextImplement) SerializeRequest() (ContextSerializeRequest, error) {
	//cache 序列化的数据
	if this.serializeRequest != nil {
		return *this.serializeRequest, nil
	}

	//生成序列化数据
	request := this.request
	result := &ContextSerializeRequest{}
	result.Method = request.Method
	result.Url = request.URL.String()
	result.Header = request.Header
	this.serializeRequest = result
	return *result, nil
}

func (this *contextImplement) DeSerializeRequest(data ContextSerializeRequest) error {
	//反序列化数据
	newRequest, err := http.NewRequest(data.Method, data.Url, nil)
	if err != nil {
		return err
	}
	newRequest.Header = data.Header
	this.request = newRequest
	this.serializeRequest = &data
	this.parseInput()
	return nil
}
