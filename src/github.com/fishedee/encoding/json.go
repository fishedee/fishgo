package encoding

import (
	"encoding/json"
	. "github.com/fishedee/language"
)

func EncodeJson(data interface{}) ([]byte, error) {
	changeValue := ArrayMapping(data)
	return json.Marshal(changeValue)
}

func DecodeJson(data []byte, value interface{}) error {
	return json.Unmarshal(data, value)
}
