package sdk

import (
	"reflect"
	"testing"
)

func assertWxSdkEqual(t *testing.T, left interface{}, right interface{}) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("assert fail: %+v != %+v", left, right)
	}
}

func TestWxSdkCheckSignature() {

}

func TestWxSdkGetSignature(t *testing.T) {
	sig := getSignature(wxSdkJsSignature{
		JsApiTicket: "ojZ8YtyVyr30HheH3CM73y7h4jJE",
		NonceStr:    "asdf",
		Timestamp:   "zxvc",
		Url:         "http://www.hongbeibang.com/123#dd?4",
	})
	target := "323d89ead3fe38064577e4d66865efaf704f5d44"
	assertWxSdkEqual(t, sig, target)
}
