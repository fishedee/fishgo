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
		target   interface{}
	}{
		{
			[]contentType{},
			func(sum int, singleData contentType) int {
				return 1
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
		{
			[]contentType{
				contentType{"s", 3, true, 0, 0, nowTime},
				contentType{"a", -1, false, 2.2, 1.1, zeroTime},
				contentType{"", 10, true, -1.1, -1.1, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			func(sum float32, singleData contentType) float32 {
				return singleData.Money + sum
			},
			0,
			(float32)(1.1),
		},
		{
			[]contentType{
				contentType{"s", 3, true, 0, 0, nowTime},
				contentType{"a", -1, false, 2.2, 1.1, zeroTime},
				contentType{"", 10, true, -1.1, -2.2, oldTime},
				contentType{"", 0, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			func(sum float64, singleData contentType) float64 {
				return singleData.CardMoney + sum
			},
			0,
			-1.1,
		},
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
	type userType struct {
		Name      string
		Age       int
		Ok        bool
		Money     float32
		CardMoney float64
		Register  time.Time
	}

	type contentType struct {
		UserName string
		Title    string
		Content  string
	}

	type resultType struct {
		UserName  string
		Age       int
		Ok        bool
		Money     float32
		CardMoney float64
		Register  time.Time
		Title     string
		Content   string
	}

	nowTime := time.Now()
	oldTime := nowTime.AddDate(-1, 0, 1)
	zeroTime := time.Time{}

	testCase := []struct {
		leftData   interface{}
		rightData  interface{}
		joinPlace  string
		joinType   string
		joinFuctor interface{}
		target     interface{}
	}{
		{

			[]userType{},
			[]userType{},
			" left ",
			"  Name  =  Name ",
			func(left userType, right userType) userType {
				return userType{}
			},
			[]userType{},
		},

		{

			[]userType{
				userType{"edward", 0, false, 1.1, 0, nowTime},
				userType{"fish", -1, true, -1.1, -1.1, zeroTime},
				userType{"jd", 1, false, 0, 1.1, oldTime},
			},
			[]contentType{
				contentType{"edward", "威风蛋糕", "威风蛋糕好好吃野！"},
				// contentType{"weinmey", "曲奇制作", "制作方法非常简单"},
				// contentType{"jd", "马卡龙", "好吃好玩"},
			},
			"left",
			"  Name  =  UserName ",
			func(left userType, right contentType) resultType {
				return resultType{
					UserName:  left.Name,
					Age:       left.Age,
					Ok:        left.Ok,
					Money:     left.Money,
					CardMoney: left.CardMoney,
					Register:  left.Register,
					Title:     right.Title,
					Content:   right.Content,
				}
			},
			[]resultType{
				resultType{"edward", 0, false, 1.1, 0, nowTime, "威风蛋糕", "威风蛋糕好好吃野！"},
				resultType{"fish", -1, true, -1.1, -1.1, zeroTime, "", ""},
				resultType{"jd", 1, false, 0, 1.1, oldTime, "", ""},
			},
		},
		{

			[]userType{
				userType{"edward", 0, false, 1.1, 0, nowTime},
				// userType{"fish", -1, true, -1.1, -1.1, zeroTime},
				// userType{"jd", 1, false, 0, 1.1, oldTime},
			},
			[]contentType{
				contentType{"edward", "威风蛋糕", "威风蛋糕好好吃野！"},
				contentType{"weinmey", "曲奇制作", "制作方法非常简单"},
				contentType{"jd", "马卡龙", "好吃好玩"},
			},
			"left",
			"  Name  =  UserName ",
			func(left userType, right contentType) resultType {
				return resultType{
					UserName:  left.Name,
					Age:       left.Age,
					Ok:        left.Ok,
					Money:     left.Money,
					CardMoney: left.CardMoney,
					Register:  left.Register,
					Title:     right.Title,
					Content:   right.Content,
				}
			},
			[]resultType{
				resultType{"edward", 0, false, 1.1, 0, nowTime, "威风蛋糕", "威风蛋糕好好吃野！"},
			},
		},
		{

			[]userType{
				userType{"edward", 0, false, 1.1, 0, nowTime},
				// userType{"fish", -1, true, -1.1, -1.1, nowTime},
				// userType{"jd", 1, false, 0, 1.1, oldTime},
			},
			[]contentType{
				contentType{"edward", "威风蛋糕", "威风蛋糕好好吃野！"},
				contentType{"weinmey", "曲奇制作", "制作方法非常简单"},
				contentType{"jd", "马卡龙", "好吃好玩"},
			},
			"right",
			"  Name  =  UserName ",
			func(left userType, right contentType) resultType {
				return resultType{
					UserName:  left.Name,
					Age:       left.Age,
					Ok:        left.Ok,
					Money:     left.Money,
					CardMoney: left.CardMoney,
					Register:  left.Register,
					Title:     right.Title,
					Content:   right.Content,
				}
			},
			[]resultType{
				resultType{"edward", 0, false, 1.1, 0, nowTime, "威风蛋糕", "威风蛋糕好好吃野！"},
				resultType{"", 0, false, 0, 0, zeroTime, "马卡龙", "好吃好玩"},
				resultType{"", 0, false, 0, 0, zeroTime, "曲奇制作", "制作方法非常简单"},
			},
		},
		{

			[]userType{
				userType{"edward", 0, false, 1.1, 0, nowTime},
				userType{"fish", -1, true, -1.1, -1.1, nowTime},
				userType{"jd", 1, false, 0, 1.1, oldTime},
			},
			[]contentType{
				contentType{"edward", "威风蛋糕", "威风蛋糕好好吃野！"},
				// contentType{"weinmey", "曲奇制作", "制作方法非常简单"},
				// contentType{"jd", "马卡龙", "好吃好玩"},
			},
			"right",
			"  Name  =  UserName ",
			func(left userType, right contentType) resultType {
				return resultType{
					UserName:  left.Name,
					Age:       left.Age,
					Ok:        left.Ok,
					Money:     left.Money,
					CardMoney: left.CardMoney,
					Register:  left.Register,
					Title:     right.Title,
					Content:   right.Content,
				}
			},
			[]resultType{
				resultType{"edward", 0, false, 1.1, 0, nowTime, "威风蛋糕", "威风蛋糕好好吃野！"},
			},
		},
		{

			[]userType{
				userType{"edward", 0, false, 1.1, 0, nowTime},
				userType{"fish", -1, true, -1.1, -1.1, nowTime},
				userType{"jd", 1, false, 0, 1.1, oldTime},
			},
			[]contentType{
				contentType{"edward", "威风蛋糕", "威风蛋糕好好吃野！"},
				// contentType{"weinmey", "曲奇制作", "制作方法非常简单"},
				// contentType{"jd", "马卡龙", "好吃好玩"},
			},
			"inner",
			"  Name  =  UserName ",
			func(left userType, right contentType) resultType {
				return resultType{
					UserName:  left.Name,
					Age:       left.Age,
					Ok:        left.Ok,
					Money:     left.Money,
					CardMoney: left.CardMoney,
					Register:  left.Register,
					Title:     right.Title,
					Content:   right.Content,
				}
			},
			[]resultType{
				resultType{"edward", 0, false, 1.1, 0, nowTime, "威风蛋糕", "威风蛋糕好好吃野！"},
			},
		},
		{

			[]userType{
				userType{"edward", 0, false, 1.1, 0, nowTime},
				userType{"fish", -1, true, -1.1, -1.1, nowTime},
				userType{"jd", 1, false, 0, 1.1, oldTime},
			},
			[]contentType{
				contentType{"edward", "威风蛋糕", "威风蛋糕好好吃野！"},
				contentType{"weinmey", "曲奇制作", "制作方法非常简单"},
				contentType{"jd", "马卡龙", "好吃好玩"},
			},
			"outer",
			"  Name  =  UserName ",
			func(left userType, right contentType) resultType {
				return resultType{
					UserName:  right.UserName,
					Age:       left.Age,
					Ok:        left.Ok,
					Money:     left.Money,
					CardMoney: left.CardMoney,
					Register:  left.Register,
					Title:     right.Title,
					Content:   right.Content,
				}
			},
			[]resultType{
				resultType{"edward", 0, false, 1.1, 0, nowTime, "威风蛋糕", "威风蛋糕好好吃野！"},
				resultType{"", -1, true, -1.1, -1.1, nowTime, "", ""},
				resultType{"jd", 1, false, 0, 1.1, oldTime, "马卡龙", "好吃好玩"},
				resultType{"weinmey", 0, false, 0, 0, zeroTime, "曲奇制作", "制作方法非常简单"},
			},
		},

		{
			[]userType{
				userType{"s", 0, false, 1.1, 0, oldTime},
				userType{"a", -1, true, -1.1, -1.1, nowTime},
				userType{"", 1, false, 1.1, 1.1, oldTime},
				userType{"", 0, true, 0, 0, zeroTime},
				userType{"z", -1, false, 0, 0, oldTime},
			},
			[]userType{
				userType{"s", -1, true, 0, 0, nowTime},
				userType{"a", 0, false, 1.1, 1.1, zeroTime},
				userType{"", -1, true, -1.1, -1.1, oldTime},
				userType{"", 1, false, 1, 1, zeroTime},
				userType{"z", 1, true, -1, -1, nowTime},
			},
			"right",
			"Age=Age",
			func(left userType, right userType) userType {
				return userType{
					Name:      left.Name,
					Age:       right.Age,
					Ok:        left.Ok,
					Money:     right.Money,
					CardMoney: left.CardMoney,
					Register:  right.Register,
				}
			},
			[]userType{
				userType{"s", 0, false, 1.1, 0, zeroTime},
				userType{"a", -1, true, 0, -1.1, nowTime},
				userType{"a", -1, true, -1.1, -1.1, oldTime},
				userType{"", 1, false, 1, 1.1, zeroTime},
				userType{"", 1, false, -1, 1.1, nowTime},
				userType{"", 0, true, 1.1, 0, zeroTime},
				userType{"z", -1, false, 0, 0, nowTime},
				userType{"z", -1, false, -1.1, 0, oldTime},
			},
		},
		{
			[]userType{
				userType{"s", 0, false, 1.1, 0, oldTime},
				userType{"a", -1, true, -1.1, -1.1, nowTime},
				userType{"", 1, false, 1.1, 1.1, oldTime},
			},
			[]userType{
				userType{"s", -1, true, 0, 0, nowTime},
				userType{"a", 0, false, 1.1, 1.1, zeroTime},
			},
			"left",
			"Ok  =  Ok",
			func(left userType, right userType) userType {
				return userType{
					Name:      left.Name,
					Age:       right.Age,
					Ok:        left.Ok,
					Money:     right.Money,
					CardMoney: left.CardMoney,
					Register:  right.Register,
				}
			},
			[]userType{
				userType{"s", 0, false, 1.1, 0, zeroTime},
				userType{"a", -1, true, 0, -1.1, nowTime},
				userType{"", 0, false, 1.1, 1.1, zeroTime},
			},
		},
		{
			[]userType{
				userType{"s", 0, false, 1.1, 0, oldTime},
				userType{"a", -1, true, -1.1, -1.1, nowTime},
				userType{"", 1, false, 0, 1.1, oldTime},
			},
			[]userType{
				userType{"s", -1, true, 0, 0, nowTime},
				userType{"a", 0, false, 1.1, 1.1, zeroTime},
			},
			"left",
			" Money=Money ",
			func(left userType, right userType) userType {
				return userType{
					Name:      left.Name,
					Age:       right.Age,
					Ok:        left.Ok,
					Money:     right.Money,
					CardMoney: left.CardMoney,
					Register:  right.Register,
				}
			},
			[]userType{
				userType{"s", 0, false, 1.1, 0, zeroTime},
				userType{"a", 0, true, 0, -1.1, zeroTime},
				userType{"", -1, false, 0, 1.1, nowTime},
			},
		},
		{
			[]userType{
				userType{"s", 0, false, 1.1, 0, oldTime},
				userType{"a", -1, true, -1.1, -1.1, nowTime},
				userType{"", 1, false, 0, 1.1, oldTime},
			},
			[]userType{
				userType{"s", -1, true, 0, 0, nowTime},
				userType{"a", 0, false, 1.1, 1.1, zeroTime},
			},
			"left",
			" CardMoney = Money ",
			func(left userType, right userType) userType {
				return userType{
					Name:      left.Name,
					Age:       right.Age,
					Ok:        left.Ok,
					Money:     right.Money,
					CardMoney: left.CardMoney,
					Register:  right.Register,
				}
			},
			[]userType{
				userType{"s", -1, false, 0, 0, nowTime},
				userType{"a", 0, true, 0, -1.1, zeroTime},
				userType{"", 0, false, 0, 1.1, zeroTime},
			},
		},
		{
			[]userType{
				userType{"s", 0, false, 1.1, 0, oldTime},
				userType{"a", -1, true, -1.1, -1.1, nowTime},
				userType{"", 1, false, 0, 1.1, oldTime},
			},
			[]userType{
				userType{"s", -1, true, 0, 0, nowTime},
				userType{"a", 0, false, 1.1, 1.1, zeroTime},
			},
			"left",
			" Register = Register ",
			func(left userType, right userType) userType {
				return userType{
					Name:      left.Name,
					Age:       right.Age,
					Ok:        left.Ok,
					Money:     right.Money,
					CardMoney: left.CardMoney,
					Register:  right.Register,
				}
			},
			[]userType{
				userType{"s", 0, false, 0, 0, zeroTime},
				userType{"a", -1, true, 0, -1.1, nowTime},
				userType{"", 0, false, 0, 1.1, zeroTime},
			},
		},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {

		result := QueryJoin(singleTestCase.leftData, singleTestCase.rightData, singleTestCase.joinPlace, singleTestCase.joinType, singleTestCase.joinFuctor)
		assertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}
}

func TestQueryGroup(t *testing.T) {

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
		data        interface{}
		groupType   string
		groupFuctor interface{}
		target      interface{}
	}{

		{
			[]contentType{},
			" Ok ",
			func(list []contentType) []contentType {
				return []contentType{}
			},
			[]contentType{},
		},
		{
			[]contentType{
				contentType{"a", 3, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 10, true, -2.2, -1.2, oldTime},
				contentType{"", -2, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			"Name",
			func(list []contentType) []contentType {
				sum := QuerySum(QueryColumn(list, "  Money  "))
				list[0].Money = sum.(float32)
				return []contentType{list[0]}
			},
			[]contentType{
				contentType{"", 10, true, -2.2, -1.2, oldTime},
				contentType{"a", 3, true, 1.1, 0, nowTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
		},
		{
			[]contentType{
				contentType{"a", 3, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 10, true, -2.2, -1.2, oldTime},
				contentType{"", -2, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			"Name",
			func(list []contentType) float32 {
				return QuerySum(QueryColumn(list, "  Money  ")).(float32)
			},
			[]float32{-2.2, 1.1, 0},
		},
		{
			[]contentType{
				contentType{"a", 3, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 10, true, -2.2, -1.2, oldTime},
				contentType{"", -2, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			"Ok",
			func(list []contentType) []contentType {
				sum := QuerySum(QueryColumn(list, "CardMoney  "))
				list[0].CardMoney = sum.(float64)
				return []contentType{list[0]}
			},
			[]contentType{
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"a", 3, true, 0, -1.2, nowTime},
			},
		},
		{
			[]contentType{
				contentType{"s", 1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 0, true, -1.1, -1.1, oldTime},
				contentType{"", -1, false, 0, 0, zeroTime},
				contentType{"z", 1, true, 0, 0, nowTime},
			},
			" Age ",
			func(list []contentType) []contentType {
				sum := QuerySum(QueryColumn(list, "  CardMoney  "))
				list[0].CardMoney = sum.(float64)
				return []contentType{list[0]}
			},
			[]contentType{
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 0, true, -1.1, -1.1, oldTime},
				contentType{"s", 1, true, 0, 0, nowTime},
			},
		},
		{
			[]contentType{
				contentType{"s", 1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 0, true, -1.1, -1.1, oldTime},
				contentType{"", -1, false, 0, 0, zeroTime},
				contentType{"z", 1, true, 0, 0, nowTime},
			},
			" Age ",
			func(list []contentType) float64 {
				return QuerySum(QueryColumn(list, "  CardMoney  ")).(float64)

			},
			[]float64{1.1, -1.1, 0},
		},
		{
			[]contentType{
				contentType{"s", 1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 0, true, -1.1, -1.1, oldTime},
				contentType{"", -1, false, 0, 0, zeroTime},
				contentType{"z", 1, true, 0, 0, nowTime},
			},
			" Age ",
			func(list []contentType) []float64 {
				sum := QuerySum(QueryColumn(list, "  CardMoney  "))
				return []float64{sum.(float64)}

			},
			[]float64{1.1, -1.1, 0},
		},
		{
			[]contentType{
				contentType{"s", 1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 0, true, -1.1, -1.1, oldTime},
				contentType{"", -1, false, 0, 0, zeroTime},
				contentType{"z", 1, true, 0, 0, nowTime},
			},
			"Register ",
			func(list []contentType) int {
				sum := QuerySum(QueryColumn(list, "  Age  "))
				return sum.(int)

			},
			[]int{-2, 0, 2},
		},
		{
			[]contentType{
				contentType{"s", 1, true, 0, 0, nowTime},
				contentType{"a", -1, false, 1.1, 1.1, zeroTime},
				contentType{"", 0, true, -1.1, -1.1, oldTime},
				contentType{"", -1, false, 0, 0, zeroTime},
				contentType{"z", 1, true, 0, 0, nowTime},
			},
			"Register ",
			func(list []contentType) []contentType {
				sum := QuerySum(QueryColumn(list, "  Age  "))
				list[0].Age = sum.(int)
				return []contentType{list[0]}
			},
			[]contentType{
				contentType{"a", -2, false, 1.1, 1.1, zeroTime},
				contentType{"", 0, true, -1.1, -1.1, oldTime},
				contentType{"s", 2, true, 0, 0, nowTime},
			},
		},
		{
			[]contentType{
				contentType{"s", 1, true, 0, 0, nowTime},
				contentType{"s", 1, true, 6.6, 6.6, nowTime},
				contentType{"", 0, true, -5.1, -5.1, oldTime},
				contentType{"", 0, true, 2.1, 2.1, oldTime},
				contentType{"", -1, false, -3.3, -3.3, zeroTime},
				contentType{"", -1, false, 4.3, 4.3, zeroTime},
			},
			" Name , Ok ",
			func(list []contentType) []contentType {
				sum := QuerySum(QueryColumn(list, "  Age  "))
				list[0].Age = sum.(int)
				return []contentType{list[0]}
			},
			[]contentType{
				contentType{"", -2, false, -3.3, -3.3, zeroTime},
				contentType{"", 0, true, -5.1, -5.1, oldTime},
				contentType{"s", 2, true, 0, 0, nowTime},
			},
		},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {

		result := QueryGroup(singleTestCase.data, singleTestCase.groupType, singleTestCase.groupFuctor)
		assertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}

}

func TestQueryColumn(t *testing.T) {

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
		data   interface{}
		Column string
		target interface{}
	}{

		{
			[]contentType{},
			" Name ",
			[]string{},
		},
		{
			[]contentType{
				contentType{"a", 3, true, 0, 0, nowTime},
				contentType{"0", -1, false, 1.1, 1.1, zeroTime},
				contentType{"1", 10, true, -2.2, -1.2, oldTime},
				contentType{"-1", -2, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			"     Name         ",
			[]string{"a", "0", "1", "-1", "z"},
		},
		{
			[]contentType{
				contentType{"a", 3, true, 0, 0, nowTime},
				contentType{"0", -1, false, 1.1, 1.1, zeroTime},
				contentType{"1", 10, true, -2.2, -1.2, oldTime},
				contentType{"-1", -2, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			"Age        ",
			[]int{3, -1, 10, -2, 3},
		},
		{
			[]contentType{
				contentType{"a", 3, true, 0, 0, nowTime},
				contentType{"0", -1, false, 1.1, 1.1, zeroTime},
				contentType{"1", 10, true, -2.2, -1.2, oldTime},
				contentType{"-1", -2, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			"Ok        ",
			[]bool{true, false, true, false, true},
		},
		{
			[]contentType{
				contentType{"a", 3, true, 0, 0, nowTime},
				contentType{"0", -1, false, 1.1, 1.1, zeroTime},
				contentType{"1", 10, true, -2.2, -1.2, oldTime},
				contentType{"-1", -2, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			"    Money  ",
			[]float32{0, 1.1, -2.2, 0, 0},
		},
		{
			[]contentType{
				contentType{"a", 3, true, 0, 0, nowTime},
				contentType{"0", -1, false, 1.1, 1.1, zeroTime},
				contentType{"1", 10, true, -2.2, -1.2, oldTime},
				contentType{"-1", -2, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			"    CardMoney",
			[]float64{0, 1.1, -1.2, 0, 0},
		},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {

		result := QueryColumn(singleTestCase.data, singleTestCase.Column)
		assertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}
}

func TestQueryReverse(t *testing.T) {

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
		data   interface{}
		target interface{}
	}{

		{
			[]contentType{},
			[]contentType{},
		},
		{
			[]contentType{
				contentType{"a", 3, true, 0, 0, nowTime},
				contentType{"0", -1, false, 1.1, 1.1, zeroTime},
				contentType{"1", 10, true, -2.2, -1.2, oldTime},
				contentType{"-1", -2, false, 0, 0, zeroTime},
				contentType{"z", 3, true, 0, 0, nowTime},
			},
			[]contentType{
				contentType{"z", 3, true, 0, 0, nowTime},
				contentType{"-1", -2, false, 0, 0, zeroTime},
				contentType{"1", 10, true, -2.2, -1.2, oldTime},
				contentType{"0", -1, false, 1.1, 1.1, zeroTime},
				contentType{"a", 3, true, 0, 0, nowTime},
			},
		},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {

		result := QueryReverse(singleTestCase.data)
		assertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}
}

func TestQueryDistinct(t *testing.T) {

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

		result := QueryDistinct(singleTestCase.origin, singleTestCase.uniqueName)
		assertEqual(t, result, singleTestCase.target, singleTestCaseIndex)

	}
}
