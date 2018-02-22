package render

import (
	"errors"
	"net/http"
)

type RedirectFormatter struct {
}

func (this *RedirectFormatter) Name() string {
	return "redirect"
}

func (this *RedirectFormatter) Format(w http.ResponseWriter, r *http.Request, data interface{}) error {
	var url string
	var code int

	if dataString, isOk := data.(string); isOk {
		url = dataString
		code = 302
	} else if dataArray, isOk := data.([]interface{}); isOk && len(dataArray) == 2 {
		url = dataArray[0].(string)
		code = dataArray[1].(int)
	} else {
		return errors.New("invalid data type for redirect formatter")
	}
	http.Redirect(w, r, url, code)
	return nil
}

func NewRedirectFormatter() (*RedirectFormatter, error) {
	return &RedirectFormatter{}, nil
}
