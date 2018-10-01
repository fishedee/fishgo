package sdk

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
)

type WxPluginSdk struct {
	ComponentAppId       string
	ComponentAppSecret   string
	ComponentAccessToken string
	MessageToken         string
	MessageAESKey        string
}

type WxPluginSdkAccessToken struct {
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int    `json:"expires_in"`
}

type WxPluginSdkPreAuthCode struct {
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int    `json:"expires_in"`
}

type WxPluginSdkSingleFuncInfo struct {
	FuncScopeCategory string `json:"funcscope_category"`
	Id                string `json:"id"`
}

type WxPluginSdkAuthorizationInfoDetail struct {
	AuthorizerAppId        string                      `json:"authorizer_appid"`
	AuthorizerAccessToken  string                      `json:"authorizer_access_token"`
	ExpiresIn              int                         `json:"expires_in"`
	AuthorizerRefreshToken string                      `json:"authorizer_refresh_token"`
	FuncInfo               []WxPluginSdkSingleFuncInfo `json:"func_info"`
}

type WxPluginSdkAuthorizationInfo struct {
	AuthorizationInfo WxPluginSdkAuthorizationInfoDetail `json:"authorization_info"`
}

type WxPluginSdkAuthorizerAccessToken struct {
	AuthorizerAccessToken  string `json:"authorizer_access_token"`
	ExpiresIn              int    `json:"expires_in"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

type WxPluginSdkAuthorizerBusinessInfo struct {
	OpenStore int `json:"open_store"`
	OpenScan  int `json:"open_scan"`
	OpenPay   int `json:"open_pay"`
	OpenCard  int `json:"open_card"`
	OpenShake int `json:"open_shake"`
}

type WxPluginSdkAuthorizerInfoDetail struct {
	NickName        string `json:"nick_name"`
	HeadImg         string `json:"head_img"`
	ServiceTypeInfo []struct {
		Id int `json:"id"`
	} `json:"service_type_info"`
	VerifyTypeInfo []struct {
		Id int `json:"id"`
	} `json:"verify_type_info"`
	UserName      string                            `json:"user_name"`
	PrincipalName string                            `json:"principal_name"`
	BusinessInfo  WxPluginSdkAuthorizerBusinessInfo `json:"business_info"`
	Alias         string                            `json:"alias"`
	QrcodeUrl     string                            `json:"qrcode_url"`
}

type WxPluginSdkAuthorizerInfo struct {
	AuthorizerInfo    WxPluginSdkAuthorizerInfoDetail    `json:"authorizer_info"`
	AuthorizationInfo WxPluginSdkAuthorizationInfoDetail `json:"authorization_info"`
}

type WxPluginSdkGetAuthorizerOptionRequest struct {
	AuthorizerAppId string `json:"authorizer_appid"`
	OptionName      string `json:"option_name"`
}

type WxPluginSdkGetAuthorizerOptionResponse struct {
	AuthorizerAppId string `json:"authorizer_appid"`
	OptionName      string `json:"option_name"`
	OptionValue     int    `json:"option_value"`
}

type WxPluginSdkSetAuthorizerOptionRequest struct {
	AuthorizerAppId string `json:"authorizer_appid"`
	OptionName      string `json:"option_name"`
	OptionValue     string `json:"option_value"`
}

type WxPluginSdkSetAuthorizerOptionResponse struct {
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

type WxPluginSdkGetAuthorizerListRequest struct {
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

type WxPluginSdkGetAuthorizerListResponseAuthorizer struct {
	AuthorizerAppId string `json:"authorizer_appid"`
	RefreshToken    string `json:"refresh_token"`
	AuthTime        string `json:"auth_time"`
}

type WxPluginSdkGetAuthorizerListResponse struct {
	TotalCount int                                              `json:"total_count"`
	List       []WxPluginSdkGetAuthorizerListResponseAuthorizer `json:"list"`
}

type WxPluginSdkMessage struct {
	ToUserName                   string
	AppId                        string `xml:"AppId"`
	CreateTime                   int    `xml:"CreateTime"`
	InfoType                     string `xml:"InfoType"`
	ComponentVerifyTicket        string `xml:"ComponentVerifyTicket"`
	AuthorizerAppId              string `xml:"AuthorizerAppid"`
	AuthorizationCode            string `xml:"AuthorizationCode"`
	AuthorizationCodeExpiredTime int    `xml:"AuthorizationCodeExpiredTime"`
	PreAuthCode                  string `xml:"PreAuthCode"`
}

type WxPluginSdkError struct {
	Code    int
	Message string
}

func (this *WxPluginSdkError) GetCode() int {
	return this.Code
}

func (this *WxPluginSdkError) GetMsg() string {
	return this.Message
}

func (this *WxPluginSdkError) Error() string {
	return fmt.Sprintf("错误码为：%v，错误描述为：%v", this.Code, this.Message)
}

func (this *WxPluginSdk) api(method string, url string, query interface{}, dataType string, data interface{}) ([]byte, error) {
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
		return nil, &WxPluginSdkError{errInfo.ErrorCode, errInfo.ErrorMsg}
	}

	return result, nil
}

func (this *WxPluginSdk) apiJson(method string, url string, query interface{}, dataType string, data interface{}, responseData interface{}) error {

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

//自身调用凭证管理
func (this *WxPluginSdk) GetToken(componentVertifyToken string) (WxPluginSdkAccessToken, error) {
	result := WxPluginSdkAccessToken{}
	err := this.apiJson(
		"POST",
		"/cgi-bin/component/api_component_token",
		map[string]string{},
		"",
		map[string]string{
			"component_appid":         this.ComponentAppId,
			"component_appsecret":     this.ComponentAppSecret,
			"component_verify_ticket": componentVertifyToken,
		},
		&result)
	if err != nil {
		return WxPluginSdkAccessToken{}, err
	}
	return result, nil
}

//授权接口
func (this *WxPluginSdk) GetPreAuthCode() (WxPluginSdkPreAuthCode, error) {
	result := WxPluginSdkPreAuthCode{}
	err := this.apiJson(
		"POST",
		"/cgi-bin/component/api_create_preauthcode",
		map[string]string{
			"component_access_token": this.ComponentAccessToken,
		},
		"",
		map[string]string{
			"component_appid": this.ComponentAppId,
		},
		&result)
	if err != nil {
		return WxPluginSdkPreAuthCode{}, err
	}
	return result, nil
}

func (this *WxPluginSdk) GetPcAuthUrl(preAuthCode string, redirectUrl string, authType string, bizAppId string) (string, error) {
	query := map[string]string{
		"component_appid": this.ComponentAppId,
		"pre_auth_code":   preAuthCode,
		"redirect_uri":    redirectUrl,
	}
	if authType != "" {
		query["auth_type"] = authType
	}
	if bizAppId != "" {
		query["biz_appid"] = bizAppId
	}
	queryString, err := EncodeUrlQuery(query)
	if err != nil {
		return "", err
	}
	return "https://mp.weixin.qq.com/cgi-bin/componentloginpage?" + string(queryString), nil
}

func (this *WxPluginSdk) GetMobileAuthUrl(preAuthCode string, redirectUrl string, authType string, bizAppId string) (string, error) {
	query := map[string]string{
		"component_appid": this.ComponentAppId,
		"pre_auth_code":   preAuthCode,
		"redirect_uri":    redirectUrl,
		"no_scan":         "1",
		"action":          "bindcomponent",
	}
	if authType != "" {
		query["auth_type"] = authType
	}
	if bizAppId != "" {
		query["biz_appid"] = bizAppId
	}
	queryString, err := EncodeUrlQuery(query)
	if err != nil {
		return "", err
	}
	return "https://mp.weixin.qq.com/safe/bindcomponent?" + string(queryString) + "#wechat_redirect", nil
}

//获取授权的结果信息
func (this *WxPluginSdk) GetAuthorizationInfo(authorizationCode string) (WxPluginSdkAuthorizationInfo, error) {
	result := WxPluginSdkAuthorizationInfo{}
	err := this.apiJson(
		"POST",
		"/cgi-bin/component/api_query_auth",
		map[string]string{
			"component_access_token": this.ComponentAccessToken,
		},
		"",
		map[string]string{
			"component_appid":    this.ComponentAppId,
			"authorization_code": authorizationCode,
		},
		&result)
	if err != nil {
		return WxPluginSdkAuthorizationInfo{}, err
	}
	return result, nil
}

//外部公众号的调用凭证管理
func (this *WxPluginSdk) GetAuthorizerAccessToken(authorizerAppId string, authorizerRefreshToken string) (WxPluginSdkAuthorizerAccessToken, error) {
	result := WxPluginSdkAuthorizerAccessToken{}
	err := this.apiJson(
		"POST",
		"/cgi-bin/component/api_authorizer_token",
		map[string]string{
			"component_access_token": this.ComponentAccessToken,
		},
		"",
		map[string]string{
			"component_appid":          this.ComponentAppId,
			"authorizer_appid":         authorizerAppId,
			"authorizer_refresh_token": authorizerRefreshToken,
		},
		&result)
	if err != nil {
		return WxPluginSdkAuthorizerAccessToken{}, err
	}
	return result, nil
}

//获取外部公众号的信息
func (this *WxPluginSdk) GetAuthorizerInfo(authorizerAppId string) (WxPluginSdkAuthorizerInfo, error) {
	result := WxPluginSdkAuthorizerInfo{}
	err := this.apiJson(
		"POST",
		"/cgi-bin/component/api_get_authorizer_info",
		map[string]string{
			"component_access_token": this.ComponentAccessToken,
		},
		"",
		map[string]string{
			"component_appid":  this.ComponentAppId,
			"authorizer_appid": authorizerAppId,
		},
		&result)
	if err != nil {
		return WxPluginSdkAuthorizerInfo{}, err
	}
	return result, nil
}

//获取外部公众号的选项信息
func (this *WxPluginSdk) GetAuthorizerOption(request WxPluginSdkGetAuthorizerOptionRequest) (WxPluginSdkGetAuthorizerOptionResponse, error) {
	result := WxPluginSdkGetAuthorizerOptionResponse{}
	err := this.apiJson(
		"POST",
		"/cgi-bin/component/api_get_authorizer_option",
		map[string]string{
			"component_access_token": this.ComponentAccessToken,
		},
		"",
		map[string]string{
			"component_appid":  this.ComponentAppId,
			"authorizer_appid": request.AuthorizerAppId,
			"option_name":      request.OptionName,
		},
		&result)
	if err != nil {
		return WxPluginSdkGetAuthorizerOptionResponse{}, err
	}
	return result, nil
}

//设置外部公众号的选项信息
func (this *WxPluginSdk) SetAuthorizerOption(request WxPluginSdkSetAuthorizerOptionRequest) (WxPluginSdkSetAuthorizerOptionResponse, error) {
	result := WxPluginSdkSetAuthorizerOptionResponse{}
	err := this.apiJson(
		"POST",
		"/cgi-bin/component/api_set_authorizer_option",
		map[string]string{
			"component_access_token": this.ComponentAccessToken,
		},
		"",
		map[string]string{
			"component_appid":  this.ComponentAppId,
			"authorizer_appid": request.AuthorizerAppId,
			"option_name":      request.OptionName,
			"option_value":     request.OptionValue,
		},
		&result)
	if err != nil {
		return WxPluginSdkSetAuthorizerOptionResponse{}, err
	}
	return result, nil
}

//获取外部公众号的列表
func (this *WxPluginSdk) GetAuthorizerList(request WxPluginSdkGetAuthorizerListRequest) (WxPluginSdkGetAuthorizerListResponse, error) {
	result := WxPluginSdkGetAuthorizerListResponse{}
	err := this.apiJson(
		"POST",
		"/cgi-bin/component/api_get_authorizer_list",
		map[string]string{
			"component_access_token": this.ComponentAccessToken,
		},
		"",
		map[string]string{
			"component_appid": this.ComponentAppId,
			"offset":          strconv.Itoa(request.Offset),
			"count":           strconv.Itoa(request.Count),
		},
		&result)
	if err != nil {
		return WxPluginSdkGetAuthorizerListResponse{}, err
	}
	return result, nil
}

//编码和解码Message
func (this *WxPluginSdk) getSignature(token string, timestamp string, nonce string, msg string) string {
	arrayInfo := []string{token, timestamp, nonce, msg}
	arrayInfo = ArraySort(arrayInfo).([]string)
	arrayInfoString := Implode(arrayInfo, "")
	return this.encodeSha1(arrayInfoString)
}

func (this *WxPluginSdk) getRandomStr(length int) []byte {
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz")
	result := make([]byte, length, length)
	for i := 0; i < length; i++ {
		key := rand.Intn(len(chars))
		result[i] = chars[key]
	}
	return result
}

func (this *WxPluginSdk) encodeSha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func (this *WxPluginSdk) decodeXml(msg []byte, data interface{}) error {
	return xml.Unmarshal(msg, data)
}

func (this *WxPluginSdk) encodeXml(encrypt string, signature string, timestamp string, nonce string) ([]byte, error) {
	return []byte(fmt.Sprintf(`<xml>
		<Encrypt><![CDATA[%s]]></Encrypt>
		<MsgSignature><![CDATA[%s]]></MsgSignature>
		<TimeStamp>%s</TimeStamp>
		<Nonce><![CDATA[%s]]></Nonce>
		</xml>`, encrypt, signature, timestamp, nonce)), nil
}

func (this *WxPluginSdk) pkcs7Unpadding(data []byte, blockSize int) []byte {
	length := len(data)
	unPadding := int(data[length-1])
	return data[:(length - unPadding)]
}

func (this *WxPluginSdk) pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func (this *WxPluginSdk) decodeAES(AESKey string, msg string) ([]byte, error) {
	aesKey, err := base64.StdEncoding.DecodeString(AESKey + "=")
	if err != nil {
		return nil, err
	}
	cipherText, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return nil, err
	}
	iv := aesKey[0:16]
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText = this.pkcs7Unpadding(cipherText, block.BlockSize())
	return cipherText, nil
}

func (this *WxPluginSdk) encodeAES(AESKey string, msg []byte) (string, error) {
	aesKey, err := base64.StdEncoding.DecodeString(AESKey + "=")
	if err != nil {
		return "", err
	}
	iv := aesKey[0:16]
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	cipherText := this.pkcs7Padding([]byte(msg), block.BlockSize())
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)
	cipherTextEncode := base64.StdEncoding.EncodeToString(cipherText)
	return cipherTextEncode, nil
}

func (this *WxPluginSdk) decodeMeta(packaget []byte) ([]byte, string) {
	//头四位随机字符串
	packaget = packaget[16:]
	//长度标记
	msgLen := binary.BigEndian.Uint32(packaget[0:4])
	packaget = packaget[4:]
	//数据
	msg := packaget[0:msgLen]
	packaget = packaget[msgLen:]
	//appId
	appId := packaget
	return msg, string(appId)
}

func (this *WxPluginSdk) encodeMeta(packaget []byte, appId string) []byte {
	var buffer bytes.Buffer
	//头四位随机字符串
	buffer.Write(this.getRandomStr(16))
	//长度标记
	lengthBuffer := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(lengthBuffer, uint32(len(packaget)))
	buffer.Write(lengthBuffer)
	//数据
	buffer.Write(packaget)
	//appId
	buffer.WriteString(appId)
	return buffer.Bytes()
}

func (this *WxPluginSdk) decodePackaget(msgSignature string, timestamp string, nonce string, msg []byte) (string, []byte, error) {
	//解包外层xml
	var encryptMessage struct {
		ToUserName string `xml:"ToUserName"`
		Encrypt    string `xml:"Encrypt"`
	}
	err := this.decodeXml(msg, &encryptMessage)
	if err != nil {
		return "", nil, err
	}
	//检查签名
	realSignature := this.getSignature(
		this.MessageToken,
		timestamp,
		nonce,
		encryptMessage.Encrypt)
	if realSignature != msgSignature {
		return "", nil, errors.New("消息签名错误")
	}
	//解包内层xml
	packaget, err := this.decodeAES(
		this.MessageAESKey,
		encryptMessage.Encrypt,
	)
	if err != nil {
		return "", nil, err
	}
	packaget, appId := this.decodeMeta(packaget)
	if appId != this.ComponentAppId {
		return "", nil, errors.New("消息appId校验错误")
	}
	return encryptMessage.ToUserName, packaget, nil
}

func (this *WxPluginSdk) DecodeMessage(msgSignature string, timestamp string, nonce string, msg []byte) (WxPluginSdkMessage, error) {
	result := WxPluginSdkMessage{}
	toUserName, packaget, err := this.decodePackaget(msgSignature, timestamp, nonce, msg)
	if err != nil {
		return WxPluginSdkMessage{}, err
	}
	err = this.decodeXml(packaget, &result)
	if err != nil {
		return WxPluginSdkMessage{}, err
	}
	result.ToUserName = toUserName
	return result, nil
}

func (this *WxPluginSdk) encodePackaget(timestamp string, nonce string, msg []byte) ([]byte, error) {
	//打包内层xml
	msgWithMeta := this.encodeMeta(
		msg,
		this.ComponentAppId)
	encodeMsg, err := this.encodeAES(
		this.MessageAESKey,
		msgWithMeta)
	if err != nil {
		return nil, err
	}
	//生成签名
	signature := this.getSignature(
		this.MessageToken,
		timestamp,
		nonce,
		encodeMsg)
	//打包外层xml
	return this.encodeXml(encodeMsg, signature, timestamp, nonce)
}

func (this *WxPluginSdk) EncodeMessage(msg []byte) ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := this.getRandomStr(32)
	return this.encodePackaget(timestamp, string(nonce), msg)
}
