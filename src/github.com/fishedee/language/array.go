package language

import (
	"reflect"
)

func ArrayKeyAndValue(data interface{}) (interface{}, interface{}) {
	//解析data
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Map {
		panic("need a map for arrayKeyAndValue")
	}
	dataKeyType := dataType.Key()
	dataValueType := dataType.Elem()

	//合并数据
	dataKeySlice := reflect.MakeSlice(reflect.SliceOf(dataKeyType), 0, 0)
	dataValueSlice := reflect.MakeSlice(reflect.SliceOf(dataValueType), 0, 0)
	dataValue := reflect.ValueOf(data)
	for _, singleKey := range dataValue.MapKeys() {
		dataKeySlice = reflect.Append(dataKeySlice, singleKey)
		dataValueSlice = reflect.Append(dataValueSlice, dataValue.MapIndex(singleKey))
	}
	return dataKeySlice.Interface(), dataValueSlice.Interface()
}

func ArrayReverse(data interface{}) interface{} {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataLen := dataValue.Len()
	result := reflect.MakeSlice(dataType, dataLen, dataLen)

	for i := 0; i != dataLen; i++ {
		result.Index(dataLen - i - 1).Set(dataValue.Index(i))
	}
	return result.Interface()
}

func ArrayIn(arrayData interface{}, findData interface{}) int {
	var findIndex int
	findIndex = -1
	arrayDataValue := reflect.ValueOf(arrayData)
	arrayDataValueLen := arrayDataValue.Len()
	for i := 0; i != arrayDataValueLen; i++ {
		singleArrayDataValue := arrayDataValue.Index(i).Interface()
		if singleArrayDataValue == findData {
			findIndex = i
			break
		}
	}
	return findIndex
}
