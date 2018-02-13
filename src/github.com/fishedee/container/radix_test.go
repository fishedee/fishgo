package container

import (
	. "github.com/fishedee/assert"
	"testing"
)

func TestRadixAllMatch(t *testing.T) {
	radixFactory := NewRadixFactory([]int{RadixMatchMode.ALL})

	testFoundData := []string{"a", "ab", "ba", "1", "2", "3c_4"}
	testNotFoundData := []string{"c", "abc", "ac", "b", "bac", "_4", "c_", "a ", " a"}

	for _, data := range testFoundData {
		radixFactory.Insert(0, data, data+"_value")
	}

	radix := radixFactory.Create()
	for _, data := range testFoundData {
		result := radix.Find(data)
		AssertEqual(t, result, []interface{}{data + "_value"})
	}
	for _, data := range testNotFoundData {
		result := radix.Find(data)
		AssertEqual(t, result, []interface{}{nil})
	}
}

func TestRadixPrefixMatch(t *testing.T) {

}
