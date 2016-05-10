package web

import (
	"errors"
	"reflect"
)

var (
	mapInjectType = map[reflect.Type]reflect.Type{}
	basicType     = reflect.TypeOf(&Basic{})
)

func addIocTarget(target interface{}) error {
	targetValue := reflect.ValueOf(target)
	targetType := targetValue.Type().Elem()
	if targetType.Kind() != reflect.Interface {
		return errors.New("ioc need a interface")
	}
	instanseType := targetValue.Elem().Elem().Type().Elem()
	mapInjectType[targetType] = instanseType
	return nil
}

func injectIocTarget(target reflect.Value, basic reflect.Value) error {
	targetType := target.Type()
	if targetType.Kind() != reflect.Struct {
		return nil
	}
	fieldNum := targetType.NumField()
	for i := 0; i != fieldNum; i++ {
		singleTarget := target.Field(i)
		singleTargetType := targetType.Field(i)
		if singleTargetType.PkgPath != "" {
			continue
		}
		if singleTarget.Type() == basicType {
			singleTarget.Set(basic)
		} else if singleTarget.Kind() == reflect.Interface &&
			singleTarget.IsNil() {
			instanseType, isExist := mapInjectType[singleTargetType.Type]
			if isExist == false {
				continue
			}
			instanseValue := reflect.New(instanseType)
			singleTarget.Set(instanseValue)
			err := injectIocTarget(instanseValue.Elem(), basic)
			if err != nil {
				return err
			}
		} else {
			err := injectIocTarget(singleTarget, basic)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func newIocInstanse(targetType reflect.Type, basic *Basic) (reflect.Value, error) {
	var instanseType reflect.Type
	if targetType.Kind() == reflect.Interface {
		var isExist bool
		instanseType, isExist = mapInjectType[targetType]
		if isExist == false {
			return reflect.Value{}, errors.New("invalid instanseType " + targetType.String())
		}
	} else {
		instanseType = targetType.Elem()
	}
	instanseValue := reflect.New(instanseType)
	err := injectIocTarget(instanseValue.Elem(), reflect.ValueOf(basic))
	if err != nil {
		return reflect.Value{}, err
	}
	return instanseValue, nil
}
