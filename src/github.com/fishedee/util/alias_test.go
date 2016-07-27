package util

import (
	"math"
	"reflect"
	"testing"
)

func assertAliasEqual(t *testing.T, left interface{}, right interface{}) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("assert fail: %+v != %+v", left, right)
	}
}

func TestAlias(t *testing.T) {
	testCase := [][]float64{
		[]float64{0.25, 0.2, 0.1, 0.05, 0.4},
		[]float64{0.25, 0.2, 0.1, 0.05, 0.1, 0.1, 0.1, 0.1},
		[]float64{0.2, 0.21, 0.11, 0.06, 0.42},
		[]float64{0.32, 0.12, 0.11, 0.05, 0.4},
	}
	probResult := [][]float64{
		[]float64{1, 0.75, 0.5, 0.25, 0.75},
		[]float64{1, 0.8, 0.8, 0.4, 0.8, 0.8, 0.8, 0.8},
		[]float64{1, 1, 0.55, 0.3, 0.95},
		[]float64{1, 0.6, 0.55, 0.25, 0.8},
	}
	aliasResult := [][]int{
		[]int{-1, 0, 4, 4, 1},
		[]int{-1, 0, 0, 0, 1, 1, 1, 1},
		[]int{-1, -1, 4, 4, 1},
		[]int{-1, 0, 4, 4, 0},
	}

	for index, singleTestCase := range testCase {
		expected := singleTestCase[0]

		alias := NewAliasMethod(singleTestCase)
		assertAliasEqual(t, alias.prob, probResult[index])
		assertAliasEqual(t, alias.alias, aliasResult[index])

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
			assertAliasEqual(t, expected, real)
		}
	}
}
