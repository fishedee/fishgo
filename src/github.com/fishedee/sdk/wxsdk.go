package sdk

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
)

type WxSdk struct {
	AppId       string
	AppSecret   string
	Token       string
	AccessToken string
	JsApiTicket string
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
	Subscribe     int    `json:"subscribe"`
	OpenId        string `json:"openid"`
	NickName      string `json:"nickname"`
	Sex           int    `json:"sex"`
	Language      string `json:"language"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	HeadImgUrl    string `json:"headimgurl"`
	SubscribeTime int    `json:"subscribe_time"`
	UnionId       string `json:"unionid"`
	Remark        string `json:"remark"`
	GroupId       int    `json:"groupid"`
}

type WxMiniProgramSession struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
}

type WxMiniProgramUserInfo struct {
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	Language  string `json:"language"`
	Watermark struct {
		Timestamp int64  `json:"timestamp"`
		AppID     string `json:"appid"`
	} `json:"watermark"`
}

type WxSdkUserList struct {
	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		OpenId []string `json:"openid"`
	} `json:"data"`
	NextOpenid string `json:"next_openid"`
}

// 微信消息模板
type WxSdkTemplateMessage struct {
	Touser      string                                     `json:"touser"`
	TemplateId  string                                     `json:"template_id"`
	Url         string                                     `json:"url"`
	Miniprogram WxSdkTemplateMessageMiniprogram            `json:"miniprogram,omitempty"`
	Data        map[string]WxSdkTemplateMessageDataContent `json:"data"`
}

type WxSdkTemplateMessageMiniprogram struct {
	Appid    string `json:"appid"`
	Pagepath string `json:"pagepath"`
}

type WxSdkTemplateMessageDataContent struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// 微信素材列表
type WxSdkSendBatchgetMaterial struct {
	Type   string `json:"type"` // 对应WxSdkMaterialType
	Offset int    `json:"offset"`
	Count  int    `json:"count"`
}

var WxSdkMaterialType = struct {
	Image string
	Video string
	Voice string
	News  string
}{
	Image: "image",
	Video: "video",
	Voice: "voice",
	News:  "news",
}

type WxSdkReceiveBatchgetMaterial struct {
	TotalCount int `json:"total_count"`
	ItemCount  int `json:"item_count"`
	Item       []struct {
		MediaId string `json:"media_id"`
		Name    string `json:"name"`
		Content []struct {
			NewsItem []struct {
				ThumbMediaId     string `json:"thumb_media_id"`
				Title            string `json:"title"`
				ShowCoverPic     string `json:"show_cover_pic"`
				Author           string `json:"author"`
				Content          string `json:"content"`
				Digest           string `json:"digest"`
				Url              string `json:"url"`
				ContentSourceUrl string `json:"content_source_url"`
			} `json:"news_item"`
		} `json:"content"`
		Url        string `json:"url"`
		UpdateTime string `json:"update_time"`
	} `json:"item"`
}

type WxSdkAddMaterial struct {
	MediaId string `json:"media_id"`
	Url     string `json:"url"`
}

type WxSdkReceiveMaterial struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DownUrl     string `json:"down_url"`
}

type WxSdkReceiveMessage struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
	MsgId        int
	MsgID        int
	//文本消息
	Content string
	//图片消息
	PicUrl  string
	MediaId string
	//语音消息
	Format      string
	Recognition string
	//视频消息
	ThumbMediaId string
	//地理位置消息
	Location_X float64
	Location_Y float64
	Scale      int
	Label      string
	//链接消息
	Title       string
	Description string
	Url         string
	//事件消息
	Event     string
	EventKey  string
	Ticket    string
	Latitude  float64
	Longitude float64
	Precision float64
}

type WxSdkSendImageMessage struct {
	MediaId string `xml:"MediaId,omitempty"`
}

type WxSdkSendVoiceMessage struct {
	MediaId string `xml:"MediaId,omitempty"`
}

type WxSdkSendVideoMessage struct {
	MediaId     string `xml:"MediaId,omitempty"`
	Title       string `xml:"MediaId,title"`
	Description string `xml:"MediaId,description"`
}

type WxSdkSendMusicMessage struct {
	Title        string `xml:"Title,omitempty"`
	Description  string `xml:"Description,omitempty"`
	MusicUrl     string `xml:"MusicUrl,omitempty"`
	HQMusicUrl   string `xml:"HQMusicUrl,omitempty"`
	ThumbMediaId string `xml:"ThumbMediaId,omitempty"`
}

type WxSdkSendArticleMessage struct {
	Title       string `xml:"Title,omitempty"`
	Description string `xml:"Description,omitempty"`
	PicUrl      string `xml:"PicUrl,omitempty"`
	Url         string `xml:"Url,omitempty"`
}

type WxSdkSendMessage struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int    `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	//文本消息
	Content string `xml:"Content,omitempty"`
	//图片消息
	Image WxSdkSendImageMessage `xml:"Image,omitempty"`
	//语音消息
	Voice WxSdkSendVoiceMessage `xml:"Voice,omitempty"`
	//视频消息
	Video WxSdkSendVideoMessage `xml:"Video,omitempty"`
	//视频消息
	Music WxSdkSendMusicMessage `xml:"Music,omitempty"`
	//图文消息
	ArticleCount int                       `xml:"ArticleCount,omitempty"`
	Articles     []WxSdkSendArticleMessage `xml:"Articles,omitempty"`
}

