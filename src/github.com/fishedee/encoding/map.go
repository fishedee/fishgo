package encoding

import (
	"fmt"
	"reflect"
)

func nameMapper(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func combileMap(result map[string]interface{}, singleResultMap interface{}) {
	singleResultMapMap, ok := singleResultMap.(map[string]interface{})
	if !ok {
		return
	}
	for key, value := range singleResultMapMap {
		result[key] = value
	}
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
	data  map[string]map[reflect.Type]arrayMappingInfo
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

func getDataTagInfoInner(dataType reflect.Type, tag string) arrayMappingInfo {
	result := arrayMappingInfo{}
	dataKind := dataType.Kind()
	if dataKind == reflect.Struct {
		if dataType == reflect.TypeOf(time.Time{}) {
			//时间类型
			result.kind = 1
		} else {
			//结构体类型
			result.kind = 2
			result.field = []arrayMappingStructInfo{}
			for i := 0; i != dataType.NumField(); i++ {
				singleDataType := dataType.Field(i)
				if singleDataType.PkgPath != "" && singleDataType.Anonymous == false {
					continue
				}
				var singleName string
				var omitempty bool
				if singleDataType.Tag.Get(tag) != "" {
					jsonTag := singleDataType.Tag.Get(tag)
					jsonTagList := strings.Split(jsonTag, ",")
					if jsonTagList[0] == "-" {
						continue
					}
					if len(jsonTagList) >= 2 && jsonTagList[1] == "omitempty" {
						omitempty = true
					}
					singleName = jsonTagList[0]
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
	} else if dataKind == reflect.Slice {
		//数组类型
		result.kind = 3
	} else if dataKind == reflect.Map{
		//映射类型
		result.kind = 4
	} else if dataKind == reflect.Int || dataKind == reflect.Int8 || dataKind == reflect.Int16 ||
		dataKind == reflect.Int32 || dataKind == reflect.Int64{
		//有符号整数类型
		result.kind = 5
	} else if dataKind == reflect.Uint || dataKind == reflect.Uint8 || dataKind == reflect.Uint16 ||
		dataKind == reflect.Uint32 || dataKind == reflect.Uint64{
		//无符号整数类型
		result.kind = 6
	}else if dataKind == reflect.Float32 || reflect.Float64{
		//浮点类型
		result.kind = 7
	}else if dataKind == reflect.Bool{
		//布尔类型
		result.kind = 8
	}else if dataKind == reflect.Uintptr{
		//指针类型
		result.kind = 9
	}else if dataKind == reflect.Interface{
		//interface类型
		result.kind = 10
	}
	return result
}

func getDataTagInfo(target reflect.Type, tag string) arrayMappingInfo {
	arrayMappingInfoMap.mutex.RLock()
	var result arrayMappingStructInfo
	var ok bool
	resultArray, okArray := arrayMappingInfoMap.data[tag]
	if okArray {
		result, ok := resultArray[target]
	}
	arrayMappingInfoMap.mutex.RUnlock()

	if ok {
		return result
	}
	result = getDataTagInfoInner(target, tag)

	arrayMappingInfoMap.mutex.Lock()
	if !okArray {
		resultArray = map[reflect.Type]arrayMappingInfo{}
		arrayMappingInfoMap.data[tag] = resultArray
	}
	resultArray[target] = result
	arrayMappingInfoMap.mutex.Unlock()

	return result
}

func encodeMapInner(data interface{}, tag string) interface{} {
	var result interface{}
	if data == nil {
		result = nil
	} else {
		dataType := getDataTagInfo(reflect.TypeOf(data), tag)
		dataValue := reflect.ValueOf(data)
		if dataType.kind == 1 {
			timeValue := data.(time.Time)
			result = timeValue.Format("2006-01-02 15:04:05")
		} else if dataType.kind == 2 {
			resultMap := map[string]interface{}{}
			for _, singleType := range dataType.field {
				singleResultMap := encodeMapInner(dataValue.Field(singleType.num).Interface(), tag)
				if singleType.anonymous == false {
					if singleType.omitempty == true && isEmptyValue(reflect.ValueOf(singleResultMap)) {
						continue
					}
					resultMap[singleType.name] = singleResultMap
				} else {
					combileMap(resultMap, singleResultMap)
				}
			}
			result = resultMap
		} else if dataType.kind == 3 {
			resultSlice := []interface{}{}
			dataLen := dataValue.Len()
			for i := 0; i != dataLen; i++ {
				singleDataValue := dataValue.Index(i)
				singleDataResult = encodeMapInner(singleDataValue.Interface(), tag)
				resultSlice = append(resultSlice, singleDataResult)
			}
			result = resultSlice
		} else {
			result = data
		}
	}
	return result, nil
}

func EncodeMap(data interface{}, tag string) (interface{}, error) {
	return encodeMapInner(data, tag), nil
}

func decodeTime(data interface{},target reflect.Value)error{
	timeValue,ok := data.(time.Time)
	if !ok{
		return errors.New(fmt.Sprintf("can't parse %s to time.time",reflect.TypeOf(data).String()))
	}
	target.Set(data)
}

func decodeInt(data interface{})
func decodeMapInner(data interface{}, target reflect.Value) error {
	if data == nil {
		return nil
	} else {
		targetType := getDataTagInfo(target.Type(), tag)
		if targetType.kind == 1 {
			
			timeValue.
		}
	}
}

func DecodeMap(data interface{}, target interface{}) error {
	return nil
}

func init() {
	arrayMappingInfoMap.data = map[string]map[reflect.Type]arrayMappingInfo{}
}
