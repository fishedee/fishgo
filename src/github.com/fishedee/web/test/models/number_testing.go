package test

import (
	. "github.com/fishedee/web"
	"github.com/fishedee/web/test/models/number"
)

type InnerTest struct {
	Test
	NumberAoTest number.NumberAoTest
}

func (this *InnerTest) TestBasic() {
	this.NumberAoTest.TestBasic()
}

func init() {
	InitTest(&InnerTest{})
}
