package web

import (
	. "github.com/fishedee/web"
	"reflect"
	"testing"
	"time"
)

func assertTimerEqual(t *testing.T, left interface{}, right interface{}, index int) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("case :%v ,%+v != %+v", index, left, right)
	}
}

func newTimerForTest(t *testing.T) Timer {
	timer, err := NewTimer()
	assertTimerEqual(t, err, nil, 0)
	return timer
}

func TestTimerClose(t *testing.T) {
	testCase := []struct {
		Handler func(timer Timer, handler func())
	}{
		{func(timer Timer, handler func()) {
			err := timer.Cron("* * * * * *", handler)
			assertTimerEqual(t, err, nil, 0)
		}},
		{func(timer Timer, handler func()) {
			err := timer.Interval(time.Millisecond, handler)
			assertTimerEqual(t, err, nil, 0)
		}},
		{func(timer Timer, handler func()) {
			err := timer.Tick(time.Millisecond, handler)
			assertTimerEqual(t, err, nil, 0)
		}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		var result int
		timer := newTimerForTest(t)
		inputEvent := make(chan bool)
		singleTestCase.Handler(timer, func() {
			inputEvent <- true
			time.Sleep(time.Second)
			result = 123
		})
		<-inputEvent
		timer.Close()
		assertTimerEqual(t, result, 123, singleTestCaseIndex)
	}
}
