package language

import (
	"fmt"
	"reflect"
	"testing"
)

func AssertArrayEqual(t *testing.T, testCaseIndex int, left interface{}, right interface{}) {
	isEqual := reflect.DeepEqual(left, right)
	if isEqual == false {
		t.Error(fmt.Sprintf("case: %v,%#v != %#v", testCaseIndex, left, right))
	}
}

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
		AssertArrayEqual(t, singleTestCaseIndex, ArraySort(key), ArraySort(singleTestCase.key))
		AssertArrayEqual(t, singleTestCaseIndex, ArraySort(value), ArraySort(singleTestCase.value))
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
		AssertArrayEqual(t, singleTestCaseIndex, target, singleTestCase.target)
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
			AssertArrayEqual(t, singleTestCaseIndex, index, i)
		}

		noOriginArray := reflect.ValueOf(singleTestCase.noOrigin)
		for i := 0; i != noOriginArray.Len(); i++ {
			index := ArrayIn(originArray.Interface(), noOriginArray.Index(i))
			AssertArrayEqual(t, singleTestCaseIndex, index, -1)
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
		AssertArrayEqual(t, singleTestCaseIndex, result, singleTestCase.target)
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
		AssertArrayEqual(t, singleTestCaseIndex, result, singleTestCase.target)
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
		AssertArrayEqual(t, singleTestCaseIndex, result, singleTestCase.target)
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
		AssertArrayEqual(t, singleTestCaseIndex, result, singleTestCase.target)
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
		AssertArrayEqual(t, singleTestCaseIndex, result, singleTestCase.target)
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
		AssertArrayEqual(t, singleTestCaseIndex, leftResult, rightResult)
	}
}
