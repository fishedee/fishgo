package render

import (
	"encoding/json"
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

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")
	err := encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func NewJsonFormatter() (*JsonFormatter, error) {
	return &JsonFormatter{}, nil
}
