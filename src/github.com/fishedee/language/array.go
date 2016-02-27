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
	omitempty bool
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
				var singleName string
				var omitempty bool
				if singleDataType.Tag.Get("json") != "" {
					jsonTag := singleDataType.Tag.Get("json")
					jsonTagList := strings.Split(jsonTag, ",")
					if jsonTagList[0] == "-" {
						continue
					}
					if len(jsonTagList) >= 2 && jsonTagList[1] == "omitempty" {
						omitempty = true
					}
					singleName = singleDataType.Tag.Get("json")
				} else {
					singleName = nameMapper(singleDataType.Name)
				}
				single := arrayMappingStructInfo{}
				single.name = singleName
				single.omitempty = omitempty
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

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
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
					if singleType.omitempty == true && isEmptyValue(reflect.ValueOf(singleResultMap)) {
						continue
					}
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
			dataLen := dataValue.Len()
			for i := 0; i != dataLen; i++ {
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

	for i := 0; i != dataLen; i++ {
		result.Index(dataLen - i - 1).Set(dataValue.Index(i))
	}
	return result.Interface()
}

func ArrayMappingByJsonOrFirstLower(data interface{}) interface{} {
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

func init() {
	arrayMappingInfoMap.data = map[reflect.Type]arrayMappingInfo{}
}
