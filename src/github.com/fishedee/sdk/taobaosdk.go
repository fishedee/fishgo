// 淘宝开放平台接口
package sdk

import (
	"errors"
	"strings"
	"time"

	"github.com/fishedee/crypto"
	"github.com/fishedee/encoding"
	"github.com/fishedee/language"
	"github.com/fishedee/util"
)

// TaoBao Sdk
type TaoBaoSdk struct {
	AppKey    string
	AppSecret string
	Type      string // 环境类型 沙箱：sandbox; 正式：online; 海外：oversea;
}

// 调用入口
// 沙箱测试环境： ISV软件的测试环境，应用创建后即可使用。此环境提供简化版的淘宝网，支持大部分场景的API调用，沙箱环境的权限和流量均无限制，可放开使用
const (
	sandboxHttpUrl  = "http://gw.api.tbsandbox.com/router/rest"
	sandboxHttpsUrl = "https://gw.api.tbsandbox.com/router/rest"
)

//
// 正式测试环境： ISV软件上线之前的正式模拟环境，应用创建成功后即可使用。此环境主要是针对部分无法在沙箱完成测试的场景使用，限制API调用为5000次/天，授权用户数量为5个，所能调用的API与应用拥有的权限能力一致。
//
// 正式环境： ISV软件上线之后使用的环境，此环境的入口与正式测试环境一致，只不过应用上线之后，流量限制会进行打开，具体流量限制与应用所属类目有关，比如服务市场类的应用，限制API调用为100万次/天。
const (
	onlineHttpUrl  = "http://gw.api.taobao.com/router/rest"
	onlineHttpsUrl = "https://eco.taobao.com/router/rest"
)

//
// 海外环境： 海外环境也属于正式环境的一种，主要是给海外（欧美国家）ISV使用，对于海外的ISV，使用海外环境会比国内环境的性能高一倍。
const (
	overseaHttpUrl  = "http://api.taobao.com/router/rest"
	overseaHttpsUrl = "https://api.taobao.com/router/rest"
)

// 淘宝接口公共参数
type TaoBaoCommonParam struct {
	// API接口名称 必须
	Method string `json:"method"`
	// TOP分配给应用的AppKey(随环境不同而不同) 必须
	AppKey string `json:"app_key"`
	// 用户登录授权成功后，TOP颁发给应用的授权信息
	Session string `json:"session"`
	// 时间戳，格式为yyyy-MM-dd HH:mm:ss，时区为GMT+8，例如：2016-01-01 12:00:00 淘宝API服务端允许客户端请求最大时间误差为10分钟 必须
	Timestamp string `json:"timestamp"`
	// 响应格式。默认为xml格式，可选值：xml，json。
	Format string `json:"format"`
	// API协议版本，可选值：2.0。 必须
	V string `json:"v"`
	// 合作伙伴身份标识。
	PartnerId string `json:"partner_id"`
	// 被调用的目标AppKey，仅当被调用的API为第三方ISV提供时有效。
	TargetAppKey string `json:"target_app_key"`
	// 是否采用精简JSON返回格式，仅当format=json时有效，默认值为：false。
	Simplify bool `json:"simplify"`
	// 签名的摘要算法，可选值为：hmac，md5。必须
	SignMethod string `json:"sign_method"`
	// API输入参数签名结果 必须
	Sign string `json:"sign"`
}

