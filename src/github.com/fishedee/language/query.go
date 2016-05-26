package language

import (
	"math"
	"reflect"
	"sort"
	"strings"
	"time"
)

//基础类函数
func QuerySelect(data interface{}, selectFuctor interface{}) interface{} {
	dataValue := reflect.ValueOf(data)
	dataLen := dataValue.Len()

	selectFuctorValue := reflect.ValueOf(selectFuctor)
	selectFuctorType := selectFuctorValue.Type()
	selectFuctorOuterType := selectFuctorType.Out(0)
	resultType := reflect.SliceOf(selectFuctorOuterType)
	resultValue := reflect.MakeSlice(resultType, dataLen, dataLen)

	for i := 0; i != dataLen; i++ {
		singleDataValue := dataValue.Index(i)
		singleResultValue := selectFuctorValue.Call([]reflect.Value{singleDataValue})[0]
		resultValue.Index(i).Set(singleResultValue)
	}
	return resultValue.Interface()
}

func QueryWhere(data interface{}, whereFuctor interface{}) interface{} {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataLen := dataValue.Len()

	whereFuctorValue := reflect.ValueOf(whereFuctor)
	resultType := reflect.SliceOf(dataType.Elem())
	resultValue := reflect.MakeSlice(resultType, 0, 0)

	for i := 0; i != dataLen; i++ {
		singleDataValue := dataValue.Index(i)
		singleResultValue := whereFuctorValue.Call([]reflect.Value{singleDataValue})[0]
		if singleResultValue.Bool() {
			resultValue = reflect.Append(resultValue, singleDataValue)
		}
	}
	return resultValue.Interface()
}

func QueryReduce(data interface{}, reduceFuctor interface{}, resultReduce interface{}) interface{} {
	dataValue := reflect.ValueOf(data)
	dataLen := dataValue.Len()

	reduceFuctorValue := reflect.ValueOf(reduceFuctor)
	resultReduceType := reduceFuctorValue.Type().In(0)
	resultReduceValue := reflect.New(resultReduceType)
	err := MapToArray(resultReduce, resultReduceValue.Interface(), "json")
	if err != nil {
		panic(err)
	}
	resultReduceValue = resultReduceValue.Elem()

	for i := 0; i != dataLen; i++ {
		singleDataValue := dataValue.Index(i)
		resultReduceValue = reduceFuctorValue.Call([]reflect.Value{resultReduceValue, singleDataValue})[0]
	}
	return resultReduceValue.Interface()
}

type queryCompare func(reflect.Value, reflect.Value) int
type querySortSlice struct {
	target         reflect.Value
	targetElemType reflect.Type
	targetCompare  []queryCompare
}

func (this *querySortSlice) Len() int {
	return this.target.Len()
}

func (this *querySortSlice) Less(i, j int) bool {
	left := this.target.Index(i)
	right := this.target.Index(j)
	for _, singleCompare := range this.targetCompare {
		compareResult := singleCompare(left, right)
		if compareResult < 0 {
			return true
		} else if compareResult > 0 {
			return false
		}
	}
	return false
}

func (this *querySortSlice) Swap(i, j int) {
	temp := reflect.New(this.targetElemType).Elem()
	left := this.target.Index(i)
	right := this.target.Index(j)
	temp.Set(left)
	left.Set(right)
	right.Set(temp)
}

func QuerySort(data interface{}, sortType string) interface{} {
	//拷贝一份
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataElemType := dataType.Elem()
	dataValueLen := dataValue.Len()

	dataResult := reflect.MakeSlice(dataType, dataValueLen, dataValueLen)
	reflect.Copy(dataResult, dataValue)

	//排序
	targetCompare := getQueryCompares(dataElemType, sortType)
	arraySlice := querySortSlice{
		target:         dataResult,
		targetElemType: dataElemType,
		targetCompare:  targetCompare,
	}
	sort.Sort(&arraySlice)

	return dataResult.Interface()
}

