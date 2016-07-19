package sdk

import (
	"errors"
	"fmt"
	. "github.com/fishedee/crypto"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	"net/url"
	"time"
)

type YouzanSdk struct {
	AppId     string
	AppSecret string
}

type YouzanSdkOauthUserInfo struct {
	AppId     string    `url:"app_id"`
	Timestamp time.Time `url:"timestamp"`
	Custom    string    `url:"custom"`
	Subscribe int       `url:"subscribe"`
	FansId    int       `url:"fans_id"`
	OpenId    string    `url:"open_id"`
	NickName  string    `url:"nickname"`
	Sex       string    `url:"sex"`
	Country   string    `url:"country"`
	Province  string    `url:"province"`
	City      string    `url:"city"`
	Avatar    string    `url:"avatar"`
	Sign      string    `url:"sign"`
}

type YouzanSdkError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type YouzanSdkTradeRequest struct {
	Tid              string `url:"tid"`
	SubTradePageSize string `url:"sub_trade_page_size,omitempty"`
	SubTradePageNo   string `url:"sub_trade_page_no,omitempty"`
	Fields           string `url:"fields,omitempty"`
}

type YouzanSdkTradeResponse struct {
	Trade YouzanSdkTradeDetail `json:"trade"`
}

type YouzanSdkTradeSoldRequest struct {
	WeixinUserType int       `url:"weixin_user_type,omitempty"`
	WeixinUserId   int       `url:"weixin_user_id,omitempty"`
	Version        string    `url:"version,omitempty"`
	UseHasNext     bool      `url:"use_has_next,omitempty"`
	Type           string    `url:"type,omitempty"`
	Status         string    `url:"status,omitempty"`
	StartUpdate    time.Time `url:"start_update,omitempty"`
	StartCreated   time.Time `url:"start_created,omitempty"`
	SenderId       int       `url:"senderId,omitempty"`
	PageSize       int       `url:"page_size,omitempty"`
	PageNo         int       `url:"page_no,omitempty"`
	Keyword        string    `url:"keyword,omitempty"`
	Fields         string    `url:"fields,omitempty"`
	EndUpdate      time.Time `url:"end_update,omitempty"`
	EndCreated     time.Time `url:"end_created,omitempty"`
	BuyerNick      string    `url:"buyer_nick,omitempty"`
	BuyerId        int       `url:"buyer_id,omitempty"`
	BuyWay         string    `url:"buy_way,omitempty"`
}

// write by lwz 2016/06/16
// 针对有赞新开放接口
// add Start
type YouzanSdkTradeSoldGetForOuterRequest struct {
	UseHasNext   bool      `url:"use_has_next,omitempty"`
	Type         string    `url:"type,omitempty"`
	Status       string    `url:"status,omitempty"`
	StartUpdate  time.Time `url:"start_update,omitempty"`
	StartCreated time.Time `url:"start_created,omitempty"`
	PageSize     int       `url:"page_size,omitempty"`
	PageNo       int       `url:"page_no,omitempty"`
	OuterUserId  string    `url:"outer_user_id"`
	OuterType    string    `url:"outer_type"`
	Fields       string    `url:"fields,omitempty"`
	EndUpdate    time.Time `url:"end_update,omitempty"`
	EndCreated   time.Time `url:"end_created,omitempty"`
}

// 订单
type YouzanSdkTradeOuterResponse struct {
	Trade YouzanSdkTradeOuterDetail `json:"trade"`
}

