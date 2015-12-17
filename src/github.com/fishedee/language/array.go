package language

import (
	"reflect"
)

func ArrayColumnKey(data interface{},name string)(interface{}){
	//提取信息
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Slice{
		panic("array column key should be a slice")
	}
	dataElemType := dataType.Elem()
	if dataElemType.Kind() != reflect.Struct{
		panic("array column key should be a struct")
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
		panic("array column key should be a slice")
	}
	dataElemType := dataType.Elem()
	if dataElemType.Kind() != reflect.Struct{
		panic("array column key should be a struct")
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
