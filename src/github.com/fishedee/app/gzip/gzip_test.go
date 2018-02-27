package gzip

import (
	. "github.com/fishedee/assert"
	. "github.com/fishedee/compress"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGzipMinSize(t *testing.T) {
	gzip, err := NewGzip(GzipConfig{
		MinSize: 5,
		Level:   5,
	})
	if err != nil {
		panic(err)
	}

	testCase := []struct {
		data       string
		isCompress bool
		handler    func(w http.ResponseWriter)
	}{
		{"1", false, func(w http.ResponseWriter) {
			w.Write([]byte("1"))
		}},
		{"1234", false, func(w http.ResponseWriter) {
			w.Write([]byte("1234"))
		}},
		{"12345", false, func(w http.ResponseWriter) {
			w.Write([]byte("12345"))
		}},
		{"123456", true, func(w http.ResponseWriter) {
			w.Write([]byte("123456"))
		}},
		{"1234", false, func(w http.ResponseWriter) {
			w.Write([]byte("123"))
			w.Write([]byte("4"))
		}},
		{"12345", false, func(w http.ResponseWriter) {
			w.Write([]byte("1234"))
			w.Write([]byte("5"))
		}},
		{"123456", true, func(w http.ResponseWriter) {
			w.Write([]byte("123456"))
		}},
		{"123456", true, func(w http.ResponseWriter) {
			w.Write([]byte("12345"))
			w.Write([]byte("6"))
		}},
		{"123456", true, func(w http.ResponseWriter) {
			w.Write([]byte("1234"))
			w.Write([]byte("56"))
		}},
		{"123456789", true, func(w http.ResponseWriter) {
			w.Write([]byte("1234"))
			w.Write([]byte("56"))
			w.Write([]byte("7"))
			w.Write([]byte("8"))
			w.Write([]byte("9"))
		}},
		{"123456789", true, func(w http.ResponseWriter) {
			w.Write([]byte("123456789"))
		}},
		{"123456789", true, func(w http.ResponseWriter) {
			w.Write([]byte("1234"))
			w.Write([]byte("5"))
			w.Write([]byte("6"))
			w.Write([]byte("7"))
			w.Write([]byte("8"))
			w.Write([]byte("9"))
		}},
	}

	for _, singleTestCase := range testCase {
		r, _ := http.NewRequest("GET", "/", nil)
		r.Header.Set("Accept-Encoding", "gzip;q=1.0, identity; q=0.5, *;q=0")
		w := httptest.NewRecorder()
		gzip.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
			singleTestCase.handler(w)
		})
		data, err := ioutil.ReadAll(w.Result().Body)
		if err != nil {
			panic(err)
		}
		if singleTestCase.isCompress == true {
			AssertEqual(t, w.Result().Header.Get("Content-Encoding"), "gzip")
			AssertEqual(t, w.Result().Header.Get("Content-Length"), "")
			data, err = DecompressGzip(data)
			if err != nil {
				panic(err)
			}
		}
		AssertEqual(t, data, []byte(singleTestCase.data))
	}
}

func TestGzipContentType(t *testing.T) {
	testCase := []struct {
		newContentType   []string
		buildContentType string
		isCompress       bool
	}{
		{nil, "text/html", true},
		{[]string{"text/html"}, "text/html", true},
		{[]string{"application/json;charset=utf-8", "text/html;charset=utf-8"}, "text/html", true},
		{[]string{"application/json;charset=utf-8", "text/html"}, "text/html;charset=utf-8", true},
		{[]string{"application/json;charset=utf-8", "text/html"}, "", false},
		{[]string{"application/json;charset=utf-8", "text/html"}, "text/plain", false},
	}

	for _, singleTestCase := range testCase {
		gzip, err := NewGzip(GzipConfig{
			MinSize:     1,
			Level:       5,
			ContentType: singleTestCase.newContentType,
		})
		if err != nil {
			panic(err)
		}
		r, _ := http.NewRequest("GET", "/", nil)
		r.Header.Set("Accept-Encoding", "gzip;q=1.0, identity; q=0.5, *;q=0")
		w := httptest.NewRecorder()
		gzip.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", singleTestCase.buildContentType)
			w.Write([]byte("12"))
		})
		data, err := ioutil.ReadAll(w.Result().Body)
		if err != nil {
			panic(err)
		}
		if singleTestCase.isCompress {
			AssertEqual(t, w.Result().Header.Get("Content-Encoding"), "gzip")
			AssertEqual(t, w.Result().Header.Get("Content-Length"), "")
			data, err = DecompressGzip(data)
			if err != nil {
				panic(err)
			}
		}
		AssertEqual(t, data, []byte("12"))
	}
}

func TestGzipNoCompress(t *testing.T) {
	testCase := []struct {
		hasAcceptEncoding bool
		writer            func(w http.ResponseWriter)
	}{
		{true, func(w http.ResponseWriter) {
			w.Write([]byte("12"))
			w.Header().Set("Content-Type", "text/html;charset=utf-8")
		}},
		{true, func(w http.ResponseWriter) {
			w.Header().Set("Content-Encoding", "deflate")
			w.Write([]byte("12"))
		}},
		{false, func(w http.ResponseWriter) {
			w.Header().Set("Content-Type", "text/html;charset=utf-8")
			w.Write([]byte("12"))
		}},
	}

	for _, singleTestCase := range testCase {
		gzip, err := NewGzip(GzipConfig{
			MinSize:     1,
			Level:       5,
			ContentType: []string{"text/html"},
		})
		if err != nil {
			panic(err)
		}
		r, _ := http.NewRequest("GET", "/", nil)
		if singleTestCase.hasAcceptEncoding {
			r.Header.Set("Accept-Encoding", "gzip;q=1.0, identity; q=0.5, *;q=0")
		}
		w := httptest.NewRecorder()
		gzip.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
			singleTestCase.writer(w)
		})
		data, err := ioutil.ReadAll(w.Result().Body)
		if err != nil {
			panic(err)
		}
		AssertEqual(t, data, []byte("12"))
	}
}

type fakeWriter struct {
	header http.Header
}

func (this *fakeWriter) Header() http.Header {
	if this.header == nil {
		this.header = http.Header{}
	}
	return this.header
}

func (this *fakeWriter) WriteHeader(status int) {

}

func (this *fakeWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func BenchmarkGzip(b *testing.B) {
	data := ""
	for i := 0; i != 1024; i++ {
		data += "Hello World" + strconv.Itoa(i)
	}
	dataByte := []byte(data)
	gzip, err := NewGzip(GzipConfig{
		Level:       9,
		MinSize:     1024,
		ContentType: nil,
	})
	if err != nil {
		panic(err)
	}
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip;q=1.0, identity; q=0.5, *;q=0")
	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		w := &fakeWriter{}
		gzip.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Write(dataByte)
		})
	}
}
