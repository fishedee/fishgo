package sdk

import (
	"fmt"
	"reflect"
	"testing"
)

func assertYunpianSdkEqual(t *testing.T, left interface{}, right interface{}) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("assert fail: %+v != %+v", left, right)
	}
}

func TestYunpianSdkSendSms(t *testing.T) {
	sdk := &YunpianSdk{
		ApiKey: "xxxx",
	}
	result, err := sdk.SendSms(
		"15018749403",
		"【云片网】您的验证码是2921",
	)
	assertYunpianSdkEqual(t, err == nil, true)
	assertYunpianSdkEqual(t, result.Count, 1)
	fmt.Println(err)
}

func TestYunpianSdkSendVoice(t *testing.T) {
	sdk := &YunpianSdk{
		ApiKey: "xxxx",
	}
	result, err := sdk.SendVoice(
		"15018749403",
		"1234",
		"",
	)
	assertYunpianSdkEqual(t, err == nil, true)
	assertYunpianSdkEqual(t, result.Count, 1)
	fmt.Println(err)
}

func TestYunpianSdkDecodeSmsStatus(t *testing.T) {
	sdk := &YunpianSdk{}
	result, err := sdk.DecodeVoiceStatus(`[{"sid":9527,"uid":null,"user_receive_time":"2014-03-17 22:55:21","error_msg":"DELIVRD","mobile":"15205201314","report_status":"SUCCESS"},{"sid":9528,"uid":null,"user_receive_time":"2014-03-17 22:55:23","error_msg":"DELIVRD","mobile":"15212341234","report_status":"SUCCESS"}]`)
	assertYunpianSdkEqual(t, err == nil, true)
	assertYunpianSdkEqual(t, result, YunpianSdkVoiceStatusResult([]YunpianSdkSingleVoiceStatusResult{
		YunpianSdkSingleVoiceStatusResult{
			Sid:             "9527",
			Uid:             "",
			UserReceiveTime: "2014-03-17 22:55:21",
			ErrorMsg:        "DELIVRD",
			Mobile:          "15205201314",
			ReportStatus:    "SUCCESS",
		},
		YunpianSdkSingleVoiceStatusResult{
			Sid:             "9528",
			Uid:             "",
			UserReceiveTime: "2014-03-17 22:55:23",
			ErrorMsg:        "DELIVRD",
			Mobile:          "15212341234",
			ReportStatus:    "SUCCESS",
		},
	}))
}
