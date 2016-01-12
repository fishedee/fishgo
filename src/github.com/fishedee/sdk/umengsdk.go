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

type common struct {
	Appkey          string      `json:"appkey"`
	Timestamp       string      `json:"timestamp"`
	Type            string      `json:"type"`
	Device_tokens   string      `json:"device_tokens"`
	Alias_type      string      `json:"alias_type"`
	Alias           string      `json:"alias"`
	File_id         string      `json:"file_id"`
	Filter          interface{} `json:"filter"`
	Production_mode string      `json:"production_mode"`
	Description     string      `json:"description"`
	Thirdparty_id   string      `json:"thirdparty_id"`
}

type androidPayload struct {
	Display_type string             `json:"display_type"`
	Body         androidPayloadBody `json:"body"`
	Extra        map[string]string  `json:"extra"`
}

type androidPayloadBody struct {
	Ticker       string `json:"ticker"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	Icon         string `json:"icon"`
	LargeIcon    string `json:"largeIcon"`
	Img          string `json:"img"`
	Sound        string `json:"sound"`
	Builder_id   string `json:"builder_id"`
	Play_vibrate bool   `json:"play_vibrate"`
	Play_lights  bool   `json:"play_lights"`
	Play_sound   bool   `json:"play_sound"`
	After_open   string `json:"after_open"`
	Url          string `json:"url"`
	Activity     string `json:"activity"`
	Custom       string `json:"custom"`
}

type androidPolicy struct {
	Start_time   string `json:"start_time"`
	Expire_time  string `json:"expire_time"`
	Max_send_num int    `json:"max_send_num"`
	Out_biz_no   string `json:"out_biz_no"`
}

type android struct {
	common
	Payload androidPayload `json:"payload"`
	Policy  androidPolicy  `json:"policy"`
}

type iOSPayload struct {
	Aps        iOSPayloadAps `json:"aps"`
	After_open string        `json:"after_open"`
	Url        string        `json:"url"`
	Activity   string        `json:"activity"`
	Custom     string        `json:"custom"`
}

type iOSPayloadAps struct {
	Alert            string `json:"alert"`
	Badge            int    `json:"badge"`
	Sound            string `json:"sound"`
	ContentAvailable string `json:"content-available"`
	Category         string `json:"category"`
}

type iOSPolicy struct {
	Start_time   string `json:"start_time"`
	Expire_time  string `json:"expire_time"`
	Max_send_num int    `json:"max_send_num"`
}

type iOS struct {
	common
	Payload iOSPayload `json:"payload"`
	Policy  iOSPolicy  `json:"policy"`
}

type UmengResult struct {
	Ret  string
	Data struct {
		Msg_id        string
		Task_id       string
		Error_code    string
		Thirdparty_id string
	}
}

func (this *UmengSdk) SendAndroidCustom(alias, aliasType, ticker, title, text, afterOpen, extra, productionMode string) (UmengResult, error) {
	sign := ""
	method := "POST"
	url := "http://msg.umeng.com/api/send"
	params := android{
		common: common{
			Appkey:          this.AccessKey,
			Timestamp:       strconv.FormatInt(time.Now().Unix(), 10),
			Type:            "customizedcast",
			Alias_type:      aliasType,
			Alias:           alias,
			Production_mode: productionMode,
		},
		Payload: androidPayload{
			"notification",
			androidPayloadBody{
				Ticker:     ticker,
				Title:      title,
				Text:       text,
				After_open: "go_app",
			},
			map[string]string{
				"after_open": afterOpen,
				"url":        extra,
				"activity":   extra,
				"custom":     extra,
			},
		},
		Policy: androidPolicy{
			Max_send_num: 100,
		},
	}

	body, err := EncodeJson(params)
	if err != nil {
		return UmengResult{}, err
	}
	sign = this.getSign(method, url, string(body))
	url = url + "?sign=" + sign

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

func (this *UmengSdk) SendIOSCustom(appMsg, alias, aliasType, afterOpen, extra string, badge int, sound, productionMode string) (UmengResult, error) {
	sign := ""
	method := "POST"
	url := "http://msg.umeng.com/api/send"
	params := iOS{
		common: common{
			Appkey:          this.AccessKey,
			Timestamp:       strconv.FormatInt(time.Now().Unix(), 10),
			Type:            "customizedcast",
			Alias_type:      aliasType,
			Alias:           alias,
			Production_mode: productionMode,
		},
		Payload: iOSPayload{
			Aps: iOSPayloadAps{
				Alert: appMsg,
				Badge: badge,
				Sound: sound,
			},
			After_open: afterOpen,
			Url:        extra,
			Activity:   extra,
			Custom:     extra,
		},
		Policy: iOSPolicy{
			Max_send_num: 100,
		},
	}

	body, err := json.Marshal(params)
	if err != nil {
		return UmengResult{}, err
	}
	sign = this.getSign(method, url, string(body))
	url = url + "?sign=" + sign

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

func (this *UmengSdk) getSign(method, url, body string) string {
	signStr := strings.ToUpper(method) + url + body + this.SecretKey
	return fmt.Sprintf("%x", md5.Sum([]byte(signStr)))
}
