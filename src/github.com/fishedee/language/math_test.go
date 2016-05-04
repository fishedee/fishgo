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
	x1 := []int{1, 0, -1}

	assert(t, AbsInt(x1[0]), 1)
	assert(t, AbsInt(x1[1]), 0)
	assert(t, AbsInt(x1[2]), 1)
}
