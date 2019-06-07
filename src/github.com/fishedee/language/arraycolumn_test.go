package language_test

import (
	. "github.com/fishedee/language"
	"testing"
	"time"
)

func TestKSort(t *testing.T) {
	//测试类型 支持bool,int,float,string和time.Time
	type contentType struct {
		Name      string
		Age       int
		Ok        bool
		Money     float32
		CardMoney float64
		Register  time.Time
	}

	nowTime := time.Now()
	oldTime := nowTime.AddDate(-1, 0, 1)
	zeroTime := time.Time{}

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
				contentType{"5", 0, true, -1.1, -1.1, oldTime},
				contentType{"z", 1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
			},
			[]contentType{
				contentType{"z", 1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"5", 0, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
			},
		},
		{
			"Age desc,Ok desc",
			[]contentType{
				contentType{"z", -1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 5, false, 0, 0, zeroTime},
			},
			[]contentType{
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 5, false, 0, 0, zeroTime},
				contentType{"z", -1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
			},
		},
		{
			"Money,Register desc",
			[]contentType{
				contentType{"z", -1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 5, false, 0, 0, zeroTime},
			},
			[]contentType{
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"z", -1, true, 0, 0, nowTime},
				contentType{"", 5, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
			},
		},
		{
			"CardMoney,Register desc",
			[]contentType{
				contentType{"z", -1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 5, false, 0, 0, zeroTime},
			},
			[]contentType{
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"z", -1, true, 0, 0, nowTime},
				contentType{"", 5, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
			},
		},
		{
			"Ok desc,Name",
			[]contentType{
				contentType{"z", -1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 5, false, 0, 0, zeroTime},
			},
			[]contentType{
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"z", -1, true, 0, 0, nowTime},
				contentType{"", 5, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
			},
		},
		{
			" Money desc,Age asc",
			[]contentType{
				contentType{"z", -1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 5, false, 0, 0, zeroTime},
			},
			[]contentType{
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"z", -1, true, 0, 0, nowTime},
				contentType{"", 5, false, 0, 0, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
			},
		},
		{
			" Money desc,Age asc,Name desc",
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 5, false, 0, 0, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
			},
			[]contentType{
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 5, false, 0, 0, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
			},
		},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {

		result := ArrayColumnSort(singleTestCase.origin, singleTestCase.sortName)
		AssertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}

}

func TestArrayColumnUnique(t *testing.T) {

	type contentType struct {
		Name      string
		Age       int
		Ok        bool
		Money     float32
		CardMoney float64
		Register  time.Time
	}

	nowTime := time.Now()
	oldTime := nowTime.AddDate(-1, 0, 1)
	zeroTime := time.Time{}

	testCase := []struct {
		uniqueName string
		origin     interface{}
		target     interface{}
	}{
		//空集
		{
			"",
			[]contentType{},
			[]contentType{},
		},
		{
			"   Name    ",
			[]contentType{},
			[]contentType{},
		},
		//默认值
		{
			"",
			[]contentType{
				contentType{"", 0, false, 0, 0, zeroTime},
			},
			[]contentType{
				contentType{"", 0, false, 0, 0, zeroTime},
			},
		},
		//单排除
		{
			"Name",
			[]contentType{
				contentType{"s", 3, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			[]contentType{
				contentType{"s", 3, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 10, true, -1.1, -1.1, oldTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
		},
		{
			"Ok",
			[]contentType{
				contentType{"b", 3, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"h", 3, true, 0, 0, nowTime},
			},
			[]contentType{
				contentType{"b", 3, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
			},
		},
		{
			"   Age   ",
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"a", 0, false, 1.1, 1.1, zeroTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
			},
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"a", 0, false, 1.1, 1.1, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
			},
		},
		{
			"   Money",
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", 0, false, 1.1, 1.1, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
			},
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"a", 0, false, 1.1, 1.1, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
			},
		},
		{
			"   CardMoney",
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", 0, false, 1.1, 1.1, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
			},
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"a", 0, false, 1.1, 1.1, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
			},
		},
		{
			"Register   ",
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"a", 0, false, 1.1, 1.1, zeroTime},
			},
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
			},
		},
		//多值传递
		{
			"  Age  ,  Money",
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"a", 0, false, 1.1, 1.1, zeroTime},
			},
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"a", 0, false, 1.1, 1.1, zeroTime},
			},
		},
		{
			"  Name  ,  Money,Register  ",
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", 15, true, 1.1, 1.1, zeroTime},
				contentType{"5", 0, false, -1.1, -1.1, oldTime},
			},
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"a", 15, true, 1.1, 1.1, zeroTime},
			},
		},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {

		result := ArrayColumnUnique(singleTestCase.origin, singleTestCase.uniqueName)
		AssertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}

}

