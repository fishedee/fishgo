package render

import (
	"errors"
	"net/http"
)

type TextFormatter struct {
}

func (this *TextFormatter) Name() string {
	return "text"
}

func (this *TextFormatter) Format(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	var result string
	if dataString, isOk := data.(string); isOk {
		result = dataString
	} else {
		return errors.New("invalid data type for text formatter")
	}
	_, err := w.Write([]byte(result))
	if err != nil {
		return err
	}
	return nil
}

func NewTextFormatter() (*TextFormatter, error) {
	return &TextFormatter{}, nil
}
