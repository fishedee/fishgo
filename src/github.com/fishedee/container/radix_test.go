package container

import (
	. "github.com/fishedee/assert"
	"math/rand"
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
	radixFactory := NewRadixFactory([]int{RadixMatchMode.PREFIX})

	testData := []string{"ab", "/ab", "cde", "abeg", "ac"}
	testFoundData := map[string]string{
		"ab":     "ab",
		"abc":    "ab",
		"/abc":   "/ab",
		"cdef/g": "cde",
		"abegh":  "abeg",
		"ac":     "ac",
		"acm":    "ac",
	}
	testNotFoundData := []string{"a", "/a", "cck", "cd"}

	for _, data := range testData {
		radixFactory.Insert(0, data, data+"_value")
	}

	radix := radixFactory.Create()
	for key, value := range testFoundData {
		result := radix.Find(key)
		AssertEqual(t, result, []interface{}{value + "_value"})
	}
	for _, data := range testNotFoundData {
		result := radix.Find(data)
		AssertEqual(t, result, []interface{}{nil})
	}
}

func TestRadixError(t *testing.T) {
	radixFactory := NewRadixFactory([]int{RadixMatchMode.ALL, RadixMatchMode.PREFIX})

	var err error

	err = radixFactory.Insert(-1, "a", "a_value")
	AssertEqual(t, err != nil, true)
	err = radixFactory.Insert(2, "a", "a_value")
	AssertEqual(t, err != nil, true)

	err = radixFactory.Insert(0, "b", "b_value")
	AssertEqual(t, err, nil)
	err = radixFactory.Insert(1, "c", "c_value")
	AssertEqual(t, err, nil)

	err = radixFactory.Insert(0, "b", "b_value")
	AssertEqual(t, err != nil, true)
	err = radixFactory.Insert(1, "c", "c_value")
	AssertEqual(t, err != nil, true)
}

func TestRadixFull(t *testing.T) {
	radixFactory := NewRadixFactory([]int{RadixMatchMode.ALL, RadixMatchMode.PREFIX})

	testData := []struct {
		mode  int
		key   string
		value string
	}{
		{0, "", "index_all"},
		{1, "", "index_prefix"},
		{0, "abc", "abc_all"},
		{0, "bc", "bc_all"},
		{1, "bc", "bc_prefix"},
		{0, "bcd", "bcd_all"},
	}
	testFindData := []struct {
		key   string
		value []interface{}
	}{
		{"", []interface{}{"index_all", "index_prefix"}},
		{"abc", []interface{}{"abc_all", "index_prefix"}},
		{"abcg", []interface{}{nil, "index_prefix"}},
		{"b", []interface{}{nil, "index_prefix"}},
		{"bc", []interface{}{"bc_all", "bc_prefix"}},
		{"bcg", []interface{}{nil, "bc_prefix"}},
		{"bcd", []interface{}{"bcd_all", "bc_prefix"}},
		{"c", []interface{}{nil, "index_prefix"}},
	}

	for _, singleTestData := range testData {
		radixFactory.Insert(singleTestData.mode, singleTestData.key, singleTestData.value)
	}

	radix := radixFactory.Create()
	for _, singleFindData := range testFindData {
		result := radix.Find(singleFindData.key)
		AssertEqual(t, result, singleFindData.value)
	}
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
	radixFactory := NewRadixFactory([]int{RadixMatchMode.ALL})
	for _, singleData := range insertData {
		radixFactory.Insert(0, singleData, true)
	}
	radix := radixFactory.Create()

	b.ResetTimer()
	b.StartTimer()
	for _, singleData := range findData {
		radix.Find(singleData)
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