// 淘宝客专用参数
type TaoBaoKeParam struct {
	// 公共参数
	TaoBaoCommonParam
	// 需返回的字段列表 必须
	Fields string `json:"fields"`
	// 查询 特殊可选
	Q string `json:"q"`
	// 后台类目ID，用,分割，最大10个 特殊可选
	Cat string `json:"cat"`
	// 所在地
	Itemloc string `json:"itemloc"`
	// 排序_des（降序），排序_asc（升序），销量（total_sales），淘客佣金比率（tk_rate）， 累计推广量（tk_total_sales），总支出佣金（tk_total_commi）
	Sort string `json:"sort"`
	// 是否商城商品，设置为true表示该商品是属于淘宝商城商品，设置为false或不设置表示不判断这个属性
	IsTmall bool `json:"is_tmall"`
	// 是否海外商品，设置为true表示该商品是属于海外商品，设置为false或不设置表示不判断这个属性
	IsOverseas bool `json:"is_overseas"`
	// 折扣价范围下限，单位：元
	StartPrice int `json:"start_price"`
	// 折扣价范围上限，单位：元
	EndPrice int `json:"end_price"`
	// 淘客佣金比率上限，如：1234表示12.34%
	StartTkRate int `json:"start_tk_rate"`
	// 淘客佣金比率下限，如：1234表示12.34%
	EndTkRate int `json:"end_tk_rate"`
	// 链接形式：1：PC，2：无线，默认：１
	Platform int `json:"platform"`
	// 第几页，默认：１
	PageNo int `json:"page_no"`
	// 页大小，默认20，1~100
	PageSize int `json:"page_size"`
	// 商品ID
	NumIIds string `json:"num_iids"`
	// 返回数量
	Count int `json:"count"`
	// 	信用等级下限
	StartCredit int `json:"start_credit"`
	// 信用等级上限
	EndCredit int `json:"end_credit"`
	// 淘客佣金比率下限
	StartCommissionRate int `json:"start_commission_rate"`
	// 淘客佣金比率上限
	EndCommissionRate int `json:"end_commission_rate"`
	// 店铺商品总数下限
	StartTotalAction int `json:"start_total_action"`
	// 店铺商品总数上限
	EndTotalAction int `json:"end_total_action"`
	// 累计推广商品下限
	StartAuctionCount int `json:"start_auction_count"`
	// 累计推广商品上限
	EndAuctionCount int `json:"end_auction_count"`
	// 推广位id
	AdzoneId string `json:"adzone_id"`
	// 自定义输入串，英文和数字组成，长度不能大于12个字符，区分不同的推广渠道
	Unid string `json:"unid"`
	// 选品库的id
	FavoritesId string `json:"favorites_id"`
	// 1: 普通选品组; 2: 高佣选品组; -1: 同时输出所有类型的选品组
	Type int `json:"type"`
	// 商品ID
	NumIId string `json:"num_iid"`
	// 招商活动ID
	EventId string `json:"event_id"`
	// 最早开团时间
	StartTime string `json:"start_time"`
	// 最晚开团时间
	EndTime string `json:"end_time"`
	// 卖家ID
	UserId string `json:"user_id"`
	// 卖家IDs
	UserIds string `json:"user_ids"`
}

// 淘宝客商品
type TaoBaoKeItem struct {
	// 商品ID
	NumIId int `json:"num_iid"`
	// 商品标题
	Title string `json:"title"`
	// 商品主图
	PictUrl string `json:"pict_url"`
	// 商品小图列表
	SmallImages map[string]string `json:"small_images"`
	// 商品一口价格
	ReservePrice string `json:"reserve_price"`
	// 商品折扣价格
	ZkFinalPrice string `json:"zk_final_price"`
	// 卖家类型，0表示集市，1表示商城
	UserType int `json:"user_type"`
	// 宝贝所在地
	ProvCity string `json:"provcity"`
	// 商品地址
	ItemUrl string `json:"item_url"`
	// 卖家昵称
	Nick string `json:"nick"`
	// 卖家id
	SellerId int `json:"seller_id"`
	// 30天销量
	Volume int `json:"volume"`
	//
	ClickUrl string `json:"click_url"`
	//
	TkRate string `json:"tk_rate"`
	//
	ZkFinalPriceWap string `json:"zk_final_price_wap"`
	//
	EventStartTime string `json:"event_start_time"`
	//
	EventEndTime string `json:"event_end_time"`
	// 宝贝描述
	Description string `json:"description"`
	// 商品淘客地址
	ItemClickUrl string `json:"item_click_url"`
	// 商铺淘客地址
	ShopClickUrl string `json:"shop_click_url"`
}

