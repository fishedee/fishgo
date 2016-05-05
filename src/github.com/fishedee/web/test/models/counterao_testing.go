package test

import (
	. "github.com/fishedee/web"
)

type CounterAoTest struct {
	Test
	CounterAo CounterAoModel
}

func (this *CounterAoTest) TestBasic() {
	total := 1000000

	testCase := []func(){
		this.CounterAo.Incr, //并行累加会失败
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
		this.Concurrent(total, 4, func() {
			singleTestCase()
		})
		this.AssertEqual(this.CounterAo.Get(), total, singleTestCaseIndex)
	}
}

func init() {
	InitTest(&CounterAoTest{})
}
