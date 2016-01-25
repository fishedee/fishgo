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
	resultReduceValue := reflect.ValueOf(resultReduce)

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
	leftDataJoinStruct, ok := leftDataElemType.FieldByName(leftJoinType)
	if !ok {
		panic(leftDataElemType.Name() + " has no field " + leftJoinType)
	}
	leftDataJoin := leftDataJoinStruct.Index

	rightData = QuerySort(rightData, rightJoinType+" asc")
	rightDataValue := reflect.ValueOf(rightData)
	rightDataType := rightDataValue.Type()
	rightDataElemType := rightDataType.Elem()
	rightDataValueLen := rightDataValue.Len()
	rightDataJoinStruct, ok := rightDataElemType.FieldByName(rightJoinType)
	if !ok {
		panic(rightDataElemType.Name() + " has no field " + rightJoinType)
	}
	rightDataJoin := rightDataJoinStruct.Index

	joinFuctorValue := reflect.ValueOf(joinFuctor)
	joinFuctorType := joinFuctorValue.Type()
	joinCompare := getSingleQueryCompare(leftDataJoinStruct.Type)
	resultValue := reflect.MakeSlice(reflect.SliceOf(joinFuctorType.Out(0)), 0, 0)

	rightHaveJoin := make([]bool, rightDataValueLen, rightDataValueLen)
	joinPlace = strings.ToLower(joinPlace)

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
	resultValue := reflect.MakeSlice(groupFuctorType.Out(0), 0, 0)

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
		resultValue = reflect.AppendSlice(resultValue, singleResult)
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
	if fieldType.Kind() == reflect.Bool {
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
	} else if fieldType.Kind() == reflect.Int {
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
	} else if fieldType.Kind() == reflect.Float32 {
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
	} else if fieldType.Kind() == reflect.Struct && fieldType == reflect.TypeOf(time.Time{}) {
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
	field, ok := dataType.FieldByName(name)
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
	return QueryReduce(data, func(sum int, single int) int {
		return sum + single
	}, 0)
}

func QueryMax(data interface{}) interface{} {
	return QueryReduce(data, func(max int, single int) int {
		if single > max {
			return single
		} else {
			return max
		}
	}, math.MinInt32)
}

func QueryMin(data interface{}) interface{} {
	return QueryReduce(data, func(min int, single int) int {
		if single < min {
			return single
		} else {
			return min
		}
	}, math.MaxInt32)
}

func QueryColumn(data interface{}, column string) interface{} {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type().Elem()
	dataLen := dataValue.Len()
	dataFieldIndexStruct, ok := dataType.FieldByName(column)
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
