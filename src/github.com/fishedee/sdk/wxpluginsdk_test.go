package sdk

import (
	"encoding/xml"
	"fmt"
	. "github.com/fishedee/assert"
	"testing"
)

func TestWxPluginSdkCrypt(t *testing.T) {
	//准备
	msg := WxPluginSdkMessage{
		ToUserName:            "fish",
		AppId:                 "123",
		CreateTime:            1413192605,
		InfoType:              "component_verify_ticket",
		ComponentVerifyTicket: "123456",
	}
	encodingAesKey := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG"
	token := "pamtest"
	appId := "wxb11529c136998cb6"
	text, err := xml.Marshal(msg)
	AssertEqual(t, err, nil)

	//打包
	wxPluginSdk := &WxPluginSdk{
		ComponentAppId: appId,
		MessageToken:   token,
		MessageAESKey:  encodingAesKey,
	}
	result, err := wxPluginSdk.EncryptMessage(text)
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
	encryptXml := fmt.Sprintf(
		"<xml><ToUserName><![CDATA[fish]]></ToUserName><Encrypt><![CDATA[%s]]></Encrypt></xml>",
		resultObject.Encrypt)
	msg2, err := wxPluginSdk.DecryptMessage(resultObject.MsgSignature, resultObject.TimeStamp, resultObject.Nonce, []byte(encryptXml))
	AssertEqual(t, err, nil)
	AssertEqual(t, msg, msg2)
}
