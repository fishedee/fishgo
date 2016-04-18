package test

import (
	. "github.com/fishedee/web"
	"testing"
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
