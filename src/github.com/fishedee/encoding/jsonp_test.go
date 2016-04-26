package encoding

import (
	"reflect"
	"testing"
)

func assertJsonpEqual(t *testing.T, left interface{}, right interface{}) {
	isEqual := reflect.DeepEqual(left, right)
	if isEqual == false {
		t.Errorf("%#v != %#v", left, right)
	}
}

func TestJsonp(t *testing.T) {
	testCase := []struct {
		origin  string
		target  interface{}
		target2 string
	}{
		{`func1({"a":"b"})`, map[string]string{
			"a": "b",
		}, "func1"},
		{`  _mmc(  {"a":"b","c":13}  )   `, struct {
			A  string
			CC int `jsonp:"c"`
		}{
			"b",
			13,
		}, "_mmc"},
	}

	//测试解码
	for _, singleTestCase := range testCase {
		targetType := reflect.TypeOf(singleTestCase.target)
		singleResult := reflect.New(targetType)
		funcName, err := DecodeJsonp([]byte(singleTestCase.origin), singleResult.Interface())

		assertJsonpEqual(t, err, nil)
		assertJsonpEqual(t, singleTestCase.target, singleResult.Elem().Interface())
		assertJsonpEqual(t, singleTestCase.target2, funcName)

	}
	//测试编码后解码
	for _, singleTestCase := range testCase {
		singleResult, err := EncodeJsonp(singleTestCase.target2, singleTestCase.target)
		assertJsonpEqual(t, err, nil)

		targetType := reflect.TypeOf(singleTestCase.target)
		singleResult2 := reflect.New(targetType)
		funcName, err := DecodeJsonp(singleResult, singleResult2.Interface())
		assertJsonpEqual(t, err, nil)
		assertJsonpEqual(t, singleTestCase.target, singleResult2.Elem().Interface())
		assertJsonpEqual(t, singleTestCase.target2, funcName)
	}
}
