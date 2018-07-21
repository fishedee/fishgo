package trigger

import (
	. "github.com/fishedee/assert"
	"testing"
)

func TestTrigger(t *testing.T) {
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

	trigger, _ := NewTrigger()
	for singleTestCaseIndex, singleTestCase := range testCase {
		result = []interface{}{}
		trigger.On(singleTestCase.name, singleTestCase.handler)
		trigger.Fire(singleTestCase.name, singleTestCase.data...)
		AssertEqual(t, result, singleTestCase.data, singleTestCaseIndex)
	}
}

func TestTriggerMulti(t *testing.T) {
	trigger, _ := NewTrigger()

	counter := 0
	trigger.On("/m1", func() {
		counter++
	})
	trigger.On("/m1", func(data int) {
		counter++
	})
	trigger.Fire("/m0", 123)
	trigger.Fire("/m1", 123)
	AssertEqual(t, counter, 2)
}

func TestTriggerReturnValue(t *testing.T) {
	trigger, _ := NewTrigger()
	trigger.On("/m1", func() {
	})
	trigger.On("/m1", func() int {
		return 2
	})
	trigger.On("/m1", func() (string, int) {
		return "3", 4
	})
	result := trigger.Fire("/m1")
	AssertEqual(t, result, []interface{}{nil, 2, "3"})
}
