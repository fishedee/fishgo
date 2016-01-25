package web

import (
	"github.com/astaxie/beego/context"
	"reflect"
	"sync"
)

var beegoValidateModelType reflect.Type
var beegoValidateModelInfo struct {
	mutex sync.RWMutex
	data  map[reflect.Type][]int
}

func init() {
	beegoValidateModelType = reflect.TypeOf(BeegoValidateModel{})
	beegoValidateModelInfo.data = map[reflect.Type][]int{}
}

type beegoValidateModelInterface interface {
	SetAppController(controller beegoValidateControllerInterface)
	SetAppModel(model beegoValidateModelInterface)
	GetSubModel() []beegoValidateModelInterface
}

type BeegoValidateModel struct {
	*BeegoValidateBasic
	AppController interface{}
	AppModel      interface{}
	Ctx           *context.Context
}

func (this *BeegoValidateModel) SetAppController(controller beegoValidateControllerInterface) {
	this.AppController = controller
	this.BeegoValidateBasic = controller.GetBasic()
	this.Ctx = this.BeegoValidateBasic.ctx
}

func (this *BeegoValidateModel) SetAppModel(model beegoValidateModelInterface) {
	this.AppModel = model
}

func (this *BeegoValidateModel) GetSubModel() []beegoValidateModelInterface {
	result := []beegoValidateModelInterface{}
	modelType := reflect.TypeOf(this.AppModel).Elem()
	modelValue := reflect.ValueOf(this.AppModel).Elem()
	modelTypeFields := getSubModuleFromType(modelType)
	for _, i := range modelTypeFields {
		result = append(
			result,
			modelValue.Field(i).Addr().Interface().(beegoValidateModelInterface),
		)
	}
	return result
}

func getSubModuleFromType(target reflect.Type) []int {
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

func getSubModuleFromTypeInner(modelType reflect.Type) []int {
	result := []int{}
	numField := modelType.NumField()
	for i := 0; i != numField; i++ {
		singleFiled := modelType.Field(i)
		if singleFiled.Anonymous {
			continue
		}
		if singleFiled.PkgPath != "" {
			continue
		}
		if isFromModelType(singleFiled.Type) {
			result = append(
				result,
				i,
			)
		}
	}
	return result
}

func isFromModelType(target reflect.Type) bool {
	if target.Kind() != reflect.Struct {
		return false
	}
	if target == beegoValidateModelType {
		return true
	}
	numField := target.NumField()
	for i := 0; i != numField; i++ {
		singleFiled := target.Field(i)
		if singleFiled.Anonymous == false {
			continue
		}
		if isFromModelType(singleFiled.Type) == true {
			return true
		}
	}
	return false
}

func prepareBeegoValidateModelInner(target beegoValidateModelInterface, controller beegoValidateControllerInterface) {
	target.SetAppController(controller)
	target.SetAppModel(target)
	for _, singleTarget := range target.GetSubModel() {
		prepareBeegoValidateModelInner(singleTarget, controller)
	}
}

func PrepareBeegoValidateModel(controller beegoValidateControllerInterface) {
	prepareBeegoValidateModelInner(controller, controller)
}
