package assert

import (
	"fmt"
	. "github.com/fishedee/language"
	"path"
	"runtime"
	"strconv"
	"testing"
)

//左右参数比较工具
func AssertEqual(t *testing.T, left interface{}, right interface{}, testcase ...interface{}) {
	//打印出错的行数
	if equalDesc, isEqual := DeepEqual(left, right); isEqual == false {
		_, filename, line, _ := runtime.Caller(1)
		t.Errorf("%+v ,testCase:%v , assert fail: %v", path.Base(filename)+":"+strconv.Itoa(line), testcase, equalDesc)
	}
}

func AssertException(t *testing.T, code int, message string, handler func(), testcase ...interface{}) {
	defer CatchCrash(func(e Exception) {
		if e.GetCode() != code {
			_, filename, line, _ := runtime.Caller(7)
			t.Errorf("%+v ,testCase:%v , assert exception code fail: %v != %v", path.Base(filename)+":"+strconv.Itoa(line), testcase, e.GetCode(), code)
		}
		if e.GetMessage() != message {
			_, filename, line, _ := runtime.Caller(7)
			t.Errorf("%+v ,testCase:%v , assert exception message fail: %v != %v", path.Base(filename)+":"+strconv.Itoa(line), testcase, e.GetMessage(), message)
		}
	})
	handler()
	_, filename, line, _ := runtime.Caller(1)
	t.Errorf("%+v ,testCase:%v , assert exception fail: no exception!", path.Base(filename)+":"+strconv.Itoa(line), testcase)
}

//抛出异常比对工具
func AssertError(t *testing.T, errorText string, function func(), testcase ...interface{}) {
	defer func() {
		r := fmt.Sprintf("%+v", recover())
		if equalDesc, isEqual := DeepEqual(r, errorText); isEqual == false {
			_, filename, line, _ := runtime.Caller(7)
			t.Errorf("%+v ,testCase:%v , assert fail: %v", path.Base(filename)+":"+strconv.Itoa(line), testcase, equalDesc)
		}
	}()
	function()

}
