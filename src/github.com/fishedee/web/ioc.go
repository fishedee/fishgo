package web

import (
	"reflect"
	"sync"
	"unsafe"
)

var (
	iocMutex          sync.RWMutex
	iocType           = map[reflect.Type][]uintptr{}
	iocTypeIndexMutex sync.RWMutex
	iocTypeIndex      = map[reflect.Type]*iocTypeIndexInfo{}
	iocBasicType      = reflect.TypeOf(&Basic{})
)

type iocTypeIndexInfoField struct {
	index    int
	children *iocTypeIndexInfo
}

type iocTypeIndexInfo struct {
	maxDepth int
	count    int
	fields   []iocTypeIndexInfoField
}

func getIocTypeIndexInner(modelType reflect.Type) *iocTypeIndexInfo {
	iocTypeIndexMutex.RLock()
	result, ok := iocTypeIndex[modelType]
	iocTypeIndexMutex.RUnlock()
	if ok {
		return result
	}

	result = &iocTypeIndexInfo{}
	numField := modelType.NumField()
	for i := 0; i != numField; i++ {
		singleFiled := modelType.Field(i)
		if singleFiled.Type == iocBasicType {
			result.fields = append(result.fields, iocTypeIndexInfoField{
				index:    i,
				children: nil,
			})
			if result.maxDepth < 1 {
				result.maxDepth = 1
			}
			result.count++
		} else if singleFiled.Type.Kind() == reflect.Struct &&
			singleFiled.PkgPath == "" {
			children := getIocTypeIndexInner(singleFiled.Type)
			if children.count == 0 {
				continue
			}
			result.fields = append(result.fields, iocTypeIndexInfoField{
				index:    i,
				children: children,
			})
			if result.maxDepth < children.maxDepth+1 {
				result.maxDepth = children.maxDepth + 1
			}
			result.count += children.count
		}
	}

	iocTypeIndexMutex.Lock()
	iocTypeIndex[modelType] = result
	iocTypeIndexMutex.Unlock()
	return result
}

func walkIocTypeIndexInfoDfs(info *iocTypeIndexInfo, currentIndex []int, currentPlace int, handler func([]int)) {
	for _, field := range info.fields {
		currentIndex[currentPlace] = field.index
		if field.children == nil {
			handler(currentIndex[0 : currentPlace+1])
		} else {
			walkIocTypeIndexInfoDfs(field.children, currentIndex, currentPlace+1, handler)
		}
	}
}

func getIocTypeIndex(target reflect.Type) []uintptr {
	iocMutex.RLock()
	result, ok := iocType[target]
	iocMutex.RUnlock()

	if ok {
		return result
	}
	info := getIocTypeIndexInner(target)
	newData := reflect.New(target).Elem()
	newDataBasicAddr := newData.UnsafeAddr()
	result = make([]uintptr, info.count, info.count)
	resultIndex := 0
	currentIndex := make([]int, info.maxDepth, info.maxDepth)
	walkIocTypeIndexInfoDfs(info, currentIndex, 0, func(singleIndex []int) {
		singleNewData := newData.FieldByIndex(singleIndex)
		singleNewDataAddr := singleNewData.UnsafeAddr()
		result[resultIndex] = singleNewDataAddr - newDataBasicAddr
		resultIndex++
	})

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
