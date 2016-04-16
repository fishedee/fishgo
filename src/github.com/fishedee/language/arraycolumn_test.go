package language

import (
	// "fmt"
	"reflect"
	"testing"
	"time"
)

func assertEqual(t *testing.T, left interface{}, right interface{}, index int) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("case :%v ,%+v != %+v", index, left, right)
	}
}

func TestArrayColumnSort(t *testing.T) {
	//测试类型 支持bool,int,float,string和time.Time
	type contentType struct {
		Name     string
		Age      int
		Ok       bool
		Money    float64
		Register time.Time
	}

	nowTime := time.Now()
	oldTime := nowTime.AddDate(-1, 0, 1)

	testCase := []struct {
		sortName string
		origin   interface{}
		target   interface{}
	}{
		//空集
		{
			"Name desc",
			[]contentType{},
			[]contentType{},
		},

		{
			"Name desc",
			[]contentType{
				contentType{"5", 0, true, -1.1, oldTime},
				contentType{"z", 1, true, 0, nowTime},
				contentType{"", 0, false, 0, time.Time{}},
				contentType{"a", -1, false, 1.1, time.Time{}},
			},
			[]contentType{
				contentType{"z", 1, true, 0, nowTime},
				contentType{"a", -1, false, 1.1, time.Time{}},
				contentType{"5", 0, true, -1.1, oldTime},
				contentType{"", 0, false, 0, time.Time{}},
			},
		},
		{
			"Age desc,Ok desc",
			[]contentType{
				contentType{"z", 3, true, 0, nowTime},
				contentType{"a", -1, false, 1.1, time.Time{}},
				contentType{"5", 10, true, -1.1, oldTime},
				contentType{"", 5, false, 0, time.Time{}},
			},
			[]contentType{
				contentType{"5", 10, true, -1.1, oldTime},
				contentType{"", 5, false, 0, time.Time{}},
				contentType{"z", 3, true, 0, nowTime},
				contentType{"a", -1, false, 1.1, time.Time{}},
			},
		},
		{
			"Money,Register desc",
			[]contentType{
				contentType{"z", 3, true, 0, nowTime},
				contentType{"a", -1, false, 1.1, time.Time{}},
				contentType{"5", 10, true, -1.1, oldTime},
				contentType{"", 5, false, 0, time.Time{}},
			},
			[]contentType{
				contentType{"5", 10, true, -1.1, oldTime},
				contentType{"z", 3, true, 0, nowTime},
				contentType{"", 5, false, 0, time.Time{}},
				contentType{"a", -1, false, 1.1, time.Time{}},
			},
		},
		{
			"Ok desc,Name",
			[]contentType{
				contentType{"z", 3, true, 0, nowTime},
				contentType{"a", -1, false, 1.1, time.Time{}},
				contentType{"5", 10, true, -1.1, oldTime},
				contentType{"", 5, false, 0, time.Time{}},
			},
			[]contentType{
				contentType{"5", 10, true, -1.1, oldTime},
				contentType{"z", 3, true, 0, nowTime},
				contentType{"", 5, false, 0, time.Time{}},
				contentType{"a", -1, false, 1.1, time.Time{}},
			},
		},
		{
			sortName: " Money desc,Age asc",
			origin: []contentType{
				contentType{"z", 3, true, 0, nowTime},
				contentType{"a", -1, false, 1.1, time.Time{}},
				contentType{"5", 10, true, -1.1, oldTime},
				contentType{"", 5, false, 0, time.Time{}},
			},
			target: []contentType{
				contentType{"a", -1, false, 1.1, time.Time{}},
				contentType{"z", 3, true, 0, nowTime},
				contentType{"", 5, false, 0, time.Time{}},
				contentType{"5", 10, true, -1.1, oldTime},
			},
		},
		{
			sortName: " Money desc,Age asc,Name desc",
			origin: []contentType{
				contentType{"b", 3, true, 0, nowTime},
				contentType{"a", -1, false, 1.1, time.Time{}},
				contentType{"5", 10, true, -1.1, oldTime},
				contentType{"", 5, false, 0, time.Time{}},
				contentType{"h", 3, true, 0, nowTime},
			},
			target: []contentType{
				contentType{"a", -1, false, 1.1, time.Time{}},
				contentType{"h", 3, true, 0, nowTime},
				contentType{"b", 3, true, 0, nowTime},
				contentType{"", 5, false, 0, time.Time{}},
				contentType{"5", 10, true, -1.1, oldTime},
			},
		},
	}
	for singleTestCaseIndex, singleTestCase := range testCase {

		result := ArrayColumnSort(singleTestCase.origin, singleTestCase.sortName)
		assertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}

	// t.Errorf("调试：")
	// fmt.Println("\n")
	// fmt.Printf("%+v", ArrayColumnSort(testCase[0].origin, testCase[0].sortName))

}
