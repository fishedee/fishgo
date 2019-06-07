package language_test

import (
	"fmt"
	. "github.com/fishedee/language"
	"reflect"
	"testing"
	"time"
)

func AssertEqual(t *testing.T, left interface{}, right interface{}, testCase ...interface{}) {
	t.Helper()
	if reflect.DeepEqual(left, right) != true {
		t.Errorf("assert equal fail testcase:%v, %v != %v", testCase, left, right)
	}
}

func AssertError(t *testing.T, errorText string, function func(), testCase ...interface{}) {
	defer func() {
		r := fmt.Sprintf("%+v", recover())
		if r != errorText {
			t.Errorf("testCase:%v , assert fail: %v != %v", testCase, errorText, r)
		}
	}()
	function()
}

func TestIsEmptyValue(t *testing.T) {
	testCase := []struct {
		data    interface{}
		isEmpty bool
	}{
		{false, true},
		{true, false},
		{0, true},
		{1, false},
		{0.0, true},
		{0.1, false},
		{"", true},
		{"a", false},
		{time.Time{}, true},
		{time.Now(), false},
	}

	for _, singleTestCase := range testCase {
		dataValue := reflect.ValueOf(singleTestCase.data)
		result := IsEmptyValue(dataValue)
		AssertEqual(t, result, singleTestCase.isEmpty)
	}
}
