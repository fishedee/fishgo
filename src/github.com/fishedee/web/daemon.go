package web

import (
	_ "github.com/a"
	. "github.com/fishedee/language"
	"reflect"
)

type daemonController struct {
	Controller
}

func (this *daemonController) startSingleTask(handler reflect.Value, handlerArgv []reflect.Value) {
	defer CatchCrash(func(exception Exception) {
		this.Log.Critical("DaemonTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	defer Catch(func(exception Exception) {
		this.Log.Error("DaemonTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	handler.Call(handlerArgv)
}

func newDaemonModel(targetType reflect.Type, controller *daemonController) reflect.Value {
	model := reflect.New(targetType.Elem())
	initModelInner(model.Interface().(beegoValidateModelInterface), controller)
	return model
}

func startSingleTask(handler interface{}) {
	controller := &daemonController{}
	controller.InitEmpty(controller, nil)

	handlerType := reflect.TypeOf(handler)
	handlerValue := reflect.ValueOf(handler)
	if handlerType.NumIn() == 0 {
		controller.startSingleTask(handlerValue, nil)
	} else {
		modelArgv := newDaemonModel(handlerType.In(0), controller)
		controller.startSingleTask(handlerValue, []reflect.Value{modelArgv})
	}
}

func InitDaemon(handler interface{}) {
	go startSingleTask(handler)
}
