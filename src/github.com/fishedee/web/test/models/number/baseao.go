package number

import (
	. "github.com/fishedee/web"
)

type BaseModel struct {
	BeegoValidateModel
}

type BaseTest struct {
	BeegoValidateTest
}

func InitTest(test BeegoValidateTestInterface) {
	InitBeegoVaildateTest(test)
}
