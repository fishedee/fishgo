package ioc

import (
	"errors"
	"reflect"
)

var (
	defaultIoc = NewIoc()
)

type typeInfo struct {
	depType []reflect.Type
	builder reflect.Value
}

type Ioc struct {
	builder map[reflect.Type]typeInfo
	cache   map[reflect.Type]reflect.Value
}

func NewIoc() *Ioc {
	ioc := &Ioc{}
	ioc.builder = map[reflect.Type]typeInfo{}
	ioc.cache = map[reflect.Type]reflect.Value{}
	return ioc
}

func (this *Ioc) dfs(t reflect.Type, visit map[reflect.Type]bool) (reflect.Value, error) {
	result, isExist := this.cache[t]
	if isExist {
		return result, nil
	}
	_, isVisit := visit[t]
	if isVisit {
		return reflect.Value{}, errors.New("loop dependence")
	}
	visit[t] = true

	info, isExist := this.builder[t]
	if isExist == false {
		return reflect.Value{}, errors.New("unknown type [" + t.String() + "]")
	}
	args := []reflect.Value{}
	for _, singleDepType := range info.depType {
		singleDepValue, err := this.dfs(singleDepType, visit)
		if err != nil {
			return reflect.Value{}, errors.New("-> " + t.Name() + " " + err.Error())
		}
		args = append(args, singleDepValue)
	}
	lastResult := info.builder.Call(args)
	this.cache[t] = lastResult[0]
	return lastResult[0], nil
}

func (this *Ioc) Register(createFun interface{}) error {
	typeValue := reflect.ValueOf(createFun)
	typeType := typeValue.Type()
	if typeType.Kind() != reflect.Func {
		return errors.New("invalid type")
	}
	numIn := []reflect.Type{}
	for i := 0; i != typeType.NumIn(); i++ {
		numIn = append(numIn, typeType.In(i))
	}
	if typeType.NumOut() != 1 {
		return errors.New("invalid num out")
	}
	numOut := typeType.Out(0)
	this.builder[numOut] = typeInfo{
		depType: numIn,
		builder: typeValue,
	}
	return nil
}

func (this *Ioc) MustRegister(createFun interface{}) {
	err := this.Register(createFun)
	if err != nil {
		panic(err)
	}
}

func (this *Ioc) Invoke(createFun interface{}) error {
	typeValue := reflect.ValueOf(createFun)
	typeType := typeValue.Type()
	if typeType.Kind() != reflect.Func {
		return errors.New("invalid type")
	}
	numIn := []reflect.Value{}
	visit := map[reflect.Type]bool{}
	for i := 0; i != typeType.NumIn(); i++ {
		singleNumInType := typeType.In(i)
		singleNumInValue, err := this.dfs(singleNumInType, visit)
		if err != nil {
			return errors.New("build [invoke] " + err.Error())
		}
		numIn = append(numIn, singleNumInValue)
	}
	typeValue.Call(numIn)
	return nil
}

func (this *Ioc) MustInvoke(createFun interface{}) {
	err := this.Invoke(createFun)
	if err != nil {
		panic(err)
	}
}

func RegisterIoc(createFun interface{}) error {
	return defaultIoc.Register(createFun)
}

func MustRegisterIoc(createFun interface{}) {
	defaultIoc.MustRegister(createFun)
}

func InvokeIoc(createFun interface{}) error {
	return defaultIoc.Invoke(createFun)
}

func MustInvokeIoc(createFun interface{}) {
	defaultIoc.MustInvoke(createFun)
}
