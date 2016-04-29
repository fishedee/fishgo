package language

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func assertTypeEqual(t *testing.T, left interface{}, right interface{}) {
	isEqual := reflect.DeepEqual(left, right)
	if isEqual == false {
		t.Error(fmt.Sprintf("%#v != %#v", left, right))
	}
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
		assertTypeEqual(t, result, singleTestCase.isEmpty)
	}
}