// 订单详细参数信息
type YouzanSdkTradeOuterDetail struct {
	Num              int                         `json:"num"`
	GoodsKind        int                         `json:"goods_kind"`
	NumIid           int                         `json:"num_iid"`
	Price            float64                     `json:"price"`
	PicPath          string                      `json:"pic_path"`
	PicThumbPath     string                      `json:"pic_thumb_path"`
	Title            string                      `json:"title"`
	Type             string                      `json:"type"`
	DiscountFee      float64                     `json:"discount_fee"`
	Status           string                      `json:"status"`
	StatusStr        string                      `json:"status_str"`
	RefundState      string                      `json:"refund_state"`
	ShippingType     string                      `json:"shipping_type"`
	PostFee          float64                     `json:"post_fee"`
	TotalFee         float64                     `json:"total_fee"`
	RefundedFee      float64                     `json:"refunded_fee"`
	Payment          float64                     `json:"payment"`
	Created          time.Time                   `json:"created"`
	UpdateTime       time.Time                   `json:"update_time"`
	PayTime          time.Time                   `json:"pay_time"`
	PayType          string                      `json:"pay_type"`
	ConsignTime      time.Time                   `json:"consign_time"`
	SignTime         time.Time                   `json:"sign_time"`
	BuyerArea        string                      `json:"buyer_area"`
	SellerFlag       int                         `json:"seller_flag"`
	BuyerMessage     string                      `json:"buyer_message"`
	Orders           []YouzanSdkTradeOrderOuter  `json:"orders"`
	FetchDetail      []YouzanSdkTradeFetch       `json:"fetch_detail"`
	CouponDetails    []YouzanSdkUmpTradeCoupon   `json:"coupon_details"`
	PromotionDetails []YouzanSdkTradePromotion   `json:"promotion_details"`
	AdjustFee        float64                     `json:"adjust_fee"`
	SubTrades        []YouzanSdkTradeOuterDetail `json:"sub_trades"`
	WeixinUserId     string                      `json:"weixin_user_id"`
	ButtonList       []YouzanSdkTradeButtonOuter `json:"button_list"`
	FeedBackNum      int                         `json:"feedback_num"`
	TradeMemo        string                      `json:"trade_memo"`
	FansInfo         YouzanSdkTradeFansOuter     `json:"fans_info"`
	BuyWayStr        string                      `json:"buy_way_str"`
	PfBuyWayStr      string                      `json:"pf_buy_way_str“`
	SendNum          int                         `json:"send_num"`
	UserId           string                      `json:"user_id"`
	Kind             int                         `json:"kind"`
	RelationType     string                      `json:"relation_type"`
	Relations        []string                    `json:"relations"`
	OutTradeNo       []string                    `json:"out_trade_no"`
	GroupNo          string                      `json:"group_no"`
	OuterUserId      int                         `json:"outer_user_id"`
	BuyerNick        string                      `json:"buyer_nick"`
	Tid              string                      `json:"tid"`
	BuyerType        int                         `json:"buyer_type"`
	BuyerId          string                      `json:"buyer_id"`
	ReceiverCity     string                      `json:"receiver_city"`
	ReceiverDistrict string                      `json:"receiver_district"`
	ReceiverName     string                      `json:"receiver_name"`
	ReceiverState    string                      `json:"receiver_state"`
	ReceiverAddress  string                      `json:"receiver_address"`
	ReceiverZip      string                      `json:"receiver_zip"`
	ReceiverMobile   string                      `json:"receiver_mobile"`
	Feedback         string                      `json:"feedback"`
	OuterTid         string                      `json:"outer_tid"`
}

// Order参数信息
type YouzanSdkTradeOrderOuter struct {
	Oid                   int                            `json:"oid"`
	OuterSkuId            string                         `json:"outer_sku_id"`
	OuterItemId           string                         `json:"outer_item_id"`
	Title                 string                         `json:"title"`
	SellerNick            string                         `json:"seller_nick"`
	FenxiaoPrice          float64                        `json:"fenxiao_price"`
	FenxiaoPayment        float64                        `json:"fenxiao_payment"`
	Price                 float64                        `json:"price"`
	TotalFee              float64                        `json:"total_fee"`
	Payment               float64                        `json:"payment"`
	DiscountFee           float64                        `json:"discount_fee"`
	SkuId                 int                            `json:"sku_id"`
	SkuUniqueCode         string                         `json:"sku_unique_code"`
	SkuPropertiesName     string                         `json:"sku_properties_name"`
	PicPath               string                         `json:"pic_path"`
	PicThumbPath          string                         `json:"pic_thumb_path"`
	ItemType              int                            `json:"item_type"`
	BuyerMessages         []YouzanSdkTradeBuyerMessage   `json:"buyer_messages"`
	OrderPromotionDetails []YouzanSdkTradeOrderPromotion `json:"order_promotion_details"`
	StateStr              string                         `json:"state_str"`
	AllowSend             int                            `json:"allow_send"`
	IsSend                int                            `json:"is_send"`
	ItemRefundState       string                         `json:"item_refund_state"`
	Num                   int                            `json:"num"`
	NumIid                int                            `json:"num_iid"`
}

