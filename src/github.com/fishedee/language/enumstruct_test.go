package language

import (
	"fmt"
	"reflect"
	"testing"
)

func AssertEnumStructEqual(t *testing.T, left interface{}, right interface{}) {
	isEqual := reflect.DeepEqual(left, right)
	if isEqual == false {
		t.Error(fmt.Sprintf("%#v != %#v", left, right))
	}
}

func TestEnumStruct(t *testing.T) {
	var testCase struct {
		EnumStruct
		ENUM1 int `enum:"1,枚举1"`
		ENUM2 int `enum:"2,枚举2"`
		ENUM3 int `enum:"3,枚举3"`
	}

	InitEnumStruct(&testCase)

	//断言基本枚举值
	AssertEnumStructEqual(t, testCase.ENUM1, 1)
	AssertEnumStructEqual(t, testCase.ENUM2, 2)
	AssertEnumStructEqual(t, testCase.ENUM3, 3)

	//断言函数
	AssertEnumStructEqual(t, testCase.Names(), map[string]string{
		"1": "枚举1",
		"2": "枚举2",
		"3": "枚举3",
	})
	AssertEnumStructEqual(t, testCase.Entrys(), map[int]string{
		1: "枚举1",
		2: "枚举2",
		3: "枚举3",
	})
	AssertEnumStructEqual(t, testCase.Datas(), []EnumData{
		{1, "枚举1"},
		{2, "枚举2"},
		{3, "枚举3"},
	})
	AssertEnumStructEqual(t, ArraySort(testCase.Keys()), []int{1, 2, 3})
	AssertEnumStructEqual(t, ArraySort(testCase.Values()), []string{"枚举1", "枚举2", "枚举3"})
}

func TestEnumStructString(t *testing.T) {
	var testCase struct {
		EnumStructString
		ENUM1 string `enum:"/content/del1,枚举1"`
		ENUM2 string `enum:"/content/del2,枚举2"`
		ENUM3 string `enum:"/content/del3,枚举3"`
	}

	InitEnumStructString(&testCase)

	//断言基本枚举值
	AssertEnumStructEqual(t, testCase.ENUM1, "/content/del1")
	AssertEnumStructEqual(t, testCase.ENUM2, "/content/del2")
	AssertEnumStructEqual(t, testCase.ENUM3, "/content/del3")

	//断言函数
	AssertEnumStructEqual(t, testCase.Names(), map[string]string{
		"/content/del1": "枚举1",
		"/content/del2": "枚举2",
		"/content/del3": "枚举3",
	})
	AssertEnumStructEqual(t, testCase.Entrys(), map[string]string{
		"/content/del1": "枚举1",
		"/content/del2": "枚举2",
		"/content/del3": "枚举3",
	})
	AssertEnumStructEqual(t, testCase.Datas(), []EnumDataString{
		{"/content/del1", "枚举1"},
		{"/content/del2", "枚举2"},
		{"/content/del3", "枚举3"},
	})
	AssertEnumStructEqual(t, ArraySort(testCase.Keys()), []string{"/content/del1", "/content/del2", "/content/del3"})
	AssertEnumStructEqual(t, ArraySort(testCase.Values()), []string{"枚举1", "枚举2", "枚举3"})
}
