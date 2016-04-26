package sdk

import (
	"reflect"
	"testing"
)

func assertQqSdkEqual(t *testing.T, left interface{}, right interface{}) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("assert fail: %+v != %+v", left, right)
	}
}

func testOauth(t *testing.T) {
	qqSdk := &QqSdk{
		AppId:     "101170133",
		AppSecret: "d5a95239584948ed04f016f69ba5d02a",
	}
	url, err := qqSdk.GetOauthUrl(
		"http://www.test.hongbeibang.com",
		"123",
		"get_user_info",
	)
	assertQqSdkEqual(t, err, nil)
	assertQqSdkEqual(t, url != "", true)
	t.Errorf("%#v", url)

	accessToken, err := qqSdk.GetOauthAccessToken(
		"http://www.test.hongbeibang.com",
		"586F54C9D41466C55DE294B47075E45B",
	)
	assertQqSdkEqual(t, err, nil)
	assertQqSdkEqual(t, accessToken.AccessToken != "", true)
	assertQqSdkEqual(t, accessToken.ExpiresIn != "", true)

	openId, err := qqSdk.GetOauthOpenId(
		accessToken.AccessToken,
	)
	assertQqSdkEqual(t, err, nil)
	assertQqSdkEqual(t, openId.ClientId != "", true)
	assertQqSdkEqual(t, openId.OpenId != "", true)

	userInfo, err := qqSdk.GetOauthUserInfo(
		accessToken.AccessToken,
		openId.OpenId,
	)
	assertQqSdkEqual(t, err, nil)
	assertQqSdkEqual(t, userInfo, QqSdkOauthUserInfo{
		NickName:        "Fish",
		Figureurl:       "http://qzapp.qlogo.cn/qzapp/101170133/FD39226A3C36F56DE959016F0F0603AB/30",
		Figureurl1:      "http://qzapp.qlogo.cn/qzapp/101170133/FD39226A3C36F56DE959016F0F0603AB/50",
		Figureurl2:      "http://qzapp.qlogo.cn/qzapp/101170133/FD39226A3C36F56DE959016F0F0603AB/100",
		FigureurlQq1:    "http://q.qlogo.cn/qqapp/101170133/FD39226A3C36F56DE959016F0F0603AB/40",
		FigureurlQq2:    "http://q.qlogo.cn/qqapp/101170133/FD39226A3C36F56DE959016F0F0603AB/100",
		Gender:          "男",
		IsYellowVip:     "0",
		Vip:             "0",
		YellowVipLevel:  "0",
		Level:           "0",
		IsYellowYearVip: "0",
	})
}

func TestOauth(t *testing.T) {
	testOauth(t)
}
