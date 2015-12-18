package util

/*
import (
	"net/http"
	"net/url"
	"time"
)

type Ajax struct {
	method string
	url    string
}

func NewAjax() *Ajax {
	return &Ajax{
		method: "get",
		url:    "",
	}
}

func (this Ajax) Get(string url) Ajax {
	this.method = "get"
	this.url = url
	return this
}

func (this Ajax) Post(string url) *Ajax {
	this.method = "post"
	this.url = url
	return this
}

func (this Ajax) SetHeader(header interface{}) Ajax {
	return this
}

func (this Ajax) SetUrl(url string) Ajax {
	this.url = url
	return this
}

func (this Ajax) SetData(dataType string, data interface{}) Ajax {
	return this
}

func (this Ajax) SetCookie(cookie interface{}) Ajax {

}

func (this Ajax) SetTimeout(timeout time.Duration) Ajax {

}

func (this Ajax) GetData(dataType string, data interface{}) Ajax {

}

func (this Ajax) GetHeader(header interface{}) Ajax {

}

func (this Ajax) GetCookie(cookie interface{}) Ajax {

}

func (this Ajax) Send() error {

}

func init() {
	var data struct {
		mm string
		xx ii
	}
	var result struct {
		uu string
	}

		err := NewAjax()
			.SetUrl("http://www.hongbeibang.com/user/get")
			.SetData('url',)
			.Send()

		err = NewAjax()
			.SetMethod("post")
			.SetUrl("post","http://www.hongbeibang.com/user/get")
			.SetData('url',data)
			.SetResponseData('url',&result)
			.Send()
		if err != nil{
			panic(err)
		}
}
*/
