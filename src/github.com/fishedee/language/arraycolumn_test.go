package language

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func assertEqual(t *testing.T, left interface{}, right interface{}) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("%v != %v", left, right)
	}
}

func TestArrayColumnSort(t *testing.T) {

	aa := []struct {
		Name     string
		Age      int
		OK       bool
		Money    float64
		Register time.Time
	}{
		{"哈哈", 15, true, 44.5, time.Now()},
		{"呵呵", 22, false, 30.1, time.Now()},
		{"啊啊", 11, true, 54.1, time.Now()},
	}
	t.Errorf("调试：")
	fmt.Printf("%+v", ArrayColumnSort(aa, "Money"))

}
