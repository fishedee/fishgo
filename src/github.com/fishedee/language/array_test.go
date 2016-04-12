package language

import (
	"fmt"
	"reflect"
	"testing"
)

func AssertArrayEqual(t *testing.T, left interface{}, right interface{}) {
	isEqual := reflect.DeepEqual(left, right)
	if isEqual == false {
		t.Error(fmt.Sprintf("%#v != %#v", left, right))
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

	for _, singleTestCase := range testCase {
		result := ArrayUnique(singleTestCase.origin)
		AssertArrayEqual(t, result, singleTestCase.target)
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

	for _, singleTestCase := range testCase {
		result := ArrayDiff(singleTestCase.origin, singleTestCase.origin2[0], singleTestCase.origin2[1:]...)
		AssertArrayEqual(t, result, singleTestCase.target)
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

	for _, singleTestCase := range testCase {
		result := ArrayIntersect(singleTestCase.origin, singleTestCase.origin2[0], singleTestCase.origin2[1:]...)
		AssertArrayEqual(t, result, singleTestCase.target)
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

	for _, singleTestCase := range testCase {
		result := ArrayMerge(singleTestCase.origin, singleTestCase.origin2[0], singleTestCase.origin2[1:]...)
		AssertArrayEqual(t, result, singleTestCase.target)
	}
}
