package language

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

func nameMapper(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func combileMap(result map[string]interface{}, singleResultMap reflect.Value) {
	singleResultMapType := singleResultMap.Type()
	if singleResultMapType.Kind() != reflect.Map {
		return
	}
	singleResultMapKeys := singleResultMap.MapKeys()
	for _, singleKey := range singleResultMapKeys {
		result[fmt.Sprintf("%v", singleKey)] = singleResultMap.MapIndex(singleKey).Interface()
	}
}

type arrayMappingStructInfo struct {
	name      string
	omitempty bool
	num       int
	anonymous bool
	index     []int
}

type arrayMappingInfo struct {
	kind       int
	isTimeType bool
	field      []arrayMappingStructInfo
}

var arrayMappingInfoMap struct {
	mutex sync.RWMutex
	data  map[string]map[reflect.Type]arrayMappingInfo
}

var interfaceType reflect.Type

func getDataTagInfoInner(dataType reflect.Type, tag string) arrayMappingInfo {
	dataTypeKind := GetTypeKind(dataType)
	result := arrayMappingInfo{}
	result.kind = dataTypeKind
	if dataTypeKind == TypeKind.STRUCT {
		if dataType == reflect.TypeOf(time.Time{}) {
			//时间类型
			result.isTimeType = true
		} else {
			//结构体类型
			result.isTimeType = false
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
				single.index = singleDataType.Index
				single.anonymous = singleDataType.Anonymous
				result.field = append(result.field, single)
			}
		}
	}
	return result
}

