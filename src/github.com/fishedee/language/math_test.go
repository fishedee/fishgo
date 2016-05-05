package language

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, x, y int) {
	if reflect.DeepEqual(x, y) == false {
		t.Error(fmt.Sprintf("%#v != %#v", x, y))
	}
}

func TestAbsInt(t *testing.T) {
	testCase := []struct {
		origin int
		target int
	}{
		{1, 1},
		{0, 0},
		{-1, 1},
	}

	for _, singleTestCase := range testCase {
		assert(t, AbsInt(singleTestCase.origin), singleTestCase.target)
	}
}
