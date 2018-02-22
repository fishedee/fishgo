package render

import (
	"errors"
	"net/http"
)

type RawFormatter struct {
}

func (this *RawFormatter) Name() string {
	return "raw"
}

func (this *RawFormatter) Format(w http.ResponseWriter, r *http.Request, data interface{}) error {
	var result []byte
	if dataByte, isOk := data.([]byte); isOk == true {
		result = dataByte
	} else {
		return errors.New("invalid data type for raw formatter")
	}
	_, err := w.Write(result)
	if err != nil {
		return err
	}
	return nil
}

func NewRawFormatter() (*RawFormatter, error) {
	return &RawFormatter{}, nil
}
