package test

import (
	. "github.com/fishedee/web"
	"time"
)

type ConfigAoTest struct {
	Test
	ConfigAo ConfigAoModel
}

func (this *ConfigAoTest) testBasicEmpty(data []string) {
	for _, singleTestCase := range data {
		data := this.ConfigAo.Get(singleTestCase)
		this.AssertEqual(data, "")
	}
}

func (this *ConfigAoTest) TestBasic() {
	testCase := []struct {
		origin string
		target string
	}{
		{this.RandomString(128), this.RandomString(10240)},
		{this.RandomString(128), this.RandomString(10240)},
		{this.RandomString(128), this.RandomString(10240)},
	}

	noTestCase := []string{
		this.RandomString(128),
		this.RandomString(128),
	}

	//清空
	for _, singleTestCase := range testCase {
		data := this.ConfigAo.Get(singleTestCase.origin)
		this.AssertEqual(data, "")
	}
	this.testBasicEmpty(noTestCase)

	//设置
	for _, singleTestCase := range testCase {
		this.ConfigAo.Set(singleTestCase.origin, singleTestCase.target)
	}
	this.testBasicEmpty(noTestCase)

	//获取
	for _, singleTestCase := range testCase {
		data := this.ConfigAo.Get(singleTestCase.origin)
		this.AssertEqual(data, singleTestCase.target)
	}
	this.testBasicEmpty(noTestCase)
}

func (this *ConfigAoTest) TestStruct() {
	this.Log.Debug("This is a log %v", "123")
	this.Log.Debug("This is a log %v", ConfigData{})
	//struct中的time.Time字段不比较
	data1 := ConfigData{
		Data: "123",
	}
	this.ConfigAo.SetStruct("test1", data1)
	this.AssertEqual(this.ConfigAo.GetStruct("test1"), data1)

	data2 := ConfigData{
		Data:       "123",
		CreateTime: time.Now().AddDate(0, -1, 0),
		ModifyTime: time.Now().AddDate(0, -1, 0),
	}
	this.AssertEqual(this.ConfigAo.GetStruct("test1"), data2)
	this.AssertEqual(data1, data2)

	//struct中的非time.Time字段会比较
	data3 := ConfigData{
		Data: "789",
	}
	this.AssertEqual(this.ConfigAo.GetStruct("test1"), data3)

	//struct里面的数组与映射nil与非nil比较
	this.AssertEqual([]int(nil), []int{})
	this.AssertEqual(map[string]string(nil), map[string]string{})
}

func init() {
	InitTest(&ConfigAoTest{})
}
