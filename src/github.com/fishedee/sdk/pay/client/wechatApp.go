package client

import (
	"bytes"
	"crypto/md5"
	// "encoding/json"
	"encoding/xml"
	//"errors"
	"errors"
	"fmt"
	"github.com/fishedee/sdk/pay/common"
	"github.com/fishedee/sdk/pay/util"
	"sort"
	"strings"
	"time"
)

var defaultWechatAppClient *WechatAppClient

func InitWxAppClient(c *WechatAppClient) {
	defaultWechatAppClient = c
}

// DefaultWechatAppClient 默认微信app客户端
func DefaultWechatAppClient() *WechatAppClient {
	return defaultWechatAppClient
}

// WechatAppClient 微信app支付
type WechatAppClient struct {
	AppID       string // AppID
	MchID       string // 商户号ID
	CallbackURL string // 回调地址
	Key         string // 密钥
	PayURL      string // 支付地址
}

// Pay 支付
func (this *WechatAppClient) Pay(charge *common.Charge) (map[string]string, error) {
	var m = make(map[string]string)
	m["appid"] = this.AppID
	m["mch_id"] = this.MchID
	m["nonce_str"] = util.RandomStr()
	m["body"] = TruncatedText(charge.Describe,128)
	m["out_trade_no"] = charge.TradeNum
	m["total_fee"] = fmt.Sprintf("%d", int(charge.MoneyFee*100))
	m["spbill_create_ip"] = util.LocalIP()
	m["notify_url"] = charge.CallbackURL
	m["trade_type"] = "APP"
	m["sign_type"] = "MD5"

	sign := this.GenSign(m)

	m["sign"] = sign
	// 转出xml结构
	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	xmlStr := fmt.Sprintf("<xml>%s</xml>", buf.String())

	re, err := HTTPSC.PostData(this.PayURL, "text/xml:charset=UTF-8", xmlStr)
	if err != nil {
		return map[string]string{}, errors.New("HTTPSC.PostData: " + err.Error())
	}

	var xmlRe common.WeChatReResult
	err = xml.Unmarshal(re, &xmlRe)
	if err != nil {
		return map[string]string{}, errors.New("xml.Unmarshal: " + err.Error())
	}

	if xmlRe.ReturnCode != "SUCCESS" {
		// 通信失败
		return map[string]string{}, errors.New("xmlRe.ReturnMsg: " + xmlRe.ReturnMsg)
	}

	if xmlRe.ResultCode != "SUCCESS" {
		// 支付失败
		return map[string]string{}, errors.New("xmlRe.ErrCodeDes: " + xmlRe.ErrCodeDes)
	}

	var c = make(map[string]string)
	c["appid"] = this.AppID
	c["partnerid"] = this.MchID
	c["prepayid"] = xmlRe.PrepayID
	c["package"] = "Sign=WXPay"
	c["noncestr"] = util.RandomStr()
	c["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())

	sign2 ,err := WechatGenSign(this.Key,m)
	if err != nil {
		return map[string]string{}, err
	}
	c["paySign"] = strings.ToUpper(sign2)

	return c, nil
}

// GenSign 产生签名
func (this *WechatAppClient) GenSign(m map[string]string) string {
	var signData []string
	for k, v := range m {
		if v != "" && k != "sign" && k != "key" {
			signData = append(signData, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + this.Key
	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return ""
	}
	signByte := c.Sum(nil)
	if err != nil {
		return ""
	}
	return strings.ToUpper(fmt.Sprintf("%x", signByte))
}
