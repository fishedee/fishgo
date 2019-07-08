package sdk

import (
	"fmt"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/util"
)

type WxQySdk struct {
	CorpId      string
	AgentId     int
	Secert      string
	AccessToken string
}

type WxQySdkToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WxQySdkUploadMedia struct {
	Type string
	Name string
	Data []byte
}

type WxQySdkUploadMediaResult struct {
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

type WxQySdkSendTextMessage struct {
	Content string `json:"content,omitempty"`
}

type WxQySdkSendImageMessage struct {
	MediaId string `json:"media_id,omitempty"`
}

type WxQySdkSendNewsArticleMessage struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	PicUrl      string `json:"picurl,omitempty"`
}

type WxQySdkSendNewsMessage struct {
	Articles []WxQySdkSendNewsArticleMessage `json:"articles,omitempty"`
}

type WxQySdkSendMessage struct {
	ToUser  string `json:"touser,omitempty"`
	ToParty string `json:"toparty,omitempty"`
	ToTag   int    `json:"totag,omitempty"`
	MsgType string `json:"msgtype"`
	Safe    int    `json:"safe,omitempty"`
	AgentId int    `json:"agentid"`
	//文本消息
	Text WxQySdkSendTextMessage `json:"text,omitempty"`
	//图片消息
	Image WxQySdkSendImageMessage `json:"image,omitempty"`
	//图文消息
	News WxQySdkSendNewsMessage `json:"news,omitempty"`
}

type WxQySdkSendMessageResult struct {
}

type WxQySdkError struct {
	Code    int
	Message string
}

func (this *WxQySdkError) GetCode() int {
	return this.Code
}

func (this *WxQySdkError) GetMsg() string {
	return this.Message
}

func (this *WxQySdkError) Error() string {
	return fmt.Sprintf("错误码为：%v，错误描述为：%v", this.Code, this.Message)
}

func (this *WxQySdk) api(method string, url string, query interface{}, dataType string, data interface{}) ([]byte, error) {
	queryInfo, err := EncodeUrlQuery(query)
	if err != nil {
		return nil, err
	}
	url = "https://qyapi.weixin.qq.com" + url
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

func (this *WxQySdk) apiJson(method string, url string, query interface{}, dataType string, data interface{}, responseData interface{}) error {
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
func (this *WxQySdk) GetAccessToken() (WxQySdkToken, error) {
	result := WxQySdkToken{}
	err := this.apiJson("GET", "/cgi-bin/gettoken", map[string]string{
		"corpid":     this.CorpId,
		"corpsecret": this.Secert,
	}, "", nil, &result)
	if err != nil {
		return WxQySdkToken{}, err
	}
	return result, nil
}

//上传临时素材
func (this *WxQySdk) UploadMedia(media WxQySdkUploadMedia) (WxQySdkUploadMediaResult, error) {
	var result WxQySdkUploadMediaResult

	data, err := this.api("POST", "/cgi-bin/media/upload", map[string]string{
		"access_token": this.AccessToken,
		"type":         media.Type,
	}, "form", map[string]interface{}{
		"media": []interface{}{media.Name, media.Data},
	})
	if err != nil {
		return result, err
	}

	err = DecodeJson(data, &result)
	if err != nil {
		return WxQySdkUploadMediaResult{}, err
	}
	return result, nil
}

//发送应用消息
func (this *WxQySdk) SendMessage(message WxQySdkSendMessage) (WxQySdkSendMessageResult, error) {
	var result WxQySdkSendMessageResult

	message.AgentId = this.AgentId
	err := this.apiJson("POST", "/cgi-bin/message/send", map[string]string{
		"access_token": this.AccessToken,
	}, "", message, &result)
	if err != nil {
		return WxQySdkSendMessageResult{}, err
	}
	return result, nil
}