func QueryJoin(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFuctor interface{}) interface{} {
	//解析配置
	leftJoinType, rightJoinType := analyseJoin(joinType)

	leftDataValue := reflect.ValueOf(leftData)
	leftDataType := leftDataValue.Type()
	leftDataElemType := leftDataType.Elem()
	leftDataValueLen := leftDataValue.Len()
	leftDataJoinStruct, ok := getFieldByName(leftDataElemType, leftJoinType)
	if !ok {
		panic(leftDataElemType.Name() + " has no field " + leftJoinType)
	}
	leftDataJoin := leftDataJoinStruct.Index

	rightData = QuerySort(rightData, rightJoinType+" asc")
	rightDataValue := reflect.ValueOf(rightData)
	rightDataType := rightDataValue.Type()
	rightDataElemType := rightDataType.Elem()
	rightDataValueLen := rightDataValue.Len()
	rightDataJoinStruct, ok := getFieldByName(rightDataElemType, rightJoinType)
	if !ok {
		panic(rightDataElemType.Name() + " has no field " + rightJoinType)
	}
	rightDataJoin := rightDataJoinStruct.Index

	joinFuctorValue := reflect.ValueOf(joinFuctor)
	joinFuctorType := joinFuctorValue.Type()
	joinCompare := getSingleQueryCompare(leftDataJoinStruct.Type)
	resultValue := reflect.MakeSlice(reflect.SliceOf(joinFuctorType.Out(0)), 0, 0)

	rightHaveJoin := make([]bool, rightDataValueLen, rightDataValueLen)
	joinPlace = strings.Trim(strings.ToLower(joinPlace), " ")
	if ArrayIn([]string{"left", "right", "inner", "outer"}, joinPlace) == -1 {
		panic("invalid joinPlace [" + joinPlace + "] ")
	}

	//开始join
	for i := 0; i != leftDataValueLen; i++ {
		//二分查找右边对应的键
		singleLeftData := leftDataValue.Index(i)
		singleLeftDataJoin := singleLeftData.FieldByIndex(leftDataJoin)
		j := sort.Search(rightDataValueLen, func(j int) bool {
			return joinCompare(rightDataValue.Index(j).FieldByIndex(rightDataJoin), singleLeftDataJoin) >= 0
		})
		//合并双边满足条件
		haveFound := false
		for ; j < rightDataValueLen; j++ {
			singleRightData := rightDataValue.Index(j)
			singleRightDataJoin := singleRightData.FieldByIndex(rightDataJoin)
			if joinCompare(singleLeftDataJoin, singleRightDataJoin) != 0 {
				break
			}
			singleResult := joinFuctorValue.Call([]reflect.Value{singleLeftData, singleRightData})[0]
			resultValue = reflect.Append(resultValue, singleResult)
			haveFound = true
			rightHaveJoin[j] = true
		}
		//合并不满足的条件
		if !haveFound && (joinPlace == "left" || joinPlace == "outer") {
			singleRightData := reflect.New(rightDataElemType).Elem()
			singleResult := joinFuctorValue.Call([]reflect.Value{singleLeftData, singleRightData})[0]
			resultValue = reflect.Append(resultValue, singleResult)
		}
	}
	//处理剩余的右侧元素
	if joinPlace == "right" || joinPlace == "outer" {
		singleLeftData := reflect.New(leftDataElemType).Elem()
		rightHaveJoinLen := len(rightHaveJoin)
		for j := 0; j != rightHaveJoinLen; j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleRightData := rightDataValue.Index(j)
			singleResult := joinFuctorValue.Call([]reflect.Value{singleLeftData, singleRightData})[0]
			resultValue = reflect.Append(resultValue, singleResult)
		}
	}
	return resultValue.Interface()
}

func QueryGroup(data interface{}, groupType string, groupFuctor interface{}) interface{} {
	//解析配置
	data = QuerySort(data, groupType)
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataValueLen := dataValue.Len()
	dataElemType := dataType.Elem()
	dataCompare := getQueryCompares(dataElemType, groupType)

	groupFuctorValue := reflect.ValueOf(groupFuctor)
	groupFuctorType := groupFuctorValue.Type()

	var resultValue reflect.Value
	resultType := groupFuctorType.Out(0)
	if resultType.Kind() == reflect.Slice {
		resultValue = reflect.MakeSlice(resultType, 0, 0)
	} else {
		resultValue = reflect.MakeSlice(reflect.SliceOf(resultType), 0, 0)
	}

	//开始group
	for i := 0; i != dataValueLen; {
		singleDataValue := dataValue.Index(i)
		j := i
		for i++; i != dataValueLen; i++ {
			singleRightDataValue := dataValue.Index(i)
			isSame := true
			for _, singleDataCompare := range dataCompare {
				if singleDataCompare(singleDataValue, singleRightDataValue) != 0 {
					isSame = false
					break
				}
			}
			if !isSame {
				break
			}
		}
		singleResult := groupFuctorValue.Call([]reflect.Value{dataValue.Slice(j, i)})[0]
		if singleResult.Kind() == reflect.Slice {
			resultValue = reflect.AppendSlice(resultValue, singleResult)
		} else {
			resultValue = reflect.Append(resultValue, singleResult)
		}
	}
	return resultValue.Interface()
}

func analyseJoin(joinType string) (string, string) {
	joinTypeArray := strings.Split(joinType, "=")
	leftJoinType := strings.Trim(joinTypeArray[0], " ")
	rightJoinType := strings.Trim(joinTypeArray[1], " ")
	return leftJoinType, rightJoinType
}

