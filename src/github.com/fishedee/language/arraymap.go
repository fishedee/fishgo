package language

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
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
	index     []int
	t         reflect.Type
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

func arrayToMapInner(data interface{}, tag string) interface{} {
	var result interface{}
	if data == nil {
		result = nil
	} else {
		dataType := getDataTagInfo(reflect.TypeOf(data), tag)
		dataValue := reflect.ValueOf(data)
		if dataType.kind == TypeKind.STRUCT && dataType.isTimeType == true {
			timeValue := data.(time.Time)
			result = timeValue.Format("2006-01-02 15:04:05")
		} else if dataType.kind == TypeKind.STRUCT && dataType.isTimeType == false {
			resultMap := map[string]interface{}{}
			for _, singleType := range dataType.field {
				singleResultMap := arrayToMapInner(dataValue.Field(singleType.num).Interface(), tag)
				if singleType.anonymous == false {
					if singleType.omitempty == true && IsEmptyValue(reflect.ValueOf(singleResultMap)) {
						continue
					}
					resultMap[singleType.name] = singleResultMap
				} else {
					combileMap(resultMap, singleResultMap)
				}
			}
			result = resultMap
		} else if dataType.kind == TypeKind.ARRAY {
			resultSlice := []interface{}{}
			dataLen := dataValue.Len()
			for i := 0; i != dataLen; i++ {
				singleDataValue := dataValue.Index(i)
				singleDataResult := arrayToMapInner(singleDataValue.Interface(), tag)
				resultSlice = append(resultSlice, singleDataResult)
			}
			result = resultSlice
		} else {
			result = data
		}
	}
	return result
}

func ArrayToMap(data interface{}, tag string) interface{} {
	return arrayToMapInner(data, tag)
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
			return err
		}
		target.SetBool(dataBool)
		return nil
	} else {
		return errors.New(fmt.Sprintf("can't parse %s to bool", dataValue.Type().String()))
	}
}

func mapToUint(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind == TypeKind.UINT {
		target.SetUint(dataValue.Uint())
		return nil
	} else if dataKind == TypeKind.STRING {
		dataUint, err := strconv.ParseUint(dataValue.String(), 10, 64)
		if err != nil {
			return err
		}
		target.SetUint(dataUint)
		return nil
	} else {
		return errors.New(fmt.Sprintf("can't parse %s to uint", dataValue.Type().String()))
	}
}

func mapToInt(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind == TypeKind.INT {
		target.SetInt(dataValue.Int())
		return nil
	} else if dataKind == TypeKind.STRING {
		dataInt, err := strconv.ParseInt(dataValue.String(), 10, 64)
		if err != nil {
			return err
		}
		target.SetInt(dataInt)
		return nil
	} else {
		return errors.New(fmt.Sprintf("can't parse %s to int", dataValue.Type().String()))
	}
}

func mapToFloat(dataValue reflect.Value, target reflect.Value) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind == TypeKind.FLOAT {
		target.SetFloat(dataValue.Float())
		return nil
	} else if dataKind == TypeKind.STRING {
		dataFloat, err := strconv.ParseFloat(dataValue.String(), 64)
		if err != nil {
			return err
		}
		target.SetFloat(dataFloat)
		return nil
	} else {
		return errors.New(fmt.Sprintf("can't parse %s to float", dataValue.Type().String()))
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
		return errors.New(fmt.Sprintf("can't parse %s to array", dataValue.Type().String()))
	}
	dataLen := dataValue.Len()
	targetType := target.Type()
	targetElemType := targetType.Elem()
	if targetType.Kind() == reflect.Slice {
		newTarget := reflect.MakeSlice(targetElemType, dataLen, dataLen)
		for i := 0; i != dataLen; i++ {
			singleData := dataValue.Index(i)
			singleDataTarget := reflect.New(targetElemType)
			err := mapToArrayInner(singleData, singleDataTarget, tag)
			if err != nil {
				return err
			}
			newTarget.Index(i).Set(singleDataTarget)
		}
		target.Set(newTarget)
	} else {
		newTarget := reflect.New(targetType)
		for i := 0; i != newTarget.Len(); i++ {
			if i >= dataLen {
				continue
			}
			singleData := dataValue.Index(i)
			singleDataTarget := reflect.New(targetElemType)
			err := mapToArrayInner(singleData, singleDataTarget, tag)
			if err != nil {
				return err
			}
			newTarget.Index(i).Set(singleDataTarget)
		}
		target.Set(newTarget)
	}
	return nil
}

func mapToMap(data reflect.Value, target reflect.Value, tag string) error {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind != TypeKind.MAP {
		return errors.New(fmt.Sprintf("can't parse %s to map", reflect.TypeOf(data).String()))
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
			return err
		}
		newTarget.SetMapIndex(singleDataTargetKey, singleDataTargetValue)
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
			return err
		}
		target.Set(reflect.ValueOf(timeValue))
		return nil
	}
	return errors.New(fmt.Sprintf("can't parse %s to time.time", dataValue.Type().String()))
}

func mapToStruct(dataValue reflect.Value, target reflect.Value, targetType arrayMappingInfo, tag string) error {
	dataType := dataValue.Type()
	dataKind := GetTypeKind(dataType)
	if dataKind != TypeKind.MAP {
		return errors.New(fmt.Sprintf("can't parse %s to struct", dataType.String()))
	}
	if dataType.Key().Kind() != reflect.String {
		return errors.New(fmt.Sprintf("can't parse %s to string", dataType.Key().String()))
	}
	for _, singleStructInfo := range targetType.field {
		singleDataKey := reflect.ValueOf(singleStructInfo.name)
		singleDataValue := target.FieldByIndex(singleStructInfo.index)
		dataValue := dataValue.MapIndex(singleDataKey)
		if dataValue.IsValid() == false {
			continue
		}
		err := mapToArrayInner(singleDataValue, dataValue, tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func mapToPtr(dataValue reflect.Value, target reflect.Value, tag string) error {
	targetElem := target.Elem()
	if targetElem.IsValid() == false {
		targetElem = reflect.New(targetElem.Type())
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
	//根据target类型的不同来设置
	targetType := getDataTagInfo(target.Type(), tag)
	if targetType.kind == TypeKind.PTR {
		return mapToPtr(data, target, tag)
	} else if targetType.kind == TypeKind.INTERFACE {
		return mapToInterface(data, target, tag)
	} else if targetType.kind == TypeKind.BOOL {
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
	return mapToArrayInner(dataValue, targetValue, tag)
}

func init() {
	arrayMappingInfoMap.data = map[string]map[reflect.Type]arrayMappingInfo{}
}