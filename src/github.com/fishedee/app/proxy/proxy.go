package proxy

import (
	"fmt"
	"reflect"
)

var (
	proxyMock = map[reflect.Type]ProxyHandler{}
)

type ProxyContext struct {
	PkgPath       string
	InterfaceName string
	MethodName    string
}

type ProxyHook = func(ctx ProxyContext, origin reflect.Value) reflect.Value

type ProxyHandler = func(originTarget reflect.Value, hookers []ProxyHook) reflect.Value

type Proxy interface {
	Hook(hooker ProxyHook)
	Proxy(originTarget interface{}) interface{}
	ProxyValue(originTarget reflect.Value) reflect.Value
}

type proxyImplement struct {
	hookers []ProxyHook
}

func NewProxy() Proxy {
	return &proxyImplement{
		hookers: []ProxyHook{},
	}
}

func (this *proxyImplement) Hook(hooker ProxyHook) {
	this.hookers = append(this.hookers, hooker)
}

func (this *proxyImplement) ProxyValue(originTarget reflect.Value) reflect.Value {
	handler, isExist := proxyMock[originTarget.Type()]
	if isExist == false {
		panic(fmt.Sprintf("can not proxy type %v,you should register proxy mock first", originTarget.Type()))
	}
	return handler(originTarget, this.hookers)
}

func (this *proxyImplement) Proxy(originTarget interface{}) interface{} {
	originValue := reflect.ValueOf(originTarget).Index(0)
	return this.ProxyValue(originValue).Interface()
}

type proxyMockInfo struct {
	ProxyMethodIndex   int
	ProxyMethodName    string
	ProxyPkgPath       string
	ProxyInterfaceName string
	MockFieldIndex     []int
}

func checkEqualType(mockMethodType reflect.Type, proxyMethodType reflect.Type) bool {
	if mockMethodType.Kind() != proxyMethodType.Kind() ||
		mockMethodType.NumIn() != proxyMethodType.NumIn() ||
		mockMethodType.NumOut() != proxyMethodType.NumOut() {
		return false
	}
	for i := 0; i != mockMethodType.NumIn(); i++ {
		if mockMethodType.In(i).String() != proxyMethodType.In(i).String() {
			return false
		}
	}
	for i := 0; i != mockMethodType.NumOut(); i++ {
		if mockMethodType.Out(i).String() != proxyMethodType.Out(i).String() {
			return false
		}
	}
	return true
}

func RegisterProxyMock(proxyTarget interface{}) {
	proxyValue := reflect.ValueOf(proxyTarget).Index(0)
	proxyType := proxyValue.Type()
	if proxyType.Kind() != reflect.Interface {
		panic(fmt.Sprintf("originTarget must be a interface,[%v]", proxyType))
	}

	mockValue := proxyValue.Elem()
	mockType := mockValue.Type()
	mockElemType := mockType.Elem()

	proxyInfos := []proxyMockInfo{}
	proxyMethodNum := proxyType.NumMethod()
	for i := 0; i != proxyMethodNum; i++ {
		proxySingleMethod := proxyType.Method(i)
		mockMethodFieldName := proxySingleMethod.Name + "Handler"
		mockMethodField, isExist := mockElemType.FieldByName(mockMethodFieldName)
		if isExist == false {
			panic(fmt.Sprintf("%v dos not have field %v", mockElemType, mockMethodFieldName))
		}
		if checkEqualType(mockMethodField.Type, proxySingleMethod.Type) == false {
			panic(fmt.Sprintf("%v.%v dos not have right tye", mockElemType, mockMethodFieldName))
		}
		proxyInfos = append(proxyInfos, proxyMockInfo{
			ProxyMethodIndex:   i,
			ProxyPkgPath:       proxyType.PkgPath(),
			ProxyInterfaceName: proxyType.Name(),
			ProxyMethodName:    proxySingleMethod.Name,
			MockFieldIndex:     mockMethodField.Index,
		})
	}

	proxyMock[proxyType] = func(originTarget reflect.Value, hookers []ProxyHook) reflect.Value {
		newMockValue := reflect.New(mockElemType)
		newMockElemValue := newMockValue.Elem()

		for _, proxyInfo := range proxyInfos {
			ctx := ProxyContext{
				PkgPath:       proxyInfo.ProxyPkgPath,
				InterfaceName: proxyInfo.ProxyInterfaceName,
				MethodName:    proxyInfo.ProxyMethodName,
			}
			method := originTarget.Method(proxyInfo.ProxyMethodIndex)
			for i := len(hookers) - 1; i >= 0; i-- {
				method = hookers[i](ctx, method)
			}
			newMockElemValue.FieldByIndex(proxyInfo.MockFieldIndex).Set(method)
		}
		return newMockValue.Convert(proxyType)
	}
}

func WrapCreatorWithProxy(originCreator interface{}) interface{} {
	proxyType := reflect.ValueOf((*Proxy)(nil)).Type().Elem()
	creator := reflect.ValueOf(originCreator)
	creatorType := creator.Type()

	numInTypes := []reflect.Type{proxyType}
	numOutTypes := []reflect.Type{}
	for i := 0; i != creatorType.NumIn(); i++ {
		numInTypes = append(numInTypes, creatorType.In(i))
	}
	for i := 0; i != creatorType.NumOut(); i++ {
		numOutTypes = append(numOutTypes, creatorType.Out(i))
	}
	newCreatorType := reflect.FuncOf(numInTypes, numOutTypes, false)
	newCreator := reflect.MakeFunc(newCreatorType, func(args []reflect.Value) []reflect.Value {
		result := creator.Call(args[1:])
		proxy := args[0].Interface().(Proxy)
		proxyResult := proxy.ProxyValue(result[0])
		return []reflect.Value{proxyResult}
	})
	return newCreator.Interface()
}
