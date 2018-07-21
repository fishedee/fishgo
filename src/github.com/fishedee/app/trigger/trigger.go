package trigger

import (
	"fmt"
	"reflect"
)

type Trigger interface {
	Fire(name string, args ...interface{}) []interface{}
	On(name string, handler interface{})
}

type triggerImplement struct {
	handlers map[string][]reflect.Value
}

func NewTrigger() (Trigger, error) {
	trigger := &triggerImplement{
		handlers: map[string][]reflect.Value{},
	}
	return trigger, nil
}

func (this *triggerImplement) On(name string, handler interface{}) {
	handlerValue := reflect.ValueOf(handler)
	handlers, isExist := this.handlers[name]
	if isExist == false {
		handlers = []reflect.Value{}
	}
	handlers = append(handlers, handlerValue)
	this.handlers[name] = handlers
}

func (this *triggerImplement) Fire(name string, args ...interface{}) []interface{} {
	handlers, isExist := this.handlers[name]
	if isExist == false {
		return nil
	}
	argsValue := []reflect.Value{}
	for _, arg := range args {
		argsValue = append(argsValue, reflect.ValueOf(arg))
	}
	result := []interface{}{}
	for _, handler := range handlers {
		numIn := handler.Type().NumIn()
		if numIn > len(argsValue) {
			panic(fmt.Sprintf("%v can't recevive %v argument", handler.Type().Name(), len(argsValue)))
		}
		resultValue := handler.Call(argsValue[0:numIn])
		if len(resultValue) == 0 {
			result = append(result, nil)
		} else {
			result = append(result, resultValue[0].Interface())
		}
	}
	return result
}
