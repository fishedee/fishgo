package sdk

import (
	"fmt"

	. "github.com/fishedee/encoding"
	. "github.com/fishedee/util"
)

type QqSdk struct {
	AppId     string
	AppSecret string
}

type QqSdkOauthAccessToken struct {
	AccessToken string `url:"access_token"`
	ExpiresIn   string `url:"expires_in"`
}

type QqSdkOauthOpenId struct {
	ClientId string `jsonp:"client_id"`
	OpenId   string `jsonp:"openid"`
}

type QqSdkOauthUserInfo struct {
	NickName        string `json:"nickname"`
	Figureurl       string `json:"figureurl"`
	Figureurl1      string `json:"figureurl_1"`
	Figureurl2      string `json:"figureurl_2"`
	FigureurlQq1    string `json:"figureurl_qq_1"`
	FigureurlQq2    string `json:"figureurl_qq_2"`
	Gender          string `json:"gender"`
	IsYellowVip     string `json:"is_yellow_vip"`
	Vip             string `json:"vip"`
	YellowVipLevel  string `json:"yellow_vip_level"`
	Level           string `json:"level"`
	IsYellowYearVip string `json:"is_yellow_year_vip"`
	Year            string `json:"year"`
	Province        string `json:"province"`
	City            string `json:"city"`
}

type QqSdkError struct {
	Code    int    `url:"error" jsonp:"error" json:"ret"`
	Message string `url:"error_description" jsonp:"error_description" json:"msg"`
}

func (this *QqSdkError) GetCode() int {
	return this.Code
}

func (this *QqSdkError) GetMsg() string {
	return this.Message
}

func (this *QqSdkError) Error() string {
	return fmt.Sprintf("错误码为：%v，错误描述为：%v", this.Code, this.Message)
}

func (this *QqSdk) api(method string, url string, query interface{}) ([]byte, error) {
	queryInfo, err := EncodeUrlQuery(query)
	if err != nil {
		return nil, err
	}
	if len(queryInfo) != 0 {
		url += "?" + string(queryInfo)
	}
	var result []byte
	ajaxOption := &Ajax{
		Url:          url,
		ResponseData: &result,
	}
	if method == "GET" {
		err = DefaultAjaxPool.Get(ajaxOption)
	} else {
		err = DefaultAjaxPool.Post(ajaxOption)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (this *QqSdk) apiUrl(method string, url string, query interface{}, responseData interface{}) error {
	result, err := this.api(method, url, query)
	if err != nil {
		return err
	}
	var sdkErr QqSdkError
	_, err = DecodeJsonp(result, &sdkErr)
	if err == nil && sdkErr.Code != 0 {
		return &sdkErr
	}
	err = DecodeUrlQuery(result, responseData)
	if err != nil {
		return err
	}
	return nil
}

func (this *QqSdk) apiJsonp(method string, url string, query interface{}, responseData interface{}) error {
	result, err := this.api(method, url, query)
	if err != nil {
		return err
	}

	var sdkErr QqSdkError
	_, err = DecodeJsonp(result, &sdkErr)
	if err == nil && sdkErr.Code != 0 {
		return &sdkErr
	}
	_, err = DecodeJsonp(result, responseData)
	if err != nil {
		return err
	}
	return nil
}

func (this *QqSdk) apiJson(method string, url string, query interface{}, responseData interface{}) error {
	result, err := this.api(method, url, query)
	if err != nil {
		return err
	}

	var sdkErr QqSdkError
	err = DecodeJson(result, &sdkErr)
	if err == nil && sdkErr.Code != 0 {
		return &sdkErr
	}
	err = DecodeJson(result, responseData)
	if err != nil {
		return err
	}
	return nil
}

//Oauth
func (this *QqSdk) GetOauthUrl(callback string, state string, scope string) (string, error) {
	query := map[string]string{
		"client_id":     this.AppId,
		"redirect_uri":  callback,
		"response_type": "code",
		"scope":         scope,
		"state":         state,
	}
	queryStr, err := EncodeUrlQuery(query)
	if err != nil {
		return "", err
	}
	return "https://graph.qq.com/oauth2.0/authorize?" + string(queryStr), nil
}

func (this *QqSdk) GetOauthAccessToken(callback string, code string) (QqSdkOauthAccessToken, error) {
	var result QqSdkOauthAccessToken
	err := this.apiUrl("GET", "https://graph.qq.com/oauth2.0/token", map[string]interface{}{
		"grant_type":    "authorization_code",
		"client_id":     this.AppId,
		"redirect_uri":  callback,
		"client_secret": this.AppSecret,
		"code":          code,
	}, &result)
	if err != nil {
		return QqSdkOauthAccessToken{}, err
	}
	return result, nil
}

func (this *QqSdk) GetOauthOpenId(accessToken string) (QqSdkOauthOpenId, error) {
	var result QqSdkOauthOpenId
	err := this.apiJsonp("GET", "https://graph.qq.com/oauth2.0/me", map[string]interface{}{
		"access_token": accessToken,
	}, &result)
	if err != nil {
		return QqSdkOauthOpenId{}, err
	}
	return result, nil
}

func (this *QqSdk) GetOauthUserInfo(accessToken string, openId string) (QqSdkOauthUserInfo, error) {
	var result QqSdkOauthUserInfo
	err := this.apiJson("GET", "https://graph.qq.com/user/get_user_info", map[string]interface{}{
		"access_token":       accessToken,
		"openid":             openId,
		"oauth_consumer_key": this.AppId,
		"format":             "json",
	}, &result)
	if err != nil {
		return QqSdkOauthUserInfo{}, err
	}
	return result, nil
}
