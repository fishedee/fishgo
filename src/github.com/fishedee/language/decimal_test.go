package language

import (
	"testing"
)

func TestDecimalOperation(t *testing.T) {
	testCase := []struct {
		handler func(Decimal, Decimal) Decimal
		left    float64
		right   float64
		result  float64
	}{
		{Decimal.Add, 1.0000, 2.0000, 3.0000},
		{Decimal.Sub, 3.0000, 2.0000, 1.0000},
		{Decimal.Mul, 1.234, 5.6789, 7.0078},
		{Decimal.Mul, 1.2330, 5.6789, 7.0021},
		{Decimal.Mul, 1.2320, 5.6789, 6.9964},
		{Decimal.Mul, 1567891.2320, 356.6789, 559233719.9494},
		{Decimal.Div, 1.0000, 2.0000, 0.5000},
		{Decimal.Div, 2.0000, 1.0000, 2.0000},
		{Decimal.Div, 1.0000, 7.0000, 0.1429},
		{Decimal.Div, 1.0000, 3.0000, 0.3333},
		{Decimal.Div, 1.0000, 9.0000, 0.1111},
		{Decimal.Div, 0.0600, 78.0000, 0.0008},
		{Decimal.Div, 0.0500, 78.0000, 0.0006},
		{Decimal.Div, 0.0010, 78.0000, 0.0000},
	}
	for _, singleTestCase := range testCase {
		leftDecimal := NewDecimal(singleTestCase.left)
		rightDecimal := NewDecimal(singleTestCase.right)
		resultDecimal := NewDecimal(singleTestCase.result)
		result := singleTestCase.handler(leftDecimal, rightDecimal)
		AssertEqual(t, result, resultDecimal)
	}
}

func TestDecimalRound(t *testing.T) {
	testCase := []struct {
		origin    float64
		precision int
		result    float64
	}{
		{2.0000, 0, 2.0000},
		{2.1234, 0, 2.0000},
		{2.4999, 0, 2.0000},
		{2.5000, 0, 3.0000},
		{2.5001, 0, 3.0000},
		{3.2000, 1, 3.2000},
		{3.2123, 1, 3.2000},
		{3.2499, 1, 3.2000},
		{3.2500, 1, 3.3000},
		{3.2501, 1, 3.3000},
		{4.3200, 2, 4.3200},
		{4.3212, 2, 4.3200},
		{4.3249, 2, 4.3200},
		{4.3250, 2, 4.3300},
		{4.3251, 2, 4.3300},
		{5.4320, 3, 5.4320},
		{5.4321, 3, 5.4320},
		{5.4324, 3, 5.4320},
		{5.4325, 3, 5.4330},
		{5.4325, 3, 5.4330},
		{5.4320, 4, 5.4320},
		{5.4321, 4, 5.4321},
		{5.4324, 4, 5.4324},
		{5.4325, 4, 5.4325},
		{5.4326, 4, 5.4326},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		originDecimal := NewDecimal(singleTestCase.origin)
		resultDecimal := NewDecimal(singleTestCase.result)
		result := originDecimal.Round(singleTestCase.precision)
		AssertEqual(t, result, resultDecimal, singleTestCaseIndex)
	}
}

func TestDecimalFromString(t *testing.T) {
	testCase := []struct {
		origin string
		result float64
	}{
		{"1", 1.0000},
		{"0.1", 0.1000},
		{"1.2345", 1.2345},
		{"1.234549", 1.2345},
		{"1.23455", 1.2346},
		{"1.23456", 1.2346},
	}
	for _, singleTestCase := range testCase {
		resultDecimal := NewDecimal(singleTestCase.result)
		result, err := NewDecimalFromString(singleTestCase.origin)
		AssertEqual(t, err, nil)
		AssertEqual(t, result, resultDecimal)
	}
}
