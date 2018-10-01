package sdk

import (
	"encoding/json"
	"encoding/xml"
	"net/url"
	"reflect"
	"testing"
)

func assertWxSdkEqual(t *testing.T, left interface{}, right interface{}) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("assert fail: %+v != %+v", left, right)
	}
}
func getWxSdk() *WxSdk {
	return &WxSdk{
		AppId:     "wxa9b4bcb4be4695d7",
		AppSecret: "e045f8cc0d5e03d5fefbe0d18d1f4ef3",
		Token:     "JxTO4tR1OJUdoKaE",
	}
}

func testBasicInner(t *testing.T, wxSdk *WxSdk) string {
	//基础接口
	urlInfo, _ := url.Parse("http://api.hongbeibang.com/weixin/getMessage?signature=82e52de654bff034e3865a8fb3e6fc3dbf3eeebd&timestamp=1461633261&nonce=1773859070")
	err := wxSdk.CheckSignature(urlInfo)
	assertWxSdkEqual(t, err == nil, true)

	accessToken, err := wxSdk.GetAccessToken()
	assertWxSdkEqual(t, err == nil, true)
	assertWxSdkEqual(t, accessToken.AccessToken != "", true)
	assertWxSdkEqual(t, accessToken.ExpiresIn != 0, true)

	serverIp, err := wxSdk.GetServerIp(accessToken.AccessToken)
	assertWxSdkEqual(t, err == nil, true)
	assertWxSdkEqual(t, len(serverIp.IpList) != 0, true)

	return accessToken.AccessToken
}

func testUser(t *testing.T, wxSdk *WxSdk, accessToken string) {
	//用户接口
	userInfo, err := wxSdk.GetUserInfo(accessToken, "oi0NVuHgWqPRWeM6Xn589VFvbaOY")
	assertWxSdkEqual(t, err == nil, true)
	assertWxSdkEqual(t, userInfo, WxSdkUserInfo{
		Subscribe:     1,
		OpenId:        "oi0NVuHgWqPRWeM6Xn589VFvbaOY",
		NickName:      "fish@烘焙帮",
		Sex:           1,
		Language:      "zh_CN",
		City:          "",
		Province:      "",
		Country:       "埃塞俄比亚",
		HeadImgUrl:    "http://wx.qlogo.cn/mmopen/TBkSnIpf1E9AP4Kbz3PHSgZNldw7ZzxyNibneguA7QcXyC7N2HbKMKj7fQoWqib4EzyQAvKBbgHXvr7syOibXAZI7ZAINoOIp6m/0",
		SubscribeTime: 1431501183,
		UnionId:       "",
		Remark:        "",
		GroupId:       0,
	})

	userInfoList, err := wxSdk.GetUserList(accessToken, "")
	assertWxSdkEqual(t, err == nil, true)
	assertWxSdkEqual(t, userInfoList.Total != 0, true)
	assertWxSdkEqual(t, userInfoList.Count != 0, true)
	assertWxSdkEqual(t, userInfoList.NextOpenid != "", true)
	for _, singleOpenId := range userInfoList.Data.OpenId {
		assertWxSdkEqual(t, singleOpenId != "", true)
	}
}

func testMessage(t *testing.T, wxSdk *WxSdk, accessToken string) {
	//反序列化
	recvMsg, err := wxSdk.ReceiveMessage([]byte(`<xml>
		 <ToUserName><![CDATA[toUser]]></ToUserName>
		 <FromUserName><![CDATA[fromUser]]></FromUserName>
		 <CreateTime>1348831860</CreateTime>
		 <MsgType><![CDATA[text]]></MsgType>
		 <Content><![CDATA[this is a test]]></Content>
		 <MsgId>1234567890123456</MsgId>
		</xml>
	`))
	assertWxSdkEqual(t, err == nil, true)
	assertWxSdkEqual(t, recvMsg, WxSdkReceiveMessage{
		ToUserName:   "toUser",
		FromUserName: "fromUser",
		CreateTime:   1348831860,
		MsgType:      "text",
		Content:      "this is a test",
		MsgId:        1234567890123456,
	})

	//序列化
	sendMessage := WxSdkSendMessage{
		ToUserName:   "toUser",
		FromUserName: "fromUser",
		CreateTime:   12345678,
		MsgType:      "news",
		ArticleCount: 2,
		Articles: []WxSdkSendArticleMessage{
			{"title", "description1", "picurl", "url"},
			{"title", "description", "picurl", "url"},
		},
	}
	data, err := wxSdk.SendMessage(sendMessage)
	assertWxSdkEqual(t, err == nil, true)
	var dataStruct struct {
		ToUserName   string
		FromUserName string
		CreateTime   int
		MsgType      string
		ArticleCount int
		Articles     []struct {
			Title       string
			Description string
			PicUrl      string
			Url         string
		} `xml:"Articles>item"`
	}
	err = xml.Unmarshal(data, &dataStruct)
	assertWxSdkEqual(t, err == nil, true)
	assertWxSdkEqual(t, sendMessage.ToUserName, dataStruct.ToUserName)
	assertWxSdkEqual(t, sendMessage.FromUserName, dataStruct.FromUserName)
	assertWxSdkEqual(t, sendMessage.CreateTime, dataStruct.CreateTime)
	assertWxSdkEqual(t, sendMessage.MsgType, dataStruct.MsgType)
	assertWxSdkEqual(t, sendMessage.ArticleCount, dataStruct.ArticleCount)
	assertWxSdkEqual(t, len(sendMessage.Articles), len(dataStruct.Articles))
}

