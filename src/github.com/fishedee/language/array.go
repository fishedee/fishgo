package language

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

type arrayColumnCompare func(reflect.Value, reflect.Value) int
type arrayColumnSlice struct {
	target        reflect.Value
	targetCompare []arrayColumnCompare
}

func (this *arrayColumnSlice) Len() int {
	return this.target.Len()
}

func (this *arrayColumnSlice) Less(i, j int) bool {
	left := this.target.Index(i)
	right := this.target.Index(j)
	return this.LessValue(left, right)
}
func (this *arrayColumnSlice) LessValue(left, right reflect.Value) bool {
	for _, singleCompare := range this.targetCompare {
		if singleCompare(left, right) < 0 {
			return true
		}
	}
	return false
}

func (this *arrayColumnSlice) Equal(i, j int) bool {
	left := this.target.Index(i)
	right := this.target.Index(j)
	return this.EqualValue(left, right)
}

func (this *arrayColumnSlice) EqualValue(left, right reflect.Value) bool {
	for _, singleCompare := range this.targetCompare {
		if singleCompare(left, right) != 0 {
			return false
		}
	}
	return true
}

func (this *arrayColumnSlice) Swap(i, j int) {
	left := this.target.Index(i)
	right := this.target.Index(j)
	right.Set(left)
	left.Set(right)
}

func (this *arrayColumnSlice) initSingleCompare(targetType reflect.Type, name string) arrayColumnCompare {
	field, ok := targetType.FieldByName(name)
	if !ok {
		panic(targetType.Name() + " has not name " + name)
	}
	fieldIndex := field.Index
	fieldType := field.Type
	if fieldType.Kind() == reflect.Bool {
		return func(left reflect.Value, right reflect.Value) int {
			leftBool := left.FieldByIndex(fieldIndex).Bool()
			rightBool := right.FieldByIndex(fieldIndex).Bool()
			if leftBool == rightBool {
				return 0
			} else if leftBool == false {
				return -1
			} else {
				return 1
			}
		}
	} else if fieldType.Kind() == reflect.Int {
		return func(left reflect.Value, right reflect.Value) int {
			leftInt := left.FieldByIndex(fieldIndex).Int()
			rightInt := right.FieldByIndex(fieldIndex).Int()
			if leftInt < rightInt {
				return -1
			} else if leftInt > rightInt {
				return 1
			} else {
				return 0
			}
		}
	} else if fieldType.Kind() == reflect.Float32 {
		return func(left reflect.Value, right reflect.Value) int {
			leftFloat := left.FieldByIndex(fieldIndex).Float()
			rightFloat := right.FieldByIndex(fieldIndex).Float()
			if leftFloat < rightFloat {
				return -1
			} else if leftFloat > rightFloat {
				return 1
			} else {
				return 0
			}
		}
	} else if fieldType.Kind() == reflect.Struct && fieldType == reflect.TypeOf(time.Time{}) {
		return func(left reflect.Value, right reflect.Value) int {
			leftTime := left.FieldByIndex(fieldIndex).Interface().(time.Time)
			rightTime := right.FieldByIndex(fieldIndex).Interface().(time.Time)
			if leftTime.Before(rightTime) {
				return -1
			} else if leftTime.After(rightTime) {
				return 1
			} else {
				return 0
			}
		}
	} else {
		panic(fieldType.Name() + " can not compare")
	}
}

func (this *arrayColumnSlice) initCompare(name []string) {
	result := []arrayColumnCompare{}
	targetType := this.target.Type().Elem()
	for _, singleName := range name {
		result = append(result, this.initSingleCompare(targetType, singleName))
	}
	this.targetCompare = result
}

func combineName(firstName string, names []string) []string {
	result := []string{}
	result = append(result, firstName)
	for _, singleName := range names {
		result = append(result, singleName)
	}
	return result
}

func ArrayColumnSort(data interface{}, firstName string, otherName ...string) interface{} {
	//建立一份拷贝数据
	name := combineName(firstName, otherName)
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataValueLen := dataValue.Len()

	dataResult := reflect.MakeSlice(dataType, dataValueLen, dataValueLen)
	reflect.Copy(dataResult, dataValue)

	//排序
	arraySlice := arrayColumnSlice{
		target: dataResult,
	}
	arraySlice.initCompare(name)
	sort.Sort(&arraySlice)
	return dataResult.Interface()
}