// YouzanSdkTradeButtonOuter参数信息
type YouzanSdkTradeButtonOuter struct {
	ToolIcon      string    `json:"tool_icon"`
	ToolTitle     string    `json:"tool_title"`
	ToolValue     string    `json:"tool_value"`
	ToolType      string    `json:"tool_type"`
	ToolParameter string    `json:"tool_parameter"`
	NewSign       string    `json:"new_sign"`
	CreateTime    time.Time `json:"create_time"`
}

// YouzanSdkTradeFansOuter参数信息
type YouzanSdkTradeFansOuter struct {
	FansNickname string `json:"fans_nickname"`
	FansId       string `json:"fans_id"`
	BuyerId      string `json:"buyer_id"`
	FansType     string `json:"fans_type"`
}

// add End

type YouzanSdkTradeBuyerMessage struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type YouzanSdkTradeOrderPromotion struct {
	PromotionType string    `json:"promotion_type"`
	ApplyAt       time.Time `json:"apply_at"`
	PromotionName string    `json:"promotion_name"`
	DiscountFee   float64   `json:"discount_fee"`
}

type YouzanSdkTradeOrder struct {
	OuterSkuId            string                         `json:"outer_sku_id"`
	SkuUniqueCode         string                         `json:"sku_unique_code"`
	OuterItemId           string                         `json:"outer_item_id"`
	PicThumbPath          string                         `json:"pic_thumb_path"`
	ItemType              int                            `json:"item_type"`
	Num                   int                            `json:"num"`
	NumIid                int                            `json:"num_iid"`
	SkuId                 int                            `json:"sku_id"`
	SkuPropertiesName     string                         `json:"sku_properties_name"`
	PicPath               string                         `json:"pic_path"`
	Oid                   int                            `json:"oid"`
	Title                 string                         `json:"title"`
	FenxiaoPayment        float64                        `json:"fenxiao_payment"`
	ItemRefundState       string                         `json:"item_refund_state"`
	DiscountFee           float64                        `json:"discount_fee"`
	BuyerMessages         []YouzanSdkTradeBuyerMessage   `json:"buyer_messages"`
	StateStr              string                         `json:"state_str"`
	OrderPromotionDetails []YouzanSdkTradeOrderPromotion `json:"order_promotion_details"`
	Price                 float64                        `json:"price"`
	FenxiaoPrice          float64                        `json:"fenxiao_price"`
	TotalFee              float64                        `json:"total_fee"`
	Payment               float64                        `json:"payment"`
	SellerNick            string                         `json:"seller_nick"`
}

type YouzanSdkTradePromotion struct {
	PromotionType      string    `json:"promotion_type"`
	UsedAt             time.Time `json:"used_at"`
	PromotionName      string    `json:"promotion_name"`
	PromotionCondition string    `json:"promotion_condition"`
	PromotionId        int       `json:"promotion_id"`
	DiscountFee        float64   `json:"discount_fee"`
}

type YouzanSdkUmpTradeCoupon struct {
	CouponDescription string    `json:"coupon_description"`
	UsedAt            time.Time `json:"used_at"`
	CouponCondition   string    `json:"coupon_condition"`
	CouponId          int       `json:"coupon_id"`
	CouponContent     string    `json:"coupon_content"`
	CouponName        string    `json:"coupon_name"`
	CouponType        string    `json:"coupon_type"`
	DiscountFee       float64   `json:"discount_fee"`
}

type YouzanSdkTradeFetch struct {
	FetcherName   string    `json:"fetcher_name"`
	ShopState     string    `json:"shop_state"`
	ShopMobile    string    `json:"shop_mobile"`
	ShopCity      string    `json:"shop_city"`
	ShopDistrict  string    `json:"shop_district"`
	FetcherMobile string    `json:"fetcher_mobile"`
	ShopName      string    `json:"shop_name"`
	ShopAddress   string    `json:"shop_address"`
	FetchTime     time.Time `json:"fetch_time"`
}