func testMenu(t *testing.T, wxSdk *WxSdk, accessToken string) {
	err := wxSdk.DelMenu(accessToken)
	assertWxSdkEqual(t, err == nil, true)

	err = wxSdk.SetMenu(accessToken, `{
	     "button":[
	     {
	          "type":"click",
	          "name":"今日歌曲",
	          "key":"V1001_TODAY_MUSIC"
	      },
	      {
	           "name":"菜单",
	           "sub_button":[
	           {
	               "type":"view",
	               "name":"搜索",
	               "url":"http://www.soso.com/"
	            },
	            {
	               "type":"view",
	               "name":"视频",
	               "url":"http://v.qq.com/"
	            },
	            {
	               "type":"click",
	               "name":"赞一下我们",
	               "key":"V1001_GOOD"
	            }]
	       }]
	}`)
	assertWxSdkEqual(t, err == nil, true)

	data, err := wxSdk.GetMenu(accessToken)
	assertWxSdkEqual(t, err == nil, true)
	var result interface{}
	err = json.Unmarshal([]byte(data), &result)
	assertWxSdkEqual(t, err == nil, true)
	assertWxSdkEqual(t, result, map[string]interface{}{
		"button": []interface{}{
			map[string]interface{}{
				"type":       "click",
				"name":       "今日歌曲",
				"key":        "V1001_TODAY_MUSIC",
				"sub_button": []interface{}{},
			},
			map[string]interface{}{
				"name": "菜单",
				"sub_button": []interface{}{
					map[string]interface{}{
						"type":       "view",
						"name":       "搜索",
						"url":        "http://www.soso.com/",
						"sub_button": []interface{}{},
					},
					map[string]interface{}{
						"type":       "view",
						"name":       "视频",
						"url":        "http://v.qq.com/",
						"sub_button": []interface{}{},
					},
					map[string]interface{}{
						"type":       "click",
						"name":       "赞一下我们",
						"key":        "V1001_GOOD",
						"sub_button": []interface{}{},
					},
				},
			},
		},
	})
}

func testWxOauth(t *testing.T, wxSdk *WxSdk, accessToken string) {
	url, err := wxSdk.GetOauthUrl("http://api.test.hongbeibang.com/login/wxcallback", "tt", "snsapi_userinfo")
	assertWxSdkEqual(t, err == nil, true)
	t.Errorf("%v", url)

	token, err := wxSdk.GetOauthToken("021HOokS0Qjj8d2BejkS0dppkS0HOokS")
	assertWxSdkEqual(t, err == nil, true)
	assertWxSdkEqual(t, token.AccessToken != "", true)
	assertWxSdkEqual(t, token.OpenId != "", true)

	userInfo, err := wxSdk.GetOauthUserInfo(token.AccessToken, token.OpenId)
	assertWxSdkEqual(t, err == nil, true)
	assertWxSdkEqual(t, userInfo, WxSdkOauthUserInfo{
		OpenId:     "oi0NVuHgWqPRWeM6Xn589VFvbaOY",
		NickName:   "fish@烘焙帮",
		Sex:        1,
		Province:   "",
		City:       "",
		Country:    "中国",
		HeadImgUrl: "http://wx.qlogo.cn/mmopen/TBkSnIpf1E9AP4Kbz3PHSgZNldw7ZzxyNibneguA7QcXyC7N2HbKMKj7fQoWqib4EzyQAvKBbgHXvr7syOibXAZI7ZAINoOIp6m/0",
		Privilege:  []string{},
		UnionId:    "",
	})
}

func testJs(t *testing.T, wxSdk *WxSdk, accessToken string) {
	ticket, err := wxSdk.GetJsApiTicket(accessToken)
	assertWxSdkEqual(t, err == nil, true)
	assertWxSdkEqual(t, ticket.Ticket != "", true)
	assertWxSdkEqual(t, ticket.ExpiresIn != 0, true)

	sig := getWxSdk().getSignature(wxSdkJsSignature{
		JsApiTicket: "sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90u5h9nbSlYy3-Sl-HhTdfl2fzFy1AOcHKP7qg",
		NonceStr:    "Wm3WZYTPz0wzccnW",
		Timestamp:   "1414587457",
		Url:         "http://mp.weixin.qq.com?params=value",
	})
	target := "0f9de62fce790f9a083d5c99e95740ceb90c27ed"
	assertWxSdkEqual(t, sig, target)
}

func testBasic(t *testing.T) {
	wxSdk := getWxSdk()

	accessToken := testBasicInner(t, wxSdk)
	testUser(t, wxSdk, accessToken)
	testMessage(t, wxSdk, accessToken)
	testMenu(t, wxSdk, accessToken)
	//testOauth(t, wxSdk, "")
	testJs(t, wxSdk, accessToken)
}
