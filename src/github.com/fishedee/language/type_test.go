package language

import (
	. "github.com/fishedee/assert"
	"reflect"
	"testing"
	"time"
)

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
