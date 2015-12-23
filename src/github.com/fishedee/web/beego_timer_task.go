package web

import (
	"time"
	"reflect"
	. "github.com/fishedee/language"
)

type timerController struct{
	BeegoValidateController
}

func getTimerController()(*BeegoValidateController){
	return &BeegoValidateController{}
}

func getTimeModel(targetType reflect.Type)(reflect.Value){
	model := reflect.New(targetType.Elem())
	controller := getTimerController()
	controllerValue := reflect.ValueOf(controller)
	prepareBeegoValidateModelInner(model,controllerValue)
	return model
}

func startSingleTimerTask(handler reflect.Value,handlerArgv []reflect.Value){
	defer CatchCrash(func(exception Exception){
		Log.Critical("TimerTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]",exception.GetCode(),exception.GetMessage(),exception.GetStackTrace())
	})
	defer Catch(func(exception Exception){
		Log.Error("TimerTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]",exception.GetCode(),exception.GetMessage(),exception.GetStackTrace())
	})
	handler.Call(handlerArgv)
}

func startSingleTask(handler interface{}){
	handlerType := reflect.TypeOf(handler)
	handlerValue := reflect.ValueOf(handler)
	if handlerType.NumIn() == 0{
		startSingleTimerTask(handlerValue,nil)
	}else{
		handlerArgv := getTimeModel(handlerType.In(0))
		startSingleTimerTask(handlerValue,[]reflect.Value{handlerArgv})
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