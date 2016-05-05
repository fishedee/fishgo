package web

import (
	"bytes"
	"fmt"
	"github.com/fishedee/encoding"
	"github.com/fishedee/language"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Testing        *testing.T
	inputData      map[string]interface{}
}

func NewContext(request *http.Request, response http.ResponseWriter, t *testing.T) Context {
	result := Context{
		Request:        request,
		ResponseWriter: response,
		Testing:        t,
	}
	result.parseInput()
	return result
}

func (this *Context) parseInput() {
	//取出get数据
	request := this.Request
	queryInput := request.URL.RawQuery

	//取出post数据
	postInput := ""
	ct := request.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	ct, _, err := mime.ParseMediaType(ct)
	if ct == "application/x-www-form-urlencoded" {
		byteArray, _ := ioutil.ReadAll(this.Request.Body)
		this.Request.Body = ioutil.NopCloser(bytes.NewReader(byteArray))
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

func (this *Context) GetParam(key string) string {
	result, isExist := this.inputData[key]
	if !isExist {
		return ""
	}
	return fmt.Sprintf("%v", result)
}

func (this *Context) GetParamToStruct(requireStruct interface{}) {
	//导出到struct
	err := language.MapToArray(this.inputData, requireStruct, "url")
	if err != nil {
		language.Throw(1, err.Error())
	}
}

func (this *Context) GetUrlParamToStruct(requireStruct interface{}) {
	if this.Request.Method != "GET" {
		language.Throw(1, "请求Method不是Get方法")
	}
	this.GetParamToStruct(requireStruct)
}

func (this *Context) GetFormParamToStruct(requireStruct interface{}) {
	if this.Request.Method != "POST" {
		language.Throw(1, "请求Method不是POST方法")
	}
	this.GetParamToStruct(requireStruct)
}

func (this *Context) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return this.Request.FormFile(key)
}

func (this *Context) GetCookie(key string) string {
	ck, err := this.Request.Cookie(key)
	if err != nil {
		return ""
	}
	return ck.Value
}

func (this *Context) GetProxy() []string {
	if ips := this.Request.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

func (this *Context) GetMethod() string {
	return this.Request.Method
}

func (this *Context) GetScheme() string {
	if this.Request.URL.Scheme != "" {
		return this.Request.URL.Scheme
	}
	if this.Request.TLS == nil {
		return "http"
	}
	return "https"
}

func (this *Context) GetHost() string {
	if this.Request.Host != "" {
		hostParts := strings.Split(this.Request.Host, ":")
		if len(hostParts) > 0 {
			return hostParts[0]
		}
		return this.Request.Host
	}
	return "localhost"
}

func (this *Context) GetPort() int {
	if this.Request.Host != "" {
		hostParts := strings.Split(this.Request.Host, ":")
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

func (this *Context) GetSite() string {
	return this.GetScheme() + "://" + this.Request.Host
}

func (this *Context) GetRemoteAddr() string {
	ips := this.GetProxy()
	if len(ips) > 0 && ips[0] != "" {
		return ips[0]
	}
	return this.Request.RemoteAddr
}

func (this *Context) GetRemoteIP() string {
	addr := this.GetRemoteAddr()
	ip := strings.Split(addr, ":")
	if len(ip) > 0 {
		if ip[0] != "" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (this *Context) GetRemotePort() int {
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

func (this *Context) GetUserAgent() string {
	return this.Request.Header.Get("User-Agent")
}

func (this *Context) GetHeader(key string) string {
	return this.Request.Header.Get(key)
}

func (this *Context) IsUpload() bool {
	return strings.Contains(this.Request.Header.Get("Content-Type"), "multipart/form-data")
}

func (this *Context) Write(data []byte) {
	writer := this.ResponseWriter
	writer.Write(data)
}

func (this *Context) WriteHeader(key string, value string) {
	writer := this.ResponseWriter
	writer.Header().Set(key, value)
}

func (this *Context) WriteStatus(code int) {
	writer := this.ResponseWriter
	writer.WriteHeader(code)
}

func (this *Context) WriteMimeHeader(mime string, title string) {
	writer := this.ResponseWriter
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
