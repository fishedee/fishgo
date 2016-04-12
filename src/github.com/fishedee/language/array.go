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

func ArrayUnique(arrayData interface{}) interface{} {
	arrayValue := reflect.ValueOf(arrayData)
	arrayType := arrayValue.Type()
	arrayLen := arrayValue.Len()

	result := reflect.MakeSlice(arrayType, 0, 0)
	resultTemp := map[interface{}]bool{}

	for i := 0; i != arrayLen; i++ {
		singleArrayDataValue := arrayValue.Index(i)
		singleArrayDataValueInterface := singleArrayDataValue.Interface()
		_, isExist := resultTemp[singleArrayDataValueInterface]
		if isExist == true {
			continue
		}
		resultTemp[singleArrayDataValueInterface] = true
		result = reflect.Append(result, singleArrayDataValue)
	}
	return result.Interface()
}

func sliceToMap(arrayData []interface{}) map[interface{}]bool {
	result := map[interface{}]bool{}

	for _, singleArray := range arrayData {
		arrayValue := reflect.ValueOf(singleArray)
		arrayLen := arrayValue.Len()

		for i := 0; i != arrayLen; i++ {
			result[arrayValue.Index(i).Interface()] = true
		}
	}

	return result
}

func ArrayDiff(arrayData interface{}, arrayData2 interface{}, arrayOther ...interface{}) interface{} {
	arrayOther = append([]interface{}{arrayData2}, arrayOther...)
	arrayOtherMap := sliceToMap(arrayOther)

	arrayValue := reflect.ValueOf(arrayData)
	arrayType := arrayValue.Type()
	arrayLen := arrayValue.Len()
	result := reflect.MakeSlice(arrayType, 0, 0)

	for i := 0; i != arrayLen; i++ {
		singleArrayDataValue := arrayValue.Index(i)
		singleArrayDataValueInterface := singleArrayDataValue.Interface()

		_, isExist := arrayOtherMap[singleArrayDataValueInterface]
		if isExist == true {
			continue
		}
		result = reflect.Append(result, singleArrayDataValue)
		arrayOtherMap[singleArrayDataValueInterface] = true
	}

	return result.Interface()
}

func ArrayIntersect(arrayData interface{}, arrayData2 interface{}, arrayOther ...interface{}) interface{} {
	arrayOther = append([]interface{}{arrayData2}, arrayOther...)
	arrayOtherMap := sliceToMap(arrayOther)

	arrayValue := reflect.ValueOf(arrayData)
	arrayType := arrayValue.Type()
	arrayLen := arrayValue.Len()
	result := reflect.MakeSlice(arrayType, 0, 0)

	for i := 0; i != arrayLen; i++ {
		singleArrayDataValue := arrayValue.Index(i)
		singleArrayDataValueInterface := singleArrayDataValue.Interface()

		isFirst, isExist := arrayOtherMap[singleArrayDataValueInterface]
		if isExist == false || isFirst == false {
			continue
		}
		result = reflect.Append(result, singleArrayDataValue)
		arrayOtherMap[singleArrayDataValueInterface] = false
	}

	return result.Interface()
}

func ArrayMerge(arrayData interface{}, arrayData2 interface{}, arrayOther ...interface{}) interface{} {
	arrayOther = append([]interface{}{arrayData2}, arrayOther...)
	arrayOtherMap := sliceToMap(arrayOther)

	arrayValue := reflect.ValueOf(arrayData)
	arrayType := arrayValue.Type()
	arrayLen := arrayValue.Len()
	result := reflect.MakeSlice(arrayType, 0, 0)

	for i := 0; i != arrayLen; i++ {
		singleArrayDataValue := arrayValue.Index(i)
		singleArrayDataValueInterface := singleArrayDataValue.Interface()

		isFirst, isExist := arrayOtherMap[singleArrayDataValueInterface]
		if isExist == true && isFirst == false {
			continue
		}
		result = reflect.Append(result, singleArrayDataValue)
		arrayOtherMap[singleArrayDataValueInterface] = false
	}

	for single, isFirst := range arrayOtherMap {
		if isFirst == false {
			continue
		}
		result = reflect.Append(result, reflect.ValueOf(single))
	}

	return result.Interface()
}
