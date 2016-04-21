package test

import (
	"github.com/fishedee/web/test/models/number"
)

type InnerTest struct {
	BaseTest
	NumberAoTest number.NumberAoTest
}

func (this *InnerTest) TestBasic() {
	this.NumberAoTest.TestBasic()
}

func init() {
	InitTest(&InnerTest{})
}
