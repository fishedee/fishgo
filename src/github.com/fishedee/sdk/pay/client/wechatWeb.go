package client

import (
	"bytes"
	// "encoding/json"
	"encoding/xml"
	//"errors"
	"fmt"
	"github.com/fishedee/sdk/pay/common"
	"github.com/fishedee/sdk/pay/util"
	"errors"
	"time"
)

var defaultWechatWebClient *WechatWebClient

func InitWxWebClient(c *WechatWebClient) {
	defaultWechatWebClient = c
}

func DefaultWechatWebClient() *WechatWebClient {
	return defaultWechatWebClient
}

// WechatWebClient 微信公众号支付
type WechatWebClient struct {
	AppID       string // 公众账号ID
	MchID       string // 商户号ID
	CallbackURL string // 回调地址
	Key         string // 密钥
	PayURL      string // 支付地址
	QueryURL    string // 查询地址
}

// Pay 支付
func (this *WechatWebClient) Pay(charge *common.Charge) (map[string]string, error) {
	var m = make(map[string]string)
	m["appid"] = this.AppID
	m["mch_id"] = this.MchID
	m["nonce_str"] = util.RandomStr()
	m["body"] = TruncatedText(charge.Describe,32)
	m["out_trade_no"] = charge.TradeNum
	m["total_fee"] = fmt.Sprintf("%d", int(charge.MoneyFee*100))
	m["spbill_create_ip"] = util.LocalIP()
	m["notify_url"] = charge.CallbackURL
	m["trade_type"] = "JSAPI"
	m["openid"] = charge.OpenID
	m["sign_type"] = "MD5"

	sign ,err := WechatGenSign(this.Key,m)
	if err != nil {
		return map[string]string{}, err
	}
	m["sign"] = sign

	// 转出xml结构
	xmlRe ,err := PostWechat(this.PayURL,m)
	if  err != nil{
		return map[string]string{}, err
	}

	var c = make(map[string]string)
	c["appId"] = this.AppID
	c["timeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
	c["nonceStr"] = util.RandomStr()
	c["package"] = fmt.Sprintf("prepay_id=%s", xmlRe.PrepayID)
	c["signType"] = "MD5"
	sign2 ,err := WechatGenSign(this.Key,c)
	if err != nil {
		return map[string]string{}, errors.New("WechatWeb: " + err.Error())
	}
	c["paySign"] = sign2

	return c, nil
}

// QueryOrder 查询订单
func (this *WechatWebClient) QueryOrder(tradeNum string) (*common.WeChatQueryResult, error) {
	var queryResult common.WeChatQueryResult
	var m = make(map[string]string)
	m["appid"] = this.AppID
	m["mch_id"] = this.MchID
	m["out_trade_no"] = tradeNum
	m["nonce_str"] = util.RandomStr()

	sign ,err := WechatGenSign(this.Key,m)
	if err != nil {
		return &queryResult, err
	}
	m["sign"] = sign

	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	xmlStr := fmt.Sprintf("<xml>%s</xml>", buf.String())

	result, err := HTTPSC.PostData(this.QueryURL, "text/xml:charset=UTF-8", xmlStr)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(result, &queryResult)
	return &queryResult, errors.New("xml.Unmarshal: " + err.Error())
}
