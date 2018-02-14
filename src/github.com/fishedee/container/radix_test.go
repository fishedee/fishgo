package container

import (
	. "github.com/fishedee/assert"
	"math/rand"
	"testing"
)

func TestRadixAllMatch(t *testing.T) {
	radixTree := NewRadixTree()

	testFoundData := []string{"a", "ab", "ba", "1", "2", "3c_4"}
	testNotFoundData := []string{"c", "abc", "ac", "b", "bac", "_4", "c_", "a ", " a"}

	for _, data := range testFoundData {
		radixTree.Set(data, data+"_value")
	}

	radix := radixTree.ToRadixArray()
	for _, data := range testFoundData {
		result := radix.ExactMatch(data)
		AssertEqual(t, result, data+"_value")
		result2 := radixTree.Get(data)
		AssertEqual(t, result2, data+"_value")
	}
	for _, data := range testNotFoundData {
		result := radix.ExactMatch(data)
		AssertEqual(t, result, nil)
		result2 := radixTree.Get(data)
		AssertEqual(t, result2, nil)
	}
}

func TestRadixPrefixMatch(t *testing.T) {
	radixTree := NewRadixTree()

	testData := []string{"ab", "/ab", "cde", "abeg", "ac"}
	testFindData := map[string][]RadixMatch{
		"ab": []RadixMatch{
			{"ab", "ab_value"},
		},
		"abc": []RadixMatch{
			{"ab", "ab_value"},
		},
		"/abc": []RadixMatch{
			{"/ab", "/ab_value"},
		},
		"cdef/g": []RadixMatch{
			{"cde", "cde_value"},
		},
		"abegh": []RadixMatch{
			{"ab", "ab_value"},
			{"abeg", "abeg_value"},
		},
		"ac": []RadixMatch{
			{"ac", "ac_value"},
		},
		"acm": []RadixMatch{
			{"ac", "ac_value"},
		},
		"a":   []RadixMatch{},
		"/a":  []RadixMatch{},
		"cck": []RadixMatch{},
		"cd":  []RadixMatch{},
	}

	for _, data := range testData {
		radixTree.Set(data, data+"_value")
	}

	radix := radixTree.ToRadixArray()
	for key, value := range testFindData {
		result := radix.PrefixMatch(key)
		AssertEqual(t, result, value)

		result2 := radix.LongestPrefixMatch(key)
		var value2 interface{}
		if len(value) == 0 {
			value2 = nil
		} else {
			value2 = value[len(value)-1].value
		}
		AssertEqual(t, result2, value2)

	}
}

func TestRadixFull(t *testing.T) {
	radixTree := NewRadixTree()

	testData := []string{"", "abc", "bc", "bcd"}
	testFindData := map[string][]RadixMatch{
		"": []RadixMatch{
			{"", "_value"},
		},
		"abc": []RadixMatch{
			{"", "_value"},
			{"abc", "abc_value"},
		},
		"abcg": []RadixMatch{
			{"", "_value"},
			{"abc", "abc_value"},
		},
		"b": []RadixMatch{
			{"", "_value"},
		},
		"bc": []RadixMatch{
			{"", "_value"},
			{"bc", "bc_value"},
		},
		"bcg": []RadixMatch{
			{"", "_value"},
			{"bc", "bc_value"},
		},
		"bcd": []RadixMatch{
			{"", "_value"},
			{"bc", "bc_value"},
			{"bcd", "bcd_value"},
		},
		"c": []RadixMatch{
			{"", "_value"},
		},
	}

	for _, data := range testData {
		radixTree.Set(data, data+"_value")
	}

	radix := radixTree.ToRadixArray()
	for key, value := range testFindData {
		result := radix.PrefixMatch(key)
		AssertEqual(t, result, value, key)

		result2 := radix.LongestPrefixMatch(key)
		var value2 interface{}
		if len(value) == 0 {
			value2 = nil
		} else {
			value2 = value[len(value)-1].value
		}
		AssertEqual(t, result2, value2, key)
	}
}

func TestRadixWalk(t *testing.T) {
	radixTree := NewRadixTree()

	testData := []string{"", "abc", "bc", "bcd"}
	walkData := []struct {
		key         string
		value       interface{}
		parentKey   string
		parentValue interface{}
	}{
		{"", "_value", "", nil},
		{"a", nil, "", "_value"},
		{"b", nil, "", "_value"},
		{"ab", nil, "a", nil},
		{"bc", "bc_value", "b", nil},
		{"abc", "abc_value", "ab", nil},
		{"bcd", "bcd_value", "bc", "bc_value"},
	}

	for _, data := range testData {
		radixTree.Set(data, data+"_value")
	}

	index := 0
	radixTree.Walk(func(key string, value interface{}, parentKey string, parentValue interface{}) {
		result := walkData[index]
		index++
		AssertEqual(t, result.key, key)
		AssertEqual(t, result.value, value)
		AssertEqual(t, result.parentKey, parentKey)
		AssertEqual(t, result.parentValue, parentValue)
	})
	AssertEqual(t, index, len(walkData))
}

func getSingleData(count int) string {
	var randStr = []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	result := make([]byte, count)
	rand.Read(result)
	for singleIndex, singleByte := range result {
		result[singleIndex] = randStr[int(singleByte)%len(randStr)]
	}
	return string(result)
}

func getData(count int, size int) []string {
	result := []string{}
	for i := 0; i != count; i++ {
		result = append(result, getSingleData(size))
	}
	return result
}

func BenchmarkRadixSpeed(b *testing.B) {
	insertData := getData(1000, 20)
	findData := getData(b.N, 20)
	radixTree := NewRadixTree()
	for _, singleData := range insertData {
		radixTree.Set(singleData, true)
	}
	radix := radixTree.ToRadixArray()

	b.ResetTimer()
	b.StartTimer()
	for _, singleData := range findData {
		radix.ExactMatch(singleData)
	}
	b.StopTimer()
}

func BenchmarkMapSpeed(b *testing.B) {
	insertData := getData(1000, 20)
	findData := getData(b.N, 20)
	mapper := map[string]bool{}
	for _, singleData := range insertData {
		mapper[singleData] = true
	}

	b.ResetTimer()
	b.StartTimer()
	for _, singleData := range findData {
		_, _ = mapper[singleData]
	}
	b.StopTimer()
}
