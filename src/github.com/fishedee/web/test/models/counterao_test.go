package test

import (
	"testing"
)

type CounterAoTest struct {
	BaseTest
	CounterAo CounterAoModel
}

func (this *CounterAoTest) TestBasic() {
	total := 1000000

	testCase := []func(){
		this.CounterAo.Incr,
		this.CounterAo.IncrAtomic,
	}
	for singleTestCaseIndex, singleTestCase := range testCase {
		//普通累加
		this.CounterAo.Reset()
		for i := 0; i != total; i++ {
			singleTestCase()
		}
		this.AssertEqual(this.CounterAo.Get(), total, singleTestCaseIndex)

		//并行累加
		this.CounterAo.Reset()
		this.Concurrent(total, 10, func() {
			singleTestCase()
		})
		this.AssertEqual(this.CounterAo.Get(), total, singleTestCaseIndex)
	}
}

func TestCounterAo(t *testing.T) {
	InitTest(t, &CounterAoTest{})
}
