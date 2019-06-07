package language_test

import (
	. "github.com/fishedee/language"
	"reflect"
	"testing"
)

func TestArrayKeyAndValue(t *testing.T) {
	testCase := []struct {
		origin interface{}
		key    interface{}
		value  interface{}
	}{
		{map[int]int{}, []int{}, []int{}},
		{map[int]string{}, []int{}, []string{}},
		{map[string]int{}, []string{}, []int{}},
		{map[int]int{1: 3}, []int{1}, []int{3}},
		{map[int]int{1: 3, 4: 5}, []int{1, 4}, []int{3, 5}},
		{map[int]int{1: 3, 4: 5, 7: 8}, []int{1, 4, 7}, []int{3, 5, 8}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		key, value := ArrayKeyAndValue(singleTestCase.origin)
		AssertEqual(t, ArraySort(key), ArraySort(singleTestCase.key), singleTestCaseIndex)
		AssertEqual(t, ArraySort(value), ArraySort(singleTestCase.value), singleTestCaseIndex)
	}
}

func TestArrayReverse(t *testing.T) {
	testCase := []struct {
		origin interface{}
		target interface{}
	}{
		{[]int{}, []int{}},
		{[]int{3}, []int{3}},
		{[]int{3, 8}, []int{8, 3}},
		{[]int{3, 8, 4}, []int{4, 8, 3}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		target := ArrayReverse(singleTestCase.origin)
		AssertEqual(t, target, singleTestCase.target, singleTestCaseIndex)
	}
}

func TestArrayIn(t *testing.T) {
	testCase := []struct {
		origin   interface{}
		noOrigin interface{}
	}{
		{[]int{}, []int{-1, 2, 4, 6}},
		{[]int{1}, []int{-1, 2, 4, 6}},
		{[]int{3, 1, 5}, []int{-1, 2, 4, 6}},
		{[]string{}, []string{"-1", "2", "4", "6"}},
		{[]string{"1"}, []string{"-1", "2", "4", "6"}},
		{[]string{"3", "1", "5"}, []string{"-1", "2", "4", "6"}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		originArray := reflect.ValueOf(singleTestCase.origin)
		for i := 0; i != originArray.Len(); i++ {
			index := ArrayIn(originArray.Interface(), originArray.Index(i).Interface())
			AssertEqual(t, index, i, singleTestCaseIndex)
		}

		noOriginArray := reflect.ValueOf(singleTestCase.noOrigin)
		for i := 0; i != noOriginArray.Len(); i++ {
			index := ArrayIn(originArray.Interface(), noOriginArray.Index(i))
			AssertEqual(t, index, -1, singleTestCaseIndex)
		}
	}
}

func TestArrayUnique(t *testing.T) {
	testCase := []struct {
		origin interface{}
		target interface{}
	}{
		{[]int{}, []int{}},
		{[]string{}, []string{}},
		{[]int{1, 3, 4, 2, 1, 5}, []int{1, 3, 4, 2, 5}},
		{[]string{"4", "1", "7", "4", "6", "7"}, []string{"4", "1", "7", "6"}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		result := ArrayUnique(singleTestCase.origin)
		AssertEqual(t, result, singleTestCase.target, singleTestCaseIndex)
	}
}

func TestArrayDiff(t *testing.T) {
	testCase := []struct {
		origin  interface{}
		origin2 []interface{}
		target  interface{}
	}{
		{[]int{}, []interface{}{[]int{}}, []int{}},
		{[]string{}, []interface{}{[]string{}}, []string{}},
		{[]int{1, 7, 8, 3, 4, 8}, []interface{}{[]int{1, 4, 6, 4}}, []int{7, 8, 3}},
		{[]int{1, 7, 8, 3, 4, 8}, []interface{}{[]int{1}, []int{4, 6, 4}}, []int{7, 8, 3}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		result := ArrayDiff(singleTestCase.origin, singleTestCase.origin2[0], singleTestCase.origin2[1:]...)
		AssertEqual(t, result, singleTestCase.target, singleTestCaseIndex)
	}
}

func TestArrayIntersect(t *testing.T) {
	testCase := []struct {
		origin  interface{}
		origin2 []interface{}
		target  interface{}
	}{
		{[]int{}, []interface{}{[]int{}}, []int{}},
		{[]string{}, []interface{}{[]string{}}, []string{}},
		{[]int{1, 7, 8, 3, 4, 8}, []interface{}{[]int{1, 4, 6, 4}}, []int{1, 4}},
		{[]int{1, 7, 8, 3, 4, 8}, []interface{}{[]int{1}, []int{4, 6, 4}}, []int{1, 4}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		result := ArrayIntersect(singleTestCase.origin, singleTestCase.origin2[0], singleTestCase.origin2[1:]...)
		AssertEqual(t, result, singleTestCase.target, singleTestCaseIndex)
	}
}

func TestArrayMerge(t *testing.T) {
	testCase := []struct {
		origin  interface{}
		origin2 []interface{}
		target  interface{}
	}{
		{[]int{}, []interface{}{[]int{}}, []int{}},
		{[]string{}, []interface{}{[]string{}}, []string{}},
		{[]int{1, 7, 8, 3, 4, 8}, []interface{}{[]int{1, 4, 6, 4}}, []int{1, 7, 8, 3, 4, 6}},
		{[]int{1, 7, 8, 3, 4, 8}, []interface{}{[]int{1}, []int{4, 6, 4}}, []int{1, 7, 8, 3, 4, 6}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		result := ArrayMerge(singleTestCase.origin, singleTestCase.origin2[0], singleTestCase.origin2[1:]...)
		AssertEqual(t, result, singleTestCase.target, singleTestCaseIndex)
	}
}

func TestArraySort(t *testing.T) {
	testCase := []struct {
		origin interface{}
		target interface{}
	}{
		{[]int{}, []int{}},
		{[]int{1}, []int{1}},
		{[]int{2, 1}, []int{1, 2}},
		{[]int{1, 2}, []int{1, 2}},
		{[]int{3, 7, 1}, []int{1, 3, 7}},
		{[]string{}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"b", "a"}, []string{"a", "b"}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]string{"a", "c", "b"}, []string{"a", "b", "c"}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		result := ArraySort(singleTestCase.origin)
		AssertEqual(t, result, singleTestCase.target, singleTestCaseIndex)
	}
}

func TestArrayShuffle(t *testing.T) {
	testCase := []struct {
		origin interface{}
	}{
		{[]int{}},
		{[]int{1}},
		{[]int{2, 1}},
		{[]int{1, 2, 3}},
		{[]string{}},
		{[]string{"a"}},
		{[]string{"b", "a"}},
		{[]string{"a", "b"}},
		{[]string{"a", "c", "b"}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		result := ArrayShuffle(singleTestCase.origin)
		leftResult := ArraySort(singleTestCase.origin)
		rightResult := ArraySort(result)
		AssertEqual(t, leftResult, rightResult, singleTestCaseIndex)
	}
}

func TestArraySlice(t *testing.T) {
	testCase := []struct {
		origin interface{}
		begin  int
		end    []int
		target interface{}
	}{
		{[]int{}, 0, nil, []int{}},
		{[]int{}, 0, []int{0}, []int{}},
		{[]int{}, 0, []int{1}, []int{}},
		{[]int{}, -1, []int{0}, []int{}},
		{[]int{}, 0, []int{1}, []int{}},
		{[]int{}, -1, []int{1}, []int{}},
		{[]int{1, 2, 3}, 0, nil, []int{1, 2, 3}},
		{[]int{1, 2, 3}, 0, []int{0}, []int{}},
		{[]int{1, 2, 3}, 0, []int{1}, []int{1}},
		{[]int{1, 2, 3}, 1, []int{3}, []int{2, 3}},
		{[]int{1, 2, 3}, 1, nil, []int{2, 3}},
		{[]int{1, 2, 3}, -5, []int{-10}, []int{}},
		{[]int{1, 2, 3}, -5, []int{0}, []int{}},
		{[]int{1, 2, 3}, -5, []int{2}, []int{1, 2}},
		{[]int{1, 2, 3}, -5, []int{2}, []int{1, 2}},
		{[]int{1, 2, 3}, 2, []int{4}, []int{3}},
		{[]int{1, 2, 3}, 3, []int{4}, []int{}},
		{[]int{1, 2, 3}, 10, []int{4}, []int{}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		result := ArraySlice(singleTestCase.origin, singleTestCase.begin, singleTestCase.end...)
		AssertEqual(t, singleTestCase.target, result, singleTestCaseIndex)
	}
}
