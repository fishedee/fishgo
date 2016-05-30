package web

import (
	"reflect"
	"sync"
	"unsafe"
)

var (
	iocMutex     sync.RWMutex
	iocType      = map[reflect.Type][]uintptr{}
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
			singleResultArray := getIocTypeIndexInner(singleFiled.Type)
			for _, singleResult := range singleResultArray {
				data := append([]int{i}, singleResult...)
				result = append(result, data)
			}
		}
	}
	return result
}

func getIocTypeIndex(target reflect.Type) []uintptr {
	iocMutex.RLock()
	result, ok := iocType[target]
	iocMutex.RUnlock()

	if ok {
		return result
	}
	index := getIocTypeIndexInner(target)
	newData := reflect.New(target).Elem()
	newDataBasicAddr := newData.UnsafeAddr()
	result = nil
	for _, singleIndex := range index {
		singleNewData := newData.FieldByIndex(singleIndex)
		singleNewDataAddr := singleNewData.UnsafeAddr()
		result = append(result, singleNewDataAddr-newDataBasicAddr)
	}

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
	targetAddr := target.UnsafeAddr()
	for _, singleIndex := range typeIndex {
		var pointer **Basic = (**Basic)(unsafe.Pointer(targetAddr + singleIndex))
		*pointer = basic
	}
}