type ArrayColumnInnerStruct struct {
	MM int
}

type ArrayColumnInnerStruct2 struct {
	ArrayColumnInnerStruct
	MM int
	DD float32
}

func TestArrayColumnKey(t *testing.T) {

	type contentType struct {
		Name      string
		Age       int
		Ok        bool
		Money     float32
		CardMoney float64
		Register  time.Time
	}

	nowTime := time.Now()
	oldTime := nowTime.AddDate(-1, 0, 1)
	zeroTime := time.Time{}

	testCase := []struct {
		keyName string
		origin  interface{}
		target  interface{}
	}{
		{
			"Name              ",
			[]contentType{},
			[]string{},
		},
		{
			"           Name              ",
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", 15, true, 1.1, 1.1, zeroTime},
				contentType{"5", 0, false, -1.1, -1.1, oldTime},
			},
			[]string{"b", "", "h", "5", "", "a", "5"},
		},
		{
			"Age              ",
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", 15, true, 1.1, 1.1, zeroTime},
				contentType{"5", 0, false, -1.1, -1.1, oldTime},
			},
			[]int{-1, 0, -1, 10, 0, 15, 0},
		},
		{
			"Ok    ",
			[]contentType{
				contentType{"b", -1, true, 0, 0, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", 15, true, 1.1, 1.1, zeroTime},
				contentType{"5", 0, false, -1.1, -1.1, oldTime},
			},
			[]bool{true, false, true, true, false, true, false},
		},
		{
			"Money    ",
			[]contentType{

				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", 15, true, 1.1, 1.1, zeroTime},
			},
			[]float32{-1.1, 0, 1.1},
		},
		{
			"CardMoney    ",
			[]contentType{

				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", 15, true, 1.1, 1.1, zeroTime},
			},
			[]float64{-1.1, 0, 1.1},
		},
		{
			" Register    ",
			[]contentType{

				contentType{"5", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"h", -1, true, 0, 0, nowTime},
			},
			[]time.Time{oldTime, zeroTime, nowTime},
		},
		{
			"MM",
			[]ArrayColumnInnerStruct2{
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{1}, 2, 1.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{2}, 4, 2.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{3}, 5, 3.1},
			},
			[]int{2, 4, 5},
		},
		{
			"DD",
			[]ArrayColumnInnerStruct2{
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{1}, 2, 1.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{2}, 4, 2.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{3}, 5, 3.1},
			},
			[]float32{1.1, 2.1, 3.1},
		},
		{
			"ArrayColumnInnerStruct.MM",
			[]ArrayColumnInnerStruct2{
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{1}, 2, 1.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{2}, 4, 2.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{3}, 5, 3.1},
			},
			[]int{1, 2, 3},
		},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {

		result := ArrayColumnKey(singleTestCase.origin, singleTestCase.keyName)
		AssertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}

}

