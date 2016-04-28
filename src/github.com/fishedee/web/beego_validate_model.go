package web

import (
	_ "github.com/a"
	"github.com/astaxie/beego/context"
	"reflect"
	"sync"
)

var beegoValidateModelType reflect.Type
var beegoValidateModelInfo struct {
	mutex sync.RWMutex
	data  map[reflect.Type][][]int
}

func init() {
	beegoValidateModelType = reflect.TypeOf(BeegoValidateModel{})
	beegoValidateModelInfo.data = map[reflect.Type][][]int{}
}

type beegoValidateModelInterface interface {
	SetAppController(controller beegoValidateControllerInterface)
}

type BeegoValidateModel struct {
	*BeegoValidateBasic
	AppController interface{}
	Ctx           *context.Context
}

func (this *BeegoValidateModel) SetAppController(controller beegoValidateControllerInterface) {
	this.AppController = controller
	this.BeegoValidateBasic = controller.GetBasic()
	this.Ctx = this.BeegoValidateBasic.ctx
}

func getSubModuleFromType(target reflect.Type) [][]int {
	beegoValidateModelInfo.mutex.RLock()
	result, ok := beegoValidateModelInfo.data[target]
	beegoValidateModelInfo.mutex.RUnlock()

	if ok {
		return result
	}
	result = getSubModuleFromTypeInner(target)

	beegoValidateModelInfo.mutex.Lock()
	beegoValidateModelInfo.data[target] = result
	beegoValidateModelInfo.mutex.Unlock()
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
	var data *beegoValidateModelInterface
	interfaceType := reflect.TypeOf(data).Elem()
	targetType := reflect.PtrTo(target)
	return targetType.Implements(interfaceType)
}

func prepareBeegoValidateModelInner(model beegoValidateModelInterface, controller beegoValidateControllerInterface) {
	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()
	modelSubModels := getSubModuleFromType(modelType)
	model.SetAppController(controller)
	for _, singleModel := range modelSubModels {
		target := modelValue.FieldByIndex(singleModel).Addr().Interface().(beegoValidateModelInterface)
		target.SetAppController(controller)
	}
}
func PrepareBeegoValidateModel(controller beegoValidateControllerInterface) {
	prepareBeegoValidateModelInner(controller, controller)
}
