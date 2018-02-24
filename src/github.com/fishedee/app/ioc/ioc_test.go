package ioc

import (
	. "github.com/fishedee/assert"
	"testing"
)

type struct1 struct {
	data *struct2
}

type struct2 struct {
	data string
}

type struct3Interface interface {
	Go()
}
type struct3 struct {
	data string
}

func (this *struct3) Go() {

}

func newStruct1(data *struct2) *struct1 {
	return &struct1{data: data}
}

func newStruct2(data string) *struct2 {
	return &struct2{data: data}
}

func newStruct3() struct3Interface {
	return &struct3{data: "Hello God"}
}

func TestIocBasic(t *testing.T) {
	ioc := NewIoc()
	ioc.MustRegister(newStruct1)
	ioc.MustRegister(newStruct2)
	ioc.MustRegister(newStruct3)
	ioc.MustRegister(func() string {
		return "Hello World"
	})

	ioc.MustInvoke(func(a *struct1, b *struct2, c string, d struct3Interface) {
		AssertEqual(t, a.data, b)
		AssertEqual(t, b.data, c)
		AssertEqual(t, c, "Hello World")
		AssertEqual(t, d.(*struct3).data, "Hello God")
	})
}

func newStruct1_Loop(data *struct2) *struct1 {
	return &struct1{data: data}
}

func newStruct2_Loop(data *struct1) *struct2 {
	return &struct2{}
}

func TestIocError(t *testing.T) {
	ioc := NewIoc()
	ioc.MustRegister(newStruct1_Loop)
	ioc.MustRegister(newStruct2_Loop)

	var err error

	err = ioc.Invoke(func(a *struct1, b *struct2) {
	})
	AssertEqual(t, err != nil, true)

	err = ioc.Invoke(func(c string) {
	})
	AssertEqual(t, err != nil, true)

	err = ioc.Invoke("123")
	AssertEqual(t, err != nil, true)

	err = ioc.Register("123")
	AssertEqual(t, err != nil, true)
}
