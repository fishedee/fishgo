package common

import (
	. "github.com/fishedee/web"
)

type BaseModel struct {
	BeegoValidateModel
}

type BaseTest struct {
	BeegoValidateTest
}

func InitTest(t *testing.T, test BeegoValidateTestInterface) {
	InitBeegoVaildateTest(t, test)
}
