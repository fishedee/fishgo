package encoding

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func AssertEqual(t *testing.T, left interface{}, right interface{}) {
	isEqual := reflect.DeepEqual(left, right)
	if isEqual == false {
		t.Error(fmt.Sprintf("%#v != %#v", left, right))
	}
}

func TestQueryDecodeBasic(t *testing.T) {
	testCase := []struct {
		origin string
		target interface{}
	}{
		//basic
		{"", ""},
		{"true", true},
		{"false", false},
		{"-1", -1},
		{"1", 1},
		{"abc", "abc"},
		//map
		{"a=3", map[string]string{
			"a": "3",
		}},
		{"a=3&b=4", map[string]string{
			"a": "3",
			"b": "4",
		}},
		{"a=3&b=4", map[string]int{
			"a": 3,
			"b": 4,
		}},
		//map slice
		{"a=3&b[]=4&b[]=5", map[string]interface{}{
			"a": "3",
			"b": []interface{}{"4", "5"},
		}},
		//map slice slice
		{"a=3&b[0][]=4&b[0][]=5&b[1][]=6&b[1][]=7", map[string]interface{}{
			"a": "3",
			"b": []interface{}{
				[]interface{}{"4", "5"},
				[]interface{}{"6", "7"},
			},
		}},
		//map slice map
		{"a=3&b[0][m1]=4&b[0][m2]=5&b[1][m1]=6&b[1][m2]=7", map[string]interface{}{
			"a": "3",
			"b": []interface{}{
				map[string]interface{}{
					"m1": "4",
					"m2": "5",
				},
				map[string]interface{}{
					"m1": "6",
					"m2": "7",
				},
			},
		}},
		//something wrong
		{"a=32&&", map[string]int{
			"a": 32,
		}},
		{"&a=32", map[string]int{
			"a": 32,
		}},
		{"&a=32&", map[string]int{
			"a": 32,
		}},
		{"&a=32&=", map[string]int{
			"a": 32,
		}},
		{"a=3&b=", map[string]int{
			"a": 3,
			"b": 0,
		}},
		{"a=3&=b", map[string]int{
			"a": 3,
		}},
		{"&a=3&b=4&", map[string]int{
			"a": 3,
			"b": 4,
		}},
		{"a=3&b=4&", map[string]int{
			"a": 3,
			"b": 4,
		}},
		//struct
		{"a=3&b[0][m1]=4&b[0][m2]=5&b[1][m1]=6&b[1][m2]=7&_c=mc", struct {
			AA string `url:"a"`
			B  []struct {
				M1 int
				M2 string
			}
			C string `url:"_c"`
			D []int
		}{
			"3",
			[]struct {
				M1 int
				M2 string
			}{
				{4, "5"},
				{6, "7"},
			},
			"mc",
			nil,
		}},
		//中文
		{"a=%e4%bd%a0%e5%a5%bd&b=%e4%b8%ad%e5%9b%bd", struct {
			AA string `url:"a"`
			BB string `url:"b"`
		}{
			"你好",
			"中国",
		}},
	}
	//测试解码
	for _, singleTestCase := range testCase {
		targetType := reflect.TypeOf(singleTestCase.target)
		singleResult := reflect.New(targetType)
		err := DecodeUrlQuery([]byte(singleTestCase.origin), singleResult.Interface())

		AssertEqual(t, err, nil)
		AssertEqual(t, singleTestCase.target, singleResult.Elem().Interface())

	}
	//测试编码后解码
	for _, singleTestCase := range testCase {
		singleResult, err := EncodeUrlQuery(singleTestCase.target)
		AssertEqual(t, err, nil)

		targetType := reflect.TypeOf(singleTestCase.target)
		singleValueResult := reflect.New(targetType)
		err = DecodeUrlQuery(singleResult, singleValueResult.Interface())
		AssertEqual(t, err, nil)
		AssertEqual(t, singleTestCase.target, singleValueResult.Elem().Interface())
	}
}

func TestQueryDecodeError(t *testing.T) {
	testCase := []struct {
		origin string
		target interface{}
		err    string
	}{
		{"a=1c&b=4", struct {
			AA int    `url:"a"`
			BB string `url:"b"`
		}{}, "参数a不是整数，其值为[1c]"},
		{"a=1.2c&b=4", struct {
			AA float32 `url:"a"`
			BB string  `url:"b"`
		}{}, "参数a不是浮点数，其值为[1.2c]"},
		{"a=2018-09-09&b=4", struct {
			AA time.Time `url:"a"`
			BB string    `url:"b"`
		}{}, "参数a不是时间，其值为[2018-09-09]"},
	}
	//测试解码
	for _, singleTestCase := range testCase {
		targetType := reflect.TypeOf(singleTestCase.target)
		singleResult := reflect.New(targetType)

		err := DecodeUrlQuery([]byte(singleTestCase.origin), singleResult.Interface())
		AssertEqual(t, err != nil, true)
		AssertEqual(t, err.Error(), singleTestCase.err)
	}
}
