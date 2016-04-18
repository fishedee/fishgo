package sdk

import (
	"reflect"
	"testing"
)

func assertOpenSearchEqual(t *testing.T, left interface{}, right interface{}) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("assert fail: %+v != %+v", left, right)
	}
}

func TestOpenSearchSignature(t *testing.T) {
	sdk := &OpenSearchSdk{
		Host:   "",
		AppId:  "testid",
		AppKey: "testsecret",
	}
	query := map[string]string{
		"Version":          "v2",
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   "14053016951271226",
		"Timestamp":        "2014-07-14T01:34:55Z",
		"query":            "config=format:json,start:0,hit:20&&query=default:'的'",
		"index_name":       "ut_3885312",
		"format":           "json",
		"fetch_fields":     "title;gmt_modified",
	}
	targetSign := "/GWWQkztlp/9Qg7rry2DuCSfKUQ="
	result, err := sdk.getSignature("GET", query)
	assertOpenSearchEqual(t, err, nil)
	assertOpenSearchEqual(t, result["Signature"], targetSign)
	assertOpenSearchEqual(t, result["AccessKeyId"], "testid")
}

func TestOpenSearchSearch(t *testing.T) {
	sdk := &OpenSearchSdk{
		Host:   "http://intranet.opensearch-cn-qingdao.aliyuncs.com",
		AppId:  "pcjpOwW9kYVgoOCP",
		AppKey: "mmFRVuxXOlRo4SSRMMy4ukkEulTkcm",
	}

	response, err := sdk.Search(OpenSearchSearchRequest{
		Query: OpenSearchQuery{
			Query: "default:'蛋糕'",
		},
		IndexName: "t_bakeweb_recipe",
	})
	assertOpenSearchEqual(t, err, nil)
	assertOpenSearchEqual(t, response.Total != 0, true)
	assertOpenSearchEqual(t, len(response.Items), response.Num)
	assertOpenSearchEqual(t, len(response.Items) != 0, true)
	for _, singleData := range response.Items {
		assertOpenSearchEqual(t, len(singleData) != 0, true)
	}
}
