package sdk

import (
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	"io"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type WxSdk struct {
	AppId     string
	AppSecret string
}

type WxSdkToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WxSdkJsTicket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

type WxSdkServerIp struct {
	IpList []string `json:"ip_list"`
}

type WxSdkUserInfo struct {
	Subscribe     int    `joson:"subscribe"`
	OpenId        string `joson:"openid"`
	NickName      string `joson:"nickname"`
	Sex           int    `joson:"sex"`
	Language      string `joson:"language"`
	City          string `joson:"city"`
	Province      string `joson:"province"`
	Country       string `joson:"country"`
	HeadImgUrl    string `joson:"headimgurl"`
	SubscribeTime int    `joson:"subscribe_time"`
	Unionid       string `joson:"unionid"`
	Remark        string `joson:"remark"`
	GroupId       int    `joson:"groupid"`
}

type WxSdkUserList struct {
	Total int `joson:"total"`
	Count int `joson:"count"`
	Data  struct {
		OpenId []string `joson:"openid"`
	} `joson:"data"`
	NextOpenid string `joson:"next_openid"`
}

type WxSdkMessage struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	MsgId        int64
	Content      string
	PicUrl       string
	MediaId      string
	Format       string
	Recognition  string
	ThumbMediaId string
	Location_X   float64
	Location_Y   float64
	Scale        int64
	Label        string
	Title        string
	Description  string
	Url          string
}

type WxSdkOauthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
}

