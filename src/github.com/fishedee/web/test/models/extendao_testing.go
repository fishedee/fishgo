package test

import ()

type ExtendAoTest struct {
	BaseTest
	ExtendAo ExtendAoModel
}

func (this *ExtendAoTest) TestBasic() {
	configAo := this.ExtendAo.BaseAoModel.ConfigAo
	configAo.Set("mm1", "mm2")
	this.AssertEqual("mm2", configAo.Get("mm1"))
}

func init() {
	InitTest(&ExtendAoTest{})
}
