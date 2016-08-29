package util

import (
	. "github.com/fishedee/assert"
	"testing"
)

func TestAjaxProxy(t *testing.T) {
	ajaxPool := NewAjaxPool(&AjaxPoolOption{
		Proxy: "http://127.0.0.1:8118",
	})
	var data string
	err := ajaxPool.Get(&Ajax{
		Url:          "http://www.google.com.hk",
		ResponseData: &data,
	})
	AssertEqual(t, err, nil)
	AssertEqual(t, len(data) != 0, true)
}