func getDataTagInfo(target reflect.Type, tag string) arrayMappingInfo {
	arrayMappingInfoMap.mutex.RLock()
	var result arrayMappingInfo
	var ok bool
	resultArray, okArray := arrayMappingInfoMap.data[tag]
	if okArray {
		result, ok = resultArray[target]
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

func arrayToMapInner(dataValue reflect.Value, tag string) (reflect.Value, bool) {
	if dataValue.IsValid() == false {
		return dataValue, true
	} else {
		var result reflect.Value
		var isEmpty bool
		dataType := getDataTagInfo(dataValue.Type(), tag)
		if dataType.kind == TypeKind.STRUCT && dataType.isTimeType == true {
			timeValue := dataValue.Interface().(time.Time)
			result = reflect.ValueOf(timeValue.Format("2006-01-02 15:04:05"))
			isEmpty = IsEmptyValue(dataValue)
		} else if dataType.kind == TypeKind.STRUCT && dataType.isTimeType == false {
			resultMap := map[string]interface{}{}
			for _, singleType := range dataType.field {
				singleResultMap, isEmptyValue := arrayToMapInner(dataValue.Field(singleType.num), tag)
				if singleType.anonymous == false {
					if singleType.omitempty == true && isEmptyValue {
						continue
					}
					if singleResultMap.IsValid() == false {
						continue
					}
					resultMap[singleType.name] = singleResultMap.Interface()
				} else {
					combileMap(resultMap, singleResultMap)
				}
			}
			result = reflect.ValueOf(resultMap)
			isEmpty = (len(resultMap) == 0)
		} else if dataType.kind == TypeKind.ARRAY {
			resultSlice := []interface{}{}
			dataLen := dataValue.Len()
			for i := 0; i != dataLen; i++ {
				singleDataValue := dataValue.Index(i)
				singleDataResult, _ := arrayToMapInner(singleDataValue, tag)
				resultSlice = append(resultSlice, singleDataResult.Interface())
			}
			result = reflect.ValueOf(resultSlice)
			isEmpty = (len(resultSlice) == 0)
		} else if dataType.kind == TypeKind.MAP {
			dataKeyType := dataValue.Type().Key()
			resultMapType := reflect.MapOf(dataKeyType, interfaceType)
			resultMap := reflect.MakeMap(resultMapType)
			dataKeys := dataValue.MapKeys()
			for _, singleDataKey := range dataKeys {
				singleDataValue := dataValue.MapIndex(singleDataKey)
				singleDataResult, _ := arrayToMapInner(singleDataValue, tag)
				resultMap.SetMapIndex(singleDataKey, singleDataResult)
			}
			result = resultMap
			isEmpty = (len(dataKeys) == 0)
		} else if dataType.kind == TypeKind.INTERFACE ||
			dataType.kind == TypeKind.PTR {
			result, isEmpty = arrayToMapInner(dataValue.Elem(), tag)
		} else {
			result = dataValue
			isEmpty = IsEmptyValue(dataValue)
		}
		return result, isEmpty
	}
}

func ArrayToMap(data interface{}, tag string) interface{} {
	dataValue, _ := arrayToMapInner(reflect.ValueOf(data), tag)
	if dataValue.IsValid() == false {
		return nil
	} else {
		return dataValue.Interface()
	}
}

func mapToBool(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind == TypeKind.BOOL {
		target.SetBool(dataValue.Bool())
		return nil
	} else if dataKind == TypeKind.STRING {
		dataBool, err := strconv.ParseBool(dataValue.String())
		if err != nil {
			return errors.New(fmt.Sprintf("不是布尔值，其值为[%s]", dataValue.String()))
		}
		target.SetBool(dataBool)
		return nil
	} else {
		return errors.New(fmt.Sprintf("不是布尔值，其类型为[%s]", dataValue.Type().String()))
	}
}

func mapToUint(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind == TypeKind.UINT {
		target.SetUint(dataValue.Uint())
		return nil
	} else if dataKind == TypeKind.INT {
		target.SetUint(uint64(dataValue.Int()))
		return nil
	} else if dataKind == TypeKind.FLOAT {
		target.SetUint(uint64(math.Floor(dataValue.Float() + 0.5)))
		return nil
	} else if dataKind == TypeKind.STRING {
		dataUint, err := strconv.ParseUint(dataValue.String(), 10, 64)
		if err != nil {
			return errors.New(fmt.Sprintf("不是无符号整数，其值为[%s]", dataValue.String()))
		}
		target.SetUint(dataUint)
		return nil
	} else {
		return errors.New(fmt.Sprintf("不是无符号整数，其类型为[%s]", dataValue.Type().String()))
	}
}

func mapToInt(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind == TypeKind.INT {
		target.SetInt(dataValue.Int())
		return nil
	} else if dataKind == TypeKind.UINT {
		target.SetInt(int64(dataValue.Uint()))
		return nil
	} else if dataKind == TypeKind.FLOAT {
		target.SetInt(int64(math.Floor(dataValue.Float() + 0.5)))
		return nil
	} else if dataKind == TypeKind.STRING {
		dataInt, err := strconv.ParseInt(dataValue.String(), 10, 64)
		if err != nil {
			return errors.New(fmt.Sprintf("不是整数，其值为[%s]", dataValue.String()))
		}
		target.SetInt(dataInt)
		return nil
	} else {
		return errors.New(fmt.Sprintf("不是整数，其类型为[%s]", dataValue.Type().String()))
	}
}

func mapToFloat(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind == TypeKind.FLOAT {
		target.SetFloat(dataValue.Float())
		return nil
	} else if dataKind == TypeKind.INT {
		target.SetFloat(float64(dataValue.Int()))
		return nil
	} else if dataKind == TypeKind.UINT {
		target.SetFloat(float64(dataValue.Uint()))
		return nil
	} else if dataKind == TypeKind.STRING {
		dataFloat, err := strconv.ParseFloat(dataValue.String(), 64)
		if err != nil {
			return errors.New(fmt.Sprintf("不是浮点数，其值为[%s]", dataValue.String()))
		}
		target.SetFloat(dataFloat)
		return nil
	} else {
		return errors.New(fmt.Sprintf("不是浮点数，其类型为[%s]", dataValue.Type().String()))
	}
}

func mapToString(dataValue reflect.Value, target reflect.Value) error {
	stringValue := fmt.Sprintf("%v", dataValue.Interface())
	target.SetString(stringValue)
	return nil
}

func mapToArray(dataValue reflect.Value, target reflect.Value, tag string) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind != TypeKind.ARRAY {
		return errors.New(fmt.Sprintf("不是数组，其类型为[%s]", dataValue.Type().String()))
	}
	dataLen := dataValue.Len()
	targetType := target.Type()
	if targetType.Kind() == reflect.Slice {
		newTarget := reflect.MakeSlice(targetType, dataLen, dataLen)
		for i := 0; i != dataLen; i++ {
			singleData := dataValue.Index(i)
			singleDataTarget := newTarget.Index(i)
			err := mapToArrayInner(singleData, singleDataTarget, tag)
			if err != nil {
				return err
			}
		}
		target.Set(newTarget)
	} else {
		newTarget := reflect.New(targetType)
		for i := 0; i != newTarget.Len(); i++ {
			if i >= dataLen {
				continue
			}
			singleData := dataValue.Index(i)
			singleDataTarget := newTarget.Index(i)
			err := mapToArrayInner(singleData, singleDataTarget, tag)
			if err != nil {
				return err
			}
		}
		target.Set(newTarget)
	}
	return nil
}

func mapToMap(dataValue reflect.Value, target reflect.Value, tag string) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind != TypeKind.MAP {
		return errors.New(fmt.Sprintf("不是映射，其类型为[%s]", dataValue.Type().String()))
	}
	dataKeys := dataValue.MapKeys()
	targetType := target.Type()
	targetKeyType := targetType.Key()
	targetValueType := targetType.Elem()
	newTarget := reflect.MakeMap(targetType)
	for _, singleDataKey := range dataKeys {
		singleDataValue := dataValue.MapIndex(singleDataKey)

		singleDataTargetKey := reflect.New(targetKeyType)
		singleDataTargetValue := reflect.New(targetValueType)
		err := mapToArrayInner(singleDataKey, singleDataTargetKey, tag)
		if err != nil {
			return err
		}
		err = mapToArrayInner(singleDataValue, singleDataTargetValue, tag)
		if err != nil {
			return errors.New(fmt.Sprintf("参数%s%s", singleDataKey, err.Error()))
		}
		newTarget.SetMapIndex(singleDataTargetKey.Elem(), singleDataTargetValue.Elem())
	}
	target.Set(newTarget)
	return nil
}

