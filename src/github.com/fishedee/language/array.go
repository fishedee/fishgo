package language

import (
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"
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

func nameMapper(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func combileMap(result map[string]interface{}, singleResultMap interface{}) error {
	singleResultMapMap, ok := singleResultMap.(map[string]interface{})
	if ok == false {
		return errors.New("Anonymous field is not a struct")
	}
	for key, value := range singleResultMapMap {
		result[key] = value
	}
	return nil
}

type arrayMappingStructInfo struct {
	name      string
	num       int
	anonymous bool
}

type arrayMappingInfo struct {
	kind  int
	field []arrayMappingStructInfo
}

var arrayMappingInfoMap struct {
	mutex sync.RWMutex
	data  map[reflect.Type]arrayMappingInfo
}

func getArrayMappingInfoInner(dataType reflect.Type) arrayMappingInfo {
	result := arrayMappingInfo{}
	if dataType.Kind() == reflect.Struct {
		if dataType == reflect.TypeOf(time.Time{}) {
			result.kind = 1
		} else {
			result.kind = 2
			result.field = []arrayMappingStructInfo{}
			for i := 0; i != dataType.NumField(); i++ {
				singleDataType := dataType.Field(i)
				if singleDataType.PkgPath != "" && singleDataType.Anonymous == false {
					continue
				}
				single := arrayMappingStructInfo{}
				single.name = nameMapper(singleDataType.Name)
				single.num = i
				single.anonymous = singleDataType.Anonymous
				result.field = append(result.field, single)
			}
		}
	} else if dataType.Kind() == reflect.Slice {
		result.kind = 3
	} else {
		result.kind = 4
	}
	return result
}

func getArrayMappingInfo(target reflect.Type) arrayMappingInfo {
	arrayMappingInfoMap.mutex.RLock()
	result, ok := arrayMappingInfoMap.data[target]
	arrayMappingInfoMap.mutex.RUnlock()

	if ok {
		return result
	}
	result = getArrayMappingInfoInner(target)

	arrayMappingInfoMap.mutex.Lock()
	arrayMappingInfoMap.data[target] = result
	arrayMappingInfoMap.mutex.Unlock()

	return result
}

func arrayMappingInner(data interface{}) (interface{}, error) {
	var result interface{}
	if data == nil {
		result = nil
	} else {
		dataType := getArrayMappingInfo(reflect.TypeOf(data))
		dataValue := reflect.ValueOf(data)
		if dataType.kind == 1 {
			timeValue := data.(time.Time)
			result = timeValue.Format("2006-01-02 15:04:05")
		} else if dataType.kind == 2 {
			resultMap := map[string]interface{}{}
			for _, singleType := range dataType.field {
				singleResultMap, err := arrayMappingInner(dataValue.Field(singleType.num).Interface())
				if err != nil {
					return result, err
				}
				if singleType.anonymous == false {
					resultMap[singleType.name] = singleResultMap
				} else {
					err := combileMap(resultMap, singleResultMap)
					if err != nil {
						return result, err
					}
				}
			}
			result = resultMap
		} else if dataType.kind == 3 {
			resultSlice := []interface{}{}
			for i := 0; i != dataValue.Len(); i++ {
				singleDataValue := dataValue.Index(i)
				singleDataResult, err := arrayMappingInner(singleDataValue.Interface())
				if err != nil {
					return result, err
				}
				resultSlice = append(resultSlice, singleDataResult)
			}
			result = resultSlice
		} else {
			result = data
		}
	}
	return result, nil
}

func ArrayReverse(data interface{}) interface{} {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataLen := dataValue.Len()
	result := reflect.MakeSlice(dataType, dataLen, dataLen)

	for i := 0; i != dataValue.Len(); i++ {
		result.Index(dataLen - i - 1).Set(dataValue.Index(i))
	}
	return result.Interface()
}

func ArrayMapping(data interface{}) interface{} {
	result, err := arrayMappingInner(data)
	if err != nil {
		panic(err)
	}
	return result
}

func ArrayIn(arrayData interface{}, findData interface{}) int {
	var findIndex int
	findIndex = -1
	arrayDataValue := reflect.ValueOf(arrayData)
	for i := 0; i != arrayDataValue.Len(); i++ {
		singleArrayDataValue := arrayDataValue.Index(i).Interface()
		if singleArrayDataValue == findData {
			findIndex = i
			break
		}
	}
	return findIndex
}

func init() {
	arrayMappingInfoMap.data = map[reflect.Type]arrayMappingInfo{}
}
