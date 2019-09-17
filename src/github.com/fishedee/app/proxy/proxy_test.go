package proxy

import (
	. "github.com/fishedee/assert"
	"reflect"
	"testing"
)

type User struct {
	UserId int
	Age    int
	Name   string
}

type IUserAo interface {
	Get(userId int) User
	Add(user User) int
}

type UserAoMock struct {
	GetHandler func(userId int) User
	AddHandler func(user User) int
}

func (this *UserAoMock) Get(userId int) User {
	return this.GetHandler(userId)
}

func (this *UserAoMock) Add(user User) int {
	return this.AddHandler(user)
}

type UserAo struct {
}

func (this *UserAo) Get(userId int) User {
	return User{
		UserId: 123,
		Age:    789,
		Name:   "fish",
	}
}

func (this *UserAo) Add(user User) int {
	return 456
}

func init() {
	RegisterProxyMock([]IUserAo{&UserAoMock{}})
}

func TestProxy(t *testing.T) {
	testLogOut := ""
	proxy := NewProxy()
	proxy.Hook(func(ctx ProxyContext, origin reflect.Value) reflect.Value {
		return reflect.MakeFunc(origin.Type(), func(in []reflect.Value) []reflect.Value {
			testLogOut += "1"
			result := origin.Call(in)
			testLogOut += "4"
			return result
		})
	})
	proxy.Hook(func(ctx ProxyContext, origin reflect.Value) reflect.Value {
		return reflect.MakeFunc(origin.Type(), func(in []reflect.Value) []reflect.Value {
			testLogOut += "2"
			result := origin.Call(in)
			testLogOut += "3"
			return result
		})
	})

	//测试1，直接用Proxy
	userAo := proxy.Proxy([]IUserAo{&UserAo{}}).(IUserAo)

	result := userAo.Get(10002)
	AssertEqual(t, result, User{
		UserId: 123,
		Age:    789,
		Name:   "fish",
	})
	AssertEqual(t, testLogOut, "1234")

	//测试2，用WrapCreatorWithProxy
	testLogOut = ""
	userAoCreator := WrapCreatorWithProxy(func() IUserAo {
		return &UserAo{}
	}).(func(p Proxy) IUserAo)
	userAo2 := userAoCreator(proxy)
	result2 := userAo2.Get(10003)
	AssertEqual(t, result2, User{
		UserId: 123,
		Age:    789,
		Name:   "fish",
	})
	AssertEqual(t, testLogOut, "1234")
}
