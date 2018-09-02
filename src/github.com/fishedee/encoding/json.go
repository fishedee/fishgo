package encoding

import (
	"bytes"
	"encoding/json"
	. "github.com/fishedee/language"
)

func EncodeJson(data interface{}) ([]byte, error) {
	changeValue := ArrayToMap(data, "json")
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")
	err := encoder.Encode(changeValue)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func DecodeJson(data []byte, value interface{}) error {
	var valueDynamic interface{}
	err := json.Unmarshal(data, &valueDynamic)
	if err != nil {
		return err
	}
	return MapToArray(valueDynamic, value, "json")
}
