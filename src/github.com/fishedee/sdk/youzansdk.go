package sdk

import (
	"fmt"
	. "github.com/fishedee/crypto"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	"reflect"
	"time"
)

type YouzanSdk struct {
	AppId     string
	AppSecret string
}

type YouzanSdkError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
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

func (this *YouzanSdk) getSign(data map[string]interface{}) string {
	dataKeysInterface, _ := ArrayKeyAndValue(data)
	dataKeys := dataKeysInterface.([]string)
	dataKeys = ArraySort(dataKeys).([]string)

	result := ""
	for _, singleDataKey := range dataKeys {
		if singleDataKey == "sign" {
			continue
		} else {
			result += fmt.Sprintf("%v%v", singleDataKey, data[singleDataKey])
		}
	}
	result = this.AppSecret + result
	return CryptoMd5([]byte(result))
}

func (this *YouzanSdk) getMethodSign(method string, data map[string]interface{}) string {
	var data2 map[string]interface{}
	reflect.Copy(reflect.ValueOf(data2), reflect.ValueOf(data))
	data2["app_id"] = this.AppId
	data2["method"] = method
	data2["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	data2["format"] = "json"
	data2["v"] = "1.0"
	data2["sign_method"] = "md5"

	return this.getSign(data2)
}

func (this *YouzanSdk) api(method string, data interface{}, responseData interface{}) error {
	var dataByte []byte
	err := DefaultAjaxPool.Get(&Ajax{
		Url:          "https://open.koudaitong.com/api/entry",
		Data:         data,
		DataType:     "url",
		ResponseData: &dataByte,
	})
	if err != nil {
		return err
	}

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
func (this *YouzanSdk) GetTradeSold(input YouzanSdkTradeSoldRequest) (YouzanSdkTradeSoldResponse, error) {
	var result YouzanSdkTradeSoldResponse
	err := this.api("kdt.trades.sold.get", input, &result)
	if err != nil {
		return YouzanSdkTradeSoldResponse{}, nil
	}
	return result, nil
}
