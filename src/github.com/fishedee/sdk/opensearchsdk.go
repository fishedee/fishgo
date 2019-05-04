package sdk

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	. "github.com/fishedee/crypto"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	"strings"
	"time"
)

type OpenSearchSdk struct {
	Host   string
	AppId  string
	AppKey string
}

type OpenSearchQuery struct {
	Config    string `query:"config,omitempty"`
	Query     string `query:"query,omitempty"`
	Sort      string `query:"sort,omitempty"`
	Filter    string `query:"filter,omitempty"`
	Aggregate string `query:"aggregate,omitempty"`
	Distinct  string `query:"distinct,omitempty"`
	Kvpairs   string `query:"kvpairs,omitempty"`
}

type OpenSearchSearchRequest struct {
	Query            OpenSearchQuery `query:"query"`
	IndexName        string          `query:"index_name"`
	FetchFields      string          `query:"fetch_fields,omitempty"`
	Qp               string          `query:"qp,omitempty"`
	Disable          string          `query:"disable,omitempty"`
	FirstFormulaName string          `query:"first_formula_name,omitempty"`
	FormulaName      string          `query:"formula_name,omitempty"`
}

/*
type OpenSearchSearchResponseItem struct {
	Fields        map[string]interface{} `json:"fields"`
	VariableValue `json:"variableValue"`
}
*/

type OpenSearchSearchResponse struct {
	SearchTime float64                  `json:"searchtime"`
	Total      int                      `json:"total"`
	Num        int                      `json:"num"`
	ViewTotal  int                      `json:"viewtotal"`
	Items      []map[string]interface{} `json:"items"`
}

type OpenSearchError struct {
	Code      int
	Message   string
	RequestId string
}

func (this *OpenSearchError) GetCode() int {
	return this.Code
}

func (this *OpenSearchError) GetMsg() string {
	return this.Message
}

func (this *OpenSearchError) Error() string {
	return fmt.Sprintf("错误码为：%v，错误描述为：%v, 请求Id: %v", this.Code, this.Message, this.RequestId)
}

func (this *OpenSearchSdk) encodeUrl(input string) (string, error) {
	output, err := EncodeUrl(input)
	if err != nil {
		return "", err
	}
	return strings.Replace(output, "+", "%20", -1), nil
}

func (this *OpenSearchSdk) getSignature(method string, query map[string]string) (map[string]string, error) {
	//参数格式化
	newQuery := map[string]string{}
	for key, value := range query {
		newQuery[key] = value
	}
	newQuery["AccessKeyId"] = this.AppId
	queryKeyInterface, _ := ArrayKeyAndValue(newQuery)
	queryKey := ArraySort(queryKeyInterface).([]string)
	stringToSignArray := []string{}
	for _, singleQueryKey := range queryKey {
		singleQueryValue := newQuery[singleQueryKey]
		singleQueryKeyEncode, err := this.encodeUrl(singleQueryKey)
		if err != nil {
			return nil, err
		}
		singleQueryValueEncode, err := this.encodeUrl(singleQueryValue)
		if err != nil {
			return nil, err
		}
		stringToSignArray = append(stringToSignArray, singleQueryKeyEncode+"="+singleQueryValueEncode)
	}
	stringToSign := strings.Join(stringToSignArray, "&")

	//字符串签名
	stringToSignEncode, err := this.encodeUrl(stringToSign)
	if err != nil {
		return nil, err
	}
	stringToSign = method + "&%2F&" + stringToSignEncode

	//hmac与base64编码
	hmacKey := this.AppKey + "&"
	mac := hmac.New(sha1.New, []byte(hmacKey))
	mac.Write([]byte(stringToSign))
	signHmac := mac.Sum(nil)
	signBase64 := base64.StdEncoding.EncodeToString(signHmac)
	newQuery["Signature"] = string(signBase64)
	return newQuery, nil
}

func (this *OpenSearchSdk) api(method string, url string, query map[string]string, bodyResult interface{}) error {
	//初始化基础参数
	newQuery := map[string]string{}
	for key, value := range query {
		newQuery[key] = value
	}
	now := time.Now().In(time.UTC)
	nowStr := now.Format("2006-01-02T15:04:05Z")
	newQuery["Version"] = "v2"
	newQuery["Timestamp"] = nowStr
	newQuery["SignatureMethod"] = "HMAC-SHA1"
	newQuery["SignatureVersion"] = "1.0"
	newQuery["SignatureNonce"] = CryptoRand(17)

	//生成签名
	newQuery, err := this.getSignature(method, newQuery)
	if err != nil {
		return err
	}

	//调用
	var result struct {
		Status    string      `json:"status"`
		RequestId string      `json:"request_id"`
		Result    interface{} `json:"result"`
		Errors    []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	result.Result = bodyResult
	ajax := Ajax{
		Url:              this.Host + url,
		Data:             newQuery,
		DataType:         "url",
		ResponseData:     &result,
		ResponseDataType: "json",
	}
	if method == "GET" {
		err = DefaultAjaxPool.Get(&ajax)
	} else {
		err = DefaultAjaxPool.Post(&ajax)
	}
	if err != nil {
		return err
	}
	if result.Status != "OK" {
		return &OpenSearchError{result.Errors[0].Code, result.Errors[0].Message, result.RequestId}
	}
	return nil
}

func (this *OpenSearchSdk) combineQuery(request interface{}) map[string]string {
	data := ArrayToMap(request, "query")
	dataMap := data.(map[string]interface{})
	result := map[string]string{}
	for singleKey, singleValue := range dataMap {
		var singleResult string
		singleMapValue, isOk := singleValue.(map[string]interface{})
		if isOk {
			singleMapValueArray := []string{}
			for key, value := range singleMapValue {
				singleMapValueArray = append(singleMapValueArray, key+"="+fmt.Sprintf("%v", value))
			}
			singleResult = strings.Join(singleMapValueArray, "&&")
		} else {
			singleResult = fmt.Sprintf("%v", singleValue)
		}
		result[singleKey] = singleResult
	}
	return result
}

func (this *OpenSearchSdk) Search(request OpenSearchSearchRequest) (OpenSearchSearchResponse, error) {
	var response OpenSearchSearchResponse
	requestQuery := this.combineQuery(request)
	err := this.api("GET", "/search", requestQuery, &response)
	if err != nil {
		return OpenSearchSearchResponse{}, err
	}
	return response, nil
}

func (this *OpenSearchSdk) Doc() {

}
