package language

import (
	"fmt"
	"reflect"
	"sort"
	"time"
)

type arrayColumnCompare func(reflect.Value, reflect.Value) int
type arrayColumnSlice struct {
	target         reflect.Value
	targetElemType reflect.Type
	targetCompare  []arrayColumnCompare
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
	temp := reflect.New(this.targetElemType).Elem()
	left := this.target.Index(i)
	right := this.target.Index(j)
	temp.Set(left)
	left.Set(right)
	right.Set(temp)
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
	} else if fieldType.Kind() == reflect.Float32 ||
		fieldType.Kind() == reflect.Float64 {
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
	} else if fieldType.Kind() == reflect.String {
		return func(left reflect.Value, right reflect.Value) int {
			leftString := left.FieldByIndex(fieldIndex).String()
			rightString := right.FieldByIndex(fieldIndex).String()
			if leftString < rightString {
				return -1
			} else if leftString > rightString {
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
	targetElemType := this.target.Type().Elem()
	for _, singleName := range name {
		result = append(result, this.initSingleCompare(targetElemType, singleName))
	}
	this.targetCompare = result
	this.targetElemType = targetElemType
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
	dataLen := dataValue.Len()

	//去重
	var lastValue reflect.Value
	arraySlice := arrayColumnSlice{
		target: dataValue,
	}
	arraySlice.initCompare(name)
	for i := 0; i != dataLen; i++ {
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
	dataLen := dataValue.Len()
	for i := 0; i != dataLen; i++ {
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
	dataLen := dataValue.Len()
	for i := 0; i != dataLen; i++ {
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

	columnMap := column.(map[string]string)
	columnKeys, _ := ArrayKeyAndValue(column)
	columnKeysReal := columnKeys.([]string)
	columnValuesReal := []string{}
	sort.Sort(sort.StringSlice(columnKeysReal))
	for _, singleKey := range columnKeysReal {
		columnValuesReal = append(columnValuesReal, columnMap[singleKey])
	}
	result = append(result, columnValuesReal)

	dataValue := reflect.ValueOf(data)
	dataLen := dataValue.Len()
	for i := 0; i != dataLen; i++ {
		singleDataValue := dataValue.Index(i)
		singleDataStringValue := ArrayToMap(singleDataValue.Interface(), "json")
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
