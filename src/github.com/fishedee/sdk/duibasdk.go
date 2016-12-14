package sdk

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	. "github.com/fishedee/crypto"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
)

type DuibaSdk struct {
	AppKey    string
	AppSecret string
}

type DuibaSdkReceiveCreditConsume struct {
	Uid         string `url:"uid"`
	Credits     int    `url:"credits"`
	AppKey      string `url:"appKey"`
	Timestamp   string `url:"timestamp"`
	Description string `url:"description"`
	OrderNum    string `url:"orderNum"`
	Type        string `url:"type"`
	FacePrice   int    `url:"facePrice"`
	ActualPrice int    `url:"actualPrice"`
	Ip          string `url:"ip"`
	WaitAudit   bool   `url:"waitAudit"`
	Params      string `url:"params"`
	Sign        string `url:"sign"`
}

type DuibaSdkSendCreditConsume struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
	BizId        string `json:"bizId"`
	Credits      string `json:"credits"`
}

type DuibaSdkReceiveCreditNotify struct {
	Uid          string `url:"uid"`
	AppKey       string `url:"appKey"`
	Timestamp    int    `url:"timestamp"`
	Success      bool   `url:"success"`
	ErrorMessage string `url:"errorMessage"`
	OrderNum     string `url:"orderNum"`
	BizId        string `url:"bizId"`
	Sign         string `url:"sign"`
}

func (this *DuibaSdk) getSign(data map[string]interface{}) string {
	dataKeysInterface, _ := ArrayKeyAndValue(data)
	dataKeys := dataKeysInterface.([]string)
	dataKeys = append(dataKeys, "appSecret")
	dataKeys = ArraySort(dataKeys).([]string)

	result := ""
	for _, singleDataKey := range dataKeys {
		if singleDataKey == "sign" {
			continue
		} else if singleDataKey == "appSecret" {
			result += this.AppSecret
		} else {
			result += fmt.Sprintf("%v", data[singleDataKey])
		}
	}
	return CryptoMd5([]byte(result))
}

func (this *DuibaSdk) checkSign(dataMap map[string]interface{}) error {
	if dataMap["appKey"] != this.AppKey {
		return errors.New(fmt.Sprintf("invalid appkey [%v != %v]", dataMap["appKey"], this.AppKey))
	}
	sign := this.getSign(dataMap)
	if sign != dataMap["sign"] {
		return errors.New(fmt.Sprintf("invalid sign [%v != %v]", dataMap["sign"], sign))
	}
	return nil
}

func (this *DuibaSdk) GetLoginUrl(clientId string, point int) (string, error) {
	query := map[string]interface{}{
		"uid":       clientId,
		"credits":   point,
		"appKey":    this.AppKey,
		"timestamp": int64(time.Now().UnixNano() / 1e6),
	}
	query["sign"] = this.getSign(query)
	queryStr, err := EncodeUrlQuery(query)
	if err != nil {
		return "", err
	}
	return this.getHttpsLoginHost() + string(queryStr), nil
}

func (this *DuibaSdk) getHttpLoginHost() string {
	return "http://www.duiba.com.cn/autoLogin/autologin?"
}

func (this *DuibaSdk) getHttpsLoginHost() string {
	return "https://www.duiba.com.cn/autoLogin/autologin?"
}

func (this *DuibaSdk) ReceiveCreditConsume(request *url.URL) (DuibaSdkReceiveCreditConsume, error) {
	var resultMap map[string]interface{}
	queryStr := request.RawQuery
	err := DecodeUrlQuery([]byte(queryStr), &resultMap)
	if err != nil {
		return DuibaSdkReceiveCreditConsume{}, err
	}
	err = this.checkSign(resultMap)
	if err != nil {
		return DuibaSdkReceiveCreditConsume{}, err
	}
	var result DuibaSdkReceiveCreditConsume
	err = MapToArray(resultMap, &result, "url")
	if err != nil {
		return DuibaSdkReceiveCreditConsume{}, err
	}
	return result, nil
}

func (this *DuibaSdk) SendCreditConsume(data DuibaSdkSendCreditConsume) ([]byte, error) {
	return EncodeJson(data)
}

func (this *DuibaSdk) ReceiveCreditNotify(request *url.URL) (DuibaSdkReceiveCreditNotify, error) {
	var resultMap map[string]interface{}
	queryStr := request.RawQuery
	err := DecodeUrlQuery([]byte(queryStr), &resultMap)
	if err != nil {
		return DuibaSdkReceiveCreditNotify{}, err
	}
	err = this.checkSign(resultMap)
	if err != nil {
		return DuibaSdkReceiveCreditNotify{}, err
	}
	var result DuibaSdkReceiveCreditNotify
	err = MapToArray(resultMap, &result, "url")
	if err != nil {
		return DuibaSdkReceiveCreditNotify{}, err
	}
	return result, nil
}

func (this *DuibaSdk) SendCreditNotify() []byte {
	return []byte("ok")
}
