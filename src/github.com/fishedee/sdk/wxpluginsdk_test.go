package sdk

import (
	"encoding/xml"
	"fmt"
	. "github.com/fishedee/assert"
	"testing"
)

func TestWxPluginSdkCrypt2(t *testing.T) {
	wxPluginSdk := &WxPluginSdk{
		ComponentAppId: "wxd2f16fa1336812b4",
		MessageToken:   "outswxtest",
		MessageAESKey:  "as23qdfeq32efge423LKIEnadsnfgnaseweewuuu22b",
	}
	encryptXml := `<xml>
    <AppId><![CDATA[wxd2f16fa1336812b4]]></AppId>
    <Encrypt><![CDATA[seGHRtsa4o0LTB4tJN309putKTKUeQUUJDh00WiHtLxrWGnOUPcCIX2GqwHoo+ScvWy+Y0Xn4QDMUwiKs2elUFOqEvE6Z0TxVaLQu4Ad5DIk18+jw0QuVC47AzUDthwMPvKOri0p2OVTdlk3S9JZwAfz2zm/RfPOsKaauQKj0NLeJGQT6ijzm3Tfn+FRH/t6IBNR2bEWg1gGCVNsP2offyNH+9XVA/PnQ5J+TQ85sGaaQYl1xjCX8+9Gfi6OanzJP/AyT5rjbv2OpKWcqp9Lsyen6bYXh6iKW9nDron0Mw3Hp7PPojE4v2cHT7SjzS7l3AaY47TG9z0rS1ppx+6Mg3uKpBdSKNm+VJVIMJMa9N6m+pyz+UvHl5v08fSFx8l6kyCF79MxvANiBZlBnjN+CBWGHIWPGlG07GHP9Ganc15IzSb0Hu1W+G4/f7ky1epCQ+Yg7A140NAtXDSNWJf+9w==]]></Encrypt>
</xml>`
	//signature := "db37ef95af48728d9ef78daf6d4ed70f9b767f42"
	timestamp := "1539147826"
	nonce := "1382272705"
	msgSignature := "f71541607f56c5f7b2638d3fda2620f525d0fb94"
	msg, err := wxPluginSdk.DecryptMessage(msgSignature,
		timestamp,
		nonce, []byte(encryptXml))
	AssertEqual(t, err, nil)
	AssertEqual(t, msg, WxPluginSdkMessage{
		AppId:                 "wxd2f16fa1336812b4",
		CreateTime:            1539147826,
		InfoType:              "component_verify_ticket",
		ComponentVerifyTicket: "ticket@@@ezQMYDedZQ_Ertr28jVFK51wmBJb0Up70LYupzCwjtS5DdRFwRI9fjFbVvUfElYZqTmtzVKcdmGFF66gl1Ex0w",
	})
}
func TestWxPluginSdkCrypt(t *testing.T) {
	//准备
	msg := WxPluginSdkMessage{
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
		"<xml><AppId><![CDATA[fish]]></AppId><Encrypt><![CDATA[%s]]></Encrypt></xml>",
		resultObject.Encrypt)
	msg2, err := wxPluginSdk.DecryptMessage(resultObject.MsgSignature, resultObject.TimeStamp, resultObject.Nonce, []byte(encryptXml))
	AssertEqual(t, err, nil)
	AssertEqual(t, msg, msg2)
}
