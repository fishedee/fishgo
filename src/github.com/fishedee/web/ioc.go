package web

import (
	"errors"
	"reflect"
)

var (
	mapInjectType = map[reflect.Type]reflect.Type{}
	basicType     = reflect.TypeOf(&Basic{})
)

func getIocPkgPath(target interface{}) string {
	targetValue := reflect.ValueOf(target)
	for targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}
	return targetValue.Type().PkgPath()
}

func getIocRealTarget(target interface{}) reflect.Value {
	targetValue := reflect.ValueOf(target)
	for targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}
	if targetValue.Kind() == reflect.Interface {
		return targetValue
	} else if targetValue.Kind() == reflect.Struct {
		return targetValue.Addr()
	} else {
		return targetValue
	}
}

func addIocTarget(target interface{}) error {
	targetValue := getIocRealTarget(target)
	if targetValue.Kind() == reflect.Interface {
		instanseType := targetValue.Elem().Type().Elem()
		mapInjectType[targetValue.Type()] = instanseType
		return nil
	} else if targetValue.Kind() == reflect.Ptr &&
		targetValue.Elem().Kind() == reflect.Struct {
		return nil
	} else {
		return errors.New("ioc need a *interface or *struct")
	}

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
