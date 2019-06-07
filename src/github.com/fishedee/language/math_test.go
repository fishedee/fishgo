package language_test

import (
	. "github.com/fishedee/language"
	"testing"
)

func TestAbsInt(t *testing.T) {
	testCase := []struct {
		origin int
		target int
	}{
		{2, 2},
		{1, 1},
		{0, 0},
		{-1, 1},
		{-2, 2},
	}

	for singleTestKey, singleTestCase := range testCase {
		AssertEqual(t, AbsInt(singleTestCase.origin), singleTestCase.target, singleTestKey)
	}
}

func TestEqualsDecimal(t *testing.T) {
	testCase := []struct {
		a float64
		b float64
	}{
		{2, 2},
		{1.1, 1.1},
		{-0.1, -0.1},
		{0, 0},
		{0.1, 0.1},
		{-1.1, -1.1},
		{-2, -2},
	}

	for singleTestKey, singleTestCase := range testCase {
		AssertEqual(t, EqualsDecimal(singleTestCase.a, singleTestCase.b), true, singleTestKey)
	}
}

func TestCmpDecimal(t *testing.T) {
	testCase := []struct {
		a      float64
		b      float64
		target int
	}{
		{2, 3, -1},
		{2, 2, 0},
		{2, 0.1, 1},
		{2, 0, 1},
		{2, -0.1, 1},
		{2, -1.1, 1},

		{1.1, 3, -1},
		{1.1, 1.1, 0},
		{1.1, 0, 1},
		{1.1, 0.1, 1},
		{1.1, -1.1, 1},

		{0.1, 3, -1},
		{0.1, 0.1, 0},
		{0.1, 0, 1},
		{0.1, -0.1, 1},
		{0, 3, -1},
		{0, 1.1, -1},
		{0, 0.1, -1},
		{0, 0, 0},
		{0, -0.1, 1},
		{0, -1.1, 1},

		{-0.1, 3, -1},
		{-0.1, 0.1, -1},
		{-0.1, 0, -1},
		{-0.1, -0.1, 0},
		{-0.1, -1.1, 1},

		{-1.1, 3, -1},
		{-1.1, 1.1, -1},
		{-1.1, 0, -1},
		{-1.1, -1.1, 0},
		{-1.1, -2.1, 1},

		{-2, 3, -1},
		{-2, 1.1, -1},
		{-2, 0, -1},
		{-2, -2, 0},
		{-2, -2.1, 1},
	}

	for singleTestKey, singleTestCase := range testCase {
		AssertEqual(t, CmpDecimal(singleTestCase.a, singleTestCase.b), singleTestCase.target, singleTestKey)
	}
}

func TestAddDecimal(t *testing.T) {
	testCase := []struct {
		a      float64
		b      float64
		target float64
	}{
		{2, 3, 5},
		{2, 2, 4},
		{2, 0, 2},
		{2, 0.1, 2.1},
		{2, -1.1, 0.9},

		{1.1, 3, 4.1},
		{1.1, 1.1, 2.2},
		{1.1, 0, 1.1},
		{1.1, 0.1, 1.2},
		{1.1, -1.1, 0},

		{0.1, 3, 3.1},
		{0.1, 0.1, 0.2},
		{0.1, 0, 0.1},
		{0.1, -0.1, 0},
		{0.1, -0.2, -0.1},

		{0, 3, 3},
		{0, 1.1, 1.1},
		{0, 0.1, 0.1},
		{0, 0, 0},
		{0, -0.1, -0.1},
		{0, -1.1, -1.1},

		{-0.1, 3, 2.9},
		{-0.1, 0.1, 0},
		{-0.1, 0, -0.1},
		{-0.1, -0.1, -0.2},
		{-0.1, -1.1, -1.2},

		{-1.1, 3, 1.9},
		{-1.1, 1.1, 0},
		{-1.1, 0, -1.1},
		{-1.1, -1.1, -2.2},
		{-1.1, -2.1, -3.2},

		{-2, 3, 1},
		{-2, 1.1, -0.9},
		{-2, 0, -2},
		{-2, -2, -4},
		{-2, -2.1, -4.1},
	}

	for singleTestKey, singleTestCase := range testCase {
		AssertEqual(t, AddDecimal(singleTestCase.a, singleTestCase.b), singleTestCase.target, singleTestKey)
	}
}

func TestSubDecimal(t *testing.T) {
	testCase := []struct {
		a      float64
		b      float64
		target float64
	}{
		{2, 3, -1},
		{2, 2, 0},
		{2, 0, 2},
		{2, 0.1, 1.9},
		{2, -1.1, 3.1},

		{1.1, 3, -1.9},
		{1.1, 1.1, 0},
		{1.1, 0, 1.1},
		{1.1, 0.1, 1},
		{1.1, -1.1, 2.2},

		{0.1, 3, -2.9},
		{0.1, 0.1, 0},
		{0.1, 0, 0.1},
		{0.1, -0.1, 0.2},
		{0.1, -0.2, 0.3},

		{0, 3, -3},
		{0, 1.1, -1.1},
		{0, 0.1, -0.1},
		{0, 0, 0},
		{0, -0.1, 0.1},
		{0, -1.1, 1.1},

		{-0.1, 3, -3.1},
		{-0.1, 0.1, -0.2},
		{-0.1, 0, -0.1},
		{-0.1, -0.1, 0},
		{-0.1, -1.1, 1},

		{-1.1, 3, -4.1},
		{-1.1, 1.1, -2.2},
		{-1.1, 0, -1.1},
		{-1.1, -1.1, 0},
		{-1.1, -2.1, 1},

		{-2, 3, -5},
		{-2, 1.1, -3.1},
		{-2, 0, -2},
		{-2, -2, 0},
		{-2, -2.1, 0.1},
	}

	for singleTestKey, singleTestCase := range testCase {
		AssertEqual(t, SubDecimal(singleTestCase.a, singleTestCase.b), singleTestCase.target, singleTestKey)
	}
}

