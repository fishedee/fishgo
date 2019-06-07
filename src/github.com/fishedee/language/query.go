package language

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"time"
)

//基础类函数QuerySelect
type QuerySelectMacroHandler func(data interface{}, selectFunctor interface{}) interface{}

func QuerySelectMacroRegister(data interface{}, selectFunctor interface{}, handler QuerySelectMacroHandler) {
	id := registerQueryTypeId([]string{reflect.TypeOf(data).String(), reflect.TypeOf(selectFunctor).String()})
	querySelectMacroMapper[id] = handler
}

func QuerySelectReflect(data interface{}, selectFuctor interface{}) interface{} {
	dataValue := reflect.ValueOf(data)
	dataLen := dataValue.Len()

	selectFuctorValue := reflect.ValueOf(selectFuctor)
	selectFuctorType := selectFuctorValue.Type()
	selectFuctorOuterType := selectFuctorType.Out(0)
	resultType := reflect.SliceOf(selectFuctorOuterType)
	resultValue := reflect.MakeSlice(resultType, dataLen, dataLen)
	callArgument := []reflect.Value{reflect.Value{}}

	for i := 0; i != dataLen; i++ {
		singleDataValue := dataValue.Index(i)
		callArgument[0] = singleDataValue
		singleResultValue := selectFuctorValue.Call(callArgument)[0]
		resultValue.Index(i).Set(singleResultValue)
	}
	return resultValue.Interface()
}

func QuerySelect(data interface{}, selectFunctor interface{}) interface{} {
	id := getQueryTypeId([]string{reflect.TypeOf(data).String(), reflect.TypeOf(selectFunctor).String()})
	handler, isExist := querySelectMacroMapper[id]
	if isExist {
		return handler(data, selectFunctor)
	} else {
		queryReflectWarn("QuerySelect")
		return QuerySelectReflect(data, selectFunctor)
	}
}

//基础类函数QueryWhere
type QueryWhereMacroHandler func(data interface{}, whereFunctor interface{}) interface{}

func QueryWhereMacroRegister(data interface{}, whereFunctor interface{}, handler QueryWhereMacroHandler) {
	id := registerQueryTypeId([]string{reflect.TypeOf(data).String(), reflect.TypeOf(whereFunctor).String()})
	queryWhereMacroMapper[id] = handler
}

