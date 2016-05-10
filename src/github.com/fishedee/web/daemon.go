package web

import (
	. "github.com/fishedee/language"
	"reflect"
)

func startSingleTaskInner(basic *Basic, handler reflect.Value, handlerArgv []reflect.Value) {
	defer CatchCrash(func(exception Exception) {
		basic.Log.Critical("DaemonTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	defer Catch(func(exception Exception) {
		basic.Log.Error("DaemonTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	handler.Call(handlerArgv)
}

func startSingleTask(handler interface{}) {
	handlerType := reflect.TypeOf(handler)
	handlerValue := reflect.ValueOf(handler)
	if handlerType.NumIn() == 0 {
		panic("Init Daemon need a type")
	}

	basic := initLocalBasic(nil)
	modelArgv, err := newIocInstanse(handlerType.In(0), basic)
	if err != nil {
		panic(err)
	}
	startSingleTaskInner(basic, handlerValue, []reflect.Value{modelArgv})
}

func InitDaemon(handler interface{}) {
	go startSingleTask(handler)
}
