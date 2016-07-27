package util

import (
	"testing"
	"time"
)

//Author:Edward
func TestStack(t *testing.T) {

	testCase := []struct {
		origin []interface{}
	}{
		{
			[]interface{}{-1, 0, 1, 2, 3, 4, 5},
		},
		{
			[]interface{}{"", "a", "哈哈", "123"},
		},
		{
			[]interface{}{-1.1, 0, 1.1},
		},
		{
			[]interface{}{false, true, true, false},
		},
		{
			[]interface{}{time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 0), time.Now().AddDate(0, 0, 1)},
		},
	}

	for _, singletest := range testCase {
		testFunc(t, singletest.origin)
	}

}

func testFunc(t *testing.T, origin []interface{}) {
	list := NewStack()

	//Push And Peak
	for _, v := range origin {
		list.Push(v)
		assertAliasEqual(t, list.Peak(), v)
	}

	getNum := 0
	for e := list.list.Front(); e != nil; e = e.Next() {
		assertAliasEqual(t, e.Value, origin[getNum])
		getNum++
	}

	//Len
	assertAliasEqual(t, list.Len(), len(origin))

	//IsEmpty false
	assertAliasEqual(t, list.IsEmpty(), false)

	//Pop
	popNum := len(origin) - 1
	for {
		data := list.Pop()
		if data == nil {
			break
		}
		assertAliasEqual(t, data, origin[popNum])
		popNum--
	}

	//IsEmpty true
	assertAliasEqual(t, list.IsEmpty(), true)
}
