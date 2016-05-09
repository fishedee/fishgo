package test

import (
	. "github.com/fishedee/web"
)

type counterAoTest struct {
	Test
	counterAo CounterAoModel
}

func (this *counterAoTest) TestBasic() {
	total := 1000000

	testCase := []func(){
		this.counterAo.Incr, //并行累加会失败
		this.counterAo.IncrAtomic,
	}
	for singleTestCaseIndex, singleTestCase := range testCase {
		//普通累加
		this.counterAo.Reset()
		for i := 0; i != total; i++ {
			singleTestCase()
		}
		this.AssertEqual(this.counterAo.Get(), total, singleTestCaseIndex)

		//并行累加
		this.counterAo.Reset()
		this.Concurrent(total, 4, func() {
			singleTestCase()
		})
		this.AssertEqual(this.counterAo.Get(), total, singleTestCaseIndex)
	}
}
