package test

import (
	. "github.com/fishedee/web"
	"github.com/fishedee/web/test/models/number"
)

type innerTest struct {
	Test
	NumberAoTest number.NumberAoTest
}

func (this *innerTest) TestBasic() {
	this.NumberAoTest.TestBasic()
}