func analyseSort(sortType string) (result1 []string, result2 []bool) {
	sortTypeArray := strings.Split(sortType, ",")
	for _, singleSortTypeArray := range sortTypeArray {
		singleSortTypeArrayTemp := strings.Split(singleSortTypeArray, " ")
		singleSortTypeArray := []string{}
		for _, singleSort := range singleSortTypeArrayTemp {
			singleSort = strings.Trim(singleSort, " ")
			if singleSort == "" {
				continue
			}
			singleSortTypeArray = append(singleSortTypeArray, singleSort)
		}
		var singleSortName string
		var singleSortType bool
		if len(singleSortTypeArray) >= 2 {
			singleSortName = singleSortTypeArray[0]
			singleSortType = (strings.ToLower(strings.Trim(singleSortTypeArray[1], " ")) == "asc")
		} else {
			singleSortName = singleSortTypeArray[0]
			singleSortType = true
		}
		result1 = append(result1, singleSortName)
		result2 = append(result2, singleSortType)
	}
	return result1, result2
}

func getQueryCompares(dataType reflect.Type, sortTypeStr string) []queryCompare {
	sortName, sortType := analyseSort(sortTypeStr)
	targetCompare := []queryCompare{}
	for index, singleSortName := range sortName {
		singleSortType := sortType[index]
		singleCompare := getQueryCompare(dataType, singleSortName)
		if !singleSortType {
			singleTempCompare := singleCompare
			singleCompare = func(left reflect.Value, right reflect.Value) int {
				return singleTempCompare(right, left)
			}
		}
		targetCompare = append(targetCompare, singleCompare)
	}
	return targetCompare
}

