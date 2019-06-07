package language_test

import (
	. "github.com/fishedee/language"
	"github.com/fishedee/language/testdata"
	"os"
	"testing"
)

func TestQueryColumn(t *testing.T) {
	testCase := testdata.GetQueryColumnTestCase()
	for singleTestCaseIndex, singleTestCase := range testCase {
		AssertEqual(t, singleTestCase.Target, singleTestCase.Handler(), singleTestCaseIndex)
	}
}

func TestQuerySelect(t *testing.T) {
	testCase := testdata.GetQuerySelectTestCase()

	for singleTestCaseIndex, singleTestCase := range testCase {
		AssertEqual(t, singleTestCase.Target, singleTestCase.Handler(), singleTestCaseIndex)
	}
}

func TestQueryWhere(t *testing.T) {
	testCase := testdata.GetQueryWhereTestCase()

	for singleTestCaseIndex, singleTestCase := range testCase {
		AssertEqual(t, singleTestCase.Target, singleTestCase.Handler(), singleTestCaseIndex)
	}
}

func TestQueryReduce(t *testing.T) {
	testCase := testdata.GetQueryReduceTestCase()

	for singleTestCaseIndex, singleTestCase := range testCase {
		AssertEqual(t, singleTestCase.Target, singleTestCase.Handler(), singleTestCaseIndex)
	}

}

func TestQuerySort(t *testing.T) {
	testCase := testdata.GetQuerySortTestCase()

	for singleTestCaseIndex, singleTestCase := range testCase {
		AssertEqual(t, singleTestCase.Target, singleTestCase.Handler(), singleTestCaseIndex)
	}

}

func TestQueryJoin(t *testing.T) {
	testCase := testdata.GetQueryJoinTestCase()

	for singleTestCaseIndex, singleTestCase := range testCase {
		AssertEqual(t, singleTestCase.Target, singleTestCase.Handler(), singleTestCaseIndex)
	}
}

func TestQueryGroup(t *testing.T) {
	testCase := testdata.GetQueryGroupTestCase()

	for singleTestCaseIndex, singleTestCase := range testCase {
		AssertEqual(t, singleTestCase.Target, singleTestCase.Handler(), singleTestCaseIndex)
	}

}

func TestQueryColumnMap(t *testing.T) {
	testCase := testdata.GetQueryColumnMapTestCase()

	for singleTestCaseIndex, singleTestCase := range testCase {
		AssertEqual(t, singleTestCase.Target, singleTestCase.Handler(), singleTestCaseIndex)
	}
}

func TestQueryReverse(t *testing.T) {
	testCase := testdata.GetQueryReverseTestCase()

	for singleTestCaseIndex, singleTestCase := range testCase {
		AssertEqual(t, singleTestCase.Target, singleTestCase.Handler(), singleTestCaseIndex)
	}
}

func TestQueryCombine(t *testing.T) {
	testCase := testdata.GetQueryCombineTestCase()

	for singleTestCaseIndex, singleTestCase := range testCase {
		AssertEqual(t, singleTestCase.Target, singleTestCase.Handler(), singleTestCaseIndex)
	}
}

func TestQueryDistinct(t *testing.T) {
	testCase := testdata.GetQueryDistinctTestCase()

	for singleTestCaseIndex, singleTestCase := range testCase {
		AssertEqual(t, singleTestCase.Target, singleTestCase.Handler(), singleTestCaseIndex)
	}
}

func init() {
	args := os.Args[len(os.Args)-1]
	if args == "macro" {
		QueryReflectWarning(true)
	} else {
		QueryReflectWarning(false)
	}
}
