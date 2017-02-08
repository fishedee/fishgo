package sdk

import (
	"fmt"

	. "github.com/fishedee/encoding"
	. "github.com/fishedee/util"
)

type YunpianSdk struct {
	ApiKey    string
	ApiSecert string
}

type YunpianSdkSendSmsResult struct {
	Count  int     `json:"count"`
	Fee    float64 `json:"count"`
	Unit   string  `json:"count"`
	Mobile string  `json:"count"`
	Sid    int     `json:"count"`
}

type YunpianSdkError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (this *YunpianSdkError) GetCode() int {
	return this.Code
}

func (this *YunpianSdkError) GetMsg() string {
	return this.Msg
}

func (this *YunpianSdkError) Error() string {
	return fmt.Sprintf("错误码为：%v，错误描述为：%v", this.Code, this.Msg)
}

func (this *YunpianSdk) api(url string, data interface{}, responseData interface{}) error {
	var dataByte []byte
	err := DefaultAjaxPool.Post(&Ajax{
		Url:          "https://sms.yunpian.com" + url,
		Data:         data,
		DataType:     "url",
		ResponseData: &dataByte,
	})
	if err != nil {
		return err
	}
	var result YunpianSdkError
	err = DecodeJson(dataByte, &result)
	if err == nil && result.Code != 0 {
		return &result
	}
	err = DecodeJson(dataByte, &responseData)
	if err != nil {
		return err
	}
	return nil
}

func (this *YunpianSdk) SendSms(mobile string, text string) (YunpianSdkSendSmsResult, error) {
	var result YunpianSdkSendSmsResult
	err := this.api("/v2/sms/single_send.json", map[string]string{
		"apikey": this.ApiKey,
		"mobile": mobile,
		"text":   text,
	}, &result)
	if err != nil {
		return YunpianSdkSendSmsResult{}, err
	}
	return result, err
}
