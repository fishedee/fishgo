package weixin

type JsConfig struct {
	AppId     string
	Noncestr  string
	Timestamp string
	Signature string
}

type JsSignature struct {
	JsApiTicket string
	Noncestr    string
	Timestamp   string
	Url         string
}
