package assert

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"testing"
)

//左右参数比较工具
func AssertEqual(t *testing.T, left interface{}, right interface{}, testcase ...interface{}) {
	//打印出错的行数
	if reflect.DeepEqual(left, right) == false {
		_, filename, line, _ := runtime.Caller(1)
		t.Errorf("%+v ,testCase:%v , assert fail: %+v != %+v", path.Base(filename)+":"+strconv.Itoa(line), testcase, left, right)
	}
}

//抛出异常比对工具
func AssertError(t *testing.T, errorText string, function func(), testcase ...interface{}) {

	defer func() {
		r := fmt.Sprintf("%+v", recover())
		if reflect.DeepEqual(r, errorText) == false {
			_, filename, line, _ := runtime.Caller(7)
			t.Errorf("%+v ,testCase:%v , assert fail: \"%+v\" != \"%+v\"", path.Base(filename)+":"+strconv.Itoa(line), testcase, r, errorText)
		}
	}()

	function()

}
