package test

import (
	. "github.com/fishedee/web"
	"github.com/fishedee/web/test/models/number"
)

type innerTest struct {
	Test
	numberAoTest number.NumberAoTest
}

func (this *innerTest) TestBasic() {
	this.numberAoTest.TestBasic()
}
