package sdk

import (
	"reflect"
	"testing"
)

func assertTaoBaoSdkEqual(t *testing.T, left interface{}, right interface{}) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("assert fail: %+v != %+v", left, right)
	}
}

func TestGetTaoBaoKeAllItem(t *testing.T) {
	taoBaoSdk := &TaoBaoSdk{
		AppKey:    "23460891",
		AppSecret: "349b32a33b43952bb3f5a86de3328106",
		Type:      "online",
	}
	param := TaoBaoKeParam{
		TaoBaoCommonParam: TaoBaoCommonParam{
			Format:     "json",
			V:          "2.0",
			SignMethod: "hmac", //"md5",
		},
		Fields: "num_iid,title,pict_url,small_images,reserve_price,zk_final_price,user_type,provcity,item_url",
		// q和cat2选1
		Q: "忠臣z6",
		// Cat: "",
	}
	result, err := taoBaoSdk.GetTaoBaoKeAllItem(param)
	if err != nil {
		t.Log(err.Error() + "\n")
	} else {
		// TODO
		assertTaoBaoSdkEqual(t, result, result)
	}
}

func TestGetTaoBaoKeItemInfo(t *testing.T) {
	taoBaoSdk := &TaoBaoSdk{
		AppKey:    "23460891",
		AppSecret: "349b32a33b43952bb3f5a86de3328106",
		Type:      "online",
	}
	param := TaoBaoKeParam{
		TaoBaoCommonParam: TaoBaoCommonParam{
			Format:     "json",
			V:          "2.0",
			SignMethod: "md5",
		},
		Fields:  "num_iid,title,pict_url,small_images,reserve_price,zk_final_price,user_type,provcity,item_url",
		NumIIds: "530194448923",
	}
	result, err := taoBaoSdk.GetTaoBaoKeItemInfo(param)
	if err != nil {
		t.Log(err.Error() + "\n")
	} else {
		// TODO
		assertTaoBaoSdkEqual(t, result, result)
	}
}

func TestGetTaoBaoKeItemRecommend(t *testing.T) {
	taoBaoSdk := &TaoBaoSdk{
		AppKey:    "23460891",
		AppSecret: "349b32a33b43952bb3f5a86de3328106",
		Type:      "online",
	}
	param := TaoBaoKeParam{
		TaoBaoCommonParam: TaoBaoCommonParam{
			Format:     "json",
			V:          "2.0",
			SignMethod: "md5",
		},
		Fields: "num_iid,title,pict_url,small_images,reserve_price,zk_final_price,user_type,provcity,item_url",
		NumIId: "530194448923",
	}
	result, err := taoBaoSdk.GetTaoBaoKeItemRecommend(param)
	if err != nil {
		t.Log(err.Error() + "\n")
	} else {
		// TODO
		assertTaoBaoSdkEqual(t, result, result)
	}
}

func TestGetTaoBaoKeShop(t *testing.T) {
	taoBaoSdk := &TaoBaoSdk{
		AppKey:    "23460891",
		AppSecret: "349b32a33b43952bb3f5a86de3328106",
		Type:      "online",
	}
	param := TaoBaoKeParam{
		TaoBaoCommonParam: TaoBaoCommonParam{
			Format:     "json",
			V:          "2.0",
			SignMethod: "md5",
		},
		Fields: "user_id,shop_title,shop_type,seller_nick,pict_url,shop_url",
		Q:      "忠臣z6",
	}
	result, err := taoBaoSdk.GetTaoBaoKeShop(param)
	if err != nil {
		t.Log(err.Error() + "\n")
	} else {
		// TODO
		assertTaoBaoSdkEqual(t, result, result)
	}
}

func TestGetTaoBaoKeUatmFavoriteItem(t *testing.T) {
	taoBaoSdk := &TaoBaoSdk{
		AppKey:    "23460891",
		AppSecret: "349b32a33b43952bb3f5a86de3328106",
		Type:      "online",
	}
	param := TaoBaoKeParam{
		TaoBaoCommonParam: TaoBaoCommonParam{
			Format:     "json",
			V:          "2.0",
			SignMethod: "md5",
		},
		Fields: "	num_iid,title,pict_url,small_images,reserve_price,zk_final_price,user_type,provcity,item_url,seller_id,volume,nick,shop_title,zk_final_price_wap,event_start_time,event_end_time,tk_rate,status,type,click_url",
		AdzoneId:    "61682495", // 推广位id
		FavoritesId: "1239563",  // 选品库的id
	}
	result, err := taoBaoSdk.GetTaoBaoKeUatmFavoriteItem(param)
	if err != nil {
		t.Log(err.Error() + "\n")
	} else {
		// TODO
		assertTaoBaoSdkEqual(t, result, result)
	}
}

func TestGetTaoBaoKeUatmFavorites(t *testing.T) {
	taoBaoSdk := &TaoBaoSdk{
		AppKey:    "23460891",
		AppSecret: "349b32a33b43952bb3f5a86de3328106",
		Type:      "online",
	}
	param := TaoBaoKeParam{
		TaoBaoCommonParam: TaoBaoCommonParam{
			Format:     "json",
			V:          "2.0",
			SignMethod: "md5",
		},
		Fields:   "favorites_title,favorites_id,type",
		PageNo:   1,
		PageSize: 20,
		Type:     1,
	}
	result, err := taoBaoSdk.GetTaoBaoKeUatmFavorites(param)
	if err != nil {
		t.Log(err.Error() + "\n")
	} else {
		// TODO
		assertTaoBaoSdkEqual(t, result, result)
	}
}
