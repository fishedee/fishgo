package web

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
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

func (this *BeegoValidateTest) AssertEqual(left interface{}, right interface{}, testCase ...interface{}) {
	isEqual := reflect.DeepEqual(left, right)
	if isEqual {
		return
	}
	backtrace := this.getBackTrace()
	if len(testCase) == 0 {
		this.t.Errorf("assertEqual Fail! %v != %v\n%s", left, right, backtrace)
	} else {
		this.t.Errorf("assertEqual Fail! %v != %v\ntestCase: %#v\n%s", left, right, testCase, backtrace)
	}

}

func (this *BeegoValidateTest) SetTesting(t *testing.T) {
	this.t = t
}

func (this *BeegoValidateTest) RandomInt() int {
	return rand.Int()
}

func (this *BeegoValidateTest) RandomString(length int) string {
	result := []rune{}
	for i := 0; i != length; i++ {
		var single rune
		randInt := rand.Int() % (10 + 26 + 26)
		if randInt < 10 {
			single = rune('0' + randInt)
		} else if randInt < 10+26 {
			single = rune('A' + (randInt - 10))
		} else {
			single = rune('a' + (randInt - 10 - 26))
		}
		result = append(result, single)
	}
	return string(result)
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

func init() {
	rand.Seed(time.Now().Unix())
}
