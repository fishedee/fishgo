package web

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"reflect"
	"sync"
)

type beegoValidateModelInfo struct{
	subModelFields []int
	setAppControllerFunc int
	setAppModelFunc int
	prepareFunc int
	finishFunc int
}

var beegoValidateModelInfoMap struct{
	mutex sync.RWMutex
	data map[reflect.Type]beegoValidateModelInfo
}

var beegoValidateModelType reflect.Type

func init(){
	beegoValidateModelInfoMap.data = map[reflect.Type]beegoValidateModelInfo{}
	beegoValidateModelType = reflect.TypeOf(BeegoValidateModel{})
}

type BeegoValidateModel struct {
	*BeegoValidateBasic
	AppController beego.ControllerInterface
	AppModel interface{}
	Ctx *context.Context
}

func (this *BeegoValidateModel)SetAppController(controller beego.ControllerInterface){
	this.AppController = controller
	controllerInfo := getControllerInfo(reflect.TypeOf(this.AppController))
	basicInfo := reflect.ValueOf(this.AppController).Method(controllerInfo.getBasicFunc).Call(nil)[0].Interface().(*BeegoValidateBasic)
	this.BeegoValidateBasic = basicInfo
	this.Ctx = this.BeegoValidateBasic.ctx
}

func (this *BeegoValidateModel)SetAppModel(model interface{}){
	this.AppModel = model
}

func isFromModelType(target reflect.Type)(bool){
	if target.Kind() != reflect.Struct{
		return false
	}
	if target == beegoValidateModelType{
		return true
	}
	for i := 0 ; i != target.NumField() ; i++{
		singleFiled := target.Field(i)
		if singleFiled.Anonymous == false{
			continue
		}
		if isFromModelType(singleFiled.Type) == true{
			return true
		}
	}
	return false
}
func getModelInfoInner(target reflect.Type)(beegoValidateModelInfo){
	//判断是否beegoModel
	isBeegoModel := false
	targetElem := target.Elem()
	if isFromModelType(targetElem){
		isBeegoModel = true
	}

	//填充数据
	result := beegoValidateModelInfo{
		[]int{},
		-1,
		-1,
		-1,
		-1,
	}
	for i := 0 ; i != targetElem.NumField() ; i++{
		singleFiled := targetElem.Field(i)
		if singleFiled.Anonymous{
			continue
		}
		if singleFiled.PkgPath != ""{
			continue
		}
		if isFromModelType(singleFiled.Type){
			result.subModelFields = append(result.subModelFields,i)
		}
	}
	for i := 0 ; i != target.NumMethod() ; i++{
		singleMethod := target.Method(i)
		if singleMethod.Name == "SetAppController" && isBeegoModel{
			result.setAppControllerFunc = i
		}
		if singleMethod.Name == "SetAppModel" && isBeegoModel{
			result.setAppModelFunc = i
		}
		if singleMethod.Name == "Prepare" && isBeegoModel{
			result.prepareFunc = i
		}
		if singleMethod.Name == "Finish" && isBeegoModel{
			result.finishFunc = i
		}
	}
	return result
}

func getModelInfo(target reflect.Type)(beegoValidateModelInfo){
	beegoValidateModelInfoMap.mutex.RLock()
	singleMethodInfo,ok := beegoValidateModelInfoMap.data[target]
	beegoValidateModelInfoMap.mutex.RUnlock()

	if ok{
		return singleMethodInfo
	}

	result := getModelInfoInner(target)

	beegoValidateModelInfoMap.mutex.Lock()
	beegoValidateModelInfoMap.data[target] = result
	beegoValidateModelInfoMap.mutex.Unlock()
	return result
}

func prepareBeegoValidateModelInner(target reflect.Value,controller reflect.Value){
	modelInfo := getModelInfo(target.Type())
	targetElem := target.Elem()
	//初始化subModel
	for _,singleSubModel := range modelInfo.subModelFields{
		prepareBeegoValidateModelInner(targetElem.Field(singleSubModel).Addr(),controller)
	}
	//设置AppController
	if modelInfo.setAppControllerFunc != -1{
		target.Method(modelInfo.setAppControllerFunc).Call([]reflect.Value{controller})
	}
	//设置AppModel
	if modelInfo.setAppModelFunc != -1{
		target.Method(modelInfo.setAppModelFunc).Call([]reflect.Value{target})
	}
	//执行prepare
	if modelInfo.prepareFunc != -1 {
		target.Method(modelInfo.prepareFunc).Call(nil)
	}
}

func finishBeegoValidateModelInner(target reflect.Value){
	modelInfo := getModelInfo(target.Type())
	targetElem := target.Elem()
	//执行finish
	if modelInfo.finishFunc != -1 {
		target.Method(modelInfo.finishFunc).Call(nil)
	}
	//初始化subModel
	for _,singleSubModel := range modelInfo.subModelFields{
		finishBeegoValidateModelInner(targetElem.Field(singleSubModel).Addr())
	}
}

func PrepareBeegoValidateModel(controller interface{}){
	controllerValue := reflect.ValueOf(controller)
	prepareBeegoValidateModelInner(controllerValue,controllerValue)
}

func FinishBeegoValidateModel(controller interface{}){
	controllerValue := reflect.ValueOf(controller)
	finishBeegoValidateModelInner(controllerValue)
}