// 淘宝客商品列表
type TaoBaoKeItems struct {
	NTbkItem []TaoBaoKeItem `json:"n_tbk_item"`
}

type TaoBaoKeItemDetails struct {
	NTbkItemDetail []TaoBaoKeItem `json:"n_tbk_item_detail"`
}

// 淘宝客商品列表及数量
type TaoBaoItemResults struct {
	Results      TaoBaoKeItems `json:"results"`
	TotalResults int           `json:"total_results"`
	RequestId    string        `json:"request_id"`
}

// 淘宝客商品详细
type TaoBaoItemDetailResults struct {
	Results      TaoBaoKeItemDetails `json:"results"`
	TotalResults int                 `json:"total_results"`
	RequestId    string              `json:"request_id"`
}

// 淘宝客搜索商品结果
type TaoBaoKeItemGetResponse struct {
	TbkItemGetResponse TaoBaoItemResults `json:"tbk_item_get_response"`
}

// 淘宝客商品详细结果-简版
type TaoBaoKeItemInfoGetResponse struct {
	TbkItemInfoGetResponse TaoBaoItemResults `json:"tbk_item_info_get_response"`
}

// 淘宝客商品详细结果-高级
type TaoBaoKeItemDetailInfoGetResponse struct {
	TbkItemDetailGetResponse TaoBaoItemDetailResults `json:"tbk_item_detail_get_response"`
}

// 淘宝客商品推荐结果
type TaoBaoKeItemRecommendGetResponse struct {
	TbkItemRecommendGetResponse TaoBaoItemResults `json:"tbk_item_recommend_get_response"`
}

type TaoBaoKeShop struct {
	UserId     int    `json:"user_id"`
	ShopTitle  string `json:"shop_title"`
	ShopType   string `json:"shop_type"`
	SellerNick string `json:"seller_nick"`
	PictUrl    string `json:"pict_url"`
	ShopUrl    string `json:"shop_url"`
	ClickUrl   string `json:"click_url"`
}

type TaoBaoKeShops struct {
	NTbkShop []TaoBaoKeShop `json:"n_tbk_shop"`
}

type TaoBaoKeShopResults struct {
	Results      TaoBaoKeShops `json:"results"`
	TotalResults int           `json:"total_results"`
	RequestId    string        `json:"request_id"`
}

// 淘宝商铺查询结果
type TaoBaoKeShopGetResponse struct {
	TbkShopGetResponse TaoBaoKeShopResults `json:"tbk_shop_get_response"`
}

// 淘宝商铺推荐结果
type TaoBaoKeShopRecommendGetResponse struct {
	TbkShopRecommendGetResponse TaoBaoKeShopResults `json:"tbk_shop_recommend_get_response"`
}

type TaoBaoKeUatmItems struct {
	UatmTbkItem []TaoBaoKeItem `json:"uatm_tbk_item"`
}

type TaoBaoKeUatmResults struct {
	Results      TaoBaoKeUatmItems `json:"results"`
	TotalResults int               `json:"total_results"`
	RequestId    string            `json:"request_id"`
}

// 淘宝联盟选品库的宝贝信息
type TaoBaoFavoriteItemGetResponse struct {
	TbkUatmFavoritesItemGetResponse TaoBaoKeUatmResults `json:"tbk_uatm_favorites_item_get_response"`
}

type TaoBaoKeFavorite struct {
	Type           int    `json:"type"`
	FavoritesId    int    `json:"favorites_id"`
	FavoritesTitle string `json:"favorites_title"`
}

type TaoBaoKeFavorites struct {
	TbkFavorites []TaoBaoKeFavorite `json:"tbk_favorites"`
}

type TaoBaoKeFavoriteResults struct {
	Results      TaoBaoKeFavorites `json:"results"`
	TotalResults int               `json:"total_results"`
	RequestId    string            `json:"request_id"`
}

// 淘宝联盟选品库的宝贝列表
type TaoBaoFavoritesGetResponse struct {
	TbkUatmFavoritesGetResponse TaoBaoKeFavoriteResults `json:"tbk_uatm_favorites_get_response"`
}