func QueryWhereReflect(data interface{}, whereFuctor interface{}) interface{} {
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

func QueryWhere(data interface{}, whereFuctor interface{}) interface{} {
	id := getQueryTypeId([]string{reflect.TypeOf(data).String(), reflect.TypeOf(whereFuctor).String()})
	handler, isExist := queryWhereMacroMapper[id]
	if isExist {
		return handler(data, whereFuctor)
	} else {
		queryReflectWarn("QueryWhere")
		return QueryWhereReflect(data, whereFuctor)
	}
}

//基础类函数QuerySort
type querySortInterface struct {
	lenHandler  func() int
	lessHandler func(i int, j int) bool
	swapHandler func(i int, j int)
}

func (this *querySortInterface) Len() int {
	return this.lenHandler()
}

func (this *querySortInterface) Less(i int, j int) bool {
	return this.lessHandler(i, j)
}

func (this *querySortInterface) Swap(i int, j int) {
	this.swapHandler(i, j)
}

func QuerySortInternal(length int, lessHandler func(i, j int) int, swapHandler func(i, j int)) {
	sort.Stable(&querySortInterface{
		lenHandler: func() int {
			return length
		},
		lessHandler: func(i int, j int) bool {
			return lessHandler(i, j) < 0
		},
		swapHandler: swapHandler,
	})

}

type QuerySortMacroHandler func(data interface{}, sortType string) interface{}

func QuerySortMacroRegister(data interface{}, sortType string, handler QuerySortMacroHandler) {
	id := registerQueryTypeId([]string{reflect.TypeOf(data).String(), sortType})
	querySortMacroMapper[id] = handler
}

func QuerySortReflect(data interface{}, sortType string) interface{} {
	//拷贝一份
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataElemType := dataType.Elem()
	dataValueLen := dataValue.Len()

	dataResult := reflect.MakeSlice(dataType, dataValueLen, dataValueLen)
	reflect.Copy(dataResult, dataValue)

	//排序
	targetCompares := getQueryExtractAndCompares(dataElemType, sortType)
	targetCompare := combineQueryCompare(targetCompares)
	result := dataResult.Interface()
	swapper := reflect.Swapper(result)

	QuerySortInternal(dataValueLen, func(i int, j int) int {
		left := dataResult.Index(i)
		right := dataResult.Index(j)
		return targetCompare(left, right)
	}, swapper)

	return result
}

func QuerySort(data interface{}, sortType string) interface{} {
	id := getQueryTypeId([]string{reflect.TypeOf(data).String(), sortType})
	handler, isExist := querySortMacroMapper[id]
	if isExist {
		return handler(data, sortType)
	} else {
		queryReflectWarn("QuerySort")
		return QuerySortReflect(data, sortType)
	}
}

//基础类函数QueryJoin
type QueryJoinMacroHandler func(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFuctor interface{}) interface{}

func QueryJoinMacroRegister(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFuctor interface{}, handler QueryJoinMacroHandler) {
	id := registerQueryTypeId([]string{reflect.TypeOf(leftData).String(), reflect.TypeOf(rightData).String(), joinPlace, joinType, reflect.TypeOf(joinFuctor).String()})
	queryJoinMacroMapper[id] = handler
}

func QueryJoinReflect(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFuctor interface{}) interface{} {
	//解析配置
	leftJoinType, rightJoinType := analyseJoin(joinType)

	leftDataValue := reflect.ValueOf(leftData)
	leftDataType := leftDataValue.Type()
	leftDataElemType := leftDataType.Elem()
	leftDataValueLen := leftDataValue.Len()
	leftDataJoinType, leftDataJoinExtract := getQueryExtract(leftDataElemType, leftJoinType)

	rightDataValue := reflect.ValueOf(rightData)
	rightDataType := rightDataValue.Type()
	rightDataElemType := rightDataType.Elem()
	rightDataValueLen := rightDataValue.Len()
	_, rightDataJoinExtract := getQueryExtract(rightDataElemType, rightJoinType)

	joinFuctorValue := reflect.ValueOf(joinFuctor)
	joinFuctorType := joinFuctorValue.Type()

	resultValue := reflect.MakeSlice(reflect.SliceOf(joinFuctorType.Out(0)), 0, 0)

	//执行join
	emptyLeftValue := reflect.New(leftDataElemType).Elem()
	emptyRightValue := reflect.New(rightDataElemType).Elem()
	joinPlace = strings.Trim(strings.ToLower(joinPlace), " ")

	nextData := make([]int, rightDataValueLen, rightDataValueLen)
	mapDataNext := reflect.MakeMapWithSize(reflect.MapOf(leftDataJoinType, reflect.TypeOf(1)), rightDataValueLen)
	mapDataFirst := reflect.MakeMapWithSize(reflect.MapOf(leftDataJoinType, reflect.TypeOf(1)), rightDataValueLen)
	tempValueInt := reflect.New(reflect.TypeOf(1)).Elem()

	for i := 0; i != rightDataValueLen; i++ {
		tempValueInt.SetInt(int64(i))
		fieldValue := rightDataJoinExtract(rightDataValue.Index(i))
		lastNextIndex := mapDataNext.MapIndex(fieldValue)
		if lastNextIndex.IsValid() {
			nextData[int(lastNextIndex.Int())] = i
		} else {
			mapDataFirst.SetMapIndex(fieldValue, tempValueInt)
		}
		nextData[i] = -1
		mapDataNext.SetMapIndex(fieldValue, tempValueInt)
	}
	rightHaveJoin := make([]bool, rightDataValueLen, rightDataValueLen)
	for i := 0; i != leftDataValueLen; i++ {
		leftValue := leftDataValue.Index(i)
		fieldValue := leftDataJoinExtract(leftDataValue.Index(i))
		rightIndex := mapDataFirst.MapIndex(fieldValue)
		if rightIndex.IsValid() {
			//找到右值
			j := int(rightIndex.Int())
			for nextData[j] != -1 {
				singleResult := joinFuctorValue.Call([]reflect.Value{leftValue, rightDataValue.Index(j)})[0]
				resultValue = reflect.Append(resultValue, singleResult)
				rightHaveJoin[j] = true
				j = nextData[j]
			}
			singleResult := joinFuctorValue.Call([]reflect.Value{leftValue, rightDataValue.Index(j)})[0]
			resultValue = reflect.Append(resultValue, singleResult)
			rightHaveJoin[j] = true
		} else {
			//找不到右值
			if joinPlace == "left" || joinPlace == "outer" {
				singleResult := joinFuctorValue.Call([]reflect.Value{leftValue, emptyRightValue})[0]
				resultValue = reflect.Append(resultValue, singleResult)
			}
		}
	}
	if joinPlace == "right" || joinPlace == "outer" {
		for j := 0; j != rightDataValueLen; j++ {
			if rightHaveJoin[j] {
				continue
			}
			singleResult := joinFuctorValue.Call([]reflect.Value{emptyLeftValue, rightDataValue.Index(j)})[0]
			resultValue = reflect.Append(resultValue, singleResult)
		}
	}
	return resultValue.Interface()
}

func QueryJoin(leftData interface{}, rightData interface{}, joinPlace string, joinType string, joinFuctor interface{}) interface{} {
	id := getQueryTypeId([]string{reflect.TypeOf(leftData).String(), reflect.TypeOf(rightData).String(), joinPlace, joinType, reflect.TypeOf(joinFuctor).String()})
	handler, isExist := queryJoinMacroMapper[id]
	if isExist {
		return handler(leftData, rightData, joinPlace, joinType, joinFuctor)
	} else {
		queryReflectWarn("QueryJoin")
		return QueryJoinReflect(leftData, rightData, joinPlace, joinType, joinFuctor)
	}
}

//基础类函数 QueryGroup
type QueryGroupMacroHandler func(data interface{}, groupType string, groupFunctor interface{}) interface{}

func QueryGroupMacroRegister(data interface{}, groupType string, groupFunctor interface{}, handler QueryGroupMacroHandler) {
	id := registerQueryTypeId([]string{reflect.TypeOf(data).String(), groupType, reflect.TypeOf(groupFunctor).String()})
	queryGroupMacroMapper[id] = handler
}

type queryGroupWalkHandler func(data reflect.Value)

func queryGroupWalkReflect(data interface{}, groupType string, groupWalkHandler queryGroupWalkHandler) {
	//解析输入数据
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type()
	dataElemType := dataType.Elem()
	dataValueLen := dataValue.Len()

	//分组操作
	groupType = strings.Trim(groupType, " ")
	dataFieldType, dataFieldExtract := getQueryExtract(dataElemType, groupType)
	findMap := reflect.MakeMapWithSize(reflect.MapOf(dataFieldType, reflect.TypeOf(1)), dataValueLen)
	bufferData := reflect.MakeSlice(dataType, dataValueLen, dataValueLen)
	tempValueInt := reflect.New(reflect.TypeOf(1)).Elem()

	nextData := make([]int, dataValueLen, dataValueLen)
	for i := 0; i != dataValueLen; i++ {
		fieldValue := dataFieldExtract(dataValue.Index(i))
		lastIndex := findMap.MapIndex(fieldValue)
		if lastIndex.IsValid() {
			nextData[int(lastIndex.Int())] = i
		}
		nextData[i] = -1
		tempValueInt.SetInt(int64(i))
		findMap.SetMapIndex(fieldValue, tempValueInt)
	}
	k := 0
	for i := 0; i != dataValueLen; i++ {
		j := i
		if nextData[j] == 0 {
			continue
		}
		kbegin := k
		for nextData[j] != -1 {
			nextJ := nextData[j]
			bufferData.Index(k).Set(dataValue.Index(j))
			nextData[j] = 0
			j = nextJ
			k++
		}
		bufferData.Index(k).Set(dataValue.Index(j))
		k++
		nextData[j] = 0
		groupWalkHandler(bufferData.Slice(kbegin, k))
	}
}

func QueryGroupReflect(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	groupFuctorValue := reflect.ValueOf(groupFunctor)
	groupFuctorType := groupFuctorValue.Type()

	//解析输入数据
	dataValueLen := reflect.ValueOf(data).Len()

	//计算最终数据
	var resultValue reflect.Value
	resultType := groupFuctorType.Out(0)
	if resultType.Kind() == reflect.Slice {
		resultValue = reflect.MakeSlice(resultType, 0, dataValueLen)
	} else {
		resultValue = reflect.MakeSlice(reflect.SliceOf(resultType), 0, dataValueLen)
	}

	//执行分组操作
	queryGroupWalkReflect(data, groupType, func(data reflect.Value) {
		singleResult := groupFuctorValue.Call([]reflect.Value{data})[0]
		if singleResult.Kind() == reflect.Slice {
			resultValue = reflect.AppendSlice(resultValue, singleResult)
		} else {
			resultValue = reflect.Append(resultValue, singleResult)
		}
	})

	return resultValue.Interface()
}

func QueryGroup(data interface{}, groupType string, groupFunctor interface{}) interface{} {
	id := getQueryTypeId([]string{reflect.TypeOf(data).String(), groupType, reflect.TypeOf(groupFunctor).String()})
	handler, isExist := queryGroupMacroMapper[id]
	if isExist {
		return handler(data, groupType, groupFunctor)
	} else {
		queryReflectWarn("QueryGroup")
		return QueryGroupReflect(data, groupType, groupFunctor)
	}
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

func getQueryExtractAndCompares(dataType reflect.Type, sortTypeStr string) []queryCompare {
	sortName, sortType := analyseSort(sortTypeStr)
	targetCompare := []queryCompare{}
	for index, singleSortName := range sortName {
		singleSortType := sortType[index]
		singleCompare := getQueryExtractAndCompare(dataType, singleSortName)
		if !singleSortType {
			//逆序
			singleTempCompare := singleCompare
			singleCompare = func(left reflect.Value, right reflect.Value) int {
				return singleTempCompare(right, left)
			}
		}
		targetCompare = append(targetCompare, singleCompare)
	}
	return targetCompare
}

func getQueryCompare(fieldType reflect.Type) queryCompare {
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
		if fieldType == reflect.TypeOf(Decimal("")) {
			return func(left reflect.Value, right reflect.Value) int {
				leftDecimal := left.Interface().(Decimal)
				rightDecimal := right.Interface().(Decimal)
				return leftDecimal.Cmp(rightDecimal)
			}
		} else {
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

type queryCompare func(reflect.Value, reflect.Value) int

func combineQueryCompare(targetCompare []queryCompare) queryCompare {
	return func(left reflect.Value, right reflect.Value) int {
		for _, singleCompare := range targetCompare {
			compareResult := singleCompare(left, right)
			if compareResult < 0 {
				return -1
			} else if compareResult > 0 {
				return 1
			}
		}
		return 0
	}
}

type queryExtract func(reflect.Value) reflect.Value

func getQueryExtract(dataType reflect.Type, name string) (reflect.Type, queryExtract) {
	if name == "." {
		return dataType, func(v reflect.Value) reflect.Value {
			return v
		}
	} else {
		field, ok := getFieldByName(dataType, name)
		if !ok {
			panic(dataType.Name() + " has not name " + name)
		}
		fieldIndex := field.Index
		fieldType := field.Type
		return fieldType, func(v reflect.Value) reflect.Value {
			return v.FieldByIndex(fieldIndex)
		}
	}
}

func getQueryExtractAndCompare(dataType reflect.Type, name string) queryCompare {
	fieldType, extract := getQueryExtract(dataType, name)
	compare := getQueryCompare(fieldType)
	return func(left reflect.Value, right reflect.Value) int {
		return compare(extract(left), extract(right))
	}
}

//扩展类函数 QueryColumn
type QueryColumnMacroHandler func(data interface{}, column string) interface{}

func QueryColumnMacroRegister(data interface{}, column string, handler QueryColumnMacroHandler) {
	id := registerQueryTypeId([]string{reflect.TypeOf(data).String(), column})
	queryColumnMacroMapper[id] = handler
}

func QueryColumnReflect(data interface{}, column string) interface{} {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type().Elem()
	dataLen := dataValue.Len()
	column = strings.Trim(column, " ")
	dataFieldType, dataFieldExtract := getQueryExtract(dataType, column)

	resultValue := reflect.MakeSlice(reflect.SliceOf(dataFieldType), dataLen, dataLen)

	for i := 0; i != dataLen; i++ {
		singleDataValue := dataValue.Index(i)
		singleResultValue := dataFieldExtract(singleDataValue)
		resultValue.Index(i).Set(singleResultValue)
	}
	return resultValue.Interface()
}

func QueryColumn(data interface{}, column string) interface{} {
	id := getQueryTypeId([]string{reflect.TypeOf(data).String(), column})
	handler, isExist := queryColumnMacroMapper[id]
	if isExist {
		return handler(data, column)
	} else {
		queryReflectWarn("QueryColumn")
		return QueryColumnReflect(data, column)
	}
}

//扩展类函数 QueryColumnMap
type QueryColumnMapMacroHandler func(data interface{}, column string) interface{}

func QueryColumnMapMacroRegister(data interface{}, column string, handler QueryColumnMapMacroHandler) {
	id := registerQueryTypeId([]string{reflect.TypeOf(data).String(), column})
	queryColumnMapMacroMapper[id] = handler
}

func QueryColumnMapReflect(data interface{}, column string) interface{} {
	column = strings.Trim(column, " ")
	if len(column) >= 2 && column[0:2] == "[]" {
		column = column[2:]
		return queryColumnMapReflectSlice(data, column)
	} else {
		return queryColumnMapReflectSingle(data, column)
	}
}

func queryColumnMapReflectSlice(data interface{}, column string) interface{} {
	dataValue := reflect.ValueOf(data)
	dataValueType := dataValue.Type()
	dataType := dataValue.Type().Elem()
	dataLen := dataValue.Len()
	column = strings.Trim(column, " ")
	dataFieldType, dataFieldExtract := getQueryExtract(dataType, column)

	resultValue := reflect.MakeMapWithSize(reflect.MapOf(dataFieldType, dataValueType), dataLen)

	queryGroupWalkReflect(data, column, func(group reflect.Value) {
		singleResultValue := dataFieldExtract(group.Index(0))
		resultValue.SetMapIndex(singleResultValue, group)
	})
	return resultValue.Interface()
}

func queryColumnMapReflectSingle(data interface{}, column string) interface{} {
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type().Elem()
	dataLen := dataValue.Len()
	column = strings.Trim(column, " ")
	dataFieldType, dataFieldExtract := getQueryExtract(dataType, column)

	resultValue := reflect.MakeMapWithSize(reflect.MapOf(dataFieldType, dataType), dataLen)
	for i := dataLen - 1; i >= 0; i-- {
		singleDataValue := dataValue.Index(i)
		singleResultValue := dataFieldExtract(singleDataValue)
		resultValue.SetMapIndex(singleResultValue, singleDataValue)
	}
	return resultValue.Interface()
}

func QueryColumnMap(data interface{}, column string) interface{} {
	id := getQueryTypeId([]string{reflect.TypeOf(data).String(), column})
	handler, isExist := queryColumnMapMacroMapper[id]
	if isExist {
		return handler(data, column)
	} else {
		queryReflectWarn("QueryColumnMap")
		return QueryColumnMapReflect(data, column)
	}
}

//扩展类函数 QueryLeftJoin,QueryRightJoin,QueryInnerJoin,QueryOuterJoin
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

//扩展累函数 QueryCombine
type QueryCombineMacroHandler func(leftData interface{}, rightData interface{}, combineFuctor interface{}) interface{}

func QueryCombineMacroRegister(leftData interface{}, rightData interface{}, combineFuctor interface{}, handler QueryCombineMacroHandler) {
	id := registerQueryTypeId([]string{reflect.TypeOf(leftData).String(), reflect.TypeOf(rightData).String(), reflect.TypeOf(combineFuctor).String()})
	queryCombineMacroMapper[id] = handler
}

func QueryCombineReflect(leftData interface{}, rightData interface{}, combineFuctor interface{}) interface{} {
	leftValue := reflect.ValueOf(leftData)
	rightValue := reflect.ValueOf(rightData)
	if leftValue.Len() != rightValue.Len() {
		panic(fmt.Sprintf("len dos not equal %v != %v", leftValue.Len(), rightValue.Len()))
	}
	dataLen := leftValue.Len()
	combineFuctorValue := reflect.ValueOf(combineFuctor)
	resultType := combineFuctorValue.Type().Out(0)
	result := reflect.MakeSlice(reflect.SliceOf(resultType), dataLen, dataLen)
	for i := 0; i != dataLen; i++ {
		singleResultValue := combineFuctorValue.Call([]reflect.Value{leftValue.Index(i), rightValue.Index(i)})
		result.Index(i).Set(singleResultValue[0])
	}
	return result.Interface()
}

func QueryCombine(leftData interface{}, rightData interface{}, combineFuctor interface{}) interface{} {
	id := getQueryTypeId([]string{reflect.TypeOf(leftData).String(), reflect.TypeOf(rightData).String(), reflect.TypeOf(combineFuctor).String()})
	handler, isExist := queryCombineMacroMapper[id]
	if isExist {
		return handler(leftData, rightData, combineFuctor)
	} else {
		queryReflectWarn("QueryCombine")
		return QueryCombineReflect(leftData, rightData, combineFuctor)
	}
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
	extractInfo := []queryExtract{}
	dataValue := reflect.ValueOf(data)
	dataType := dataValue.Type().Elem()
	for _, singleName := range name {
		_, extract := getQueryExtract(dataType, singleName)
		extractInfo = append(extractInfo, extract)
	}

	//整合map
	existsMap := map[interface{}]bool{}
	result := reflect.MakeSlice(dataValue.Type(), 0, 0)
	dataLen := dataValue.Len()
	for i := 0; i != dataLen; i++ {
		singleValue := dataValue.Index(i)
		newData := reflect.New(dataType).Elem()
		for _, singleExtract := range extractInfo {
			singleField := singleExtract(singleValue)
			singleExtract(newData).Set(singleField)
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

func registerQueryTypeId(data []string) int64 {
	var result int64
	for _, m := range data {
		id, isExist := queryTypeIdMapper[m]
		if isExist == false {
			id = int64(len(queryTypeIdMapper)) + 1
			queryTypeIdMapper[m] = id
		}
		result = result<<10 + id
	}
	return result
}

func getQueryTypeId(data []string) int64 {
	var result int64
	for _, m := range data {
		id, isExist := queryTypeIdMapper[m]
		if isExist == false {
			return -1
		}
		result = result<<10 + id
	}
	return result
}

func queryReflectWarn(funcName string) {
	if queryReflectWarning {
		_, file, line, _ := runtime.Caller(2)
		log.Printf("%s:%d use %v reflect version,you should use querygen to avoid this warning", file, line, funcName)
	}

}
func QueryReflectWarning(isWarning bool) {
	queryReflectWarning = isWarning
}

var (
	querySelectMacroMapper    = map[int64]QuerySelectMacroHandler{}
	queryWhereMacroMapper     = map[int64]QueryWhereMacroHandler{}
	querySortMacroMapper      = map[int64]QuerySortMacroHandler{}
	queryJoinMacroMapper      = map[int64]QueryJoinMacroHandler{}
	queryGroupMacroMapper     = map[int64]QueryGroupMacroHandler{}
	queryColumnMacroMapper    = map[int64]QueryColumnMacroHandler{}
	queryColumnMapMacroMapper = map[int64]QueryColumnMapMacroHandler{}
	queryCombineMacroMapper   = map[int64]QueryCombineMacroHandler{}
	queryTypeIdMapper         = map[string]int64{}
	queryReflectWarning       = false
)
