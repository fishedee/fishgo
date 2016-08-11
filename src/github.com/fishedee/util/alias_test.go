package util

import (
	. "github.com/fishedee/assert"
	"math"
	"testing"
)

func TestAlias(t *testing.T) {
	testCase := [][]float64{
		[]float64{0.25, 0.2, 0.1, 0.05, 0.4},
		[]float64{0.25, 0.2, 0.1, 0.05, 0.1, 0.1, 0.1, 0.1},
		[]float64{0.2, 0.21, 0.11, 0.06, 0.42},
		[]float64{0.32, 0.12, 0.11, 0.05, 0.4},
		[]float64{0.22, 0.44, 0.34},
		[]float64{0, 0, 0.01, 0.33, 0.55, 0, 0.11, 0},
	}
	probResult := [][]float64{
		[]float64{1, 0.75, 0.5, 0.25, 0.75},
		[]float64{1, 0.8, 0.8, 0.4, 0.8, 0.8, 0.8, 0.8},
		[]float64{1, 1, 0.55, 0.3, 0.95},
		[]float64{1, 0.6, 0.55, 0.25, 0.8},
		[]float64{0.66, 1, 0.6801},
		[]float64{0, 0, 0.08, 1, 0.36, 0, 0.88, 0},
	}
	aliasResult := [][]int{
		[]int{-1, 0, 4, 4, 1},
		[]int{-1, 0, 0, 0, 1, 1, 1, 1},
		[]int{-1, -1, 4, 4, 1},
		[]int{-1, 0, 4, 4, 0},
		[]int{2, -1, 1},
		[]int{3, 4, 4, -1, 3, 4, 4, 4},
	}

	for index, singleTestCase := range testCase {
		expected := singleTestCase[0]

		alias := NewAliasMethod(singleTestCase)
		AssertEqual(t, alias.prob, probResult[index])
		AssertEqual(t, alias.alias, aliasResult[index])

		sum := 0.0
		testNum := 10000
		for i := 0; i < testNum; i++ {
			rand := alias.Rand()
			if rand == 0 {
				sum++
			}
		}
		real := sum / float64(testNum)
		if math.Abs(expected-real) >= 0.05 {
			AssertEqual(t, expected, real)
		}
	}

	//异常抛出测试
	testErrorCase := []struct {
		in  []float64
		out string
	}{
		{
			[]float64{0, 0.01, 0.08, 0.2, 0.7, 0, 0},
			"传入概率数组之和不为1～[0 0.01 0.08 0.2 0.7 0 0]",
		},
		{
			[]float64{0, 0.01, 0.1, 0.2, 0.7, 0, 0},
			"传入概率数组之和不为1～[0 0.01 0.1 0.2 0.7 0 0]",
		},
		{
			[]float64{0, 0.17, 0.3, 0.7, 0.11, 0, 0.11},
			"传入概率数组之和不为1～[0 0.17 0.3 0.7 0.11 0 0.11]",
		},
		{
			[]float64{0, 0.17, 0.3, 0.7, -0.17, 0, 0},
			"传入概率不能小于0,你输入了:-0.17",
		},
	}

	for singleTestKey, singleTestCase := range testErrorCase {
		AssertError(t, singleTestCase.out, func() {
			NewAliasMethod(singleTestCase.in)
		}, singleTestKey)
	}

}

func TestInit(t *testing.T) {
	alias := &AliasMethod{}
	num := alias.Rand()
	AssertEqual(t, num, -1)
}
