package number

import (
	. "github.com/fishedee/web"
)

type NumberAoTest struct {
	Test
	NumberAo NumberAoModel
}

func (this *NumberAoTest) BenchmarkBasic() {
	i := 0
	this.Benchmark(1000, 100, func() {
		i++
	}, "testCaseDesc")
}

func (this *NumberAoTest) TestBasic() {
	testCase := []struct {
		origin  int
		origin2 int
		target  int
	}{
		{-1, 0, -1},
		{0, -1, -1},
		{0, 0, 0},
		{1, 0, 1},
		{0, 1, 1},
		{1, 1, 2},
		{2, 3, 5},
		{-1, 3, 2},
		{2, -3, -1},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		target := this.NumberAo.Add(singleTestCase.origin, singleTestCase.origin2)
		this.AssertEqual(target, singleTestCase.target, singleTestCaseIndex)
	}
}

func init() {
	InitTest(&NumberAoTest{})
}
