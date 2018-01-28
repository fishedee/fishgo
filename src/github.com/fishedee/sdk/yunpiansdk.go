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

type YunpianSdkSendVoiceResult struct {
	Count int     `json:"count"`
	Fee   float64 `json:"fee"`
	Sid   string  `json:"sid"`
}

type YunpianSdkSendSmsResult struct {
	Count  int     `json:"count"`
	Fee    float64 `json:"fee"`
	Unit   string  `json:"unit"`
	Mobile string  `json:"mobile"`
	Sid    int     `json:"sid"`
}

type YunpianSdkSingleVoiceStatusResult struct {
	Sid             string `json:"sid"`
	Uid             string `json:"uid"`
	UserReceiveTime string `json:"user_receive_time"`
	ErrorMsg        string `json:"error_msg"`
	Mobile          string `json:"mobile"`
	ReportStatus    string `json:"report_status"`
}

type YunpianSdkVoiceStatusResult []YunpianSdkSingleVoiceStatusResult

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
		Url:          url,
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

func (this *YunpianSdk) DecodeVoiceStatus(voiceStatus string) (YunpianSdkVoiceStatusResult, error) {
	var result YunpianSdkVoiceStatusResult
	err := DecodeJson([]byte(voiceStatus), &result)
	if err != nil {
		return YunpianSdkVoiceStatusResult{}, err
	}
	return result, err
}

func (this *YunpianSdk) SendVoice(mobile string, code string, callBackUrl string) (YunpianSdkSendVoiceResult, error) {
	var result YunpianSdkSendVoiceResult
	err := this.api("https://voice.yunpian.com/v2/voice/send.json", map[string]string{
		"apikey":       this.ApiKey,
		"mobile":       mobile,
		"code":         code,
		"callback_url": callBackUrl,
	}, &result)
	if err != nil {
		return YunpianSdkSendVoiceResult{}, err
	}
	return result, err
}

func (this *YunpianSdk) SendSms(mobile string, text string) (YunpianSdkSendSmsResult, error) {
	var result YunpianSdkSendSmsResult
	err := this.api("https://sms.yunpian.com/v2/sms/single_send.json", map[string]string{
		"apikey": this.ApiKey,
		"mobile": mobile,
		"text":   text,
	}, &result)
	if err != nil {
		return YunpianSdkSendSmsResult{}, err
	}
	return result, err
}

// 批量发送短信
func (this *YunpianSdk) SendMsgs(mobile string, text string) (YunpianSdkSendSmsResult, error) {
	var result YunpianSdkSendSmsResult
	err := this.api("https://sms.yunpian.com/v2/sms/batch_send.json", map[string]string{
		"apikey": this.ApiKey,
		"mobile": mobile,
		"text":   text,
	}, &result)
	if err != nil {
		return YunpianSdkSendSmsResult{}, err
	}
	return result, err
}
