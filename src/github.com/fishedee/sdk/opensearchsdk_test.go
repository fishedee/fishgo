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

func TestOpenSearchSignature2(t *testing.T) {
	sdk := &OpenSearchSdk{
		Host:   "",
		AppId:  "pcjpOwW9kYVgoOCP",
		AppKey: "mmFRVuxXOlRo4SSRMMy4ukkEulTkcm",
	}
	query := map[string]string{
		"query":            "query=default:'搜 索'",
		"index_name":       "t_bakeweb_recipe",
		"Version":          "v2",
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   "14645766476223910",
		"Timestamp":        "2016-05-30T02:50:47Z",
	}
	targetSign := "HZRs+/BmgI3c5cAh1xfzY5cZ5N4="
	result, err := sdk.getSignature("GET", query)
	assertOpenSearchEqual(t, err, nil)
	assertOpenSearchEqual(t, result["Signature"], targetSign)
	assertOpenSearchEqual(t, result["AccessKeyId"], "pcjpOwW9kYVgoOCP")
}

func TestOpenSearchSearch(t *testing.T) {
	sdk := &OpenSearchSdk{
		Host:   "http://intranet.opensearch-cn-qingdao.aliyuncs.com",
		AppId:  "pcjpOwW9kYVgoOCP",
		AppKey: "mmFRVuxXOlRo4SSRMMy4ukkEulTkcm",
	}

	response, err := sdk.Search(OpenSearchSearchRequest{
		Query: OpenSearchQuery{
			Query: "default:'蛋 糕'",
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