type WxSdkOauthUserInfo struct {
	OpenId     string `joson:"openid"`
	NickName   string `joson:"nickname"`
	Sex        int    `joson:"sex"`
	Province   string `joson:"province"`
	City       string `joson:"city"`
	Country    string `joson:"country"`
	HeadImgUrl string `joson:"headimgurl"`
	Privilege  string `joson:"privilege"`
	Unionid    string `joson:"unionid"`
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

type wxSdkDownload struct {
	Errcode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}

type WxSdkError struct {
	Code    int
	Message string
}

func (this *WxSdkError) GetCode() int {
	return this.Code
}

func (this *WxSdkError) GetMsg() string {
	return this.Message
}

func (this *WxSdkError) Error() string {
	return fmt.Sprintf("错误码为：%v，错误描述为：%v", this.Code, this.Message)
}

func (this *WxSdk) api(method string, url string, query interface{}, data interface{}) ([]byte, error) {
	queryInfo, err := EncodeUrlQuery(query)
	if err != nil {
		return nil, err
	}
	url = "https://api.weixin.qq.com" + url
	if len(queryInfo) != 0 {
		url += "?" + string(queryInfo)
	}
	var result []byte
	if method == "GET" {
		ajaxOption := &Ajax{
			Url:          url,
			ResponseData: &result,
		}
		err = DefaultAjaxPool.Get(ajaxOption)
	} else {
		ajaxOption := &Ajax{
			Url:          url,
			Data:         data,
			ResponseData: &result,
		}
		err = DefaultAjaxPool.Post(ajaxOption)
	}
	if err != nil {
		return nil, err
	}
	var errInfo struct {
		ErrorCode int    `json:"errcode"`
		ErrorMsg  string `json:"errmsg"`
	}
	err = DecodeJson(result, &errInfo)
	if err == nil && errInfo.ErrorCode != 0 {
		return nil, &WxSdkError{errInfo.ErrorCode, errInfo.ErrorMsg}
	}

	return result, nil
}

func (this *WxSdk) apiJson(method string, url string, query interface{}, data interface{}, responseData interface{}) error {
	data, err := EncodeJson(data)
	if err != nil {
		return err
	}

	result, err := this.api(url, method, query, data)
	if err != nil {
		return err
	}
	err = DecodeJson(result, responseData)
	if err != nil {
		return err
	}
	return nil
}

//基础接口
func (this *WxSdk) CheckSignature(accessToken string, requestUrl url.URL) error {
	var input struct {
		Signature string `url:"signature"`
		Timestamp string `url:"timestamp"`
		Nonce     string `url:"nonce"`
	}
	queryString := requestUrl.RawQuery
	err := DecodeUrlQuery([]byte(queryString), &input)
	if err != nil {
		return err
	}
	arrayInfo := []string{accessToken, input.Timestamp, input.Nonce}
	arrayInfo = ArraySort(arrayInfo).([]string)
	arrayInfoString := Implode(arrayInfo, "")
	curSignature := this.sha1(arrayInfoString)
	if curSignature != input.Signature {
		return errors.New("invalid signature")
	} else {
		return nil
	}
}

func (this *WxSdk) GetAccessToken() (WxSdkToken, error) {
	result := WxSdkToken{}
	err := this.apiJson("GET", "/cgi-bin/token", map[string]string{
		"grant_type": "client_credential",
		"appid":      this.AppId,
		"secret":     this.AppSecret,
	}, nil, &result)
	if err != nil {
		return WxSdkToken{}, err
	}
	return result, nil
}

func (this *WxSdk) GetServerIp(accessToken string) (WxSdkServerIp, error) {
	result := WxSdkServerIp{}
	err := this.apiJson("GET", "/cgi-bin/getcallbackip", map[string]string{
		"access_token": accessToken,
	}, nil, &result)
	if err != nil {
		return WxSdkServerIp{}, err
	}
	return result, nil
}

//素材接口
func (this *WxSdk) DownloadMedia(accessToken, mediaId string) ([]byte, error) {
	return this.api("GET", "/cgi-bin/media/get", map[string]string{
		"access_token": accessToken,
		"media_id":     mediaId,
	}, nil)
}

//用户接口
func (this *WxSdk) getUserInfo(accessToken, openId string) (WxSdkUserInfo, error) {
	result := WxSdkUserInfo{}
	err := this.apiJson("GET", "/cgi-bin/user/info", map[string]string{
		"access_token": accessToken,
		"openid":       openId,
		"lang":         "zh_CN",
	}, nil, &result)
	if err != nil {
		return WxSdkUserInfo{}, err
	}
	return result, nil
}

func (this *WxSdk) getUserList(accessToken, next_openid string) (WxSdkUserList, error) {
	argv := map[string]string{
		"access_token": accessToken,
	}
	if next_openid != "" {
		argv["next_openid"] = next_openid
	}
	result := WxSdkUserList{}
	err := this.apiJson("GET", "/cgi-bin/user/get", argv, nil, &result)
	if err != nil {
		return WxSdkUserList{}, err
	}
	return result, nil
}

//Message接口
func (this *WxSdk) receiveMessage(data []byte) (WxSdkMessage, error) {
	var result WxSdkMessage
	err := xml.Unmarshal(data, &result)
	return result, err
}

func (this *WxSdk) sendMessage(message WxSdkMessage) ([]byte, error) {
	return xml.Marshal(message)
}

//Menu接口
func (this *WxSdk) setMenu(accessToken string, data string) error {
	_, err := this.api("POST", "/cgi-bin/menu/create", map[string]string{
		"access_token": accessToken,
	}, data)
	if err != nil {
		return err
	}
	return nil
}

func (this *WxSdk) getMenu(accessToken string) (string, error) {
	data, err := this.api("GET", "/cgi-bin/menu/get", map[string]string{
		"access_token": accessToken,
	}, nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (this *WxSdk) delMenu(accessToken string) error {
	_, err := this.api("GET", "/cgi-bin/menu/delete", map[string]string{
		"access_token": accessToken,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

//OAuth接口
func (this *WxSdk) getLoginUrl(callback string, state string, scope string) (string, error) {
	query := map[string]string{
		"appid":         this.AppId,
		"redirect_uri":  callback,
		"response_type": "code",
		"scope":         scope,
		"state":         state,
	}
	queryStr, err := EncodeUrlQuery(query)
	if err != nil {
		return "", err
	}
	return "https://open.weixin.qq.com/connect/oauth2/authorize?" + string(queryStr), nil
}

func (this *WxSdk) getPcLoginUrl(callback string, state string, scope string) (string, error) {
	query := map[string]string{
		"appid":         this.AppId,
		"redirect_uri":  callback,
		"response_type": "code",
		"scope":         scope,
		"state":         state,
	}
	queryStr, err := EncodeUrlQuery(query)
	if err != nil {
		return "", err
	}
	return "https://open.weixin.qq.com/connect/qrconnect?" + string(queryStr), nil
}

func (this *WxSdk) getOauhToken(code string) (WxSdkOauthToken, error) {
	result := WxSdkOauthToken{}
	err := this.apiJson("GET", "/sns/oauth2/access_token", map[string]string{
		"appid":      this.AppId,
		"secret":     this.AppSecret,
		"code":       code,
		"grant_type": "authorization_code",
	}, nil, &result)
	if err != nil {
		return WxSdkOauthToken{}, err
	}
	return result, nil
}

func (this *WxSdk) getOauthUserInfo(accessToken, openid string) (WxSdkOauthUserInfo, error) {
	result := WxSdkOauthUserInfo{}
	err := this.apiJson("GET", "/sns/userinfo", map[string]string{
		"access_token": accessToken,
		"openid":       openid,
		"lang":         "zh_CN",
	}, nil, &result)
	if err != nil {
		return WxSdkOauthUserInfo{}, err
	}
	return result, nil
}

//Js接口
func (this *WxSdk) GetJsApiTicket(accessToken string) (WxSdkJsTicket, error) {
	result := WxSdkJsTicket{}
	err := this.apiJson("GET", "/cgi-bin/ticket/getticket", map[string]string{
		"access_token": accessToken,
		"type":         "jsapi",
	}, nil, &result)
	if err != nil {
		return WxSdkJsTicket{}, err
	}
	return result, nil
}

func (this *WxSdk) GetJsConfig(jsApiTicket string, url string) (WxSdkJsConfig, error) {
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
	return jsConfig, nil
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