func TestArrayColumnMap(t *testing.T) {

	type contentType struct {
		Name      string
		Age       int
		Ok        bool
		Money     float32
		CardMoney float64
		Register  time.Time
	}

	nowTime := time.Now()
	oldTime := nowTime.AddDate(-1, 0, 1)
	zeroTime := time.Time{}

	testCase := []struct {
		mapName string
		origin  interface{}
		target  interface{}
	}{
		{
			"Name              ",
			[]contentType{},
			map[string]contentType{},
		},
		{
			"           Name              ",
			[]contentType{
				contentType{"b", 3, true, 1.1, 1.1, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"5", -1, false, -1.1, -1.1, oldTime},
				contentType{"", -10, true, -1.1, -1.1, zeroTime},
			},
			map[string]contentType{
				"b": contentType{"b", 3, true, 1.1, 1.1, nowTime},
				"":  contentType{"", 0, false, 0, 0, zeroTime},
				"5": contentType{"5", -1, false, -1.1, -1.1, oldTime},
			},
		},
		{
			"           Age",
			[]contentType{
				contentType{"b", 3, true, 1.1, 1.1, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"5", -1, false, -1.1, -1.1, oldTime},
				contentType{"ss", 0, true, 0, 0, zeroTime},
			},
			map[int]contentType{
				3:  contentType{"b", 3, true, 1.1, 1.1, nowTime},
				0:  contentType{"", 0, false, 0, 0, zeroTime},
				-1: contentType{"5", -1, false, -1.1, -1.1, oldTime},
			},
		},
		{
			"Ok    ",
			[]contentType{
				contentType{"b", 3, true, 1.1, 1.1, nowTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"5", -1, false, -1.1, -1.1, oldTime},
			},
			map[bool]contentType{
				true:  contentType{"b", 3, true, 1.1, 1.1, nowTime},
				false: contentType{"", 0, false, 0, 0, zeroTime},
			},
		},
		{
			"Money",
			[]contentType{
				contentType{"b", 3, true, 1.1, 1.1, nowTime},
				contentType{"5", -1, false, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"", -10, true, -1.1, -1.1, zeroTime},
			},
			map[float32]contentType{
				1.1:  contentType{"b", 3, true, 1.1, 1.1, nowTime},
				-1.1: contentType{"5", -1, false, -1.1, -1.1, oldTime},
				0:    contentType{"", 0, false, 0, 0, zeroTime},
			},
		},
		{
			"CardMoney",
			[]contentType{
				contentType{"b", 3, true, 1.1, 1.1, nowTime},
				contentType{"5", -1, false, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"", -10, true, -1.1, -1.1, zeroTime},
			},
			map[float64]contentType{
				1.1:  contentType{"b", 3, true, 1.1, 1.1, nowTime},
				-1.1: contentType{"5", -1, false, -1.1, -1.1, oldTime},
				0:    contentType{"", 0, false, 0, 0, zeroTime},
			},
		},
		{
			" Register    ",
			[]contentType{
				contentType{"b", 3, true, 1.1, 1.1, nowTime},
				contentType{"5", -1, false, -1.1, -1.1, oldTime},
				contentType{"", -10, true, -1.1, -1.1, zeroTime},
				contentType{"", 0, false, 0, 0, zeroTime},
			},
			map[time.Time]contentType{
				nowTime:  contentType{"b", 3, true, 1.1, 1.1, nowTime},
				oldTime:  contentType{"5", -1, false, -1.1, -1.1, oldTime},
				zeroTime: contentType{"", -10, true, -1.1, -1.1, zeroTime},
			},
		},
		{
			"  Age  ,  Ok    ",
			[]contentType{
				contentType{"b", 3, true, 1.1, 1.1, nowTime},
				contentType{"5", -1, false, -1.1, -1.1, oldTime},
				contentType{"5", -1, false, -1.1, -1.1, oldTime},
				contentType{"5", -1, true, -1.1, -1.1, oldTime},
				contentType{"5", -1, false, -1.1, -1.1, oldTime},
				contentType{"", 0, true, -1.1, -1.1, zeroTime},
				contentType{"", 0, false, -1.1, -1.1, zeroTime},
				contentType{"", 0, true, -1.1, -1.1, zeroTime},
			},
			map[int]map[bool]contentType{
				3: map[bool]contentType{
					true: {"b", 3, true, 1.1, 1.1, nowTime},
				},
				-1: map[bool]contentType{
					false: {"5", -1, false, -1.1, -1.1, oldTime},
					true:  {"5", -1, true, -1.1, -1.1, oldTime},
				},
				0: map[bool]contentType{
					true:  {"", 0, true, -1.1, -1.1, zeroTime},
					false: {"", 0, false, -1.1, -1.1, zeroTime},
				},
			},
		},
		{
			"MM",
			[]ArrayColumnInnerStruct2{
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{1}, 2, 1.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{2}, 4, 2.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{3}, 5, 3.1},
			},
			map[int]ArrayColumnInnerStruct2{
				2: ArrayColumnInnerStruct2{ArrayColumnInnerStruct{1}, 2, 1.1},
				4: ArrayColumnInnerStruct2{ArrayColumnInnerStruct{2}, 4, 2.1},
				5: ArrayColumnInnerStruct2{ArrayColumnInnerStruct{3}, 5, 3.1},
			},
		},
		{
			"DD",
			[]ArrayColumnInnerStruct2{
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{1}, 2, 1.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{2}, 4, 2.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{3}, 5, 3.1},
			},
			map[float32]ArrayColumnInnerStruct2{
				1.1: ArrayColumnInnerStruct2{ArrayColumnInnerStruct{1}, 2, 1.1},
				2.1: ArrayColumnInnerStruct2{ArrayColumnInnerStruct{2}, 4, 2.1},
				3.1: ArrayColumnInnerStruct2{ArrayColumnInnerStruct{3}, 5, 3.1},
			},
		},
		{
			"ArrayColumnInnerStruct.MM",
			[]ArrayColumnInnerStruct2{
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{1}, 2, 1.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{2}, 4, 2.1},
				ArrayColumnInnerStruct2{ArrayColumnInnerStruct{3}, 5, 3.1},
			},
			map[int]ArrayColumnInnerStruct2{
				1: ArrayColumnInnerStruct2{ArrayColumnInnerStruct{1}, 2, 1.1},
				2: ArrayColumnInnerStruct2{ArrayColumnInnerStruct{2}, 4, 2.1},
				3: ArrayColumnInnerStruct2{ArrayColumnInnerStruct{3}, 5, 3.1},
			},
		},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {

		result := ArrayColumnMap(singleTestCase.origin, singleTestCase.mapName)
		AssertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}

}
