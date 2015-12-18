package language

import (
	"reflect"
	"fmt"
	"strings"
	"errors"
	"time"
)

func ArrayColumnKey(data interface{},name string)(interface{}){
	//提取信息
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Slice{
		panic("array column should be a slice")
	}
	dataElemType := dataType.Elem()
	if dataElemType.Kind() != reflect.Struct{
		panic("array column element should be a struct")
	}
	dataElemFieldType,ok := dataElemType.FieldByName(name)
	if !ok{
		panic("dataElemFieldType has not filed "+name)
	}

	//整合slice
	resultType := reflect.SliceOf(dataElemFieldType.Type)
	result := reflect.MakeSlice(resultType,0,0)
	dataValue := reflect.ValueOf(data)
	for i := 0 ; i != dataValue.Len() ; i++{
		singleDataValue := dataValue.Index(i)
		singleDataFieldValue := singleDataValue.FieldByName(name)
		result = reflect.Append(result,singleDataFieldValue)
	}
	return result.Interface()
}

func ArrayColumnMap(data interface{},name string)(interface{}){
	//提取信息
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Slice{
		panic("array column should be a slice")
	}
	dataElemType := dataType.Elem()
	if dataElemType.Kind() != reflect.Struct{
		panic("array column element should be a struct")
	}
	dataElemFieldType,ok := dataElemType.FieldByName(name)
	if !ok{
		panic("dataElemFieldType has not filed "+name)
	}

	//整合map
	resultType := reflect.MapOf(dataElemFieldType.Type,dataElemType)
	result := reflect.MakeMap(resultType)
	dataValue := reflect.ValueOf(data)
	for i := 0 ; i != dataValue.Len() ; i++{
		singleDataValue := dataValue.Index(i)
		singleDataFieldValue := singleDataValue.FieldByName(name)
		result.SetMapIndex(singleDataFieldValue,singleDataValue)
	}
	return result.Interface()
}

func ArrayColumnTable(column interface{},data interface{})([][]string){
	result := [][]string{}

	columnKeys,columnValues := ArrayKeyAndValue(column)
	columnKeysReal := columnKeys.([]string)
	columnValuesReal := columnValues.([]string)
	result = append(result,columnValuesReal)

	dataValue := reflect.ValueOf(data)
	for i := 0 ; i != dataValue.Len() ; i++{
		singleDataValue := dataValue.Index(i)
		singleDataStringValue := ArrayMapping(singleDataValue.Interface())
		singleDataStringValueData := reflect.ValueOf(singleDataStringValue)
		singleResult := []string{}
		for _,singleColumn := range columnKeysReal{
			singleResultString := ""
			singleValue := singleDataStringValueData.MapIndex(reflect.ValueOf(singleColumn))
			if singleValue.IsValid() == false{
				singleResultString = ""
			}else{
				singleResultString = fmt.Sprintf("%v",singleValue)
			}
			singleResult = append(
				singleResult,
				singleResultString,
			)
		}
		result = append(result,singleResult)
	}
	return result
}

func ArrayKeyAndValue(data interface{})(interface{},interface{}){
	//解析data
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Map{
		panic("need a map for arrayKeyAndValue")
	}
	dataKeyType := dataType.Key()
	dataValueType := dataType.Elem()

	//合并数据
	dataKeySlice := reflect.MakeSlice(reflect.SliceOf(dataKeyType),0,0)
	dataValueSlice := reflect.MakeSlice(reflect.SliceOf(dataValueType),0,0)
	dataValue := reflect.ValueOf(data)
	for _,singleKey := range dataValue.MapKeys(){
		dataKeySlice = reflect.Append(dataKeySlice,singleKey)
		dataValueSlice = reflect.Append(dataValueSlice,dataValue.MapIndex(singleKey))
	}
	return dataKeySlice.Interface(),dataValueSlice.Interface()
}

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

func ArrayMapping(data interface{})(interface{}){
	result,err := structToMap(data)
	if err != nil{
		panic(err)
	}
	return result
}

