package container

import (
	. "github.com/fishedee/assert"
	"math/rand"
	"testing"
)

func TestTrieAllMatch(t *testing.T) {
	trieTree := NewTrieTree()

	testFoundData := []string{"a", "ab", "ba", "1", "2", "3c_4"}
	testNotFoundData := []string{"c", "abc", "ac", "b", "bac", "_4", "c_", "a ", " a"}

	for _, data := range testFoundData {
		trieTree.Set(data, data+"_value")
	}

	trie := trieTree.ToTrieArray()
	for _, data := range testFoundData {
		result := trie.ExactMatch(data)
		AssertEqual(t, result, data+"_value")
		result2 := trieTree.Get(data)
		AssertEqual(t, result2, data+"_value")
	}
	for _, data := range testNotFoundData {
		result := trie.ExactMatch(data)
		AssertEqual(t, result, nil)
		result2 := trieTree.Get(data)
		AssertEqual(t, result2, nil)
	}
}

func TestTriePrefixMatch(t *testing.T) {
	trieTree := NewTrieTree()

	testData := []string{"ab", "/ab", "cde", "abeg", "ac"}
	testFindData := map[string][]TrieMatch{
		"ab": []TrieMatch{
			{"ab", "ab_value"},
		},
		"abc": []TrieMatch{
			{"ab", "ab_value"},
		},
		"/abc": []TrieMatch{
			{"/ab", "/ab_value"},
		},
		"cdef/g": []TrieMatch{
			{"cde", "cde_value"},
		},
		"abegh": []TrieMatch{
			{"ab", "ab_value"},
			{"abeg", "abeg_value"},
		},
		"ac": []TrieMatch{
			{"ac", "ac_value"},
		},
		"acm": []TrieMatch{
			{"ac", "ac_value"},
		},
		"a":   []TrieMatch{},
		"/a":  []TrieMatch{},
		"cck": []TrieMatch{},
		"cd":  []TrieMatch{},
	}

	for _, data := range testData {
		trieTree.Set(data, data+"_value")
	}

	trie := trieTree.ToTrieArray()
	for key, value := range testFindData {
		result := trie.PrefixMatch(key)
		AssertEqual(t, result, value)

		result2 := trie.LongestPrefixMatch(key)
		var value2 interface{}
		if len(value) == 0 {
			value2 = nil
		} else {
			value2 = value[len(value)-1].value
		}
		AssertEqual(t, result2, value2)

	}
}

func TestTrieFull(t *testing.T) {
	trieTree := NewTrieTree()

	testData := []string{"", "abc", "bc", "bcd"}
	testFindData := map[string][]TrieMatch{
		"": []TrieMatch{
			{"", "_value"},
		},
		"abc": []TrieMatch{
			{"", "_value"},
			{"abc", "abc_value"},
		},
		"abcg": []TrieMatch{
			{"", "_value"},
			{"abc", "abc_value"},
		},
		"b": []TrieMatch{
			{"", "_value"},
		},
		"bc": []TrieMatch{
			{"", "_value"},
			{"bc", "bc_value"},
		},
		"bcg": []TrieMatch{
			{"", "_value"},
			{"bc", "bc_value"},
		},
		"bcd": []TrieMatch{
			{"", "_value"},
			{"bc", "bc_value"},
			{"bcd", "bcd_value"},
		},
		"c": []TrieMatch{
			{"", "_value"},
		},
	}

	for _, data := range testData {
		trieTree.Set(data, data+"_value")
	}

	trie := trieTree.ToTrieArray()
	for key, value := range testFindData {
		result := trie.PrefixMatch(key)
		AssertEqual(t, result, value, key)

		result2 := trie.LongestPrefixMatch(key)
		var value2 interface{}
		if len(value) == 0 {
			value2 = nil
		} else {
			value2 = value[len(value)-1].value
		}
		AssertEqual(t, result2, value2, key)
	}
}

func TestTrieWalk(t *testing.T) {
	trieTree := NewTrieTree()

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
		trieTree.Set(data, data+"_value")
	}

	index := 0
	trieTree.Walk(func(key string, value interface{}, parentKey string, parentValue interface{}) {
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

func BenchmarkTrieSpeed(b *testing.B) {
	insertData := getData(1000, 20)
	findData := getData(b.N, 20)
	trieTree := NewTrieTree()
	for _, singleData := range insertData {
		trieTree.Set(singleData, true)
	}
	trie := trieTree.ToTrieArray()

	b.ResetTimer()
	b.StartTimer()
	for _, singleData := range findData {
		trie.ExactMatch(singleData)
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
