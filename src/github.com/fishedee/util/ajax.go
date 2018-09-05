package util

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/fishedee/encoding"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type AjaxPoolOption struct {
	Timeout         time.Duration
	MaxIdleConnect  int
	HasCookieJar    bool
	TLSClientConfig *tls.Config
	Proxy           string
}

type AjaxPool struct {
	client *http.Client
}

type AjaxStatusCodeError struct {
	statusCode int
	body       []byte
}

func (this *AjaxStatusCodeError) Error() string {
	return "返回码不是200，而是" + strconv.Itoa(this.statusCode) + ",数据为：[" + string(this.body) + "]"
}

type Ajax struct {
	Method   string
	Url      string
	Header   interface{}
	DataType string
	Data     interface{}
	Cookie   interface{}

	ResponseDataType string
	ResponseData     interface{}
	ResponseHeader   interface{}
	ResponseCookie   interface{}
}

var DefaultAjaxPool *AjaxPool

func init() {
	DefaultAjaxPool = NewAjaxPool(nil)
}

func NewAjaxPool(option *AjaxPoolOption) *AjaxPool {
	if option == nil {
		option = &AjaxPoolOption{
			Timeout:         30 * time.Second,
			HasCookieJar:    false,
			TLSClientConfig: nil,
		}
	}
	client := &http.Client{}
	if option.Timeout <= 0 {
		option.Timeout = 30 * time.Second
	}
	client.Timeout = option.Timeout
	if option.HasCookieJar {
		jar, err := cookiejar.New(nil)
		if err != nil {
			panic(err)
		}
		client.Jar = jar
	}
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if option.MaxIdleConnect <= 0 {
		option.MaxIdleConnect = 100
	}
	transport.MaxIdleConns = option.MaxIdleConnect
	if option.TLSClientConfig != nil {
		transport.TLSClientConfig = option.TLSClientConfig
	}
	if option.Proxy != "" {
		proxyUrl, err := url.Parse(option.Proxy)
		if err != nil {
			panic(err)
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	client.Transport = transport
	return &AjaxPool{
		client: client,
	}
}

func (this *AjaxPool) Do(option *Ajax) error {
	return option.Send(this.client)
}

func (this *AjaxPool) Get(option *Ajax) error {
	option.Method = "GET"
	return option.Send(this.client)
}

func (this *AjaxPool) Post(option *Ajax) error {
	option.Method = "POST"
	return option.Send(this.client)
}

func (this *AjaxPool) Del(option *Ajax) error {
	option.Method = "DEL"
	return option.Send(this.client)
}

func (this *AjaxPool) Put(option *Ajax) error {
	option.Method = "PUT"
	return option.Send(this.client)
}

func (this *Ajax) createRequestMethod() (string, error) {
	httpMethod := map[string]string{
		"get":  "GET",
		"post": "POST",
		"del":  "DEL",
		"put":  "PUT",
	}
	httpMethodInfo, ok := httpMethod[strings.ToLower(this.Method)]
	if ok == false {
		return "", errors.New("invalid http method : " + this.Method)
	}
	return httpMethodInfo, nil
}

func (this *Ajax) createRequestUrl() (string, error) {
	requestUrlInfo, err := url.Parse(this.Url)
	if err != nil {
		return "", err
	}
	if requestUrlInfo.Scheme != "http" &&
		requestUrlInfo.Scheme != "https" {
		return "", errors.New("invalid http url scheme: " + requestUrlInfo.Scheme)
	}
	return this.Url, nil
}

func (this *Ajax) createRequestPlainData() (string, []byte, error) {
	var result []byte
	dataString, ok1 := this.Data.(string)
	dataByteSlice, ok2 := this.Data.([]byte)
	if ok1 {
		result = []byte(dataString)
	} else if ok2 {
		result = dataByteSlice
	} else {
		return "", nil, errors.New("invalid plain data type " + fmt.Sprintf("%v", this.Data))
	}
	return "text/plain; charset=UTF-8", result, nil
}

func (this *Ajax) createRequestUrlData() (string, []byte, error) {
	var result []byte
	var err error
	urlValues := url.Values{}
	dataUrlValues, ok1 := this.Data.(url.Values)
	if ok1 {
		urlValues = dataUrlValues
		result = []byte(urlValues.Encode())
	} else {
		result, err = EncodeUrlQuery(this.Data)
		if err != nil {
			return "", nil, err
		}
	}
	return "application/x-www-form-urlencoded; charset=UTF-8", result, nil
}

func (this *Ajax) createRequestJsonData() (string, []byte, error) {
	result, err := json.Marshal(this.Data)
	if err != nil {
		return "", nil, err
	}
	return "application/json; charset=UTF-8", result, nil
}

func (this *Ajax) createRequestFormData() (string, []byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	formData, ok := this.Data.(map[string]interface{})
	if !ok {
		return "", nil, errors.New("invalid form data type is not map[string]interface{}")
	}
	for key, singleFormData := range formData {
		singleFormDataString, ok1 := singleFormData.(string)
		singleFormDataFile, ok2 := singleFormData.(*os.File)
		singleFormDataSlice, ok3 := singleFormData.([]interface{})
		if ok1 {
			err := bodyWriter.WriteField(key, singleFormDataString)
			if err != nil {
				return "", nil, err
			}
		} else if ok2 {
			fileWriter, err := bodyWriter.CreateFormFile(key, singleFormDataFile.Name())
			if err != nil {
				return "", nil, err
			}
			_, err = io.Copy(fileWriter, singleFormDataFile)
			if err != nil {
				return "", nil, err
			}
		} else if ok3 {
			if len(singleFormDataSlice) != 2 {
				return "", nil, errors.New("invalid form data value,should has two argument")
			}
			singleFormDataSliceName := singleFormDataSlice[0].(string)
			singleFormDataSliceValue := singleFormDataSlice[1]
			fileWriter, err := bodyWriter.CreateFormFile(key, singleFormDataSliceName)
			if err != nil {
				return "", nil, err
			}
			singleFormDataSliceValueByte, ok1 := singleFormDataSliceValue.([]byte)
			singleFormDataSliceValueReader, ok2 := singleFormDataSliceValue.(io.Reader)
			if ok1 {
				_, err := fileWriter.Write(singleFormDataSliceValueByte)
				if err != nil {
					return "", nil, err
				}
			} else if ok2 {
				_, err := io.Copy(fileWriter, singleFormDataSliceValueReader)
				if err != nil {
					return "", nil, err
				}
			} else {
				return "", nil, errors.New("invalid form data slice,not byte[] or reader : " + singleFormDataSliceName)
			}
		} else {
			return "", nil, errors.New("invalid form data: " + key)
		}
	}
	err := bodyWriter.Close()
	if err != nil {
		return "", nil, err
	}
	return bodyWriter.FormDataContentType(), bodyBuf.Bytes(), nil
}

func (this *Ajax) createRequestData() (string, []byte, error) {
	if this.Data == nil {
		return "", nil, nil
	} else {
		if this.DataType == "" {
			_, ok1 := this.Data.(string)
			_, ok2 := this.Data.([]byte)
			if ok1 || ok2 {
				return this.createRequestPlainData()
			} else {
				return this.createRequestUrlData()
			}
		} else if this.DataType == "plain" {
			return this.createRequestPlainData()
		} else if this.DataType == "url" {
			return this.createRequestUrlData()
		} else if this.DataType == "form" {
			return this.createRequestFormData()
		} else if this.DataType == "json" {
			return this.createRequestJsonData()
		} else {
			return "", nil, errors.New("invalid dataType : " + this.DataType)
		}
	}
}

func (this *Ajax) createRequestHeader(request *http.Request) error {
	if this.Header == nil {
		return nil
	}
	map1, ok1 := this.Header.(map[string]string)
	map2, ok2 := this.Header.(http.Header)
	if ok1 {
		for key, value := range map1 {
			request.Header.Set(key, value)
		}
	} else if ok2 {
		for key, value := range map2 {
			request.Header[key] = value
		}
	} else {
		return errors.New("invalid header type " + fmt.Sprintf("%v", this.Header))
	}
	return nil
}

func (this *Ajax) createRequestCookie(request *http.Request) error {
	if this.Cookie == nil {
		return nil
	}
	cookie1, ok1 := this.Cookie.([]*http.Cookie)
	if ok1 {
		for _, value := range cookie1 {
			request.AddCookie(value)
		}
	} else {
		return errors.New("invalid cookie type " + fmt.Sprintf("%v", this.Cookie))
	}
	return nil
}

func (this *Ajax) createRequest() (*http.Request, error) {
	method, err := this.createRequestMethod()
	if err != nil {
		return nil, err
	}

	url, err := this.createRequestUrl()
	if err != nil {
		return nil, err
	}

	contentType, data, err := this.createRequestData()
	if err != nil {
		return nil, err
	}

	var dataReader io.Reader
	if method == "GET" || method == "DEL" {
		if len(data) != 0 {
			url = url + "?" + string(data)
		}
		dataReader = nil
	} else {
		dataReader = bytes.NewReader(data)
	}
	request, err := http.NewRequest(method, url, dataReader)
	if err != nil {
		return nil, err
	}
	if dataReader != nil {
		request.Header.Set("Content-Type", contentType)
		request.Header.Set("Content-Length", strconv.Itoa(len(data)))
	}

	err = this.createRequestHeader(request)
	if err != nil {
		return nil, err
	}

	err = this.createRequestCookie(request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (this *Ajax) saveResponsePlainData(respString []byte) error {
	dataString, ok1 := this.ResponseData.(*string)
	dataByteSlice, ok2 := this.ResponseData.(*[]byte)
	if ok1 {
		*dataString = string(respString)
	} else if ok2 {
		*dataByteSlice = respString
	} else {
		return errors.New("invalid response data plain type " + fmt.Sprintf("%v", this.ResponseData))
	}
	return nil
}

func (this *Ajax) saveResponseUrlData(respString []byte) error {
	values, err := url.ParseQuery(string(respString))
	if err != nil {
		return err
	}

	dataValues, ok1 := this.ResponseData.(*url.Values)
	dataMap, ok2 := this.ResponseData.(*map[string]string)
	if ok1 {
		*dataValues = values
	} else if ok2 {
		for key, _ := range values {
			(*dataMap)[key] = values.Get(key)
		}
	} else {
		return errors.New("invalid response data url type " + fmt.Sprintf("%v", this.ResponseData))
	}
	return nil
}

func (this *Ajax) saveResponseJsonData(respString []byte) error {
	err := json.Unmarshal(respString, this.ResponseData)
	if err != nil {
		return err
	}
	return nil
}

func (this *Ajax) decodeResponseBody(response *http.Response) ([]byte, error) {
	contentEncoding := response.Header.Get("Content-Encoding")
	var reader io.ReadCloser
	var decodeReader io.ReadCloser
	if strings.Contains(contentEncoding, "gzip") {
		gzipReader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		decodeReader = gzipReader
		reader = gzipReader
	} else if strings.Contains(contentEncoding, "deflate") {
		decodeReader := flate.NewReader(response.Body)
		reader = decodeReader
	} else {
		reader = response.Body
	}
	respString, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if decodeReader != nil {
		err := decodeReader.Close()
		if err != nil {
			return nil, err
		}
	}
	return respString, nil
}

func (this *Ajax) saveResponseData(response *http.Response) ([]byte, error) {
	if this.ResponseData == nil {
		return []byte{}, nil
	}
	respString, err := this.decodeResponseBody(response)
	if err != nil {
		return []byte{}, err
	}

	dataType := this.ResponseDataType
	if dataType == "" {
		_, ok1 := this.ResponseData.(*string)
		_, ok2 := this.ResponseData.(*[]byte)
		if ok1 || ok2 {
			return respString, this.saveResponsePlainData(respString)
		} else {
			return respString, this.saveResponseJsonData(respString)
		}
	} else if dataType == "plain" {
		return respString, this.saveResponsePlainData(respString)
	} else if dataType == "url" {
		return respString, this.saveResponseUrlData(respString)
	} else if dataType == "json" {
		return respString, this.saveResponseJsonData(respString)
	} else {
		return respString, errors.New("invalid response data type " + dataType)
	}
}

func (this *Ajax) saveResponseHeader(response *http.Response) error {
	if this.ResponseHeader == nil {
		return nil
	}
	dataHeader, ok1 := this.ResponseHeader.(*http.Header)
	dataMap, ok2 := this.ResponseHeader.(*map[string]string)
	if ok1 {
		*dataHeader = response.Header
	} else if ok2 {
		for key, _ := range response.Header {
			(*dataMap)[key] = response.Header.Get(key)
		}
	} else {
		return errors.New("invalid response header type " + fmt.Sprintf("%v", this.ResponseHeader))
	}
	return nil
}

func (this *Ajax) saveResponseCookie(response *http.Response) error {
	if this.ResponseCookie == nil {
		return nil
	}
	dataCookie, ok1 := this.ResponseCookie.(*[]*http.Cookie)
	if ok1 {
		*dataCookie = response.Cookies()
	} else {
		return errors.New("invalid response cookie type " + fmt.Sprintf("%v", this.ResponseCookie))
	}
	return nil
}

func (this *Ajax) saveResponse(response *http.Response) error {
	var err error
	var data []byte

	data, err = this.saveResponseData(response)
	if err != nil {
		return err
	}

	err = this.saveResponseHeader(response)
	if err != nil {
		return err
	}

	err = this.saveResponseCookie(response)
	if err != nil {
		return err
	}

	if response.StatusCode < 200 ||
		response.StatusCode >= 300 {
		return &AjaxStatusCodeError{
			statusCode: response.StatusCode,
			body:       data,
		}
	}
	return nil
}

func (this *Ajax) Send(client *http.Client) error {
	//创建请求
	request, err := this.createRequest()
	if err != nil {
		return err
	}

	//执行请求
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//保存请求
	err = this.saveResponse(response)
	if err != nil {
		return err
	}
	return nil
}
