package language

import (
	. "github.com/fishedee/assert"
	"testing"
)

func getCatchMessage(fun func()) (_last string) {
	defer Catch(func(e Exception) {
		_last = e.GetMessage()
	})
	fun()
	return ""
}

func getCatchCrashMessage(fun func()) (_last string) {
	defer CatchCrash(func(e Exception) {
		_last = e.GetMessage()
	})
	fun()
	return ""
}

type errorStruct struct {
}

func (this *errorStruct) Error() string {
	return "m2"
}

func TestCatch(t *testing.T) {
	testCase := []struct {
		origin func()
		target string
	}{
		{func() {
			panic("m1")
		}, "m1"},
		{func() {
			panic(&errorStruct{})
		}, "m2"},
		{func() {
			Throw(1, "m3")
		}, "m3"},
	}

	for singleIndex, singleTestCase := range testCase {
		msg := getCatchCrashMessage(singleTestCase.origin)
		AssertEqual(t, msg, singleTestCase.target, singleIndex)
	}
}

//FIXME check exception stack
