package sdk

import (
	"encoding/xml"
	"fmt"
	. "github.com/fishedee/assert"
	"testing"
)

func TestWxCryptSdk(t *testing.T) {
	//打包
	encodingAesKey := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG"
	token := "pamtest"
	timeStamp := "1409304348"
	nonce := "xxxxxx"
	appId := "wxb11529c136998cb6"
	text := "<xml><ToUserName><![CDATA[oia2Tj我是中文jewbmiOUlr6X-1crbLOvLw]]></ToUserName><FromUserName><![CDATA[gh_7f083739789a]]></FromUserName><CreateTime>1407743423</CreateTime><MsgType><![CDATA[video]]></MsgType><Video><MediaId><![CDATA[eYJ1MbwPRJtOvIEabaxHs7TX2D-HV71s79GUxqdUkjm6Gs2Ed1KF3ulAOA9H1xG0]]></MediaId><Title><![CDATA[testCallBackReplyVideo]]></Title><Description><![CDATA[testCallBackReplyVideo]]></Description></Video></xml>"
	wxCryptSdk := &WxCryptSdk{
		AppId:  appId,
		Token:  token,
		AESKey: encodingAesKey,
	}
	result, err := wxCryptSdk.Encrypt(timeStamp, nonce, []byte(text))
	AssertEqual(t, err, nil)

	//分析
	var resultObject struct {
		MsgSignature string `xml:"MsgSignature"`
		Encrypt      string `xml:"Encrypt"`
		TimeStamp    string `xml:"TimeStamp"`
		Nonce        string `xml:"Nonce"`
	}
	err = xml.Unmarshal(result, &resultObject)
	AssertEqual(t, err, nil)

	//解包
	encryptXml := fmt.Sprintf("<xml><ToUserName><![CDATA[toUser]]></ToUserName><Encrypt><![CDATA[%s]]></Encrypt></xml>", resultObject.Encrypt)
	toUserName, text2, err := wxCryptSdk.Decrypt(resultObject.MsgSignature, resultObject.TimeStamp, resultObject.Nonce, []byte(encryptXml))
	AssertEqual(t, err, nil)
	AssertEqual(t, toUserName, "toUser")
	AssertEqual(t, string(text2), text)
}
