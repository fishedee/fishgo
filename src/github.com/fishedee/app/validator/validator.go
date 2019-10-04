package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/fishedee/app/quicktag"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Validator interface {
	//输入参数
	Url() *url.URL
	Body() io.ReadCloser

	Param(name string) (string, error)
	MustParam(name string) string

	Query(name string) (string, error)
	MustQuery(name string) string
	QueryArray(name string) ([]string, error)
	MustQueryArray(name string) []string

	Form(name string) (string, error)
	MustForm(name string) string
	FormArray(name string) ([]string, error)
	MustFormArray(name string) []string

	File(key string) (*multipart.FileHeader, error)
	MustFile(key string) *multipart.FileHeader

	BindParam(obj interface{}) error
	MustBindParam(obj interface{})
	BindQuery(obj interface{}) error
	MustBindQuery(obj interface{})
	BindForm(obj interface{}) error
	MustBindForm(obj interface{})
	BindJson(obj interface{}) error
	MustBindJson(obj interface{})
	Bind(obj interface{}) error
	MustBind(obj interface{})

	//元信息
	Request() *http.Request
	Cookie(key string) string
	Proxy() []string
	Method() string
	Scheme() string
	Host() string
	Port() int
	Site() string
	RemoteAddr() string
	RemoteIP() string
	RemotePort() int
	UserAgent() string
	Header(key string) string
	IsUpload() bool
	IsLocal() bool
}

type ValidatorFactory interface {
	Create(r *http.Request, param map[string]string) Validator
}

type ValidatorConfig struct {
	MaxFormSize       int `config:"maxformsize"`
	MaxFileSize       int `config:"maxfilesize"`
	MaxFileMemorySize int `config:"maxfilememorysize"`
}

type validatorFactoryImplement struct {
	config ValidatorConfig
}

func NewValidatorFactory(config ValidatorConfig) (ValidatorFactory, error) {
	return &validatorFactoryImplement{
		config: config,
	}, nil
}

func (this *validatorFactoryImplement) Create(r *http.Request, param map[string]string) Validator {
	return newValidator(r, param, this.config)
}

type validatorImplement struct {
	config       ValidatorConfig
	request      *http.Request
	param        map[string]string
	query        map[string]interface{}
	form         map[string]interface{}
	file         map[string]*multipart.FileHeader
	jsonForm     []byte
	jsonQuickTag *QuickTag
	parseResult  error
	hasParse     bool
}

func newValidator(r *http.Request, param map[string]string, config ValidatorConfig) Validator {
	if param == nil {
		param = map[string]string{}
	}
	if config.MaxFormSize <= 0 {
		config.MaxFormSize = 1024 * 1024 * 10
	}
	if config.MaxFileSize <= 0 {
		config.MaxFileSize = 1024 * 1024 * 20
	}
	if config.MaxFileMemorySize <= 0 {
		config.MaxFileMemorySize = 1024 * 1024 * 10
	}
	return &validatorImplement{
		request:      r,
		param:        param,
		hasParse:     false,
		config:       config,
		jsonQuickTag: NewQuickTag("json"),
	}
}

func (this *validatorImplement) parse() error {
	if this.hasParse {
		return this.parseResult
	}
	this.hasParse = true
	this.parseResult = this.parseInner()
	return this.parseResult
}

type LimitedReader struct {
	R io.Reader
	N int64
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, errors.New("文件太大，拒绝读取")
	}
	if int64(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(p)
	l.N -= int64(n)
	return
}

