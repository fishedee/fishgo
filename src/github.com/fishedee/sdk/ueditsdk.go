package sdk

import (
	"fmt"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	"io/ioutil"
	"net/http"
	"path"
)

type UeditSdk struct {
}

type UeditSdkUploadCallback func([]byte) (string, error)

type UeditSdkConfig struct {
	//图片配置项
	ImageMaxSize        int                    `json:"imageMaxSize"`
	ImageAllowFiles     []string               `json:"imageAllowFiles"`
	ImageCompressBorder int                    `json:"imageCompressBorder"`
	ImageInsertAlign    string                 `json:"imageInsertAlign"`
	ImageUploadCallback UeditSdkUploadCallback `json:"-"`

	//涂鸦配置项
	ScrawlMaxSize        int                    `json:"ScrawlMaxSize"`
	ScrawlInsertAlign    string                 `json:"scrawlInsertAlign"`
	ScrawlUploadCallback UeditSdkUploadCallback `json:"-"`
}

type UeditSdkError struct {
	Data string
}

func (this *UeditSdkError) Error() string {
	return this.Data
}

type ueditSdkConfigReal struct {
	UeditSdkConfig
	ImageActionName string `json:"imageActionName"`
	ImageFieldName  string `json:"imageFieldName"`
	ImageUrlPrefix  string `json:"imageUrlPrefix"`

	ScrawlActionName string `json:"scrawlActionName"`
	ScrawlFieldName  string `json:"scrawlFieldName"`
	ScrawlUrlPrefix  string `json:"scrawlUrlPrefix"`
}

type ueditSdkUploadConfig struct {
	MaxSize    int
	AllowFiles []string
	OriginName string
	Format     string
	CallBack   UeditSdkUploadCallback
}

func (this *UeditSdk) getConfig(config UeditSdkConfig) ueditSdkConfigReal {
	result := ueditSdkConfigReal{
		UeditSdkConfig: config,
	}
	//上传图片配置项
	result.ImageActionName = "uploadimage"
	result.ImageFieldName = "upfile"
	result.ImageUrlPrefix = ""
	if result.ImageMaxSize == 0 {
		result.ImageMaxSize = 2048000
	}
	if result.ImageAllowFiles == nil {
		result.ImageAllowFiles = []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}
	}
	if result.ImageCompressBorder == 0 {
		result.ImageCompressBorder = 640
	}
	if result.ImageInsertAlign == "" {
		result.ImageInsertAlign = "none"
	}

	//涂鸦配置项
	result.ScrawlActionName = "uploadscrawl"
	result.ScrawlFieldName = "upfile"
	result.ScrawlUrlPrefix = ""
	if result.ScrawlMaxSize == 0 {
		result.ScrawlMaxSize = 2048000
	}
	if result.ScrawlInsertAlign == "" {
		result.ScrawlInsertAlign = "none"
	}
	return result
}

func (this *UeditSdk) Handle(config UeditSdkConfig, request *http.Request) ([]byte, error) {
	realConfig := this.getConfig(config)

	var input struct {
		Action   string `url:"action"`
		Callback string `url:"callback"`
	}

	queryStr := request.URL.RawQuery
	err := DecodeUrlQuery([]byte(queryStr), &input)
	if err != nil {
		return nil, err
	}

	result, err := this.handleAction(realConfig, input.Action, request)
	if err != nil {
		return nil, err
	}

	resultByte, err := EncodeJson(result)
	if err != nil {
		return nil, err
	}
	if input.Callback != "" {
		return []byte(input.Callback + "(" + string(resultByte) + ")"), nil
	} else {
		return resultByte, nil
	}
}

func (this *UeditSdk) handleAction(config ueditSdkConfigReal, action string, request *http.Request) (interface{}, error) {
	var result interface{}
	var err error
	if action == "config" {
		result = config
	} else if action == "uploadimage" {
		uploadConfig := ueditSdkUploadConfig{
			MaxSize:    config.ImageMaxSize,
			AllowFiles: config.ImageAllowFiles,
			OriginName: "",
			Format:     "binary",
			CallBack:   config.ImageUploadCallback,
		}
		result, err = this.handleUploadAction(uploadConfig, request)
	} else if action == "uploadscrawl" {
		uploadConfig := ueditSdkUploadConfig{
			MaxSize:    config.ScrawlMaxSize,
			AllowFiles: nil,
			OriginName: "scrawl.png",
			Format:     "base64",
			CallBack:   config.ScrawlUploadCallback,
		}
		result, err = this.handleUploadAction(uploadConfig, request)
	} else {
		err = &UeditSdkError{"请求地址出错"}
	}
	if err != nil {
		sdkError, isOk := err.(*UeditSdkError)
		if isOk {
			return map[string]interface{}{
				"state": sdkError.Error(),
			}, nil
		} else {
			return nil, err
		}
	} else {
		return result, nil
	}
}

func (this *UeditSdk) handleUploadAction(config ueditSdkUploadConfig, request *http.Request) (interface{}, error) {
	var fileData []byte
	var fileName string

	if config.Format == "binary" {
		file, header, err := request.FormFile("upfile")
		if err != nil {
			return nil, err
		}
		fileName = header.Filename
		fileData, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		file := request.FormValue("upfile")
		fileData, err = DecodeBase64(file)
		if err != nil {
			return nil, err
		}
		fileName = config.OriginName
	}

	fileType := path.Ext(fileName)
	fileSize := len(fileData)

	if config.MaxSize != 0 && len(fileData) > config.MaxSize {
		return nil, &UeditSdkError{fmt.Sprintf("超出文件大小 %v", config.MaxSize)}
	}
	if len(config.AllowFiles) != 0 && ArrayIn(config.AllowFiles, fileType) == -1 {
		return nil, &UeditSdkError{fmt.Sprintf("文件类型不允许 [%v],%v", fileType, config.AllowFiles)}
	}
	url, err := config.CallBack(fileData)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"state":    "SUCCESS",
		"url":      url,
		"title":    fileName,
		"original": fileName,
		"type":     fileType,
		"size":     fileSize,
	}, nil
}
