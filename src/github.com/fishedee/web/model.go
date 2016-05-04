package web

import (
	"reflect"
	"sync"
)

var modelInfo struct {
	mutex sync.RWMutex
	data  map[reflect.Type][][]int
}

func init() {
	modelInfo.data = map[reflect.Type][][]int{}
}

type ModelInterface interface {
	SetAppController(controller ControllerInterface)
	GetBasic() *Basic
}

type Model struct {
	*Basic
	appController ControllerInterface
}

func (this *Model) SetAppController(controller ControllerInterface) {
	this.appController = controller
	this.Basic = controller.GetBasic()
}

func (this *Model) GetBasic() *Basic {
	return this.Basic
}

func getSubModuleFromType(target reflect.Type) [][]int {
	modelInfo.mutex.RLock()
	result, ok := modelInfo.data[target]
	modelInfo.mutex.RUnlock()

	if ok {
		return result
	}
	result = getSubModuleFromTypeInner(target)

	modelInfo.mutex.Lock()
	modelInfo.data[target] = result
	modelInfo.mutex.Unlock()
	return result
}

func getSubModuleFromTypeInner(modelType reflect.Type) [][]int {
	result := [][]int{}
	numField := modelType.NumField()
	for i := 0; i != numField; i++ {
		singleFiled := modelType.Field(i)
		if singleFiled.PkgPath != "" {
			continue
		}
		if isFromModelType(singleFiled.Type) == false {
			continue
		}
		singleResultArray := getSubModuleFromType(singleFiled.Type)
		for _, singleResult := range singleResultArray {
			data := append([]int{i}, singleResult...)
			result = append(result, data)
		}
		if singleFiled.Anonymous == false {
			data := []int{i}
			result = append(result, data)
		}
	}
	return result
}

func isFromModelType(target reflect.Type) bool {
	var data *ModelInterface
	interfaceType := reflect.TypeOf(data).Elem()
	targetType := reflect.PtrTo(target)
	return targetType.Implements(interfaceType)
}

func initModelInner(model ModelInterface, controller ControllerInterface) {
	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()
	modelSubModels := getSubModuleFromType(modelType)
	model.SetAppController(controller)
	for _, singleModel := range modelSubModels {
		target := modelValue.FieldByIndex(singleModel).Addr().Interface().(ModelInterface)
		target.SetAppController(controller)
	}
}
func initModel(controller ControllerInterface) {
	initModelInner(controller, controller)
}
