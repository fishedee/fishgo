package web

import (
	. "github.com/fishedee/web"
	"reflect"
	"testing"
	"time"
)

type timerModel struct {
}

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
		Handler func(timer Timer, handler interface{})
	}{
		{func(timer Timer, handler interface{}) {
			timer.Cron("* * * * * *", handler)
		}},
		{func(timer Timer, handler interface{}) {
			timer.Interval(time.Millisecond, handler)
		}},
		{func(timer Timer, handler interface{}) {
			timer.Tick(time.Millisecond, handler)
		}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		var result int
		timer := newTimerForTest(t)
		inputEvent := make(chan bool, 10)
		singleTestCase.Handler(timer, func(this *timerModel) {
			inputEvent <- true
			time.Sleep(time.Second)
			result = 123
		})
		<-inputEvent
		timer.Close()
		assertTimerEqual(t, result, 123, singleTestCaseIndex)
	}
}