type YouzanSdkTradeDetail struct {
	ConsignTime      time.Time                 `json:"consign_time"`
	BuyerArea        string                    `json:"buyer_area"`
	Num              int                       `json:"num"`
	AdjustFee        float64                   `json:"adjust_fee"`
	RelationType     string                    `json:"relation_type"`
	Type             string                    `json:"type"`
	BuyerId          int                       `json:"buyer_id"`
	Tid              string                    `json:"tid"`
	Feedback         int                       `json:"feedback"`
	Price            float64                   `json:"price"`
	TotalFee         float64                   `json:"total_fee"`
	Payment          float64                   `json:"payment"`
	WeixinUserId     int                       `json:"weixin_user_id"`
	SubTrades        []YouzanSdkTradeDetail    `json:"sub_trades"`
	BuyerMessage     string                    `json:"buyer_message"`
	Created          time.Time                 `json:"created"`
	PayTime          time.Time                 `json:"pay_time"`
	OutTradeNo       []string                  `json:"out_trade_no"`
	Orders           []YouzanSdkTradeOrder     `json:"orders"`
	PromotionDetails []YouzanSdkTradePromotion `json:"promotion_details"`
	RefundState      string                    `json:"refund_state"`
	Status           string                    `json:"status"`
	PostFee          float64                   `json:"post_fee"`
	PicThumbPath     string                    `json:"pic_thumb_path"`
	ReceiverCity     string                    `json:"receiver_city"`
	ShippingType     string                    `json:"shipping_type"`
	RefundedFee      float64                   `json:"refunded_fee"`
	NumIid           int                       `json:"num_iid"`
	Title            string                    `json:"title"`
	DiscountFee      float64                   `json:"discount_fee"`
	ReceiverState    string                    `json:"receiver_state"`
	UpdateTime       time.Time                 `json:"update_time"`
	CouponDetails    []YouzanSdkUmpTradeCoupon `json:"coupon_details"`
	ReceiverZip      string                    `json:"receiver_zip"`
	ReceiverName     string                    `json:"receiver_name"`
	PayType          string                    `json:"pay_type"`
	Profit           float64                   `json:"profit"`
	FetchDetail      []YouzanSdkTradeFetch     `json:"fetch_detail"`
	BuyerType        int                       `json:"buyer_type"`
	ReceiverDistrict string                    `json:"receiver_district"`
	PicPath          string                    `json:"pic_path"`
	ReceiverMobile   string                    `json:"receiver_mobile"`
	SignTime         time.Time                 `json:"sign_time"`
	SellerFlag       int                       `json:"seller_flag"`
	BuyerNick        string                    `json:"buyer_nick"`
	Handled          int                       `json:"handled"`
	ReceiverAddress  string                    `json:"receiver_address"`
	TradeMemo        string                    `json:"trade_memo"`
	Relations        []string                  `json:"relations"`
	OuterTid         string                    `json:"outer_tid"`
}

type YouzanSdkTradeSoldResponse struct {
	TotalResults int                    `json:"total_results"`
	Trades       []YouzanSdkTradeDetail `json:"trades"`
	HasNext      bool                   `json:"has_next"`
}

func (this *YouzanSdkError) GetCode() int {
	return this.Code
}

func (this *YouzanSdkError) GetMsg() string {
	return this.Message
}

func (this *YouzanSdkError) Error() string {
	return fmt.Sprintf("错误码为：%v，错误描述为：%v", this.Code, this.Message)
}

func (this *YouzanSdk) getNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (this *YouzanSdk) getSign(data map[string]interface{}) (string, error) {
	dataKeysInterface, _ := ArrayKeyAndValue(data)
	dataKeys := dataKeysInterface.([]string)
	dataKeys = ArraySort(dataKeys).([]string)

	result := ""
	for _, singleDataKey := range dataKeys {
		if singleDataKey == "sign" {
			continue
		} else {
			singleDataValue := fmt.Sprintf("%v", data[singleDataKey])
			result += singleDataKey + singleDataValue
		}
	}
	result = this.AppSecret + result + this.AppSecret
	return CryptoMd5([]byte(result)), nil
}

func (this *YouzanSdk) getMethodSign(method string, data map[string]interface{}) (map[string]interface{}, error) {
	data["app_id"] = this.AppId
	data["method"] = method
	data["timestamp"] = this.getNowTime()
	data["format"] = "json"
	data["v"] = "1.0"
	data["sign_method"] = "md5"
	sign, err := this.getSign(data)
	if err != nil {
		return nil, err
	}
	data["sign"] = sign
	return data, nil
}

