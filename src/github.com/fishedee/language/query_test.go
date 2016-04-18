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

func TestQuerySelect(t *testing.T) {
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
		origin   interface{}
		function interface{}
		target   interface{}
	}{
		{
			[]contentType{},
			func(singleData contentType) contentType {
				return singleData
			},
			[]contentType{},
		},
		{
			[]contentType{
				contentType{"5", 1, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, nowTime},
			},
			func(singleData contentType) contentType {

				singleData.Name += "Edward"
				return singleData
			},
			[]contentType{
				contentType{"5Edward", 1, true, -1.1, -1.1, oldTime},
				contentType{"Edward", 0, false, 0, 0, zeroTime},
				contentType{"aEdward", -1, false, 1.1, 1.1, nowTime},
			},
		},
		{
			[]contentType{
				contentType{"5", 1, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, nowTime},
			},
			func(singleData contentType) contentType {

				singleData.Name += "Edward"
				return singleData
			},
			[]contentType{
				contentType{"5Edward", 1, true, -1.1, -1.1, oldTime},
				contentType{"Edward", 0, false, 0, 0, zeroTime},
				contentType{"aEdward", -1, false, 1.1, 1.1, nowTime},
			},
		},
		{
			[]contentType{
				contentType{"5", 1, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, nowTime},
			},
			func(singleData contentType) string {

				return singleData.Name
			},
			[]string{"5", "", "a"},
		},
		{
			[]contentType{
				contentType{"5", 1, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, nowTime},
			},
			func(singleData contentType) int {

				return singleData.Age
			},
			[]int{1, 0, -1},
		},
		{
			[]contentType{
				contentType{"5", 1, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, nowTime},
			},
			func(singleData contentType) bool {

				return singleData.Ok
			},
			[]bool{true, false, false},
		},
		{
			[]contentType{
				contentType{"5", 1, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, nowTime},
			},
			func(singleData contentType) float32 {

				return singleData.Money
			},
			[]float32{-1.1, 0, 1.1},
		},
		{
			[]contentType{
				contentType{"5", 1, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, nowTime},
			},
			func(singleData contentType) float64 {

				return singleData.CardMoney
			},
			[]float64{-1.1, 0, 1.1},
		},
		{
			[]contentType{
				contentType{"5", 1, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, nowTime},
			},
			func(singleData contentType) time.Time {

				return singleData.Register
			},
			[]time.Time{oldTime, zeroTime, nowTime},
		},
		{
			[]contentType{
				contentType{"5", 1, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"a", -1, false, 1.1, 1.1, nowTime},
			},
			func(singleData contentType) map[string]int {

				return map[string]int{singleData.Name: singleData.Age}
			},
			[]map[string]int{{"5": 1}, {"": 0}, {"a": -1}},
		},
	}

	// t.Error(QuerySelect(testCase[0].origin, testCase[0].function))

	for singleTestCaseIndex, singleTestCase := range testCase {

		result := QuerySelect(singleTestCase.origin, singleTestCase.function)
		assertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}
}

func TestQueryWhere(t *testing.T) {
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
		origin   interface{}
		function interface{}
		target   interface{}
	}{
		{
			[]contentType{},
			func(singleData contentType) bool {
				return true
			},
			[]contentType{},
		},
		{
			[]contentType{
				contentType{"s", 3, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			func(singleData contentType) bool {
				return singleData.Age >= 1
			},
			[]contentType{
				contentType{"s", 3, true, 0, 0, nowTime},
				contentType{"", 10, true, -1.1, -1.1, oldTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
		},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {

		result := QueryWhere(singleTestCase.origin, singleTestCase.function)
		assertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}

	// t.Error(QueryWhere(
	// 	[]contentType{
	// 		contentType{"s", 3, true, 0, 0, nowTime},
	// 		contentType{"a", -1, false, 1.1, 1.1, zeroTime},
	// 		contentType{"", 10, true, -1.1, -1.1, oldTime},
	// 		contentType{"", 0, false, 0, 0, zeroTime},
	// 		contentType{"z", 3, true, 0, 0, nowTime},
	// 	},
	// 	func(singleData contentType) bool {
	// 		return singleData.Age >= 1
	// 	}))
}

func TestQueryReduce(t *testing.T) {
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
		origin   interface{}
		function interface{}
		initNum  int
		target   int
	}{
		{
			[]contentType{},
			func(singleData contentType) bool {
				return true
			},
			0,
			0,
		},
		{
			[]contentType{
				contentType{"s", 3, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			func(sum int, singleData contentType) int {
				return singleData.Age + sum
			},
			0,
			15,
		},
		// {
		// 	[]contentType{
		// 		contentType{"s", 3, true, 0, 0, nowTime},
		// 		contentType{"a", -1, false, 1.1, 1.1, zeroTime},
		// 		contentType{"", 10, true, -1.1, -1.1, oldTime},
		// 		contentType{"", 0, false, 0, 0, zeroTime},
		// 		contentType{"z", 3, true, 0, 0, nowTime},
		// 	},
		// 	func(sum float32, singleData contentType) float32 {
		// 		return singleData.Money + sum
		// 	},
		// 	0,
		// 	15,
		// },
	}

	for singleTestCaseIndex, singleTestCase := range testCase {

		result := QueryReduce(singleTestCase.origin, singleTestCase.function, singleTestCase.initNum)
		assertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}

}

func TestQuerySort(t *testing.T) {
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

		result := QuerySort(singleTestCase.origin, singleTestCase.sortName)
		assertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}

}

func TestQueryJoin(t *testing.T) {
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
	t.Errorf("%+v", QueryJoin([]contentType{
		contentType{"s", 3, false, 0, 0, nowTime},
		contentType{"a", -1, true, 1.1, 1.1, zeroTime},
		contentType{"", 1, false, -1.1, -1.1, oldTime},
		contentType{"", 3, true, 0, 0, zeroTime},
		contentType{"z", 3, false, 0, 0, nowTime},
	}, []contentType{
		contentType{"s", 3, true, 0, 0, nowTime},
		contentType{"a", -1, false, 1.1, 1.1, zeroTime},
		contentType{"", 2, true, -1.1, -1.1, oldTime},
		contentType{"", 4, false, 0, 0, zeroTime},
		contentType{"z", 3, true, 0, 0, nowTime},
	}, "left", "  Name  =  Name ", func(left contentType, right contentType) contentType {
		return contentType{
			Name:      left.Name,
			Age:       left.Age,
			Ok:        left.Ok,
			Money:     left.Money,
			CardMoney: left.CardMoney,
			Register:  left.Register,
		}
	}))

	// [{Name:s Age:3 Ok:false Money:0 CardMoney:0 Register:2016-04-18 18:08:43.240014 +0800 CST}
	//  {Name:a Age:-1 Ok:true Money:1.1 CardMoney:1.1 Register:0001-01-01 00:00:00 +0000 UTC}
	//   {Name: Age:1 Ok:false Money:-1.1 CardMoney:-1.1 Register:2015-04-19 18:08:43.240014 +0800 CST}
	//   {Name: Age:1 Ok:false Money:-1.1 CardMoney:-1.1 Register:2015-04-19 18:08:43.240014 +0800 CST}
	//    {Name: Age:3 Ok:true Money:0 CardMoney:0 Register:0001-01-01 00:00:00 +0000 UTC}
	//     {Name: Age:3 Ok:true Money:0 CardMoney:0 Register:0001-01-01 00:00:00 +0000 UTC}
	// {Name:z Age:3 Ok:false Money:0 CardMoney:0 Register:2016-04-18 18:08:43.240014 +0800 CST}]
}
