package web

import (
	. "github.com/fishedee/language"
	"reflect"
)

func startSingleTaskInner(handler reflect.Value, handlerArgv []reflect.Value, basic *Basic) {
	defer CatchCrash(func(exception Exception) {
		basic.Log.Critical("DaemonTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	defer Catch(func(exception Exception) {
		basic.Log.Error("DaemonTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	handler.Call(handlerArgv)
}

func startSingleTask(handler interface{}) {
	basic := initEmptyBasic(nil)
	handlerType := reflect.TypeOf(handler)
	handlerValue := reflect.ValueOf(handler)
	if handlerType.NumIn() == 0 {
		startSingleTaskInner(handlerValue, nil, basic)
	} else {
		target := reflect.New(handlerType.In(0).Elem())
		injectIoc(target, basic)
		startSingleTaskInner(handlerValue, []reflect.Value{target}, basic)
	}
}

func InitDaemon(handler interface{}) {
	startSingleTask(handler)
}
