package web

import (
	"net/http"
	"testing"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Testing        *testing.T
}

func NewContext(request *http.Request, response http.ResponseWriter, t *testing.T) Context {
	return Context{
		Request:        request,
		ResponseWriter: response,
		Testing:        t,
	}
}