func getSingleQueryCompare(fieldType reflect.Type) queryCompare {
	typeKind := GetTypeKind(fieldType)
	if typeKind == TypeKind.BOOL {
		return func(left reflect.Value, right reflect.Value) int {
			leftBool := left.Bool()
			rightBool := right.Bool()
			if leftBool == rightBool {
				return 0
			} else if leftBool == false {
				return -1
			} else {
				return 1
			}
		}
	} else if typeKind == TypeKind.INT {
		return func(left reflect.Value, right reflect.Value) int {
			leftInt := left.Int()
			rightInt := right.Int()
			if leftInt < rightInt {
				return -1
			} else if leftInt > rightInt {
				return 1
			} else {
				return 0
			}
		}
	} else if typeKind == TypeKind.UINT {
		return func(left reflect.Value, right reflect.Value) int {
			leftUint := left.Uint()
			rightUint := right.Uint()
			if leftUint < rightUint {
				return -1
			} else if leftUint > rightUint {
				return 1
			} else {
				return 0
			}
		}
	} else if typeKind == TypeKind.FLOAT {
		return func(left reflect.Value, right reflect.Value) int {
			leftFloat := left.Float()
			rightFloat := right.Float()
			if leftFloat < rightFloat {
				return -1
			} else if leftFloat > rightFloat {
				return 1
			} else {
				return 0
			}
		}
	} else if typeKind == TypeKind.STRING {
		return func(left reflect.Value, right reflect.Value) int {
			leftString := left.String()
			rightString := right.String()
			if leftString < rightString {
				return -1
			} else if leftString > rightString {
				return 1
			} else {
				return 0
			}
		}
	} else if typeKind == TypeKind.STRUCT && fieldType == reflect.TypeOf(time.Time{}) {
		return func(left reflect.Value, right reflect.Value) int {
			leftTime := left.Interface().(time.Time)
			rightTime := right.Interface().(time.Time)
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

func getQueryCompare(dataType reflect.Type, name string) queryCompare {
	field, ok := getFieldByName(dataType, name)
	if !ok {
		panic(dataType.Name() + " has not name " + name)
	}
	fieldIndex := field.Index
	fieldType := field.Type
	compare := getSingleQueryCompare(fieldType)
	return func(left reflect.Value, right reflect.Value) int {
		return compare(left.FieldByIndex(fieldIndex), right.FieldByIndex(fieldIndex))
	}
}

//扩展类函数
func QueryLeftJoin(leftData interface{}, rightData interface{}, joinType string, joinFuctor interface{}) interface{} {
	return QueryJoin(leftData, rightData, "left", joinType, joinFuctor)
}

func QueryRightJoin(leftData interface{}, rightData interface{}, joinType string, joinFuctor interface{}) interface{} {
	return QueryJoin(leftData, rightData, "right", joinType, joinFuctor)
}

func QueryInnerJoin(leftData interface{}, rightData interface{}, joinType string, joinFuctor interface{}) interface{} {
	return QueryJoin(leftData, rightData, "inner", joinType, joinFuctor)
}

func QueryOuterJoin(leftData interface{}, rightData interface{}, joinType string, joinFuctor interface{}) interface{} {
	return QueryJoin(leftData, rightData, "outer", joinType, joinFuctor)
}

func QuerySum(data interface{}) interface{} {
	dataType := reflect.TypeOf(data).Elem()
	if dataType.Kind() == reflect.Int {
		return QueryReduce(data, func(sum int, single int) int {
			return sum + single
		}, 0)
	} else if dataType.Kind() == reflect.Float32 {
		return QueryReduce(data, func(sum float32, single float32) float32 {
			return sum + single
		}, (float32)(0.0))
	} else if dataType.Kind() == reflect.Float64 {
		return QueryReduce(data, func(sum float64, single float64) float64 {
			return sum + single
		}, 0.0)
	} else {
		panic("invalid type " + dataType.String())
	}
}

func QueryMax(data interface{}) interface{} {
	dataType := reflect.TypeOf(data).Elem()
	if dataType.Kind() == reflect.Int {
		return QueryReduce(data, func(max int, single int) int {
			if single > max {
				return single
			} else {
				return max
			}
		}, math.MinInt32)
	} else if dataType.Kind() == reflect.Float32 {
		return QueryReduce(data, func(max float32, single float32) float32 {
			if single > max {
				return single
			} else {
				return max
			}
		}, math.SmallestNonzeroFloat32)
	} else if dataType.Kind() == reflect.Float64 {
		return QueryReduce(data, func(max float64, single float64) float64 {
			if single > max {
				return single
			} else {
				return max
			}
		}, math.SmallestNonzeroFloat64)
	} else {
		panic("invalid type " + dataType.String())
	}
}

func QueryMin(data interface{}) interface{} {
	dataType := reflect.TypeOf(data).Elem()
	if dataType.Kind() == reflect.Int {
		return QueryReduce(data, func(min int, single int) int {
			if single < min {
				return single
			} else {
				return min
			}
		}, math.MaxInt32)
	} else if dataType.Kind() == reflect.Float32 {
		return QueryReduce(data, func(min float32, single float32) float32 {
			if single < min {
				return single
			} else {
				return min
			}
		}, math.MaxFloat32)
	} else if dataType.Kind() == reflect.Float64 {
		return QueryReduce(data, func(min float64, single float64) float64 {
			if single < min {
				return single
			} else {
				return min
			}
		}, math.MaxFloat64)
	} else {
		panic("invalid type " + dataType.String())
	}
}

func QueryColumn(data interface{}, column string) interface{} {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type().Elem()
	dataLen := dataValue.Len()
	column = strings.Trim(column, " ")
	dataFieldIndexStruct, ok := getFieldByName(dataType, column)
	if !ok {
		panic(dataType.Name() + " has no field " + column)
	}
	dataFieldIndex := dataFieldIndexStruct.Index

	resultValue := reflect.MakeSlice(reflect.SliceOf(dataFieldIndexStruct.Type), dataLen, dataLen)

	for i := 0; i != dataLen; i++ {
		singleDataValue := dataValue.Index(i)
		singleResultValue := singleDataValue.FieldByIndex(dataFieldIndex)
		resultValue.Index(i).Set(singleResultValue)
	}
	return resultValue.Interface()
}

func QueryReverse(data interface{}) interface{} {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataLen := dataValue.Len()
	result := reflect.MakeSlice(dataType, dataLen, dataLen)

	for i := 0; i != dataLen; i++ {
		result.Index(dataLen - i - 1).Set(dataValue.Index(i))
	}
	return result.Interface()
}

func QueryDistinct(data interface{}, columnNames string) interface{} {
	//提取信息
	name := Explode(columnNames, ",")
	nameInfo := []arrayColumnMapInfo{}
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type().Elem()
	for _, singleName := range name {
		singleField, ok := getFieldByName(dataType, singleName)
		if !ok {
			panic(dataType.Name() + " struct has not field " + singleName)
		}
		nameInfo = append(nameInfo, arrayColumnMapInfo{
			Index: singleField.Index,
			Type:  singleField.Type,
		})
	}

	//整合map
	existsMap := map[interface{}]bool{}
	result := reflect.MakeSlice(dataValue.Type(), 0, 0)
	dataLen := dataValue.Len()
	for i := 0; i != dataLen; i++ {
		singleValue := dataValue.Index(i)
		newData := reflect.New(dataType).Elem()
		for _, singleNameInfo := range nameInfo {
			singleField := singleValue.FieldByIndex(singleNameInfo.Index)
			newData.FieldByIndex(singleNameInfo.Index).Set(singleField)
		}
		newDataValue := newData.Interface()
		_, isExist := existsMap[newDataValue]
		if isExist {
			continue
		}
		result = reflect.Append(result, singleValue)
		existsMap[newDataValue] = true
	}
	return result.Interface()
}
