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
	//调整为testing.T的输出行数的方式
	t.Helper()
	if equalDesc, isEqual := DeepEqual(left, right); isEqual == false {
		output := ""
		if len(testcase) != 0 {
			output = fmt.Sprintf("testCase: %v , ", testcase)
		}
		t.Error(output + "assert equal fail: " + equalDesc)
	}
}

func AssertException(t *testing.T, code int, message string, handler func(), testcase ...interface{}) {
	//调整为testing.T的输出行数的方式
	t.Helper()
	failDesc := ""
	func() {
		defer CatchCrash(func(e Exception) {
			if e.GetCode() != code {
				failDesc = fmt.Sprintf("assert exception code fail: %v != %v", code, e.GetCode())
			}
			if e.GetMessage() != message {
				failDesc = fmt.Sprintf("assert exception message fail: %v != %v", message, e.GetMessage())
			}
		})
		handler()
		failDesc = "assert exception fail: no exception!"
	}()
	if failDesc != "" {
		output := ""
		if len(testcase) != 0 {
			output = fmt.Sprintf("testCase: %v , ", testcase)
		}
		t.Error(output + failDesc)
	}

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