func TestMulDecimal(t *testing.T) {
	testCase := []struct {
		a      float64
		b      float64
		target float64
	}{
		{2, 3, 6},
		{2, 2, 4},
		{2, 0, 0},
		{2, 0.1, 0.2},
		{2, -1.1, -2.2},

		{1.1, 3, 3.3},
		{1.1, 1.1, 1.21},
		{1.1, 0, 0},
		{1.1, 0.1, 0.11},
		{1.1, -1.1, -1.21},

		{0.1, 3, 0.3},
		{0.1, 0.1, 0.01},
		{0.1, 0, 0},
		{0.1, -0.1, -0.01},
		{0.1, -0.2, -0.02},

		{0, 3, 0},
		{0, 1.1, 0},
		{0, 0.1, 0},
		{0, 0, 0},
		{0, -0.1, 0},
		{0, -1.1, 0},

		{-0.1, 3, -0.3},
		{-0.1, 0.1, -0.01},
		{-0.1, 0, 0},
		{-0.1, -0.1, 0.01},
		{-0.1, -1.1, 0.11},

		{-1.1, 3, -3.3},
		{-1.1, 1.1, -1.21},
		{-1.1, 0, 0},
		{-1.1, -1.1, 1.21},
		{-1.1, -2.1, 2.31},

		{-2, 3, -6},
		{-2, 1.1, -2.2},
		{-2, 0, 0},
		{-2, -2, 4},
		{-2, -2.1, 4.2},
	}

	for singleTestKey, singleTestCase := range testCase {
		AssertEqual(t, MulDecimal(singleTestCase.a, singleTestCase.b), singleTestCase.target, singleTestKey)
	}
}

func TestDivDecimal_Normal(t *testing.T) {
	testCase := []struct {
		a      float64
		b      float64
		target float64
	}{
		{7, 3, 2.3333333333333335},
		{7, -3, -2.3333333333333335},
		{3, 1, 3},
		{2, 2, 1},
		{2, 0.1, 20},
		{2, -0.1, -20},
		{0.1, 3, 0.0333333333333333},
		{0.1, 0.01, 10},
		{0.1, -0.01, -10},
		{-0.1, -0.01, 10},
		{-0.1, 0.01, -10},
		{-1.1, -0.01, 110},
		{-1.1, 0.01, -110},
		{-1.1, -0.01, 110},
		{-1.1, 0.01, -110},
		{-1.1, 3, -0.3666666666666667},
		{-7, 3, -2.3333333333333335},
		{-7, -3, 2.3333333333333335},
	}

	for singleTestKey, singleTestCase := range testCase {
		AssertEqual(t, DivDecimal(singleTestCase.a, singleTestCase.b), singleTestCase.target, singleTestKey)
	}
}

func TestDivDecimal_Error(t *testing.T) {
	testCase := []struct {
		a float64
		b float64
	}{
		{1, 0},
		{0.1, 0},
		{0, 0},
		{-0.1, 0},
	}

	for singleTestKey, singleTestCase := range testCase {
		AssertError(t, "decimal division by 0", func() {
			DivDecimal(singleTestCase.a, singleTestCase.b)
		}, singleTestKey)
	}
}

func TestModDecimal(t *testing.T) {
	testCase := []struct {
		a      float64
		b      float64
		target float64
	}{
		{7, 3, 1},
		{7, -3, 1},
		{3, 1, 0},
		{2, 2, 0},
		{2, 0.1, 0},
		{2, -0.1, 0},
		{1.1, 3, 1.1},
		{0.1, 3, 0.1},
		{0.1, 0.01, 0},
		{0.1, -0.01, 0},
		{-0.1, -0.01, 0},
		{-0.1, 0.01, 0},
		{-1.1, -0.01, 0},
		{-1.1, 0.01, 0},
		{-1.1, -0.01, 0},
		{-1.1, 0.01, 0},
		{-1.1, 3, -1.1},
		{-7, 3, -1},
		{-7, -3, -1},
	}

	for singleTestKey, singleTestCase := range testCase {
		AssertEqual(t, ModDecimal(singleTestCase.a, singleTestCase.b), singleTestCase.target, singleTestKey)
	}
}

func TestAbsDecimal(t *testing.T) {
	testCase := []struct {
		origin float64
		target float64
	}{
		{2, 2},
		{1.1, 1.1},
		{0.1, 0.1},
		{0, 0},
		{-0.1, 0.1},
		{-1.1, 1.1},
		{-7, 7},
	}

	for singleTestKey, singleTestCase := range testCase {
		AssertEqual(t, AbsDecimal(singleTestCase.origin), singleTestCase.target, singleTestKey)
	}
}
