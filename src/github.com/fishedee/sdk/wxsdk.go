package weixin

import (
	"crypto/sha1"
	"fmt"
	"errors"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type WxSdk struct{
	AppId string
	AppSecret string
}

type WxSdkToken struct{
	Errcode int `json:"errcode,omitempty"`
	Errmsg string `json:"errmsg,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn int `json:"expires_in,omitempty"`
}

type WxSdkJsTicket struct{
	Errcode int `json:"errcode,omitempty"`
	Errmsg string `json:"errmsg,omitempty"`
	Ticket string `json:"ticket,omitempty"`
	ExpiresIn int `json:"expires_in,omitempty"`
}

type WxSdkJsConfig struct {
	AppId     string
	NonceStr  string
	Timestamp string
	Signature string
}

type wxSdkJsSignature struct {
	JsApiTicket string
	NonceStr    string
	Timestamp   string
	Url         string
}

func (this *WxSdk) GetAccessToken() (WxSdkToken,error){
	appId := this.AppId
	appSecret := this.AppSecret
	result := WxSdkToken{}

	response, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appId + "&secret=" + appSecret)
	if err != nil {
		return result,err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result,err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result,err
	}
	if result.Errcode != 0{
		return result,errors.New("微信接口返回失败"+result.Errmsg)
	}

	return result,nil
}

func (this *WxSdk) GetJsApiTicket(accessToken string)(WxSdkJsTicket,error) {
	result := WxSdkJsTicket{}

	response, err := http.Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=" + accessToken + "&type=jsapi")
	if err != nil {
		return result,err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result,err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result,err
	}
	if result.Errcode != 0{
		return result,errors.New("微信接口返回失败"+result.Errmsg)
	}

	return result,nil
}

func (this *WxSdk) GetJsConfig(jsApiTicket string, url string)(WxSdkJsConfig,error) {
	appId := this.AppId
	
	if strings.Index(url, "#") != -1 {
		data := strings.Split(url, "#")
		url = data[0]
	}
	noncestr := this.createNoncestr()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	jsSignature := wxSdkJsSignature{
		JsApiTicket: jsApiTicket,
		NonceStr:    noncestr,
		Timestamp:   timestamp,
		Url:         url,
	}
	jsConfig := WxSdkJsConfig{
		AppId:     appId,
		NonceStr:  noncestr,
		Timestamp: timestamp,
		Signature: this.getSignature(jsSignature),
	}
	return jsConfig,nil
}

func (this *WxSdk) getSignature(jsSignature wxSdkJsSignature) string {
	signature := ""
	signature = "jsapi_ticket=" + jsSignature.JsApiTicket + "&noncestr=" + jsSignature.NonceStr + "&timestamp=" + jsSignature.Timestamp + "&url=" + jsSignature.Url
	signature = this.sha1(signature)
	return signature
}

func (this *WxSdk) createNoncestr() string {
	chars := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	result := ""
	for i := 0; i < 32; i++ {
		key := rand.Intn(len(chars))
		result += string(chars[key])
	}
	return result
}

func (this *WxSdk) sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