func ArrayColumnUnique(data interface{}, firstName string, otherName ...string) interface{} {
	//建立一份数据
	name := combineName(firstName, otherName)
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataResult := reflect.MakeSlice(dataType, 0, 0)

	//去重
	var lastValue reflect.Value
	arraySlice := arrayColumnSlice{
		target: dataValue,
	}
	arraySlice.initCompare(name)
	for i := 0; i != dataValue.Len(); i++ {
		singleValue := dataValue.Index(i)
		if i != 0 && arraySlice.EqualValue(lastValue, singleValue) {
			continue
		}
		lastValue = singleValue
		dataResult = reflect.Append(dataResult, singleValue)
	}
	return dataResult.Interface()
}

func ArrayColumnKey(data interface{}, name string) interface{} {
	//提取信息
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Slice {
		panic("array column should be a slice")
	}
	dataElemType := dataType.Elem()
	if dataElemType.Kind() != reflect.Struct {
		panic("array column element should be a struct")
	}
	dataElemFieldType, ok := dataElemType.FieldByName(name)
	if !ok {
		panic("dataElemFieldType has not filed " + name)
	}

	//整合slice
	resultType := reflect.SliceOf(dataElemFieldType.Type)
	result := reflect.MakeSlice(resultType, 0, 0)
	dataValue := reflect.ValueOf(data)
	for i := 0; i != dataValue.Len(); i++ {
		singleDataValue := dataValue.Index(i)
		singleDataFieldValue := singleDataValue.FieldByName(name)
		result = reflect.Append(result, singleDataFieldValue)
	}
	return result.Interface()
}

type arrayColumnMapInfo struct {
	Index   []int
	Type    reflect.Type
	MapType reflect.Type
}

func ArrayColumnMap(data interface{}, firstName string, otherName ...string) interface{} {
	//提取信息
	name := combineName(firstName, otherName)
	nameInfo := []arrayColumnMapInfo{}
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type().Elem()
	for _, singleName := range name {
		singleField, ok := dataType.FieldByName(singleName)
		if !ok {
			panic(dataType.Name() + " struct has not field " + singleName)
		}
		nameInfo = append(nameInfo, arrayColumnMapInfo{
			Index: singleField.Index,
			Type:  singleField.Type,
		})
	}
	prevType := dataType
	for i := len(nameInfo) - 1; i >= 0; i-- {
		nameInfo[i].MapType = reflect.MapOf(
			nameInfo[i].Type,
			prevType,
		)
		prevType = nameInfo[i].MapType
	}

	//整合map
	result := reflect.MakeMap(nameInfo[0].MapType)
	for i := 0; i != dataValue.Len(); i++ {
		singleValue := dataValue.Index(i)
		prevValue := result
		for singleNameIndex, singleNameInfo := range nameInfo {
			var nextValue reflect.Value
			singleField := singleValue.FieldByIndex(singleNameInfo.Index)
			nextValue = prevValue.MapIndex(singleField)
			if !nextValue.IsValid() {
				if singleNameIndex+1 < len(nameInfo) {
					nextValue = reflect.MakeMap(nameInfo[singleNameIndex+1].MapType)
				} else {
					nextValue = singleValue
				}
				prevValue.SetMapIndex(singleField, nextValue)
			}
			prevValue = nextValue
		}
	}
	return result.Interface()
}

func ArrayColumnTable(column interface{}, data interface{}) [][]string {
	result := [][]string{}

	columnKeys, columnValues := ArrayKeyAndValue(column)
	columnKeysReal := columnKeys.([]string)
	columnValuesReal := columnValues.([]string)
	result = append(result, columnValuesReal)

	dataValue := reflect.ValueOf(data)
	for i := 0; i != dataValue.Len(); i++ {
		singleDataValue := dataValue.Index(i)
		singleDataStringValue := ArrayMapping(singleDataValue.Interface())
		singleDataStringValueData := reflect.ValueOf(singleDataStringValue)
		singleResult := []string{}
		for _, singleColumn := range columnKeysReal {
			singleResultString := ""
			singleValue := singleDataStringValueData.MapIndex(reflect.ValueOf(singleColumn))
			if singleValue.IsValid() == false {
				singleResultString = ""
			} else {
				singleResultString = fmt.Sprintf("%v", singleValue)
			}
			singleResult = append(
				singleResult,
				singleResultString,
			)
		}
		result = append(result, singleResult)
	}
	return result
}

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
	dataElemType := dataValue.Elem().Type()
	dataLen := dataValue.Len()
	result := reflect.MakeSlice(dataElemType, dataLen, dataLen)

	for i := dataLen - 1; i >= 0; i-- {
		result = reflect.Append(result, dataValue.Index(i))
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
