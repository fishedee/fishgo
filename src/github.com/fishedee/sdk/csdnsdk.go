package sdk

import (
	"errors"
	"fmt"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/util"
	"strings"
)

type CsdnSdk struct {
	AppKey    string
	AppSecert string
}

type CsdnBaseResponse struct {
	Request   string `json:"request"`
	ErrorCode string `json:"error_code"`
	Error     string `json:"error"`
}

type CsdnSdkAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserName    string `json:"username"`
}

type CsdnSdkBaseArticle struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

type CsdnSdkArticle struct {
	CsdnSdkBaseArticle
	Title          string `json:"title"`
	CreateAt       string `json:"create_at"`
	ViewCount      int    `json:"view_count"`
	CommentCount   int    `json:"comment_count"`
	CommentAllowed int    `json:"comment_allowed"`
	Type           string `json:"type"`
	Channel        int    `json:"channel"`
	Digg           int    `json:"digg"`
	Bury           int    `json:"bury"`
	Description    string `json:"description"`
}

type CsdnSdkDetailArticle struct {
	CsdnSdkArticle
	Categories string `json:"categories"`
	Tags       string `json:"tags"`
	Content    string `json:"content"`
}

type CsdnSdkCategory struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Hide         bool   `json:"hide"`
	ArticleCount int    `json:"article_count"`
}

type CsdnSdkArticleList struct {
	Page  int              `json:"page"`
	Count int              `json:"count"`
	Size  int              `json:"size"`
	List  []CsdnSdkArticle `json:"list"`
}

type CsdnSdkGetArticleListRequest struct {
	AccessToken string `url:"access_token"`
	Status      string `url:"status,omitempty"`
	Page        int    `url:"page,omitempty"`
	Size        int    `url:"size,omitempty"`
}

type CsdnSdkGetArticleRequest struct {
	AccessToken string `url:"access_token"`
	Id          int    `url:"id"`
}

type CsdnSdkGetCategoryListRequest struct {
	AccessToken string `url:"access_token"`
}

type CsdnSdkSaveArticleRequest struct {
	AccessToken string `url:"access_token"`
	Id          int    `url:"id,omitempty"`
	Title       string `url:"title"`
	Type        string `url:"type,omitempty"`
	Description string `url:"description,omitempty"`
	Content     string `url:"content"`
	Categories  string `url:"categories,omitempty"`
	Tags        string `url:"tags,omitempty"`
	Ip          string `url:"ip,omitempty"`
}

func (this *CsdnSdk) GetAuthUrl(redirectUrl string) (string, error) {
	appKeyEncode, err := EncodeUrl(this.AppKey)
	if err != nil {
		return "", err
	}
	redirectUrlEncode, err := EncodeUrl(redirectUrl)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"http://api.csdn.net/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code",
		appKeyEncode,
		redirectUrlEncode,
	), nil
}

func (this *CsdnSdk) ajax(url string, method string, data interface{}, responseData interface{}) error {
	//请求网络
	var responseDataByte []byte
	option := Ajax{
		Url:          url,
		Data:         data,
		ResponseData: &responseDataByte,
	}
	var err error
	if strings.ToLower(method) == "get" {
		err = DefaultAjaxPool.Get(&option)
	} else {
		err = DefaultAjaxPool.Post(&option)
	}
	if err != nil {
		return err
	}
	fmt.Println(string(responseDataByte))

	//判断基础失败
	var baseResponse CsdnBaseResponse
	err = DecodeJson(responseDataByte, &baseResponse)
	if err == nil && baseResponse.ErrorCode != "" {
		return errors.New(fmt.Sprintf("调用失败，错误码：%v,错误描述:%v", baseResponse.ErrorCode, baseResponse.Error))
	}

	//转换数据
	err = DecodeJson(responseDataByte, &responseData)
	if err != nil {
		return err
	}
	return nil
}

func (this *CsdnSdk) GetAccessToken(redirectUrl string, code string) (CsdnSdkAccessToken, error) {
	var result CsdnSdkAccessToken
	err := this.ajax(
		"http://api.csdn.net/oauth2/access_token",
		"get",
		map[string]interface{}{
			"client_id":     this.AppKey,
			"client_secret": this.AppSecert,
			"grant_type":    "authorization_code",
			"redirect_uri":  redirectUrl,
			"code":          code,
		},
		&result,
	)
	return result, err
}

func (this *CsdnSdk) GetArticleList(request CsdnSdkGetArticleListRequest) (CsdnSdkArticleList, error) {
	var result CsdnSdkArticleList
	err := this.ajax(
		"http://api.csdn.net/blog/getarticlelist",
		"get",
		request,
		&result,
	)
	return result, err
}

func (this *CsdnSdk) GetArticle(request CsdnSdkGetArticleRequest) (CsdnSdkDetailArticle, error) {
	var result CsdnSdkDetailArticle
	err := this.ajax(
		"http://api.csdn.net/blog/getarticle",
		"get",
		request,
		&result,
	)
	return result, err
}

func (this *CsdnSdk) GetCategoryList(request CsdnSdkGetCategoryListRequest) ([]CsdnSdkCategory, error) {
	var result []CsdnSdkCategory
	err := this.ajax(
		"http://api.csdn.net/blog/getcategorylist",
		"get",
		request,
		&result,
	)
	return result, err
}

func (this *CsdnSdk) SaveArticle(request CsdnSdkSaveArticleRequest) (CsdnSdkBaseArticle, error) {
	var result CsdnSdkBaseArticle
	err := this.ajax(
		"http://api.csdn.net/blog/savearticle",
		"post",
		request,
		&result,
	)
	return result, err
}
