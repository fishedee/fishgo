package test

import (
	"testing"
)

type ConfigAoTest struct {
	BaseTest
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

func TestConfigAo(t *testing.T) {
	InitTest(t, &ConfigAoTest{})
}
