package sdk

import (
	"net/url"
	"reflect"
	"testing"
)

func assertDuibaSdkEqual(t *testing.T, left interface{}, right interface{}) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("assert fail: %+v != %+v", left, right)
	}
}

func TestDuibaSdkGetSign(t *testing.T) {
	testCase := []struct {
		origin map[string]interface{}
		target string
	}{
		{map[string]interface{}{
			"appKey": "testappkey",
		}, "925e68506cf5a9ac740aedc2bd78f5e4"},
		{map[string]interface{}{
			"appKey":    "testappkey",
			"timestamp": 1415250311646,
		}, "868a4339ea0400aec9b7a4742e06120e"},
		{map[string]interface{}{
			"appKey":      "testappkey",
			"description": "兑吧",
		}, "54432c8d6f76246d5890f05b7df0182f"},
		{map[string]interface{}{
			"appKey":      "testappkey",
			"description": "兑吧",
			"sign":        "123",
		}, "54432c8d6f76246d5890f05b7df0182f"},
	}

	sdk := &DuibaSdk{
		AppKey:    "testappkey",
		AppSecret: "testappsecret",
	}

	for _, singleTestCase := range testCase {
		target := sdk.getSign(singleTestCase.origin)
		assertDuibaSdkEqual(t, target, singleTestCase.target)
	}
}

func TestDuibaLoginUrl(t *testing.T) {
	sdk := &DuibaSdk{
		AppKey:    "appKey",
		AppSecret: "appSecret",
	}

	url, err := sdk.GetLoginUrl("userId001", 100)
	assertDuibaSdkEqual(t, err, nil)
	t.Errorf("%v", url)
}

func TestDuibaReceiveCreditConsume(t *testing.T) {
	urlInfo, err := url.Parse("http://api.test.hongbeibang.com/client/exchange?uid=10014&credits=1&orderNum=2016042710002341000238533&params=&type=coupon&ip=61.145.97.52&sign=87445c85ea9e88ee85b5ebbb53a1c284&waitAudit=false&timestamp=1461722423785&actualPrice=0&description=%E6%B5%8B%E8%AF%95%E4%B8%93%E7%94%A8%E4%BC%98%E6%83%A0%E5%88%B8&facePrice=1&appKey=24usbX8LZZscYRHGBLhhe5zszrg3&")
	assertDuibaSdkEqual(t, err, nil)

	sdk := &DuibaSdk{
		AppKey:    "24usbX8LZZscYRHGBLhhe5zszrg3",
		AppSecret: "4UFKw81Rev615GLtjPK4wf3hemE9",
	}
	result, err := sdk.ReceiveCreditConsume(urlInfo)
	assertDuibaSdkEqual(t, err, nil)
	assertDuibaSdkEqual(t, result, DuibaSdkReceiveCreditConsume{
		Uid:         "10014",
		Credits:     1,
		AppKey:      "24usbX8LZZscYRHGBLhhe5zszrg3",
		Timestamp:   "1461722423785",
		Description: "测试专用优惠券",
		OrderNum:    "2016042710002341000238533",
		Type:        "coupon",
		FacePrice:   1,
		ActualPrice: 0,
		Ip:          "61.145.97.52",
		WaitAudit:   false,
		Params:      "",
		Sign:        "87445c85ea9e88ee85b5ebbb53a1c284",
	})
}

func TestDuibaReceiveCreditNotify(t *testing.T) {
	urlInfo, err := url.Parse("http://api.test.hongbeibang.com/client/exchangeResult?sign=b3eb7e8bf443e334356cb65b46494f79&uid=10014&timestamp=1461726097620&errorMessage=&orderNum=2016042711013690600695942&bizId=10227&success=true&appKey=24usbX8LZZscYRHGBLhhe5zszrg3&")

	sdk := &DuibaSdk{
		AppKey:    "24usbX8LZZscYRHGBLhhe5zszrg3",
		AppSecret: "4UFKw81Rev615GLtjPK4wf3hemE9",
	}
	result, err := sdk.ReceiveCreditNotify(urlInfo)
	assertDuibaSdkEqual(t, err, nil)
	assertDuibaSdkEqual(t, result, DuibaSdkReceiveCreditNotify{
		Uid:          "10014",
		AppKey:       "24usbX8LZZscYRHGBLhhe5zszrg3",
		Timestamp:    1461726097620,
		Success:      true,
		ErrorMessage: "",
		OrderNum:     "2016042711013690600695942",
		BizId:        "10227",
		Sign:         "b3eb7e8bf443e334356cb65b46494f79",
	})
}
