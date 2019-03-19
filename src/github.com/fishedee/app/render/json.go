package render

import (
	. "github.com/fishedee/encoding"
	"net/http"
)

type JsonFormatter struct {
}

func (this *JsonFormatter) Name() string {
	return "json"
}

func (this *JsonFormatter) Format(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Cache-Control", "private, no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Add("Cache-Control", "post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	dataByte, err := EncodeJson(data)
	if err != nil {
		return err
	}
	w.Write(dataByte)
	return nil
}

func NewJsonFormatter() (*JsonFormatter, error) {
	return &JsonFormatter{}, nil
}