func (this *YouzanSdk) api(method string, data interface{}, responseData interface{}) error {
	//调整输入参数
	queryParamsInterface := ArrayToMap(data, "url")
	queryParams := queryParamsInterface.(map[string]interface{})
	queryParams, err := this.getMethodSign(method, queryParams)
	if err != nil {
		return err
	}

	//请求
	var dataByte []byte
	queryBytes, err := EncodeUrlQuery(queryParams)
	if err != nil {
		return err
	}

	err = DefaultAjaxPool.Get(&Ajax{
		Url:          "https://open.koudaitong.com/api/entry",
		Data:         queryBytes,
		ResponseData: &dataByte,
	})
	if err != nil {
		return err
	}

	//分析输出参数
	var responseMap map[string]interface{}
	err = DecodeJson(dataByte, &responseMap)
	if err != nil {
		return err
	}
	if errorData, isErr := responseMap["error_response"]; isErr {
		var youzanErr YouzanSdkError
		err := MapToArray(errorData, &youzanErr, "json")
		if err != nil {
			return err
		} else {
			return &youzanErr
		}
	}

	err = MapToArray(responseMap["response"], responseData, "json")
	if err != nil {
		return err
	}
	return nil
}

//交易接口
func (this *YouzanSdk) GetTrade(input YouzanSdkTradeRequest) (YouzanSdkTradeResponse, error) {
	var result YouzanSdkTradeResponse
	err := this.api("kdt.trade.get", input, &result)
	if err != nil {
		return YouzanSdkTradeResponse{}, err
	}
	return result, nil
}

func (this *YouzanSdk) GetTradeSold(input YouzanSdkTradeSoldRequest) (YouzanSdkTradeSoldResponse, error) {
	var result YouzanSdkTradeSoldResponse
	err := this.api("kdt.trades.sold.get", input, &result)
	if err != nil {
		return YouzanSdkTradeSoldResponse{}, err
	}
	return result, nil
}

// write by lwz 2016/0615
//根据第三方用户userId获取交易订单列表
func (this *YouzanSdk) GetTradeSoldFoRouter(input YouzanSdkTradeSoldGetForOuterRequest) (YouzanSdkTradeSoldResponse, error) {
	var result YouzanSdkTradeSoldResponse
	err := this.api("kdt.trades.sold.getforouter", input, &result)
	if err != nil {
		return YouzanSdkTradeSoldResponse{}, err
	}

	return result, nil
}

//授权接口
func (this *YouzanSdk) GetOauthUrl(redirectUrl string, scope string, state string) (string, error) {
	inputData := map[string]interface{}{
		"app_id":       this.AppId,
		"redirect_url": redirectUrl,
		"scope":        scope,
		"timestamp":    this.getNowTime(),
		"custom":       state,
	}
	sign, err := this.getSign(inputData)
	if err != nil {
		return "", err
	}
	inputData["sign"] = sign

	inputDataStr, err := EncodeUrlQuery(inputData)
	if err != nil {
		return "", err
	}
	return "http://wap.koudaitong.com/v2/open/weixin/auth?" + string(inputDataStr), nil
}

func (this *YouzanSdk) GetOauthUserInfo(requestUrl *url.URL) (YouzanSdkOauthUserInfo, error) {
	//获取输入
	urlQuery := requestUrl.RawQuery
	var inputRequest map[string]interface{}
	err := DecodeUrlQuery([]byte(urlQuery), &inputRequest)
	if err != nil {
		return YouzanSdkOauthUserInfo{}, err
	}
	if msg, isErr := inputRequest["msg"]; isErr {
		return YouzanSdkOauthUserInfo{}, &YouzanSdkError{1, fmt.Sprintf("%v", msg)}
	}

	//验证签名
	sign, err := this.getSign(inputRequest)
	if err != nil {
		return YouzanSdkOauthUserInfo{}, err
	}
	if inputRequest["sign"] != sign {
		return YouzanSdkOauthUserInfo{}, errors.New(fmt.Sprintf("签名失败 %v != %v", inputRequest["sign"], sign))
	}

	//提取输出
	var result YouzanSdkOauthUserInfo
	err = MapToArray(inputRequest, &result, "url")
	if err != nil {
		return YouzanSdkOauthUserInfo{}, err
	}
	return result, nil
}