type WxSdkSendQrcode struct {
	ExpireSeconds int    `json:"expire_seconds"`
	ActionName    string `json:"action_name"`
	ActionInfo    struct {
		Scene struct {
			SceneId  int    `json:"scene_id,omitempty"`
			SceneStr string `json:"scene_str,omitempty"`
		} `json:"scene"`
	} `json:"action_info"`
}

type WxSdkReceiveQrcode struct {
	Ticket        string `json:"ticket"`
	ExpireSeconds int    `json:"expire_seconds"`
	Url           string `json:"url"`
}

type WxSdkOauthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
}

type WxSdkOauthUserInfo struct {
	OpenId     string   `json:"openid"`
	NickName   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgUrl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionId    string   `json:"unionid"`
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

type WxSdkCommonResult struct {
	Errcode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
	MsgID   int64  `json:"msgid,omitempty"`
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

func (this *WxSdk) api(method string, url string, query interface{}, dataType string, data interface{}) ([]byte, error) {
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
			DataType:     dataType,
		}
		err = DefaultAjaxPool.Get(ajaxOption)
	} else {
		ajaxOption := &Ajax{
			Url:          url,
			Data:         data,
			ResponseData: &result,
			DataType:     dataType,
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

func (this *WxSdk) apiJson(method string, url string, query interface{}, dataType string, data interface{}, responseData interface{}) error {

	data, err := EncodeJson(data)
	if err != nil {
		return err
	}

	result, err := this.api(method, url, query, dataType, data)
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
func (this *WxSdk) CheckSignature(requestUrl *url.URL) error {
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
	arrayInfo := []string{this.Token, input.Timestamp, input.Nonce}
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
	}, "", nil, &result)
	if err != nil {
		return WxSdkToken{}, err
	}
	return result, nil
}

func (this *WxSdk) GetServerIp(accessToken string) (WxSdkServerIp, error) {
	result := WxSdkServerIp{}
	err := this.apiJson("GET", "/cgi-bin/getcallbackip", map[string]string{
		"access_token": accessToken,
	}, "", nil, &result)
	if err != nil {
		return WxSdkServerIp{}, err
	}
	return result, nil
}

// 获取素材列表
func (this *WxSdk) GetBatchgetMaterial(data WxSdkSendBatchgetMaterial) (WxSdkReceiveBatchgetMaterial, error) {
	var result WxSdkReceiveBatchgetMaterial
	err := this.apiJson("POST", "/cgi-bin/material/batchget_material", map[string]string{
		"access_token": this.AccessToken,
	}, "", data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// 添加素材
// 参数 fileAddress 上传文件的本地地址
// 参数 materialType WxSdkMaterialType类型
// 参数 title 多媒体标题
// 参数 introduction 多媒体描述
func (this *WxSdk) AddMaterial(fileAddress, materialType, title, introduction string) (WxSdkAddMaterial, error) {

	var result WxSdkAddMaterial

	description := struct {
		Title        string `json:"title"`
		Introduction string `json:"introduction"`
	}{
		Title:        title,
		Introduction: introduction,
	}

	descriptionJson, err := EncodeJson(description)
	if err != nil {
		return result, err
	}

	cmd := exec.Command("curl", "https://api.weixin.qq.com/cgi-bin/material/add_material?access_token="+this.AccessToken+"&type="+materialType,
		"-F", "media=@"+fileAddress,
		"-F", "description="+string(descriptionJson))
	cmd.Stdin = strings.NewReader("some input")
	var outCmd bytes.Buffer
	var errCmd bytes.Buffer
	cmd.Stdout = &outCmd
	cmd.Stderr = &errCmd
	err = cmd.Run()

	if err != nil {
		return result, errors.New("err:" + err.Error() + ",errCmd:" + string(errCmd.Bytes()) + ",outCmd:" + string(outCmd.Bytes()))
	}

	var errInfo struct {
		ErrorCode int    `json:"errcode"`
		ErrorMsg  string `json:"errmsg"`
	}

	err = DecodeJson(outCmd.Bytes(), &errInfo)
	if err != nil {
		return result, err
	}

	if errInfo.ErrorCode != 0 {
		return result, &WxSdkError{Code: errInfo.ErrorCode, Message: errInfo.ErrorMsg}
	}

	err = DecodeJson(outCmd.Bytes(), &result)
	return result, err

}

// 获取素材接口
// 参数 mediaId 多媒体的media_id
func (this *WxSdk) GetMaterial(mediaId string) (WxSdkReceiveMaterial, error) {
	var result WxSdkReceiveMaterial
	err := NewAjaxPool(&AjaxPoolOption{}).Post(&Ajax{
		Url: "https://api.weixin.qq.com/cgi-bin/material/get_material?access_token=" + this.AccessToken,
		Data: map[string]interface{}{
			"media_id": mediaId,
		},
		ResponseDataType: "json",
		ResponseData:     &result,
	})
	return result, err
}

// 删除素材列表
func (this *WxSdk) DelMaterial(mediaId string) ([]byte, error) {
	return this.api("POST", "/cgi-bin/material/del_material",
		map[string]string{
			"access_token": this.AccessToken,
		}, "json",
		map[string]string{
			"media_id": mediaId,
		})
}

//获取素材接口
func (this *WxSdk) DownloadMedia(accessToken, mediaId string) ([]byte, error) {
	return this.api("GET", "/cgi-bin/media/get", map[string]string{
		"access_token": accessToken,
		"media_id":     mediaId,
	}, "", nil)
}

//用户接口
func (this *WxSdk) GetUserInfo(accessToken, openId string) (WxSdkUserInfo, error) {
	result := WxSdkUserInfo{}
	err := this.apiJson("GET", "/cgi-bin/user/info", map[string]string{
		"access_token": accessToken,
		"openid":       openId,
		"lang":         "zh_CN",
	}, "", nil, &result)
	if err != nil {
		return WxSdkUserInfo{}, err
	}
	return result, nil
}

func (this *WxSdk) GetUserList(accessToken, next_openid string) (WxSdkUserList, error) {
	argv := map[string]string{
		"access_token": accessToken,
	}
	if next_openid != "" {
		argv["next_openid"] = next_openid
	}
	result := WxSdkUserList{}
	err := this.apiJson("GET", "/cgi-bin/user/get", argv, "", nil, &result)
	if err != nil {
		return WxSdkUserList{}, err
	}
	return result, nil
}

//消息接口
func (this *WxSdk) ReceiveMessage(data []byte) (WxSdkReceiveMessage, error) {
	var result WxSdkReceiveMessage
	err := xml.Unmarshal(data, &result)
	return result, err
}

func (this *WxSdk) generateXml(data interface{}) string {
	result := ""
	if mapData, isOk := data.(map[string]interface{}); isOk {
		for key, value := range mapData {
			result += "<" + key + ">"
			result += this.generateXml(value)
			result += "</" + key + ">"
		}
		return result
	} else if arrayData, isOk := data.([]interface{}); isOk {
		for _, singleData := range arrayData {
			result += "<item>"
			result += this.generateXml(singleData)
			result += "</item>"
		}
		return result
	} else if stringData, isOk := data.(string); isOk {
		return "<![CDATA[" + stringData + "]]>"
	} else {
		return fmt.Sprintf("%v", data)
	}
}

func (this *WxSdk) SendMessage(message WxSdkSendMessage) ([]byte, error) {
	data := ArrayToMap(message, "xml")
	body := this.generateXml(data)
	result := []byte("<xml>" + body + "</xml>")
	return result, nil
}

func (this *WxSdk) SendPairMessage(requestUrl *url.URL) ([]byte, error) {
	var input struct {
		EchoStr string `url:"echostr"`
	}
	queryString := requestUrl.RawQuery
	err := DecodeUrlQuery([]byte(queryString), &input)
	if err != nil {
		return nil, err
	}
	return []byte(input.EchoStr), nil
}

//菜单接口
func (this *WxSdk) SetMenu(accessToken string, data string) error {
	_, err := this.api("POST", "/cgi-bin/menu/create", map[string]string{
		"access_token": accessToken,
	}, "", data)
	if err != nil {
		return err
	}
	return nil
}

func (this *WxSdk) GetMenu(accessToken string) (string, error) {
	data, err := this.api("GET", "/cgi-bin/menu/get", map[string]string{
		"access_token": accessToken,
	}, "", nil)
	if err != nil {
		return "", err
	}
	var result interface{}
	err = DecodeJson(data, &result)
	if err != nil {
		return "", err
	}
	resultJson := result.(map[string]interface{})
	data, err = EncodeJson(resultJson["menu"])
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (this *WxSdk) DelMenu(accessToken string) error {
	_, err := this.api("GET", "/cgi-bin/menu/delete", map[string]string{
		"access_token": accessToken,
	}, "", nil)
	if err != nil {
		return err
	}
	return nil
}

// 发送消息模板
func (this *WxSdk) SendTemplateMessage(accessToken string, msgData WxSdkTemplateMessage) (WxSdkCommonResult, error) {

	result := WxSdkCommonResult{}

	msgJson, err := EncodeJson(msgData)
	if err != nil {
		return result, err
	}

	data, err := this.api("POST", "/cgi-bin/message/template/send", map[string]string{
		"access_token": accessToken,
	}, "", msgJson)
	if err != nil {
		return result, err
	}

	err = DecodeJson(data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

//手动拼接参数
func (this *WxSdk) getOauthUrlQuery(query map[string]string) string {
	queryString := ""

	sorted_keys := make([]string, 0)
	for k, _ := range query {
		sorted_keys = append(sorted_keys, k)
	}

	//对key排序
	sort.Strings(sorted_keys)

	for _, key := range sorted_keys {
		keyEncode, err := EncodeUrl(key)
		if err != nil {
			continue
		}
		dataEncode, err := EncodeUrl(query[key])
		if err != nil {
			continue
		}
		queryString += keyEncode + "=" + dataEncode + "&"
	}
	return strings.Trim(queryString, "&")
}

//OAuth接口
func (this *WxSdk) GetOauthUrl(callback string, state string, scope string) (string, error) {
	query := map[string]string{
		"appid":         this.AppId,
		"redirect_uri":  callback,
		"response_type": "code",
		"scope":         scope,
		"state":         state,
	}
	queryStr := this.getOauthUrlQuery(query)
	return "https://open.weixin.qq.com/connect/oauth2/authorize?" + string(queryStr), nil
}

func (this *WxSdk) GetPcOauthUrl(callback string, state string, scope string) (string, error) {
	query := map[string]string{
		"appid":         this.AppId,
		"redirect_uri":  callback,
		"response_type": "code",
		"scope":         scope,
		"state":         state,
	}
	queryStr := this.getOauthUrlQuery(query)
	return "https://open.weixin.qq.com/connect/qrconnect?" + string(queryStr), nil
}

func (this *WxSdk) GetOauthToken(code string) (WxSdkOauthToken, error) {
	result := WxSdkOauthToken{}
	err := this.apiJson("GET", "/sns/oauth2/access_token", map[string]string{
		"appid":      this.AppId,
		"secret":     this.AppSecret,
		"code":       code,
		"grant_type": "authorization_code",
	}, "", nil, &result)
	if err != nil {
		return WxSdkOauthToken{}, err
	}
	return result, nil
}

func (this *WxSdk) GetOauthUserInfo(accessToken, openid string) (WxSdkOauthUserInfo, error) {
	result := WxSdkOauthUserInfo{}
	err := this.apiJson("GET", "/sns/userinfo", map[string]string{
		"access_token": accessToken,
		"openid":       openid,
		"lang":         "zh_CN",
	}, "", nil, &result)
	if err != nil {
		return WxSdkOauthUserInfo{}, err
	}
	return result, nil
}

// 创建二维码
func (this *WxSdk) AddQrcode(data WxSdkSendQrcode) (WxSdkReceiveQrcode, error) {
	result := WxSdkReceiveQrcode{}

	err := this.apiJson("POST", "/cgi-bin/qrcode/create", map[string]string{
		"access_token": this.AccessToken,
	}, "", data, &result)
	return result, err
}

//Js接口
func (this *WxSdk) GetJsApiTicket(accessToken string) (WxSdkJsTicket, error) {
	result := WxSdkJsTicket{}
	err := this.apiJson("GET", "/cgi-bin/ticket/getticket", map[string]string{
		"access_token": accessToken,
		"type":         "jsapi",
	}, "", nil, &result)
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

func (this *WxSdk) GetSessionByMiniProgram(code string) (WxMiniProgramSession, error) {
	result := WxMiniProgramSession{}
	err := this.apiJson("GET", "/sns/jscode2session", map[string]string{
		"appid":      this.AppId,
		"secret":     this.AppSecret,
		"js_code":    code,
		"grant_type": "authorization_code",
	}, "", nil, &result)
	if err != nil {
		return WxMiniProgramSession{}, err
	}
	return result, nil
}

// 小程序解密用户信息
func (this *WxSdk) DecryptByMiniProgram(sessionKey, encryptedData, iv string) (*WxMiniProgramUserInfo, error) {
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, errors.New("sessionKey:" + sessionKey + ",err:" + err.Error())
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, errors.New("encryptedData:" + sessionKey + ",err:" + err.Error())
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, errors.New("iv:" + sessionKey + ",err:" + err.Error())
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, err = pkcs7Unpad(cipherText, block.BlockSize())
	if err != nil {
		return nil, err
	}
	var userInfo WxMiniProgramUserInfo
	err = json.Unmarshal(cipherText, &userInfo)
	if err != nil {
		return nil, err
	}
	if userInfo.Watermark.AppID != this.AppId {
		return nil, errors.New("app id not match")
	}
	return &userInfo, nil
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

// pkcs7Unpad returns slice of the original data without padding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, errors.New("invalid block size")
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, errors.New("invalid PKCS7 data")
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, errors.New("invalid padding on input")
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, errors.New("invalid padding on input")
		}
	}
	return data[:len(data)-n], nil
}
