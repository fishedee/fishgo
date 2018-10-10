package sdk

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	. "github.com/fishedee/encoding"
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
	FuncScopeCategory struct {
		Id int `json:"id"`
	} `json:"funcscope_category"`
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
	ServiceTypeInfo struct {
		Id int `json:"id"`
	} `json:"service_type_info"`
	VerifyTypeInfo struct {
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

//解密Message
func (this *WxPluginSdk) DecryptMessage(msgSignature string, timestamp string, nonce string, msg []byte) (WxPluginSdkMessage, error) {
	result := WxPluginSdkMessage{}
	wxCryptSdk := &WxCryptSdk{
		AppId:  this.ComponentAppId,
		Token:  this.MessageToken,
		AESKey: this.MessageAESKey,
	}
	_, packaget, err := wxCryptSdk.Decrypt(msgSignature, timestamp, nonce, msg)
	if err != nil {
		return WxPluginSdkMessage{}, err
	}
	err = xml.Unmarshal(packaget, &result)
	if err != nil {
		return WxPluginSdkMessage{}, err
	}
	return result, nil
}

//加密Message
func (this *WxPluginSdk) EncryptMessage(msg []byte) ([]byte, error) {
	wxCryptSdk := &WxCryptSdk{
		AppId:  this.ComponentAppId,
		Token:  this.MessageToken,
		AESKey: this.MessageAESKey,
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := wxCryptSdk.getRandomStr(32)
	return wxCryptSdk.Encrypt(timestamp, string(nonce), msg)
}
