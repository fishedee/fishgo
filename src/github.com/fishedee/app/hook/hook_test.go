package hook

import (
	. "github.com/fishedee/assert"
	"testing"
)

func TestHook(t *testing.T) {
	result := []interface{}{}
	testCase := []struct {
		name    string
		handler interface{}
		data    []interface{}
	}{
		{"m1", func() {
		}, []interface{}{}},
		{"m2", func(data int) {
			result = append(result, data)
		}, []interface{}{1}},
		{"m3", func(data string) {
			result = append(result, data)
		}, []interface{}{"abc"}},
		{"m4", func(data string, data2 int) {
			result = append(result, data)
			result = append(result, data2)
		}, []interface{}{"abcd", 3}},
	}

	hook, _ := NewHook()
	for singleTestCaseIndex, singleTestCase := range testCase {
		result = []interface{}{}
		hook.Register(singleTestCase.name, singleTestCase.handler)
		hook.Trigger(singleTestCase.name, singleTestCase.data...)
		AssertEqual(t, result, singleTestCase.data, singleTestCaseIndex)
	}
}

func TestHookMulti(t *testing.T) {
	hook, _ := NewHook()

	counter := 0
	hook.Register("/m1", func() {
		counter++
	})
	hook.Register("/m1", func(data int) {
		counter++
	})
	hook.Trigger("/m1", 123)
	AssertEqual(t, counter, 2)
}

func TestHookReturnValue(t *testing.T) {
	hook, _ := NewHook()
	hook.Register("/m1", func() {
	})
	hook.Register("/m1", func() int {
		return 2
	})
	hook.Register("/m1", func() (string, int) {
		return "3", 4
	})
	result := hook.Trigger("/m1")
	AssertEqual(t, result, []interface{}{nil, 2, "3"})
}
