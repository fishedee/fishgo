package web

import (
	"reflect"
	"sync"
)

var (
	iocMutex     sync.RWMutex
	iocType      = map[reflect.Type][][]int{}
	iocBasicType = reflect.TypeOf(&Basic{})
)

func getIocTypeIndexInner(modelType reflect.Type) [][]int {
	result := [][]int{}
	numField := modelType.NumField()
	for i := 0; i != numField; i++ {
		singleFiled := modelType.Field(i)
		if singleFiled.Type == iocBasicType {
			result = append(result, []int{i})
		} else if singleFiled.Type.Kind() == reflect.Struct &&
			singleFiled.PkgPath == "" {
			singleResultArray := getIocTypeIndex(singleFiled.Type)
			for _, singleResult := range singleResultArray {
				data := append([]int{i}, singleResult...)
				result = append(result, data)
			}
		}
	}
	return result
}

func getIocTypeIndex(target reflect.Type) [][]int {
	iocMutex.RLock()
	result, ok := iocType[target]
	iocMutex.RUnlock()

	if ok {
		return result
	}
	result = getIocTypeIndexInner(target)

	iocMutex.Lock()
	iocType[target] = result
	iocMutex.Unlock()

	return result
}

func injectIoc(target reflect.Value, basic *Basic) {
	for target.Kind() == reflect.Ptr {
		target = target.Elem()
	}
	typeIndex := getIocTypeIndex(target.Type())
	basicValue := reflect.ValueOf(basic)
	for _, singleIndex := range typeIndex {
		target.FieldByIndex(singleIndex).Set(basicValue)
	}
}
