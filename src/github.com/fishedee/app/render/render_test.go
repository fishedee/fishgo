package render

import (
	. "github.com/fishedee/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderBasic(t *testing.T) {
	testCase := []struct {
		name        string
		data        interface{}
		output      string
		contentType string
	}{
		{"raw", []byte("123"), "123", "text/plain; charset=utf-8"},
		{"text", "456", "456", "text/plain; charset=utf-8"},
		{"json", map[string]string{"a": "123", "b": "456"}, "{\"a\":\"123\",\"b\":\"456\"}\n", "application/json; charset=utf-8"},
		{"html", []interface{}{"index.html", map[string]interface{}{"Name": "Fish"}}, "Hello Fish", "text/html; charset=utf-8"},
		{"html", []interface{}{"index2.html", map[string]interface{}{"User": "Fish"}}, "<html><body><div>Fish</div></body></html>", "text/html; charset=utf-8"},
	}

	renderFactory, err := NewRenderFactory(RenderConfig{TemplateDir: "testdata"})
	if err != nil {
		panic(err)
	}
	for index, singleTestCase := range testCase {
		r, _ := http.NewRequest("GET", "http://www.baidu.com/", nil)
		w := httptest.NewRecorder()
		render := renderFactory.Create(w, r)
		err := render.Format(singleTestCase.name, singleTestCase.data)
		AssertEqual(t, err, nil, index)

		result := w.Result()
		bodyArray, _ := ioutil.ReadAll(result.Body)
		AssertEqual(t, string(bodyArray), singleTestCase.output, index)
		AssertEqual(t, result.Header.Get("Content-Type"), singleTestCase.contentType, index)
	}
}

func TestRenderRedirect(t *testing.T) {
	testCase := []struct {
		data interface{}
		code int
		url  string
	}{
		{"/", 302, "/"},
		{[]interface{}{"http://www.baidu.com/a", 301}, 301, "http://www.baidu.com/a"},
		{[]interface{}{"http://www.qq.com/b", 302}, 302, "http://www.qq.com/b"},
	}
	renderFactory, _ := NewRenderFactory(RenderConfig{})
	for index, singleTestCase := range testCase {
		r, _ := http.NewRequest("GET", "http://www.baidu.com/c", nil)
		w := httptest.NewRecorder()
		render := renderFactory.Create(w, r)
		err := render.Format("redirect", singleTestCase.data)
		AssertEqual(t, err, nil, index)

		result := w.Result()
		location, _ := result.Location()
		AssertEqual(t, location.String(), singleTestCase.url, index)
		AssertEqual(t, result.StatusCode, singleTestCase.code, index)
	}
}
