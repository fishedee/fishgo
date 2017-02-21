package web

import (
	"errors"
	"reflect"
)

func RunMain(handler interface{}) error {
	handlerValue := reflect.ValueOf(handler)
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return errors.New("invalid handler not function")
	}
	if handlerType.NumIn() != 1 {
		return errors.New("invalid handler should has a parameter")
	}
	handlerFirstArgv := handlerType.In(0)
	if handlerFirstArgv.Kind() != reflect.Ptr {
		return errors.New("invalid handler first parameter is not a ptr")
	}
	handlerFirstArgv = handlerFirstArgv.Elem()
	target := reflect.New(handlerFirstArgv)
	basic := initEmptyBasic(nil)
	injectIoc(target, basic)
	handlerValue.Call([]reflect.Value{target})
	return nil
}
