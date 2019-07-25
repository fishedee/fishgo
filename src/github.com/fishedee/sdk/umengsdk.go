package sdk

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	. "github.com/fishedee/encoding"
	. "github.com/fishedee/util"
)

type UmengSdk struct {
	AccessKey string
	SecretKey string
}

type AnaymonusMap map[string]string

const (
	sendUrl   = "http://msg.umeng.com/api/send"
	statusUrl = "http://msg.umeng.com/api/status"
	uploadUrl = "http://msg.umeng.com/upload"
)

type UmengCommon struct {
	Appkey         string      `json:"appkey"`
	Timestamp      string      `json:"timestamp"`
	Type           string      `json:"type"`
	DeviceTokens   string      `json:"device_tokens"`
	AliasType      string      `json:"alias_type"`
	Alias          string      `json:"alias"`
	FileId         string      `json:"file_id"`
	Filter         interface{} `json:"filter"`
	ProductionMode string      `json:"production_mode"`
	Description    string      `json:"description"`
	ThirdpartyId   string      `json:"thirdparty_id"`
}

type UmengAndroidPayload struct {
	DisplayType string                  `json:"display_type"`
	Body        UmengAndroidPayloadBody `json:"body"`
	Extra       map[string]string       `json:"extra"`
}

type UmengAndroidPayloadBody struct {
	Ticker      string `json:"ticker"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	Icon        string `json:"icon"`
	LargeIcon   string `json:"largeIcon"`
	Img         string `json:"img"`
	Sound       string `json:"sound"`
	BuilderId   string `json:"builder_id"`
	PlayVibrate bool   `json:"play_vibrate"`
	PlayLights  bool   `json:"play_lights"`
	PlaySound   bool   `json:"play_sound"`
	AfterOpen   string `json:"after_open"`
	Url         string `json:"url"`
	Activity    string `json:"activity"`
	Custom      string `json:"custom"`
}

type UmengAndroidPolicy struct {
	StartTime  string `json:"start_time"`
	ExpireTime string `json:"expire_time"`
	//Max_send_num int    `json:"max_send_num"`
	OutBizNo string `json:"out_biz_no"`
}

type UmengAndroid struct {
	UmengCommon
	Payload		UmengAndroidPayload `json:"payload"`
	Policy		UmengAndroidPolicy  `json:"policy"`
	Mipush		string				`json:"mipush"`
	MiActivity	string				`json:"mi_activity"`
}

type UmengIOSPayload struct {
	Aps UmengIOSPayloadAps `json:"aps"`
	AnaymonusMap
}

type UmengIOSPayloadAps struct {
	Alert            UmengIOSPayloadApsAlert `json:"alert"`
	Badge            int                     `json:"badge"`
	Sound            string                  `json:"sound"`
	ContentAvailable string                  `json:"content-available"`
	Category         string                  `json:"category"`
}

type UmengIOSPayloadApsAlert struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Body     string `json:"body"`
}

type UmengIOSPolicy struct {
	StartTime  string `json:"start_time"`
	ExpireTime string `json:"expire_time"`
	//Max_send_num int    `json:"max_send_num"`
}

type UmengIOS struct {
	UmengCommon
	Payload UmengIOSPayload `json:"payload"`
	Policy  UmengIOSPolicy  `json:"policy"`
}

type UmengResult struct {
	Ret  string `json:"ret"`
	Data struct {
		MsgId        string `json:"msg_id"`
		TaskId       string `json:"task_id"`
		ErrorCode    string `json:"error_code"`
		ThirdpartyId string `json:"thirdparty_id"`
	}
}

type UmengStatus struct {
	Appkey    string `json:"appkey"`
	Timestamp string `json:"timestamp"`
	TaskId    string `json:"task_id"`
}

type UmengStatusResult struct {
	Ret  string `json:"ret"`
	Data struct {
		TaskId string `json:"task_id"`
		Status int    `json:"status"` // 消息状态: 0-排队中, 1-发送中，2-发送完成，3-发送失败，4-消息被撤销，
		// 5-消息过期, 6-筛选结果为空，7-定时任务尚未开始处理
		TotalCount   int `json:"total_count"`   // 消息总数
		AcceptCount  int `json:"accept_count"`  // 消息受理数
		SentCount    int `json:"sent_count"`    // 消息实际发送数
		OpenCount    int `json:"open_count"`    //打开数
		DismissCount int `json:"dismiss_count"` //忽略数

		ErrorCode string `json:"error_code"`
	}
}

type UmengFile struct {
	Appkey    string `json:"appkey"`
	Timestamp string `json:"timestamp"`
	Content   string `json:"content"`
}

type UmengFileResult struct {
	Ret  string `json:"ret"`
	Data struct {
		FileId string `json:"file_id"`
	}
}

func (this *UmengSdk) SendAndroid(umengAndroid UmengAndroid) (UmengResult, error) {
	sign := ""
	method := "POST"

	body, err := EncodeJson(umengAndroid)
	if err != nil {
		return UmengResult{}, err
	}
	sign = this.getSign(method, sendUrl, string(body))
	url := sendUrl + "?sign=" + sign

	var result []byte
	err = DefaultAjaxPool.Post(&Ajax{
		Url:          url,
		Data:         body,
		ResponseData: &result,
	})
	if err != nil {
		if _, ok := err.(*AjaxStatusCodeError); !ok {
			return UmengResult{}, err
		}
	}

	var finalResult UmengResult
	err = DecodeJson(result, &finalResult)
	if err != nil {
		return UmengResult{}, err
	}
	return finalResult, nil
}

func (this *UmengSdk) SendIOS(umengIOS UmengIOS) (UmengResult, error) {
	sign := ""
	method := "POST"

	body, err := EncodeJson(umengIOS)
	if err != nil {
		return UmengResult{}, err
	}
	sign = this.getSign(method, sendUrl, string(body))
	url := sendUrl + "?sign=" + sign

	var result []byte
	err = DefaultAjaxPool.Post(&Ajax{
		Url:          url,
		Data:         body,
		ResponseData: &result,
	})
	if err != nil {
		if _, ok := err.(*AjaxStatusCodeError); !ok {
			return UmengResult{}, err
		}
	}

	var finalResult UmengResult
	err = json.Unmarshal(result, &finalResult)
	if err != nil {
		return UmengResult{}, err
	}
	return finalResult, nil
}

func (this *UmengSdk) GetFileId(deviceToken string) (UmengFileResult, error) {
	sign := ""
	method := "POST"

	body, err := json.Marshal(UmengFile{
		Appkey:    this.AccessKey,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		Content:   deviceToken,
	})
	if err != nil {
		return UmengFileResult{}, err
	}
	sign = this.getSign(method, uploadUrl, string(body))
	url := uploadUrl + "?sign=" + sign

	var result []byte
	err = DefaultAjaxPool.Post(&Ajax{
		Url:          url,
		Data:         body,
		ResponseData: &result,
	})
	if err != nil {
		return UmengFileResult{}, err
	}

	var finalResult UmengFileResult
	err = json.Unmarshal(result, &finalResult)
	if err != nil {
		return UmengFileResult{}, err
	}
	return finalResult, nil
}

func (this *UmengSdk) GetStatus(taskId string) (UmengStatusResult, error) {
	sign := ""
	method := "POST"

	body, err := json.Marshal(UmengStatus{
		Appkey:    this.AccessKey,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		TaskId:    taskId,
	})
	if err != nil {
		return UmengStatusResult{}, err
	}
	sign = this.getSign(method, statusUrl, string(body))
	url := statusUrl + "?sign=" + sign

	var result []byte
	err = DefaultAjaxPool.Post(&Ajax{
		Url:          url,
		Data:         body,
		ResponseData: &result,
	})
	if err != nil {
		return UmengStatusResult{}, err
	}

	var finalResult UmengStatusResult
	err = json.Unmarshal(result, &finalResult)
	if err != nil {
		return UmengStatusResult{}, err
	}
	return finalResult, nil
}

func (this *UmengSdk) getSign(method, url, body string) string {
	signStr := strings.ToUpper(method) + url + body + this.SecretKey
	return fmt.Sprintf("%x", md5.Sum([]byte(signStr)))
}