func mapToTime(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	if dataType == reflect.TypeOf(time.Time{}) {
		target.Set(dataValue)
	} else if dataType.Kind() == reflect.String {
		timeValue, err := time.ParseInLocation("2006-01-02 15:04:05", dataValue.String(), time.Now().Local().Location())
		if err != nil {
			return errors.New(fmt.Sprintf("不是时间，其值为[%s]", dataValue.String()))
		}
		target.Set(reflect.ValueOf(timeValue))
		return nil
	}
	return errors.New(fmt.Sprintf("不是时间，其类型为[%s]", dataValue.Type().String()))
}

func mapToStruct(dataValue reflect.Value, target reflect.Value, targetType arrayMappingInfo, tag string) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind != TypeKind.MAP {
		return errors.New(fmt.Sprintf("不是映射，其类型为[%s]", dataValue.Type().String()))
	}
	dataTypeKey := dataType.Key()
	for _, singleStructInfo := range targetType.field {
		if singleStructInfo.anonymous == true {
			//FIXME 暂不考虑匿名结构体的覆盖问题
			singleDataValue := target.FieldByIndex(singleStructInfo.index)
			err := mapToArrayInner(dataValue, singleDataValue, tag)
			if err != nil {
				return errors.New(fmt.Sprintf("参数%s%s", singleStructInfo.name, err.Error()))
			}
		} else {
			singleMapKey := reflect.New(dataTypeKey)
			singleDataKey := reflect.ValueOf(singleStructInfo.name)
			err := mapToArrayInner(singleDataKey, singleMapKey, tag)
			if err != nil {
				return err
			}

			singleDataValue := target.FieldByIndex(singleStructInfo.index)
			singleMapResult := dataValue.MapIndex(singleMapKey.Elem())
			if singleMapResult.IsValid() == false {
				continue
			}
			err = mapToArrayInner(singleMapResult, singleDataValue, tag)
			if err != nil {
				return errors.New(fmt.Sprintf("参数%s%s", singleDataKey, err.Error()))
			}
		}
	}
	return nil
}

func mapToPtr(dataValue reflect.Value, target reflect.Value, tag string) error {
	targetElem := target.Elem()
	if targetElem.IsValid() == false {
		targetElem = reflect.New(target.Type().Elem())
		target.Set(targetElem)
	}
	return mapToArrayInner(dataValue, targetElem, tag)
}

func mapToInterface(dataValue reflect.Value, target reflect.Value, tag string) error {
	targetElem := target.Elem()
	if targetElem.IsValid() == false {
		target.Set(dataValue)
		return nil
	}
	return mapToArrayInner(dataValue, targetElem, tag)
}

func mapToArrayInner(data reflect.Value, target reflect.Value, tag string) error {
	//处理data是个nil的问题
	if data.IsValid() == false {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}
	//处理data是多层嵌套的问题
	dataKind := data.Type().Kind()
	if dataKind == reflect.Interface {
		return mapToArrayInner(data.Elem(), target, tag)
	} else if dataKind == reflect.Ptr {
		return mapToArrayInner(data.Elem(), target, tag)
	}
	//根据target是多层嵌套的问题
	targetType := getDataTagInfo(target.Type(), tag)
	if targetType.kind == TypeKind.PTR {
		return mapToPtr(data, target, tag)
	} else if targetType.kind == TypeKind.INTERFACE {
		return mapToInterface(data, target, tag)
	}
	//处理data是个空字符串
	if dataKind == reflect.String && data.String() == "" {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}
	if targetType.kind == TypeKind.BOOL {
		return mapToBool(data, target)
	} else if targetType.kind == TypeKind.INT {
		return mapToInt(data, target)
	} else if targetType.kind == TypeKind.UINT {
		return mapToUint(data, target)
	} else if targetType.kind == TypeKind.FLOAT {
		return mapToFloat(data, target)
	} else if targetType.kind == TypeKind.STRING {
		return mapToString(data, target)
	} else if targetType.kind == TypeKind.ARRAY {
		return mapToArray(data, target, tag)
	} else if targetType.kind == TypeKind.MAP {
		return mapToMap(data, target, tag)
	} else if targetType.kind == TypeKind.STRUCT {
		if targetType.isTimeType {
			return mapToTime(data, target)
		} else {
			return mapToStruct(data, target, targetType, tag)
		}
	} else {
		return errors.New("unkown target type " + target.Type().String())
	}
}

func MapToArray(data interface{}, target interface{}, tag string) error {
	dataValue := reflect.ValueOf(data)
	targetValue := reflect.ValueOf(target)
	if targetValue.IsValid() == false {
		return errors.New("target is nil")
	}
	if targetValue.Kind() != reflect.Ptr {
		return errors.New("invalid target is not ptr")
	}
	return mapToArrayInner(dataValue, targetValue, tag)
}

func init() {
	arrayMappingInfoMap.data = map[string]map[reflect.Type]arrayMappingInfo{}
	var mm struct {
		Test interface{}
	}
	interfaceType = reflect.TypeOf(mm).Field(0).Type
}
