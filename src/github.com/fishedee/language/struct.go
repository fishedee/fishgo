package language

import (
	"reflect"
	"strings"
	"errors"
	"time"
)

func nameMapper(name string)(string){
	return strings.ToLower(name[0:1])+name[1:]
}

func isPublic(name string)(bool){
	first := name[0:1]
	return first >= "A" && first <= "Z"
}

func combileMap(result map[string]interface{},singleResultMap interface{})(error){
	singleResultMapMap,ok := singleResultMap.(map[string]interface{})
	if ok == false{
		return errors.New("Anonymous field is not a struct")
	}
	for key,value := range singleResultMapMap{
		result[key] = value
	}
	return nil
}

func structToMap(data interface{})(interface{},error){
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
				if isPublic( singleDataType.Name) == false{
					continue
				}
				singleName := nameMapper(singleDataType.Name)
				singleResultMap,err := structToMap(singleDataValue.Interface())
				if err != nil{
					return nil,err
				}
				if singleDataType.Anonymous == false{
					resultMap[singleName] = singleResultMap
				}else{
					err := combileMap(resultMap,singleResultMap)
					if err != nil{
						return nil,err
					}
				}
			}
			result = resultMap
		}
	}else if( dataType.Kind() == reflect.Slice ){
		resultSlice := []interface{}{}
		for i := 0 ; i != dataValue.Len() ; i++{
			singleDataValue := dataValue.Index(i)
			singleDataResult,err := structToMap(singleDataValue.Interface())
			if err != nil{
				return nil,err
			}
			resultSlice = append(resultSlice,singleDataResult)
		}
		result = resultSlice
	}else{
		result = data
	}
	return result,nil
}

func StructToMap(data interface{})(interface{}){
	result,err := structToMap(data)
	if err != nil{
		panic(err)
	}
	return result
}