func (this *validatorImplement) parseInner() error {
	//解析query参数
	queryInput := this.request.URL.RawQuery
	this.query = map[string]interface{}{}
	err := DecodeUrlQuery([]byte(queryInput), &this.query)
	if err != nil {
		return err
	}

	//解析body参数
	this.jsonForm = []byte{}
	this.form = map[string]interface{}{}
	this.file = map[string]*multipart.FileHeader{}
	if this.request.Body == nil {
		return nil
	}
	ct := this.request.Header.Get("Content-Type")
	ct, ctParam, err := mime.ParseMediaType(ct)

	if ct == "application/x-www-form-urlencoded" || ct == "" {
		bodyReader := &LimitedReader{this.request.Body, int64(this.config.MaxFormSize)}
		byteArray, err := ioutil.ReadAll(bodyReader)
		if err != nil {
			return err
		}
		err = DecodeUrlQuery(byteArray, &this.form)
		if err != nil {
			return err
		}
	} else if ct == "application/json" {
		bodyReader := &LimitedReader{this.request.Body, int64(this.config.MaxFormSize)}
		byteArray, err := ioutil.ReadAll(bodyReader)
		if err != nil {
			return err
		}
		this.jsonForm = byteArray
	} else if ct == "multipart/form-data" {
		boundary, ok := ctParam["boundary"]
		if !ok {
			return errors.New("multipart has not boundary!")
		}
		bodyReader := &LimitedReader{this.request.Body, int64(this.config.MaxFileSize)}
		multipartReader := multipart.NewReader(bodyReader, boundary)
		form, err := multipartReader.ReadForm(int64(this.config.MaxFileMemorySize))
		if err != nil {
			return err
		}
		for key, value := range form.Value {
			if len(value) == 0 {
				continue
			}
			this.form[key] = value[0]
		}
		for key, value := range form.File {
			if len(value) == 0 {
				continue
			}
			this.file[key] = value[0]
		}
	} else {
		return fmt.Errorf("unspport content-type:[%v]", ct)
	}
	return nil
}

func (this *validatorImplement) getData(data map[string]interface{}, name string) ([]string, error) {
	result, isExist := data[name]
	if isExist == false {
		return nil, nil
	}
	resultArray := []string{}
	if resultArrayInterface, isOk := result.([]interface{}); isOk {
		for i := 0; i != len(resultArrayInterface); i++ {
			resultArray = append(resultArray, fmt.Sprintf("%v", resultArrayInterface[i]))
		}
	} else {
		resultArray = append(resultArray, fmt.Sprintf("%v", result))
	}
	return resultArray, nil
}

func (this *validatorImplement) Url() *url.URL {
	return this.request.URL
}

func (this *validatorImplement) Body() io.ReadCloser {
	return this.request.Body
}

func (this *validatorImplement) Param(name string) (string, error) {
	return this.param[name], nil
}

func (this *validatorImplement) MustParam(name string) string {
	result, err := this.Param(name)
	if err != nil {
		Throw(1, err.Error())
	}
	return result
}

func (this *validatorImplement) Query(name string) (string, error) {
	result, err := this.QueryArray(name)
	if err != nil {
		return "", err
	}
	if len(result) == 0 {
		return "", nil
	}
	return result[0], nil
}

func (this *validatorImplement) MustQuery(name string) string {
	result, err := this.Query(name)
	if err != nil {
		Throw(1, err.Error())
	}
	return result
}

func (this *validatorImplement) QueryArray(name string) ([]string, error) {
	err := this.parse()
	if err != nil {
		return nil, err
	}
	return this.getData(this.query, name)
}

func (this *validatorImplement) MustQueryArray(name string) []string {
	result, err := this.QueryArray(name)
	if err != nil {
		Throw(1, err.Error())
	}
	return result
}

func (this *validatorImplement) Form(name string) (string, error) {
	result, err := this.FormArray(name)
	if err != nil {
		return "", err
	}
	if len(result) == 0 {
		return "", nil
	}
	return result[0], nil
}

func (this *validatorImplement) MustForm(name string) string {
	result, err := this.Form(name)
	if err != nil {
		Throw(1, err.Error())
	}
	return result
}

func (this *validatorImplement) FormArray(name string) ([]string, error) {
	err := this.parse()
	if err != nil {
		return nil, err
	}
	return this.getData(this.form, name)
}

func (this *validatorImplement) MustFormArray(name string) []string {
	result, err := this.FormArray(name)
	if err != nil {
		Throw(1, err.Error())
	}
	return result
}

func (this *validatorImplement) File(key string) (*multipart.FileHeader, error) {
	err := this.parse()
	if err != nil {
		return nil, err
	}
	return this.file[key], nil
}

func (this *validatorImplement) MustFile(key string) *multipart.FileHeader {
	result, err := this.File(key)
	if err != nil {
		Throw(1, err.Error())
	}
	return result
}

