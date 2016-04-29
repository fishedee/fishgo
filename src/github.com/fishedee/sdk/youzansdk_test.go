package sdk

import (
	"reflect"
	"testing"
)

func assertYouzanSdkEqual(t *testing.T, left interface{}, right interface{}) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("assert fail: %+v != %+v", left, right)
	}
}

func TestYouzanSdkGetSign(t *testing.T) {
	sdk := &YouzanSdk{
		AppId:     "test",
		AppSecret: "test",
	}
	sign, err := sdk.getSign(map[string]interface{}{
		"method":      "kdt.item.get",
		"timestamp":   "2013-05-06 13:52:03",
		"format":      "json",
		"app_id":      "test",
		"v":           "1.0",
		"sign_method": "md5",
		"num_iid":     3838293428,
	})
	assertYouzanSdkEqual(t, err, nil)
	assertYouzanSdkEqual(t, sign, "74d4c18b9f077ed998feb10e96c58497")
}

func getYouzanSdk() *YouzanSdk {
	return &YouzanSdk{
		AppId:     "6d41aa6d167d19739c",
		AppSecret: "9ecf9ce46781840790990ee144220776",
	}
}

func TestYouzanSdkTrade(t *testing.T) {
	sdk := getYouzanSdk()
	result, err := sdk.GetTrade(YouzanSdkTradeRequest{
		Tid: "E20160412142348004391833",
	})
	assertYouzanSdkEqual(t, err, nil)
	t.Errorf("%#v", result)
}

func TestYouzanSdkTradeSold(t *testing.T) {
	sdk := getYouzanSdk()
	result, err := sdk.GetTradeSold(YouzanSdkTradeSoldRequest{
		//WeixinUserId: 44351341,
		//Keyword: "10006",
		//Status:    "TRADE_CLOSED_BY_USER",
		//WeixinUserId:   10006,
		PageSize: 10,
		PageNo:   0,
	})
	assertYouzanSdkEqual(t, err, nil)
	t.Errorf("%#v", result)
}

func TestYouzanSdkOauth(t *testing.T) {
	sdk := getYouzanSdk()
	loginUrl, err := sdk.GetOauthUrl("http://www.test.hongbeibang.com", "snsapi_userinfo", "")
	assertYouzanSdkEqual(t, err, nil)
	t.Errorf("%v", loginUrl)
}
