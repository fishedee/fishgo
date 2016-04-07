package sdk

import (
	. "github.com/fishedee/util"
)

type OpenSearchSdk struct {
	AppId  string
	AppKey string
}

type OpenSearchCommon struct {
	Version          string `json:"Version"`
	AccessKeyId      string `json:"AccessKeyId"`
	Signature        string `json:"Signature"`
	SignatureMethod  string `json:"SignatureMethod"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	SignatureNonce   string `json:"SignatureNonce"`
}

type OpenSearchOption struct {
	Query            OpenSearchQuery   `json:"query"`
	IndexName        []string          `json:"index_name"`
	FetchFields      []string          `json:"fetch_fields"`
	Qp               []string          `json:"qp"`
	Disable          string            `json:"disable"`
	FirstFormulaName string            `json:"first_formula_name"`
	FormulaName      string            `json:"formula_name"`
	Summary          OpenSearchSummary `json:"summary"`
}

type OpenSearchQuery struct {
	Config    OpenSearchConfig   `json:"config"`
	Query     map[string]string  `json:"query"`
	Sort      []string           `json:"sort"`
	Filter    map[string]string  `json:"filter"`
	Aggregate map[string]string  `json:"aggregate"`
	Distinct  OpenSearchDistinct `json:"distinct"`
	Kvpairs   map[string]string  `json:"kvpairs"`
}

type OpenSearchConfig struct {
	Start      int    `json:"start"`
	Hit        int    `json:"hit"`
	Format     string `json:"format"`
	RerankSize int    `json:"rerank_size"`
}

type OpenSearchDistinct struct {
	DistKey        string  `json:"dist_key"`
	DistTimes      int     `json:"dist_times"`
	DistCount      int     `json:"dist_count"`
	Reserved       bool    `json:"reserved"`
	UpdateTotalHit bool    `json:"update_total_hit"`
	DistFilter     string  `json:"dist_filter"`
	Grade          float64 `json:"grade"`
}

type OpenSearchSummary struct {
	SummaryField          string `json:"summary_field"`
	SummaryElement        string `json:"summary_element"`
	SummaryEllipsis       string `json:"summary_ellipsis"`
	SummarySnipped        int    `json:"summary_snipped"`
	SummaryLen            string `json:"summary_len"`
	SummaryElementPrefix  string `json:"summary_element_prefix"`
	SummaryElementPostfix string `json:"summary_element_postfix"`
}

type OpenSearchResult struct {
	Status    string                `json:"status"`
	RequestId string                `json:"request_id"`
	Result    OpenSearchInnerResult `json:"result"`
	Errors    []OpenSearchError     `json:"errors"`
	Tracer    string                `json:"tracer"`
}

type OpenSearchInnerResult struct {
	Searchtime float64          `json:"searchtime"`
	Total      int              `json:"total"`
	Num        int              `json:"num"`
	Viewtotal  int              `json:"viewtotal"`
	Items      []OpenSearchItem `json:"items"`
	Facet      []interface{}    `json:"facet"`
}

type OpenSearchItem struct {
	Fields        OpenSearchField `json:"fields"`
	VariableValue interface{}     `json:"variableValue"`
}

type OpenSearchField struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	Title           string `json:"title"`
	Body            string `json:"body"`
	Url             string `json:"url"`
	CreateTimestamp string `json:"create_timestamp"`
	IndexName       string `json:"index_name"`
}

type OpenSearchError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

/*
* 签名
 */
func (this *OpenSearchSdk) getSignature(url string) string {
	//参数排序
	//url = this.sortParams(url)

	//名称和值分别url编码

	//连接

	//hmac & base64

	//url+signature
	return ""
}

/**
* 搜索
 */
func (this *OpenSearchSdk) GetByKeyword(option OpenSearchOption) (OpenSearchResult, error) {
	//参数
	url := ""
	body := ""

	//请求
	var result []byte
	err := DefaultAjaxPool.Get(&Ajax{
		Url:          url,
		Data:         body,
		DataType:     "json",
		ResponseData: &result,
	})
	if err != nil {
		return OpenSearchResult{}, err
	}

	//结果
	return OpenSearchResult{}, nil
}

func (this *OpenSearchSdk) getRequestUrl(option OpenSearchOption) string {
	example := `
http://$host/search?
index_name=bbs&
query=config=start:0,hit:10,format=fulljson&&query=default:'的'&&filter=create_timestamp>1423000000&&sort=+type;-RANK&
fetch_fields=id;title;body;url;type;create_timestamp&
first_formula_name=first_bbs&
formula_name=second_bbs&
summary=summary_snipped:1,summary_field:title,summary_element:high,summary_len:32,summary_ellipsis:...;summary_snipped:2,summary_field:body,summary_element:high,summary_len:60,summary_ellipsis:...
`
	return example
}
