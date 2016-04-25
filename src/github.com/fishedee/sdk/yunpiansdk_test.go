package sdk

import (
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
		ApiKey: "850bfe241f3d8c3ca9ecde53d161b209",
	}
	result, err := sdk.SendSms(
		"15018749403",
		"【烘焙帮】您的验证码是zzmce，打死都不要告诉别人噢～～。感谢您使用我们的app!",
	)
	assertYunpianSdkEqual(t, err == nil, true)
	assertYunpianSdkEqual(t, result.Count, 1)
}
