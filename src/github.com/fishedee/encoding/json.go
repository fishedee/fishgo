package encoding

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"
)
func nameMapper(name string)(string){
	return strings.ToLower(name[0:1])+name[1:]
}
func changeToValue(data interface{})(interface{}){
	var result interface{}
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	if data == nil{
		result = data
	}else if dataType.Kind() == reflect.Struct{
		if dataType == reflect.TypeOf(time.Time{}){
			timeValue := data.(time.Time)
			result = timeValue.Format("2006-01-02 15:04:05")
		}else{
			resultMap := map[string]interface{}{}
			for i := 0 ; i != dataValue.NumField() ; i++{
				singleDataType := dataType.Field(i)
				singleDataValue := dataValue.Field(i)
				singleName := nameMapper(singleDataType.Name)
				resultMap[singleName] = changeToValue(singleDataValue.Interface())
			}
			result = resultMap
		}
	}else if( dataType.Kind() == reflect.Slice ){
		resultSlice := []interface{}{}
		for i := 0 ; i != dataValue.Len() ; i++{
			singleDataValue := dataValue.Index(i)
			resultSlice = append(resultSlice,changeToValue(singleDataValue.Interface()) )
		}
		result = resultSlice
	}else{
		result = data
	}
	return result
}

func EncodeJson(data interface{})([]byte,error){
	return json.Marshal(changeToValue(data))
}

func DecodeJson(data []byte,value interface{})(error){
	return json.Unmarshal(data,value)
}