package language

import (
	"fmt"
	"reflect"
	"testing"
)

func AssertEqualPlode(t *testing.T, left interface{}, right interface{}, testCase interface{}) {
	isEqual := reflect.DeepEqual(left, right)
	if isEqual == false {
		t.Error(fmt.Sprintf("%#v != %#v,testCase:%v", left, right, testCase))
	}
}

func TestPlode(t *testing.T) {
	testCase := []struct {
		origin    string
		seperator string
		data      interface{}
	}{
		{"", ",", []int{}},
		{"", ",", []string{}},
		{"1", ",", []int{1}},
		{"mmx", ",", []string{"mmx"}},
		{"1,2,3", ",", []int{1, 2, 3}},
		{"mmx,mmd,xxu", ",", []string{"mmx", "mmd", "xxu"}},
		{"1_2_3", "_", []int{1, 2, 3}},
		{"mmx_mmd_xxu", "_", []string{"mmx", "mmd", "xxu"}},
	}

	//test explode
	for _, singleTestCase := range testCase {
		dataType := reflect.TypeOf(singleTestCase.data)
		dataValue := reflect.New(dataType)
		Explode(singleTestCase.origin, singleTestCase.seperator, dataValue.Interface())
		AssertEqualPlode(t, dataValue.Elem().Interface(), singleTestCase.data, singleTestCase)
	}

	//test implode
	for _, singleTestCase := range testCase {
		singleOrigin := Implode(singleTestCase.data, singleTestCase.seperator)
		AssertEqualPlode(t, singleOrigin, singleTestCase.origin, singleTestCase)
	}
}