type TaoBaoKeUatmEvent struct {
	EventId    int    `json:"event_id"`
	EventTitle string `json:"event_title"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

type TaoBaoKeUatmEvents struct {
	TbkEvent []TaoBaoKeUatmEvent `json:"tbk_event"`
}

type TaoBaoKeUatmEventResults struct {
	Results      TaoBaoKeUatmEvents `json:"results"`
	TotalResults int                `json:"total_results"`
	RequestId    string             `json:"request_id"`
}

// 淘宝客定向招商活动基本信息
type TaoBaoKeUatmEventGetResponse struct {
	TbkUatmEventGetResponse TaoBaoKeUatmEventResults `json:"tbk_uatm_event_get_response"`
}

type TaoBaoKeJuTqg struct {
	Title        string `json:"title"`
	TotalAmount  int    `json:"total_amount"`
	ClickUrl     string `json:"click_url"`
	CategoryName string `json:"category_name"`
	ZkFinalPrice string `json:"zk_final_price"`
	EndTime      string `json:"end_time"`
	SoldNum      int    `json:"sold_num"`
	StartTime    string `json:"start_time"`
	ReservePrice string `json:"reserve_price"`
	PicUrl       string `json:"pic_url"`
	NumIId       int    `json:"num_iid"`
}

type TaoBaoKeJuTqgs struct {
	Results []TaoBaoKeJuTqg `json:"results"`
}

type TaoBaoKeJuTqgResults struct {
	Results      TaoBaoKeJuTqgs `json:"results"`
	TotalResults int            `json:"total_results"`
	RequestId    string         `json:"request_id"`
}

// 淘抢购商品列表
type TaoBaoKeJuTqgGetResponse struct {
	TbkJuTqgGetResponse TaoBaoKeJuTqgResults `json:"tbk_ju_tqg_get_response"`
}

// 淘宝客商品链接转换
type TaoBaoKeItemConvertResponse struct {
	TbkItemConvertResponse TaoBaoItemResults `json:"tbk_item_convert_response"`
}

// 淘宝客店铺链接转换
type TaoBaoKeShopConvertResponse struct {
	TbkShopConvertResponse TaoBaoKeShopResults `json:"tbk_shop_convert_response"`
}

// 淘宝客错误
type TaoBaoKeError struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

// 淘宝客错误结果
type TaoBaoKeErrorResponse struct {
	ErrorResponse TaoBaoKeError `json:"error_response"`
}

// 搜索淘宝客商品
func (this *TaoBaoSdk) GetTaoBaoKeAllItem(param TaoBaoKeParam) (TaoBaoKeItemGetResponse, error) {
	method := "taobao.tbk.item.get"
	result := TaoBaoKeItemGetResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 取淘宝客商品详细--简版
func (this *TaoBaoSdk) GetTaoBaoKeItemInfo(param TaoBaoKeParam) (TaoBaoKeItemInfoGetResponse, error) {
	method := "taobao.tbk.item.info.get"
	result := TaoBaoKeItemInfoGetResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 取淘宝客商品详细--高级
func (this *TaoBaoSdk) GetTaoBaoKeItemDetailInfo(param TaoBaoKeParam) (TaoBaoKeItemDetailInfoGetResponse, error) {
	method := "taobao.tbk.item.detail.get"
	result := TaoBaoKeItemDetailInfoGetResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 淘宝客商品链接转换--高级
func (this *TaoBaoSdk) ConvertTaoBaoKeItem(param TaoBaoKeParam) (TaoBaoKeItemConvertResponse, error) {
	method := "taobao.tbk.item.convert"
	result := TaoBaoKeItemConvertResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 淘宝客店铺链接转换--高级
func (this *TaoBaoSdk) ConvertTaoBaoKeShop(param TaoBaoKeParam) (TaoBaoKeShopConvertResponse, error) {
	method := "taobao.tbk.shop.convert"
	result := TaoBaoKeShopConvertResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 淘宝客商品关联推荐查询
func (this *TaoBaoSdk) GetTaoBaoKeItemRecommend(param TaoBaoKeParam) (TaoBaoKeItemRecommendGetResponse, error) {
	method := "taobao.tbk.item.recommend.get"
	result := TaoBaoKeItemRecommendGetResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 淘宝客店铺查询
func (this *TaoBaoSdk) GetTaoBaoKeShop(param TaoBaoKeParam) (TaoBaoKeShopGetResponse, error) {
	method := "taobao.tbk.shop.get"
	result := TaoBaoKeShopGetResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 淘宝客店铺推荐
func (this *TaoBaoSdk) GetTaoBaoKeShopRecommend(param TaoBaoKeParam) (TaoBaoKeShopRecommendGetResponse, error) {
	method := "taobao.tbk.shop.recommend.get"
	result := TaoBaoKeShopRecommendGetResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 获取淘宝联盟选品库的宝贝信息
func (this *TaoBaoSdk) GetTaoBaoKeUatmFavoriteItem(param TaoBaoKeParam) (TaoBaoFavoriteItemGetResponse, error) {
	method := "taobao.tbk.uatm.favorites.item.get"
	result := TaoBaoFavoriteItemGetResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 获取淘宝联盟选品库列表
func (this *TaoBaoSdk) GetTaoBaoKeUatmFavorites(param TaoBaoKeParam) (TaoBaoFavoritesGetResponse, error) {
	method := "taobao.tbk.uatm.favorites.get"
	result := TaoBaoFavoritesGetResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 淘客自己发起的，正在进行中的定向招商的活动列表
func (this *TaoBaoSdk) GetTaoBaoKeUatmEvents(param TaoBaoKeParam) (TaoBaoKeUatmEventGetResponse, error) {
	method := "taobao.tbk.uatm.event.get"
	result := TaoBaoKeUatmEventGetResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 通过指定定向招商活动id，获取该活动id下的宝贝信息
func (this *TaoBaoSdk) GetTaoBaoKeUatmEventItem(param TaoBaoKeParam) (TaoBaoFavoriteItemGetResponse, error) {
	method := "taobao.tbk.uatm.event.item.get"
	result := TaoBaoFavoriteItemGetResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 获取淘抢购的数据，淘客商品转淘客链接，非淘客商品输出普通链接
func (this *TaoBaoSdk) GetTaoBaoKeJuTqg(param TaoBaoKeParam) (TaoBaoKeJuTqgGetResponse, error) {
	method := "taobao.tbk.ju.tqg.get"
	result := TaoBaoKeJuTqgGetResponse{}
	err := this.getTaoBaoKe(method, param, &result)
	return result, err
}

// 淘宝客方法调用
func (this *TaoBaoSdk) getTaoBaoKe(method string, param TaoBaoKeParam, result interface{}) error {

	// 参数
	param.Method = method
	param.AppKey = this.AppKey
	param.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	// 参数map
	paramMap, err := this.getParamMap(param)
	if err != nil {
		return err
	}

	// 签名和参数字符串
	sign, paramStr, err := this.getSignature(paramMap)
	if err != nil {
		return err
	}
	paramStr += "&sign=" + sign

	// host
	host := ""
	switch this.Type {
	case "sandbox":
		host = sandboxHttpUrl
	case "online":
		host = onlineHttpUrl
	case "oversea":
		host = overseaHttpUrl
	default:
		return errors.New("请声明应用类型！")
	}

	// 拼接url
	url := host + "?" + paramStr

	// 请求
	err = this.api(url, "get", &result)
	if err != nil {
		return err
	}

	return nil
}

// 将参数结构体转为map
func (this *TaoBaoSdk) getParamMap(param TaoBaoKeParam) (map[string]string, error) {
	paramMap := map[string]string{}

	jsonParam, err := encoding.EncodeJson(param)
	if err != nil {
		return paramMap, err
	}
	err = encoding.DecodeJson(jsonParam, &paramMap)
	if err != nil {
		return paramMap, err
	}

	return paramMap, nil
}

// ajax请求
func (this *TaoBaoSdk) api(url, method string, result interface{}) error {

	// Ajax请求
	var dataByte []byte
	err := util.DefaultAjaxPool.Get(&util.Ajax{
		Method: method,
		Url:    url,
		Header: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded;charset=utf-8",
		},
		ResponseData: &dataByte,
	})
	if err != nil {
		return err
	}
	// fmt.Printf("%+v\n", string(dataByte))

	// 错误结果
	errorResult := TaoBaoKeErrorResponse{}
	err = encoding.DecodeJson(dataByte, &errorResult)
	if err == nil && errorResult.ErrorResponse.Code != 0 {
		return errors.New(string(dataByte))
	}

	// 正常结果
	err = encoding.DecodeJson(dataByte, &result)
	if err != nil {
		return err
	}

	return nil
}

// 根据传入参数，计算签名，拼凑url query字符串
// 参照http://open.taobao.com/docs/doc.htm?spm=a219a.7629140.0.0.rSsjZq&treeId=1&articleId=101617&docType=1#s0
func (this *TaoBaoSdk) getSignature(param map[string]string) (string, string, error) {

	// 除去sign参数和byte[]类型的参数，所有公共参数和业务参数按字母序排序
	// 如：foo=1, bar=2, foo_bar=3, foobar=4排序后的顺序是bar=2, foo=1, foo_bar=3, foobar=4。
	keys, _ := language.ArrayKeyAndValue(param)
	sortKeys := language.ArraySort(keys).([]string)

	// 将排序好的参数名和参数值拼装在一起
	// 由上可得：bar2foo1foo_bar3foobar4。
	sortStr := ""
	paramStr := ""
	newMap := map[string]string{}
	for _, single := range sortKeys {
		v := param[single]
		// 过滤默认值
		if v == "" || v == "0" || v == "false" {
			continue
		}

		newMap[single] = v
		encodeValue, err := encoding.EncodeUrl(v)
		if err != nil {
			return "", "", err
		}
		paramStr += single + "=" + encodeValue + "&"
		sortStr += single + v
	}
	paramStr = strings.TrimRight(paramStr, "&")

	// 摘要算法求签名
	sign := ""
	if v, ok := param["sign_method"]; ok {
		switch v {
		case "md5":
			sign = this.getMd5Signature(sortStr)
		case "hmac":
			sign = this.getHMACMd5Signature(sortStr)
		default:
			return "", "", errors.New("请传入合法摘要算法名！")
		}
	} else {
		return "", "", errors.New("请传入摘要算法名！")
	}

	// 大写
	upperSign := strings.ToUpper(sign)

	return upperSign, paramStr, nil
}

// md5
func (this *TaoBaoSdk) getMd5Signature(sortStr string) string {
	sign := ""

	// MD5算法，则需要在拼装的字符串前后加上app的secret后，再进行摘要
	// 如：md5(secret+bar2foo1foo_bar3foobar4+secret)
	// 将摘要得到的字节流结果使用十六进制表示
	// 如：hex(“helloworld”.getBytes(“utf-8”)) = “68656C6C6F776F726C64”
	sortStrWithSecret := this.AppSecret + sortStr + this.AppSecret
	sign = crypto.CryptoMd5([]byte(sortStrWithSecret))

	return sign
}

// hmacmd5
func (this *TaoBaoSdk) getHMACMd5Signature(sortStr string) string {
	sign := ""

	// HMAC_MD5算法，则需要用app的secret初始化摘要算法后，再进行摘要
	// 如：hmac_md5(bar2foo1foo_bar3foobar4)
	// 将摘要得到的字节流结果使用十六进制表示
	// 如：hex(“helloworld”.getBytes(“utf-8”)) = “68656C6C6F776F726C64”
	sign = crypto.CryptoHMACMd5([]byte(sortStr), []byte(this.AppSecret))

	return sign
}
