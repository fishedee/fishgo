package util

import (
	"compress/flate"
	"compress/gzip"
	. "github.com/fishedee/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentEncodingBaidu(t *testing.T) {
	testCase := []struct {
		encoding string
	}{
		{"gzip"},
	}

	for _, singleTestCase := range testCase {
		header := map[string]string{}
		var data string
		DefaultAjaxPool.Get(&Ajax{
			Url: "http://www.baidu.com",
			Header: map[string]string{
				"Accept-Encoding": singleTestCase.encoding,
			},
			ResponseHeader: &header,
			ResponseData:   &data,
		})
		AssertEqual(t, header["Content-Encoding"], singleTestCase.encoding)
		AssertEqual(t, strings.Contains(data, "<html>"), true)
	}
}

func TestAjaxContentEncoding(t *testing.T) {
	testData := "Hello World"
	testCase := []struct {
		acceptEncoding string
		handler        http.HandlerFunc
	}{
		{"gzip", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Encoding", "gzip")
			writer := gzip.NewWriter(w)
			writer.Write([]byte(testData))
			writer.Close()
		}},
		{"deflate", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Encoding", "deflate")
			writer, err := flate.NewWriter(w, flate.DefaultCompression)
			if err != nil {
				panic(err)
			}
			writer.Write([]byte(testData))
			writer.Close()
		}},
		{"", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(testData))
		}},
	}

	for index, singleTestCase := range testCase {
		ts := httptest.NewServer(singleTestCase.handler)
		defer ts.Close()

		header := map[string]string{}
		if singleTestCase.acceptEncoding != "gzip" {
			header["Accept-Encoding"] = singleTestCase.acceptEncoding
		}
		ajaxPool := NewAjaxPool(&AjaxPoolOption{})
		var data string
		err := ajaxPool.Get(&Ajax{
			Url:          ts.URL,
			Header:       header,
			ResponseData: &data,
		})
		if err != nil {
			panic(err)
		}
		AssertEqual(t, data, testData, index)
	}
}

func GGTestAjaxProxy(t *testing.T) {
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
