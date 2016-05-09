package test

import (
	. "github.com/fishedee/web"
)

type extendAoTest struct {
	Test
	extendAo ExtendAoModel
}

func (this *extendAoTest) TestBasic() {
	configAo := this.extendAo.BaseAoModel.ConfigAo
	configAo.Set("mm1", "mm2")
	this.AssertEqual("mm2", configAo.Get("mm1"))
}
