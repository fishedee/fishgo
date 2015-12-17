package encoding

import (
	. "github.com/fishedee/language"
	"encoding/json"
)

func EncodeJson(data interface{})([]byte,error){
	changeValue := StructToMap(data)
	return json.Marshal(changeValue)
}

func DecodeJson(data []byte,value interface{})(error){
	return json.Unmarshal(data,value)
}