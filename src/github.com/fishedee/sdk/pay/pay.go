package gopay

import (
	"errors"
	"github.com/fishedee/sdk/pay/client"
	"github.com/fishedee/sdk/pay/common"
	"github.com/fishedee/sdk/pay/constant"
	//"github.com/fishedee/sdk/pay/util"
	//"strconv"
)

// 获取支付接口
func Pay(charge *common.Charge) (map[string]string, error) {
	err := checkCharge(charge)
	if err != nil {
		//log.Error(err, charge)
		panic(err)
	}

	ct := getPayClient(charge.PayMethod)
	re, err := ct.Pay(charge)
	if err != nil {
		//log.Error("支付失败:", err, charge)
		panic(err)
	}
	return re, err
}

// 验证内容
func checkCharge(charge *common.Charge) error {
	if charge.PayMethod < 0 {
		return errors.New("payMethod less than 0")
	}
	if charge.MoneyFee < 0 {
		return errors.New("totalFee less than 0")
	}
	return nil
}

// getPayClient 得到需要支付的客户端
func getPayClient(payMethod int64) common.PayClient {
	//如果使用余额支付
	switch payMethod {
	 case constant.ALI_WEB:
	 	return client.DefaultAliWebClient()
	 case constant.ALI_APP:
	 	return client.DefaultAliAppClient()
	case constant.WECHAT_WEB:
		return client.DefaultWechatWebClient()
	case constant.WECHAT_APP:
		return client.DefaultWechatAppClient()
	}
	return nil
}
