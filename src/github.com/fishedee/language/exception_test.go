package language_test

import (
	. "github.com/fishedee/language"
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

func getLastStackTraceLine(e Exception) string {
	lines := Explode(e.GetStackTraceLine(0), "/")
	return lines[len(lines)-1]
}

func TestCatchStack1(t *testing.T) {
	defer Catch(func(e Exception) {
		AssertEqual(t, getLastStackTraceLine(e), "exception_test.go:62")
	})
	Throw(1, "test1")
}

func TestCatchStack2(t *testing.T) {
	defer CatchCrash(func(e Exception) {
		AssertEqual(t, getLastStackTraceLine(e), "exception_test.go:69")
	})
	Throw(1, "test2")
}

func TestCatchStack3(t *testing.T) {
	defer CatchCrash(func(e Exception) {
		AssertEqual(t, getLastStackTraceLine(e), "exception_test.go:76")
	})
	panic("test3")
}

func TestCatchStack4(t *testing.T) {
	defer CatchCrash(func(e Exception) {
		AssertEqual(t, getLastStackTraceLine(e), "exception_test.go:86")
	})
	defer Catch(func(e Exception) {
		AssertEqual(t, "should not be here!", false)
	})
	panic("test4")
}

func TestCatchStack5(t *testing.T) {
	defer CatchCrash(func(e Exception) {
		AssertEqual(t, getLastStackTraceLine(e), "exception_test.go:96")
	})
	defer Catch(func(e Exception) {
		panic(&e)
	})
	Throw(1, "test5")
}

func TestCatchStack6(t *testing.T) {
	defer Catch(func(e Exception) {
		AssertEqual(t, getLastStackTraceLine(e), "exception_test.go:106")
	})
	defer Catch(func(e Exception) {
		panic(&e)
	})
	Throw(1, "test6")
}

func TestCatchStack7(t *testing.T) {
	defer CatchCrash(func(e Exception) {
		AssertEqual(t, getLastStackTraceLine(e), "exception_test.go:116")
	})
	defer CatchCrash(func(e Exception) {
		panic(&e)
	})
	panic("test7")
}
