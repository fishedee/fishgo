package language

import (
	"reflect"
)

func MapToSlice(data interface{})(interface{},interface{}){
	//解析data
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Map{
		panic("need a map for maptoslice")
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
