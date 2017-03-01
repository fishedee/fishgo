package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fishedee/sdk/pay/common"
	"net/url"
	"sort"
	"strings"
	"time"
	"github.com/go-errors/errors"
)

var defaultAliAppClient *AliAppClient

type AliAppClient struct {
	SellerID   string //合作者ID
	AppID      string // 应用ID
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func InitAliAppClient(c *AliAppClient) {
	defaultAliAppClient = c
}

// DefaultAliAppClient 得到默认支付宝app客户端
func DefaultAliAppClient() *AliAppClient {
	return defaultAliAppClient
}

func (this *AliAppClient) Pay(charge *common.Charge) (map[string]string, error) {
	var m = make(map[string]string)
	var bizContent = make(map[string]string)
	m["app_id"] = this.AppID
	m["method"] = "alipay.trade.app.pay"
	m["format"] = "JSON"
	m["charset"] = "utf-8"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["notify_url"] = charge.CallbackURL
	m["sign_type"] = "RSA"
	//m["subject"] = charge.Describe
	//m["out_trade_no"] = charge.TradeNum
	//m["product_code"] = "QUICK_MSECURITY_PAY"
	//m["total_amount"] = fmt.Sprintf("%.2f", charge.MoneyFee)
	bizContent["subject"] = TruncatedText(charge.Describe,64)
	bizContent["out_trade_no"] = charge.TradeNum
	bizContent["product_code"] = "QUICK_MSECURITY_PAY"
	bizContent["total_amount"] = fmt.Sprintf("%.2f", charge.MoneyFee)

	bizContentJson, err := json.Marshal(bizContent)
	if err != nil {
		return map[string]string{}, errors.New("json.Marshal: "+err.Error())
	}
	m["biz_content"] = string(bizContentJson)

	m["sign"] = this.GenSign(m)

	fmt.Printf("%+v", m)
	return m, nil
}

// GenSign 产生签名
func (this *AliAppClient) GenSign(m map[string]string) string {
	var data []string
	for k, v := range m {
		if v != "" && k != "sign" {
			data = append(data, fmt.Sprintf(`%s=%s`, k, v))
		}
	}
	sort.Strings(data)
	signData := strings.Join(data, "&")
	fmt.Println(signData)
	s := sha1.New()
	_, err := s.Write([]byte(signData))
	if err != nil {
		panic(err)
	}
	hashByte := s.Sum(nil)
	signByte, err := this.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA1)
	if err != nil {
		panic(err)
	}
	return url.QueryEscape(base64.StdEncoding.EncodeToString(signByte))
}

// CheckSign 检测签名
func (this *AliAppClient) CheckSign(signData, sign string) {
	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		panic(err)
	}
	s := sha1.New()
	_, err = s.Write([]byte(signData))
	if err != nil {
		panic(err)
	}
	hash := s.Sum(nil)
	err = rsa.VerifyPKCS1v15(this.PublicKey, crypto.SHA1, hash, signByte)
	if err != nil {
		panic(err)
	}
}

// ToURL
func (this *AliAppClient) ToURL(m map[string]string) string {
	var buf []string
	for k, v := range m {
		buf = append(buf, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	return strings.Join(buf, "&")
}
