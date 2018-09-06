package web

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/assert"
	"testing"
	"time"
)

func newTimerForTest(t *testing.T) Timer {
	log, err := NewLog(LogConfig{
		Driver: "console",
	})
	AssertEqual(t, err, nil, 0)
	timer, err := NewTimer(log)
	AssertEqual(t, err, nil, 0)
	return timer
}

func TestTimerBasic(t *testing.T) {
	testCase := []struct {
		Cron    string
		Counter int
	}{
		{"* * * * * *", 4},
		{"*/2 * * * * *", 2},
		{"*/3 * * * * *", 1},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		timer := newTimerForTest(t)
		counter := 0
		timer.Cron(singleTestCase.Cron, func() {
			counter++
		})
		go timer.Run()
		time.Sleep(time.Second * 4)
		timer.Close()
		AssertEqual(t, counter, singleTestCase.Counter, singleTestCaseIndex)
	}
}

func TestTimerException(t *testing.T) {
	timer := newTimerForTest(t)
	counter := 0
	timer.Cron("* * * * * *", func() {
		counter++
		panic("fuck")
	})
	go timer.Run()
	time.Sleep(time.Second * 2)
	timer.Close()
	AssertEqual(t, counter, 2)
}
