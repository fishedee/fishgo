package web

import (
	"time"
	"reflect"
	. "github.com/fishedee/language"
)

type timerController struct{
	BeegoValidateController
}

func (this *timerController)startSingleTimerTask(handler reflect.Value,handlerArgv []reflect.Value){
	defer CatchCrash(func(exception Exception){
		this.Log.Critical("TimerTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]",exception.GetCode(),exception.GetMessage(),exception.GetStackTrace())
		if this.Monitor != nil{
			this.Monitor.AscCriticalCount()
		}
	})
	defer Catch(func(exception Exception){
		this.Log.Error("TimerTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]",exception.GetCode(),exception.GetMessage(),exception.GetStackTrace())
		if this.Monitor != nil{
			this.Monitor.AscErrorCount()
		}
	})
	handler.Call(handlerArgv)
}

func newTimerController()(*timerController){
	controller := &timerController{}
	controller.AppController = controller
	controller.Prepare()
	return controller
}

func newTimeModel(targetType reflect.Type,controller *timerController)(reflect.Value){
	model := reflect.New(targetType.Elem())
	controllerValue := reflect.ValueOf(controller)
	prepareBeegoValidateModelInner(model,controllerValue)
	return model
}

func startSingleTask(handler interface{}){
	controller := newTimerController()
	handlerType := reflect.TypeOf(handler)
	handlerValue := reflect.ValueOf(handler)
	if handlerType.NumIn() == 0{
		controller.startSingleTimerTask(handlerValue,nil)
	}else{
		modelArgv := newTimeModel(handlerType.In(0),controller)
		controller.startSingleTimerTask(handlerValue,[]reflect.Value{modelArgv})
	}
}

func InitTimerTask(duraction time.Duration,handler interface{}){
	//带有延后属性的定时器
	go func(){
		timeChan := time.After(duraction)
		for {
			<- timeChan
			startSingleTask(handler)
			timeChan = time.After(duraction)
		}
	}();
}

func InitTickTask(duraction time.Duration,handler interface{}){
	//没有延后属性的定时器
	go func(){
		tickChan := time.Tick(duraction)
		for {
			<- tickChan
			startSingleTask(handler)
		}
	}();
}