package web

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type BeegoValidateTestInterface interface {
	beegoValidateControllerInterface
	SetTesting(*testing.T)
}

type BeegoValidateTest struct {
	BeegoValidateController
	t *testing.T
}

func (this *BeegoValidateTest) getBackTrace() string {
	stack := ""
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
	}
	return stack
}

func (this *BeegoValidateTest) AssertEqual(left interface{}, right interface{}) {
	isEqual := reflect.DeepEqual(left, right)
	if isEqual {
		return
	}
	backtrace := this.getBackTrace()
	this.t.Errorf("assertEqual Fail! %v != %v\n%s", left, right, backtrace)
}

func (this *BeegoValidateTest) SetTesting(t *testing.T) {
	this.t = t
}

func InitBeegoVaildateTest(t *testing.T, test BeegoValidateTestInterface) {
	//初始化test
	test.SetTesting(t)
	test.SetAppControllerInner(test)
	test.Prepare()

	//遍历test，执行测试
	testType := reflect.TypeOf(test)
	testValue := reflect.ValueOf(test)
	testMethodNum := testType.NumMethod()
	for i := 0; i != testMethodNum; i++ {
		singleValueMethodType := testType.Method(i)
		if strings.HasPrefix(singleValueMethodType.Name, "Test") == false {
			continue
		}
		//执行测试
		singleValueMethodType.Func.Call([]reflect.Value{testValue})
	}
}
