package encoding

import (
	"encoding/json"
	. "github.com/fishedee/language"
)

func EncodeJson(data interface{}) ([]byte, error) {
	changeValue := ArrayToMap(data, "json")
	return json.Marshal(changeValue)
}

func DecodeJson(data []byte, value interface{}) error {
	var valueDynamic interface{}
	err := json.Unmarshal(data, &valueDynamic)
	if err != nil {
		return err
	}
	return MapToArray(valueDynamic, value, "json")
}
