package middleware

import (
	"net/http"
)

type fakeWriter struct {
	header http.Header
	result string
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
	this.result += string(data)
	return len(data), nil
}

func (this *fakeWriter) Read() string {
	return this.result
}