func (this *validatorImplement) BindParam(obj interface{}) error {
	if len(this.param) == 0 {
		return nil
	}
	return MapToArray(this.param, obj, "validator")
}

func (this *validatorImplement) MustBindParam(obj interface{}) {
	err := this.BindParam(obj)
	if err != nil {
		Throw(1, err.Error())
	}
}

func (this *validatorImplement) BindQuery(obj interface{}) error {
	err := this.parse()
	if err != nil {
		return err
	}
	if len(this.query) == 0 {
		return nil
	}
	return MapToArray(this.query, obj, "validator")
}

func (this *validatorImplement) MustBindQuery(obj interface{}) {
	err := this.BindQuery(obj)
	if err != nil {
		Throw(1, err.Error())
	}
}

func (this *validatorImplement) BindForm(obj interface{}) error {
	err := this.parse()
	if err != nil {
		return err
	}
	if len(this.form) == 0 {
		return nil
	}
	return MapToArray(this.form, obj, "validator")
}

func (this *validatorImplement) MustBindForm(obj interface{}) {
	err := this.BindForm(obj)
	if err != nil {
		Throw(1, err.Error())
	}
}

func (this *validatorImplement) BindJson(obj interface{}) error {
	err := this.parse()
	if err != nil {
		return err
	}
	if len(this.jsonForm) == 0 {
		return nil
	}
	quickTagObj := this.jsonQuickTag.GetTagInstance(obj)

	return json.Unmarshal(this.jsonForm, quickTagObj)
}

func (this *validatorImplement) MustBindJson(obj interface{}) {
	err := this.BindJson(obj)
	if err != nil {
		Throw(1, err.Error())
	}
}

func (this *validatorImplement) Bind(obj interface{}) error {
	var err error
	err = this.parse()
	if err != nil {
		return err
	}
	err = this.BindParam(obj)
	if err != nil {
		return err
	}
	err = this.BindQuery(obj)
	if err != nil {
		return err
	}
	err = this.BindForm(obj)
	if err != nil {
		return err
	}
	err = this.BindJson(obj)
	if err != nil {
		return err
	}
	return nil
}

func (this *validatorImplement) MustBind(obj interface{}) {
	err := this.Bind(obj)
	if err != nil {
		Throw(1, err.Error())
	}
}

//元信息
func (this *validatorImplement) Request() *http.Request {
	return this.request
}

func (this *validatorImplement) Cookie(key string) string {
	ck, err := this.request.Cookie(key)
	if err != nil {
		return ""
	}
	return ck.Value
}

func (this *validatorImplement) Proxy() []string {
	if ips := this.request.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

func (this *validatorImplement) Method() string {
	return this.request.Method
}

func (this *validatorImplement) Scheme() string {
	if this.request.URL.Scheme != "" {
		return this.request.URL.Scheme
	}
	if this.request.TLS == nil {
		return "http"
	}
	return "https"
}

func (this *validatorImplement) Host() string {
	if this.request.Host != "" {
		hostParts := strings.Split(this.request.Host, ":")
		if len(hostParts) > 0 {
			return hostParts[0]
		}
		return this.request.Host
	}
	return "localhost"
}

func (this *validatorImplement) Port() int {
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

func (this *validatorImplement) Site() string {
	return this.Scheme() + "://" + this.request.Host
}

func (this *validatorImplement) RemoteAddr() string {
	ips := this.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		return ips[0]
	}
	return this.request.RemoteAddr
}

func (this *validatorImplement) RemoteIP() string {
	addr := this.RemoteAddr()
	ip := strings.Split(addr, ":")
	if len(ip) > 0 {
		if ip[0] != "" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (this *validatorImplement) RemotePort() int {
	addr := this.RemoteAddr()
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

func (this *validatorImplement) UserAgent() string {
	return this.request.Header.Get("User-Agent")
}

func (this *validatorImplement) Header(key string) string {
	return this.request.Header.Get(key)
}

func (this *validatorImplement) IsUpload() bool {
	return strings.Contains(this.request.Header.Get("Content-Type"), "multipart/form-data")
}

func (this *validatorImplement) IsLocal() bool {
	return this.RemoteIP() == "127.0.0.1"
}
