package weixin

import (
	"crypto/sha1"
	"fmt"
	. "github.com/fishedee/encoding"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type JsSdkModel struct {
}

var JsSdk = &JsSdkModel{}

func (this *JsSdkModel) GetAccessToken(appId string, appSecret string) string {
	accessToken := ""
	response, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appId + "&secret=" + appSecret)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var value map[string]interface{}
	err = DecodeJson(body, &value)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if v, ok := value["access_token"]; ok {
		accessToken = v.(string)
	} else {
		panic(err)
	}

	return accessToken
}

func (this *JsSdkModel) GetJsApiTicket(accessToken string) string {
	jsApiTicket := ""
	response, err := http.Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + accessToken + "&type=jsapi")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var value map[string]interface{}
	err = DecodeJson(body, &value)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if v, ok := value["ticket"]; ok {
		jsApiTicket = v.(string)
	} else {
		panic(err)
	}
	return jsApiTicket
}

func (this *JsSdkModel) GetJsConfig(appId string, jsApiTicket string, url string) JsConfig {
	if strings.Index(url, "#") != -1 {
		data := strings.Split(url, "#")
		url = data[0]
	}
	var jsConfig JsConfig
	noncestr := this.createNoncestr()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	jsSignature := JsSignature{
		JsApiTicket: noncestr,
		Noncestr:    jsApiTicket,
		Timestamp:   timestamp,
		Url:         url,
	}
	jsConfig = JsConfig{
		AppId:     appId,
		Noncestr:  noncestr,
		Timestamp: timestamp,
		Signature: this.getSignature(jsSignature),
	}
	return jsConfig
}

func (this *JsSdkModel) getSignature(jsSignature JsSignature) string {
	signature := ""
	signature = "jsapi_ticket=" + jsSignature.JsApiTicket + "&noncestr=" + jsSignature.Noncestr + "&timestamp=" + jsSignature.Timestamp + "&url=" + jsSignature.Url
	signature = this.sha1(signature)
	return signature
}

func (this *JsSdkModel) createNoncestr() string {
	chars := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	result := ""
	for i := 0; i < 32; i++ {
		key := rand.Intn(len(chars))
		result += string(chars[key])
	}
	return result
}

func (this *JsSdkModel) sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
